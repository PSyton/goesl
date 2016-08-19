// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

// LoggerInterface base logger interface
type LoggerInterface interface {
	Debug(message string, args ...interface{})
	Error(message string, args ...interface{})
	Notice(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warning(message string, args ...interface{})
}

var (
	logger localLogger
)

// SetLogger set global library logger
func SetLogger(l LoggerInterface) {
	logger.impl = l
}

type localLogger struct {
	LoggerInterface
	impl LoggerInterface
}

func (l *localLogger) isValid() bool {
	return l.impl != nil
}

func (l *localLogger) Debug(message string, args ...interface{}) {
	if l.isValid() {
		l.impl.Debug(message, args)
	}
}

func (l *localLogger) Error(message string, args ...interface{}) {
	if l.isValid() {
		l.impl.Error(message, args)
	}
}

func (l *localLogger) Notice(message string, args ...interface{}) {
	if l.isValid() {
		l.impl.Notice(message, args)
	}
}

func (l *localLogger) Info(message string, args ...interface{}) {
	if l.isValid() {
		l.impl.Info(message, args)
	}
}

func (l *localLogger) Warning(message string, args ...interface{}) {
	if l.isValid() {
		l.impl.Warning(message, args)
	}
}

func init() {
	logger = localLogger{
		impl: nil,
	}
}
