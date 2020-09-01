package groups

import (
	"net/http"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/gin-gonic/gin"
)

type indexResponse struct {
	Total  uint
	Groups []models.Group
}

// Index handle GET /groups which simply list out all groups
//
// Index godoc
// @Summary All Groups
// @Description Simply list out all groups
// @ID groups.Index
// @Security ApiKeyAuth
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Access Token"
// @Success 200 {object} indexResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /groups [get]
func Index(c *gin.Context) {
	var groups []models.Group
	if err := orm.DB.Find(&groups).Order("updated_at desc").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, indexResponse{
		Total:  uint(len(groups)),
		Groups: groups,
	})
}
