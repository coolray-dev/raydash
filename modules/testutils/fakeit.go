package testutils

import "github.com/brianvoe/gofakeit/v5"

func FakePassword() string {
	return gofakeit.Password(true, true, true, true, true, 8)
}
