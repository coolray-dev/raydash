package users

import (
	"net/http"

	"github.com/coolray-dev/raydash/api/v1/handler"
	orm "github.com/coolray-dev/raydash/database"
	model "github.com/coolray-dev/raydash/models"
	"github.com/gin-gonic/gin"
)

type indexResponse struct {
	Total uint
	Users []model.User `json:"users"`
}

// Index handle GET /users which simply list out all users
// accept gid param as filter
//
// Index godoc
// @Summary All Users
// @Description Simply list out all users
// @ID users.Index
// @Security ApiKeyAuth
// @Tags Users
// @Accept  json
// @Produce  json
// @Param gid query uint false "Group ID"
// @Param Authorization header string false "Node Token"
// @Success 200 {object} indexResponse
// @Failure 403 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /users [get]
func Index(c *gin.Context) {
	var u []model.User
	users := &u

	if err := orm.DB.Preload("Groups").Find(users).Order("updated_at desc").Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			&handler.ErrorResponse{Error: err.Error()})
		return
	}

	// Apply GID Filter
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

	c.JSON(http.StatusOK, &indexResponse{
		Total: uint(len(*users)),
		Users: *users,
	})
}
