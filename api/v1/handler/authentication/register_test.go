package authentication_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/testutils"
	"github.com/coolray-dev/raydash/modules/utils"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/google/uuid"
	assertlib "github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	teardown := testutils.Setup()
	defer teardown()
	router := testutils.GetRouter()

	// Create a fake user for testing
	var user models.User
	gofakeit.Struct(&user)

	orm.DB.Create(&user)

	cases := []struct {
		Name     string
		Email    string
		Username string
		Password string
		Status   int
	}{
		{
			"Normal Register",
			gofakeit.Email(),
			gofakeit.Username(),
			"password",
			http.StatusCreated,
		},
		{
			"With existing username",
			gofakeit.Email(),
			user.Username,
			"password",
			http.StatusConflict,
		},
		{
			"With empty username",
			gofakeit.Email(),
			"",
			"password",
			http.StatusBadRequest,
		},
		{
			"With empty email",
			"",
			gofakeit.Username(),
			testutils.FakePassword(),
			http.StatusBadRequest,
		},
		{
			"With invalid email",
			"testreg05test.org",
			gofakeit.Username(),
			testutils.FakePassword(),
			http.StatusBadRequest,
		},
		{
			"With empty password",
			gofakeit.Email(),
			gofakeit.Username(),
			"",
			http.StatusBadRequest,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			assert := assertlib.New(t)

			body := map[string]interface{}{
				"email":    c.Email,
				"username": c.Username,
				"password": c.Password,
			}
			bodyjson, _ := json.Marshal(body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(bodyjson)))

			router.ServeHTTP(w, req)

			type Response struct {
				ErrorMessage string       `json:"error"`
				User         *models.User `json:"user"`
			}
			var response Response
			resParseErr := json.Unmarshal([]byte(w.Body.String()), &response)

			assert.Equal(c.Status, w.Code)
			assert.Nil(resParseErr)

			if c.Status == http.StatusCreated {
				_, UUIDParseErr := uuid.Parse(response.User.UUID)

				assert.Nil(UUIDParseErr)
				assert.Equal(c.Username, response.User.Username)
				assert.Equal(c.Email, response.User.Email)

				var user models.User
				assert.False(
					orm.DB.Where("username = ?", c.Username).
						Where("email = ?", c.Email).
						Where("password = ?", utils.Hash(c.Password)).
						First(&user).
						RecordNotFound(),
				)

				_, UUIDParseErr = uuid.Parse(user.UUID)
				assert.Nil(UUIDParseErr)
			} else {
				assert.Nil(response.User)
			}

			return
		})
	}

	return
}
