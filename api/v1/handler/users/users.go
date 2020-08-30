package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
)

// Index handle GET /users which simply list out all users
// accept gid param as filter
// Index godoc
// @Summary Simply list out all users
// @Description Simply list out all users
// @ID users.Index
// @Accept  json
// @Produce  json
// @Param gid query int false "Group ID"
// @Success 200 {object} []models.User
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /users [get]
func Index(c *gin.Context) {
	var u []model.User
	users := &u
	if err := orm.DB.Preload("Groups").Find(users).Order("updated_at desc").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if gid, exists := c.Get("gid"); exists {
		var t []model.User
		for _, i := range *users {
			for _, group := range i.Groups {
				if gid == group.ID {
					t = append(t, i)
					break
				}
			}
		}
		users = &t
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(*users),
		"users": *users,
	})
}

// Show query a user and return it using url param "username" as condition
func Show(c *gin.Context) {
	var user model.User
	username := c.Param("username")
	if err := orm.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
	return
}

// Update receive a user object and update it
func Update(c *gin.Context) {
	var user model.User
	username := c.Param("username")
	if err := orm.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := orm.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusNotModified, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
	return
}

// Destroy delete a user from db
func Destroy(c *gin.Context) {
	username := c.Param("username")
	tokenUsername := c.MustGet("username").(string)
	isAdmin := c.MustGet("isAdmin").(bool)
	if (username != tokenUsername) && !isAdmin {

		c.JSON(http.StatusForbidden, gin.H{"error": "No permission"})
		return
	}
	if err := orm.DB.Where("username = ?", username).Delete(model.User{}).Error; err != nil {
		c.JSON(http.StatusNotModified, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": "",
	})
}

// Nodes return all nodes a user has
func Nodes(c *gin.Context) {
	var user model.User
	username := c.Param("username")
	if err := orm.DB.Preload("Groups").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var nodes []*model.Node
	for _, g := range user.Groups {
		if err := orm.DB.Preload("Nodes").Where("ID = ?", g.ID).First(&g).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		nodes = append(nodes, g.Nodes...)

	}
	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
	})
	return
}

// Groups shows all groups a user has
func Groups(c *gin.Context) {
	var user model.User
	username := c.Param("username")
	if err := orm.DB.Preload("Groups").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"groups": user.Groups,
	})
	return
}

// Services returns all services a user has
func Services(c *gin.Context) {
	var user model.User
	username := c.Param("username")
	if err := orm.DB.Preload("Groups").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var nodes []*model.Node
	for _, g := range user.Groups {
		if err := orm.DB.Preload("Nodes").Where("ID = ?", g.ID).First(&g).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		nodes = append(nodes, g.Nodes...)

	}
	var services []*model.Service
	for _, n := range nodes {
		if err := orm.DB.Preload("Services").Where("ID = ?", n.ID).First(&n).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		services = append(services, n.Services...)
	}
	c.JSON(http.StatusOK, gin.H{
		"services": services,
	})
	return
}
