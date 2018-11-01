package log

import (
	"fmt"
	"time"

	"github.com/go-stack/stack"
)

const timeKey = "t"
const lvlKey = "lvl"
const msgKey = "msg"
const ctxKey = "ctx"
const errorKey = "LOG15_ERROR"
const skipLevel = 2

type Lvl int

const (
	LvlAll Lvl = iota
	LvlDebug
	LvlInfo
	LvlWarn
	LvlError
	LvlOff
)

// AlignedString returns a 5-character string containing the name of a Lvl.
func (l Lvl) AlignedString() string {
	switch l {
	case LvlAll:
		return "ALL"
	case LvlDebug:
		return "DEBUG"
	case LvlInfo:
		return "INFO "
	case LvlWarn:
		return "WARN "
	case LvlError:
		return "ERROR"
	case LvlOff:
		return "OFF "
	default:
		panic("Unknown level")
	}
}

// Strings returns the name of a Lvl.
func (l Lvl) String() string {
	switch l {
	case LvlAll:
		return "all"
	case LvlDebug:
		return "dbug"
	case LvlInfo:
		return "info"
	case LvlWarn:
		return "warn"
	case LvlError:
		return "eror"
	case LvlOff:
		return "off"
	default:
		panic("Unknown level")
	}
}

// LvlFromString returns the appropriate Lvl from a string name.
// Useful for parsing command line args and configuration files.
func LvlFromString(lvlString string) (Lvl, error) {
	switch lvlString {
	case "all":
		return LvlAll, nil
	case "debug", "dbug":
		return LvlDebug, nil
	case "info":
		return LvlInfo, nil
	case "warn":
		return LvlWarn, nil
	case "error", "eror":
		return LvlError, nil
	case "off":
		return LvlOff, nil
	default:
		return LvlDebug, fmt.Errorf("Unknown level: %v", lvlString)
	}
}

// A Record is what a Logger asks its handler to write
type Record struct {
	Time     time.Time
	Lvl      Lvl
	Msg      string
	Ctx      []interface{}
	Call     stack.Call
	KeyNames RecordKeyNames
}

// RecordKeyNames gets stored in a Record when the write function is executed.
type RecordKeyNames struct {
	Time string
	Msg  string
	Lvl  string
	Ctx  string
}

// A Logger writes key/value pairs to a Handler
type Logger interface {
	// New returns a new Logger that has this logger's context plus the given context
	New(ctx ...interface{}) Logger

	// GetHandler gets the handler associated with the logger.
	GetHandler() Handler

	// SetHandler updates the logger to write records to the specified handler.
	SetHandler(h Handler)

	// Log a message at the given level with context key/value pairs
	Debug(format string, ctx ...interface{})
	Info(format string, ctx ...interface{})
	Warn(format string, ctx ...interface{})
	Error(format string, ctx ...interface{})
	//Crit(format string, ctx ...interface{})

}

type logger struct {
	ctx []interface{}
	h   *swapHandler
}

func (l *logger) write(lvl Lvl, msg string, skip int) {
	l.h.Log(&Record{
		Time: time.Now(),
		Lvl:  lvl,
		Msg:  msg,
		Call: stack.Caller(skip),
		KeyNames: RecordKeyNames{
			Time: timeKey,
			Msg:  msgKey,
			Lvl:  lvlKey,
			Ctx:  ctxKey,
		},
	})
}

func (l *logger) New(ctx ...interface{}) Logger {
	child := &logger{newContext(l.ctx, ctx), new(swapHandler)}
	child.SetHandler(l.h)
	return child
}

func newContext(prefix []interface{}, suffix []interface{}) []interface{} {
	normalizedSuffix := normalize(suffix)
	newCtx := make([]interface{}, len(prefix)+len(normalizedSuffix))
	n := copy(newCtx, prefix)
	copy(newCtx[n:], normalizedSuffix)
	return newCtx
}

func (l *logger) Debug(format string, v ...interface{}) {
	l.write(LvlDebug, fmt.Sprintf(format, v...), skipLevel)
}

func (l *logger) Info(format string, v ...interface{}) {
	l.write(LvlInfo, fmt.Sprintf(format, v...), skipLevel)
}

func (l *logger) Warn(format string, v ...interface{}) {
	l.write(LvlWarn, fmt.Sprintf(format, v...), skipLevel)
}

func (l *logger) Error(format string, v ...interface{}) {
	root.write(LvlError, fmt.Sprintf(format, v...), skipLevel)
}

//func (l *logger) Crit(format string, v ...interface{}) {
//	root.write(LvlOff,fmt.Sprintf(format,v...),skipLevel)
//	os.Exit(1)
//}
func (l *logger) GetHandler() Handler {
	return l.h.Get()
}

func (l *logger) SetHandler(h Handler) {
	l.h.Swap(h)
}

func normalize(ctx []interface{}) []interface{} {
	// if the caller passed a Ctx object, then expand it
	if len(ctx) == 1 {
		if ctxMap, ok := ctx[0].(Ctx); ok {
			ctx = ctxMap.toArray()
		}
	}

	// ctx needs to be even because it's a series of key/value pairs
	// no one wants to check for errors on logging functions,
	// so instead of erroring on bad input, we'll just make sure
	// that things are the right length and users can fix bugs
	// when they see the output looks wrong
	if len(ctx)%2 != 0 {
		//ctx = append(ctx, nil, errorKey, "Normalized odd number of arguments by adding nil")
		ctx = append(ctx, "")
	}

	return ctx
}

// Lazy allows you to defer calculation of a logged value that is expensive
// to compute until it is certain that it must be evaluated with the given filters.
//
// Lazy may also be used in conjunction with a Logger's New() function
// to generate a child logger which always reports the current value of changing
// state.
//
// You may wrap any function which takes no arguments to Lazy. It may return any
// number of values of any type.
type Lazy struct {
	Fn interface{}
}

// Ctx is a map of key/value pairs to pass as context to a log function
// Use this only if you really need greater safety around the arguments you pass
// to the logging functions.
type Ctx map[string]interface{}

func (c Ctx) toArray() []interface{} {
	arr := make([]interface{}, len(c)*2)

	i := 0
	for k, v := range c {
		arr[i] = k
		arr[i+1] = v
		i += 2
	}

	return arr
}
