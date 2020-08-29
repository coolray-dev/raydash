package authentication

import (
	"fmt"
	"time"

	orm "github.com/coolray-dev/raydash/api/database"
	model "github.com/coolray-dev/raydash/api/models"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/google/uuid"
)

// TokenPayload is a convented JWT Payload struct
type TokenPayload struct {
	jwt.Payload
	UID      uint64 `json:"uid"`
	Username string `json:"username"`
}

// Verify validate token with given key
func Verify(token []byte, key []byte) (*TokenPayload, error) {

	var plain TokenPayload

	// Register expiration and issuedAt validator
	now := time.Now()
	iatValidator := jwt.IssuedAtValidator(now)
	expValidator := jwt.ExpirationTimeValidator(now)
	validatePayload := jwt.ValidatePayload(&plain.Payload, iatValidator, expValidator)

	// Registor Algorithm
	hs := jwt.NewHS512(key)

	// Call jwt package to verify
	_, err := jwt.Verify(token, hs, &plain, validatePayload)
	if err != nil {
		return nil, fmt.Errorf("Token verifying error: %w", err)
	}
	return &plain, nil
}

// SignRefreshToken signs a refresh token of a user
func SignRefreshToken(user *model.User) (token string, err error) {
	var key []byte
	key, err = user.GetJwtKey()
	if err != nil {
		return
	}
	var hs = jwt.NewHS512(key)
	now := time.Now()
	plain := TokenPayload{ // defined in jwthelper.go
		Payload: jwt.Payload{
			Issuer:         "RayDash",
			Subject:        "RefreshToken",
			Audience:       jwt.Audience{},
			ExpirationTime: jwt.NumericDate(now.Add(24 * time.Hour)), // Hard code 1 day refreshToken expire time for now
			NotBefore:      jwt.NumericDate(now),
			IssuedAt:       jwt.NumericDate(now),
			JWTID:          uuid.New().String(),
		},
		UID:      user.ID,
		Username: user.Username,
	}
	var tokenb []byte
	tokenb, err = jwt.Sign(plain, hs)
	token = string(tokenb)
	if user.Token == nil {
		user.Token = make(map[string]time.Time)
	}
	user.Token[token] = now.Add(24 * time.Hour)
	if err := orm.DB.Save(&user).Error; err != nil {
		return "", fmt.Errorf("Database error: %w", err)
	}
	return token, err
}

// SignAccessToken signs a access token of a user
func SignAccessToken(user *model.User) (token string, err error) {
	var key []byte
	key, err = user.GetJwtKey()
	if err != nil {
		return
	}
	var hs = jwt.NewHS512(key)
	now := time.Now()
	plain := TokenPayload{ // defined in jwthelper
		Payload: jwt.Payload{
			Issuer:         "RayDash",
			Subject:        "AccessToken",
			Audience:       jwt.Audience{},
			ExpirationTime: jwt.NumericDate(now.Add(5 * time.Minute)), // Hard code 5 min accessToken expire time for now
			NotBefore:      jwt.NumericDate(now),
			IssuedAt:       jwt.NumericDate(now),
			JWTID:          uuid.New().String(),
		},
		UID:      user.ID,
		Username: user.Username,
	}
	var tokenb []byte
	tokenb, err = jwt.Sign(plain, hs)
	token = string(tokenb)
	return token, err
}
