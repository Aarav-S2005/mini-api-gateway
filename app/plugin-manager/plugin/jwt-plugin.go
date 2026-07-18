package plugin

import (
	"crypto"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Source string

const (
	SourceCookie Source = "cookie"
	SourceHeader Source = "header"
)

type JwtPlugin struct{}

type JwtConfig struct {
	Enabled     bool   `json:"enabled"`
	Algorithm   string `json:"algorithm"`
	PublicKey   string `json:"public_key"`
	UserIDClaim string `json:"user_id_claim"`
	TokenSource Source `json:"token_source"`

	HeaderName string `json:"header_name,omitempty"`
	Prefix     string `json:"prefix,omitempty"`

	CookieName string `json:"cookie_name,omitempty"`
}

func (j JwtPlugin) Name() string {
	return "jwt-auth"
}

func (j JwtPlugin) Validate(config map[string]interface{}) error {
	var cfg JwtConfig
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return err
	}
	if cfg.Algorithm == "" || cfg.PublicKey == "" || cfg.UserIDClaim == "" || cfg.TokenSource == "" {
		return errors.New("invalid jwt config")
	}
	if cfg.TokenSource == SourceHeader && (cfg.HeaderName == "" || cfg.Prefix == "") {
		return errors.New("invalid jwt config")
	}
	if cfg.TokenSource == SourceCookie && cfg.CookieName == "" {
		return errors.New("invalid jwt config")
	}
	return nil
}

func (j JwtPlugin) CreateMiddleware(config map[string]interface{}) (func(next http.Handler) http.Handler, error) {
	var cfg JwtConfig
	b, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return nil, err
	}
	var pubKey crypto.PublicKey
	if cfg.Algorithm == "RS256" {
		pubKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.PublicKey))
	} else if cfg.Algorithm == "ES256" {
		pubKey, err = jwt.ParseECPublicKeyFromPEM([]byte(cfg.PublicKey))
	}
	if err != nil {
		return nil, err
	}

	return func(next http.Handler) http.Handler {
		if !cfg.Enabled {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenStr string
			if cfg.TokenSource == SourceCookie {
				cookie, err := r.Cookie(cfg.CookieName)
				if err != nil {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
				tokenStr = cookie.Value
			} else if cfg.TokenSource == SourceHeader {
				auth := r.Header.Get(cfg.HeaderName)
				if !strings.HasPrefix(auth, cfg.Prefix+" ") {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
				tokenStr = strings.TrimPrefix(auth, cfg.Prefix+" ")
			}
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				if token.Method.Alg() != cfg.Algorithm {
					return nil, errors.New("invalid jwt config")
				}
				return pubKey, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			userID, ok := claims[cfg.UserIDClaim].(string)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			r.Header.Set("X-User-Id", userID)
			next.ServeHTTP(w, r)
		})
	}, nil
}
