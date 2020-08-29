package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ParseParams parse a set of parameters in url and pass them to handler through gin.Context
func ParseParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		parseParamUint64("limit", "0", "_limit", c)
		parseParamUint64("page", "0", "_page", c)
		parseParamUint64("gid", "0", "gid", c)
		parseParamUint64("uid", "0", "uid", c)
		parseParamUint64("nid", "0", "nid", c)
		if _, err := parseParamTime("before", "0", "_before", c); err != nil {
			c.Abort()
			return
		}
		if _, err := parseParamTime("after", "0", "_after", c); err != nil {
			c.Abort()
			return
		}
		c.Next()
	}
}

func parseParamUint64(varName string, defaultParam string, paramName string, c *gin.Context) {
	if param, err := strconv.ParseUint(c.DefaultQuery(paramName, defaultParam), 10, 0); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Param '" + paramName + "' not uint"})
		c.Abort()
		return
	} else if param == 0 {
		return
	} else if param < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid param " + paramName})
		c.Abort()
		return
	} else {
		c.Set(varName, param)
		return
	}
}

func parseParamTime(varName string, defaultParam string, paramName string, c *gin.Context) (param time.Time, err error) {
	var unixTimeStamp int64
	if unixTimeStamp, err = strconv.ParseInt(c.DefaultQuery(paramName, defaultParam), 10, 0); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Param '" + paramName + "' not unix timestamp"})
		return
	} else if unixTimeStamp == 0 {
		return
	}
	param = time.Unix(unixTimeStamp, 0)
	c.Set(varName, param)
	return

}
