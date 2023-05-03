package queries

import (
	"context"
	"crowdfunding/pkg/utils"

	"gorm.io/gorm"
)

type CampaignsPostgreQuery struct {
	db    *gorm.DB
	table string
}

func NewPostgreQuery(db *gorm.DB) *CampaignsPostgreQuery {
	return &CampaignsPostgreQuery{
		db:    db,
		table: "campaigns",
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

func (c *CampaignsPostgreQuery) FindManyJoin(payload *QueryPayload) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		var data = make([]map[string]interface{}, 0)

		result := c.db.Table(payload.Table).Select(payload.Select).Where(payload.Where).Offset(payload.Offset).Limit(payload.Limit).Joins(payload.Join).Order(payload.Order).Find(&data)
		if result.Error != nil {
			output <- utils.Result{
				Error: result.Error,
			}
		}
		output <- utils.Result{Data: data, Count: result.RowsAffected}
	}()

	return output
}

func (q *CampaignsPostgreQuery) CountDataJoin(payload *QueryPayload) <-chan utils.ResultCount {
	output := make(chan utils.ResultCount)

	go func() {
		defer close(output)

		var data int64
		result := q.db.Table(payload.Table).Select(payload.Select).Where(payload.Where).Limit(payload.Limit).Offset(payload.Offset).Joins(payload.Join).Count(&data)
		if result.Error != nil {
			output <- utils.ResultCount{
				Error: "Data Not Found",
			}
		}
		output <- utils.ResultCount{Data: data}
	}()

	return output
}
