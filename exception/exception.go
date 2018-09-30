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

func (t *logMessage) Message() string {
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

func (e *TimeoutException) Code() ExcTypes { return timeoutExceptionCode }
func (e *TimeoutException) What() string   { return "Timeout" }

type FileNotFoundException struct{ logMessage }

func (e *FileNotFoundException) Code() ExcTypes { return fileNotFoundExceptionCode }
func (e *FileNotFoundException) What() string   { return "File Not Found" }

/**
 * @brief report's parse errors
 */

type ParseErrorException struct{ logMessage }

func (e *ParseErrorException) Code() ExcTypes { return parseErrorExceptionCode }
func (e *ParseErrorException) What() string   { return "Parse Error" }

type InvalidArgException struct{ logMessage }

func (e *InvalidArgException) Code() ExcTypes { return invalidArgExceptionCode }
func (e *InvalidArgException) What() string   { return "Key Not Found" }

/**
 * @brief reports when a key, guid, or other item is not found.
 */

type KeyNotFoundException struct{ logMessage }

func (e *KeyNotFoundException) Code() ExcTypes { return keyNotFoundExceptionCode }
func (e *KeyNotFoundException) What() string   { return "Key Not Found" }

type BadCastException struct{ logMessage }

func (e *BadCastException) Code() ExcTypes { return badCastExceptionCode }
func (e *BadCastException) What() string   { return "Bad Cast" }

type OutOfRangeException struct{ logMessage }

func (e *OutOfRangeException) Code() ExcTypes { return outOfRangeExceptionCode }
func (e *OutOfRangeException) What() string   { return "Out of Range" }

/** @brief if an operation is unsupported or not valid this may be thrown */
type InvalidOperationException struct{ logMessage }

func (e *InvalidOperationException) Code() ExcTypes { return invalidOperationExceptionCode }
func (e *InvalidOperationException) What() string   { return "Invalid Operation" }

/** @brief if an host name can not be resolved this may be thrown */
type UnknownHostException struct{ logMessage }

func (e *UnknownHostException) Code() ExcTypes { return unknownHostExceptionCode }
func (e *UnknownHostException) What() string   { return "Unknown Host" }

/**
 *  @brief used to report a canceled Operation
 */
type CanceledException struct{ logMessage }

func (e *CanceledException) Code() ExcTypes { return canceledExceptionCode }
func (e *CanceledException) What() string   { return "Canceled" }

/**
 *  @brief used inplace of assert() to report violations of pre conditions.
 */