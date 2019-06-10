// The clock package provides a wrapper around the global time functions to 
// allow for easier testing without needing to inject a clock implementation 
// everywhere
//
// Application Usage:
// Import this package and use the clock.Now(), clock.After(d time.Duration), 
// etc functions instead of the time.Now(), time.After(d time.Duration), etc 
// functions respectively. The standard clock will use the underlying time 
// functions.
// 
// Testing Usage:
// Create a mock clock (either the basic one provided in mock.go, or roll your
// own that implements clock.Time) and call clock.Set(mockClock) before running
// your tests. For instance, this will create a static clock who's Now() 
// function will always return 2000-01-01 12:00:00 UTC:
// 
// testNow := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
// clockMock := clock.NewStaticMock(testNow)
// clock.Set(clockMock)
// 
// You can use clock.Get() and defer clock.Set() to push and pop clock
// implementations for nested tests:
// 
// prevClock := clock.Get()
// clockMock := clock.NewStaticMock(time.Now())
// clock.Set(clockMock)
// defer clock.Set(prevClock)
// ... run outer test ...
//   ... run nested test with its own clock set/get ...
package clock

import "time"

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
