package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coolray-dev/raydash/api/v1/handler"
	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
)

type servicesResponse struct {
	Services []*model.Service
}

// Services returns all services a user has
//
// Services godoc
// @Summary List all services
// @Description Return a list of services of a user
// @ID users.Services
// @Tags Users
// @Accept  json
// @Produce  json
// @Param username path string true "Username"
// @Param Authorization header string false "Node Token"
// @Success 200 {object} servicesResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /users/{username}/services [get]
func Services(c *gin.Context) {
	var user model.User
	username := c.Param("username")
	if err := orm.DB.Preload("Groups").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, &handler.ErrorResponse{Error: err.Error()})
		return
	}
	var nodes []*model.Node
	for _, g := range user.Groups {
		if err := orm.DB.Preload("Nodes").Where("ID = ?", g.ID).First(&g).Error; err != nil {
			c.JSON(http.StatusInternalServerError, &handler.ErrorResponse{Error: err.Error()})
			return
		}
		nodes = append(nodes, g.Nodes...)

	}
	var services []*model.Service
	for _, n := range nodes {
		if err := orm.DB.Preload("Services").Where("ID = ?", n.ID).First(&n).Error; err != nil {
			c.JSON(http.StatusInternalServerError, &handler.ErrorResponse{Error: err.Error()})
			return
		}
		services = append(services, n.Services...)
	}
	c.JSON(http.StatusOK, &servicesResponse{
		Services: services,
	})
	return
}
