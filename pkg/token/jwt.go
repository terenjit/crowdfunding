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
