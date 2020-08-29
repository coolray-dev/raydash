package options

import (
	"net/http"

	orm "github.com/coolray-dev/raydash/api/database"
	model "github.com/coolray-dev/raydash/api/models"
	"github.com/gin-gonic/gin"
)

// Index list out all options stored in DB
func Index(c *gin.Context) {
	var opts []model.Option
	if err := orm.DB.Find(&opts).Order("updated_at desc").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total":   len(opts),
		"options": opts,
	})
}

// Update receive a option object and store it
func Update(c *gin.Context) {
	name := c.Param("name")
	type Request struct {
		Value string `binding:"required"`
	}
	var json Request
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var opt model.Option
	opt.Value = json.Value
	opt.Name = name

	if err := orm.DB.Save(&opt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"option": opt,
	})
	return
}
