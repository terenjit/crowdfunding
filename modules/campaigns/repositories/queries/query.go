package queries

import "crowdfunding/pkg/utils"

type CampaignsPostgre interface {
	CountData(payload *QueryPayload) <-chan utils.ResultCount
	FindManyJoin(payload *QueryPayload) <-chan utils.Result
}
