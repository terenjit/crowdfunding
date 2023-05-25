package usecases

import (
	"context"
	models "crowdfunding/modules/transactions/models/domain"
	"crowdfunding/modules/transactions/repositories/commands"
	"crowdfunding/modules/transactions/repositories/queries"
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
	midtrans "github.com/veritrans/go-midtrans"
)

type commandUsecase struct {
	postgreQuery   queries.TransactionsPostgre
	postgreCommand commands.CommandPostgre
}

func NewCommandUsecasePostgre(postgreQuery queries.TransactionsPostgre, postgreCommand commands.CommandPostgre) *commandUsecase {
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

	transactionId := uuid.NewString()

	data := models.TransactionModel{
		ID:         transactionId,
		CampaignID: payload.CampaignID,
		UserID:     payload.Opts.UserID,
		Amount:     payload.Amount,
		Status:     "pending",
		Code:       transactionId,
		CreatedAt:  time.Now(),
	}

	insertCampaign := <-c.postgreCommand.InsertOne("transactions", &data)
	if insertCampaign.Error != nil {
		errObj := httpError.NewConflict()
		errObj.Message = "pembuatan campaign gagal"
		result.Error = errObj
		return result
	}

	paymentURL, err := c.GetPaymentURL(&data)
	if err != nil {
		errObj := httpError.NewConflict()
		errObj.Message = "pembuatan payment url gagal di create campaign"
		result.Error = errObj
		return result
	}

	data.PaymentURL = paymentURL

	insertCampaign = <-c.postgreCommand.Update("transactions", &data)
	if insertCampaign.Error != nil {
		errObj := httpError.NewConflict()
		errObj.Message = "update paymentURL gagal"
		result.Error = errObj
		return result
	}

	result.Data = data
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

func (c commandUsecase) GetPaymentURL(payload *models.TransactionModel) (string, error) {
	var result utils.Result

	requestUser := c.HttpRequest(http.MethodGet, fmt.Sprintf("localhost:9000/crowdfunding/v1/users/%s", payload.Opts.UserID), payload.Opts.Authorization, nil)
	if requestUser.Error == false {
		errObj := httpError.NewConflict()
		errObj.Message = "Data pengguna tidak ditemukan"
		result.Error = errObj
		return "", nil
	}

	midclient := midtrans.NewClient()
	midclient.ServerKey = ""
	midclient.ClientKey = ""
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: payload.Opts.UserID,
			FName: payload.Opts.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  payload.ID,
			GrossAmt: payload.Amount,
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		errObj := httpError.NewConflict()
		errObj.Message = "pembuatan payment URL gagal"
	}

	return snapTokenResp.RedirectURL, nil
}

func (c commandUsecase) ProcessPayment(ctx context.Context, payload *models.TransactionNotificationInput) error {
	transaction_id := payload.OrderID

	transaction, err := c.postgreCommand.FindByID(transaction_id)
	if err != nil {
		return err
	}

	if payload.PaymentType == "credit_card" && payload.TransactionStatus == "capture" && payload.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if payload.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if payload.TransactionStatus == "deny" || payload.TransactionStatus == "expire" || payload.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := c.postgreCommand.UpdateTransaction(transaction)
	if err != nil {
		return err
	}

	campaign, err := c.postgreCommand.FindCampaignByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		updateCampaign := <-c.postgreCommand.Update("campaigns", campaign)
		if updateCampaign.Error != nil {
			errObj := httpError.NewConflict()
			errObj.Message = "update campaign terakhir gagal"
		}
	}

	return nil
}
