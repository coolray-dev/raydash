package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/modules/setting"
	"github.com/coolray-dev/raydash/modules/utils"
)

type BaseModel struct {
	ID        uint64    `gorm:"primary_key" json:"id" fake:"skip"`
	CreatedAt time.Time `json:"created_at" fake:"skip"`
	UpdatedAt time.Time `json:"updated_at" fake:"skip"`
}

func Seed() {
	if !orm.DB.HasTable(&Group{}) {
		orm.DB.AutoMigrate(&Group{})
		var adminGroup Group
		adminGroup.ID = 1
		adminGroup.Name = "Admin"
		if err := orm.DB.Create(&adminGroup).Error; err != nil {
			fmt.Println("Database Seeding Error: ", err)
		}
	}
	if !orm.DB.HasTable(&User{}) {
		orm.DB.AutoMigrate(&User{})
		var admin User = User{
			UUID:     uuid.New().String(),
			Username: "admin",
			Password: utils.Hash(setting.Config.GetString("app.adminpassword")),
		}
		var adminGroup Group
		if err := orm.DB.Where("ID = ?", "1").First(&adminGroup).Error; err != nil {
			fmt.Println("Database Seeding Error: ", err)
		}
		admin.Groups = append(admin.Groups, &adminGroup)
		if err := orm.DB.Create(&admin).Error; err != nil {
			fmt.Println("Database Seeding Error: ", err)
		}
	}
}

func Migrate() {
	orm.DB.AutoMigrate(&Group{})
	orm.DB.AutoMigrate(&User{})
	orm.DB.AutoMigrate(&Node{})
}
