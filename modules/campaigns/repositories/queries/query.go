package queries

import (
	"context"
	"crowdfunding/pkg/utils"
)

type CampaignsPostgre interface {
	CountData(payload *QueryPayload) <-chan utils.ResultCount
	FindManyJoin(payload *QueryPayload) <-chan utils.Result
	FindOneJoin(payload *QueryPayload) <-chan utils.Result
	FindOne(ctx context.Context, id string) <-chan utils.Result
}
