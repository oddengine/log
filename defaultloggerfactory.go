package log

import (
	"io"
	"strings"

	"github.com/oddcancer/log/level"
)

// DefaultLoggerFactory creates new DefaultLogger.
type DefaultLoggerFactory struct {
	out   io.Writer
	level level.Level
}

// Init this class.
func (me *DefaultLoggerFactory) Init(out io.Writer, n level.Level) *DefaultLoggerFactory {
	me.out = out
	me.level = n
	return me
}

// NewLogger returns a configured ILogger for the given scope.
func (me *DefaultLoggerFactory) NewLogger(scope string) ILogger {
	return NewDefaultLogger(me.out, me.level, strings.ToUpper(scope), DEFAULT_DEPTH)
}

// NewDefaultLoggerFactory creates a new DefaultLoggerFactory.
func NewDefaultLoggerFactory(constraints *DefaultWriterConstraints) *DefaultLoggerFactory {
	w := new(DefaultWriter).Init(constraints)
	n := level.Parse(constraints.Level, "|")
	return new(DefaultLoggerFactory).Init(w, n)
}
