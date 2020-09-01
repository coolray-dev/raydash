package models

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/modules/setting"
	"github.com/coolray-dev/raydash/modules/utils"
)

type BaseModel struct {
	ID        uint64    `gorm:"primaryKey" json:"id" fake:"skip"`
	CreatedAt time.Time `json:"created_at" fake:"skip"`
	UpdatedAt time.Time `json:"updated_at" fake:"skip"`
}

func init() {
	migrate()
	if err := orm.DB.Where("name = ?", "admin").First(&Group{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		var adminGroup Group
		adminGroup.Name = "admin"
		if err := orm.DB.Create(&adminGroup).Error; err != nil {
			fmt.Println("Database Seeding Error: ", err.Error())
			os.Exit(1)
		}
	} else if err != nil {
		fmt.Println("Database Seeding Error: ", err.Error())
		os.Exit(1)
	}
	if err := orm.DB.Where("username = ?", "admin").First(&User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		var admin User = User{
			UUID:     uuid.New().String(),
			Username: "admin",
			Password: utils.Hash(setting.Config.GetString("app.adminpassword")),
		}
		var adminGroup Group
		if err := orm.DB.Where("name = ?", "admin").First(&adminGroup).Error; err != nil {
			fmt.Println("Database Seeding Error: ", err)
			os.Exit(1)
		}
		admin.Groups = append(admin.Groups, &adminGroup)
		if err := orm.DB.Create(&admin).Error; err != nil {
			fmt.Println("Database Seeding Error: ", err)
			os.Exit(1)
		}
	}
}

func migrate() {
	orm.DB.AutoMigrate(&Group{},
		&User{},
		&Node{},
		&ForgetPassword{},
		&Option{},
		&Service{},
		&Announcement{})

}
