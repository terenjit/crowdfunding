package token

import (
	"context"
	"crowdfunding/pkg/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secret_key = []byte("CrowdFunding_s3cr3t_k3y")

func Generate(ctx context.Context, UserId string, expired time.Duration) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		now := time.Now()
		exp := now.Add(expired)

		claims := jwt.MapClaims{
			"exp": exp.Unix(),
			"iat": now.Unix(),
			"id":  UserId,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(secret_key)
		if err != nil {
			output <- utils.Result{Error: err}
			return
		}

		output <- utils.Result{Data: tokenString}

	}()

	return output
}

func Validate(ctx context.Context, tokenString string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		tokenParse, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret_key), nil
		})

		var errToken string
		switch ve := err.(type) {
		case *jwt.ValidationError:
			if ve.Errors == jwt.ValidationErrorExpired {
				errToken = "token has been expired"
			} else {
				errToken = "token parsing error"
			}
		}

		if len(errToken) > 0 {
			output <- utils.Result{Error: errToken}
			return
		}

		if !tokenParse.Valid {
			output <- utils.Result{Error: "token parsing error"}
			return
		}

		mapClaims, _ := tokenParse.Claims.(jwt.MapClaims)

		tokenClaim := Claim{
			UserID: mapClaims["userId"].(string),
			Key:    mapClaims["key"].(string),
		}

		if mapClaims["rt"] != nil {
			tokenClaim.RefreshToken = mapClaims["rt"].(string)
		}

		output <- utils.Result{Data: tokenClaim}
	}()

	return output
}
