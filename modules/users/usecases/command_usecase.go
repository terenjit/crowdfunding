package usecases

import (
	"context"
	models "crowdfunding/modules/users/models/domain"
	"crowdfunding/modules/users/repositories/commands"
	"crowdfunding/modules/users/repositories/queries"
	"crowdfunding/pkg/utils"
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
		Occupation: payload.Occupation,
		Email:      payload.Email,
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
