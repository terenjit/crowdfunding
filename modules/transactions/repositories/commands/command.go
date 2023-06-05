package commands

import (
	"context"
	models "crowdfunding/modules/transactions/models/domain"
	"crowdfunding/pkg/utils"
)

type CommandPostgre interface {
	InsertOne(table string, document interface{}) <-chan utils.Result
	Update(ID string, document interface{}) <-chan utils.Result
	FindByID(ID string) (models.TransactionModel, error)
	FindCampaignByID(ID string) (models.CampaignModel, error)
	UpdateTransaction(data models.TransactionModel) (models.TransactionModel, error)
	FindOneByID(ctx context.Context, ID string) <-chan utils.Result
	FindOneCampaignByID(ctx context.Context, ID string) <-chan utils.Result
	UpdateCampaign(ID string, document interface{}) <-chan utils.Result
}
