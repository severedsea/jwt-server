package web

import (
	"net/http"

	"github.com/pkg/errors"
)

// Error represents a handler error. It contains web-related information such as
// HTTP status code, error code, description
// Implements standard error interface
type Error struct {
	Status   int                    `json:"-"`
	Code     string                 `json:"code"`
	Desc     string                 `json:"description"`
	Err      error                  `json:"-"`
	Metadata map[string]interface{} `json:"-"`
}

func (e Error) Error() string {
	return e.Desc
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) WithMetadata(metadata map[string]interface{}) *Error {
	e.Metadata = metadata

	return &e
}

// WithStack adds stack trace into the Error object
func WithStack(err error) *Error {
	var result *Error

	if !errors.As(err, &result) {
		result = &Error{Status: http.StatusInternalServerError, Code: "internal_error", Desc: err.Error()}
	}

	if result.Err == nil {
		result.Err = errors.WithStack(errors.New(result.Error()))
	}

	return result
}

// NewError returns a new Error object based on the provided err and message
func NewError(err error, message string) *Error {
	var result *Error

	switch errors.As(err, &result) {
	case true:
		result.Desc = message

	default:
		result = &Error{Status: http.StatusInternalServerError, Code: "internal_error", Desc: message}
	}
	result.Err = errors.WithStack(errors.New(message))

	return result
}
