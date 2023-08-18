package ratelimit

import "time"

// realClock implements Clock in terms of standard time functions.
type realClock struct{}

// Now implements Clock.Now by calling time.Now.
func (realClock) Now() time.Time {
	return time.Now()
}

// Sleep implements Clock.Sleep by calling time.Sleep.
func (realClock) Sleep(d time.Duration) {
	time.Sleep(d)
}
