package usecases

import (
	"context"
	models "crowdfunding/modules/transactions/models/domain"
	"crowdfunding/pkg/utils"
)

type QueryUsecase interface {
	ListTransactions(ctx context.Context, payload *models.TransactionList) utils.Result
	ListUserTransactions(ctx context.Context, payload *models.TransactionList) utils.Result
}

// CommandUsecase interface
type CommandUsecase interface {
	Create(ctx context.Context, payload *models.CreateRequest) utils.Result
	ProcessPayment(ctx context.Context, payload *models.TransactionNotificationInput) error
	// GetPaymentURL(ctx context.Context, payload *models.TransactionModel) (string, error)
}
