package usecases

import (
	"context"
	models "crowdfunding/modules/campaigns/models/domain"
	"crowdfunding/modules/campaigns/repositories/queries"
	"crowdfunding/pkg/utils"
	"encoding/json"
	"math"
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
	var queryRes utils.Result
	var count utils.ResultCount
	queryRes.Data = []models.Campaign{}
	count.Data = 0
	quantity := payload.Quantity
	page := payload.Page
	offset := quantity * (page - 1)

	var query string
	where := make(map[string]interface{})

	if payload.UserID != "" {
		query = "c.user_id = @c.user_id"
		where["c.user_id"] = payload.UserID
	}

	queryPayload := queries.QueryPayload{
		Table:  "campaigns c",
		Query:  query,
		Where:  where,
		Select: "c.id",
		Output: []models.Campaign{},
	}

	count = <-q.campaignPostgreQuery.CountDataJoin(&queryPayload)
	if count.Error != nil {
		return result
	}

	if payload.Quantity == 0 {
		payload.Quantity = int(count.Data)
	}

	if count.Error == nil || count.Data > 0 {
		queryPayload.Select = "c.*"
		queryPayload.Offset = offset
		queryPayload.Limit = payload.Quantity
		queryPayload.Order = "c.created_at asc "

		queryRes = <-q.campaignPostgreQuery.FindManyJoin(&queryPayload)
		if queryRes.Error != nil {
			queryRes.Data = []models.Campaign{}
			count.Data = 0
		}
	}

	var data []models.Campaign
	byteStatus, _ := json.Marshal(queryRes.Data)
	json.Unmarshal(byteStatus, &data)

	result.Data = data
	totalData := count.Data
	result.MetaData = map[string]interface{}{
		"page":      page,
		"quantity":  queryRes.Count,
		"totalPage": math.Ceil(float64((totalData + int64(payload.Quantity) - 1) / int64(payload.Quantity))),
		"totalData": totalData,
	}
	return result
}
