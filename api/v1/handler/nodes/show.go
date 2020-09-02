package nodes

import (
	"errors"
	"net/http"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type showResponse struct {
	Node models.Node
}

// Show receive a id from request url and return the node of the specific id
//
// Show godoc
// @Summary Show Node
// @Description Show Node according to nid
// @ID Nodes.Show
// @Security ApiKeyAuth
// @Tags Nodes
// @Accept  json
// @Produce  json
// @Param nid path uint true "Node ID"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} showResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /nodes/{nid} [get]
func Show(c *gin.Context) {

	// Get Node ID
	nid, err := parseNID(c)
	if err != nil {
		log.Log.WithError(err).Warn("Error Getting Node ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var node models.Node
	node.ID = nid

	if err := orm.DB.First(&node).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		log.Log.WithField("nodeID", nid).Warn("Node Not Found")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		log.Log.WithError(err).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, showResponse{
		Node: node,
	})
	return
}
