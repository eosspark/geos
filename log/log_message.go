package log

import (
	"fmt"
	"runtime/debug"
)

type Context struct {
	LogLevel  Lvl
	StackInfo []byte
}

func (c Context) String() string {
	return string(c.StackInfo)
}

type Message struct {
	context Context
	Format  string
	Args    []interface{}
}

func FcLogMessage(level Lvl, format string, args ...interface{}) Message {
	return Message{
		context: Context{
			LogLevel:  level,
			StackInfo: debug.Stack(),
		},
		Format: format,
		Args:   args,
	}
}

func (l Message) GetContext() Context {
	return l.context
}

func (l Message) GetMessage() string {
	return fmt.Sprintf(l.Format, l.Args...)
}

type Messages = []Message
