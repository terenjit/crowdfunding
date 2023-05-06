package handlers

import (
	"context"
	"crowdfunding/config"
	"crowdfunding/middleware"
	models "crowdfunding/modules/users/models/domain"
	"crowdfunding/modules/users/repositories/commands"
	"crowdfunding/modules/users/repositories/queries"
	"crowdfunding/modules/users/usecases"
	database "crowdfunding/pkg/databases"
	"crowdfunding/pkg/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// HTTP Handler
type HTTPHandler struct {
	commandUsecase usecases.CommandUsecase
	queryUsecase   usecases.QueryUsecase
}

// New initiation
func New() *HTTPHandler {

	postgreDb := database.Initpostgre(context.Background())
	utils.LoadValidatorSchemas(config.GlobalEnv.RootApp + "/modules/users/models/json-schema")

	postgreQuery := queries.NewUserPostgreQuery(postgreDb)
	postgreCommand := commands.NewUserPostgreCommand(postgreDb)
	queryUsecase := usecases.NewUserQueryUsecase(postgreQuery)
	commandUsecase := usecases.NewUserCommandUsecase(postgreQuery, postgreCommand, queryUsecase)

	return &HTTPHandler{
		commandUsecase: commandUsecase,
		queryUsecase:   queryUsecase,
	}
}

// Mount function
func (h *HTTPHandler) Mount(echoGroup *echo.Group) {
	echoGroup.POST("/v1/users/register", h.Register)
	echoGroup.POST("/v1/users/login", h.Login, middleware.VerifyBasicAuth())
	echoGroup.POST("/v1/users/avatars/:id", h.UploadAvatar, middleware.VerifyBearer())
	echoGroup.GET("/v1/users", h.getList)
	echoGroup.GET("/v1/users/:id", h.getDetail)

}

// Register Function
func (h *HTTPHandler) Register(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)

	if err := utils.ValidateDocument("register", body); err != nil {
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}

	var data models.Register
	json.Unmarshal(body, &data)

	result := h.commandUsecase.Register(c.Request().Context(), &data)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Register User", http.StatusCreated, c)
}

// Login Function
func (h *HTTPHandler) Login(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)

	if err := utils.ValidateDocument("login", body); err != nil {
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}

	var data models.LoginRequest
	json.Unmarshal(body, &data)

	result := h.commandUsecase.Login(c.Request().Context(), &data)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Login User", http.StatusOK, c)

}

func (h *HTTPHandler) UploadAvatar(c echo.Context) error {
	ID := c.Param("id")

	file, header, err := c.Request().FormFile("avatar")

	if err != nil {
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}

	result := h.commandUsecase.SaveAvatar(c.Request().Context(), file, header, ID)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Success update avatar", http.StatusOK, c)
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
	var data models.UsersGetList
	json.Unmarshal(payload, &data)

	result := h.queryUsecase.GetList(c.Request().Context(), &data)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.PaginationResponse(result.Data, result.MetaData, "List All Users", http.StatusOK, c)
}

func (h *HTTPHandler) getDetail(c echo.Context) error {
	id := c.Param("id")

	result := h.queryUsecase.GetDetail(c.Request().Context(), id)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Detail User", http.StatusOK, c)
}
