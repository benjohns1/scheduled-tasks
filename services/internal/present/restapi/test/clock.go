package test

import (
	"time"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/clock"
	format "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/json"
)

// SetStaticClock sets the clock to a static time for testing
func SetStaticClock(now time.Time) (nowStr string, reset func()) {
	nowStr = FormatTime(now)
	prevClock := clock.Get()
	clockMock := clock.NewStaticMock(now)
	clock.Set(clockMock)
	return nowStr, func() {
		clock.Set(prevClock)
	}
}

// FormatTime formats a time per the presentation output
func FormatTime(t time.Time) string {
	return t.Format(format.OutTimeFormat)
}
