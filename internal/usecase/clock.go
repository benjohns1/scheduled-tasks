package usecase

import "time"

// Clock describes the time functions needed by the scheduler
type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
	Sleep(d time.Duration)
	Tick(d time.Duration) <-chan time.Time
	Since(t time.Time) time.Duration
	Until(t time.Time) time.Duration
}
