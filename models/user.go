package models

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"

	orm "github.com/coolray-dev/raydash/database"
	"gorm.io/gorm"
)

// User table model
type User struct {
	BaseModel
	UUID              string               `gorm:"unique" json:"uuid" fake:"{uuid}"`
	Token             map[string]time.Time `json:"-" gorm:"-" fake:"skip"`            // gorm doesn't support complex type so hav to marshal it
	TokenStr          string               `json:"-" gorm:"column:token" fake:"skip"` // the actual data is stored here
	JwtKey            []byte               `json:"-" fake:"skip"`                     // Do not export it due to leak risk
	Email             string               `json:"email" fake:"{email}"`
	Username          string               `gorm:"unique" json:"username" fake:"{username}"`
	Password          string               `json:"-" fake:"{password:true,true,true,true,true,8}"`
	SubscriptionToken string               `json:"subscription_token"`
	CurrentTraffic    int64                `json:"current_traffic"`
	MaxTraffic        int64                `json:"max_traffic"`
	Groups            []*Group             `gorm:"many2many:groups_users;" json:"-" fake:"skip"`
}

// GetJwtKey provide access to private var jwtKey, if jwtKey is nil then generate it
func (user *User) GetJwtKey() (key []byte, err error) {
	if user.JwtKey == nil {
		key = make([]byte, 128) // 128 bits key seems to be secure enough for hs512 algorithm
		if _, err := rand.Read(key); err != nil {
			return nil, fmt.Errorf("Jwt Key Generation error: %w", err)
		}
		user.JwtKey = key
		if err := orm.DB.Save(&user).Error; err != nil {
			return nil, fmt.Errorf("Error saving key to db: %w", err)
		}
		return key, nil
	}
	key = user.JwtKey
	return key, nil
}

// BeforeSave marshal the token map
// This is a GORM feature called hook
func (user *User) BeforeSave(*gorm.DB) error {
	if user.Token == nil {
		return nil
	}

	b, err := json.Marshal(&user.Token)
	if err != nil {
		return err
	}

	user.TokenStr = string(b)
	return nil
}

// AfterFind unmarshal token map
func (user *User) AfterFind(*gorm.DB) error {
	if user.TokenStr == "" {
		return nil
	}

	return json.Unmarshal([]byte(user.TokenStr), &user.Token)
}
