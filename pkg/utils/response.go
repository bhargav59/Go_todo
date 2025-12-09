package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standardized API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents error details in API response
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Common error codes
const (
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeForbidden      = "FORBIDDEN"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeConflict       = "CONFLICT"
	ErrCodeInternal       = "INTERNAL_ERROR"
	ErrCodeBadRequest     = "BAD_REQUEST"
)

// Success sends a successful response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, code string, message string, details interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, details interface{}) {
	Error(c, http.StatusBadRequest, ErrCodeValidation, "Validation failed", details)
}

// UnauthorizedError sends an unauthorized error response
func UnauthorizedError(c *gin.Context, message string) {
	if message == "" {
		message = "Authentication required"
	}
	Error(c, http.StatusUnauthorized, ErrCodeUnauthorized, message, nil)
}

// ForbiddenError sends a forbidden error response
func ForbiddenError(c *gin.Context, message string) {
	if message == "" {
		message = "Access denied"
	}
	Error(c, http.StatusForbidden, ErrCodeForbidden, message, nil)
}

// NotFoundError sends a not found error response
func NotFoundError(c *gin.Context, resource string) {
	Error(c, http.StatusNotFound, ErrCodeNotFound, resource+" not found", nil)
}

// ConflictError sends a conflict error response
func ConflictError(c *gin.Context, message string) {
	Error(c, http.StatusConflict, ErrCodeConflict, message, nil)
}

// InternalError sends an internal server error response
func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = "An internal error occurred"
	}
	Error(c, http.StatusInternalServerError, ErrCodeInternal, message, nil)
}

// BadRequestError sends a bad request error response
func BadRequestError(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, ErrCodeBadRequest, message, nil)
}

// Created sends a 201 created response
func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data)
}

// OK sends a 200 OK response
func OK(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusOK, message, data)
}

// NoContent sends a 204 No Content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
