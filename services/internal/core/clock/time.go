package clock

import "time"

// Time describes the time functions needed by the clock
type Time interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
	Sleep(d time.Duration)
	Since(t time.Time) time.Duration
	Until(t time.Time) time.Duration
}

var clock Time

// Default to using standard system time clock for production
func init() {
	clock = New()
}

// Set allows setting the clock to another implementation (e.g. for testing)
func Set(newClock Time) {
	clock = newClock
}

// Get returns the current clock implementation being used
func Get() Time {
	return clock
}

// Now returns the current time
func Now() time.Time {
	return clock.Now()
}

// After waits for the duration to elapse and then sends the current time on the returned channel
func After(d time.Duration) <-chan time.Time {
	return clock.After(d)
}

// Sleep pauses the current goroutine for at least the duration d
func Sleep(d time.Duration) {
	clock.Sleep(d)
}

// Since returns the time elapsed since t
func Since(t time.Time) time.Duration {
	return clock.Since(t)
}

// Until returns the duration until t
func Until(t time.Time) time.Duration {
	return clock.Until(t)
}
