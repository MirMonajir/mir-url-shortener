package domain

import "fmt"

// ErrorType defines the type of error
type ErrorType string

const (
	ErrInvalidURL      ErrorType = "invalid_url"
	ErrURLNotFound     ErrorType = "url_not_found"
	ErrShortCodeExists ErrorType = "short_code_exists"
	ErrEmptyURL        ErrorType = "empty_url"
	ErrInvalidInput    ErrorType = "invalid_input"
	ErrInternal        ErrorType = "internal_error"
)

// AppError represents an application error with type and message
type AppError struct {
	Type    ErrorType `json:"error_type"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
	Code    int       `json:"code"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(errorType ErrorType, message string, code int) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Code:    code,
	}
}

// NewAppErrorWithDetails creates an error with additional details
func NewAppErrorWithDetails(errorType ErrorType, message, details string, code int) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Details: details,
		Code:    code,
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Type    ErrorType         `json:"error_type"`
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
	Code    int               `json:"code"`
}

// Error implements the error interface
func (e *ValidationErrors) Error() string {
	return fmt.Sprintf("validation error: %s", e.Message)
}

// NewValidationErrors creates validation errors
func NewValidationErrors(errors []ValidationError) *ValidationErrors {
	return &ValidationErrors{
		Type:    ErrInvalidInput,
		Message: "Request validation failed",
		Errors:  errors,
		Code:    400,
	}
}
