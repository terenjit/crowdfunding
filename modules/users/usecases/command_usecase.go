package usecases

import (
	"context"
	"crowdfunding/modules/users/helpers"
	models "crowdfunding/modules/users/models/domain"
	"crowdfunding/modules/users/repositories/commands"
	"crowdfunding/modules/users/repositories/queries"
	"crowdfunding/pkg/utils"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	httpError "crowdfunding/pkg/http-error"

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

	data := models.UserFormatter{
		ID:         findUser.ID,
		Name:       findUser.Name,
		Username:   findUser.Username,
		Email:      findUser.Email,
		Occupation: findUser.Occupation,
		Token:      "tokentokentoken",
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

	user.AvatarFileName = header.Filename
	size := header.Size
	ext := filepath.Ext(user.AvatarFileName)
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
