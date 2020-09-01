package testutils

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/brianvoe/gofakeit/v5"
	v1 "github.com/coolray-dev/raydash/api/v1"
	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/modules/setting"
)

// Setup add hook to gorm:create and record all insertion, return a teardown function that clean all test record and remove the hook, refer from https://jarifibrahim.github.io/blog/test-cleanup-with-gorm-hooks/
func Setup() (*gorm.DB, func(*gorm.DB)) {
	tx := orm.DB.Begin()
	tx.SavePoint("Origin")
	gofakeit.Seed(0)
	return tx, func(tx *gorm.DB) {
		tx.RollbackTo("Origin")
		orm.DB = tx.Commit()
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
