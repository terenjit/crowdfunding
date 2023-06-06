package commands

import (
	models "crowdfunding/modules/campaigns/models/domain"
	"crowdfunding/pkg/utils"

	"gorm.io/gorm"
)

type PostgreCommand struct {
	db *gorm.DB
}

type CommandPayload struct {
	Table     string
	Select    string
	Join      string
	Query     interface{}
	Parameter map[string]interface{}
	Document  interface{}
	Output    interface{}
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

		result := c.db.Table(table).Create(document)
		if result.Error != nil {
			output <- utils.Result{Error: result}
		}

		output <- utils.Result{Data: document}
	}()

	return output
}

func (c *PostgreCommand) Update(payload *CommandPayload) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result := c.db.Debug().Table(payload.Table).Where(payload.Query, payload.Parameter).Updates(payload.Document)
		if result.Error != nil {
			output <- utils.Result{Error: result}
		}

	}()

	return output
}

func (c *PostgreCommand) UploadImages(param string, data *models.CampaignImages) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)
		var images models.CampaignImages
		res := c.db.Table("campaign_images").Model(&images).Where(param).Create(data)
		if res.Error != nil {
			output <- utils.Result{Error: res.Error}
		}

		output <- utils.Result{Data: res.RowsAffected}
	}()
	return output
}

func (c *PostgreCommand) MarkAllImagesAsNonPrimary(id string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)
		var images models.CampaignImagesFormatter
		res := c.db.Table("campaign_images").Model(&images).Where("campaign_id = ?", id).Update("is_primary", 0)
		if res.Error != nil {
			output <- utils.Result{Error: res.Error}
		}

		output <- utils.Result{Data: res.RowsAffected}
	}()
	return output
}

func (c *PostgreCommand) FindOne(payload *CommandPayload) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)
		var data models.Campaign
		result := c.db.Debug().Preload("User").Table(payload.Table).Select(payload.Select).Where(payload.Query, payload.Parameter).Joins(payload.Join).Limit(1).Find(&data)
		if result.Error != nil || result.RowsAffected == 0 {
			output <- utils.Result{
				Error: "Data Not Found",
			}
		}
		output <- utils.Result{Data: data}
	}()

	return output
}
