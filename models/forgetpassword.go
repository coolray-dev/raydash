package models

import (
	orm "github.com/coolray-dev/raydash/database"
)

type ForgetPassword struct {
	BaseModel
	Token  string `json:"token" fake:"{uuid}"`
	UserID uint   `fake:"skip"`
	User   *User
}

func init() {
	orm.DB.AutoMigrate(&ForgetPassword{})
}
