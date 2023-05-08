package commands

import "crowdfunding/pkg/utils"

type CommandPostgre interface {
	Update(payload *CommandPayload) <-chan utils.Result
}
