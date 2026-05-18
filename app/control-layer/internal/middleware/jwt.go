package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func InitAuth(secret string) {
	tokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

func Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(tokenAuth)
}

func Authenticator() func(http.Handler) http.Handler {
	return jwtauth.Authenticator(tokenAuth)
}

func SignJwt(username string) (string, error) {
	if tokenAuth == nil {
		return "", errors.New("auth not initialized")
	}

	now := time.Now()

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"user_id": username,
		"exp":     now.Add(24 * time.Hour).Unix(),
		"iat":     now.Unix(),
	})

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
