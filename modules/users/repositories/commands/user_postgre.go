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

type CommandPayload struct {
	Table     string
	Query     string
	Parameter map[string]interface{}
	Where     map[string]interface{}
	Select    string
	Join      string
	Limit     int
	Offset    int
	Order     string
	Group     string
	Distinct  string
	Document  interface{}
	Output    interface{}
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

func (c *UserpostgreCommand) UpdatedUser(payload *CommandPayload) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result := c.db.Table(payload.Table).Where(payload.Query, payload.Parameter).Updates(payload.Document)
		if result.Error != nil {
			output <- utils.Result{Error: result}
		}

	}()

	return output
}

func (c *UserpostgreCommand) UpdateProfile(ctx context.Context, data *models.UpdatedUser) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		res := c.db.Model("users").Where("id = ?", data.ID).Updates(map[string]interface{}{
			"name":       data.Name,
			"email":      data.Email,
			"password":   data.Password,
			"avatar":     data.Avatar,
			"occupation": data.Occupation,
			"updated_at": data.UpdatedAt,
			"updated_by": data.UpdatedBy,
		})

		if res.Error != nil {
			output <- utils.Result{Error: res.Error}
		}

		output <- utils.Result{Data: res.RowsAffected}

	}()
	return output
}
