package usecases

import (
	"context"
	models "crowdfunding/modules/campaigns/models/domain"
	"crowdfunding/pkg/utils"
	"mime/multipart"
)

type QueryUsecase interface {
	GetList(ctx context.Context, payload *models.CampainGetList) utils.Result
	GetDetail(ctx context.Context, payload *models.CampaignGetDetail) utils.Result
}

// CommandUsecase interface
type CommandUsecase interface {
	Create(ctx context.Context, payload *models.CreateRequest) utils.Result
	Update(ctx context.Context, payload *models.UpdateCampaign) utils.Result
	UploadCampaignImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, payload *models.UploadImageRequest) utils.Result
}
