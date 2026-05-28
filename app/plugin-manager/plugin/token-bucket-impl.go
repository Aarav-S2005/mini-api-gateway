package plugin

import (
	"math"
	"net/http"
	"sync"
	"time"
)

type Bucket struct {
	Tokens     float64
	LastRefill time.Time
}

type TokenBucketRateLimiter struct {
	buckets    map[string]*Bucket
	mu         sync.Mutex
	capacity   float64
	refillRate float64
	KeyBy      KeyBy
}

func NewTokenBucketRateLimiter(capacity, refillRate float64, keyBy KeyBy) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		buckets:    make(map[string]*Bucket),
		capacity:   capacity,
		refillRate: refillRate,
		KeyBy:      keyBy,
	}
}

func (rl *TokenBucketRateLimiter) Refill(bucket *Bucket) {
	now := time.Now()
	elapsed := now.Sub(bucket.LastRefill).Seconds()
	newTokens := rl.refillRate * elapsed
	bucket.Tokens = math.Min(newTokens+bucket.Tokens, rl.capacity)
	bucket.LastRefill = now
}

func (rl *TokenBucketRateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, ok := rl.buckets[key]
	if !ok {
		bucket = &Bucket{
			Tokens:     rl.capacity,
			LastRefill: time.Now(),
		}
		rl.buckets[key] = bucket
	}
	rl.Refill(bucket)
	if bucket.Tokens >= 1 {
		bucket.Tokens--
		return true
	}
	return false
}

func (rl *TokenBucketRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rl.KeyBy == KeyByIP {
			ip, err := getClientIP(r)
			if err != nil {
				http.Error(w, "could not parse ip", http.StatusBadRequest)
				return
			}
			if !rl.allow(ip) {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		} else if rl.KeyBy == KeyByAPIKey {
			apiKey := r.Header.Get("X-Api-Key")
			if apiKey == "" {
				http.Error(w, "could not parse api key", http.StatusBadRequest)
				return
			}
			if !rl.allow(apiKey) {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		}
	})
}
