package services

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//Index list out all services and return
func Index(c *gin.Context) {
	var services []model.Service
	uid, uexists := c.Get("uid")
	nid, nexists := c.Get("nid")
	switch {
	case uexists:
		{
			var user model.User
			user.ID = uid.(uint64)
			if err := orm.DB.Model(&user).Association("Services").Find(&services); errors.Is(err, gorm.ErrRecordNotFound) {
				log.Log.WithFields(logrus.Fields{
					"uid": uid,
				}).Info("No Services in Database")
				c.Status(http.StatusNotFound)
				return
			} else if err != nil {
				log.Log.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("Database Error")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	case nexists:
		{
			var node model.Node
			node.ID = nid.(uint64)
			if err := orm.DB.Model(&node).Association("Services").Find(&services); errors.Is(err, gorm.ErrRecordNotFound) {
				log.Log.WithFields(logrus.Fields{
					"nid": nid,
				}).Info("No Services in Database")
				c.Status(http.StatusNotFound)
				return
			} else if err != nil {
				log.Log.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("Database Error")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	default:
		{
			if err := orm.DB.Find(&services).Error; errors.Is(err, gorm.ErrRecordNotFound) {
				log.Log.Info("No Services in Database")
				c.Status(http.StatusNotFound)
				return
			} else if err != nil {
				log.Log.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("Database Error")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	for i, s := range services {
		var node model.Node
		if err := orm.DB.Where("id = ?", s.NodeID).Find(&node).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Log.WithFields(logrus.Fields{
				"nid": s.NodeID,
			}).Warning("No Such Node in Database")
			continue
		} else if err != nil {
			log.Log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Database Error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if node.HasMultiPort {

		} else {
			var user model.User
			if err2 := orm.DB.Where("id = ?", s.UserID).Find(&user).Error; errors.Is(err2, gorm.ErrRecordNotFound) {
				log.Log.WithFields(logrus.Fields{
					"uid": s.UserID,
				}).Warning("No Such User in Database")
				continue
			} else if err2 != nil {
				log.Log.WithFields(logrus.Fields{
					"error": err2.Error(),
				}).Error("Database Error")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
				return
			}
			services[i].Protocol = "vmess"
			services[i].Host = node.Host
			services[i].Port = node.Settings.Port
			services[i].VmessSetting = node.Settings.VmessSetting
			services[i].VmessUser = model.VmessUser{
				Email:    user.Email,
				UUID:     user.UUID,
				AlterID:  64,
				Security: "auto",
			}
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"total":    len(services),
		"services": services,
	})
	return
}

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

// Update receive a id and a service object from request and update the specific record in DB
func Update(c *gin.Context) {
	sid, err := parseSID(c)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("URL Param Invalid")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var service model.Service
	service.ID = sid
	if err = c.ShouldBindJSON(&service); err != nil {
		log.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Request Binding Error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = orm.DB.Save(&service).Error; err != nil {
		log.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"service": service,
	})
	return
}

// Destroy delete a service from db
func Destroy(c *gin.Context) {
	sid, err := parseSID(c)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("URL Param Invalid")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var service model.Service
	service.ID = sid

	if err = orm.DB.Delete(&service).Error; err != nil {
		log.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"node": "",
	})
	return
}

func parseSID(c *gin.Context) (sid uint64, err error) {
	sid, err = strconv.ParseUint(c.Param("sid"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid SID: %w", err)
	}
	return
}
