package testutils

import (
	"github.com/coolray-dev/raydash/api/v1/handler/authentication"
	"github.com/coolray-dev/raydash/models"
)

func SignAccessToken(user *models.User) string {
	token, err := authentication.SignAccessToken(user)

	if err != nil {
		panic(err)
	}

	return token
}

func SignRefreshToken(user *models.User) string {
	token, err := authentication.SignRefreshToken(user)

	if err != nil {
		panic(err)
	}

	return token
}
