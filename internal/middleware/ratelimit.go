package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimiter is a small in-memory fixed-window limiter keyed by client IP.
// It protects login, signup, reset and contact endpoints from brute force and spam.
// State is per process, which is fine for a single instance; use Redis if you scale out.
type RateLimiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	hits   map[string]*bucket
}

type bucket struct {
	count   int
	resetAt time.Time
}

// NewRateLimiter allows `limit` requests per `window` for each client IP.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{limit: limit, window: window, hits: make(map[string]*bucket)}
	go rl.sweep()
	return rl
}

// Allow records a hit and reports whether the key is still within its quota.
func (rl *RateLimiter) Allow(key string, now time.Time) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, ok := rl.hits[key]
	if !ok || now.After(b.resetAt) {
		rl.hits[key] = &bucket{count: 1, resetAt: now.Add(rl.window)}
		return true
	}
	b.count++
	return b.count <= rl.limit
}

func (rl *RateLimiter) sweep() {
	t := time.NewTicker(5 * time.Minute)
	defer t.Stop()
	for now := range t.C {
		rl.mu.Lock()
		for k, b := range rl.hits {
			if now.After(b.resetAt) {
				delete(rl.hits, k)
			}
		}
		rl.mu.Unlock()
	}
}

// Wrap returns a handler that answers 429 once the client exceeds the quota.
func (rl *RateLimiter) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.Allow(ClientIP(r), time.Now()) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":"too many requests, please try again shortly"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ClientIP prefers the first X-Forwarded-For hop (set by Render/Vercel proxies) and falls back to the socket address.
func ClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if first := strings.TrimSpace(strings.Split(xff, ",")[0]); first != "" {
			return first
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
