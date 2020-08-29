package authentication

import (
	"net/http"

	"github.com/coolray-dev/raydash/modules/setting"

	orm "github.com/coolray-dev/raydash/api/database"
	"github.com/coolray-dev/raydash/api/models"
	model "github.com/coolray-dev/raydash/api/models"
	"github.com/coolray-dev/raydash/modules/mail"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	if query := orm.DB.Where("email = ?", json.Email).First(&user); query.RecordNotFound() {
		c.JSON(http.StatusOK, gin.H{
			"message": "Password reset token has been sent to you if you already registered",
		})
		return
	} else if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": query.Error.Error(),
		})
		return
	}

	if orm.DB.Where("user_id = ?", user.ID).First(&fp).RecordNotFound() {
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
