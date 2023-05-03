package commands

import "gorm.io/gorm"

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
