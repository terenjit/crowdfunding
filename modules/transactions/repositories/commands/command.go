package commands

import (
	"crowdfunding/pkg/utils"
)

type CommandPostgre interface {
	InsertOne(table string, document interface{}) <-chan utils.Result
	Update(table string, document interface{}) <-chan utils.Result
}
