package users

import (
	"github.com/coolray-dev/raydash/models"
)

type userResponse struct {
	User models.User `json:"user"`
}
