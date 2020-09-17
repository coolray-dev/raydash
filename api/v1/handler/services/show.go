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

// Show a Service
//
// Show godoc
// @Summary Show Service
// @Description Show a service according to nid
// @ID Services.Show
// @Security ApiKeyAuth
// @Tags Services
// @Accept  json
// @Produce  json
// @Param nid path uint true "Services ID"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} serviceResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /services/{nid} [get]
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

	c.JSON(http.StatusOK, serviceResponse{
		Service: service,
	})
}
