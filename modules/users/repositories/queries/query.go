package queries

import (
	"context"
	"crowdfunding/pkg/utils"
)

type UserMongo interface {
	// FindOne(ctx context.Context, parameter bson.M) <-chan utils.Result
	// FindOneByUsername(ctx context.Context, username, phoneNumber string) <-chan utils.Result
	// FindOneByUserId(ctx context.Context, userId string) <-chan utils.Result
	// FindOneByPhoneNumber(ctx context.Context, phoneNumber string) <-chan utils.Result
	// FindOneByEmail(ctx context.Context, email string) <-chan utils.Result
	// FindOneByNIK(ctx context.Context, nik string) <-chan utils.Result
}

type UserPostgre interface {
	// FindOneByUsername(ctx context.Context, username, phoneNumber string) <-chan utils.Result
	// FindOneByUserId(ctx context.Context, userId string) <-chan utils.Result
	// FindOneByPhoneNumber(ctx context.Context, phoneNumber string) <-chan utils.Result
	// FindOneByEmail(ctx context.Context, email string) <-chan utils.Result
	// FindOneByNIK(ctx context.Context, nik string) <-chan utils.Result
	// ViewProfile(ctx context.Context, userId string) <-chan utils.Result
	FindOne(ctx context.Context, query string, parameter map[string]interface{}) <-chan utils.Result
}
