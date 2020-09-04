package nodes

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/coolray-dev/raydash/modules/utils"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// AccessToken receive a id from request url and return the access token of the specific id
//
// AccessToken godoc
// @Summary Node AccessToken
// @Description Node AccessToken according to nid
// @ID Nodes.AccessToken
// @Security ApiKeyAuth
// @Tags Nodes
// @Accept  json
// @Produce  json
// @Param nid path uint true "Node ID"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} accessTokenResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /nodes/{nid}/token [get]
func AccessToken(c *gin.Context) {

	// Get Node ID
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

	c.JSON(http.StatusOK, accessTokenResponse{
		AccessToken: node.AccessToken,
	})
	return
}

// GenerateToken receive a id from request url and return the access token of the specific id
//
// GenerateToken godoc
// @Summary Generate Node AccessToken
// @Description Generate Node AccessToken according to nid
// @ID Nodes.GenerateToken
// @Security ApiKeyAuth
// @Tags Nodes
// @Accept  plain
// @Produce  json
// @Param nid path uint true "Node ID"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} accessTokenResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /nodes/{nid}/token [post]
func GenerateToken(c *gin.Context) {

	// Get Node ID
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

	node.AccessToken = utils.RandString(64)

	// Save to DB
	if err := orm.DB.Save(&node).Error; err != nil {
		log.Log.WithError(err).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, accessTokenResponse{
		AccessToken: node.AccessToken,
	})
	return
}

func parseNID(c *gin.Context) (nid uint64, err error) {
	nid, err = strconv.ParseUint(c.Param("nid"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid NID: %w", err)
	}
	return
}
