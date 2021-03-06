package authentication_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/testutils"
	"github.com/coolray-dev/raydash/modules/utils"
	assertlib "github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestResetpassword(t *testing.T) {

	router := testutils.GetRouter()

	// Create a fake forget-password record for testing
	var fp1 models.ForgetPassword
	gofakeit.Struct(&fp1)
	fp1.User.UUID = gofakeit.UUID()
	orm.DB.Create(&fp1)

	var fp2 models.ForgetPassword
	gofakeit.Struct(&fp2)
	fp2.User.UUID = gofakeit.UUID()
	orm.DB.Create(&fp2)

	//defer orm.DB.Rollback()

	cases := []struct {
		Name     string
		Token    string
		Record   *models.ForgetPassword
		Password string
		Status   int
	}{
		{
			"Normal reset",
			fp1.Token,
			&fp1,
			testutils.FakePassword(),
			http.StatusNoContent,
		},
		{
			"With non-existing uuid token",
			gofakeit.UUID(),
			nil, // useless due to no matching user
			testutils.FakePassword(),
			http.StatusNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			assert := assertlib.New(t)

			body := map[string]string{
				"token":    c.Token,
				"password": c.Password,
			}
			bodyjson, _ := json.Marshal(body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/v1/password/reset", bytes.NewBuffer([]byte(bodyjson)))

			router.ServeHTTP(w, req)

			var response map[string]interface{}

			assert.Equal(c.Status, w.Code)

			if c.Status == http.StatusOK {
				assert.Nil(w.Body)
				assert.True(errors.Is(orm.DB.Where("id = ?", c.Record.ID).First(c.Record).Error, gorm.ErrRecordNotFound))

				var user models.User
				orm.DB.Where("id = ?", c.Record.UserID).First(&user)
				assert.Equal(utils.Hash(c.Password), user.Password)
			} else {
				json.Unmarshal([]byte(w.Body.String()), &response)
				assert.NotEqual("", response["error"])
			}
		})
	}
}
