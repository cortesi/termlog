package termlog

import (
	"math/rand"
	"time"
)

type stream struct {
	header string
	quiet  bool
	id     string
	log    *Log
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func rndstr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Say logs a line
func (s *stream) Say(format string, args ...interface{}) {
	s.log.output(s.quiet, &line{"", s.log.format(true, say, format, args), s})
}

// Notice logs a line with the Notice color
func (s *stream) Notice(format string, args ...interface{}) {
	s.log.output(s.quiet, &line{"", s.log.format(true, notice, format, args), s})
}

// Warn logs a line with the Warn color
func (s *stream) Warn(format string, args ...interface{}) {
	s.log.output(s.quiet, &line{"", s.log.format(true, warn, format, args), s})
}

// Shout logs a line with the Shout color
func (s *stream) Shout(format string, args ...interface{}) {
	s.log.output(s.quiet, &line{"", s.log.format(true, shout, format, args), s})
}

// SayAs logs a line
func (s *stream) SayAs(name string, format string, args ...interface{}) {
	s.log.output(s.quiet, &line{name, s.log.format(true, say, format, args), s})
}

// NoticeAs logs a line with the Notice color
func (s *stream) NoticeAs(name string, format string, args ...interface{}) {
	s.log.output(s.quiet, &line{name, s.log.format(true, notice, format, args), s})
}

// WarnAs logs a line with the Warn color
func (s *stream) WarnAs(name string, format string, args ...interface{}) {
	s.log.output(s.quiet, &line{name, s.log.format(true, warn, format, args), s})
}

// ShoutAs logs a line with the Shout color
func (s *stream) ShoutAs(name string, format string, args ...interface{}) {
	s.log.output(s.quiet, &line{name, s.log.format(true, shout, format, args), s})
}

// Quiet disables output for this subgroup
func (s *stream) Quiet() {
	s.quiet = true
}

// Header immedately outputs the stream header
func (s *stream) Header() {
	s.log.mu.Lock()
	defer s.log.mu.Unlock()
	s.log.header(s)
}

func (s *stream) getID() string {
	if s.id == "" {
		s.id = rndstr(16)
	}
	return s.id
}

func (s *stream) getHeader() string {
	return s.header
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
