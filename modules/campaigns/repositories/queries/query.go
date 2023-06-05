package queries

import (
	"context"
	"crowdfunding/pkg/utils"
)

type CampaignsPostgre interface {
	CountData(payload *QueryPayload) <-chan utils.ResultCount
	FindManyJoin(payload *QueryPayload) <-chan utils.Result
	FindOneJoin(payload *QueryPayload) <-chan utils.Result
	FindOne(payload *QueryPayload) <-chan utils.Result
	FindOneByID(ctx context.Context, ID string) <-chan utils.Result
	FindManyJoinDetail(payload *QueryPayload) <-chan utils.Result
}
