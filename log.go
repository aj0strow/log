package log

import (
	"fmt"
	"os"
	"time"
)

// Level represents a scale of importance. You can filter log messages
// by importance to limit the volume of messages.
type Level int

func (level Level) String() string {
	switch level {
	case Trace:
		return "TRACE"
	case Info:
		return "INFO"
	case Error:
		return "ERROR"
	default:
		return "NONE"
	}
}

const (
	Trace Level = 0
	Info  Level = 1
	Error Level = 2
)

// Message is one discrete log event.
type Message struct {
	Level   Level
	Time    time.Time
	Message string
}

// Sink is an interface for any log output.
type Sink interface {
	Append(*Message) error
}

func New() *Logger {
	return &Logger{
		TimeSource: &LocalTime{},
	}
}

// Logger is a sink itself that appends to a filtered collection
// of output sinks. It includes helper methods to make creating
// and appending messages easy.
type Logger struct {
	Sinks      []Sink
	TimeSource TimeSource
}

func (logger *Logger) AddSink(level Level, sink Sink) {
	logger.Sinks = append(logger.Sinks, &Filter{
		Sink:  sink,
		Level: level,
	})
}

func (logger *Logger) Append(msg *Message) error {
	for _, sink := range logger.Sinks {
		if err := sink.Append(msg); err != nil {
			if msg.Level < Error {
				logger.Error(err)
			}
		}
	}
	return nil
}

func (logger *Logger) Appendf(level Level, msg string, args ...interface{}) error {
	return logger.Append(&Message{
		Level:   level,
		Time:    logger.TimeSource.Now(),
		Message: fmt.Sprintf(msg, args...),
	})
}

func (logger *Logger) Trace(msg string, args ...interface{}) error {
	return logger.Appendf(Trace, msg, args...)
}

func (logger *Logger) Info(msg string, args ...interface{}) error {
	return logger.Appendf(Info, msg, args...)
}

func (logger *Logger) Errorf(msg string, args ...interface{}) error {
	return logger.Appendf(Error, msg, args...)
}

func (logger *Logger) Error(err error) error {
	return logger.Errorf(err.Error())
}

func (logger *Logger) Fatalf(msg string, args ...interface{}) {
	logger.Errorf(msg, args...)
	os.Exit(1)
}

func (logger *Logger) Fatal(err error) {
	logger.Error(err)
	os.Exit(1)
}

var _ Sink = (*Logger)(nil)

// Filter is a sink that wraps another sink and only
// appends messages of equal or higher level.
type Filter struct {
	Sink  Sink
	Level Level
}

func (filter *Filter) Append(msg *Message) error {
	if msg.Level >= filter.Level {
		return filter.Sink.Append(msg)
	}
	return nil
}

var _ Sink = (*Filter)(nil)

// Stub time package to allow testing. If you don't like
// local time, feel free to implement TimeSource.
type TimeSource interface {
	Now() time.Time
}

// Default time source.
type LocalTime struct{}

func (*LocalTime) Now() time.Time {
	return time.Now()
}

var _ TimeSource = (*LocalTime)(nil)
