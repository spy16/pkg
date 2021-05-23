package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error implements the error interface and provides certain additional fields.
// Generally, instead of directly creating instance of Error object, common
// helper methods such as `ErrUnauthenticated` are defined in this package which
// can be used easily.
type Error struct {
	// Code is generally a 4xx or 5xx HTTP Standard code
	Code int `json:"code"`

	// Type is a single word, camel-cased name for the error. This helps in
	// differentiating between common error codes. For example, code can be
	// 400 to represent bad request and type name an be InvalidValue, InvalidType
	// MissingParameter etc. which are all bad requests
	Type string `json:"type"`

	// Details provides any additional information about the error. For example,
	// in case of MissingParameter error, details will contain the name of the
	// missing parameter.
	Details map[string]interface{} `json:"details,omitempty"`

	// Message is a string which can be directly shown to user. This should
	// not contain any technical errors.
	Message string `json:"message,omitempty"`
}

func (e Error) String() string {
	msg := fmt.Sprintf("%s (code: %d)", e.Message, e.Code)
	return msg
}

func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// Write appropriately formats the error object and writes it to the
// ResponseWriter. The `Code` field will also be sent as StatusCode in
// the response.
func (e Error) Write(w http.ResponseWriter) error {
	code := e.Code
	if code < 400 {
		code = http.StatusInternalServerError
	}
	return writeJSON(w, code, e)
}

// SetCode sets the value of code field
func (e *Error) SetCode(code int) *Error {
	e.Code = code
	return e
}

// SetMessage sets the value of message field
func (e *Error) SetMessage(msg string) *Error {
	e.Message = msg
	return e
}

// writeJSON serializes given interface using JSON encoder and writes it
// to given http.ResponseWriter object with StatusCode set to `code`.
func writeJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

// New returns a new instance of Error object
func New(format string, args ...interface{}) *Error {
	return ErrUnexpected(fmt.Errorf(format, args...))
}
