package groups

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func parseGID(c *gin.Context) (gid uint64, err error) {
	gid, err = strconv.ParseUint(c.Param("gid"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid GID: %w", err)
	}
	return
}
