package usecases

import (
	"context"
	models "crowdfunding/modules/campaigns/models/domain"
	"crowdfunding/modules/campaigns/repositories/queries"
	"crowdfunding/pkg/utils"
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
	//var count utils.ResultCount
	// queryRes.Data = []models.Campaign{}
	//count.Data = 0
	// quantity := payload.Quantity
	// page := payload.Page
	// offset := quantity * (page - 1)

	var query string
	parameter := make(map[string]interface{})

	query = "c.is_deleted = @is_deleted"
	parameter["is_deleted"] = false

	if payload.Page == 0 {
		payload.Page = 1
	}
	if payload.UserID != "" {
		query = query + " " + "AND user_id = @user_id"
		parameter["user_id"] = payload.UserID
	}

	join := `left join campaign_images ci on ci.campaign_id = c.id`

	queryPayload := queries.QueryPayload{
		Table:     "campaigns c",
		Query:     query,
		Parameter: parameter,
		Select:    "c.id",
		Join:      join,
		Output:    []models.Campaign{},
	}

	count := <-q.campaignPostgreQuery.CountData(&queryPayload)
	if count.Error != nil {
		result.Data = []models.Campaign{}
		return result
	}

	if payload.Quantity == 0 {
		payload.Quantity = int(count.Data)
	}

	offset := payload.Quantity * (payload.Page - 1)
	queryPayload.Select = "c.*"
	queryPayload.Offset = offset
	queryPayload.Limit = payload.Quantity
	if count.Error == nil || count.Data > 0 {
		queryRes = <-q.campaignPostgreQuery.FindManyJoin(&queryPayload)
		if queryRes.Error != nil {
			queryRes.Data = []models.Campaign{}
			count.Data = 0
		}
	}

	dataCampaign := queryRes.Data.([]models.Campaign)
	for i := 0; i < len(dataCampaign); i++ {
		parameter["id"] = dataCampaign[i].ID

		queryPayload := queries.QueryPayload{
			Table:     "campaign_images ci",
			Select:    "ci.id , ci.campaign_id, ci.file_name, ci.is_primary",
			Query:     "ci.campaign_id = @id AND ci.is_primary = 1",
			Parameter: parameter,
			Output:    []models.CampaignImages{},
		}

		campaignImages := <-q.campaignPostgreQuery.FindManyJoin(&queryPayload)
		if campaignImages.Error == nil {
			dataCampaign[i].Images = campaignImages.Data.([]models.CampaignImages)
		}
	}

	result.Data = dataCampaign
	totalData := count.Data
	result.MetaData = map[string]interface{}{
		"page":      payload.Page,
		"quantity":  queryRes.Count,
		"totalPage": math.Ceil(float64((totalData + int64(payload.Quantity) - 1) / int64(payload.Quantity))),
		"totalData": totalData,
	}
	return result
}
