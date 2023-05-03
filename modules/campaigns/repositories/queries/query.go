package queries

import "crowdfunding/pkg/utils"

type CampaignsPostgre interface {
	CountDataJoin(payload *QueryPayload) <-chan utils.ResultCount
	FindManyJoin(payload *QueryPayload) <-chan utils.Result
}
