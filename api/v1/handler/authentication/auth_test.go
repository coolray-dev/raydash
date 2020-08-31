package authentication_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/coolray-dev/raydash/modules/casbin"

	"github.com/brianvoe/gofakeit/v5"
	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/setting"
	"github.com/coolray-dev/raydash/modules/testutils"
	assertlib "github.com/stretchr/testify/assert"
)

// Create a fake user for testing
var user models.User

func TestMain(m *testing.M) {

	tx, teardown := testutils.Setup()
	defer teardown(tx)

	gofakeit.Struct(&user)

	orm.DB.Create(&user)
	casbin.AddDefaultUserPolicy(&user)
	code := m.Run()
	os.Exit(code)
}

func TestLogin(t *testing.T) {
	router := testutils.GetRouter()

	cases := []struct {
		Name     string
		Username string
		Password string
		Code     int
	}{
		{"Correct", "admin", setting.Config.GetString("app.adminpassword"), http.StatusOK},
		{"WrongPassword", "admin", testutils.FakePassword(), http.StatusUnauthorized},
		{"EmptyPassword", gofakeit.Username(), "", http.StatusUnauthorized},
		{"BothEmpty", "", "", http.StatusUnauthorized},
		{"EmptyUsername", "", testutils.FakePassword(), http.StatusUnauthorized},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			assert := assertlib.New(t)

			body := map[string]interface{}{
				"username": c.Username,
				"password": c.Password,
			}
			bodyjson, _ := json.Marshal(body)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/v1/login", bytes.NewBuffer([]byte(bodyjson)))
			router.ServeHTTP(w, req)
			var response map[string]string
			err := json.Unmarshal([]byte(w.Body.String()), &response)
			assert.Nil(err)
			_, AccessTokenExists := response["access_token"]
			_, RefreshTokenExists := response["refresh_token"]
			if w.Code == http.StatusOK {
				assert.True(AccessTokenExists)
				assert.True(RefreshTokenExists)
				return
			}
			assert.False(AccessTokenExists)
			assert.False(RefreshTokenExists)
			assert.Equal(http.StatusUnauthorized, w.Code)
			return
		})
	}
	return
}

func TestLogout(t *testing.T) {

	router := testutils.GetRouter()

	refreshToken := testutils.SignRefreshToken(&user)

	cases := []struct {
		Name         string
		AuthToken    string
		RefreshToken string
		Status       int
	}{
		{
			"Normal logout with access token",
			testutils.SignAccessToken(&user),
			testutils.SignRefreshToken(&user),
			http.StatusNoContent,
		},
		{
			"Nomal logout with refresh token",
			refreshToken,
			refreshToken,
			http.StatusForbidden,
		},
		{
			"Logout with non-existing refresh token in body",
			testutils.SignAccessToken(&user),
			gofakeit.Word(),
			http.StatusNotFound,
		},
		{
			"Logout with no refresh token provided",
			testutils.SignAccessToken(&user),
			"",
			http.StatusBadRequest,
		},
		{
			"Logout without auth header",
			"",
			testutils.SignRefreshToken(&user),
			http.StatusForbidden,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			assert := assertlib.New(t)

			body := make(map[string]string)

			if c.RefreshToken != "" {
				body["refresh_token"] = c.RefreshToken
			}
			bodyjson, _ := json.Marshal(body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/v1/logout", bytes.NewBuffer([]byte(bodyjson)))

			if c.AuthToken != "" {
				req.Header.Add("Authorization", "Bearer "+c.AuthToken)
			}

			router.ServeHTTP(w, req)

			assert.Equal(c.Status, w.Code)

			if c.Status%100 == 2 { // if status_code starts with 2, which means success
				orm.DB.Where("id = ?", user.ID).First(&user)
				_, tokenExists := user.Token[c.RefreshToken]
				assert.False(tokenExists)
			}
		})
	}
}
