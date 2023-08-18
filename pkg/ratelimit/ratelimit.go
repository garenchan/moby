package ratelimit

import (
	"sync"
	"time"

	"github.com/juju/ratelimit"
)

// Limiter rate limiter
type Limiter struct {
	mu     sync.RWMutex
	rate   float64
	clock  ratelimit.Clock
	bucket *ratelimit.Bucket
}

// NewLimiter returns a Limiter that will limit to the given RPS.
func NewLimiter(rate float64) *Limiter {
	limiter := &Limiter{}
	limiter.clock = realClock{}

	limiter.setRate(rate)

	return limiter
}

// Take takes count tokens from the bucket without blocking. It returns
// the time that the caller should wait until the tokens are actually
// available.
func (rl *Limiter) Take(count int64) time.Duration {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	if rl.bucket == nil {
		return 0
	}

	return rl.bucket.Take(count)
}

// Wait takes count tokens from the bucket, waiting until they are available.
func (rl *Limiter) Wait(count int64) {
	if d := rl.Take(count); d > 0 {
		rl.clock.Sleep(d)
	}
}

// GetRate get rate
func (rl *Limiter) GetRate() float64 {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return rl.rate
}

// SetRate set rate with lock
func (rl *Limiter) SetRate(rate float64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.setRate(rate)
}

// setRate set rate
func (rl *Limiter) setRate(rate float64) {
	rl.rate = rate

	// When rate is less than or equal to 0, it means that there is no limit.
	if rate > 0 {
		rl.bucket = ratelimit.NewBucketWithRate(rate, int64(rate))
	} else {
		rl.bucket = nil
	}
}
