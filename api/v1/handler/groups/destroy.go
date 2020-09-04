package groups

import (
	"net/http"
	"strconv"

	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/casbin"
	"github.com/gin-gonic/gin"
)

type destroyResponse struct {
	Group string `json:"group"`
}

// Destroy receive a id from request and delete in from DB
//
// Destroy godoc
// @Summary Destroy Group
// @Description Destroy Group according to gid
// @ID groups.Destroy
// @Security ApiKeyAuth
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param gid path uint true "Group ID"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} destroyResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /groups/{gid} [delete]
func Destroy(c *gin.Context) {
	gid, err := parseGID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var group model.Group
	group.ID = gid
	if err = orm.DB.Where("id = ?", gid).First(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	oldname := group.Name
	if err = orm.DB.Delete(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Remove Group from casbin
	casbin.Enforcer.RemovePolicy("group::"+oldname, "/*/groups/"+strconv.Itoa(int(group.ID))+"*", "*")

	c.JSON(http.StatusOK, destroyResponse{
		Group: "",
	})
	return
}
