package authentication

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/coolray-dev/raydash/modules/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	model "github.com/coolray-dev/raydash/models"
	tool "github.com/coolray-dev/raydash/modules/utils"
)

// Login check username and password in request body json and returns access_token and refresh_token
func Login(c *gin.Context) {
	type Request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var json Request

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	var user model.User
	if err := orm.DB.Where("username = ?", json.Username).
		Where("password = ?", tool.Hash(json.Password)).
		First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	accessToken, err := jwt.SignAccessToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	var refreshToken string
	refreshToken, err = jwt.SignRefreshToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})

}

// Logout checks both accessToken in "Authorization" Header ( use middleware ) and refreshToken in request body json
// and delete the refreshToken from database
func Logout(c *gin.Context) {
	type Request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	var json Request
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, JWTErr := jwt.ParseUID(json.RefreshToken)
	if JWTErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": JWTErr.Error(),
		})
		return
	}

	var user models.User
	if err := orm.DB.Where("id = ?", uid).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if _, ok := user.Token[json.RefreshToken]; !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "RefreshToken not found",
		})
		return
	}
	delete(user.Token, json.RefreshToken)
	if err := orm.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
	return
}

// RefreshToken checks refreshToken and return a new accessToken
func RefreshToken(c *gin.Context) {
	// Bind request to a json
	type Request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwtTokens := strings.Split(req.RefreshToken, ".")
	if len(jwtTokens) != 3 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid AccessToken"})
		return
	}

	var payload jwt.TokenPayload
	dec, _ := base64.StdEncoding.DecodeString(jwtTokens[1] + "==")
	if err := json.Unmarshal(dec, &payload); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid RefreshToken: " + err.Error()})
		return
	}

	var user models.User
	if err := orm.DB.Where("id = ?", payload.UID).
		Where("username = ?", payload.Username).
		First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		// user not found
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Invalid RefreshToken: " + err.Error(),
		})
		return
	} else if err != nil {
		// unknown error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	key, jwtKeyErr := user.GetJwtKey()

	if jwtKeyErr != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Invalid RefreshToken: cannot get jwt key from current user",
		})
		return
	}

	_, jwtParseErr := jwt.Verify([]byte(req.RefreshToken), key)
	if jwtParseErr != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Invalid RefreshToken: " + jwtParseErr.Error(),
		})
		return
	}

	now := time.Now()

	// Validate refresh token
	if user.Token[req.RefreshToken].IsZero() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Token not found"})
		return
	}
	if now.After(user.Token[req.RefreshToken]) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Token expired"})
		return
	}

	// Renew its expiry time
	user.Token[req.RefreshToken] = now.Add(24 * time.Hour)
	if err := orm.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Sign access token and return
	accessToken, err := jwt.SignAccessToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
	return
}
