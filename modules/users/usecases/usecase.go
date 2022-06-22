package usecases

import (
	"context"
	models "crowdfunding/modules/users/models/domain"
	"crowdfunding/pkg/utils"
)

// CommandUsecase interface
type CommandUsecase interface {
	// UpdateProfile(ctx context.Context, payload *models.UpdatedUser) utils.Result
	Register(ctx context.Context, payload *models.Register) utils.Result
	// Login(ctx context.Context, payload *models.LoginRequest) utils.Result
	// Logout(ctx context.Context, payload *token.Claim) utils.Result
	// RefreshToken(ctx context.Context, payload *models.RefreshTokenRequest) utils.Result
	// ChangePassword(ctx context.Context, payload *models.ChangePassword) utils.Result
	// VerifyOtpCode(ctx context.Context, payload *models.VerifyOtp) utils.Result
	// ForgotAccount(ctx context.Context, payload *models.ForgotAccount) utils.Result
	// UpdateAvatar(ctx context.Context, file multipart.File, header *multipart.FileHeader, id string) utils.Result
	// GenerateOTP(ctx context.Context, payload *models.GenerateOTP) utils.Result
}

type QueryUsecase interface {
	// ViewProfile(ctx context.Context, userId string) utils.Result
}
