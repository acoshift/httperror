// Package httperror is the reusable http error collection
package httperror

import (
	"errors"
	"fmt"
	"net/http"
)

// Error is the httperror's Error
type Error struct {
	Status  int    `json:"status"`  // http status code
	Code    string `json:"code"`    // error code
	Message string `json:"message"` // error message
}

// Error implements error interface
func (err *Error) Error() string {
	return fmt.Sprintf("%s: [%d] %s", err.Code, err.Status, err.Message)
}

// Clone error
func (err Error) Clone() *Error {
	return &err
}

// Func is the error creator function
type Func func(error) error

// New is the helper function for create Func
func New(status int, code string) Func {
	return func(err error) error {
		return &Error{status, code, err.Error()}
	}
}

// StatusFunc is the error creator function pre-defined status
type StatusFunc func(string, error) error

// NewWithStatus is the helper function for create StatusFunc
func NewWithStatus(status int) StatusFunc {
	return func(code string, err error) error {
		return &Error{status, code, err.Error()}
	}
}

// CodeFunc is the error creator function pre-defined code
type CodeFunc func(int, error) error

// NewWithCode is the helper function for create CodeFunc
func NewWithCode(code string) CodeFunc {
	return func(status int, err error) error {
		return &Error{status, code, err.Error()}
	}
}

// NewHTTPError is the helper function for create http error
func NewHTTPError(status int, code string) error {
	return &Error{status, code, http.StatusText(status)}
}

// Pre-defined errors
var (
	BadRequest          = NewHTTPError(http.StatusBadRequest, "bad_request")
	Unauthorized        = NewHTTPError(http.StatusUnauthorized, "unauthorized")
	Forbidden           = NewHTTPError(http.StatusForbidden, "forbidden")
	NotFound            = NewHTTPError(http.StatusNotFound, "not_found")
	MethodNotAllowed    = NewHTTPError(http.StatusMethodNotAllowed, "method_not_allowed")
	RequestTimeout      = NewHTTPError(http.StatusRequestTimeout, "request_timeout")
	Conflict            = NewHTTPError(http.StatusConflict, "conflict")
	Gone                = NewHTTPError(http.StatusGone, "gone")
	InternalServerError = NewHTTPError(http.StatusInternalServerError, "internal_server_error")
	NotImplemented      = NewHTTPError(http.StatusNotImplemented, "not_implemented")
)

// Merge an error with other error
// if one or both errors are Error type, result will be an Error
// if none is Error, result will be native go's error
func Merge(err, other error) error {
	if other == nil {
		return err
	}
	if err == nil {
		return other
	}
	if e, ok := err.(*Error); ok {
		r := e.Clone()
		r.Message += "; " + other.Error()
		return r
	}
	if e, ok := other.(*Error); ok {
		r := e.Clone()
		r.Message += "; " + err.Error()
		return r
	}
	return errors.New(err.Error() + "; " + other.Error())
}

// BadRequestWith merges error with bad request
func BadRequestWith(err error) error {
	return Merge(BadRequest, err)
}

// InternalServerErrorWith merges error with internal server error
func InternalServerErrorWith(err error) error {
	return Merge(InternalServerError, err)
}

// ConflictWith merges error with conflict
func ConflictWith(err error) error {
	return Merge(Conflict, err)
}
