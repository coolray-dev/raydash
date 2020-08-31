package models

type Option struct {
	BaseModel
	Name  string `json:"name"`
	Value string `json:"value"`
}
