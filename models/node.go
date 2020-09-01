package models

// Node is a struct of node info
type Node struct {
	BaseModel
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Groups         []*Group   `gorm:"many2many:groups_nodes;" json:"-"`
	Services       []*Service `json:"-"`
	Host           string     `json:"host"` // The Host to access v2ray
	Ports          string     `json:"ports"`
	AccessToken    string     `json:"access_token,omitempty"`
	CurrentTraffic uint64     `json:"current_traffic"`
	MaxTraffic     uint64     `json:"max_traffic"`
	HasUDP         bool       `json:"hasUDP"`
	HasMultiPort   bool       `json:"hasMultiPort"`
	Settings       `json:"settings"`
}

type Settings struct {
	Listen             string `json:"listen"`
	Port               uint   `json:"port"`
	VmessSetting       `json:"vmessSettings"`
	ShadowsocksSetting `json:"shadowsocksSettings"`
}
