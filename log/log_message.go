package log

import "fmt"

type LogLevel = Lvl

type LogMessage struct {
	Level  LogLevel
	Format string
	Args   []interface{}
}

func FcLogMessage(level LogLevel, format string, args ...interface{}) LogMessage {
	return LogMessage{
		Level:  level,
		Format: format,
		Args:   args,
	}
}

func (l LogMessage) Message() string {
	return fmt.Sprintf(l.Format, l.Args...)
}
