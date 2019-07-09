package usecase

import "fmt"

// Error defines the custom use case error interface
type Error interface {
	Code() ErrorCode
	Error() string
	Prefix(format string, a ...interface{}) Error
}

type errorData struct {
	code ErrorCode
	text string
}

// NewError formats and returns a use case error
func NewError(code ErrorCode, format string, a ...interface{}) Error {
	text := fmt.Sprintf(format, a...)
	return &errorData{code, text}
}

// Prefix adds the formatted string to the beginning of an existing error string
func (e *errorData) Prefix(format string, a ...interface{}) Error {
	e.text = fmt.Sprintf(format, a...) + ": " + e.text
	return e
}

// Code returns the error code
func (e *errorData) Code() ErrorCode {
	return e.code
}

// Error returns the error text
func (e *errorData) Error() string {
	return fmt.Sprintf("%v: %v", e.code, e.text)
}

// ErrorCode use case error type
type ErrorCode uint8

// Use case error types
const (
	ErrNone    ErrorCode = 0
	ErrUnknown ErrorCode = 1 << iota
	ErrRecordNotFound
	ErrDuplicateRecord
	ErrInvalidID
)

func (ec ErrorCode) String() string {
	switch ec {
	case ErrNone:
		return "None"
	case ErrUnknown:
		return "Unknown"
	case ErrRecordNotFound:
		return "Record not found"
	case ErrDuplicateRecord:
		return "Duplicate record"
	case ErrInvalidID:
		return "Invalid ID"
	}
	return "[Invalid error code]"
}
