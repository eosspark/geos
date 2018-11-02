package log

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

const (
	timeFormat     = "2006-01-02T15:04:05-0700"
	termTimeFormat = "2006-01-02T15:04:05.000"
	floatFormat    = 'f'
	termMsgJust    = 40
)

type Format interface {
	Format(r *Record) []byte
}

func FormatFunc(f func(*Record) []byte) Format {
	return formatFunc(f)
}

type formatFunc func(*Record) []byte

func (f formatFunc) Format(r *Record) []byte {
	return f(r)
}

func TerminalFormat(useColor bool) Format {
	return FormatFunc(func(r *Record) []byte {
		var color = 0
		if useColor {
			switch r.Lvl {
			case LvlAll:
				color = 35
			case LvlError:
				color = 31
			case LvlWarn:
				color = 33
			case LvlInfo:
				color = 0
			case LvlDebug:
				color = 36
			case LvlOff:
				color = 34
			}
		}

		b := &bytes.Buffer{}
		if color > 0 {
			fmt.Fprintf(b, "\x1b[%dm[%s] %v %s \x1b[0m", color, r.Time.Format(termTimeFormat), r.Call, r.Name)
		} else {
			fmt.Fprintf(b, "[%s] %v %s ", r.Time.Format(termTimeFormat), r.Call, r.Name)
		}

		length := utf8.RuneCountInString(r.Call.String() + r.Name)
		if len(r.Msg) > 0 && length < termMsgJust {
			b.Write(bytes.Repeat([]byte{' '}, termMsgJust-length))
		}

		if color > 0 {
			fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m", color, r.Msg)
		} else {
			b.WriteString(r.Msg)
		}
		b.WriteByte('\n')

		return b.Bytes()

	})
}

func LogfmtFormat() Format {
	return FormatFunc(func(r *Record) []byte {
		b := &bytes.Buffer{}
		lvl := r.Lvl.AlignedString()
		fmt.Fprintf(b, "%s[%s] %v %s %s", lvl, r.Time.Format(termTimeFormat), r.Call, r.Name, r.Msg)
		b.WriteByte('\n')
		return b.Bytes()
	})
}
