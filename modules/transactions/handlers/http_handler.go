package handlers

import (
	"context"
	models "crowdfunding/modules/transactions/models/domain"
	"crowdfunding/modules/transactions/repositories/commands"
	"crowdfunding/modules/transactions/repositories/queries"
	"crowdfunding/modules/transactions/usecases"
	database "crowdfunding/pkg/databases"
	"crowdfunding/pkg/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	queryUsecase   usecases.QueryUsecase
	commandUsecase usecases.CommandUsecase
}

func New() *HTTPHandler {
	postgreDb := database.Initpostgre(context.Background())

	PostgreQuery := queries.NewPostgreQuery(postgreDb)
	QueryUsecase := usecases.NewQueryUsecase(PostgreQuery)

	postgreCommand := commands.NewPostgreCommand(postgreDb)
	commandUsecase := usecases.NewCommandUsecasePostgre(PostgreQuery, postgreCommand)
	return &HTTPHandler{
		queryUsecase:   QueryUsecase,
		commandUsecase: commandUsecase,
	}
}

func (h *HTTPHandler) Mount(echoGroup *echo.Group) {
	echoGroup.GET("/v1/campaigns/:campaign_id/transactions", h.ListofTransactions)
}

func (h *HTTPHandler) ListofTransactions(c echo.Context) error {
	campaignid := c.Param("campaign_id")

	query := make(map[string]interface{})

	for key, value := range c.QueryParams() {
		if key != "page" && key != "quantity" {
			query[key] = value[0]
		} else {
			query[key] = value[0]
			v, err := strconv.Atoi(value[0])
			if err == nil {
				query[key] = v
			}
		}
	}

	payload, _ := json.Marshal(query)
	header, _ := json.Marshal(c.Get("opts"))
	var data models.TransactionList
	json.Unmarshal(payload, &data)
	json.Unmarshal(header, &data.Opts)
	data.CampaignID = campaignid

	result := h.queryUsecase.ListTransactions(c.Request().Context(), &data)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.PaginationResponse(result.Data, result.MetaData, "List All Campaigns", http.StatusOK, c)
}