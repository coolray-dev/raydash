package users

import (
	"net/http"

	"github.com/coolray-dev/raydash/api/v1/handler"
	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/gin-gonic/gin"
)

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
// @Summary User traffic
// @Description Update user traffic
// @ID users.Traffic
// @Security ApiKeyAuth
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body trafficRequest true "User Object"
// @Param username path string true "Username"
// @Param Authorization header string false "Node Token"
// @Success 200 {object} userResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /users/{username}/traffic [patch]
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

	c.JSON(http.StatusOK, &userResponse{
		User: user,
	})
	return
}
