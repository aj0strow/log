package log

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type MockSink struct {
	Error    error
	Messages []*Message
}

func (mock *MockSink) Len() int {
	return len(mock.Messages)
}

func (mock *MockSink) Append(msg *Message) error {
	mock.Messages = append(mock.Messages, msg)
	return mock.Error
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		Level  Level
		String string
	}{
		{Error, "ERROR"},
		{Info, "INFO"},
		{Trace, "TRACE"},
	}
	for _, tt := range tests {
		if tt.Level.String() != tt.String {
			t.Errorf("wrong level string: %v != %s", tt.Level, tt.String)
		}
	}
}

func TestFilter(t *testing.T) {
	s1 := &MockSink{}
	f1 := &Filter{
		Sink:  s1,
		Level: Error,
	}
	f1.Append(&Message{
		Level: Trace,
	})
	if s1.Len() > 0 {
		t.Errorf("filter should have 0 writes")
	}
	f1.Append(&Message{
		Level: Info,
	})
	if s1.Len() > 0 {
		t.Errorf("filter should have 0 writes")
	}
	f1.Append(&Message{
		Level: Error,
	})
	if s1.Len() != 1 {
		t.Errorf("filter should have 1 write")
	}

	s2 := &MockSink{}
	f2 := &Filter{
		Sink:  s2,
		Level: Info,
	}
	f2.Append(&Message{
		Level: Trace,
	})
	f2.Append(&Message{
		Level: Info,
	})
	f2.Append(&Message{
		Level: Error,
	})
	if s2.Len() != 2 {
		t.Errorf("filter should have 2 writes")
	}
}

func TestLoggerSinks(t *testing.T) {
	s1 := &MockSink{}
	s2 := &MockSink{}
	logger := New()
	logger.AddSink(Error, s1)
	logger.AddSink(Info, s2)
	logger.Append(&Message{
		Level: Info,
	})
	if s1.Len() != 0 {
		t.Errorf("filter should have 0 writes")
	}
	if s2.Len() != 1 {
		t.Errorf("filter should have 1 write")
	}
}

type FixedTime struct {
	Time time.Time
}

func (fixed *FixedTime) Now() time.Time {
	return fixed.Time
}

func TestLoggerMethods(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		Log     func(*Logger)
		Message *Message
	}{
		{
			func(logger *Logger) {
				logger.Error(fmt.Errorf("error message"))
			},
			&Message{
				Level:   Error,
				Time:    now,
				Message: "error message",
			},
		},
		{
			func(logger *Logger) {
				logger.Errorf("error %s", "message")
			},
			&Message{
				Level:   Error,
				Time:    now,
				Message: "error message",
			},
		},
		{
			func(logger *Logger) {
				logger.Info("info %s", "message")
			},
			&Message{
				Level:   Info,
				Time:    now,
				Message: "info message",
			},
		},
		{
			func(logger *Logger) {
				logger.Trace("trace %s", "message")
			},
			&Message{
				Level:   Trace,
				Time:    now,
				Message: "trace message",
			},
		},
	}
	for _, test := range tests {
		logger := New()
		logger.TimeSource = &FixedTime{now}
		s := &MockSink{}
		logger.AddSink(Trace, s)
		test.Log(logger)
		if !reflect.DeepEqual(s.Messages[0], test.Message) {
			t.Errorf("incorrect message:\nwant:%# v\nhave:%# v\n", test.Message, s.Messages[0])
		}
	}
}

func TestLoggerPropagate(t *testing.T) {
	logger := New()
	s1 := &MockSink{}
	logger.AddSink(Trace, s1)
	s2 := &MockSink{}
	logger.AddSink(Error, s2)
	s1.Error = fmt.Errorf("no more space")
	logger.Info("nothing serious")
	if s2.Len() != 1 {
		t.Errorf("log errors should propagate")
	}
	if s2.Messages[0].Message != "no more space" {
		t.Errorf("wrong message")
	}
}
