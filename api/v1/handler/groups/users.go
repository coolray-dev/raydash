package groups

import (
	"errors"
	"net/http"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/casbin"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type usersResponse struct {
	Users []*models.User
}

type userModificationRequest struct {
	Username string `binding:"required"`
}

// Users return all users a group have
//
// Users godoc
// @Summary Get Group Users
// @Description Simply list out all users belong to a certain group
// @ID groups.Users
// @Security ApiKeyAuth
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param gid path uint true "Group ID"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} usersResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /groups/{gid}/users [get]
func Users(c *gin.Context) {
	var group models.Group
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
	c.JSON(http.StatusOK, usersResponse{
		Users: group.Users,
	})
	return
}

// AppendUser add a user to specific group
//
// AppendUser godoc
// @Summary Append User
// @Description Add a user to a group
// @ID groups.AppendUser
// @Security ApiKeyAuth
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param user body userModificationRequest true "Username"
// @Param gid path uint true "Group ID"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} usersResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /groups/{gid}/users [patch]
func AppendUser(c *gin.Context) {
	var group models.Group
	gid, err := parseGID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var json userModificationRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
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

	// Add User to Group in casbin
	casbin.Enforcer.AddGroupingPolicy(user.Username, group.Name)

	c.JSON(http.StatusOK, usersResponse{
		Users: group.Users,
	})
	return
}

// RemoveUser remove a user from a specific group
//
// RemoveUser godoc
// @Summary Remove User
// @Description Remove a user from a group
// @ID groups.RemoveUser
// @Security ApiKeyAuth
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param gid path uint true "Group ID"
// @Param username path string true "Username"
// @Param Authorization header string true "Access Token"
// @Success 200 {object} usersResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /groups/{gid}/users/{username} [delete]
func RemoveUser(c *gin.Context) {
	var group models.Group
	username := c.Param("username")
	gid, err := parseGID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
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

	// Remove User from Group in casbin
	casbin.Enforcer.RemoveGroupingPolicy(username, group.Name)

	c.JSON(http.StatusOK, usersResponse{
		Users: group.Users,
	})
	return
}
