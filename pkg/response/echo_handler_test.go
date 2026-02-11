package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCustomHTTPErrorHandler(t *testing.T) {
	e := echo.New()

	t.Run("handle ErrorResponse", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := ErrorResponse{
			HTTPCode: http.StatusTeapot,
			Message:  "I'm a teapot",
		}

		CustomHTTPErrorHandler(err, c)

		assert.Equal(t, http.StatusTeapot, rec.Code)
		assert.Contains(t, rec.Body.String(), "I'm a teapot")
	})

	t.Run("handle generic error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := errors.New("something went wrong")

		CustomHTTPErrorHandler(err, c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Bad request")
	})
}

func TestEchoMiddleware(t *testing.T) {
	e := echo.New()
	middleware := EchoMiddleware()

	t.Run("successful request with request ID", func(t *testing.T) {
		handler := middleware(func(c echo.Context) error {
			return OK(c, "success", nil)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)
		assert.NoError(t, err)
		assert.NotEmpty(t, rec.Header().Get("X-Request-ID"))
	})

	t.Run("handler returns error", func(t *testing.T) {
		handler := middleware(func(c echo.Context) error {
			return errors.New("oops")
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandleError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	t.Run("nil error", func(t *testing.T) {
		err := HandleError(c, nil)
		assert.NoError(t, err)
	})

	t.Run("ErrorResponse", func(t *testing.T) {
		origErr := ErrorResponse{HTTPCode: http.StatusForbidden, Message: "forbidden"}
		err := HandleError(c, origErr)
		assert.Equal(t, origErr, err)
	})

	t.Run("generic error", func(t *testing.T) {
		origErr := errors.New("bad thing")
		err := HandleError(c, origErr)
		assert.IsType(t, ErrorResponse{}, err)
		assert.Equal(t, http.StatusBadRequest, err.(ErrorResponse).HTTPCode)
	})
}

func TestSuccessResponses(t *testing.T) {
	e := echo.New()

	t.Run("OK response", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)

		err := OK(c, "all good", map[string]string{"foo": "bar"})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "all good")
		assert.Contains(t, rec.Body.String(), "bar")
	})

	t.Run("Created response", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c := e.NewContext(httptest.NewRequest(http.MethodPost, "/", nil), rec)

		err := Created(c, "created", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("NoContent response", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c := e.NewContext(httptest.NewRequest(http.MethodDelete, "/", nil), rec)

		err := NoContent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}

func TestValidationError(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodPost, "/", nil), rec)

	err := ValidationError(c, errors.New("invalid"), map[string]string{"email": "required"})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "required")
}
