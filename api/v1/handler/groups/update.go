package groups

import (
	"net/http"
	"strconv"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/casbin"
	"github.com/gin-gonic/gin"
)

type updateResponse struct {
	Group models.Group
}

// Update receive a id and a group object from request and update the specific record in DB
//
// Update godoc
// @Summary Update Group
// @Description Update a group
// @ID groups.Update
// @Security ApiKeyAuth
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param gid path uint true "Group ID"
// @Param group body models.Group true "Group Object"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} updateResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /groups/{gid} [patch]
func Update(c *gin.Context) {
	gid, err := parseGID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var group models.Group
	group.ID = gid
	if err = orm.DB.Where("id = ?", gid).First(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	oldname := group.Name
	if err = c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = orm.DB.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Remove old policy and add new policy
	if oldname != group.Name {
		casbin.Enforcer.RemovePolicy("group::"+oldname, "/*/groups/"+strconv.Itoa(int(group.ID))+"*", "*")
		casbin.Enforcer.AddPolicy("group::"+group.Name, "/*/groups/"+strconv.Itoa(int(group.ID))+"*", "*")
	}

	c.JSON(http.StatusOK, updateResponse{
		Group: group,
	})
	return
}
