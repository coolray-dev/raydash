package testutils

import (
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/jwt"
)

func SignAccessToken(user *models.User) string {
	token, err := jwt.SignAccessToken(user)

	if err != nil {
		panic(err)
	}

	return token
}

func SignRefreshToken(user *models.User) string {
	token, err := jwt.SignRefreshToken(user)

	if err != nil {
		panic(err)
	}

	return token
}
