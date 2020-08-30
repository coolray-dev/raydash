package users

import (
	"net/http"

	"github.com/coolray-dev/raydash/api/v1/handler"
	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
	"github.com/gin-gonic/gin"
)

type groupsResponse struct {
	Groups []*model.Group
}

// Groups shows all groups a user has
//
// Groups godoc
// @Summary List all groups
// @Description Return a list of groups of a user
// @ID users.Groups
// @Tags Users
// @Accept  json
// @Produce  json
// @Param username path string true "Username"
// @Param Authorization header string false "Node Token"
// @Success 200 {object} groupsResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /users/{username}/groups [get]
func Groups(c *gin.Context) {
	var user model.User
	username := c.Param("username")
	if err := orm.DB.Preload("Groups").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, &handler.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, &groupsResponse{
		Groups: user.Groups,
	})
	return
}
