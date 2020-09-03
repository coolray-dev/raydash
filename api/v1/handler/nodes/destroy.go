package nodes

import (
	"net/http"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
)

type destroyResponse struct {
	Node string `json:"node"`
}

// Destroy receive a id from request and delete in from DB
//
// Destroy godoc
// @Summary Destroy Node
// @Description Destroy Node according to nid
// @ID Nodes.Destroy
// @Security ApiKeyAuth
// @Tags Nodes
// @Accept  json
// @Produce  json
// @Param nid path uint true "Node ID"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} destroyResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /nodes/{nid} [delete]
func Destroy(c *gin.Context) {

	nid, err := parseNID(c)
	if err != nil {
		log.Log.WithError(err).Warn("Error Getting Node ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var node models.Node
	node.ID = nid

	if err = orm.DB.Delete(&node).Error; err != nil {
		log.Log.WithError(err).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, destroyResponse{
		Node: "",
	})
	return
}
