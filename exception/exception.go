package exception

import (
	. "github.com/eosspark/eos-go/log"
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
	FcLogMessage(lv LogLevel, format string, args ...interface{})
}

type FcException struct { LogMessage }
func (FcException) Code() ExcTypes { return UnspecifiedExceptionCode }
func (FcException) What() string   { return "unspecified" }

type UnHandledException struct{ LogMessage }

func (UnHandledException) Code() ExcTypes { return UnhandledExceptionCode }
func (UnHandledException) What() string   { return "" }

type TimeoutException struct{ LogMessage }

func (TimeoutException) Code() ExcTypes { return TimeoutExceptionCode }
func (TimeoutException) What() string   { return "Timeout" }

type FileNotFoundException struct{ LogMessage }

func (FileNotFoundException) Code() ExcTypes { return FileNotFoundExceptionCode }
func (FileNotFoundException) What() string   { return "File Not Found" }

/**
 * @brief report's parse errors
 */

type ParseErrorException struct{ LogMessage }

func (ParseErrorException) Code() ExcTypes { return ParseErrorExceptionCode }
func (ParseErrorException) What() string   { return "Parse Error" }

type InvalidArgException struct{ LogMessage }

func (InvalidArgException) Code() ExcTypes { return InvalidArgExceptionCode }
func (InvalidArgException) What() string   { return "Key Not Found" }

/**
 * @brief reports when a key, guid, or other item is not found.
 */

type KeyNotFoundException struct{ LogMessage }

func (KeyNotFoundException) Code() ExcTypes { return KeyNotFoundExceptionCode }
func (KeyNotFoundException) What() string   { return "Key Not Found" }

type BadCastException struct{ LogMessage }

func (BadCastException) Code() ExcTypes { return BadCastExceptionCode }
func (BadCastException) What() string   { return "Bad Cast" }

type OutOfRangeException struct{ LogMessage }

func (OutOfRangeException) Code() ExcTypes { return OutOfRangeExceptionCode }
func (OutOfRangeException) What() string   { return "Out of Range" }

/** @brief if an operation is unsupported or not valid this may be thrown */
type InvalidOperationException struct{ LogMessage }

func (InvalidOperationException) Code() ExcTypes { return InvalidOperationExceptionCode }
func (InvalidOperationException) What() string   { return "Invalid Operation" }

/** @brief if an host name can not be resolved this may be thrown */
type UnknownHostException struct{ LogMessage }

func (UnknownHostException) Code() ExcTypes { return UnknownHostExceptionCode }
func (UnknownHostException) What() string   { return "Unknown Host" }

/**
 *  @brief used to report a canceled Operation
 */
type CanceledException struct{ LogMessage }

func (CanceledException) Code() ExcTypes { return CanceledExceptionCode }
func (CanceledException) What() string   { return "Canceled" }

/**
 *  @brief used inplace of assert() to report violations of pre conditions.
 */
type AssertException struct{ LogMessage }

func (AssertException) Code() ExcTypes { return AssertExceptionCode }
func (AssertException) What() string   { return "Assert Exception" }

type EofException struct{ LogMessage }

func (EofException) Code() ExcTypes { return EofExceptionCode }
func (EofException) What() string   { return "End Of File" }

type NullOptional struct{ LogMessage }

func (NullOptional) Code() ExcTypes { return NullOptionalCode }
func (NullOptional) What() string   { return "null optional" }

type UdtException struct{ LogMessage }

func (UdtException) Code() ExcTypes { return UdtErrorCode }
func (UdtException) What() string   { return "UDT error" }

type AesException struct{ LogMessage }

func (AesException) Code() ExcTypes { return AesErrorCode }
func (AesException) What() string   { return "AES error" }

type OverflowException struct{ LogMessage }

func (OverflowException) Code() ExcTypes { return OverflowCode }
func (OverflowException) What() string   { return "Integer Overflow" }

type UnderflowException struct{ LogMessage }

func (UnderflowException) Code() ExcTypes { return UnderflowCode }
func (UnderflowException) What() string   { return "Integer Underflow" }

type DivideByZeroException struct{ LogMessage }

func (DivideByZeroException) Code() ExcTypes { return DivideByZeroCode }
func (DivideByZeroException) What() string   { return "Integer Divide By Zero" }
