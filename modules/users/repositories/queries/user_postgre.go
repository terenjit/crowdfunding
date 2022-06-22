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
