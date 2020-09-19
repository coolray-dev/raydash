package services

import (
	"fmt"
	"strconv"

	"github.com/coolray-dev/raydash/models"
	"github.com/gin-gonic/gin"
)

// Unified Service Object
type serviceRequest struct {
	Name        string                    `json:"name" `
	Description string                    `json:"description"`
	Host        string                    `json:"host"`
	Port        uint                      `json:"port"`
	Protocol    string                    `json:"protocol"`
	NID         uint64                    `json:"nid" binding:"required"`
	UID         uint64                    `json:"uid" binding:"required"`
	VS          models.VmessSetting       `json:"vmessSettings"`
	SS          models.ShadowsocksSetting `json:"shadowsocksSettings"`
}

// Unified response for only one service
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
