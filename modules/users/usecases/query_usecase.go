package usecases

import "crowdfunding/modules/users/repositories/queries"

type userQueryUsecase struct {
	userPostgreQuery queries.UserPostgre
}

func NewUserQueryUsecase(userPostgreQuery queries.UserPostgre) *userQueryUsecase {
	return &userQueryUsecase{
		userPostgreQuery: userPostgreQuery,
	}
}
