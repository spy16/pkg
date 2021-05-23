package errors

import (
	"net/http"
)

// ErrConversion can be used to represent conversion errors
func ErrConversion(err error) *Error {
	return &Error{
		Code:    1,
		Type:    "ConversionError",
		Details: map[string]interface{}{},
		Message: err.Error(),
	}
}

// ErrConnection is returned when connection fails
func ErrConnection(err error) *Error {
	return &Error{
		Code: 1,
		Type: "ConnectionError",
		Details: map[string]interface{}{
			"error": err.Error(),
		},
		Message: "Failed to connect to server",
	}
}

// ErrUnexpected represents an unexpected internal error
// Err field will be populated in this case
func ErrUnexpected(err error) *Error {
	st := &Error{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
		Type:    "Unknown",
		Details: map[string]interface{}{
			"error": err.Error(),
		},
	}
	return st
}
