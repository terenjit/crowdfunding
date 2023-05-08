package usecases

import (
	"context"
	models "crowdfunding/modules/campaigns/models/domain"
	"crowdfunding/modules/campaigns/repositories/commands"
	"crowdfunding/modules/campaigns/repositories/queries"
	"crowdfunding/pkg/utils"
	"time"
)

type commandUsecase struct {
	postgreQuery   queries.CampaignsPostgre
	postgreCommand commands.CommandPostgre
}

func NewCommandUsecasePostgre(postgreQuery queries.CampaignsPostgre, postgreCommand commands.CommandPostgre) *commandUsecase {
	return &commandUsecase{
		postgreQuery:   postgreQuery,
		postgreCommand: postgreCommand,
	}
}

func (c commandUsecase) Update(ctx context.Context, payload *models.UpdateCampaign) utils.Result {
	var result utils.Result

	var query string
	parameter := make(map[string]interface{})

	query = "c.id = @id"
	parameter["id"] = payload.ID

	currentTime := time.Now()
	payload.UpdatedAt = currentTime
	payload.UpdatedBy = payload.UserID

	queryPayload := commands.CommandPayload{
		Table:     "campaigns c",
		Query:     query,
		Parameter: parameter,
		Document:  payload,
	}

	c.postgreCommand.Update(&queryPayload)

	result.Data = payload
	return result
}
