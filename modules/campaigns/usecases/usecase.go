package usecases

import (
	"context"
	models "crowdfunding/modules/campaigns/models/domain"
	"crowdfunding/pkg/utils"
)

type CommandUsecase interface {
}

type QueryUsecase interface {
	GetList(ctx context.Context, payload *models.CampainGetList) utils.Result
	GetDetail(ctx context.Context, id string) utils.Result
}
