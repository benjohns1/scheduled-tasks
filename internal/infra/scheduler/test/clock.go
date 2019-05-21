package test

import "time"

type ClockMock struct {
	mockNow func() time.Time
}

func NewClockMock(now func() time.Time) *ClockMock {
	return &ClockMock{mockNow: now}
}

func NewStaticClockMock(now time.Time) *ClockMock {
	return &ClockMock{mockNow: func() time.Time { return now }}
}

// Now implementes the standard time function
func (c *ClockMock) Now() time.Time {
	return c.mockNow()
}

// After implementes the standard time function
func (c *ClockMock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// Sleep implementes the standard time function
func (c *ClockMock) Sleep(d time.Duration) {
	time.Sleep(d)
}

// Tick implementes the standard time function
func (c *ClockMock) Tick(d time.Duration) <-chan time.Time {
	return time.Tick(d)
}

// Since implementes the standard time function
func (c *ClockMock) Since(t time.Time) time.Duration {
	return c.Now().Sub(t)
}

// Until implementes the standard time function
func (c *ClockMock) Until(t time.Time) time.Duration {
	return t.Sub(c.Now())
}
