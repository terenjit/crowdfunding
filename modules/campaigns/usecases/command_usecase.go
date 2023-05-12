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
	"net/http"
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
	payload.UpdatedBy = payload.UserID

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
