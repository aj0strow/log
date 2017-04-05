package log

import (
	"os"
	"sync"
)

type Std struct {
	mu     sync.Mutex
	buf    []byte
}

func (s *Std) Append(msg *Message) error {
	return s.Write(msg.Message)
}

func (s *Std) Write(msg string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buf = s.buf[:0]
	s.buf = append(s.buf, msg...)
	if len(msg) == 0 || msg[len(msg)-1] != '\n' {
		s.buf = append(s.buf, '\n')
	}
	if _, err := os.Stderr.Write(s.buf); err != nil {
		return err
	}
	return nil
}
