package queries

import (
	"context"
	models "crowdfunding/modules/campaigns/models/domain"
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

func (c *CampaignsPostgreQuery) FindOne(payload *QueryPayload) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)
		var data models.Campaign
		result := c.db.Debug().Preload("User").Table(payload.Table).Select(payload.Select).Where(payload.Query, payload.Parameter).Limit(1).Find(&data)
		if result.Error != nil || result.RowsAffected == 0 {
			output <- utils.Result{
				Error: "Data Not Found",
			}
		}
		output <- utils.Result{Data: data}
	}()

	return output
}

func (c *CampaignsPostgreQuery) FindOneJoin(payload *QueryPayload) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		var data map[string]interface{}
		result := c.db.Table(payload.Table).Select(payload.Select).Where(payload.Parameter).Joins(payload.Join).Find(&data)
		if result.Error != nil {
			output <- utils.Result{
				Error: "Data Not Found",
			}
		}

		output <- utils.Result{Data: data}
	}()

	return output
}

func (c *CampaignsPostgreQuery) FindManyJoin(payload *QueryPayload) <-chan utils.Result {
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

func (q *CampaignsPostgreQuery) CountData(payload *QueryPayload) <-chan utils.ResultCount {
	output := make(chan utils.ResultCount)

	go func() {
		defer close(output)

		var data int64
		result := q.db.Table(payload.Table).Select(payload.Select).Where(payload.Query, payload.Parameter).Limit(payload.Limit).Offset(payload.Offset).Count(&data)
		if result.Error != nil {
			output <- utils.ResultCount{
				Error: "Data Not Found",
			}
		}
		output <- utils.ResultCount{Data: data}
	}()

	return output
}
