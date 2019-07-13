package test

import (
	"fmt"
	"testing"
)

type loggerStub struct{}

func (l *loggerStub) Printf(format string, v ...interface{}) {
	if testing.Verbose() {
		fmt.Printf(fmt.Sprintf("    LOG: %v\n", format), v...)
	}
}
