package models

// Announcement is used by admin to annouce sth
type Announcement struct {
	BaseModel
	Level   string `json:"level"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
