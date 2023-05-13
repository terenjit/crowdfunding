package commands

import (
	models "crowdfunding/modules/campaigns/models/domain"
	"crowdfunding/pkg/utils"
)

type CommandPostgre interface {
	InsertOne(table string, document interface{}) <-chan utils.Result
	Update(payload *CommandPayload) <-chan utils.Result
	UploadImages(param string, data *models.CampaignImages) <-chan utils.Result
	MarkAllImagesAsNonPrimary(id string) <-chan utils.Result
}
