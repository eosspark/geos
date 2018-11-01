package log

import (
	"os"

	"fmt"
)

var (
	root          = &logger{[]interface{}{}, new(swapHandler)}
	StdoutHandler = StreamHandler(os.Stdout, LogfmtFormat())
	StderrHandler = StreamHandler(os.Stderr, LogfmtFormat())
)

func init() {
	//root.SetHandler(DiscardHandler())
	//root.SetHandler(StdoutHandler)
	//root.SetHandler(LvlFilterHandler(LvlError, StdoutHandler))
	root.SetHandler(StreamHandler(os.Stdout, TerminalFormat(true)))

}

// New returns a new logger with the given context.
// New is a convenient alias for Root().New
func New(ctx ...interface{}) Logger {
	return root.New(ctx...)
}

// Root returns the root logger
func Root() Logger {
	return root
}

// The following functions bypass the exported logger methods (logger.Debug,
// etc.) to keep the call depth the same for all paths to logger.write so
// runtime.Caller(2) always refers to the call site in client code.

// Debug is a convenient alias for Root().Debug
func Debug(format string, v ...interface{}) {
	root.write(LvlDebug, fmt.Sprintf(format, v...), skipLevel)
}

// Info is a convenient alias for Root().Info
func Info(format string, v ...interface{}) {
	root.write(LvlInfo, fmt.Sprintf(format, v...), skipLevel)
}

// Warn is a convenient alias for Root().Warn
func Warn(format string, v ...interface{}) {
	root.write(LvlWarn, fmt.Sprintf(format, v...), skipLevel)
}

// Error is a convenient alias for Root().Error
func Error(format string, v ...interface{}) {
	root.write(LvlError, fmt.Sprintf(format, v...), skipLevel)
}

//// Crit is a convenient alias for Root().Crit
//func Crit(format string, v ...interface{}) {
//	root.write(LvlOff, fmt.Sprintf(format,v...), skipLevel)
//	os.Exit(1)
//}
