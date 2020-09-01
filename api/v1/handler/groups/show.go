package groups

import (
	"errors"
	"net/http"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type showResponse struct {
	Group models.Group
}

// Show receive a id from request url and return the group of the specific id
//
// Show godoc
// @Summary Show Group
// @Description Show Group according to gid
// @ID groups.Show
// @Security ApiKeyAuth
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param gid path uint true "Group ID"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} showResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /groups/{gid} [get]
func Show(c *gin.Context) {
	gid, err := parseGID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var group models.Group
	group.ID = gid

	if err := orm.DB.First(&group).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, showResponse{
		Group: group,
	})
	return
}
