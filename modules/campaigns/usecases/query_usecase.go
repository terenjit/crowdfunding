package usecases

import (
	"context"
	models "crowdfunding/modules/campaigns/models/domain"
	"crowdfunding/modules/campaigns/repositories/queries"
	httpError "crowdfunding/pkg/http-error"
	"crowdfunding/pkg/utils"
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

type campaignQueryUsecase struct {
	campaignPostgreQuery queries.CampaignsPostgre
}

func NewQueryUsecase(campaignPostgreQuery queries.CampaignsPostgre) *campaignQueryUsecase {
	return &campaignQueryUsecase{
		campaignPostgreQuery: campaignPostgreQuery,
	}
}

func (q campaignQueryUsecase) GetDetail(ctx context.Context, payload *models.CampaignGetDetail) utils.Result {
	var result utils.Result

	parameter := make(map[string]interface{})
	parameter["is_deleted"] = false
	parameter["id"] = payload.ID

	queryPayload := queries.QueryPayload{
		Table:     "campaigns c",
		Select:    "c.* ",
		Query:     "c.id = @id AND c.is_deleted = @is_deleted",
		Parameter: parameter,
		Output:    models.Campaign{},
	}

	queryRes := <-q.campaignPostgreQuery.FindOne(&queryPayload)
	if queryRes.Error != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "Campaign tidak ditemukan"
		result.Error = errObj
		return result
	}
	var campaign models.Campaign
	jsonCampaign, _ := json.Marshal(queryRes.Data)
	json.Unmarshal(jsonCampaign, &campaign)

	var perks []string

	for _, perk := range strings.Split(campaign.Perks, ",") {
		perks = append(perks, strings.TrimSpace(perk))
	}

	var dataCampaign models.CampaignFormatter
	jsonCampaignFormatter, _ := json.Marshal(queryRes.Data)
	json.Unmarshal(jsonCampaignFormatter, &dataCampaign)

	dataCampaign.Perks = perks

	queryPayload = queries.QueryPayload{
		Table:     "campaign_images ci",
		Select:    "ci.file_name, ci.is_primary",
		Query:     "ci.campaign_id = @id",
		Parameter: parameter,
		Output:    []models.CampaignImagesFormatter{},
	}
	campaignImages := <-q.campaignPostgreQuery.FindManyJoin(&queryPayload)
	if campaignImages.Error == nil {
		dataCampaign.Images = campaignImages.Data.([]models.CampaignImagesFormatter)
	}

	var data []models.CampaignImagesFormatter
	byteImages, _ := json.Marshal(queryRes.Data)
	json.Unmarshal(byteImages, &data)

	for _, v := range campaign.Images {
		CampaignImagesFormatter := models.CampaignImagesFormatter{}
		CampaignImagesFormatter.FileName = v.FileName
		isPrimary := false

		if v.IsPrimary == 1 {
			isPrimary = true
		}
		CampaignImagesFormatter.IsPrimary = isPrimary
		data = append(data, CampaignImagesFormatter)
	}

	result.Data = dataCampaign
	return result
}

func (q campaignQueryUsecase) GetList(ctx context.Context, payload *models.CampainGetList) utils.Result {
	var result utils.Result
	var queryRes utils.Result

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
	if payload.Search != "" {
		query = query + " AND (c.name ILike @search)"
		parameter["search"] = fmt.Sprintf("%v%v%v", "%", payload.Search, "%")
	}

	join := `left join campaign_images ci on ci.campaign_id = c.id`

	queryPayload := queries.QueryPayload{
		Table:     "campaigns c",
		Query:     query,
		Parameter: parameter,
		Select:    "c.id",
		Join:      join,
		Output:    []models.CampaignsFormat{},
	}

	count := <-q.campaignPostgreQuery.CountData(&queryPayload)
	if count.Error != nil {
		result.Data = []models.CampaignsFormat{}
		return result
	}

	if payload.Quantity == 0 {
		payload.Quantity = int(count.Data)
	}

	offset := payload.Quantity * (payload.Page - 1)
	queryPayload.Select = "c.*, ci.file_name as images_url"
	queryPayload.Offset = offset
	queryPayload.Limit = payload.Quantity
	if count.Error == nil || count.Data > 0 {
		queryRes = <-q.campaignPostgreQuery.FindManyJoin(&queryPayload)
		if queryRes.Error != nil {
			queryRes.Data = []models.CampaignsFormat{}
			count.Data = 0
		}
	}

	dataCampaign := queryRes.Data.([]models.CampaignsFormat)
	// for i := 0; i < len(dataCampaign); i++ {
	// 	parameter["id"] = dataCampaign[i].ID

	// 	queryPayload := queries.QueryPayload{
	// 		Table:     "campaign_images ci",
	// 		Select:    "ci.id , ci.campaign_id, ci.file_name, ci.is_primary",
	// 		Query:     "ci.campaign_id = @id AND ci.is_primary = 1",
	// 		Parameter: parameter,
	// 		Output:    []models.CampaignImages{},
	// 	}

	// 	campaignImages := <-q.campaignPostgreQuery.FindManyJoin(&queryPayload)
	// 	if campaignImages.Error == nil {
	// 		dataCampaign[i].Images = campaignImages.Data.([]models.CampaignImages)
	// 	}
	// }

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
