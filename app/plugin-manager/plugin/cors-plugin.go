package plugin

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type CorsPlugin struct{}

type CorsConfig struct {
	Enabled        bool     `json:"enabled"`
	AllowedOrigins []string `json:"allowed_origins" bson:"allowed_origins"`
	AllowedMethods []string `json:"allowed_methods" bson:"allowed_methods"`
	AllowedHeaders []string `json:"allowed_headers" bson:"allowed_headers"`
}

func (c *CorsPlugin) Name() string {
	return "cors"
}

func (c *CorsPlugin) Validate(config map[string]interface{}) error {
	var cfg CorsConfig
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(b), &cfg)
	if err != nil {
		return err
	}
	if len(cfg.AllowedOrigins) == 0 {
		return errors.New("allowed_origins required")
	}
	return nil
}

func (c *CorsPlugin) CreateMiddleware(config map[string]interface{}) (func(next http.Handler) http.Handler, error) {
	var cfg CorsConfig
	b, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(b), &cfg)
	if err != nil {
		return nil, err
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !cfg.Enabled {
				next.ServeHTTP(w, r)
				return
			}
			origin := r.Header.Get("Origin")

			for _, allowedOrigin := range cfg.AllowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
					break
				}
			}
			w.Header().Set(
				"Access-Control-Allow-Methods",
				strings.Join(cfg.AllowedMethods, ","),
			)
			w.Header().Set(
				"Access-Control-Allow-Headers",
				strings.Join(cfg.AllowedHeaders, ","),
			)
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}, nil
}
