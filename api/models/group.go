package models

type Group struct {
	BaseModel
	Name        string  `json:"groupname"`
	Description string  `json:"description"`
	Users       []*User `gorm:"many2many:groups_users;" json:"-"`
	Nodes       []*Node `gorm:"many2many:groups_nodes;" json:"-"`
}
