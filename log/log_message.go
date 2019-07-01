package log

import (
	"fmt"
	"runtime/debug"

	"github.com/go-stack/stack"
)

type Context struct {
	LogLevel  Lvl
	Call      stack.Call
	StackInfo []byte
}

func (c Context) String() string {
	if len(c.StackInfo) > 0 {
		return string(c.StackInfo)
	}
	return c.Call.String()
}

type Message struct {
	context Context
	Format  string
	Args    []interface{}
}

func FcLogMessage(level Lvl, format string, args ...interface{}) Message {
	return LogMessage(level, format, args, 2)
}

func LogMessage(level Lvl, format string, args []interface{}, skip ...int) Message {
	msg := Message{
		context: Context{
			LogLevel: level,
		},
		Format: format,
		Args:   args,
	}

	if len(skip) > 0 {
		msg.context.Call = stack.Caller(skip[0])
	} else {
		msg.context.StackInfo = debug.Stack()
	}

	return msg
}

func (l Message) GetContext() Context {
	return l.context
}

func (l Message) GetMessage() string {
	return fmt.Sprintf(l.Format, l.Args...)
}

type Messages = []Message
