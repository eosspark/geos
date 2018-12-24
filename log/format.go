package log

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

const (
	eosTimeFormat = "2006-01-02T15:04:05.999"
	floatFormat   = 'f'
	termMsgJust   = 70
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
			fmt.Fprintf(b, "\x1b[%dm%s %v %n\x1b[0m", color, r.Time.Format(eosTimeFormat), r.Call, r.Call)

		} else {
			fmt.Fprintf(b, "%s %v %n", r.Time.Format(eosTimeFormat), r.Call, r.Call)
		}

		length := utf8.RuneCountInString(r.Call.String()) + utf8.RuneCountInString(fmt.Sprintf("%n", r.Call))
		if len(r.Msg) > 0 && length < termMsgJust {
			b.Write(bytes.Repeat([]byte{' '}, termMsgJust-length))
		}
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m", color, "] ")
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
		fmt.Fprintf(b, "%s %s %v %n %s", lvl, r.Time.Format(eosTimeFormat), r.Call, r.Call, r.Msg)

		b.WriteByte('\n')
		return b.Bytes()
	})
}
