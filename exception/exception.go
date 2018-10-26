package exception

import (
	"fmt"
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
	Message() string
	setMessage(s string)
}

// Exception log manager
type logMessage struct {
	message string
}

func (t logMessage) Message() string {
	return t.message
}

func (t *logMessage) setMessage(message string) {
	t.message = message
}

func EosAssert(expr bool, exception Exception, format string, args ...interface{}) {
	if !expr {
		throwException(exception, format, args...)
	}
}

func EosThrow(exception Exception, format string, args ...interface{}) {
	EosAssert(false, exception, format, args...)
}

func throwException(exception Exception, format string, args ...interface{}) {
	formatMessage(exception, format, args...)
	makeLog(exception)

	//throw := reflect.ValueOf(exception).Elem().Interface()
	panic(exception)
}

func formatMessage(exception Exception, format string, args ...interface{}) {
	exception.setMessage(fmt.Sprintf(format, args...))
}

func makeLog(exception Exception) {
	println(exception.Message())
}

type UnHandledException struct{ logMessage }

func (UnHandledException) Code() ExcTypes { return UnhandledExceptionCode }
func (UnHandledException) What() string   { return "" }

type TimeoutException struct{ logMessage }

func (TimeoutException) Code() ExcTypes { return TimeoutExceptionCode }
func (TimeoutException) What() string   { return "Timeout" }

type FileNotFoundException struct{ logMessage }

func (FileNotFoundException) Code() ExcTypes { return FileNotFoundExceptionCode }
func (FileNotFoundException) What() string   { return "File Not Found" }

/**
 * @brief report's parse errors
 */

type ParseErrorException struct{ logMessage }

func (ParseErrorException) Code() ExcTypes { return ParseErrorExceptionCode }
func (ParseErrorException) What() string   { return "Parse Error" }

type InvalidArgException struct{ logMessage }

func (InvalidArgException) Code() ExcTypes { return InvalidArgExceptionCode }
func (InvalidArgException) What() string   { return "Key Not Found" }

/**
 * @brief reports when a key, guid, or other item is not found.
 */

type KeyNotFoundException struct{ logMessage }

func (KeyNotFoundException) Code() ExcTypes { return KeyNotFoundExceptionCode }
func (KeyNotFoundException) What() string   { return "Key Not Found" }

type BadCastException struct{ logMessage }

func (BadCastException) Code() ExcTypes { return BadCastExceptionCode }
func (BadCastException) What() string   { return "Bad Cast" }

type OutOfRangeException struct{ logMessage }

func (OutOfRangeException) Code() ExcTypes { return OutOfRangeExceptionCode }
func (OutOfRangeException) What() string   { return "Out of Range" }

/** @brief if an operation is unsupported or not valid this may be thrown */
type InvalidOperationException struct{ logMessage }

func (InvalidOperationException) Code() ExcTypes { return InvalidOperationExceptionCode }
func (InvalidOperationException) What() string   { return "Invalid Operation" }

/** @brief if an host name can not be resolved this may be thrown */
type UnknownHostException struct{ logMessage }

func (UnknownHostException) Code() ExcTypes { return UnknownHostExceptionCode }
func (UnknownHostException) What() string   { return "Unknown Host" }

/**
 *  @brief used to report a canceled Operation
 */
type CanceledException struct{ logMessage }

func (CanceledException) Code() ExcTypes { return CanceledExceptionCode }
func (CanceledException) What() string   { return "Canceled" }

/**
 *  @brief used inplace of assert() to report violations of pre conditions.
 */
type AssertException struct{ logMessage }

func (AssertException) Code() ExcTypes { return AssertExceptionCode }
func (AssertException) What() string   { return "Assert Exception" }

type EofException struct{ logMessage }

func (EofException) Code() ExcTypes { return EofExceptionCode }
func (EofException) What() string   { return "End Of File" }

type NullOptional struct{ logMessage }

func (NullOptional) Code() ExcTypes { return NullOptionalCode }
func (NullOptional) What() string   { return "null optional" }

type UdtException struct{ logMessage }

func (UdtException) Code() ExcTypes { return UdtErrorCode }
func (UdtException) What() string   { return "UDT error" }

type AesException struct{ logMessage }

func (AesException) Code() ExcTypes { return AesErrorCode }
func (AesException) What() string   { return "AES error" }

type OverflowException struct{ logMessage }

func (OverflowException) Code() ExcTypes { return OverflowCode }
func (OverflowException) What() string   { return "Integer Overflow" }

type UnderflowException struct{ logMessage }

func (UnderflowException) Code() ExcTypes { return UnderflowCode }
func (UnderflowException) What() string   { return "Integer Underflow" }

type DivideByZeroException struct{ logMessage }

func (DivideByZeroException) Code() ExcTypes { return DivideByZeroCode }
func (DivideByZeroException) What() string   { return "Integer Divide By Zero" }
