package announcements

import (
	"fmt"
	"net/http"
	"strconv"

	orm "github.com/coolray-dev/raydash/api/database"
	model "github.com/coolray-dev/raydash/api/models"
	"github.com/gin-gonic/gin"
)

// Index gets all anns and return to http in json
func Index(c *gin.Context) {
	var anns []model.Announcement
	query := orm.DB
	if before, exists := c.Get("before"); exists {
		query = query.Where(" updated_at <= ?", before)
	}
	if after, exists := c.Get("after"); exists {
		query = query.Where(" updated_at >= ?", after)
	}

	const defaultPage uint64 = 1
	const defaultLimit uint64 = 10
	limit, limitexists := c.Get("limit")
	if !limitexists {
		limit = defaultLimit
	}
	page, pageexists := c.Get("page")
	if !pageexists {
		page = defaultPage
	}

	offset := limit.(uint64) * (page.(uint64) - 1)
	query = query.Limit(limit).Offset(offset).Order("updated_at desc")

	if err := query.Find(&anns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total":         len(anns),
		"announcements": anns,
	})
	return
}

// Store parse a ann from request and create a new record in DB
func Store(c *gin.Context) {
	type Request struct {
		Content string `binding:"required"`
		Level   string `binding:"required"`
		Title   string `binding:"required"`
	}

	var json Request
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ann model.Announcement
	ann.Content = json.Content
	ann.Level = json.Level
	ann.Title = json.Title

	if err := orm.DB.Create(&ann).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"annoucement": ann,
	})
	return

}

// Show receive a id from request and find the ann of the specific id
func Show(c *gin.Context) {
	aid, err := parseAID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ann model.Announcement

	ann.ID = aid

	if query := orm.DB.First(&ann); query.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": query.Error.Error()})
		return
	} else if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"announcement": ann,
	})
	return

}

// Update receive a id and a ann from request and update the specific record in DB
func Update(c *gin.Context) {
	aid, err := parseAID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	type Request struct {
		Content string `binding:"required"`
		Level   string `binding:"required"`
		Title   string `binding:"required"`
	}
	var json Request
	if err = c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ann model.Announcement
	ann.ID = aid

	if query := orm.DB.First(&ann); query.RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": query.Error.Error()})
		return
	} else if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	ann.Content = json.Content
	ann.Level = json.Level
	ann.Title = json.Title

	if err = orm.DB.Save(&ann).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"announcement": ann,
	})
	return
}

// Destroy receive a id from request and delete in from DB
func Destroy(c *gin.Context) {
	aid, err := parseAID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var ann model.Announcement
	ann.ID = aid

	if err = orm.DB.Delete(&ann).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"announcement": "",
	})
	return
}

func parseAID(c *gin.Context) (aid uint64, err error) {
	aid, err = strconv.ParseUint(c.Param("aid"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid AID: %w", err)
	}
	return
}
