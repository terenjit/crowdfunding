package main

import (
	"crowdfunding/config"
	"crowdfunding/pkg/utils"
	"fmt"
	"net/http"

	campaignHTTPHandler "crowdfunding/modules/campaigns/handlers"
	userHTTPHandler "crowdfunding/modules/users/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	// Echo instance
	e := echo.New()
	e.Validator = utils.NewValidationUtil()
	e.Use(middleware.CORS())

	e.GET("users", func(c echo.Context) error {
		return c.String(http.StatusOK, "This service is running properly")
	})

	crowdfundingGroup := e.Group("/crowdfunding")

	//initiate user http handler
	userHTTP := userHTTPHandler.New()
	userHTTP.Mount(crowdfundingGroup)

	//initiate campaign http handler
	campaignsHTTP := campaignHTTPHandler.New()
	campaignsHTTP.Mount(crowdfundingGroup)

	listenerPort := fmt.Sprintf("localhost:%d", config.GlobalEnv.HTTPPort)
	e.Logger.Fatal(e.Start(listenerPort))
}
