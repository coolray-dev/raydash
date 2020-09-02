package middleware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/casbin"
	"github.com/coolray-dev/raydash/modules/jwt"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Authorize Parse identity info from header and check permission
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Store role within this whole function
		var role string
		var subject string // casbin auth subject

		// Check Authorization Header
		kind, token, headerErr := checkHeader(c)
		if headerErr != nil {
			log.Log.Info(headerErr.Error())
			role = "anonymous"
			goto casbin
		}

		switch kind {
		case "node":
			var node models.Node
			if err := orm.DB.Where("access_token = ?", token).First(&node).Error; errors.Is(err, gorm.ErrRecordNotFound) {
				log.Log.Debug("Token Not Matching Any Node")
				role = "anonymous"
				break
			} else if err != nil {
				log.Log.WithError(err).Info("Database Error")
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				c.Abort()
				return
			} else {
				log.Log.Debug("Token Match")
				role = "node"
				subject = "node::" + strconv.Itoa(int(node.ID))
			}

		case "jwt":
			tokenSplit := strings.Split(token, ".")
			var payload jwt.TokenPayload
			dec, base64err := base64.RawURLEncoding.DecodeString(tokenSplit[1])
			if base64err != nil {
				log.Log.WithError(base64err).Info("Invalid JWT Token")
				role = "anonymous"
				break
			}
			if err := json.Unmarshal(dec, &payload); err != nil {
				log.Log.WithError(err).Info("Invalid JWT Token")
				role = "anonymous"
				break
			}

			var user models.User
			if err := orm.DB.Where("id = ?", payload.UID).
				Where("username = ?", payload.Username).
				First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {

				log.Log.Info("Invalid User")
				role = "anonymous"
				break
			} else if err != nil {
				log.Log.WithError(err).Error("Database Error")
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				c.Abort()
				return
			}

			key, jwtKeyErr := user.GetJwtKey()
			if jwtKeyErr != nil {
				log.Log.WithError(jwtKeyErr).Info("Error Getting User JWT Key")
				role = "anonymous"
				break
			}

			if plain, err := jwt.Verify([]byte(token), key); err != nil {
				log.Log.WithError(err).Info("JWT Verification Failed")
				role = "anonymous"
				break
			} else if plain.Subject != "AccessToken" {
				log.Log.Info("JWT Subject not matching 'AccessToken', might have used refresh token")
				role = "anonymous"
				break
			} else {
				log.Log.WithField("Expire", plain.ExpirationTime).Debug("JWT Verification Success")
				role = "user"
				subject = plain.Username
			}

		default:
			log.Log.Panic("Unknown Error")
		}

	casbin:
		allow, err := casbinAuthorize(role, subject, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		if !allow {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Permission Denied",
			})
			c.Abort()
			return
		}
	}
}

func checkHeader(c *gin.Context) (kind string, token string, err error) {

	// Get Header
	header := c.GetHeader("Authorization")
	if header == "" {
		err = errors.New("Empty Authorization Header")
		return
	}

	headerList := strings.Split(header, " ")
	if len(headerList) != 2 {
		err = errors.New("Invalid Authorization Header")
		return
	}

	t := headerList[0]
	content := headerList[1]
	if t != "Bearer" {
		err = errors.New("Only Support Bearer Authorization")
		return
	}
	contentList := strings.Split(content, ".")

	switch len(contentList) {
	case 2:
		if contentList[0] != "node" {
			err = errors.New("Invalid Token Structure")
			return
		}
		kind = "node"
		token = contentList[1]
		return
	case 3:
		kind = "jwt"
		token = content
		return
	default:
		err = errors.New("Invalid Authorization Token")
		return
	}
}

func casbinAuthorize(role, sub, obj, act string) (bool, error) {
	var res bool
	var err error

	switch role {
	case "anonymous":
		res, err = casbin.Enforcer.Enforce("role::anonymous", obj, act)
	case "node":
		res, err = casbin.Enforcer.Enforce(sub, obj, act)
	case "user":
		res, err = casbin.Enforcer.Enforce(sub, obj, act)
	default:
		return false, errors.New("Invalid role")
	}

	if err != nil {
		return false, err
	}
	if !res {
		log.Log.WithField("Role", role).WithField("Subject", sub).Debug("Access Denied")
	} else {
		log.Log.WithField("Role", role).WithField("Subject", sub).Debug("Access Allowed")
	}
	return res, nil
}
