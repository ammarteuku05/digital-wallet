package response

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// CustomHTTPErrorHandler is a custom error handler for Echo framework
func CustomHTTPErrorHandler(err error, c echo.Context) {
	var errorResponse ErrorResponse

	// Check if it's already our custom ErrorResponse
	if customErr, ok := err.(ErrorResponse); ok {
		errorResponse = customErr
	} else {
		// Convert standard errors to our custom format
		errorResponse = GenerateResponseFromIError(err)
	}

	// Log the error for debugging (you might want to use a proper logger)
	if errorResponse.Internal != nil {
		log.Printf("Error: %v", errorResponse.Internal)
	}

	// Send the error response
	if !c.Response().Committed {
		if err := c.JSON(errorResponse.StatusCode(), errorResponse); err != nil {
			log.Printf("Failed to send error response: %v", err)
		}
	}
}

// EchoMiddleware creates an Echo middleware for error handling
func EchoMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Add request ID to context if not already present
			if c.Request().Header.Get("X-Request-ID") == "" {
				requestID := uuid.New().String()
				c.Request().Header.Set("X-Request-ID", requestID)
				c.Response().Header().Set("X-Request-ID", requestID)
			}

			// Call the next handler
			err := next(c)
			if err != nil {
				// Handle the error using our custom handler
				CustomHTTPErrorHandler(err, c)
				return nil // Don't propagate the error to Echo's default handler
			}

			return nil
		}
	}
}

// Helper functions for controllers

// HandleError is a helper function for controllers to handle errors
func HandleError(c echo.Context, err error) error {
	if err == nil {
		return nil
	}

	// Check if it's already our custom ErrorResponse
	if _, ok := err.(ErrorResponse); ok {
		return err
	}

	// Convert to our custom error format
	return GenerateResponseFromIError(err)
}

// SendSuccessResponse is a helper function for controllers to send success responses
func SendSuccessResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	response := map[string]interface{}{
		"code":    "20000",
		"message": message,
	}

	if data != nil {
		response["data"] = data
	}

	return c.JSON(statusCode, response)
}

// ValidationError sends a 400 Bad Request response with validation details
func ValidationError(c echo.Context, err error, validationData interface{}) error {
	errorResponse := ErrWithData(err, validationData, http.StatusBadRequest)
	return c.JSON(errorResponse.StatusCode(), errorResponse)
}

// Created sends a 201 Created response
func Created(c echo.Context, message string, data interface{}) error {
	return SendSuccessResponse(c, http.StatusCreated, message, data)
}

// OK sends a 200 OK response
func OK(c echo.Context, message string, data interface{}) error {
	return SendSuccessResponse(c, http.StatusOK, message, data)
}

// NoContent sends a 204 No Content response
func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
