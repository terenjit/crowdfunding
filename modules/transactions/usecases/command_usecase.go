package usecases

import (
	"crowdfunding/modules/transactions/repositories/commands"
	"crowdfunding/modules/transactions/repositories/queries"
)

type commandUsecase struct {
	postgreQuery   queries.TransactionsPostgre
	postgreCommand commands.CommandPostgre
}

func NewCommandUsecasePostgre(postgreQuery queries.TransactionsPostgre, postgreCommand commands.CommandPostgre) *commandUsecase {
	return &commandUsecase{
		postgreQuery:   postgreQuery,
		postgreCommand: postgreCommand,
	}
}
