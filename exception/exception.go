package exception

import (
	. "github.com/eosspark/eos-go/log"
	"bytes"
	"strconv"
	"reflect"
)

type ExcTypes int

const (
	UnspecifiedExceptionCode = ExcTypes(iota)
	UnhandledExceptionCode    ///< for unhandled 3rd party exceptions
	TimeoutExceptionCode      ///< timeout exceptions
	FileNotFoundExceptionCode
	ParseErrorExceptionCode
	InvalidArgExceptionCode
	KeyNotFoundExceptionCode
	BadCastExceptionCode
	OutOfRangeExceptionCode
	CanceledExceptionCode
	AssertExceptionCode
	_
	EofExceptionCode
	StdExceptionCode
	InvalidOperationExceptionCode
	UnknownHostExceptionCode
	NullOptionalCode
	UdtErrorCode
	AesErrorCode
	OverflowCode
	UnderflowCode
	DivideByZeroCode
)

// base eos Exception interface, every Exception need to implements
type Exception interface {
	Code() ExcTypes
	What() string
	String() string
	AppendLog(l Message)
	GetLog() []Message

	message(e Exception) string
}

func GetDetailMessage(ex Exception) string {
	return ex.message(ex)
}

type ELog []Message

func NewELog(l Message) ELog {
	e := ELog{}
	e.AppendLog(l)
	return e
}

func (e *ELog) AppendLog(l Message) {
	*e = append(*e, l)
}

func (e ELog) GetLog() []Message {
	return e
}

func (e ELog) String() string {
	var buffer bytes.Buffer
	for i := range e {
		buffer.WriteString(e[i].GetMessage())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (e ELog) message(ex Exception) string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(ex.Code())))
	buffer.WriteString(" ")
	buffer.WriteString(reflect.TypeOf(ex).String())
	buffer.WriteString(": ")
	buffer.WriteString(ex.What())
	buffer.WriteString("\n")
	for i := range e {
		buffer.WriteString("[")
		buffer.WriteString(e[i].GetMessage())
		buffer.WriteString("]")
		buffer.WriteString("\n")
		buffer.WriteString(e[i].GetContext().String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

type FcException struct{ ELog }

func (FcException) Code() ExcTypes { return UnspecifiedExceptionCode }
func (FcException) What() string   { return "unspecified" }

type UnHandledException struct{ ELog }

func (UnHandledException) Code() ExcTypes { return UnhandledExceptionCode }
func (UnHandledException) What() string   { return "" }

type TimeoutException struct{ ELog }

func (TimeoutException) Code() ExcTypes { return TimeoutExceptionCode }
func (TimeoutException) What() string   { return "Timeout" }

type FileNotFoundException struct{ ELog }

func (FileNotFoundException) Code() ExcTypes { return FileNotFoundExceptionCode }
func (FileNotFoundException) What() string   { return "File Not Found" }

/**
 * @brief report's parse errors
 */

type ParseErrorException struct{ ELog }

func (ParseErrorException) Code() ExcTypes { return ParseErrorExceptionCode }
func (ParseErrorException) What() string   { return "Parse Error" }

type InvalidArgException struct{ ELog }

func (InvalidArgException) Code() ExcTypes { return InvalidArgExceptionCode }
func (InvalidArgException) What() string   { return "Key Not Found" }

/**
 * @brief reports when a key, guid, or other item is not found.
 */

type KeyNotFoundException struct{ ELog }

func (KeyNotFoundException) Code() ExcTypes { return KeyNotFoundExceptionCode }
func (KeyNotFoundException) What() string   { return "Key Not Found" }

type BadCastException struct{ ELog }

func (BadCastException) Code() ExcTypes { return BadCastExceptionCode }
func (BadCastException) What() string   { return "Bad Cast" }

type OutOfRangeException struct{ ELog }

func (OutOfRangeException) Code() ExcTypes { return OutOfRangeExceptionCode }
func (OutOfRangeException) What() string   { return "Out of Range" }

/** @brief if an operation is unsupported or not valid this may be thrown */
type InvalidOperationException struct{ ELog }

func (InvalidOperationException) Code() ExcTypes { return InvalidOperationExceptionCode }
func (InvalidOperationException) What() string   { return "Invalid Operation" }

/** @brief if an host name can not be resolved this may be thrown */
type UnknownHostException struct{ ELog }

func (UnknownHostException) Code() ExcTypes { return UnknownHostExceptionCode }
func (UnknownHostException) What() string   { return "Unknown Host" }

/**
 *  @brief used to report a canceled Operation
 */
type CanceledException struct{ ELog }

func (CanceledException) Code() ExcTypes { return CanceledExceptionCode }
func (CanceledException) What() string   { return "Canceled" }

/**
 *  @brief used inplace of assert() to report violations of pre conditions.
 */
type AssertException struct{ ELog }

func (AssertException) Code() ExcTypes { return AssertExceptionCode }
func (AssertException) What() string   { return "Assert Exception" }

type EofException struct{ ELog }

func (EofException) Code() ExcTypes { return EofExceptionCode }
func (EofException) What() string   { return "End Of File" }

type NullOptional struct{ ELog }

func (NullOptional) Code() ExcTypes { return NullOptionalCode }
func (NullOptional) What() string   { return "null optional" }

type UdtException struct{ ELog }

func (UdtException) Code() ExcTypes { return UdtErrorCode }
func (UdtException) What() string   { return "UDT error" }

type AesException struct{ ELog }

func (AesException) Code() ExcTypes { return AesErrorCode }
func (AesException) What() string   { return "AES error" }

type OverflowException struct{ ELog }

func (OverflowException) Code() ExcTypes { return OverflowCode }
func (OverflowException) What() string   { return "Integer Overflow" }

type UnderflowException struct{ ELog }

func (UnderflowException) Code() ExcTypes { return UnderflowCode }
func (UnderflowException) What() string   { return "Integer Underflow" }

type DivideByZeroException struct{ ELog }

func (DivideByZeroException) Code() ExcTypes { return DivideByZeroCode }
func (DivideByZeroException) What() string   { return "Integer Divide By Zero" }
