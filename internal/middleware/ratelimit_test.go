package middleware

import (
	"testing"
	"time"
)

func TestRateLimiterBlocksAfterLimitAndResets(t *testing.T) {
	rl := &RateLimiter{limit: 3, window: time.Minute, hits: make(map[string]*bucket)}
	now := time.Now()
	for i := 0; i < 3; i++ {
		if !rl.Allow("1.2.3.4", now) {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
	if rl.Allow("1.2.3.4", now) {
		t.Fatal("fourth request must be blocked")
	}
	if !rl.Allow("5.6.7.8", now) {
		t.Fatal("a different IP must have its own quota")
	}
	if !rl.Allow("1.2.3.4", now.Add(2*time.Minute)) {
		t.Fatal("quota must reset after the window")
	}
}
