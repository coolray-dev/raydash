package nodes

import (
	"errors"
	"net/http"

	"github.com/coolray-dev/raydash/modules/utils"

	"github.com/coolray-dev/raydash/api/v1/handler"
	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type usersResponse struct {
	Total uint
	Users []models.User `json:"users"`
}

// Users receive a id from request url and return all the Users a node has
//
// Users godoc
// @Summary Node Users
// @Description Show Users of a node
// @ID Nodes.Users
// @Security ApiKeyAuth
// @Tags Nodes
// @Accept  json
// @Produce  json
// @Param nid path uint true "Node ID"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} usersResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /nodes/{nid}/users [get]
func Users(c *gin.Context) {

	// Get Node ID
	nid, err := parseNID(c)
	if err != nil {
		log.Log.WithError(err).Warn("Error Getting Node ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var node models.Node

	// Find node
	if err := orm.DB.Preload("Services").First(&node, nid).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		log.Log.WithField("nodeID", nid).Warn("Node Not Found")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		log.Log.WithError(err).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get all users' id
	var userList []uint64
	for _, s := range node.Services {
		userList = append(userList, s.UserID)
	}
	userList = utils.UInt64SliceDeDuplicate(userList)

	// Get users
	var users []models.User
	if err := orm.DB.Find(&users, userList).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Log.WithError(err).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, usersResponse{
		Total: uint(len(users)),
		Users: users,
	})
	return
}

type trafficResponse struct {
	User models.User `json:"user"`
}

type trafficRequest struct {
	CurrentTraffic int64 `json:"current_traffic"`
	MaxTraffic     int64 `json:"max_traffic"`
}

// Traffic receive traffic info and update it
//
// Traffic godoc
// @Summary Update user traffic
// @Description Update user traffic
// @ID nodes.Traffic
// @Security ApiKeyAuth
// @Tags Nodes
// @Accept  json
// @Produce  json
// @Param traffic body trafficRequest true "Traffic Info Object"
// @Param nid path uint true "Node ID"
// @Param username path string true "Username"
// @Param Authorization header string false "Node Token"
// @Success 200 {object} trafficResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /nodes/{nid}/users/{username}/traffic [patch]
func Traffic(c *gin.Context) {
	var user models.User
	username := c.Param("username")
	if err := orm.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, &handler.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var json trafficRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, &handler.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	if json.CurrentTraffic != 0 {
		user.CurrentTraffic = json.CurrentTraffic
	}
	if json.MaxTraffic != 0 {
		user.MaxTraffic = json.MaxTraffic
	}

	if err := orm.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, &handler.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &trafficResponse{
		User: user,
	})
	return
}
