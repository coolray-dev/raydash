package authentication_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/coolray-dev/raydash/modules/mail"
	"github.com/coolray-dev/raydash/modules/setting"
	"github.com/coolray-dev/raydash/modules/testutils"
	assertlib "github.com/stretchr/testify/assert"
)

func TestForgetpassword(t *testing.T) {

	router := testutils.GetRouter()

	// Create a fake user for testing
	var user models.User
	gofakeit.Struct(&user)

	orm.DB.Create(&user)

	mail.MailChan = make(chan *models.Mail, 5)

	// init mail
	var wg sync.WaitGroup
	var mailCfg models.MailConfig
	if err := setting.Config.Unmarshal(setting.Config.Sub("mail")); err != nil {
		log.Log.WithError(err).Fatal("Error Unmarshalling Mail Config")
	}
	mailWorker := mail.NewWorker(&mailCfg, mail.MailChan, &wg)
	mailWorker.Start()

	cases := []struct {
		Name   string
		Email  string
		Status int
	}{
		{
			"Normal forget",
			user.Email,
			http.StatusOK,
		},
		{
			"With non-existing email",
			gofakeit.Email(),
			http.StatusOK,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			assert := assertlib.New(t)

			body := map[string]string{
				"email": c.Email,
			}
			bodyjson, _ := json.Marshal(body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/v1/password/forget", bytes.NewBuffer([]byte(bodyjson)))

			router.ServeHTTP(w, req)

			type Response struct {
				ErrorMessage string `json:"error"`
				Message      string `json:"message"`
			}
			var response Response
			resParseErr := json.Unmarshal([]byte(w.Body.String()), &response)

			assert.Nil(resParseErr)
			assert.Equal(c.Status, w.Code)

			if c.Status == http.StatusOK {
				var fp models.ForgetPassword
				err := orm.DB.Where("user_id = ?", user.ID).First(&fp).Error
				assert.Nil(err)
			}
		})
	}

	mailWorker.Stop()
	wg.Wait()
}
