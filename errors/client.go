package errors

import (
	"fmt"
	"net/http"
)

// ErrEmptyParam is returned when a required field has no value
func ErrEmptyParam(param string) *Error {
	return &Error{
		Code:    http.StatusBadRequest,
		Type:    "EmptyField",
		Message: fmt.Sprintf("%s cannot be empty", param),
		Details: map[string]interface{}{
			"param": param,
		},
	}
}

// ErrNoSuchUser represents 401
func ErrNoSuchUser(userid string) *Error {
	return &Error{
		Code:    http.StatusUnauthorized,
		Type:    "NoSuchUser",
		Details: map[string]interface{}{},
		Message: fmt.Sprintf("No user found with id '%s'", userid),
	}
}

// ErrMissingCredentials can be used when auth header is missing
func ErrMissingCredentials() *Error {
	return &Error{
		Code:    http.StatusUnauthorized,
		Type:    "MissingCredentials",
		Details: map[string]interface{}{},
		Message: fmt.Sprintf("Authentication credentials are missing"),
	}
}

// ErrUnauthenticated represents 401
func ErrUnauthenticated() *Error {
	return &Error{
		Code:    http.StatusUnauthorized,
		Type:    "Unauthenticated",
		Details: map[string]interface{}{},
		Message: fmt.Sprintf("Wrong credentials provided or user does not exist"),
	}
}

// ErrUnauthorized represents 403
func ErrUnauthorized(action, resource, user string) *Error {
	return &Error{
		Code: http.StatusForbidden,
		Type: "Unauthorized",
		Details: map[string]interface{}{
			"action":   action,
			"resource": resource,
			"user":     user,
		},
		Message: fmt.Sprintf("You do not have permission to %s resource %s", action, resource),
	}
}

// ErrMethodNotAllowed represents 405
func ErrMethodNotAllowed() *Error {
	return &Error{
		Code: http.StatusMethodNotAllowed,
		Type: "MethodNotAllowed",
	}
}

// ErrMissingParam represents a missing parameter
func ErrMissingParam(param string) *Error {
	st := &Error{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("A value is required for parameter '%s'", param),
		Type:    "MissingParameter",
		Details: map[string]interface{}{
			"param": param,
		},
	}
	return st
}

// ErrBadData represents a bad request with non-parsable data
func ErrBadData() *Error {
	reason := "Failed to parse data"
	st := &Error{
		Code:    http.StatusBadRequest,
		Message: reason,
		Type:    "BadRequest",
		Details: map[string]interface{}{},
	}
	return st
}

// ErrBadSpec is returned when certain specification (e.g. request
// body, data format in a file etc.) is invalid. `v` is an example
// specification format (e.g. a struct of the data model)
func ErrBadSpec(v interface{}) *Error {
	st := &Error{
		Code:    http.StatusBadRequest,
		Message: "Specification is invalid",
		Type:    "BadSpecification",
		Details: map[string]interface{}{
			"expected": v,
		},
	}
	return st
}

// ErrBadRequest represents a generic bad request
func ErrBadRequest(format string, args ...interface{}) *Error {
	st := &Error{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf(format, args...),
		Type:    "BadRequest",
		Details: map[string]interface{}{},
	}
	return st
}

// ErrNotFound represents a generic not found error object
func ErrNotFound(msg string) *Error {
	st := &Error{
		Code:    http.StatusNotFound,
		Message: msg,
		Type:    "NotFound",
		Details: map[string]interface{}{},
	}
	return st
}

// ErrResourceNotFound represents an access to non-existent resource.
// rid is resource-id (e.g. Bob) and rtype is resource type (e.g. User)
func ErrResourceNotFound(rid string, rtype string) *Error {
	st := &Error{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s not found with id '%s'", rtype, rid),
		Type:    "ResourceNotFound",
		Details: map[string]interface{}{
			"id":   rid,
			"type": rtype,
		},
	}
	return st
}

// ErrPathNotFound represents an access to non-existent path
func ErrPathNotFound(path string) *Error {
	st := &Error{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("Path '%s' not found", path),
		Type:    "NotFound",
		Details: map[string]interface{}{
			"path": path,
		},
	}
	return st
}

// ErrConflict represents a conflicting resource. rid and rtype
// are resource-id and resource-type respectively
func ErrConflict(rid string, rtype string) *Error {
	st := &Error{
		Code:    http.StatusConflict,
		Message: fmt.Sprintf("Resource of type '%s' already exists with id '%s'", rtype, rid),
		Type:    "Conflict",
		Details: map[string]interface{}{
			"id":   rid,
			"type": rtype,
		},
	}
	return st
}
