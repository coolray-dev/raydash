package users

import (
	"net/http"

	"github.com/coolray-dev/raydash/api/v1/handler"
	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
	"github.com/gin-gonic/gin"
)

type updateResponse struct {
	User model.User
}

// Update receive a user object and update it
//
// Update godoc
// @Summary Update user
// @Description Update the user provided
// @ID users.Update
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body models.User true "User Object"
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

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, &handler.ErrorResponse{
			Error: err.Error(),
		})
		return
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
