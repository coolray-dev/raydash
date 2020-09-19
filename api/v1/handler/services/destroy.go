package services

import (
	"net/http"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
)

type destroyResponse struct {
	Service string `json:"service"`
}

// Destroy delete a service from db
//
// Destroy godoc
// @Summary Destroy Service
// @Description Destroy a service according to sid
// @ID Services.Destroy
// @Security ApiKeyAuth
// @Tags Services
// @Accept  json
// @Produce  json
// @Param sid path uint true "Service ID"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} destroyResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /services/{sid} [delete]
func Destroy(c *gin.Context) {
	sid, err := parseSID(c)
	if err != nil {
		log.Log.WithError(err).Error("URL Param Invalid")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = orm.DB.Delete(&models.Service{}, sid).Error; err != nil {
		log.Log.WithError(err).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"service": "",
	})
	return
}
