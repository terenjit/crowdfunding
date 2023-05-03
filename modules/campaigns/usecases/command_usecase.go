package usecases

import (
	"crowdfunding/modules/campaigns/repositories/commands"
	"crowdfunding/modules/campaigns/repositories/queries"
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
