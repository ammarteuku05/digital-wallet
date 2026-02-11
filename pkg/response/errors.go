package response

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// ErrorResponse is the response that represents an error.
type ErrorResponse struct {
	HTTPCode  int         `json:"-"`
	Message   string      `json:"message"`
	Code      string      `json:"code,omitempty"`
	RequestID string      `json:"request_id"`
	Internal  error       `json:"-"`
	Data      interface{} `json:"data,omitempty"`
}

// Error is required by the error interface.
func (e ErrorResponse) Error() string {
	return e.Message
}

// StatusCode is required by CustomHTTPErrorHandler
func (e ErrorResponse) StatusCode() int {
	return e.HTTPCode
}

// IError represents a custom error with code and message
type IError struct {
	Code    string
	Message string
}

func (e IError) Error() string {
	return e.Message
}

// Predefined error types
var (
	ErrInternal                = IError{Code: "50000", Message: "Internal server error"}
	ErrUnauthorizedType        = IError{Code: "40100", Message: "Unauthorized access"}
	ErrForbiddenType           = IError{Code: "40300", Message: "Access forbidden"}
	ErrSessionExpiredType      = IError{Code: "44000", Message: "Session has expired"}
	ErrResourceNotFound        = IError{Code: "40400", Message: "Resource not found"}
	ErrBadRequestType          = IError{Code: "40000", Message: "Bad request"}
	ErrValidationFailed        = IError{Code: "40001", Message: "Validation failed"}
	ErrDuplicateEntry          = IError{Code: "40002", Message: "Duplicate entry"}
	ErrInvalidCredentials      = IError{Code: "40003", Message: "Invalid credentials"}
	ErrAccountDeactivated      = IError{Code: "40004", Message: "Account is deactivated"}
	ErrTokenExpired            = IError{Code: "40005", Message: "Token has expired"}
	ErrInvalidToken            = IError{Code: "40006", Message: "Invalid token"}
	ErrInsufficientPermissions = IError{Code: "40007", Message: "Insufficient permissions"}
	ErrMaxFileSizeExceed       = IError{Code: "40008", Message: "file size exceed, maxmimum is 10MB"}
	ErrFileExtNotAllowed       = IError{Code: "40009", Message: "file extension is not allowed"}
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// ErrInternalServerError creates a new error response representing an internal server error (HTTP 500)
func ErrInternalServerError(err error) ErrorResponse {
	if err == nil {
		err = ErrInternal
	}

	if _, ok := err.(stackTracer); !ok {
		err = errors.WithStack(err)
	}

	originalErr := errors.Cause(err)
	var Code string
	var errorMessage string

	if val, ok := originalErr.(IError); ok {
		Code = val.Code
		errorMessage = val.Message
	} else {
		Code = ErrInternal.Code
		errorMessage = ErrInternal.Message
	}

	return ErrorResponse{
		HTTPCode:  http.StatusInternalServerError,
		Message:   errorMessage,
		Code:      Code,
		Internal:  err,
		RequestID: uuid.New().String(),
	}
}

// ErrUnauthorized creates a new error response representing an unauthorized access (HTTP 401)
func ErrUnauthorized(err error) ErrorResponse {
	if err == nil {
		err = ErrUnauthorizedType
	}

	if _, ok := err.(stackTracer); !ok {
		err = errors.WithStack(err)
	}

	originalErr := errors.Cause(err)
	var Code string
	var errorMessage string

	if val, ok := originalErr.(IError); ok {
		Code = val.Code
		errorMessage = val.Message
	} else {
		Code = ErrUnauthorizedType.Code
		errorMessage = ErrUnauthorizedType.Message
	}

	return ErrorResponse{
		HTTPCode:  http.StatusUnauthorized,
		Message:   errorMessage,
		Code:      Code,
		Internal:  err,
		RequestID: uuid.New().String(),
	}
}

// ErrForbidden creates a new error response representing an authorization failure (HTTP 403)
func ErrForbidden(err error) ErrorResponse {
	if err == nil {
		err = ErrForbiddenType
	}

	if _, ok := err.(stackTracer); !ok {
		err = errors.WithStack(err)
	}

	originalErr := errors.Cause(err)
	var Code string
	var errorMessage string

	if val, ok := originalErr.(IError); ok {
		Code = val.Code
		errorMessage = val.Message
	} else {
		Code = ErrForbiddenType.Code
		errorMessage = ErrForbiddenType.Message
	}

	return ErrorResponse{
		HTTPCode:  http.StatusForbidden,
		Message:   errorMessage,
		Code:      Code,
		Internal:  err,
		RequestID: uuid.New().String(),
	}
}

// ErrSessionExpired creates a new error response representing an session expired error
func ErrSessionExpired(err error) ErrorResponse {
	if err == nil {
		err = ErrSessionExpiredType
	}

	if _, ok := err.(stackTracer); !ok {
		err = errors.WithStack(err)
	}

	originalErr := errors.Cause(err)
	var Code string
	var errorMessage string

	if val, ok := originalErr.(IError); ok {
		Code = val.Code
		errorMessage = val.Message
	} else {
		Code = ErrSessionExpiredType.Code
		errorMessage = ErrSessionExpiredType.Message
	}

	return ErrorResponse{
		HTTPCode:  440,
		Message:   errorMessage,
		Code:      Code,
		Internal:  err,
		RequestID: uuid.New().String(),
	}
}

// ErrNotFound creates a new error response representing a resource not found (HTTP 404)
func ErrNotFound(err error) ErrorResponse {
	if err == nil {
		err = ErrResourceNotFound
	}

	if _, ok := err.(stackTracer); !ok {
		err = errors.WithStack(err)
	}

	originalErr := errors.Cause(err)
	var Code string
	var errorMessage string

	if val, ok := originalErr.(IError); ok {
		Code = val.Code
		errorMessage = val.Message
	} else {
		Code = ErrResourceNotFound.Code
		errorMessage = ErrResourceNotFound.Message
	}

	return ErrorResponse{
		HTTPCode:  http.StatusNotFound,
		Message:   errorMessage,
		Code:      Code,
		Internal:  err,
		RequestID: uuid.New().String(),
	}
}

// ErrBadRequest creates a new error response representing a bad request (HTTP 400)
func ErrBadRequest(err error) ErrorResponse {
	if err == nil {
		err = ErrBadRequestType
	}

	if _, ok := err.(stackTracer); !ok {
		err = errors.WithStack(err)
	}

	originalErr := errors.Cause(err)
	var Code string
	var errorMessage string

	if val, ok := originalErr.(IError); ok {
		Code = val.Code
		errorMessage = val.Message
	} else {
		Code = ErrBadRequestType.Code
		errorMessage = ErrBadRequestType.Message
	}

	return ErrorResponse{
		HTTPCode:  http.StatusBadRequest,
		Message:   errorMessage,
		Code:      Code,
		Internal:  err,
		RequestID: uuid.New().String(),
	}
}

// ErrWithData creates an error response with additional data
func ErrWithData(err error, data interface{}, statusCode int) ErrorResponse {
	if err == nil {
		err = ErrBadRequestType
	}

	if _, ok := err.(stackTracer); !ok {
		err = errors.WithStack(err)
	}

	originalErr := errors.Cause(err)
	var Code string
	var errorMessage string

	if val, ok := originalErr.(IError); ok {
		Code = val.Code
		errorMessage = val.Message
	} else {
		Code = ErrBadRequestType.Code
		errorMessage = ErrBadRequestType.Message
	}

	if data == nil {
		data = map[string]interface{}{}
	}

	return ErrorResponse{
		HTTPCode:  statusCode,
		Message:   errorMessage,
		Code:      Code,
		Internal:  err,
		Data:      data,
		RequestID: uuid.New().String(),
	}
}

// HTTPError creates a custom error response with specified status code and error code
func HTTPError(err error, statusCode int, Code string, message string) ErrorResponse {
	if err == nil {
		err = ErrInternal
	}

	if _, ok := err.(stackTracer); !ok {
		err = errors.WithStack(err)
	}

	return ErrorResponse{
		HTTPCode:  statusCode,
		Message:   message,
		Code:      Code,
		Internal:  err,
		RequestID: uuid.New().String(),
	}
}

// GenerateResponseFromIError generates error response based on error code
func GenerateResponseFromIError(err error) ErrorResponse {
	if err == nil {
		return ErrInternalServerError(nil)
	}

	// Check if it's already an IError
	if iErr, ok := err.(IError); ok {
		switch iErr.Code {
		case ErrUnauthorizedType.Code:
			return ErrUnauthorized(err)
		case ErrForbiddenType.Code:
			return ErrForbidden(err)
		case ErrSessionExpiredType.Code:
			return ErrSessionExpired(err)
		case ErrResourceNotFound.Code:
			return ErrNotFound(err)
		case ErrInternal.Code:
			return ErrInternalServerError(err)
		default:
			return ErrBadRequest(err)
		}
	}

	// Try to parse error code from error message (legacy support)
	if len(err.Error()) >= 3 {
		switch err.Error()[0:3] {
		case "401":
			return ErrUnauthorized(err)
		case "403":
			return ErrForbidden(err)
		case "404":
			return ErrNotFound(err)
		case "500":
			return ErrInternalServerError(err)
		default:
			return ErrBadRequest(err)
		}
	}

	// Default to bad request
	return ErrBadRequest(err)
}

// New creates a new IError with custom message
func New(code, message string) IError {
	return IError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// Wrapf wraps an error with formatted additional context
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// Cause returns the underlying cause of the error
func Cause(err error) error {
	return errors.Cause(err)
}

// WithStack annotates err with a stack trace
func WithStack(err error) error {
	return errors.WithStack(err)
}

// WithMessage annotates err with a new message
func WithMessage(err error, message string) error {
	return errors.WithMessage(err, message)
}

// WithMessagef annotates err with the format specifier
func WithMessagef(err error, format string, args ...interface{}) error {
	return errors.WithMessagef(err, format, args...)
}

// Is reports whether any error in err's chain matches target
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Custom error creation functions for common scenarios

// NewValidationError creates a validation error
func NewValidationError(message string) IError {
	return IError{
		Code:    ErrValidationFailed.Code,
		Message: message,
	}
}

// NewDuplicateEntryError creates a duplicate entry error
func NewDuplicateEntryError(message string) IError {
	return IError{
		Code:    ErrDuplicateEntry.Code,
		Message: message,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) IError {
	return IError{
		Code:    ErrResourceNotFound.Code,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) IError {
	if message == "" {
		message = ErrUnauthorizedType.Message
	}
	return IError{
		Code:    ErrUnauthorizedType.Code,
		Message: message,
	}
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) IError {
	if message == "" {
		message = ErrForbiddenType.Message
	}
	return IError{
		Code:    ErrForbiddenType.Code,
		Message: message,
	}
}
