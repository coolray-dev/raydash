package models

import (
	orm "github.com/coolray-dev/raydash/api/database"
)

type Service struct {
	BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      uint64 `json:"uid"`
	NodeID      uint64 `json:"nid"`

	Host      string `json:"host"`
	Port      uint   `json:"port"`
	Protocol  string `json:"protocol"`
	VmessUser `json:"vmessUser"`
	VmessSetting
	ShadowsocksSetting
}

type ShadowsocksSetting struct{}

type VmessUser struct {
	Email    string `json:"email"`
	UUID     string `json:"uuid"`
	AlterID  uint   `json:"alterid"`
	Security string `json:"security"`
}
type VmessSetting struct {
	StreamSettings `json:"streamSettings"`
	SniffingSettings
	allocate struct{}
}

type StreamSettings struct {
	TransportProtocol string `json:"protocol"`
}

type SniffingSettings struct{}

func init() {
	orm.DB.AutoMigrate(&Service{})
}
