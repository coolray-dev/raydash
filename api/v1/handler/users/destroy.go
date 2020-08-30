package users

import (
	"net/http"

	"github.com/coolray-dev/raydash/api/v1/handler"
	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
	"github.com/gin-gonic/gin"
)

type destroyResponse struct {
	User string
}

// Destroy delete a user from db
//
// Destroy godoc
// @Summary Delete a user
// @Description Delete a user according to username
// @ID users.Destroy
// @Tags Users
// @Accept  json
// @Produce  json
// @Param username path string true "Username"
// @Param Authorization header string false "Node Token"
// @Success 200 {object} destroyResponse
// @Failure 500 {object} handler.ErrorResponse
// @Failure 403 {object} handler.ErrorResponse
// @Router /users/{username} [delete]
func Destroy(c *gin.Context) {
	username := c.Param("username")
	tokenUsername := c.MustGet("username").(string)
	isAdmin := c.MustGet("isAdmin").(bool)
	if (username != tokenUsername) && !isAdmin {
		c.JSON(http.StatusForbidden, &handler.ErrorResponse{Error: "No permission"})
		return
	}
	if err := orm.DB.Where("username = ?", username).Delete(model.User{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, &handler.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, &destroyResponse{
		User: "",
	})
}
