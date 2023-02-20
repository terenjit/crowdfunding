package middleware

import (
	"crowdfunding/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func VerifyBasicAuth() echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, context echo.Context) (bool, error) {
		if username == config.GlobalEnv.BasicAuthUsername && password == config.GlobalEnv.BasicAuthPassword {
			return true, nil
		}

		return false, nil
	})
}
