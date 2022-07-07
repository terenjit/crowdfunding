package usecases

import (
	"context"
	models "crowdfunding/modules/users/models/domain"
	"crowdfunding/pkg/utils"
	"mime/multipart"
)

// CommandUsecase interface
type CommandUsecase interface {
	Register(ctx context.Context, payload *models.Register) utils.Result
	Login(ctx context.Context, payload *models.LoginRequest) utils.Result
	SaveAvatar(ctx context.Context, file multipart.File, header *multipart.FileHeader, ID string) utils.Result
}

type QueryUsecase interface {
	// ViewProfile(ctx context.Context, userId string) utils.Result
}
