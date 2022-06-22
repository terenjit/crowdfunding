package handlers

import (
	"context"
	"crowdfunding/config"
	models "crowdfunding/modules/users/models/domain"
	"crowdfunding/modules/users/repositories/commands"
	"crowdfunding/modules/users/repositories/queries"
	"crowdfunding/modules/users/usecases"
	database "crowdfunding/pkg/databases"
	"crowdfunding/pkg/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

//HTTP Handler
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
	// echoGroup.POST("/v1/users/register", h.Register, middleware.VerifyBasicAuth())
	// echoGroup.GET("/v1/users/:id", h.ViewProfile, middleware.VerifyBearer())
	// echoGroup.PUT("/v1/users/:id", h.UpdateProfile, middleware.VerifyBearer())
	// echoGroup.PUT("/v1/users/avatar/:id", h.UpdateAvatar, middleware.VerifyBearer())
	// echoGroup.POST("/v1/users/refresh-token", h.RefreshToken, middleware.VerifyBasicAuth())
	// echoGroup.POST("/v1/users/logout", h.Logout, middleware.VerifyBearer())
	// echoGroup.PUT("/v1/users/change-password", h.ChangePassword, middleware.VerifyBearer())
	// echoGroup.POST("/v1/users/otp", h.GenerateOTP, middleware.VerifyBasicAuth())
	// echoGroup.POST("/v1/users/verify-otp-code", h.VerifyOtpCode, middleware.VerifyBasicAuth())
	// echoGroup.POST("/v1/users/forgot-account", h.ForgotAccount, middleware.VerifyBasicAuth())
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
