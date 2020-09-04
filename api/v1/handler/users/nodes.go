package users

import (
	"net/http"

	"github.com/coolray-dev/raydash/api/v1/handler"
	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
	"github.com/gin-gonic/gin"
)

type nodesResponse struct {
	Nodes []*model.Node `json:"nodes"`
}

// Nodes return all nodes a user has
//
// Nodes godoc
// @Summary List all nodes
// @Description Return a list of nodes of a user
// @ID users.Nodes
// @Security ApiKeyAuth
// @Tags Users
// @Accept  json
// @Produce  json
// @Param username path string true "Username"
// @Param Authorization header string false "Node Token"
// @Success 200 {object} nodesResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /users/{username}/nodes [get]
func Nodes(c *gin.Context) {
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
	c.JSON(http.StatusOK, &nodesResponse{
		Nodes: nodes,
	})
	return
}
