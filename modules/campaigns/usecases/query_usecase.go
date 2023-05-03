package usecases

import (
	"context"
	models "crowdfunding/modules/campaigns/models/domain"
	"crowdfunding/modules/campaigns/repositories/queries"
	"crowdfunding/pkg/utils"
)

type campaignQueryUsecase struct {
	campaignPostgreQuery queries.CampaignsPostgre
}

func NewQueryUsecase(campaignPostgreQuery queries.CampaignsPostgre) *campaignQueryUsecase {
	return &campaignQueryUsecase{
		campaignPostgreQuery: campaignPostgreQuery,
	}
}

func (q campaignQueryUsecase) GetList(ctx context.Context, payload *models.CampainGetList) utils.Result {
	var result utils.Result

	return result
}
