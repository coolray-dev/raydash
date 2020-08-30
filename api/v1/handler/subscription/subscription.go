package subscription

import (
	"bytes"
	"net/http"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/clash"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Clash(c *gin.Context) {

	subToken := c.Query("token")
	var user models.User
	if query := orm.DB.Preload("Groups").Where("subscription_token = ?", subToken).First(&user); query.RecordNotFound() {
		log.Log.Info("No Such User in Database")
		c.Status(http.StatusNotFound)
		return
	} else if query.Error != nil {
		log.Log.WithFields(logrus.Fields{
			"error": query.Error.Error(),
		}).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}
	var nodes []*models.Node
	for _, g := range user.Groups {
		if err := orm.DB.Preload("Nodes").Where("ID = ?", g.ID).First(&g).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		nodes = append(nodes, g.Nodes...)

	}
	var services []*models.Service
	for _, n := range nodes {
		if err := orm.DB.Preload("Services").Where("ID = ?", n.ID).First(&n).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		services = append(services, n.Services...)
	}

	var b bytes.Buffer
	clash.Generate(services, &b)
	c.String(http.StatusOK, b.String())
	return
}
