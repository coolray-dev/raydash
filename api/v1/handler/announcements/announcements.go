package announcements

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	model "github.com/coolray-dev/raydash/models"
	"github.com/gin-gonic/gin"
)

type indexResponse struct {
	Total         uint64                `json:"total"`
	Announcements []models.Announcement `json:"announcements"`
}

// Index gets all anns and return to http in json
//
// Index godoc
// @Summary All Announcements
// @Description Simply list out all announcements
// @ID announcements.Index
// @Security ApiKeyAuth
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Access Token"
// @Success 200 {object} indexResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /announcements [get]
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
	query = query.Limit(int(limit.(uint64))).Offset(int(offset)).Order("updated_at desc")

	if err := query.Find(&anns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, indexResponse{
		Total:         uint64(len(anns)),
		Announcements: anns,
	})
	return
}

type createResponse struct {
	Announcement models.Announcement `json:"announcements"`
}

// Create parse a ann from request and create a new record in DB
//
// Create godoc
// @Summary Create Announcement
// @Description Create an announcement from post json object
// @ID Announcements.Create
// @Security ApiKeyAuth
// @Tags Announcements
// @Accept  json
// @Produce  json
// @Param ann body models.Announcement true "Announcement Object"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} createResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /announcements [post]
func Create(c *gin.Context) {
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

	if err := orm.DB.First(&ann).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	if err := orm.DB.First(&ann).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
