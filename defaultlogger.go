package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/oddcancer/log/level"
)

const (
	// DEFAULT_DEPTH = this frame + wrapper func + caller
	DEFAULT_DEPTH = 3
)

var (
	std = NewDefaultLogger(os.Stdout, level.TRACE, "CORE", DEFAULT_DEPTH+1)
)

// DefaultLogger encapsulates functionality for providing logging at user-defined levels.
type DefaultLogger struct {
	sync.RWMutex

	level     level.Level
	scope     string
	calldepth int
	trace     *log.Logger
	debug     *log.Logger
	info      *log.Logger
	warn      *log.Logger
	error     *log.Logger
}

// Init this class.
func (me *DefaultLogger) Init(n level.Level, scope string, calldepth int) *DefaultLogger {
	me.level = n
	me.scope = scope
	me.calldepth = calldepth
	return me
}

// WithTrace is a chainable configuration function which sets the trace-level logger.
func (me *DefaultLogger) WithTrace(logger *log.Logger) *DefaultLogger {
	me.level |= level.TRACE
	me.trace = logger
	return me
}

// WithDebug is a chainable configuration function which sets the debug-level logger.
func (me *DefaultLogger) WithDebug(logger *log.Logger) *DefaultLogger {
	me.debug = logger
	return me
}

// WithInfo is a chainable configuration function which sets the info-level logger.
func (me *DefaultLogger) WithInfo(logger *log.Logger) *DefaultLogger {
	me.info = logger
	return me
}

// WithWarn is a chainable configuration function which sets the warn-level logger.
func (me *DefaultLogger) WithWarn(logger *log.Logger) *DefaultLogger {
	me.warn = logger
	return me
}

// WithError is a chainable configuration function which sets the error-level logger.
func (me *DefaultLogger) WithError(logger *log.Logger) *DefaultLogger {
	me.error = logger
	return me
}

func (me *DefaultLogger) log(logger *log.Logger, n level.Level, s string) {
	me.Lock()
	defer me.Unlock()

	err := logger.Output(me.calldepth, s)
	if err != nil {
		Warnf("Failed to log: %s", err)
		return
	}

	if logger.Writer() != me.trace.Writer() && (me.level.Get()&level.TRACE) != 0 {
		me.trace.SetPrefix(getPrefix(n, me.scope))
		me.trace.Output(me.calldepth, s)
	}
}

func (me *DefaultLogger) logf(logger *log.Logger, n level.Level, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)

	me.Lock()
	defer me.Unlock()

	err := logger.Output(me.calldepth, s)
	if err != nil {
		Warnf("Failed to log: %s", err)
		return
	}

	if logger.Writer() != me.trace.Writer() && (me.level.Get()&level.TRACE) != 0 {
		me.trace.SetPrefix(getPrefix(n, me.scope))
		me.trace.Output(me.calldepth, s)
	}
}

// Trace emits the preformatted message if the logger is at or below trace-level.
func (me *DefaultLogger) Trace(s string) {
	if (me.level.Get() & level.TRACE) != 0 {
		me.log(me.trace, level.TRACE, s)
	}
}

// Tracef formats and emits a message if the logger is at or below trace-level.
func (me *DefaultLogger) Tracef(format string, args ...interface{}) {
	if (me.level.Get() & level.TRACE) != 0 {
		me.logf(me.trace, level.TRACE, format, args...)
	}
}

// Debug emits the preformatted message if the logger is at or below debug-level.
func (me *DefaultLogger) Debug(n uint32, s string) {
	if (me.level.Get() & level.DEBUG) <= level.DEBUG0<<n {
		me.log(me.debug, level.DEBUG0<<n, s)
	}
}

// Debugf formats and emits a message if the logger is at or below debug-level.
func (me *DefaultLogger) Debugf(n uint32, format string, args ...interface{}) {
	if (me.level.Get() & level.DEBUG) <= level.DEBUG0<<n {
		me.logf(me.debug, level.DEBUG0<<n, format, args...)
	}
}

// Info emits the preformatted message if the logger is at or below info-level.
func (me *DefaultLogger) Info(s string) {
	if (me.level.Get() & level.INFO) != 0 {
		me.log(me.info, level.INFO, s)
	}
}

// Infof formats and emits a message if the logger is at or below info-level.
func (me *DefaultLogger) Infof(format string, args ...interface{}) {
	if (me.level.Get() & level.INFO) != 0 {
		me.logf(me.info, level.INFO, format, args...)
	}
}

// Warn emits the preformatted message if the logger is at or below warn-level.
func (me *DefaultLogger) Warn(s string) {
	if (me.level.Get() & level.WARN) != 0 {
		me.log(me.warn, level.WARN, s)
	}
}

// Warnf formats and emits a message if the logger is at or below warn-level.
func (me *DefaultLogger) Warnf(format string, args ...interface{}) {
	if (me.level.Get() & level.WARN) != 0 {
		me.logf(me.warn, level.WARN, format, args...)
	}
}

// Error emits the preformatted message if the logger is at or below error-level.
func (me *DefaultLogger) Error(s string) {
	if (me.level.Get() & level.ERROR) != 0 {
		me.log(me.error, level.ERROR, s)
	}
}

// Errorf formats and emits a message if the logger is at or below error-level.
func (me *DefaultLogger) Errorf(format string, args ...interface{}) {
	if (me.level.Get() & level.ERROR) != 0 {
		me.logf(me.error, level.ERROR, format, args...)
	}
}

// NewDefaultLogger returns a configured ILogger.
func NewDefaultLogger(out io.Writer, n level.Level, scope string, calldepth int) *DefaultLogger {
	return new(DefaultLogger).Init(n, scope, calldepth).
		WithTrace(log.New(os.Stdout, getPrefix(level.TRACE, scope), log.LstdFlags|log.Lshortfile)).
		WithDebug(log.New(out, getPrefix(level.DEBUG, scope), log.LstdFlags|log.Lshortfile)).
		WithInfo(log.New(out, getPrefix(level.INFO, scope), log.LstdFlags|log.Lshortfile)).
		WithWarn(log.New(out, getPrefix(level.WARN, scope), log.LstdFlags|log.Lshortfile)).
		WithError(log.New(out, getPrefix(level.ERROR, scope), log.LstdFlags|log.Lshortfile))
}

func getPrefix(n level.Level, scope string) string {
	if n >= level.ERROR {
		return fmt.Sprintf("[ERROR] %s ", scope)
	}
	if n >= level.WARN {
		return fmt.Sprintf("[WARN ] %s ", scope)
	}
	if n >= level.INFO {
		return fmt.Sprintf("[INFO ] %s ", scope)
	}
	if (n & level.DEBUG) != 0 {
		return fmt.Sprintf("[DEBUG] %s ", scope)
	}
	if (n & level.TRACE) != 0 {
		return fmt.Sprintf("[TRACE] %s ", scope)
	}
	return fmt.Sprintf("%s ", scope)
}
