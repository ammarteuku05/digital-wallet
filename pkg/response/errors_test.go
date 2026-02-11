package response

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIError_Error(t *testing.T) {
	err := IError{Code: "40000", Message: "test error"}
	assert.Equal(t, "test error", err.Error())
}

func TestErrorResponse_Error(t *testing.T) {
	errResp := ErrorResponse{Message: "error message"}
	assert.Equal(t, "error message", errResp.Error())
}

func TestErrorResponse_StatusCode(t *testing.T) {
	errResp := ErrorResponse{HTTPCode: http.StatusBadRequest}
	assert.Equal(t, http.StatusBadRequest, errResp.StatusCode())
}

func TestErrInternalServerError(t *testing.T) {
	t.Run("with nil error", func(t *testing.T) {
		result := ErrInternalServerError(nil)

		assert.Equal(t, http.StatusInternalServerError, result.HTTPCode)
		assert.Equal(t, ErrInternal.Message, result.Message)
		assert.Equal(t, ErrInternal.Code, result.Code)
		assert.NotEmpty(t, result.RequestID)
	})

	t.Run("with custom error", func(t *testing.T) {
		customErr := errors.New("database connection failed")
		result := ErrInternalServerError(customErr)

		assert.Equal(t, http.StatusInternalServerError, result.HTTPCode)
		assert.NotEmpty(t, result.RequestID)
		assert.NotNil(t, result.Internal)
	})

	t.Run("with IError", func(t *testing.T) {
		customErr := IError{Code: "50001", Message: "custom internal error"}
		result := ErrInternalServerError(customErr)

		assert.Equal(t, http.StatusInternalServerError, result.HTTPCode)
		assert.Equal(t, "50001", result.Code)
		assert.Equal(t, "custom internal error", result.Message)
	})
}

func TestErrUnauthorized(t *testing.T) {
	t.Run("with nil error", func(t *testing.T) {
		result := ErrUnauthorized(nil)

		assert.Equal(t, http.StatusUnauthorized, result.HTTPCode)
		assert.Equal(t, ErrUnauthorizedType.Message, result.Message)
		assert.Equal(t, ErrUnauthorizedType.Code, result.Code)
	})

	t.Run("with custom error", func(t *testing.T) {
		customErr := errors.New("invalid token")
		result := ErrUnauthorized(customErr)

		assert.Equal(t, http.StatusUnauthorized, result.HTTPCode)
		assert.NotEmpty(t, result.RequestID)
	})
}

func TestErrForbidden(t *testing.T) {
	t.Run("with nil error", func(t *testing.T) {
		result := ErrForbidden(nil)

		assert.Equal(t, http.StatusForbidden, result.HTTPCode)
		assert.Equal(t, ErrForbiddenType.Message, result.Message)
	})

	t.Run("with IError", func(t *testing.T) {
		customErr := IError{Code: "40301", Message: "insufficient permissions"}
		result := ErrForbidden(customErr)

		assert.Equal(t, http.StatusForbidden, result.HTTPCode)
		assert.Equal(t, "40301", result.Code)
	})
}

func TestErrSessionExpired(t *testing.T) {
	result := ErrSessionExpired(nil)

	assert.Equal(t, 440, result.HTTPCode)
	assert.Equal(t, ErrSessionExpiredType.Message, result.Message)
	assert.Equal(t, ErrSessionExpiredType.Code, result.Code)
}

func TestErrNotFound(t *testing.T) {
	t.Run("with nil error", func(t *testing.T) {
		result := ErrNotFound(nil)

		assert.Equal(t, http.StatusNotFound, result.HTTPCode)
		assert.Equal(t, ErrResourceNotFound.Message, result.Message)
	})

	t.Run("with custom error", func(t *testing.T) {
		customErr := errors.New("user not found")
		result := ErrNotFound(customErr)

		assert.Equal(t, http.StatusNotFound, result.HTTPCode)
		assert.NotEmpty(t, result.RequestID)
	})
}

func TestErrBadRequest(t *testing.T) {
	t.Run("with nil error", func(t *testing.T) {
		result := ErrBadRequest(nil)

		assert.Equal(t, http.StatusBadRequest, result.HTTPCode)
		assert.Equal(t, ErrBadRequestType.Message, result.Message)
	})

	t.Run("with IError", func(t *testing.T) {
		customErr := IError{Code: "40010", Message: "invalid input"}
		result := ErrBadRequest(customErr)

		assert.Equal(t, http.StatusBadRequest, result.HTTPCode)
		assert.Equal(t, "40010", result.Code)
		assert.Equal(t, "invalid input", result.Message)
	})
}

func TestErrWithData(t *testing.T) {
	t.Run("with data", func(t *testing.T) {
		customErr := errors.New("validation failed")
		data := map[string]string{"field": "email", "error": "invalid format"}

		result := ErrWithData(customErr, data, http.StatusBadRequest)

		assert.Equal(t, http.StatusBadRequest, result.HTTPCode)
		assert.NotNil(t, result.Data)
		assert.Equal(t, data, result.Data)
	})

	t.Run("with nil data", func(t *testing.T) {
		customErr := errors.New("error")
		result := ErrWithData(customErr, nil, http.StatusBadRequest)

		assert.NotNil(t, result.Data)
		assert.Equal(t, map[string]interface{}{}, result.Data)
	})
}

func TestHTTPError(t *testing.T) {
	err := errors.New("custom error")
	result := HTTPError(err, http.StatusTeapot, "41800", "I'm a teapot")

	assert.Equal(t, http.StatusTeapot, result.HTTPCode)
	assert.Equal(t, "41800", result.Code)
	assert.Equal(t, "I'm a teapot", result.Message)
	assert.NotEmpty(t, result.RequestID)
}

func TestGenerateResponseFromIError(t *testing.T) {
	t.Run("with nil error", func(t *testing.T) {
		result := GenerateResponseFromIError(nil)
		assert.Equal(t, http.StatusInternalServerError, result.HTTPCode)
	})

	t.Run("with unauthorized IError", func(t *testing.T) {
		err := ErrUnauthorizedType
		result := GenerateResponseFromIError(err)
		assert.Equal(t, http.StatusUnauthorized, result.HTTPCode)
	})

	t.Run("with forbidden IError", func(t *testing.T) {
		err := ErrForbiddenType
		result := GenerateResponseFromIError(err)
		assert.Equal(t, http.StatusForbidden, result.HTTPCode)
	})

	t.Run("with session expired IError", func(t *testing.T) {
		err := ErrSessionExpiredType
		result := GenerateResponseFromIError(err)
		assert.Equal(t, 440, result.HTTPCode)
	})

	t.Run("with not found IError", func(t *testing.T) {
		err := ErrResourceNotFound
		result := GenerateResponseFromIError(err)
		assert.Equal(t, http.StatusNotFound, result.HTTPCode)
	})

	t.Run("with internal error IError", func(t *testing.T) {
		err := ErrInternal
		result := GenerateResponseFromIError(err)
		assert.Equal(t, http.StatusInternalServerError, result.HTTPCode)
	})

	t.Run("with unknown IError code", func(t *testing.T) {
		err := IError{Code: "99999", Message: "unknown error"}
		result := GenerateResponseFromIError(err)
		assert.Equal(t, http.StatusBadRequest, result.HTTPCode)
	})

	t.Run("with error message starting with 401", func(t *testing.T) {
		err := errors.New("401 unauthorized access")
		result := GenerateResponseFromIError(err)
		assert.Equal(t, http.StatusUnauthorized, result.HTTPCode)
	})

	t.Run("with error message starting with 403", func(t *testing.T) {
		err := errors.New("403 forbidden")
		result := GenerateResponseFromIError(err)
		assert.Equal(t, http.StatusForbidden, result.HTTPCode)
	})

	t.Run("with error message starting with 404", func(t *testing.T) {
		err := errors.New("404 not found")
		result := GenerateResponseFromIError(err)
		assert.Equal(t, http.StatusNotFound, result.HTTPCode)
	})

	t.Run("with error message starting with 500", func(t *testing.T) {
		err := errors.New("500 internal server error")
		result := GenerateResponseFromIError(err)
		assert.Equal(t, http.StatusInternalServerError, result.HTTPCode)
	})

	t.Run("with generic error", func(t *testing.T) {
		err := errors.New("some error")
		result := GenerateResponseFromIError(err)
		assert.Equal(t, http.StatusBadRequest, result.HTTPCode)
	})
}

func TestNew(t *testing.T) {
	err := New("12345", "custom error message")

	assert.Equal(t, "12345", err.Code)
	assert.Equal(t, "custom error message", err.Message)
	assert.Equal(t, "custom error message", err.Error())
}

func TestWrap(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := Wrap(originalErr, "additional context")

	require.NotNil(t, wrappedErr)
	assert.Contains(t, wrappedErr.Error(), "additional context")
	assert.Contains(t, wrappedErr.Error(), "original error")
}

func TestWrapf(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := Wrapf(originalErr, "context with %s", "formatting")

	require.NotNil(t, wrappedErr)
	assert.Contains(t, wrappedErr.Error(), "context with formatting")
}

func TestCause(t *testing.T) {
	originalErr := errors.New("root cause")
	wrappedErr := Wrap(originalErr, "wrapper")

	cause := Cause(wrappedErr)
	assert.Equal(t, originalErr, cause)
}

func TestWithStack(t *testing.T) {
	err := errors.New("error")
	stackedErr := WithStack(err)

	require.NotNil(t, stackedErr)
	assert.NotEqual(t, err, stackedErr)
}

func TestWithMessage(t *testing.T) {
	err := errors.New("original")
	annotatedErr := WithMessage(err, "additional message")

	require.NotNil(t, annotatedErr)
	assert.Contains(t, annotatedErr.Error(), "additional message")
}

func TestWithMessagef(t *testing.T) {
	err := errors.New("original")
	annotatedErr := WithMessagef(err, "formatted %s", "message")

	require.NotNil(t, annotatedErr)
	assert.Contains(t, annotatedErr.Error(), "formatted message")
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("field is required")

	assert.Equal(t, ErrValidationFailed.Code, err.Code)
	assert.Equal(t, "field is required", err.Message)
}

func TestNewDuplicateEntryError(t *testing.T) {
	err := NewDuplicateEntryError("email already exists")

	assert.Equal(t, ErrDuplicateEntry.Code, err.Code)
	assert.Equal(t, "email already exists", err.Message)
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("User")

	assert.Equal(t, ErrResourceNotFound.Code, err.Code)
	assert.Equal(t, "User not found", err.Message)
}

func TestNewUnauthorizedError(t *testing.T) {
	t.Run("with custom message", func(t *testing.T) {
		err := NewUnauthorizedError("invalid credentials")

		assert.Equal(t, ErrUnauthorizedType.Code, err.Code)
		assert.Equal(t, "invalid credentials", err.Message)
	})

	t.Run("with empty message", func(t *testing.T) {
		err := NewUnauthorizedError("")

		assert.Equal(t, ErrUnauthorizedType.Code, err.Code)
		assert.Equal(t, ErrUnauthorizedType.Message, err.Message)
	})
}

func TestNewForbiddenError(t *testing.T) {
	t.Run("with custom message", func(t *testing.T) {
		err := NewForbiddenError("access denied")

		assert.Equal(t, ErrForbiddenType.Code, err.Code)
		assert.Equal(t, "access denied", err.Message)
	})

	t.Run("with empty message", func(t *testing.T) {
		err := NewForbiddenError("")

		assert.Equal(t, ErrForbiddenType.Code, err.Code)
		assert.Equal(t, ErrForbiddenType.Message, err.Message)
	})
}
