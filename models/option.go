package models

import (
	orm "github.com/coolray-dev/raydash/database"
)

type Option struct {
	BaseModel
	Name  string `json:"name"`
	Value string `json:"value"`
}

func init() {
	orm.DB.AutoMigrate(&Option{})
}
