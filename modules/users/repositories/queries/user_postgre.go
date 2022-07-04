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
