package clock

import "time"

// Mock provides a standard testing mock for a clock instance
type Mock struct {
	now func() time.Time
}

// NewMock creates a new clock mock that uses the 'now' function to generate Now() responses
func NewMock(now func() time.Time) *Mock {
	return &Mock{now: now}
}

// NewStaticMock creates a new clock mock who's Now() function always returns the 'now'
func NewStaticMock(now time.Time) *Mock {
	return &Mock{now: func() time.Time { return now }}
}

// Now implementes the standard time function
func (c *Mock) Now() time.Time {
	return c.now()
}

// After implementes the standard time function
func (c *Mock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// Sleep implementes the standard time function
func (c *Mock) Sleep(d time.Duration) {
	time.Sleep(d)
}

// Since implementes the standard time function
func (c *Mock) Since(t time.Time) time.Duration {
	return c.Now().Sub(t)
}

// Until implementes the standard time function
func (c *Mock) Until(t time.Time) time.Duration {
	return t.Sub(c.Now())
}
