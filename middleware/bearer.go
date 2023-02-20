package middleware

import (
	"crowdfunding/pkg/token"
	"crowdfunding/pkg/utils"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func VerifyBearer() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := strings.TrimPrefix(c.Request().Header.Get(echo.HeaderAuthorization), "Bearer ")

			if len(tokenString) == 0 {
				return utils.Response(nil, "Invalid token!", http.StatusUnauthorized, c)
			}

			parsedToken := <-token.Validate(c.Request().Context(), tokenString)
			if parsedToken.Error != nil {
				return utils.Response(nil, parsedToken.Error.(string), http.StatusUnauthorized, c)
			}

			parsedByte, err := json.Marshal(parsedToken.Data)
			if err != nil {
				return utils.Response(nil, err.Error(), http.StatusUnauthorized, c)
			}

			var claimToken token.Claim
			json.Unmarshal(parsedByte, &claimToken)

			c.Set("opts", claimToken)

			return next(c)
		}
	}
}
