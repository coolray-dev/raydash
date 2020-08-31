package database

import (
	"fmt"

	"github.com/coolray-dev/raydash/modules/utils"

	setting "github.com/coolray-dev/raydash/modules/setting"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is a instance of database connection
var DB *gorm.DB

func init() {
	var err error
	switch setting.Config.GetString("database.type") {
	case "sqlite3":
		DB, err = gorm.Open(sqlite.Open(utils.AbsPath(setting.Config.GetString("database.path"))), &gorm.Config{})
	}

	if err != nil {
		fmt.Printf("sqlite creation error %v", err)
	}

	if DB.Error != nil {
		fmt.Printf("database error %v", DB.Error)
	}
}
