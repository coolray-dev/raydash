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
// @ID Announcements.Index
// @Security ApiKeyAuth
// @Tags Announcements
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
// @Param Annoucement body annRequest true "Announcement Object"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} annResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /announcements [post]
func Create(c *gin.Context) {

	var json annRequest
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

type annResponse struct {
	Announcement models.Announcement `json:"announcement"`
}

// Show receive a id from request and find the ann of the specific id
//
// Show godoc
// @Summary Show Announcements
// @Description Show a announcement according to id
// @ID Announcements.Show
// @Security ApiKeyAuth
// @Tags Announcements
// @Accept  json
// @Produce  json
// @Param aid path uint true "Announcement ID"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} annResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /announcements/{aid} [get]
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
	c.JSON(http.StatusOK, annResponse{
		Announcement: ann,
	})
	return

}

type annRequest struct {
	Content string `binding:"required"`
	Level   string `binding:"required"`
	Title   string `binding:"required"`
}

// Update receive a id and a ann from request and update the specific record in DB
//
// Update godoc
// @Summary Update Announcement
// @Description Update a announcement
// @ID Announcements.Update
// @Security ApiKeyAuth
// @Tags Announcements
// @Accept  json
// @Produce  json
// @Param aid path uint true "Announcement ID"
// @Param Announcement body annRequest true "Announcement Object"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} annResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /announcements/{aid} [patch]
func Update(c *gin.Context) {
	aid, err := parseAID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var json annRequest
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
	c.JSON(http.StatusOK, annResponse{
		Announcement: ann,
	})
	return
}

type destroyResponse struct {
	Announcement string `json:"announcement"`
}

// Destroy receive a id from request and delete in from DB
//
// Destroy godoc
// @Summary Destroy Announcement
// @Description Destroy an announcement according to nid
// @ID Announcements.Destroy
// @Security ApiKeyAuth
// @Tags Announcements
// @Accept  json
// @Produce  json
// @Param aid path uint true "Announcement ID"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} destroyResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /announcements/{aid} [delete]
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
	c.JSON(http.StatusOK, destroyResponse{
		Announcement: "",
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
