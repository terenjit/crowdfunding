package queries

import (
	"context"
	models "crowdfunding/modules/users/models/domain"
	"crowdfunding/pkg/utils"

	"gorm.io/gorm"
)

// UserpostgreQuery model
type UserpostgreQuery struct {
	db    *gorm.DB
	table string
}

// NewUserPostgreQuery create new user query
func NewUserPostgreQuery(db *gorm.DB) *UserpostgreQuery {
	return &UserpostgreQuery{
		db:    db,
		table: "users",
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

func (q *UserpostgreQuery) CountData(payload *QueryPayload) <-chan utils.ResultCount {
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

func (c *UserpostgreQuery) FindOne(ctx context.Context, query string, parameter map[string]interface{}) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		var data models.User
		result := c.db.Table(c.table).Model(&data).Where(query, parameter).Find(&data)
		if result.Error != nil {
			output <- utils.Result{
				Error: result.Error,
			}
		}
		output <- utils.Result{Data: data}
	}()

	return output
}

// find one by Username
func (c *UserpostgreQuery) FindOneByUsername(ctx context.Context, username string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		var user models.User
		result := c.db.Model(&user).Where(c.db.Where("username = ?", username)).Find(&user)
		if result.Error != nil {
			output <- utils.Result{
				Error: result.Error,
			}
		}
		output <- utils.Result{Data: user}
	}()

	return output
}

func (c *UserpostgreQuery) FindOneByEmail(ctx context.Context, email string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		var user models.User

		result := c.db.Where("email = ?", email).Find(&user)
		if result.Error != nil {
			output <- utils.Result{
				Error: result.Error,
			}
		}
		output <- utils.Result{Data: user}
	}()

	return output
}

func (c *UserpostgreQuery) FindOneByID(ctx context.Context, ID string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		var user models.User

		result := c.db.Where("id = ?", ID).Find(&user)
		if result.Error != nil {
			output <- utils.Result{
				Error: result.Error,
			}
		}
		output <- utils.Result{Data: user}
	}()

	return output
}

func (c *UserpostgreQuery) FindMany(payload *QueryPayload) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result := c.db.Debug().Table(payload.Table).Select(payload.Select).Where(payload.Query, payload.Parameter).Offset(payload.Offset).Limit(payload.Limit).Order(payload.Order).Find(&payload.Output)
		if result.Error != nil {
			output <- utils.Result{
				Error: result.Error,
			}
		}
		output <- utils.Result{Data: payload.Output, Count: result.RowsAffected}
	}()

	return output
}
