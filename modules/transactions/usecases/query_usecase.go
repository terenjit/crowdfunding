package usecases

import (
	"context"
	models "crowdfunding/modules/transactions/models/domain"
	"crowdfunding/modules/transactions/repositories/queries"
	"crowdfunding/pkg/utils"
	"encoding/json"
	"math"
)

type transactionsQueryUsecase struct {
	transactionsPostgreQuery queries.TransactionsPostgre
}

func NewQueryUsecase(transactionsPostgreQuery queries.TransactionsPostgre) *transactionsQueryUsecase {
	return &transactionsQueryUsecase{
		transactionsPostgreQuery: transactionsPostgreQuery,
	}
}

func (q transactionsQueryUsecase) ListTransactions(ctx context.Context, payload *models.TransactionList) utils.Result {
	var result utils.Result
	var queryRes utils.Result

	parameter := make(map[string]interface{})
	parameter["is_deleted"] = false
	parameter["campaign_id"] = payload.CampaignID

	queryPayload := queries.QueryPayload{
		Table:     "transactions t",
		Select:    "t.id",
		Query:     "t.campaign_id = @campaign_id AND t.is_deleted = @is_deleted",
		Parameter: parameter,
		Join:      "left join users u on u.id = t.user_id",
		Output:    []models.Transaction{},
	}

	count := <-q.transactionsPostgreQuery.CountData(&queryPayload)
	if count.Error != nil {
		result.Data = []models.Transaction{}
		return result
	}

	if payload.Quantity == 0 {
		payload.Quantity = int(count.Data)
	}

	offset := payload.Quantity * (payload.Page - 1)
	queryPayload.Select = "t.*, u.name as name"
	queryPayload.Offset = offset
	queryPayload.Limit = payload.Quantity
	if count.Error == nil || count.Data > 0 {
		queryRes = <-q.transactionsPostgreQuery.FindManyJoin(&queryPayload)
		if queryRes.Error != nil {
			queryRes.Data = []models.Transaction{}
			count.Data = 0
		}
	}

	var data []models.Transaction
	byteStatus, _ := json.Marshal(queryRes.Data)
	json.Unmarshal(byteStatus, &data)

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

func (q transactionsQueryUsecase) ListUserTransactions(ctx context.Context, payload *models.TransactionList) utils.Result {
	var result utils.Result
	var queryRes utils.Result

	parameter := make(map[string]interface{})
	parameter["is_deleted"] = false
	parameter["user_id"] = payload.UserID

	join := `left join users u on u.id = t.user_id
			left join campaigns c on c.id = t.campaign_id
			left join campaign_images ci on ci.campaign_id = c.id`

	queryPayload := queries.QueryPayload{
		Table:     "transactions t",
		Select:    "t.id",
		Query:     "t.user_id = @user_id AND t.is_deleted = @is_deleted AND ci.is_primary = 1",
		Parameter: parameter,
		Join:      join,
		Output:    []models.Transaction{},
	}

	count := <-q.transactionsPostgreQuery.CountData(&queryPayload)
	if count.Error != nil {
		result.Data = []models.Transaction{}
		return result
	}

	if payload.Quantity == 0 {
		payload.Quantity = int(count.Data)
	}

	offset := payload.Quantity * (payload.Page - 1)
	queryPayload.Select = "t.*, u.name as name, ci.file_name as file_name, ci.is_primary"
	queryPayload.Offset = offset
	queryPayload.Limit = payload.Quantity
	if count.Error == nil || count.Data > 0 {
		queryRes = <-q.transactionsPostgreQuery.FindManyJoin(&queryPayload)
		if queryRes.Error != nil {
			queryRes.Data = []models.Transaction{}
			count.Data = 0
		}
	}

	var data []models.Transaction
	byteStatus, _ := json.Marshal(queryRes.Data)
	json.Unmarshal(byteStatus, &data)

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
