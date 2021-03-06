package authentication

import (
	"errors"
	"net/http"

	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ResetPassword(c *gin.Context) {
	type Request struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	var json Request

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	var fp model.ForgetPassword

	if err := orm.DB.Where("token = ?", json.Token).
		Preload("User").
		First(&fp).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "invalid token",
		})
		return
	}

	orm.DB.Model(fp.User).Update("password", utils.Hash(json.Password))
	orm.DB.Unscoped().Delete(&fp)

	c.Status(http.StatusNoContent)
}
