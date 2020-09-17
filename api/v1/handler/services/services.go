package services

import (
	"fmt"
	"strconv"

	"github.com/coolray-dev/raydash/models"
	"github.com/gin-gonic/gin"
)

type serviceResponse struct {
	Service models.Service `json:"service"`
}

func parseSID(c *gin.Context) (sid uint64, err error) {
	sid, err = strconv.ParseUint(c.Param("sid"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid SID: %w", err)
	}
	return
}
