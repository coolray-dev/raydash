package services

import (
	"errors"
	"net/http"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type indexResponse struct {
	Total    uint
	Services []models.Service `json:"services"`
}

// Index list out all services and return
//
// Index godoc
// @Summary All Services
// @Description Simply list out all Services
// @ID Services.Index
// @Security ApiKeyAuth
// @Tags Services
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Access Token"
// @Success 200 {object} indexResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /services [get]
func Index(c *gin.Context) {
	var services []models.Service
	uid, uexists := c.Get("uid")
	nid, nexists := c.Get("nid")
	switch {
	case uexists:
		{
			var user models.User
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
			var node models.Node
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
		var node models.Node
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
			var user models.User
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
			services[i].VmessUser = models.VmessUser{
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
