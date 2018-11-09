package log

import "fmt"

type LogLevel = Lvl

type LogMessage struct {
	Level  LogLevel
	Format string
	Args   []interface{}
}

func (l LogMessage) FcLogMessage(level LogLevel, format string, args ...interface{}) {
	l.Level = level
	l.Format = format
	l.Args = args
}

func (l LogMessage) Message() string {
	return fmt.Sprintf(l.Format, l.Args...)
}
