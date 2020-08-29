package middleware

import (
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/coolray-dev/raydash/modules/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Log adds a few useful entries to logrus Log
func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Log.Debug("Adding Logger Hook")
		log.Log.AddHook(&ginHook{
			context: c,
		})
		c.Next()
		utils.RemoveLoggerHook(log.Log, &ginHook{})
		return
	}
}

type ginHook struct {
	context *gin.Context
}

func (h *ginHook) Fire(entry *logrus.Entry) error {
	entry.Data["Path"] = h.context.Request.URL.Path
	return nil
}

func (h *ginHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
