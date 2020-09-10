package services

import (
	"errors"
	"net/http"

	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Store recieve a service object and store it in DB
func Store(c *gin.Context) {
	type Request struct {
		Name        string                   `json:"name" `
		Description string                   `json:"description"`
		Host        string                   `json:"host"`
		Port        uint                     `json:"port"`
		Protocol    string                   `json:"protocol"`
		NID         uint64                   `json:"nid" binding:"required"`
		UID         uint64                   `json:"uid" binding:"required"`
		VS          model.VmessSetting       `json:"vmessSettings"`
		SS          model.ShadowsocksSetting `json:"shadowsocksSettings"`
	}
	var json Request
	if err := c.ShouldBindJSON(&json); err != nil {
		log.Log.Error("Request Binding Error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var service model.Service
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
