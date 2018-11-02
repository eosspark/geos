package log

import (
	"fmt"
	"os"
)

var (
	root            = &logger{"", new(swapHandler)}
	StdoutHandler   = StreamHandler(os.Stdout, LogfmtFormat())
	TerminalHandler = StreamHandler(os.Stdout, TerminalFormat(true))
)

func init() {
	root.SetHandler(DiscardHandler())

	//root.SetHandler(TerminalHandler)

	//root.SetHandler(LvlFilterHandler(LvlError,TerminalHandler))

	//h,_ := FileHandler("./log.log",LogfmtFormat())
	//root.SetHandler(h)

}

func New(name ...string) Logger {
	return root.New(name...)
}

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
