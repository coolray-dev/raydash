package middleware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	auth "github.com/coolray-dev/raydash/api/v1/handler/authentication"
	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthNodeToken check if request is from a node instead of a user
func AuthNodeToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			log.Log.Info("Empty Authorization Header")
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid AccessToken"})
			c.Abort()
			return
		}
		headerList := strings.Split(header, " ")
		if len(headerList) != 2 {
			log.Log.Info("Invalid Authorization Header")
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid AccessToken"})
			c.Abort()
			return
		}
		t := headerList[0]
		token := headerList[1]
		if t != "Bearer" {
			log.Log.Info("Only Support Bearer Authorization")
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid AccessToken"})
			c.Abort()
			return
		}
		var node models.Node
		if err := orm.DB.Where("access_token = ?", token).First(&node).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Log.Debug("Token Not Matching Any Node")
			c.Set("isNode", false)
			return
		} else if err != nil {
			log.Log.WithError(err).Info("Database Error")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		log.Log.Debug("Token Match")
		c.Set("isNode", true)
		c.Set("nodeID", node.ID)
		return
	}
}

// AuthAccessToken parse the Authorization Header from request
// and pass
//
// *models.User
// uid
// username
//
// to gin Context
func AuthAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Check if is a node access
		// if yes then should not parse jwt, directly return
		isNode, isNodeExists := c.Get("isNode")
		if isNodeExists {
			if isNode.(bool) {
				return
			}
		}

		header := c.GetHeader("Authorization")
		if header == "" {
			log.Log.Info("Empty Authorization Header")
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid AccessToken"})
			c.Abort()
			return
		}
		headerList := strings.Split(header, " ")
		if len(headerList) != 2 {
			log.Log.Info("Invalid Authorization Header")
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid AccessToken"})
			c.Abort()
			return
		}
		t := headerList[0]
		content := headerList[1]
		if t != "Bearer" {
			log.Log.Info("Only Support Bearer Authorization")
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid AccessToken"})
			c.Abort()
			return
		}
		contentList := strings.Split(content, ".")
		if len(contentList) != 3 {
			log.Log.Info("Invalid JWT Token")
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid AccessToken"})
			c.Abort()
			return
		}

		var payload auth.TokenPayload
		dec, _ := base64.StdEncoding.DecodeString(contentList[1] + "==")
		if err := json.Unmarshal(dec, &payload); err != nil {
			log.Log.WithError(err).Info("Invalid JWT Token")
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid AccessToken: " + err.Error()})
			c.Abort()
			return
		}

		var user models.User
		if err := orm.DB.Where("id = ?", payload.UID).
			Where("username = ?", payload.Username).
			First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {

			log.Log.Info("Invalid User")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			c.Abort()
		} else if err != nil {
			log.Log.WithError(err).Error("Database Error")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			c.Abort()
		}

		key, jwtKeyErr := user.GetJwtKey()
		if jwtKeyErr != nil {
			log.Log.WithError(jwtKeyErr).Info("Error Getting User JWT Key")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": jwtKeyErr.Error(),
			})
		}

		if plain, err := auth.Verify([]byte(content), key); err != nil {
			log.Log.WithError(err).Info("JWT Verification Failed")
			c.JSON(http.StatusForbidden, gin.H{
				"error": "AccessToken invalid: " + err.Error(),
			})
			c.Abort()
			return
		} else if plain.Subject != "AccessToken" {
			log.Log.Info("JWT Subject not matching 'AccessToken', might have used refresh token")
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Not Access Token",
			})
			c.Abort()
			return
		} else {
			log.Log.WithField("Expire", plain.ExpirationTime).Debug("JWT Verification Success")
		}

		c.Set("user", &user)
		c.Set("uid", payload.UID)
		c.Set("username", payload.Username)
		return
	}
}

// AuthAdmin required AuthToken before it
// it passes a bool called isAdmin to Context
func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Check if is a node access
		// if yes then not admin
		isNode, isNodeExists := c.Get("isNode")
		if isNodeExists {
			if isNode.(bool) {
				c.Set("isAdmin", false)
				return
			}
		}

		var user models.User
		tokenUsername := c.MustGet("username")
		if err := orm.DB.Preload("Groups").Where("username = ?", tokenUsername).First(&user).Error; err != nil {
			log.Log.WithError(err).Error("Database Error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		for _, g := range user.Groups {
			if g.ID == 1 {
				c.Set("isAdmin", true)
				return
			}
		}
		c.Set("isAdmin", false)
	}
}
