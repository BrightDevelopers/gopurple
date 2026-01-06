package services

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/brightdevelopers/gopurple/internal/auth"
	"github.com/brightdevelopers/gopurple/internal/config"
	"github.com/brightdevelopers/gopurple/internal/errors"
	"github.com/brightdevelopers/gopurple/internal/http"
	"github.com/brightdevelopers/gopurple/internal/types"
)

// PresentationService provides presentation management operations.
type PresentationService interface {
	// List retrieves a list of presentations on the network with optional filtering and pagination
	List(ctx context.Context, opts ...ListOption) (*types.PresentationList, error)

	// GetCount retrieves the number of presentations on the network
	GetCount(ctx context.Context) (*types.PresentationCount, error)

	// GetByID returns the specified presentation by ID
	GetByID(ctx context.Context, id int) (*types.Presentation, error)

	// GetByName returns the specified presentation by name
	GetByName(ctx context.Context, name string) (*types.Presentation, error)

	// Create creates a new presentation on the network
	Create(ctx context.Context, request *types.PresentationCreateRequest) (*types.Presentation, error)

	// Update modifies an existing presentation by ID
	Update(ctx context.Context, id int, request *types.PresentationCreateRequest) (*types.Presentation, error)

	// DeleteByID removes the specified presentation by ID
	DeleteByID(ctx context.Context, id int) error

	// DeleteByName removes the specified presentation by name
	DeleteByName(ctx context.Context, name string) error

	// DeleteByFilter removes presentations specified by a filter from the network
	DeleteByFilter(ctx context.Context, filter string) (*types.PresentationDeleteResult, error)
}

// presentationService implements the PresentationService interface.
type presentationService struct {
	config      *config.Config
	httpClient  *http.HTTPClient
	authManager *auth.AuthManager
}

// NewPresentationService creates a new presentation service.
func NewPresentationService(cfg *config.Config, httpClient *http.HTTPClient, authManager *auth.AuthManager) PresentationService {
	return &presentationService{
		config:      cfg,
		httpClient:  httpClient,
		authManager: authManager,
	}
}

// List retrieves a list of presentations with optional filtering and pagination.
func (s *presentationService) List(ctx context.Context, opts ...ListOption) (*types.PresentationList, error) {
	// Ensure we have authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return nil, err
	}

	// Build query parameters from options
	config := &listConfig{}
	for _, opt := range opts {
		opt.apply(config)
	}

	params := url.Values{}
	if config.pageSize > 0 {
		params.Set("pageSize", strconv.Itoa(config.pageSize))
	}
	if config.marker != "" {
		params.Set("marker", config.marker)
	}
	if config.filter != "" {
		params.Set("filter", config.filter)
	}
	if config.sort != "" {
		params.Set("sort", config.sort)
	}

	// Build URL
	baseURL := fmt.Sprintf("%s/%s/Presentations", s.config.BSNBaseURL, s.config.APIVersion)
	if len(params) > 0 {
		baseURL += "?" + params.Encode()
	}

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Make the request
	var result types.PresentationList
	if err := s.httpClient.GetWithAuth(ctx, token, baseURL, &result); err != nil {
		return nil, errors.NewAPIError(0, "presentation_list_failed", "Failed to list presentations", err.Error())
	}

	return &result, nil
}

// GetCount retrieves the number of presentations on the network.
func (s *presentationService) GetCount(ctx context.Context) (*types.PresentationCount, error) {
	// Ensure we have authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return nil, err
	}

	// Build URL
	countURL := fmt.Sprintf("%s/%s/Presentations/Count", s.config.BSNBaseURL, s.config.APIVersion)

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Make the request - API returns a plain number, not a JSON object
	var count int
	if err := s.httpClient.GetWithAuth(ctx, token, countURL, &count); err != nil {
		return nil, errors.NewAPIError(0, "presentation_count_failed", "Failed to get presentation count", err.Error())
	}

	return &types.PresentationCount{Count: count}, nil
}

// GetByID returns the specified presentation by ID.
func (s *presentationService) GetByID(ctx context.Context, id int) (*types.Presentation, error) {
	if id <= 0 {
		return nil, errors.NewValidationError("id", id, "presentation ID must be positive")
	}

	// Ensure we have authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return nil, err
	}

	// Build URL
	presentationURL := fmt.Sprintf("%s/%s/Presentations/%d",
		s.config.BSNBaseURL, s.config.APIVersion, id)

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Make the request
	var result types.Presentation
	if err := s.httpClient.GetWithAuth(ctx, token, presentationURL, &result); err != nil {
		return nil, errors.NewAPIError(0, "presentation_get_failed",
			fmt.Sprintf("Failed to get presentation with ID %d", id), err.Error())
	}

	return &result, nil
}

// GetByName returns the specified presentation by name.
func (s *presentationService) GetByName(ctx context.Context, name string) (*types.Presentation, error) {
	if name == "" {
		return nil, errors.NewValidationError("name", name, "presentation name cannot be empty")
	}

	// Ensure we have authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return nil, err
	}

	// Build URL - name needs to be URL encoded
	presentationURL := fmt.Sprintf("%s/%s/Presentations/%s",
		s.config.BSNBaseURL, s.config.APIVersion, url.PathEscape(name))

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Make the request
	var result types.Presentation
	if err := s.httpClient.GetWithAuth(ctx, token, presentationURL, &result); err != nil {
		return nil, errors.NewAPIError(0, "presentation_get_failed",
			fmt.Sprintf("Failed to get presentation with name %s", name), err.Error())
	}

	return &result, nil
}

// DeleteByID removes the specified presentation by ID.
func (s *presentationService) DeleteByID(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.NewValidationError("id", id, "presentation ID must be positive")
	}

	// Ensure we have authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return err
	}

	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return err
	}

	// Build URL
	presentationURL := fmt.Sprintf("%s/%s/Presentations/%d",
		s.config.BSNBaseURL, s.config.APIVersion, id)

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return err
	}

	// Make the DELETE request
	if err := s.httpClient.DeleteWithAuth(ctx, token, presentationURL, nil); err != nil {
		return errors.NewAPIError(0, "presentation_delete_failed",
			fmt.Sprintf("Failed to delete presentation with ID %d", id), err.Error())
	}

	return nil
}

// DeleteByName removes the specified presentation by name.
func (s *presentationService) DeleteByName(ctx context.Context, name string) error {
	if name == "" {
		return errors.NewValidationError("name", name, "presentation name cannot be empty")
	}

	// Ensure we have authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return err
	}

	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return err
	}

	// Build URL - name needs to be URL encoded
	presentationURL := fmt.Sprintf("%s/%s/Presentations/%s",
		s.config.BSNBaseURL, s.config.APIVersion, url.PathEscape(name))

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return err
	}

	// Make the DELETE request
	if err := s.httpClient.DeleteWithAuth(ctx, token, presentationURL, nil); err != nil {
		return errors.NewAPIError(0, "presentation_delete_failed",
			fmt.Sprintf("Failed to delete presentation with name '%s'", name), err.Error())
	}

	return nil
}

// Create creates a new presentation on the network.
func (s *presentationService) Create(ctx context.Context, request *types.PresentationCreateRequest) (*types.Presentation, error) {
	if request == nil {
		return nil, errors.NewValidationError("request", request, "presentation create request cannot be nil")
	}

	if request.Name == "" {
		return nil, errors.NewValidationError("name", request.Name, "presentation name cannot be empty")
	}

	// Ensure we have authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return nil, err
	}

	// Build URL
	createURL := fmt.Sprintf("%s/%s/Presentations", s.config.BSNBaseURL, s.config.APIVersion)

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Make the POST request
	var result types.Presentation
	if err := s.httpClient.PostWithAuth(ctx, token, createURL, request, &result); err != nil {
		return nil, errors.NewAPIError(0, "presentation_create_failed",
			fmt.Sprintf("Failed to create presentation '%s'", request.Name), err.Error())
	}

	return &result, nil
}

// Update modifies an existing presentation by ID.
func (s *presentationService) Update(ctx context.Context, id int, request *types.PresentationCreateRequest) (*types.Presentation, error) {
	if id <= 0 {
		return nil, errors.NewValidationError("id", id, "presentation ID must be greater than 0")
	}

	if request == nil {
		return nil, errors.NewValidationError("request", request, "presentation update request cannot be nil")
	}

	if request.Name == "" {
		return nil, errors.NewValidationError("name", request.Name, "presentation name cannot be empty")
	}

	// Ensure we have authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return nil, err
	}

	// Build URL
	updateURL := fmt.Sprintf("%s/%s/Presentations/%d", s.config.BSNBaseURL, s.config.APIVersion, id)

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Make the PUT request
	var result types.Presentation
	if err := s.httpClient.PutWithAuth(ctx, token, updateURL, request, &result); err != nil {
		return nil, errors.NewAPIError(0, "presentation_update_failed",
			fmt.Sprintf("Failed to update presentation %d ('%s')", id, request.Name), err.Error())
	}

	return &result, nil
}

// DeleteByFilter removes presentations specified by a filter from the network.
func (s *presentationService) DeleteByFilter(ctx context.Context, filter string) (*types.PresentationDeleteResult, error) {
	// Ensure we have authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return nil, err
	}

	if filter == "" {
		return nil, errors.NewValidationError("filter", filter, "filter cannot be empty for bulk delete operation")
	}

	// Build URL with filter parameter
	params := url.Values{}
	params.Set("filter", filter)

	deleteURL := fmt.Sprintf("%s/%s/Presentations?%s", s.config.BSNBaseURL, s.config.APIVersion, params.Encode())

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Make the DELETE request
	var result types.PresentationDeleteResult
	if err := s.httpClient.DeleteWithAuth(ctx, token, deleteURL, &result); err != nil {
		return nil, errors.NewAPIError(0, "presentation_delete_failed",
			fmt.Sprintf("Failed to delete presentations with filter '%s'", filter), err.Error())
	}

	return &result, nil
}
