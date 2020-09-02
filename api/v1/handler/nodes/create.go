package nodes

import (
	"net/http"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
)

type createResponse struct {
	Node models.Node
}

// Create receive a id and a node object from request and update the specific record in DB
//
// Create godoc
// @Summary Create Node
// @Description Create node from post json object
// @ID Nodes.Create
// @Security ApiKeyAuth
// @Tags Nodes
// @Accept  json
// @Produce  json
// @Param node body models.Node true "Node Object"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} createResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /nodes [post]
func Create(c *gin.Context) {

	// Bind Request
	var node models.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		log.Log.WithError(err).Warning("Could not bind request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save to DB
	if err := orm.DB.Save(&node).Error; err != nil {
		log.Log.WithError(err).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return result
	log.Log.Debug("Success")
	c.JSON(http.StatusOK, createResponse{
		Node: node,
	})
	return
}
