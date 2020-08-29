package nodes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
)

// Index handle GET /nodes which simply list out all nodes
// accept gid param as filter
func Index(c *gin.Context) {
	var n []model.Node
	nodes := &n
	if err := orm.DB.Preload("Groups").Find(nodes).Order("updated_at desc").Error; err != nil {
		log.Log.WithError(err).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if gid, exists := c.Get("gid"); exists {
		var t []model.Node
		for _, i := range *nodes {
			for _, group := range i.Groups {
				if gid == group.ID {
					t = append(t, i)
					break
				}
			}
		}
		nodes = &t
	}

	// Check Admin
	// Show Access Token if isAdmin
	if isAdmin := c.MustGet("isAdmin").(bool); !isAdmin {
		for i := range *nodes {
			(*nodes)[i].AccessToken = ""
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(*nodes),
		"nodes": nodes,
	})
}

// Create receive a id and a node object from request and update the specific record in DB
func Create(c *gin.Context) {

	// Check Admin
	if isAdmin := c.MustGet("isAdmin").(bool); !isAdmin {
		log.Log.WithFields(logrus.Fields{
			"isAdmin": isAdmin,
		}).Warning("Node Creation Failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only Admin Can Add Node"})
		return
	}

	// Bind Request
	var node model.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		log.Log.WithError(err).Warning("Could not bind request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save to DB
	if err := orm.DB.Save(&node).Error; err != nil {
		log.Log.WithError(err).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return result
	log.Log.Debug("Success")
	c.JSON(http.StatusOK, gin.H{
		"node": node,
	})
	return
}

// Show receive a id from request url and return the node of the specific id
func Show(c *gin.Context) {

	// Get Node ID
	nid, err := parseNID(c)
	if err != nil {
		log.Log.WithError(err).Warn("Error Getting Node ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if is a node access
	isNode, isNodeExists := c.Get("isNode")
	if isNodeExists {
		if isNode.(bool) {
			if nid != c.MustGet("nodeID").(uint64) {
				log.Log.WithField("nodeID", c.MustGet("nodeID").(uint64)).Warn("Node Not Authorized")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Not Authorized"})
				return
			}
		}
	}

	var node model.Node
	node.ID = nid

	if query := orm.DB.First(&node); query.RecordNotFound() {
		log.Log.WithField("nodeID", nid).Warn("Node Not Found")
		c.JSON(http.StatusNotFound, gin.H{"error": query.Error.Error()})
		return
	} else if query.Error != nil {
		log.Log.WithError(err).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"node": node,
	})
	log.Log.Debug("Success")
	return
}

// Update receive a id and a node object from request and update the specific record in DB
func Update(c *gin.Context) {

	// Check Admin
	if isAdmin := c.MustGet("isAdmin").(bool); !isAdmin {
		log.Log.WithFields(logrus.Fields{
			"isAdmin": isAdmin,
		}).Warning("Node Creation Failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only Admin Can Modify Node"})
		return
	}

	nid, err := parseNID(c)
	if err != nil {
		log.Log.WithError(err).Warn("Error Getting Node ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var node model.Node
	node.ID = nid

	// Bind Request
	if err = c.ShouldBindJSON(&node); err != nil {
		log.Log.WithError(err).Warn("Error Binding Request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save to DB
	if err = orm.DB.Save(&node).Error; err != nil {
		log.Log.WithError(err).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"node": node,
	})
	log.Log.Debug("Success")
	return
}

// Destroy receive a id from request and delete in from DB
func Destroy(c *gin.Context) {

	// Check Admin
	if isAdmin := c.MustGet("isAdmin").(bool); !isAdmin {
		log.Log.WithFields(logrus.Fields{
			"isAdmin": isAdmin,
		}).Warning("Node Creation Failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only Admin Can Delete Node"})
		return
	}

	nid, err := parseNID(c)
	if err != nil {
		log.Log.WithError(err).Warn("Error Getting Node ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var node model.Node
	node.ID = nid

	if err = orm.DB.Delete(&node).Error; err != nil {
		log.Log.WithError(err).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"node": "",
	})
	log.Log.Debug("Success")
	return
}

func parseNID(c *gin.Context) (nid uint64, err error) {
	nid, err = strconv.ParseUint(c.Param("nid"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid NID: %w", err)
	}
	return
}
