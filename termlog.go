// Package termlog provides facilities for logging to a terminal geared towards
// interactive use.
package termlog

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
)

const defaultTimeFmt = "15:04:05: "
const indent = "  "

// Palette defines the colour of output
type Palette struct {
	Timestamp *color.Color
	Say       *color.Color
	Notice    *color.Color
	Warn      *color.Color
	Shout     *color.Color
}

// DefaultPalette is a sensbile default palette, with the following foreground
// colours:
//
// 	Say: Terminal default
// 	Notice: Blue
// 	Warn: Yellow
// 	Shout: Red
// 	Timestamp: Cyan
var DefaultPalette = Palette{
	Say:       color.New(),
	Notice:    color.New(color.FgBlue),
	Warn:      color.New(color.FgYellow),
	Shout:     color.New(color.FgRed),
	Timestamp: color.New(color.FgCyan),
}

// Logger logs things
type Logger interface {
	Say(format string, args ...interface{})
	Notice(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Shout(format string, args ...interface{})

	SayAs(name string, format string, args ...interface{})
	NoticeAs(name string, format string, args ...interface{})
	WarnAs(name string, format string, args ...interface{})
	ShoutAs(name string, format string, args ...interface{})
}

// Group is a collected group of log entries
type Group interface {
	Logger
	Done()
	Quiet()
}

// TermLog is the top-level termlog interface
type TermLog interface {
	Logger
	Group() Group
	Quiet()
}

type line struct {
	name   string
	color  *color.Color
	format string
	args   []interface{}
}

// Log is the top-level log structure
type Log struct {
	mu      sync.Mutex
	Palette *Palette
	TimeFmt string
	enabled map[string]bool
	quiet   bool
}

// NewLog creates a new Log instance and initialises it with a set of defaults.
func NewLog() *Log {
	l := &Log{
		Palette: &DefaultPalette,
		enabled: make(map[string]bool),
		TimeFmt: defaultTimeFmt,
	}
	l.enabled[""] = true
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		l.Color(false)
	}
	return l
}

// Color sets the state of colour output - true to turn on, false to disable.
func (*Log) Color(state bool) {
	color.NoColor = !state
}

// Enable logging for a specified name
func (l *Log) Enable(name string) {
	l.enabled[name] = true
}

// Quiet disables all output
func (l *Log) Quiet() {
	l.quiet = true
}

func (l *Log) output(quiet bool, lines ...*line) {
	if quiet {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(lines) == 0 {
		return
	}
	first := true
	for _, line := range lines {
		if _, ok := l.enabled[line.name]; !ok {
			continue
		}
		var format string
		if first {
			l.Palette.Timestamp.Printf(
				"%s", time.Now().Format(l.TimeFmt),
			)
			first = false
			format = line.format + "\n"
		} else {
			format = indent + line.format + "\n"
		}
		line.color.Printf(format, line.args...)
	}
}

// Say logs a line
func (l *Log) Say(format string, args ...interface{}) {
	l.output(l.quiet, &line{"", l.Palette.Say, format, args})
}

// Notice logs a line with the Notice color
func (l *Log) Notice(format string, args ...interface{}) {
	l.output(l.quiet, &line{"", l.Palette.Notice, format, args})
}

// Warn logs a line with the Warn color
func (l *Log) Warn(format string, args ...interface{}) {
	l.output(l.quiet, &line{"", l.Palette.Warn, format, args})
}

// Shout logs a line with the Shout color
func (l *Log) Shout(format string, args ...interface{}) {
	l.output(l.quiet, &line{"", l.Palette.Shout, format, args})
}

// SayAs logs a line
func (l *Log) SayAs(name string, format string, args ...interface{}) {
	l.output(l.quiet, &line{name, l.Palette.Say, format, args})
}

// NoticeAs logs a line with the Notice color
func (l *Log) NoticeAs(name string, format string, args ...interface{}) {
	l.output(l.quiet, &line{name, l.Palette.Notice, format, args})
}

// WarnAs logs a line with the Warn color
func (l *Log) WarnAs(name string, format string, args ...interface{}) {
	l.output(l.quiet, &line{name, l.Palette.Warn, format, args})
}

// ShoutAs logs a line with the Shout color
func (l *Log) ShoutAs(name string, format string, args ...interface{}) {
	l.output(l.quiet, &line{name, l.Palette.Shout, format, args})
}

// Group creates a new log group
func (l *Log) Group() Group {
	return &group{
		palette: l.Palette,
		lines:   make([]*line, 0),
		log:     l,
		quiet:   l.quiet,
	}
}

// Group is a group of lines that constitue a single log entry that won't be
// split. Lines in a group are indented.
type group struct {
	palette *Palette
	lines   []*line
	log     *Log
	quiet   bool
}

func (g *group) addLine(name string, color *color.Color, format string, args []interface{}) {
	g.lines = append(g.lines, &line{name, color, format, args})
}

// Say logs a line
func (g *group) Say(format string, args ...interface{}) {
	g.addLine("", g.palette.Say, format, args)
}

// Notice logs a line with the Notice color
func (g *group) Notice(format string, args ...interface{}) {
	g.addLine("", g.palette.Notice, format, args)
}

// Warn logs a line with the Warn color
func (g *group) Warn(format string, args ...interface{}) {
	g.addLine("", g.palette.Warn, format, args)
}

// Shout logs a line with the Shout color
func (g *group) Shout(format string, args ...interface{}) {
	g.addLine("", g.palette.Shout, format, args)
}

// SayAs logs a line
func (g *group) SayAs(name string, format string, args ...interface{}) {
	g.addLine(name, g.palette.Say, format, args)
}

// NoticeAs logs a line with the Notice color
func (g *group) NoticeAs(name string, format string, args ...interface{}) {
	g.addLine(name, g.palette.Notice, format, args)
}

// WarnAs logs a line with the Warn color
func (g *group) WarnAs(name string, format string, args ...interface{}) {
	g.addLine(name, g.palette.Warn, format, args)
}

// ShoutAs logs a line with the Shout color
func (g *group) ShoutAs(name string, format string, args ...interface{}) {
	g.addLine(name, g.palette.Shout, format, args)
}

// Done outputs the group to screen
func (g *group) Done() {
	g.log.output(g.quiet, g.lines...)
}

// Quiet disables output for this subgroup
func (g *group) Quiet() {
	g.quiet = true
}

// NewContext creates a new context with an included Logger
func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, "termlog", logger)
}

// FromContext retrieves a Logger from a context. If no logger is present, we
// return a new silenced logger that will produce no output.
func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value("termlog").(Logger)
	if !ok {
		l := NewLog()
		l.Quiet()
		return l
	}
	return logger
}

// SetOutput sets the output writer for termlog (stdout by default).
func SetOutput(w io.Writer) {
	color.Output = w
}
