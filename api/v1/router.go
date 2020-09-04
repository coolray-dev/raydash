package v1

import (
	"github.com/coolray-dev/raydash/api/v1/handler/announcements"
	"github.com/coolray-dev/raydash/api/v1/handler/authentication"
	"github.com/coolray-dev/raydash/api/v1/handler/groups"
	"github.com/coolray-dev/raydash/api/v1/handler/nodes"
	"github.com/coolray-dev/raydash/api/v1/handler/options"
	"github.com/coolray-dev/raydash/api/v1/handler/services"
	"github.com/coolray-dev/raydash/api/v1/handler/subscription"
	"github.com/coolray-dev/raydash/api/v1/handler/users"
	"github.com/coolray-dev/raydash/api/v1/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter config router routes and middlewares
func SetupRouter(router *gin.Engine, c *cors.Config) {

	// Middleware must be registered before routes initial
	// Otherwise middlewares will not be used

	// CORS Middleware
	setCORSHeaders(c)
	router.Use(cors.New(*c))

	// Log Middleware
	router.Use(middleware.Log())

	router.Use(middleware.Authorize())

	v1 := router.Group("/v1")

	// Finally Setup Routes
	setupRoutes(v1)
}

// SetupRoutes initialized a gin router to route request
func setupRoutes(router *gin.RouterGroup) {

	usersAPI := router.Group("/users")
	{
		usersAPI.GET("", middleware.ParseParams(), users.Index)
		usersAPI.GET("/:username", users.Show)
		usersAPI.PATCH("/:username", users.Update)
		usersAPI.DELETE("/:username", users.Destroy)
		usersAPI.PATCH("/:username/traffic", users.Traffic)
		usersAPI.GET("/:username/groups", users.Groups)
		usersAPI.GET("/:username/nodes", users.Nodes)
		usersAPI.GET("/:username/services", users.Services)
	}

	router.POST("/register", authentication.Register)
	router.POST("/login", authentication.Login)
	router.DELETE("/logout", authentication.Logout)
	router.POST("/refresh", authentication.RefreshToken)

	passwordAPI := router.Group("/password")
	{
		passwordAPI.POST("/reset", authentication.ResetPassword)
		passwordAPI.POST("/forget", authentication.ForgetPassword)
	}
	nodesAPI := router.Group("/nodes")
	{
		nodesAPI.GET("", nodes.Index)
		nodesAPI.POST("", nodes.Create)
		nodesAPI.GET("/:nid", nodes.Show)
		nodesAPI.PATCH("/:nid", nodes.Update)
		nodesAPI.DELETE("/:nid", nodes.Destroy)
		nodesAPI.GET("/:nid/users", nodes.Users)
		nodesAPI.PATCH("/:nid/users/:username/traffic", nodes.Traffic)
		nodesAPI.GET("/:nid/services", nodes.Services)
		nodesAPI.GET("/:nid/token", nodes.AccessToken)
		nodesAPI.POST("/:nid/token", nodes.GenerateToken)
	}
	servicesAPI := router.Group("/services")
	{
		servicesAPI.GET("", middleware.ParseParams(), services.Index)
		servicesAPI.POST("", services.Store)
		servicesAPI.PATCH("/:sid", services.Update)
		servicesAPI.DELETE("/:sid", services.Destroy)
	}
	groupsAPI := router.Group("/groups")
	{
		groupsAPI.GET("", groups.Index)
		groupsAPI.POST("", groups.Create)
		groupsAPI.GET("/:gid", groups.Show)
		groupsAPI.PATCH("/:gid", groups.Update)
		groupsAPI.DELETE("/:gid", groups.Destroy)
		groupsAPI.GET("/:gid/users", groups.Users)
		groupsAPI.PATCH("/:gid/users", groups.AppendUser)
		groupsAPI.DELETE("/:gid/users/:username", groups.RemoveUser)
	}

	optionsAPI := router.Group("/options")
	{
		optionsAPI.GET("", options.Index)
		optionsAPI.PUT("/:name", options.Update)
	}

	announcementsAPI := router.Group("/announcements")
	{
		announcementsAPI.GET("", middleware.ParseParams(), announcements.Index)
		announcementsAPI.POST("", announcements.Store)
		announcementsAPI.GET("/:aid", announcements.Show)
		announcementsAPI.PATCH("/:aid", announcements.Update)
		announcementsAPI.DELETE("/:aid", announcements.Destroy)
	}
	subscriptionAPI := router.Group("/subscription")
	{
		subscriptionAPI.GET("/clash", subscription.Clash)
	}

	return
}

func setCORSHeaders(config *cors.Config) {
	config.AllowHeaders = append(config.AllowHeaders, "authorization")
}
