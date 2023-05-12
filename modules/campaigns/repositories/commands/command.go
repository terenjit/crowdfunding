package commands

import "crowdfunding/pkg/utils"

type CommandPostgre interface {
	InsertOne(table string, document interface{}) <-chan utils.Result
	Update(payload *CommandPayload) <-chan utils.Result
}
