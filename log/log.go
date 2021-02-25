package log

import (
	"fmt"
	"log"
)

var (
	_ Logger = (*StdLogger)(nil)
	_ Logger = (*NoOpLogger)(nil)
)

// Logger is responsible for providing logging facilities to bot instance.
type Logger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
}

// StdLogger implements the Logger using standard library log package.
type StdLogger struct{}

func (s StdLogger) Debugf(msg string, args ...interface{}) { stdLog("DEBUG", msg, args...) }
func (s StdLogger) Infof(msg string, args ...interface{})  { stdLog("INFO ", msg, args...) }
func (s StdLogger) Warnf(msg string, args ...interface{})  { stdLog("WARN ", msg, args...) }
func (s StdLogger) Errorf(msg string, args ...interface{}) { stdLog("ERR  ", msg, args...) }

// NoOpLogger implements a Logger that simply ignores all the log entries.
type NoOpLogger struct{}

func (n NoOpLogger) Debugf(string, ...interface{}) {}
func (n NoOpLogger) Infof(string, ...interface{})  {}
func (n NoOpLogger) Warnf(string, ...interface{})  {}
func (n NoOpLogger) Errorf(string, ...interface{}) {}

func stdLog(level, msg string, args ...interface{}) {
	log.Printf("[" + level + "] " + fmt.Sprintf(msg, args...))
}
