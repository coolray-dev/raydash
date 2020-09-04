package users

import (
	"net/http"

	"github.com/coolray-dev/raydash/api/v1/handler"
	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
	"github.com/gin-gonic/gin"
)

type showResponse struct {
	User model.User `json:"user"`
}

// Show query a user and return it using url param "username" as condition
//
// Show godoc
// @Summary Show the required user
// @Description Return user according to username in url
// @ID users.Show
// @Security ApiKeyAuth
// @Tags Users
// @Accept  json
// @Produce  json
// @Param username path string true "Username"
// @Param Authorization header string false "Node Token"
// @Success 200 {object} showResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /users/{username} [get]
func Show(c *gin.Context) {
	var user model.User
	username := c.Param("username")
	if err := orm.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, &handler.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, &showResponse{
		User: user,
	})
	return
}
