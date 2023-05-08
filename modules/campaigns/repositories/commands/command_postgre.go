package commands

import (
	"crowdfunding/pkg/utils"

	"gorm.io/gorm"
)

type PostgreCommand struct {
	db *gorm.DB
}

type CommandPayload struct {
	Table     string
	Query     interface{}
	Parameter map[string]interface{}
	Document  interface{}
	Raw       string
}

func NewPostgreCommand(db *gorm.DB) *PostgreCommand {
	return &PostgreCommand{
		db: db,
	}
}

func (c *PostgreCommand) Update(payload *CommandPayload) <-chan utils.Result {
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
