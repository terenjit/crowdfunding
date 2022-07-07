package commands

import (
	"context"
	models "crowdfunding/modules/users/models/domain"
	"crowdfunding/pkg/utils"

	"gorm.io/gorm"
)

// UserpostgreCommand model
type UserpostgreCommand struct {
	db *gorm.DB
}

// NewUserPostgreQuery create new user query
func NewUserPostgreCommand(db *gorm.DB) *UserpostgreCommand {
	return &UserpostgreCommand{
		db: db,
	}
}

func (c *UserpostgreCommand) InsertOneUser(ctx context.Context, data *models.User) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)
		result := c.db.Create(&data)
		if result.Error != nil {
			output <- utils.Result{Error: result}
		}
		output <- utils.Result{Data: data}
	}()

	return output
}

func (c *UserpostgreCommand) Update(param string, data *models.User) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)
		var user models.User
		res := c.db.Model(&user).Where(param).Save(data)
		if res.Error != nil {
			output <- utils.Result{Error: res.Error}
		}

		output <- utils.Result{Data: res.RowsAffected}
	}()
	return output
}
