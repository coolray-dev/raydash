package users

import (
	"net/http"

	"github.com/coolray-dev/raydash/api/v1/handler"
	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
	"github.com/gin-gonic/gin"
)

type updateResponse struct {
	User model.User `json:"user"`
}

type updateRequest struct {
	UUID              string `json:"uuid" fake:"{uuid}"`
	Email             string `json:"email" fake:"{email}"`
	SubscriptionToken string `json:"subscription_token"`
}

// Update receive a user object and update it
//
// Update godoc
// @Summary Update user
// @Description Update the user provided
// @ID users.Update
// @Security ApiKeyAuth
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body updateRequest true "User Object"
// @Param username path string true "Username"
// @Param Authorization header string false "Node Token"
// @Success 200 {object} updateResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /users/{username} [patch]
func Update(c *gin.Context) {
	var user model.User
	username := c.Param("username")
	if err := orm.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, &handler.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var json updateRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, &handler.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	if json.UUID != "" {
		user.UUID = json.UUID
	}
	if json.Email != "" {
		user.Email = json.Email
	}
	if json.SubscriptionToken != "" {
		user.SubscriptionToken = json.SubscriptionToken
	}

	if err := orm.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, &handler.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &updateResponse{
		User: user,
	})
	return
}
