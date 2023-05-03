package queries

import (
	"context"

	"gorm.io/gorm"
)

type SubvesselsPostgreQuery struct {
	db    *gorm.DB
	table string
}

func NewPostgreQuery(db *gorm.DB) *SubvesselsPostgreQuery {
	return &SubvesselsPostgreQuery{
		db:    db,
		table: "campaigns",
	}
}

type QueryPayload struct {
	Ctx       context.Context
	Table     string
	Query     string
	Parameter map[string]interface{}
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
