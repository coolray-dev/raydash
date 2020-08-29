package authentication

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	orm "github.com/coolray-dev/raydash/api/database"
	"github.com/coolray-dev/raydash/api/models"
	model "github.com/coolray-dev/raydash/api/models"
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
	if query := orm.DB.Where("username = ?", json.Username).Where("password = ?", tool.Hash(json.Password)).First(&user); query.RecordNotFound() {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	} else if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": query.Error.Error(),
		})
		return
	}
	accessToken, err := SignAccessToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	var refreshToken string
	refreshToken, err = SignRefreshToken(&user)
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
	uid, _ := c.Get("uid")
	var user model.User
	user.ID = uid.(uint64)

	if query := orm.DB.First(&user); query.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{
			"error": query.Error.Error(),
		})
		return
	} else if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": query.Error.Error(),
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

	var payload TokenPayload
	dec, _ := base64.StdEncoding.DecodeString(jwtTokens[1] + "==")
	if err := json.Unmarshal(dec, &payload); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid RefreshToken: " + err.Error()})
		return
	}

	var user models.User
	if query := orm.DB.Where("id = ?", payload.UID).
		Where("username = ?", payload.Username).First(&user); query.RecordNotFound() {
		// user not found
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Invalid RefreshToken: " + query.Error.Error(),
		})
		return
	} else if query.Error != nil {
		// unknown error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": query.Error.Error(),
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

	_, jwtParseErr := Verify([]byte(req.RefreshToken), key)
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
	accessToken, err := SignAccessToken(&user)
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
