package queries

import (
	"context"
	"crowdfunding/pkg/utils"

	"gorm.io/gorm"
)

type TransactionsPostgreQuery struct {
	db    *gorm.DB
	table string
}

func NewPostgreQuery(db *gorm.DB) *TransactionsPostgreQuery {
	return &TransactionsPostgreQuery{
		db:    db,
		table: "transactions",
	}
}

type QueryPayload struct {
	Ctx       context.Context
	Table     string
	Query     string
	Parameter map[string]interface{}
	Where     map[string]interface{}
	Select    string
	Join      string
	Limit     int
	Offset    int
	Order     string
	Id        string
	Output    interface{}
	Group     string
	Distinct  string
}

func (q *TransactionsPostgreQuery) CountData(payload *QueryPayload) <-chan utils.ResultCount {
	output := make(chan utils.ResultCount)

	go func() {
		defer close(output)

		var data int64
		result := q.db.Debug().Table(payload.Table).Select(payload.Select).Where(payload.Query, payload.Parameter).Limit(payload.Limit).Offset(payload.Offset).Count(&data)
		if result.Error != nil {
			output <- utils.ResultCount{
				Error: "Data Not Found",
			}
		}
		output <- utils.ResultCount{Data: data}
	}()

	return output
}

func (c *TransactionsPostgreQuery) FindManyJoin(payload *QueryPayload) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result := c.db.Debug().Table(payload.Table).Select(payload.Select).Where(payload.Query, payload.Parameter).Offset(payload.Offset).Limit(payload.Limit).Joins(payload.Join).Order(payload.Order).Find(&payload.Output)
		if result.Error != nil {
			output <- utils.Result{
				Error: result.Error,
			}
		}
		output <- utils.Result{Data: payload.Output, Count: result.RowsAffected}
	}()

	return output
}
