package clock

import "time"

// Time describes the time functions needed by the scheduler
type Time interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
	Sleep(d time.Duration)
	Tick(d time.Duration) <-chan time.Time
	Since(t time.Time) time.Duration
	Until(t time.Time) time.Duration
}

// Clock is a concrete implementation of standard time functions for production usage
type Clock struct{}

// New creates a new instance of a Clock
func New() *Clock {
	return &Clock{}
}

// Now wraps the standard time function
func (c *Clock) Now() time.Time {
	return time.Now()
}

// After wraps the standard time function
func (c *Clock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// Sleep wraps the standard time function
func (c *Clock) Sleep(d time.Duration) {
	time.Sleep(d)
}

// Tick wraps the standard time function
func (c *Clock) Tick(d time.Duration) <-chan time.Time {
	return time.Tick(d)
}

// Since wraps the standard time function
func (c *Clock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// Until wraps the standard time function
func (c *Clock) Until(t time.Time) time.Duration {
	return time.Until(t)
}
