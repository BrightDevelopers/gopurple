package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/brightsign/gopurple/internal/config"
	"github.com/brightsign/gopurple/internal/errors"
	"github.com/go-resty/resty/v2"
)

// HTTPClient wraps the resty client with BSN.cloud-specific functionality.
type HTTPClient struct {
	client *resty.Client
	config *config.Config
}

// NewHTTPClient creates a new HTTP client with the given configuration.
func NewHTTPClient(cfg *config.Config) *HTTPClient {
	client := resty.New().
		SetTimeout(cfg.Timeout).
		SetRetryCount(cfg.RetryCount).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			// Exponential backoff for retries
			return time.Duration(resp.Request.Attempt) * time.Second, nil
		}).
		AddRetryCondition(func(resp *resty.Response, err error) bool {
			// Retry on network errors
			if err != nil {
				return true
			}
			// Retry on server errors and rate limiting
			return resp.StatusCode() >= 500 || resp.StatusCode() == http.StatusTooManyRequests
		}).
		SetHeaders(map[string]string{
			"Accept":     "application/json",
			"User-Agent": "gopurple-sdk/1.0",
		}).
		SetAuthScheme("Bearer"). // Required for SetAuthToken() to work correctly
		SetDebug(cfg.Debug)      // Enable debug logging if configured

	return &HTTPClient{
		client: client,
		config: cfg,
	}
}

// Request represents an HTTP request to be made.
type Request struct {
	Method string
	URL    string
	Body   interface{}
	Result interface{}
	Headers map[string]string
	QueryParams map[string]string
}

// Do executes an HTTP request with error handling and response parsing.
func (h *HTTPClient) Do(ctx context.Context, req *Request) error {
	request := h.client.R().SetContext(ctx)

	// Set headers
	if req.Headers != nil {
		request.SetHeaders(req.Headers)
	}

	// Set query parameters
	if req.QueryParams != nil {
		request.SetQueryParams(req.QueryParams)
	}

	// Set body
	if req.Body != nil {
		request.SetBody(req.Body)
		request.SetHeader("Content-Type", "application/json")
	}

	// Set result structure
	if req.Result != nil {
		request.SetResult(req.Result)
	}

	// Set error structure to capture API errors
	var apiError struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
		Details          string `json:"details"`
	}
	request.SetError(&apiError)

	// Execute request
	resp, err := request.Execute(req.Method, req.URL)
	if err != nil {
		return errors.NewNetworkError(fmt.Sprintf("%s %s", req.Method, req.URL), err)
	}

	// Check for API errors
	if !resp.IsSuccess() {
		return h.handleAPIError(resp, &apiError)
	}

	return nil
}

// DoWithAuth executes an HTTP request with authentication headers.
func (h *HTTPClient) DoWithAuth(ctx context.Context, token string, req *Request) error {
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}
	req.Headers["Authorization"] = "Bearer " + token
	return h.Do(ctx, req)
}

// handleAPIError converts HTTP error responses to appropriate error types.
func (h *HTTPClient) handleAPIError(resp *resty.Response, apiError *struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	Details          string `json:"details"`
}) error {
	statusCode := resp.StatusCode()

	// Try to get error details from response
	errorCode := apiError.Error
	errorMessage := apiError.ErrorDescription
	errorDetails := apiError.Details

	// If we couldn't parse the error, include the raw response body
	if errorCode == "" && errorMessage == "" {
		bodyStr := string(resp.Body())
		if bodyStr != "" && len(bodyStr) < 500 {
			errorDetails = bodyStr
		}
	}

	// Fallback to status text if no error details
	if errorCode == "" {
		errorCode = http.StatusText(statusCode)
	}
	if errorMessage == "" {
		errorMessage = "Request failed"
	}

	// Create specific error types based on status code
	switch statusCode {
	case http.StatusUnauthorized:
		return errors.NewAuthError("invalid or expired token", errors.NewAPIError(statusCode, errorCode, errorMessage, errorDetails))
	case http.StatusForbidden:
		return errors.NewAuthError("insufficient permissions", errors.NewAPIError(statusCode, errorCode, errorMessage, errorDetails))
	default:
		return errors.NewAPIError(statusCode, errorCode, errorMessage, errorDetails)
	}
}

// Get performs a GET request.
func (h *HTTPClient) Get(ctx context.Context, url string, result interface{}) error {
	return h.Do(ctx, &Request{
		Method: "GET",
		URL:    url,
		Result: result,
	})
}

// GetWithAuth performs a GET request with authentication.
func (h *HTTPClient) GetWithAuth(ctx context.Context, token, url string, result interface{}) error {
	return h.DoWithAuth(ctx, token, &Request{
		Method: "GET",
		URL:    url,
		Result: result,
	})
}

// GetBytesWithAuth performs a GET request and returns raw bytes (for downloading files).
func (h *HTTPClient) GetBytesWithAuth(ctx context.Context, token, url string) ([]byte, error) {
	request := h.client.R().
		SetContext(ctx).
		SetAuthToken(token)

	// Set error structure
	var apiError struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
		Details          string `json:"details"`
	}
	request.SetError(&apiError)

	// Execute request
	resp, err := request.Get(url)
	if err != nil {
		return nil, errors.NewNetworkError(fmt.Sprintf("GET %s", url), err)
	}

	// Check for HTTP errors
	if resp.StatusCode() >= 400 {
		return nil, h.handleAPIError(resp, &apiError)
	}

	return resp.Body(), nil
}

// Post performs a POST request.
func (h *HTTPClient) Post(ctx context.Context, url string, body, result interface{}) error {
	return h.Do(ctx, &Request{
		Method: "POST",
		URL:    url,
		Body:   body,
		Result: result,
	})
}

// PostWithAuth performs a POST request with authentication.
func (h *HTTPClient) PostWithAuth(ctx context.Context, token, url string, body, result interface{}) error {
	return h.DoWithAuth(ctx, token, &Request{
		Method: "POST",
		URL:    url,
		Body:   body,
		Result: result,
	})
}

// Put performs a PUT request.
func (h *HTTPClient) Put(ctx context.Context, url string, body, result interface{}) error {
	return h.Do(ctx, &Request{
		Method: "PUT",
		URL:    url,
		Body:   body,
		Result: result,
	})
}

// PutWithAuth performs a PUT request with authentication.
func (h *HTTPClient) PutWithAuth(ctx context.Context, token, url string, body, result interface{}) error {
	return h.DoWithAuth(ctx, token, &Request{
		Method: "PUT",
		URL:    url,
		Body:   body,
		Result: result,
	})
}

// DeleteWithAuth performs a DELETE request with authentication.
func (h *HTTPClient) DeleteWithAuth(ctx context.Context, token, url string, result interface{}) error {
	return h.DoWithAuth(ctx, token, &Request{
		Method: "DELETE",
		URL:    url,
		Result: result,
	})
}

// PostForm performs a POST request with form data (for OAuth token requests).
func (h *HTTPClient) PostForm(ctx context.Context, url string, data map[string]string, result interface{}) error {
	request := h.client.R().SetContext(ctx)
	
	// Set form data
	request.SetFormData(data)
	request.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	
	// Set result
	if result != nil {
		request.SetResult(result)
	}
	
	// Set error structure
	var apiError struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	request.SetError(&apiError)
	
	resp, err := request.Post(url)
	if err != nil {
		return errors.NewNetworkError(fmt.Sprintf("POST %s", url), err)
	}
	
	if !resp.IsSuccess() {
		return h.handleAPIError(resp, &struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
			Details          string `json:"details"`
		}{
			Error:            apiError.Error,
			ErrorDescription: apiError.ErrorDescription,
		})
	}
	
	return nil
}

// PostFormWithAuth performs a POST request with form data and basic auth.
func (h *HTTPClient) PostFormWithAuth(ctx context.Context, clientID, clientSecret, url string, data map[string]string, result interface{}) error {
	request := h.client.R().SetContext(ctx)

	// Set basic auth
	request.SetBasicAuth(clientID, clientSecret)

	// Set form data
	request.SetFormData(data)
	request.SetHeader("Content-Type", "application/x-www-form-urlencoded")

	// Set result
	if result != nil {
		request.SetResult(result)
	}

	// Set error structure
	var apiError struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	request.SetError(&apiError)

	resp, err := request.Post(url)
	if err != nil {
		return errors.NewNetworkError(fmt.Sprintf("POST %s", url), err)
	}

	if !resp.IsSuccess() {
		return h.handleAPIError(resp, &struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
			Details          string `json:"details"`
		}{
			Error:            apiError.Error,
			ErrorDescription: apiError.ErrorDescription,
		})
	}

	return nil
}

// GetClient returns the underlying resty client for advanced usage.
func (h *HTTPClient) GetClient() *resty.Client {
	return h.client
}