package usecases

import (
	"context"
	models "crowdfunding/modules/campaigns/models/domain"
	"crowdfunding/modules/campaigns/repositories/commands"
	"crowdfunding/modules/campaigns/repositories/queries"
	httpError "crowdfunding/pkg/http-error"
	"crowdfunding/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type commandUsecase struct {
	postgreQuery   queries.CampaignsPostgre
	postgreCommand commands.CommandPostgre
}

func NewCommandUsecasePostgre(postgreQuery queries.CampaignsPostgre, postgreCommand commands.CommandPostgre) *commandUsecase {
	return &commandUsecase{
		postgreQuery:   postgreQuery,
		postgreCommand: postgreCommand,
	}
}

func (c commandUsecase) Create(ctx context.Context, payload *models.CreateRequest) utils.Result {
	var result utils.Result

	requestUser := c.HttpRequest(http.MethodGet, fmt.Sprintf("localhost:9000/crowdfunding/v1/users/%s", payload.Opts.UserID), payload.Opts.Authorization, nil)
	if requestUser.Error == false {
		errObj := httpError.NewConflict()
		errObj.Message = "Data pengguna tidak ditemukan"
		result.Error = errObj
		return result
	}

	campaignId := uuid.NewString()

	createdSlug := fmt.Sprintf("%s %s", payload.Name, payload.Opts.UserID)

	data := models.Campaign{
		ID:               campaignId,
		UserID:           payload.Opts.UserID,
		Name:             payload.Name,
		ShortDescription: payload.ShortDescription,
		Description:      payload.Description,
		GoalAmount:       payload.GoalAmount,
		Slug:             slug.Make(createdSlug),
		Perks:            payload.Perks,
		CreatedAt:        time.Now(),
	}

	insertCampaign := <-c.postgreCommand.InsertOne("campaigns", &data)
	if insertCampaign.Error != nil {
		errObj := httpError.NewConflict()
		errObj.Message = "pembuatan campaign gagal"
		result.Error = errObj
		return result
	}

	result.Data = data
	return result
}

func (c commandUsecase) Update(ctx context.Context, payload *models.UpdateCampaign) utils.Result {
	var result utils.Result

	var query string
	parameter := make(map[string]interface{})

	query = "c.id = @id"
	parameter["id"] = payload.ID

	currentTime := time.Now()
	payload.UpdatedAt = currentTime
	payload.UpdatedBy = payload.Opts.UserID

	queryPayload := commands.CommandPayload{
		Table:     "campaigns c",
		Query:     query,
		Parameter: parameter,
		Document:  payload,
	}

	c.postgreCommand.Update(&queryPayload)

	result.Data = payload
	return result
}
func (c commandUsecase) HttpRequest(method string, url string, auth string, requestBody io.Reader) utils.Result {
	var result utils.Result

	request, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		errObj := httpError.NewConflict()
		errObj.Message = "URL tidak ditemukan"
		result.Error = errObj
		return result
	}

	request.Header.Add("Authorization", auth)
	var client = &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		errObj := httpError.NewConflict()
		errObj.Message = "Gagal fetch URL"
		result.Error = errObj
		return result
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var http struct {
		Success bool        `json:"success" bson:"success" default:"false"`
		Data    interface{} `json:"data" bson:"data"`
		Message string      `json:"message" bson:"message"`
		Code    int         `json:"code" bson:"code"`
	}
	json.Unmarshal([]byte(body), &http)

	result.Data = http.Data
	result.Error = http.Success
	return result
}
func (c commandUsecase) UploadCampaignImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, payload *models.UploadImageRequest) utils.Result {
	var result utils.Result

	res := <-c.postgreQuery.FindOneByID(ctx, payload.CampaignID)
	campaignImages := res.Data.(models.CampaignImages)
	if campaignImages.ID == "" {
		errObj := httpError.NewNotFound()
		errObj.Message = "campaign images tidak ditemukan"
		result.Error = errObj
		return result
	}

	campaignImages.FileName = "images/" + header.Filename
	ext := filepath.Ext(campaignImages.FileName)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		errObj := httpError.NewBadRequest()
		errObj.Message = "Format file tidak valid, format yang valid: jpg, png, jpeg"
		result.Error = errObj
		return result
	}

	isPrimary := 0
	if payload.IsPrimary {
		isPrimary = 1

		imagesPrimary := <-c.postgreCommand.MarkAllImagesAsNonPrimary(payload.CampaignID)
		if imagesPrimary.Error != nil {
			errObj := httpError.NewConflict()
			errObj.Message = "perubahan is_primary gagal"
			result.Error = errObj
			return result
		}
	}

	campaignImage := models.CampaignImages{
		ID:         uuid.NewString(),
		CampaignID: payload.CampaignID,
		IsPrimary:  isPrimary,
		FileName:   campaignImages.FileName,
	}

	newCampaignImages := <-c.postgreCommand.UploadImages(campaignImage.CampaignID, &campaignImage)
	if newCampaignImages.Error != nil {
		errObj := httpError.NewConflict()
		errObj.Message = "upload images campaign gagal"
		result.Error = errObj
		return result
	}

	result.Data = campaignImage
	return result
}
