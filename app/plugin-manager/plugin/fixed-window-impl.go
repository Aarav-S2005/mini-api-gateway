package plugin

import (
	"net/http"
	"sync"
	"time"
)

type Window struct {
	count      int
	windowEnds time.Time
}

type FixedWindowRateLimiter struct {
	requestWindows map[string]*Window
	mu             sync.Mutex
	keyBy          KeyBy
	limit          int
	windowSize     time.Duration
}

func NewFixedWindowRateLimiter(limit int, windowSize time.Duration, keyBy KeyBy) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		requestWindows: make(map[string]*Window),
		limit:          limit,
		windowSize:     windowSize,
		keyBy:          keyBy,
	}
}

func (rl *FixedWindowRateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	w, ok := rl.requestWindows[key]
	if !ok || now.After(w.windowEnds) {
		w = &Window{
			count:      0,
			windowEnds: now.Add(rl.windowSize),
		}
		rl.requestWindows[key] = w
	}
	if w.count < rl.limit {
		w.count++
		return true
	}
	return false
}

func (rl *FixedWindowRateLimiter) Middleware(next http.Handler) http.Handler {
	go func() {
		for {
			time.Sleep(time.Minute)
			rl.mu.Lock()
			now := time.Now()
			for k, v := range rl.requestWindows {
				if now.After(v.windowEnds.Add(time.Minute)) {
					delete(rl.requestWindows, k)
				}
			}
			rl.mu.Unlock()
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := ""
		if rl.keyBy == KeyByIP {
			ip, err := getClientIP(r)
			if err != nil {
				http.Error(w, "could not parse ip", http.StatusBadRequest)
				return
			}
			key = ip
		} else if rl.keyBy == KeyByAPIKey {
			key = r.Header.Get("X-Api-GW-Key")
			if key == "" {
				http.Error(w, "could not parse api key", http.StatusBadRequest)
				return
			}
		}
		if !rl.allow(key) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
