package groups

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
)

// Index handle GET /groups which simply list out all groups
func Index(c *gin.Context) {
	var groups []model.Group
	if err := orm.DB.Find(&groups).Order("updated_at desc").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":  len(groups),
		"groups": groups,
	})
}

// Show receive a id from request url and return the group of the specific id
func Show(c *gin.Context) {
	gid, err := parseGID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var group model.Group
	group.ID = gid

	if err := orm.DB.First(&group).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"group": group,
	})
	return
}

// Create receive a group object from request and update the specific record in DB
func Create(c *gin.Context) {
	var group model.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := orm.DB.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"group": group,
	})
	return
}

// Update receive a id and a group object from request and update the specific record in DB
func Update(c *gin.Context) {
	gid, err := parseGID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var group model.Group
	group.ID = gid
	if err = c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = orm.DB.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"group": group,
	})
	return
}

// Destroy receive a id from request and delete in from DB
func Destroy(c *gin.Context) {
	gid, err := parseGID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var group model.Group
	group.ID = gid

	if err = orm.DB.Delete(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"group": "",
	})
	return
}

// Users return all users a group have
func Users(c *gin.Context) {
	var group model.Group
	gid, err := parseGID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := orm.DB.Preload("Users").Where("ID = ?", gid).First(&group).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"users": "[]",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": group.Users,
	})
	return
}

// AppendUser add a user to specific group
func AppendUser(c *gin.Context) {
	var group model.Group
	gid, err := parseGID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type Request struct {
		Username string `json:"username" binding:"required"`
	}
	var json Request

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user model.User
	if err := orm.DB.Where("username = ?", json.Username).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"users": "[]",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := orm.DB.Preload("Users").Where("ID = ?", gid).First(&group).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"users": "[]",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	group.Users = append(group.Users, &user)
	if err = orm.DB.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": group.Users,
	})
	return
}

// RemoveUser remove a user from a specific group
func RemoveUser(c *gin.Context) {
	var group model.Group
	username := c.Param("username")
	gid, err := parseGID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := orm.DB.Where("username = ?", username).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"users": "[]",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := orm.DB.Preload("Users").Where("ID = ?", gid).First(&group).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"users": "[]",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var index int
	for i := range group.Users {
		if group.Users[i].ID == user.ID {
			index = i
			break
		}
	}
	group.Users = append(group.Users[:index], group.Users[index+1:]...)
	if err = orm.DB.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": group.Users,
	})
	return
}

func parseGID(c *gin.Context) (gid uint64, err error) {
	gid, err = strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid GID: %w", err)
	}
	return
}
