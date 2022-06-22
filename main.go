package main

import (
	"crowdfunding/config"
	"fmt"
	"net/http"

	userHTTPHandler "crowdfunding/modules/users/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	// Echo instance
	e := echo.New()

	e.Use(middleware.CORS())

	e.GET("users", func(c echo.Context) error {
		return c.String(http.StatusOK, "This service is running properly")
	})

	userGroup := e.Group("/users")

	//initiate user http handler
	userHTTP := userHTTPHandler.New()
	userHTTP.Mount(userGroup)

	listenerPort := fmt.Sprintf("localhost:%d", config.GlobalEnv.HTTPPort)
	e.Logger.Fatal(e.Start(listenerPort))
}
