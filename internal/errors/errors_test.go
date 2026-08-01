package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestAPIError(t *testing.T) {
	err := NewAPIError(400, "invalid_request", "Invalid request format", "Missing required field")
	
	expected := "API error 400 (invalid_request): Invalid request format - Missing required field"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
	
	// Test without details
	err = NewAPIError(404, "not_found", "Resource not found", "")
	expected = "API error 404 (not_found): Resource not found"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestAuthenticationError(t *testing.T) {
	// Test with wrapped error
	innerErr := errors.New("connection failed")
	err := NewAuthError("token expired", innerErr)
	
	expected := "authentication failed: token expired: connection failed"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
	
	// Test unwrapping
	if unwrapped := err.Unwrap(); unwrapped != innerErr {
		t.Errorf("Expected unwrapped error to be %v, got %v", innerErr, unwrapped)
	}
	
	// Test without wrapped error
	err = NewAuthError("invalid credentials", nil)
	expected = "authentication failed: invalid credentials"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestNetworkError(t *testing.T) {
	innerErr := errors.New("connection timeout")
	err := NewNetworkError("GET /devices", innerErr)
	
	expected := "network error during GET /devices: connection timeout"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
	
	// Test unwrapping
	if unwrapped := err.Unwrap(); unwrapped != innerErr {
		t.Errorf("Expected unwrapped error to be %v, got %v", innerErr, unwrapped)
	}
}

func TestConfigurationError(t *testing.T) {
	err := NewConfigError("ClientID", "field is required", "set BS_CLIENT_ID environment variable")
	
	expected := "configuration error: ClientID - field is required (suggestion: set BS_CLIENT_ID environment variable)"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
	
	// Test without suggestion
	err = NewConfigError("Timeout", "must be positive", "")
	expected = "configuration error: Timeout - must be positive"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestValidationError(t *testing.T) {
	err := NewValidationError("networkName", "", "cannot be empty")
	
	expected := "validation error: networkName= - cannot be empty"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestIsAuthenticationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "AuthenticationError",
			err:      NewAuthError("test", nil),
			expected: true,
		},
		{
			name:     "APIError 401",
			err:      NewAPIError(http.StatusUnauthorized, "unauthorized", "Invalid token", ""),
			expected: true,
		},
		{
			name:     "APIError 403",
			err:      NewAPIError(http.StatusForbidden, "forbidden", "Insufficient permissions", ""),
			expected: true,
		},
		{
			name:     "APIError 400",
			err:      NewAPIError(http.StatusBadRequest, "bad_request", "Invalid request", ""),
			expected: false,
		},
		{
			name:     "NetworkError",
			err:      NewNetworkError("test", errors.New("connection failed")),
			expected: false,
		},
		{
			name:     "Generic error",
			err:      errors.New("some error"),
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAuthenticationError(tt.err)
			if result != tt.expected {
				t.Errorf("IsAuthenticationError() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "APIError 500",
			err:      NewAPIError(http.StatusInternalServerError, "internal_error", "Server error", ""),
			expected: true,
		},
		{
			name:     "APIError 502",
			err:      NewAPIError(http.StatusBadGateway, "bad_gateway", "Bad gateway", ""),
			expected: true,
		},
		{
			name:     "APIError 429",
			err:      NewAPIError(http.StatusTooManyRequests, "rate_limit", "Rate limit exceeded", ""),
			expected: true,
		},
		{
			name:     "APIError 400",
			err:      NewAPIError(http.StatusBadRequest, "bad_request", "Invalid request", ""),
			expected: false,
		},
		{
			name:     "APIError 404",
			err:      NewAPIError(http.StatusNotFound, "not_found", "Not found", ""),
			expected: false,
		},
		{
			name:     "NetworkError",
			err:      NewNetworkError("test", errors.New("connection failed")),
			expected: true,
		},
		{
			name:     "AuthenticationError",
			err:      NewAuthError("test", nil),
			expected: false,
		},
		{
			name:     "Generic error",
			err:      errors.New("some error"),
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRetryableError(tt.err)
			if result != tt.expected {
				t.Errorf("IsRetryableError() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
func TestNewAPIErrorFrom_InheritsStatusFromCause(t *testing.T) {
	// A service method re-wrapping a transport failure must not lose the HTTP
	// status: without it, callers are left matching on message text to tell a
	// rate limit apart from a bad credential.
	cause := NewAPIError(http.StatusTooManyRequests, "Too Many Requests", "Request failed", "slow down")

	wrapped := NewAPIErrorFrom("network_context_failed", "Failed to set network context to 'lab'", cause)

	if wrapped.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status code 429, got %d", wrapped.StatusCode)
	}
	if wrapped.Code != "network_context_failed" {
		t.Errorf("Expected code 'network_context_failed', got %q", wrapped.Code)
	}
	if wrapped.Details != cause.Error() {
		t.Errorf("Expected details %q, got %q", cause.Error(), wrapped.Details)
	}

	var unwrapped *APIError
	if !errors.As(errors.Unwrap(wrapped), &unwrapped) {
		t.Fatal("Expected the cause to be reachable by unwrapping")
	}
	if unwrapped != cause {
		t.Error("Expected to unwrap to the original cause")
	}

	if !IsRateLimited(wrapped) {
		t.Error("Expected IsRateLimited to report true for a wrapped 429")
	}
	if !IsRetryableError(wrapped) {
		t.Error("Expected IsRetryableError to report true for a wrapped 429")
	}
}

func TestNewAPIErrorFrom_InheritsStatusThroughAuthenticationError(t *testing.T) {
	// The HTTP client wraps 401 and 403 in an AuthenticationError, so the
	// status is two levels down once a service re-wraps it.
	cause := NewAuthError("invalid or expired token",
		NewAPIError(http.StatusUnauthorized, "Unauthorized", "Request failed", ""))

	wrapped := NewAPIErrorFrom("rdws_info_failed", "Failed to get info", cause)

	if wrapped.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code 401, got %d", wrapped.StatusCode)
	}
	if got := StatusCodeOf(wrapped); got != http.StatusUnauthorized {
		t.Errorf("Expected StatusCodeOf to report 401, got %d", got)
	}
	if !IsUnauthorizedError(wrapped) {
		t.Error("Expected IsUnauthorizedError to report true for a wrapped 401")
	}
	if !IsAuthenticationError(wrapped) {
		t.Error("Expected IsAuthenticationError to report true for a wrapped 401")
	}
	if IsRateLimited(wrapped) {
		t.Error("Expected IsRateLimited to report false for a 401")
	}
}

func TestNewAPIErrorFrom_NoStatusAvailable(t *testing.T) {
	cause := errors.New("connection reset by peer")

	wrapped := NewAPIErrorFrom("rdws_logs_failed", "Failed to get logs", cause)

	if wrapped.StatusCode != 0 {
		t.Errorf("Expected status code 0 when the cause carries none, got %d", wrapped.StatusCode)
	}
	if StatusCodeOf(wrapped) != 0 {
		t.Errorf("Expected StatusCodeOf to report 0, got %d", StatusCodeOf(wrapped))
	}
	if IsRateLimited(wrapped) {
		t.Error("Expected IsRateLimited to report false when no status is available")
	}
	if !errors.Is(wrapped, cause) {
		t.Error("Expected the cause to remain reachable with errors.Is")
	}
}

func TestNewAPIErrorFrom_NilCause(t *testing.T) {
	wrapped := NewAPIErrorFrom("some_failure", "Something went wrong", nil)

	if wrapped.StatusCode != 0 {
		t.Errorf("Expected status code 0 for a nil cause, got %d", wrapped.StatusCode)
	}
	if wrapped.Details != "" {
		t.Errorf("Expected empty details for a nil cause, got %q", wrapped.Details)
	}
	if wrapped.Unwrap() != nil {
		t.Error("Expected Unwrap to return nil for a nil cause")
	}
}

func TestStatusCodeOf_DirectAPIError(t *testing.T) {
	if got := StatusCodeOf(NewAPIError(http.StatusNotFound, "Not Found", "Request failed", "")); got != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", got)
	}
	if got := StatusCodeOf(nil); got != 0 {
		t.Errorf("Expected 0 for a nil error, got %d", got)
	}
}
