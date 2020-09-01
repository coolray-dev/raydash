package groups

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/casbin"
)

type createResponse struct {
	Group models.Group
}

// Create receive a group object from request and update the specific record in DB
//
// Create godoc
// @Summary Create Group
// @Description Create a group
// @ID groups.Create
// @Security ApiKeyAuth
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param group body models.Group true "Group Object"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} createResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /groups [post]
func Create(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := orm.DB.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Deal With Access Control
	casbin.Enforcer.AddPolicy("group::"+group.Name, "/*/groups/"+strconv.Itoa(int(group.ID))+"*", "*")

	c.JSON(http.StatusOK, createResponse{
		Group: group,
	})
	return
}
