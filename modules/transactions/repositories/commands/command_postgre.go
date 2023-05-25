package commands

import (
	models "crowdfunding/modules/transactions/models/domain"
	"crowdfunding/pkg/utils"

	"gorm.io/gorm"
)

type PostgreCommand struct {
	db *gorm.DB
}

type CommandPayload struct {
	Table     string
	Query     interface{}
	Parameter map[string]interface{}
	Document  interface{}
	Raw       string
}

func NewPostgreCommand(db *gorm.DB) *PostgreCommand {
	return &PostgreCommand{
		db: db,
	}
}

func (c *PostgreCommand) InsertOne(table string, document interface{}) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result := c.db.Debug().Table(table).Create(document)
		if result.Error != nil {
			output <- utils.Result{Error: result}
		}

		output <- utils.Result{Data: document}
	}()

	return output
}

func (c *PostgreCommand) Update(table string, document interface{}) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)
		result := c.db.Debug().Table(table).Updates(document)
		if result.Error != nil {
			output <- utils.Result{Error: result}
		}
	}()

	return output
}

func (c *PostgreCommand) FindByID(ID string) (models.TransactionModel, error) {

	var transaction models.TransactionModel

	err := c.db.Where("id = ?", ID).Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (c *PostgreCommand) FindCampaignByID(ID string) (models.CampaignModel, error) {

	var campaign models.CampaignModel

	err := c.db.Where("id = ?", ID).Find(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (c *PostgreCommand) UpdateTransaction(data models.TransactionModel) (models.TransactionModel, error) {

	var transaction models.TransactionModel

	result := c.db.Debug().Updates(&transaction).Error
	if result != nil {
		return transaction, result
	}

	return transaction, result
}
