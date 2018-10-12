package exception

import (
	"fmt"
)

type ExcTypes int

const (
	unspecifiedExceptionCode = ExcTypes(iota)
	unhandledExceptionCode    ///< for unhandled 3rd party exceptions
	timeoutExceptionCode      ///< timeout exceptions
	fileNotFoundExceptionCode
	parseErrorExceptionCode
	invalidArgExceptionCode
	keyNotFoundExceptionCode
	badCastExceptionCode
	outOfRangeExceptionCode
	canceledExceptionCode
	assertExceptionCode
	_
	eofExceptionCode
	stdExceptionCode
	invalidOperationExceptionCode
	unknownHostExceptionCode
	nullOptionalCode
	udtErrorCode
	aesErrorCode
	overflowCode
	underflowCode
	divideByZeroCode
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

type TimeoutException struct{ logMessage }

func (TimeoutException) Code() ExcTypes { return timeoutExceptionCode }
func (TimeoutException) What() string   { return "Timeout" }

type FileNotFoundException struct{ logMessage }

func (FileNotFoundException) Code() ExcTypes { return fileNotFoundExceptionCode }
func (FileNotFoundException) What() string   { return "File Not Found" }

/**
 * @brief report's parse errors
 */

type ParseErrorException struct{ logMessage }

func (ParseErrorException) Code() ExcTypes { return parseErrorExceptionCode }
func (ParseErrorException) What() string   { return "Parse Error" }

type InvalidArgException struct{ logMessage }

func (InvalidArgException) Code() ExcTypes { return invalidArgExceptionCode }
func (InvalidArgException) What() string   { return "Key Not Found" }

/**
 * @brief reports when a key, guid, or other item is not found.
 */

type KeyNotFoundException struct{ logMessage }

func (KeyNotFoundException) Code() ExcTypes { return keyNotFoundExceptionCode }
func (KeyNotFoundException) What() string   { return "Key Not Found" }

type BadCastException struct{ logMessage }

func (BadCastException) Code() ExcTypes { return badCastExceptionCode }
func (BadCastException) What() string   { return "Bad Cast" }

type OutOfRangeException struct{ logMessage }

func (OutOfRangeException) Code() ExcTypes { return outOfRangeExceptionCode }
func (OutOfRangeException) What() string   { return "Out of Range" }

/** @brief if an operation is unsupported or not valid this may be thrown */
type InvalidOperationException struct{ logMessage }

func (InvalidOperationException) Code() ExcTypes { return invalidOperationExceptionCode }
func (InvalidOperationException) What() string   { return "Invalid Operation" }

/** @brief if an host name can not be resolved this may be thrown */
type UnknownHostException struct{ logMessage }

func (UnknownHostException) Code() ExcTypes { return unknownHostExceptionCode }
func (UnknownHostException) What() string   { return "Unknown Host" }

/**
 *  @brief used to report a canceled Operation
 */
type CanceledException struct{ logMessage }

func (CanceledException) Code() ExcTypes { return canceledExceptionCode }
func (CanceledException) What() string   { return "Canceled" }

/**
 *  @brief used inplace of assert() to report violations of pre conditions.
 */
type AssertException struct{ logMessage }

func (AssertException) Code() ExcTypes { return assertExceptionCode }
func (AssertException) What() string   { return "Assert Exception" }

type EofException struct{ logMessage }

func (EofException) Code() ExcTypes { return eofExceptionCode }
func (EofException) What() string   { return "End Of File" }

type NullOptional struct{ logMessage }

func (NullOptional) Code() ExcTypes { return nullOptionalCode }
func (NullOptional) What() string   { return "null optional" }

type UdtException struct{ logMessage }

func (UdtException) Code() ExcTypes { return udtErrorCode }
func (UdtException) What() string   { return "UDT error" }

type AesException struct{ logMessage }

func (AesException) Code() ExcTypes { return aesErrorCode }
func (AesException) What() string   { return "AES error" }

type OverflowException struct{ logMessage }

func (OverflowException) Code() ExcTypes { return overflowCode }
func (OverflowException) What() string   { return "Integer Overflow" }

type UnderflowException struct{ logMessage }

func (UnderflowException) Code() ExcTypes { return underflowCode }
func (UnderflowException) What() string   { return "Integer Underflow" }

type DivideByZeroException struct{ logMessage }

func (DivideByZeroException) Code() ExcTypes { return divideByZeroCode }
func (DivideByZeroException) What() string   { return "Integer Divide By Zero" }
