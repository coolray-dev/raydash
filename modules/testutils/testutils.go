package testutils

import (
	"database/sql"

	"github.com/brianvoe/gofakeit/v5"
	v1 "github.com/coolray-dev/raydash/api/v1"
	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/setting"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Setup add hook to gorm:create and record all insertion, return a teardown function that clean all test record and remove the hook, refer from https://jarifibrahim.github.io/blog/test-cleanup-with-gorm-hooks/
func Setup() func() {
	type Entity struct {
		Table string
		Key   string
		Value interface{}
	}

	gofakeit.Seed(0)
	models.Migrate()

	var entities []Entity
	var hookName string

	hookName = "RecordTestAndClean"

	orm.DB.Callback().Create().After("gorm:create").Register(hookName, func(scope *gorm.Scope) {
		entities = append(entities, Entity{
			Table: scope.TableName(),
			Key:   scope.PrimaryKey(),
			Value: scope.PrimaryKeyValue(),
		})
	})

	return func() {
		// Remove hook
		defer orm.DB.Callback().Create().Remove(hookName)
		// Find out if db is in transaction
		_, inTransaction := orm.DB.CommonDB().(*sql.Tx)
		tx := orm.DB
		if !inTransaction {
			tx = orm.DB.Begin()
		}

		// Remove reversed entites
		for i := len(entities) - 1; i >= 0; i-- {
			entity := entities[i]
			tx.Table(entity.Table).Where(entity.Key+"= ?", entity.Value).Delete("")
		}

		if !inTransaction {
			tx.Commit()
		}
	}
}

// GetRouter return a router for test
func GetRouter() *gin.Engine {

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// CORS config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = setting.Config.GetStringSlice("app.frontend")

	v1.SetupRouter(router, &corsConfig)
	return router
}
