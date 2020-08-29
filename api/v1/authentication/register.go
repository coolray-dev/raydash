package authentication

import (
	"net/http"
	"strings"

	"github.com/coolray-dev/raydash/modules/utils"

	orm "github.com/coolray-dev/raydash/api/database"
	model "github.com/coolray-dev/raydash/api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Register handles POST /register which create a user
func Register(c *gin.Context) {
	type Request struct {
		Username string `binding:"required"`
		Password string `binding:"required"`
		Email    string `binding:"required"`
	}

	var json Request
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := utils.VerifyEmailFormat(json.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	user.Username = json.Username
	user.Password = utils.Hash(json.Password)
	user.Email = json.Email
	user.UUID = uuid.New().String()
	user.CurrentTraffic = 0
	user.MaxTraffic = 0

	if err := orm.DB.Create(&user).Error; err != nil {
		if strings.ContainsAny(err.Error(), "UNIQUE constraint failed:") {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
	return
}
