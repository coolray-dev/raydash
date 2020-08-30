package clash

import (
	"bytes"
	"text/template"

	"github.com/coolray-dev/raydash/modules/utils"

	"github.com/coolray-dev/raydash/models"
)

type clashNode struct {
	Name      string
	Type      string
	Server    string
	Port      uint
	UUID      string
	AlterID   uint `yaml:"alterId"`
	Cipher    string
	UDP       bool
	Network   string
	WSPath    string
	WSHeaders map[string]string
}

func Generate(services []*models.Service, b *bytes.Buffer) string {
	var nodeYAML []clashNode
	for _, s := range services {
		nodeYAML = append(nodeYAML, *convert(s))
	}
	tmpl, err := template.ParseFiles(utils.AbsPath("template/clash.tmpl"))
	if err != nil {
		panic(err)
	}
	tmpl.Execute(b, nodeYAML)
	return ""
}

func convert(s *models.Service) *clashNode {
	var node clashNode
	node.Name = s.Name
	node.Type = "vmess"
	node.Server = s.Host
	node.Port = s.Port
	node.UUID = s.UUID
	node.AlterID = s.AlterID
	node.Cipher = "auto"
	node.UDP = true
	node.Network = s.TransportProtocol
	node.WSPath = "/"
	return &node
}
