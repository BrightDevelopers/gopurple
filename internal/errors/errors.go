package errors

import (
	"fmt"
	"net/http"
)

// Error types for different categories of API errors

// APIError represents an error response from the BSN.cloud API.
type APIError struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"error"`
	Message    string `json:"error_description"`
	Details    string `json:"details,omitempty"`

	// Err is the underlying error this one was built from, if any. It is kept
	// so that the HTTP status of a transport-level failure survives being
	// re-wrapped by a service method, instead of only appearing in Details.
	Err error `json:"-"`
}

func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("API error %d (%s): %s - %s", e.StatusCode, e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("API error %d (%s): %s", e.StatusCode, e.Code, e.Message)
}

// Unwrap returns the underlying error, so errors.As and errors.Is can reach it.
func (e *APIError) Unwrap() error {
	return e.Err
}

// AuthenticationError indicates authentication-related errors.
type AuthenticationError struct {
	Reason string
	Err    error
}

func (e *AuthenticationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("authentication failed: %s: %v", e.Reason, e.Err)
	}
	return fmt.Sprintf("authentication failed: %s", e.Reason)
}

func (e *AuthenticationError) Unwrap() error {
	return e.Err
}

// NetworkError indicates network-related errors.
type NetworkError struct {
	Operation string
	Err       error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error during %s: %v", e.Operation, e.Err)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

// ConfigurationError indicates configuration-related errors.
type ConfigurationError struct {
	Field   string
	Reason  string
	Suggest string
}

func (e *ConfigurationError) Error() string {
	msg := fmt.Sprintf("configuration error: %s - %s", e.Field, e.Reason)
	if e.Suggest != "" {
		msg += fmt.Sprintf(" (suggestion: %s)", e.Suggest)
	}
	return msg
}

// ValidationError indicates data validation errors.
type ValidationError struct {
	Field  string
	Value  interface{}
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s=%v - %s", e.Field, e.Value, e.Reason)
}

// Helper functions for creating common errors

// NewAPIError creates an APIError from an HTTP response.
func NewAPIError(statusCode int, code, message, details string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Details:    details,
	}
}

// NewAPIErrorFrom creates an APIError that keeps cause reachable through
// errors.As and inherits cause's HTTP status when cause carries one.
//
// Prefer this over NewAPIError(0, ...) when re-wrapping an error returned by
// the HTTP client: with NewAPIError the real status is discarded and survives
// only as text inside Details, which forces callers to match on strings to
// tell a 429 from a 401.
func NewAPIErrorFrom(code, message string, cause error) *APIError {
	apiErr := &APIError{
		Code:    code,
		Message: message,
		Err:     cause,
	}
	if cause != nil {
		apiErr.Details = cause.Error()
		apiErr.StatusCode = StatusCodeOf(cause)
	}
	return apiErr
}

// StatusCodeOf returns the HTTP status code carried by err or any error it
// wraps, or 0 when no status is available.
func StatusCodeOf(err error) int {
	for err != nil {
		if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode != 0 {
			return apiErr.StatusCode
		}
		unwrapper, ok := err.(interface{ Unwrap() error })
		if !ok {
			return 0
		}
		err = unwrapper.Unwrap()
	}
	return 0
}

// IsRateLimited reports whether err was caused by an HTTP 429 response.
func IsRateLimited(err error) bool {
	return StatusCodeOf(err) == http.StatusTooManyRequests
}

// IsUnauthorizedError reports whether err was caused by an HTTP 401 or 403 response.
func IsUnauthorizedError(err error) bool {
	status := StatusCodeOf(err)
	return status == http.StatusUnauthorized || status == http.StatusForbidden
}

// NewAuthError creates an AuthenticationError.
func NewAuthError(reason string, err error) *AuthenticationError {
	return &AuthenticationError{
		Reason: reason,
		Err:    err,
	}
}

// NewNetworkError creates a NetworkError.
func NewNetworkError(operation string, err error) *NetworkError {
	return &NetworkError{
		Operation: operation,
		Err:       err,
	}
}

// NewConfigError creates a ConfigurationError.
func NewConfigError(field, reason, suggest string) *ConfigurationError {
	return &ConfigurationError{
		Field:   field,
		Reason:  reason,
		Suggest: suggest,
	}
}

// NewValidationError creates a ValidationError.
func NewValidationError(field string, value interface{}, reason string) *ValidationError {
	return &ValidationError{
		Field:  field,
		Value:  value,
		Reason: reason,
	}
}

// IsAuthenticationError checks if an error is authentication-related.
func IsAuthenticationError(err error) bool {
	_, ok := err.(*AuthenticationError)
	if ok {
		return true
	}
	
	// Also check for API errors with authentication status codes, including
	// ones only reachable by unwrapping.
	return IsUnauthorizedError(err)
}

// IsNetworkError checks if an error is network-related.
func IsNetworkError(err error) bool {
	_, ok := err.(*NetworkError)
	return ok
}

// IsConfigurationError checks if an error is configuration-related.
func IsConfigurationError(err error) bool {
	_, ok := err.(*ConfigurationError)
	return ok
}

// IsRetryableError checks if an error might succeed on retry.
func IsRetryableError(err error) bool {
	// Retry on server errors and rate limiting, including statuses only
	// reachable by unwrapping.
	if status := StatusCodeOf(err); status >= 500 || status == http.StatusTooManyRequests {
		return true
	}

	if _, ok := err.(*NetworkError); ok {
		// Retry network errors (connection failures, etc.)
		return true
	}
	
	return false
}