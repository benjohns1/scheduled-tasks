package scheduler

import "time"

// Clock is a concrete implementation of standard time functions to inject into the scheduler for production
type Clock struct{}

// Now implementes the standard time function
func (c *Clock) Now() time.Time {
	return time.Now()
}

// After implementes the standard time function
func (c *Clock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// Sleep implementes the standard time function
func (c *Clock) Sleep(d time.Duration) {
	time.Sleep(d)
}

// Tick implementes the standard time function
func (c *Clock) Tick(d time.Duration) <-chan time.Time {
	return time.Tick(d)
}

// Since implementes the standard time function
func (c *Clock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// Until implementes the standard time function
func (c *Clock) Until(t time.Time) time.Duration {
	return time.Until(t)
}
