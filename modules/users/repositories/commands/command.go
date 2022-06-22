package commands

import (
	"context"
	models "crowdfunding/modules/users/models/domain"
	"crowdfunding/pkg/utils"
)

type UserMongo interface {
	// InsertOneUser(ctx context.Context, data *models.User) <-chan utils.Result
	// UpdateProfile(ctx context.Context, data *models.User) <-chan utils.Result
	// UpdatePasswordUser(ctx context.Context, data *models.ChangePassword) <-chan utils.Result
}

type UserPostgre interface {
	InsertOneUser(ctx context.Context, data *models.User) <-chan utils.Result
	// UpdateProfile(ctx context.Context, data *models.User) <-chan utils.Result
	// UpdatePasswordUser(ctx context.Context, data *models.ChangePassword) <-chan utils.Result
	// UpdateOne(param string, data *models.User) <-chan utils.Result
}
