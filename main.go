// @title RayDash API
// @version 1.0.0
// @description A Swagger UI For RayDash API

// @license.name GPLv3
// @license.url https://www.gnu.org/licenses/gpl-3.0.html

// @host localhost
// @BasePath /v1
// @query.collection.format multi

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles"

	v1 "github.com/coolray-dev/raydash/api/v1"
	_ "github.com/coolray-dev/raydash/docs"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/coolray-dev/raydash/modules/mail"
	"github.com/coolray-dev/raydash/modules/setting"
)

func main() {

	// Setup Log
	setupLog()

	// waitgroup for goroutine
	var wg sync.WaitGroup

	// init mail
	// Do not use Config.Sub("mail") and Unmarshal
	// Config.Sub do check RAYDASH_MAIL, which is not expected
	var mailCfg *models.MailConfig = &models.MailConfig{
		Host:          setting.Config.GetString("mail.host"),
		Port:          setting.Config.GetInt("mail.port"),
		Username:      setting.Config.GetString("mail.username"),
		Password:      setting.Config.GetString("mail.password"),
		AllowInsecure: setting.Config.GetBool("mail.allowinsecure"),
	}
	mail.MailChan = make(chan *models.Mail, 5)
	mailWorker := mail.NewWorker(mailCfg, mail.MailChan, &wg)
	mailWorker.Start()

	// init router
	router := gin.Default()

	// CORS config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = setting.Config.GetStringSlice("app.frontend")

	v1.SetupRouter(router, &corsConfig)

	// Swagger UI
	router.GET("/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup Swagger
	setupSwagger()

	// Get bind address from config and setup server
	bindPort := setting.Config.GetString("app.port")
	bindAddr := setting.Config.GetString("app.address")
	server := &http.Server{
		Addr:    bindAddr + ":" + bindPort,
		Handler: router,
	}

	// Create channel to catch system signal
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Monitor signal from channel sigs
	go func() {
		wg.Add(1)
		sig := <-sigs

		// Do graceful shutdown
		log.Log.Infof("Received signal %s", sig)
		log.Log.Info("Shutting Down")
		log.Log.Info("Stopping MailWorker")
		mailWorker.Stop()
		log.Log.Info("Stopping Gin")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Log.Errorf("Server Close With Error: %s", err.Error())
		}
		wg.Done()
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed && err != nil {
		log.Log.WithError(err).Fatal("HTTP Server Listen Error")
	}

	// wait for all goroutine to exit
	wg.Wait()
	log.Log.Info("RayDash Shutdown Success")
	return
}

func setupLog() {
	switch setting.Config.GetString("app.log.level") {
	case "debug":
		log.Log.SetLevel(logrus.DebugLevel)
	case "info":
		log.Log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.Log.SetLevel(logrus.WarnLevel)
	case "error":
		log.Log.SetLevel(logrus.ErrorLevel)
	default:
		log.Log.SetLevel(logrus.InfoLevel)
	}

}
