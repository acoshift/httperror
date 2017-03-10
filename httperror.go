// Package httperror is the reusable http error collection
package httperror

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

// NewError creates new Error
func NewError(status int, code string, message string) error {
	return &Error{Status: status, Code: code, Message: message}
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

var mapHTTPStatusCode = map[int]string{
	http.StatusBadRequest:          "bad_request",
	http.StatusUnauthorized:        "unauthorized",
	http.StatusForbidden:           "forbidden",
	http.StatusNotFound:            "not_found",
	http.StatusMethodNotAllowed:    "method_not_allowed",
	http.StatusRequestTimeout:      "request_timeout",
	http.StatusConflict:            "conflict",
	http.StatusGone:                "gone",
	http.StatusInternalServerError: "internal_server_error",
	http.StatusNotImplemented:      "not_implemented",
}

func newPreDefinedHTTPError(status int) error {
	return NewHTTPError(status, mapHTTPStatusCode[status])
}

func newPreDefinedEmptyHTTPError(status int) error {
	return NewError(status, mapHTTPStatusCode[status], "")
}

// Pre-defined errors
var (
	BadRequest          = newPreDefinedHTTPError(http.StatusBadRequest)
	Unauthorized        = newPreDefinedHTTPError(http.StatusUnauthorized)
	Forbidden           = newPreDefinedHTTPError(http.StatusForbidden)
	NotFound            = newPreDefinedHTTPError(http.StatusNotFound)
	MethodNotAllowed    = newPreDefinedHTTPError(http.StatusMethodNotAllowed)
	RequestTimeout      = newPreDefinedHTTPError(http.StatusRequestTimeout)
	Conflict            = newPreDefinedHTTPError(http.StatusConflict)
	Gone                = newPreDefinedHTTPError(http.StatusGone)
	InternalServerError = newPreDefinedHTTPError(http.StatusInternalServerError)
	NotImplemented      = newPreDefinedHTTPError(http.StatusNotImplemented)

	// Empty message errors
	emptyBadRequest          = newPreDefinedEmptyHTTPError(http.StatusBadRequest)
	emptyUnauthorized        = newPreDefinedEmptyHTTPError(http.StatusUnauthorized)
	emptyForbidden           = newPreDefinedEmptyHTTPError(http.StatusForbidden)
	emptyNotFound            = newPreDefinedEmptyHTTPError(http.StatusNotFound)
	emptyMethodNotAllowed    = newPreDefinedEmptyHTTPError(http.StatusMethodNotAllowed)
	emptyRequestTimeout      = newPreDefinedEmptyHTTPError(http.StatusRequestTimeout)
	emptyConflict            = newPreDefinedEmptyHTTPError(http.StatusConflict)
	emptyGone                = newPreDefinedEmptyHTTPError(http.StatusGone)
	emptyInternalServerError = newPreDefinedEmptyHTTPError(http.StatusInternalServerError)
	emptyNotImplemented      = newPreDefinedEmptyHTTPError(http.StatusNotImplemented)
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
		if len(r.Message) > 0 {
			r.Message += "; "
		}
		r.Message += other.Error()
		return r
	}
	if e, ok := other.(*Error); ok {
		r := e.Clone()
		if len(r.Message) > 0 {
			r.Message += "; "
		}
		r.Message += err.Error()
		return r
	}
	return errors.New(err.Error() + "; " + other.Error())
}

// BadRequestWith merges error with bad request
func BadRequestWith(err error) error {
	return Merge(emptyBadRequest, err)
}

// UnauthorizedWith merges error with unauthorized
func UnauthorizedWith(err error) error {
	return Merge(emptyUnauthorized, err)
}

// ForbiddenWith merges error with forbidden
func ForbiddenWith(err error) error {
	return Merge(emptyForbidden, err)
}

// NotFoundWith merges error with not found
func NotFoundWith(err error) error {
	return Merge(emptyNotFound, err)
}

// MethodNotAllowedWith merges error with method not allowed
func MethodNotAllowedWith(err error) error {
	return Merge(emptyMethodNotAllowed, err)
}

// RequestTimeoutWith merges error with request timeout
func RequestTimeoutWith(err error) error {
	return Merge(emptyRequestTimeout, err)
}

// ConflictWith merges error with conflict
func ConflictWith(err error) error {
	return Merge(emptyConflict, err)
}

// GoneWith merges error with gone
func GoneWith(err error) error {
	return Merge(emptyGone, err)
}

// InternalServerErrorWith merges error with internal server error
func InternalServerErrorWith(err error) error {
	return Merge(emptyInternalServerError, err)
}

// GRPC maps grpc error to http error
func GRPC(err error) error {
	if err == nil {
		return nil
	}
	// check is err grpc error
	desc := grpc.ErrorDesc(err)
	switch grpc.Code(err) {
	case codes.OK:
		return nil
	case codes.Canceled:
		return NewError(http.StatusRequestTimeout, "canceled", desc)
	case codes.Unknown:
		return NewError(http.StatusInternalServerError, "unknown", desc)
	case codes.InvalidArgument:
		return NewError(http.StatusBadRequest, "invalid_argument", desc)
	case codes.DeadlineExceeded:
		return NewError(http.StatusRequestTimeout, "deadline_exceeded", desc)
	case codes.NotFound:
		return NewError(http.StatusNotFound, "not_found", desc)
	case codes.AlreadyExists:
		return NewError(http.StatusConflict, "already_exists", desc)
	case codes.PermissionDenied:
		return NewError(http.StatusForbidden, "permission_denied", desc)
	case codes.Unauthenticated:
		return NewError(http.StatusUnauthorized, "unauthenticated", desc)
	case codes.ResourceExhausted:
		return NewError(http.StatusForbidden, "resource_exhausted", desc)
	case codes.FailedPrecondition:
		return NewError(http.StatusPreconditionFailed, "failed_precondition", desc)
	case codes.Aborted:
		return NewError(http.StatusConflict, "aborted", desc)
	case codes.OutOfRange:
		return NewError(http.StatusBadRequest, "out_of_range", desc)
	case codes.Unimplemented:
		return NewError(http.StatusNotImplemented, "unimplemented", desc)
	case codes.Internal:
		return NewError(http.StatusInternalServerError, "internal", desc)
	case codes.Unavailable:
		return NewError(http.StatusServiceUnavailable, "service_unavailable", desc)
	case codes.DataLoss:
		return NewError(http.StatusInternalServerError, "data_loss", desc)
	default:
		return err
	}
}
