package database

import (
	"fmt"
	"os"

	"github.com/coolray-dev/raydash/modules/setting"
	"github.com/coolray-dev/raydash/modules/utils"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is a instance of database connection
var DB *gorm.DB

func init() {

	var err error

	// Get Config & Connect
	switch setting.Config.GetString("database.type") {
	case "sqlite3":
		DB, err = gorm.Open(sqlite.Open(utils.AbsPath(setting.Config.GetString("database.path"))), &gorm.Config{})

	case "mysql":
		dsn := setting.Config.GetString("database.username") +
			":" +
			setting.Config.GetString("database.password") +
			"@tcp(" +
			setting.Config.GetString("database.host") +
			")/" +
			setting.Config.GetString("database.dbname") +
			"?charset=utf8mb4"
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}

	// Handle Error Now
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if DB.Error != nil {
		fmt.Printf("database error %v", DB.Error)
	}
}
