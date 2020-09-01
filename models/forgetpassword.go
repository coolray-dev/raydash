package models

type ForgetPassword struct {
	BaseModel
	Token  string `json:"token" fake:"{uuid}"`
	UserID uint   `fake:"skip"`
	User   *User
}
