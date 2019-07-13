package test

import (
	"time"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/clock"
	format "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/json"
)

// SetStaticClock sets the clock to a static time for testing
func SetStaticClock(now time.Time) (nowStr string, reset func()) {
	nowStr = now.Format(format.OutTimeFormat)
	prevClock := clock.Get()
	clockMock := clock.NewStaticMock(now)
	clock.Set(clockMock)
	return nowStr, func() {
		clock.Set(prevClock)
	}
}
