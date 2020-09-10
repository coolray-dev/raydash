package services

import (
	"net/http"

	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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
