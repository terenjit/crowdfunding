package usecases

import (
	"context"
	"crowdfunding/config"
	"crowdfunding/modules/users/helpers"
	models "crowdfunding/modules/users/models/domain"
	"crowdfunding/modules/users/repositories/commands"
	"crowdfunding/modules/users/repositories/queries"
	httpError "crowdfunding/pkg/http-error"
	"crowdfunding/pkg/token"
	"crowdfunding/pkg/utils"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userCommandUsecase struct {
	userPostgreQuery   queries.UserPostgre
	userPostgreCommand commands.UserPostgre
	userQueryUsecase   QueryUsecase
}

// NewUserCommandUsecase initiation

func NewUserCommandUsecase(userPostgreQuery queries.UserPostgre, userPostgreCommand commands.UserPostgre, userQueryUsecase QueryUsecase) *userCommandUsecase {
	return &userCommandUsecase{
		userPostgreQuery:   userPostgreQuery,
		userPostgreCommand: userPostgreCommand,
		userQueryUsecase:   userQueryUsecase,
	}
}

func (c userCommandUsecase) Register(ctx context.Context, payload *models.Register) utils.Result {
	var result utils.Result

	name := payload.Name
	password := payload.Password
	parameter := make(map[string]interface{})
	var query string

	query = "name = @name"
	parameter["name"] = name

	if payload.Email != "" {
		query = query + " " + "OR email = @email"
		parameter["email"] = payload.Email
	}

	queryRes := <-c.userPostgreQuery.FindOne(ctx, query, parameter)
	data := queryRes.Data.(models.User)
	if data.Name != "" {
		errObj := httpError.NewConflict()
		errObj.Message = "Account already exist"
		result.Error = errObj
		return result
	}

	UserId := uuid.New().String()

	genPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		errObj := httpError.NewConflict()
		errObj.Message = "Failed to hash password"
		result.Error = errObj
		return result
	}

	var initials string
	nameSplit := strings.Split(payload.Name, " ")
	for i := 0; i < len(nameSplit); i++ {
		if i >= 2 {
			break
		}
		initials += strings.ToUpper(string(nameSplit[i][0]))
	}

	var user = models.User{
		ID:         UserId,
		Name:       payload.Name,
		Username:   payload.Username,
		Occupation: payload.Occupation,
		Email:      payload.Email,
		Role:       "user",
		Password:   string(genPassword),
	}

	result = <-c.userPostgreCommand.InsertOneUser(ctx, &user)
	if result.Error != nil {
		errObj := httpError.NewConflict()
		errObj.Message = "Failed insert user"
		result.Error = errObj
		return result
	}

	claim := token.Claim{
		Username: user.Username,
		UserID:   user.ID,
	}

	jwt := <-token.Generate(ctx, &claim, config.GlobalEnv.AccessTokenExpired)
	if jwt.Error != nil {
		errObj := httpError.NewBadRequest()
		errObj.Message = fmt.Sprintf("%v", queryRes.Error)
		result.Error = errObj
		return result
	}

	dataRegister := models.UserFormatter{
		ID:         user.ID,
		Name:       user.Name,
		Username:   user.Username,
		Email:      user.Email,
		Occupation: user.Occupation,
		Token:      jwt.Data.(string),
	}

	result.Data = dataRegister
	return result
}

func (c userCommandUsecase) Login(ctx context.Context, payload *models.LoginRequest) utils.Result {
	var result utils.Result

	email := payload.Email
	password := payload.Password

	queryRes := <-c.userPostgreQuery.FindOneByEmail(ctx, email)

	findUser := queryRes.Data.(models.User)
	if findUser.ID == "" {
		errObj := httpError.NewNotFound()
		errObj.Message = "Email belum terdaftar"
		result.Error = errObj
		return result
	}

	err := helpers.VerifyPassword(password, findUser.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		errObj := httpError.NewNotFound()
		errObj.Message = "Password salah"
		result.Error = errObj
		return result
	}

	claim := token.Claim{
		Username: findUser.Username,
		UserID:   findUser.ID,
	}

	jwt := <-token.Generate(ctx, &claim, config.GlobalEnv.AccessTokenExpired)
	if jwt.Error != nil {
		errObj := httpError.NewBadRequest()
		errObj.Message = fmt.Sprintf("%v", queryRes.Error)
		result.Error = errObj
		return result
	}

	data := models.UserFormatter{
		ID:         findUser.ID,
		Name:       findUser.Name,
		Username:   findUser.Username,
		Email:      findUser.Email,
		Occupation: findUser.Occupation,
		Token:      jwt.Data.(string),
	}

	result.Data = data
	return result
}

func (c userCommandUsecase) SaveAvatar(ctx context.Context, file multipart.File, header *multipart.FileHeader, ID string) utils.Result {
	var result utils.Result

	param := fmt.Sprintf("id = '%v'", ID)

	res := <-c.userPostgreQuery.FindOneByID(ctx, ID)
	user := res.Data.(models.User)
	if user.ID == "" {
		errObj := httpError.NewNotFound()
		errObj.Message = "Akun tidak ditemukan"
		result.Error = errObj
		return result
	}

	user.Avatar = "images/" + header.Filename
	size := header.Size
	ext := filepath.Ext(user.Avatar)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		errObj := httpError.NewBadRequest()
		errObj.Message = "Format file tidak valid, format yang valid: jpg, png, jpeg"
		result.Error = errObj
		return result
	}

	if size > int64(1024*1024*2) {
		errObj := httpError.NewBadRequest()
		errObj.Message = "size photo tidak boleh lebih dari 2MB"
		result.Error = errObj
		return result
	}

	res = <-c.userPostgreCommand.Update(param, &user)
	if res.Error != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = "gagal update avatar error: "
		result.Error = errObj
		return result
	}

	result.Data = user
	return result
}

func (c userCommandUsecase) Update(ctx context.Context, payload *models.UpdatedUser) utils.Result {
	var result utils.Result
	var query string
	parameter := make(map[string]interface{})
	query = "u.is_deleted = @is_deleted"
	parameter["is_deleted"] = false
	query = "u.id = @id"
	parameter["id"] = payload.ID

	queryPayload := queries.QueryPayload{
		Table:     "users u",
		Select:    "u.*",
		Query:     query,
		Parameter: parameter,
		Output:    models.UpdatedUser{},
	}

	queryRes := <-c.userPostgreQuery.FindOneUser(&queryPayload)
	if queryRes.Error != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "users tidak ditemukan"
		result.Error = errObj
		return result
	}

	var data = models.UpdatedUser{}
	byteStatus, _ := json.Marshal(queryRes.Data)
	json.Unmarshal(byteStatus, &data)

	var document = make(map[string]interface{})
	result.Data = &document

	currentTime := time.Now()
	document["updated_at"] = currentTime
	document["updated_by"] = payload.Opts.UserID
	document["name"] = payload.Name
	document["occupation"] = payload.Occupation
	document["email"] = data.Email
	document["password"] = data.Password
	document["avatar"] = data.Avatar

	commandPayload := commands.CommandPayload{
		Table:     "users u",
		Query:     "u.id = @id AND u.is_deleted = @is_deleted",
		Parameter: parameter,
		Document:  document,
	}

	var update = <-c.userPostgreCommand.UpdatedUser(&commandPayload)
	if update.Error != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "Gagal memperbarui campaign"
		result.Error = errObj
		return result
	}
	return result
}
