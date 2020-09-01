package announcements_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	orm "github.com/coolray-dev/raydash/database"
	"github.com/coolray-dev/raydash/models"
	"github.com/coolray-dev/raydash/modules/casbin"
	"github.com/coolray-dev/raydash/modules/testutils"
	assertlib "github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestIndex(t *testing.T) {
	testutils.Setup()

	router := testutils.GetRouter()

	var user models.User
	gofakeit.Struct(&user)
	orm.DB.Save(&user)
	casbin.AddDefaultUserPolicy(&user)

	cases := []struct {
		Name   string
		Token  string
		Status int
	}{
		{
			"Index with authorized user",
			testutils.SignAccessToken(&user),
			http.StatusOK,
		},
		{
			"Index without authorized user",
			"",
			http.StatusForbidden,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			assert := assertlib.New(t)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/v1/announcements", nil)

			if c.Token != "" {
				req.Header.Add("Authorization", "Bearer "+c.Token)
			}

			router.ServeHTTP(w, req)

			assert.Equal(c.Status, w.Code)

			if c.Status%100 == 2 {
				var response map[string]interface{}
				err := json.Unmarshal([]byte(w.Body.String()), &response)

				assert.Nil(err)
			}

		})
	}
}
