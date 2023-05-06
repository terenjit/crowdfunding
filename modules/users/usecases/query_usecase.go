package usecases

import (
	"context"
	models "crowdfunding/modules/users/models/domain"
	"crowdfunding/modules/users/repositories/queries"
	httpError "crowdfunding/pkg/http-error"
	"crowdfunding/pkg/utils"
	"encoding/json"
	"math"
)

type userQueryUsecase struct {
	userPostgreQuery queries.UserPostgre
}

func NewUserQueryUsecase(userPostgreQuery queries.UserPostgre) *userQueryUsecase {
	return &userQueryUsecase{
		userPostgreQuery: userPostgreQuery,
	}
}

func (q userQueryUsecase) GetDetail(ctx context.Context, id string) utils.Result {
	var result utils.Result

	queryRes := <-q.userPostgreQuery.FindOneByID(ctx, id)
	if queryRes.Error != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "User tidak ditemukan"
		result.Error = errObj
		return result
	}

	var dataUser models.User
	jsonCampaign, _ := json.Marshal(queryRes.Data)
	json.Unmarshal(jsonCampaign, &dataUser)

	parameter := make(map[string]interface{})
	parameter["id"] = id

	result.Data = dataUser
	return result
}

func (q userQueryUsecase) GetList(ctx context.Context, payload *models.UsersGetList) utils.Result {
	var result utils.Result
	var queryRes utils.Result

	var query string
	parameter := make(map[string]interface{})

	query = "u.is_deleted = @is_deleted"
	parameter["is_deleted"] = false

	if payload.Page == 0 {
		payload.Page = 1
	}
	if payload.UserID != "" {
		query = query + " " + "AND user_id = @user_id"
		parameter["user_id"] = payload.UserID
	}

	queryPayload := queries.QueryPayload{
		Table:     "users u",
		Query:     query,
		Parameter: parameter,
		Select:    "u.id",
		Output:    []models.User{},
	}

	count := <-q.userPostgreQuery.CountData(&queryPayload)
	if count.Error != nil {
		result.Data = []models.User{}
		return result
	}

	if payload.Quantity == 0 {
		payload.Quantity = int(count.Data)
	}

	offset := payload.Quantity * (payload.Page - 1)
	queryPayload.Select = "u.*"
	queryPayload.Offset = offset
	queryPayload.Limit = payload.Quantity
	queryPayload.Order = "u asc"
	if count.Error == nil || count.Data > 0 {
		queryRes = <-q.userPostgreQuery.FindMany(&queryPayload)
		if queryRes.Error != nil {
			queryRes.Data = []models.User{}
			count.Data = 0
		}
	}

	var data = []models.User{}
	jsonCampaign, _ := json.Marshal(queryRes.Data)
	json.Unmarshal(jsonCampaign, &data)

	result.Data = data
	totalData := count.Data
	result.MetaData = map[string]interface{}{
		"page":      payload.Page,
		"quantity":  queryRes.Count,
		"totalPage": math.Ceil(float64((totalData + int64(payload.Quantity) - 1) / int64(payload.Quantity))),
		"totalData": totalData,
	}
	return result
}
