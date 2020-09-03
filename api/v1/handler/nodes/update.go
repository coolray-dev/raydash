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

type updateResponse struct {
	Node models.Node `json:"node"`
}

// Update receive a id and a node object from request and update the specific record in DB
//
// Update godoc
// @Summary Update Node
// @Description Update a Node
// @ID Nodes.Update
// @Security ApiKeyAuth
// @Tags Nodes
// @Accept  json
// @Produce  json
// @Param nid path uint true "Node ID"
// @Param node body models.Node true "Node Object"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} updateResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /nodes/{nid} [patch]
func Update(c *gin.Context) {

	nid, err := parseNID(c)
	if err != nil {
		log.Log.WithError(err).Warn("Error Getting Node ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var node models.Node

	if err := orm.DB.First(&node, nid).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		log.Log.WithField("nodeID", nid).Warn("Node Not Found")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		log.Log.WithError(err).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Bind Request
	if err = c.ShouldBindJSON(&node); err != nil {
		log.Log.WithError(err).Warn("Error Binding Request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save to DB
	if err = orm.DB.Save(&node).Error; err != nil {
		log.Log.WithError(err).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, updateResponse{
		Node: node,
	})
	return
}
