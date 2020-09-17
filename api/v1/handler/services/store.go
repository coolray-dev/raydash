package services

import (
	"errors"
	"net/http"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	model "github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type storeRequest struct {
	Name        string                    `json:"name" `
	Description string                    `json:"description"`
	Host        string                    `json:"host"`
	Port        uint                      `json:"port"`
	Protocol    string                    `json:"protocol"`
	NID         uint64                    `json:"nid" binding:"required"`
	UID         uint64                    `json:"uid" binding:"required"`
	VS          models.VmessSetting       `json:"vmessSettings"`
	SS          models.ShadowsocksSetting `json:"shadowsocksSettings"`
}

// Store recieve a service object and store it in DB
//
// Store godoc
// @Summary Create Service
// @Description Create a service from post json object
// @ID Services.Store
// @Security ApiKeyAuth
// @Tags Services
// @Accept  json
// @Produce  json
// @Param service body storeRequest true "Service Object"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} serviceResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /services [post]
func Store(c *gin.Context) {

	var json storeRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		log.Log.Error("Request Binding Error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var service models.Service
	service.Name = json.Name
	service.Description = json.Description
	service.NodeID = json.NID
	service.UserID = json.UID
	var node model.Node
	if err := orm.DB.Where("id = ?", json.NID).Find(&node).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		log.Log.WithFields(logrus.Fields{
			"nid": json.NID,
		}).Warning("No Such Node in Database")
		c.Status(http.StatusNotFound)
		return
	} else if err != nil {
		log.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if node.HasMultiPort {
		service.Host = json.Host
		service.Port = json.Port
		service.Protocol = json.Protocol
		service.VmessSetting = json.VS
		service.ShadowsocksSetting = json.SS
	} else {
		var user model.User
		if err2 := orm.DB.Where("id = ?", json.UID).Find(&user).Error; errors.Is(err2, gorm.ErrRecordNotFound) {
			log.Log.WithFields(logrus.Fields{
				"uid": json.UID,
			}).Warning("No Such User in Database")
			c.Status(http.StatusNotFound)
			return
		} else if err2 != nil {
			log.Log.WithFields(logrus.Fields{
				"error": err2.Error(),
			}).Error("Database Error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
			return
		}
		service.Protocol = "vmess"
		service.Host = node.Host
		service.Port = node.Settings.Port
		service.VmessSetting = node.Settings.VmessSetting
		service.VmessUser = model.VmessUser{
			Email:    user.Email,
			UUID:     user.UUID,
			AlterID:  64,
			Security: "auto",
		}
	}

	if err := orm.DB.Create(&service).Error; err != nil {
		log.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"service": service,
	})
	return
}
