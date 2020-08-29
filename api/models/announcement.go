package models

import (
	orm "github.com/coolray-dev/raydash/api/database"
)

// Announcement is used by admin to annouce sth
type Announcement struct {
	BaseModel
	Level   string `json:"level"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func init() {
	orm.DB.AutoMigrate(&Announcement{})
}
