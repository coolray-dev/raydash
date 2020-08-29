package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"

	v1 "github.com/coolray-dev/raydash/api/v1"
	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/coolray-dev/raydash/modules/mail"
	"github.com/coolray-dev/raydash/modules/setting"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {

	// Setup Log
	setupLog()

	// init DB
	defer orm.DB.Close()
	models.Seed()
	models.Migrate()

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

	// Get bind address from config and setup server
	bindAddr := setting.Config.GetString("app.bind")
	server := &http.Server{
		Addr:    bindAddr,
		Handler: router,
	}

	// Create channel to catch system signal
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt)

	// Monitor signal from channel sigs
	go func() {
		sig := <-sigs

		// Do graceful shutdown
		fmt.Println("Received signal", sig)
		fmt.Println("Shutting Down")
		fmt.Print("Stopping MailWorker...")
		mailWorker.Stop()
		fmt.Println("Done")
		fmt.Print("Stopping Gin...")
		if err := server.Close(); err != nil {
			log.Log.Error("Server Close:", err)
		}
		fmt.Println("Done")
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Log.Info("HTTP Server closed")
		} else {
			log.Log.WithError(err).Fatal("HTTP Server closed unexpect")
		}
	}

	// wait for all goroutine to exit
	wg.Wait()
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
