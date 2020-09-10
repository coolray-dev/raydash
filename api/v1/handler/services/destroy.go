package services

import (
	"net/http"

	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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
