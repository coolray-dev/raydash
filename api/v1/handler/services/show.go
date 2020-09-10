package services

import (
	"errors"
	"net/http"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type showResponse struct {
	Service models.Service `json:"service"`
}

func Show(c *gin.Context) {
	sid, err := parseSID(c)

	if err != nil {
		log.Log.WithError(err).Warn("Error Getting Service ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var service models.Service
	service.ID = sid

	if err := orm.DB.First(&service).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		log.Log.WithField("serviceID", sid).Warn("Service Not Found")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		log.Log.WithError(err).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, showResponse{
		Service: service,
	})
}
