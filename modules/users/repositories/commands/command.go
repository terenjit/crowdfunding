package commands

import (
	"context"
	models "crowdfunding/modules/users/models/domain"
	"crowdfunding/pkg/utils"
)

type UserPostgre interface {
	InsertOneUser(ctx context.Context, data *models.User) <-chan utils.Result
	Update(param string, data *models.User) <-chan utils.Result
}
