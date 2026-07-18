package plugin

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type RateLimitPlugin struct{}

type RateLimitStrategy string
type KeyBy string

const (
	TokenBucket RateLimitStrategy = "token_bucket"
	FixedWindow RateLimitStrategy = "fixed_window"
)

const (
	KeyByIP     KeyBy = "ip"
	KeyByAPIKey KeyBy = "api_key"
)

type RateLimitConfig struct {
	Enabled     bool               `json:"enabled"`
	Strategy    RateLimitStrategy  `json:"strategy"`
	KeyBy       KeyBy              `json:"key_by,omitempty"`
	TokenBucket *TokenBucketConfig `json:"token_bucket,omitempty"`
	FixedWindow *FixedWindowConfig `json:"fixed_window,omitempty"`
}

type TokenBucketConfig struct {
	Capacity   int `json:"capacity"`
	RefillRate int `json:"refill_rate"`
}

type FixedWindowConfig struct {
	Limit         int `json:"limit"`
	WindowSeconds int `json:"window_seconds"`
}

func (r *RateLimitPlugin) Name() string {
	return "rate-limit"
}

func (r *RateLimitPlugin) Validate(config map[string]interface{}) error {
	var cfg RateLimitConfig
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return err
	}
	if cfg.Strategy == FixedWindow && cfg.FixedWindow != nil && (cfg.FixedWindow.Limit < 1 || cfg.FixedWindow.WindowSeconds < 1) {
		return errors.New("fixed window config error")
	}
	if cfg.Strategy == TokenBucket && cfg.TokenBucket != nil && (cfg.TokenBucket.Capacity < 1 || cfg.TokenBucket.RefillRate < 1) {
		return errors.New("token bucket config error")
	}
	return nil
}

func (r *RateLimitPlugin) CreateMiddleware(config map[string]interface{}) (func(next http.Handler) http.Handler, error) {
	var cfg RateLimitConfig
	b, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return nil, err
	}
	return func(next http.Handler) http.Handler {
		if !cfg.Enabled {
			return next
		}
		if cfg.Strategy == TokenBucket {
			rl := NewTokenBucketRateLimiter(float64(cfg.TokenBucket.Capacity), float64(cfg.TokenBucket.RefillRate), cfg.KeyBy)
			return rl.Middleware(next)
		} else if cfg.Strategy == FixedWindow {
			rl := NewFixedWindowRateLimiter(cfg.FixedWindow.Limit, time.Duration(cfg.FixedWindow.WindowSeconds)*time.Second, cfg.KeyBy)
			return rl.Middleware(next)
		}
		return next
	}, nil
}

// Implemnetation of ratelimiting
// diff files
