package queries

import "crowdfunding/pkg/utils"

type TransactionsPostgre interface {
	CountData(payload *QueryPayload) <-chan utils.ResultCount
	FindManyJoin(payload *QueryPayload) <-chan utils.Result
}
