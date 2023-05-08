package handlers

import (
	"context"
	"crowdfunding/middleware"
	models "crowdfunding/modules/campaigns/models/domain"
	"crowdfunding/modules/campaigns/repositories/commands"
	"crowdfunding/modules/campaigns/repositories/queries"
	"crowdfunding/modules/campaigns/usecases"
	database "crowdfunding/pkg/databases"
	"crowdfunding/pkg/utils"
	"encoding/json"
	"log"
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
	echoGroup.GET("/v1/campaigns", h.getList, middleware.VerifyBearer())
	echoGroup.GET("/v1/campaigns/:id", h.getDetail, middleware.VerifyBearer())
	echoGroup.PUT("/v1/campaigns/:id", h.Update, middleware.VerifyBearer())

}

func (h *HTTPHandler) getList(c echo.Context) error {

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
	var data models.CampainGetList
	json.Unmarshal(payload, &data)
	json.Unmarshal(header, &data.Opts)

	result := h.queryUsecase.GetList(c.Request().Context(), &data)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.PaginationResponse(result.Data, result.MetaData, "List All Campaigns", http.StatusOK, c)
}

func (h *HTTPHandler) getDetail(c echo.Context) error {
	id := c.Param("id")

	result := h.queryUsecase.GetDetail(c.Request().Context(), id)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Detail Campaign", http.StatusOK, c)
}

func (h *HTTPHandler) Update(c echo.Context) error {
	var data = new(models.UpdateCampaign)

	if err := utils.BindValidate(c, data); err != nil {
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}

	result := h.commandUsecase.Update(c.Request().Context(), data)
	if result.Error != nil {
		log.Println(result.Error)
		return utils.ResponseError(result.Error, c)
	}
	return utils.Response(result.Data, "User success updated", http.StatusOK, c)
}
