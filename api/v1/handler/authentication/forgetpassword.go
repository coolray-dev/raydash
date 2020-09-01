package authentication

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	model "github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/mail"
	"github.com/coolray-dev/raydash/modules/setting"
)

func ForgetPassword(c *gin.Context) {
	type Request struct {
		Email string `json:"email"`
	}
	var json Request

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	var user model.User
	var fp model.ForgetPassword

	if err := orm.DB.Where("email = ?", json.Email).
		First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Password reset token has been sent to you if you already registered",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := orm.DB.Where("user_id = ?", user.ID).First(&fp).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		fp.Token = uuid.New().String()
		fp.User = &user
		orm.DB.Create(&fp)
	} else {
		orm.DB.Model(&fp).Update("token", uuid.New().String())
	}

	mail.MailChan <- &models.Mail{
		From:        setting.Config.GetString("mail.from"),
		To:          user.Email,
		Subject:     "Password Reset",
		ContentType: "text/html",
		Content:     fp.Token,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset token has been sent to you if you already registered",
	})
}
