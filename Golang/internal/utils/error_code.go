package utils

import "fmt"

/*
Used to define error codes and to debug what the problem is
*/

// Predefined error codes
const (
	ErrFileAlreadyExists  = 1001
	ErrFileCreation       = 1002
	ErrSchemaRegistration = 1003
)

// Error represents a custom error type with a code and message
type Error struct {
	Code    int
	Message string
}

// NewError creates a new custom error
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Error implements the built-in error interface
func (e *Error) Error() string {
	return fmt.Sprintf("[Error Code %d]: %s", e.Code, e.Message)
}
