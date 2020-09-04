package nodes

import (
	"net/http"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/gin-gonic/gin"
)

type indexResponse struct {
	Total uint
	Nodes []models.Node `json:"nodes"`
}

// Index handle GET /nodes which simply list out all nodes
// accept gid param as filter
//
// Index godoc
// @Summary All Nodes
// @Description Simply list out all Nodes
// @ID Nodes.Index
// @Security ApiKeyAuth
// @Tags Nodes
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Access Token"
// @Success 200 {object} indexResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /nodes [get]
func Index(c *gin.Context) {
	var n []models.Node
	nodes := &n
	if err := orm.DB.Preload("Groups").Find(nodes).Order("updated_at desc").Error; err != nil {
		log.Log.WithError(err).Error("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if gid, exists := c.Get("gid"); exists {
		var t []models.Node
		for _, n := range *nodes {
			for _, group := range n.Groups {
				if gid == group.ID {
					t = append(t, n)
					break
				}
			}
		}
		nodes = &t
	}

	c.JSON(http.StatusOK, indexResponse{
		Total: uint(len(*nodes)),
		Nodes: *nodes,
	})
}
