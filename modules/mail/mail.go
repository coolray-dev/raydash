package mail

import (
	"crypto/tls"
	"errors"
	"sync"
	"time"

	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"gopkg.in/gomail.v2"
)

// MailChan is a public channel that recieve mail from other modules
var MailChan chan *models.Mail

// Worker is a mail handler
type Worker struct {
	host          string
	port          int
	username      string
	password      string
	allowInsecure bool
	client        *gomail.Dialer
	MailChannel   chan *models.Mail
	WaitGroup     *sync.WaitGroup
}

// NewWorker returns a Worker instance
func NewWorker(config *models.MailConfig, task chan *models.Mail, wg *sync.WaitGroup) *Worker {
	var worker Worker
	worker.config(config)
	worker.MailChannel = task
	worker.WaitGroup = wg
	return &worker
}

// Start starts a worker instance
func (w *Worker) Start() {
	w.init()
	w.WaitGroup.Add(1)
	go w.startWorker()
	log.Log.Info("MailWorker Started")
	return
}

// Stop stops a worker instance
func (w *Worker) Stop() {
	close(w.MailChannel)
	w.WaitGroup.Done()
	return
}

func (w *Worker) config(c *models.MailConfig) {
	w.host = c.Host
	w.port = c.Port
	w.username = c.Username
	w.password = c.Password
	w.allowInsecure = c.AllowInsecure
}

func (w *Worker) init() {
	w.client = gomail.NewDialer(w.host, w.port, w.username, w.password)

	if w.allowInsecure {
		w.client.TLSConfig = &tls.Config{InsecureSkipVerify: w.allowInsecure}
	}
	return
}

func (w *Worker) startWorker() {
	for mail := range w.MailChannel {
		time.Sleep(5 * time.Second)
		if err := w.send(mail); err != nil {
			log.Log.WithError(err).Error("Error Sending Email")
		}
	}
}

func (w *Worker) send(mail *models.Mail) error {
	if w.client == nil {
		return errors.New("Mail client has not been initialized")
	}

	message := gomail.NewMessage()

	message.SetHeader("From", mail.From)
	message.SetHeader("To", mail.To)
	message.SetHeader("Subject", mail.Subject)
	message.SetBody(mail.ContentType, mail.Content)

	if err := w.client.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
