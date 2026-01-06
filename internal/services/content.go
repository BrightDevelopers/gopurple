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

// ContentService provides content file management operations.
type ContentService interface {
	// List retrieves a list of content files on the network
	List(ctx context.Context, opts ...ListOption) (*types.ContentFileList, error)

	// Delete removes content files specified by a filter from the network
	Delete(ctx context.Context, filter string) (*types.ContentDeleteResult, error)

	// DeleteByID removes a specific content file by ID
	DeleteByID(ctx context.Context, id int) error

	// GetCount retrieves the number of content files on the network
	GetCount(ctx context.Context) (*types.ContentFileCount, error)

	// GetByID retrieves the specified content file metadata
	GetByID(ctx context.Context, id int) (*types.ContentFile, error)

	// Download downloads the actual content file data by ID
	Download(ctx context.Context, id int) ([]byte, error)
}

// contentService implements the ContentService interface.
type contentService struct {
	config      *config.Config
	httpClient  *http.HTTPClient
	authManager *auth.AuthManager
}

// NewContentService creates a new content service.
func NewContentService(cfg *config.Config, httpClient *http.HTTPClient, authManager *auth.AuthManager) ContentService {
	return &contentService{
		config:      cfg,
		httpClient:  httpClient,
		authManager: authManager,
	}
}

// List retrieves a list of content files on the network with optional filtering and pagination.
func (s *contentService) List(ctx context.Context, opts ...ListOption) (*types.ContentFileList, error) {
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

	// Build URL - Content API base URL
	baseURL := fmt.Sprintf("%s/%s/Content", s.config.BSNBaseURL, s.config.APIVersion)
	if len(params) > 0 {
		baseURL += "?" + params.Encode()
	}

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Make the request
	var result types.ContentFileList
	if err := s.httpClient.GetWithAuth(ctx, token, baseURL, &result); err != nil {
		return nil, errors.NewAPIError(0, "content_list_failed", "Failed to list content files", err.Error())
	}

	return &result, nil
}

// Delete removes content files specified by a filter from the network.
func (s *contentService) Delete(ctx context.Context, filter string) (*types.ContentDeleteResult, error) {
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

	deleteURL := fmt.Sprintf("%s/%s/Content?%s", s.config.BSNBaseURL, s.config.APIVersion, params.Encode())

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Make the DELETE request
	var result types.ContentDeleteResult
	if err := s.httpClient.DeleteWithAuth(ctx, token, deleteURL, &result); err != nil {
		return nil, errors.NewAPIError(0, "content_delete_failed",
			fmt.Sprintf("Failed to delete content files with filter '%s'", filter), err.Error())
	}

	return &result, nil
}

// DeleteByID removes a specific content file by ID.
func (s *contentService) DeleteByID(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.NewValidationError("id", strconv.Itoa(id), "content ID must be positive")
	}

	// Ensure we have authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return err
	}

	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return err
	}

	// Build URL
	deleteURL := fmt.Sprintf("%s/%s/Content/%d", s.config.BSNBaseURL, s.config.APIVersion, id)

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return err
	}

	// Make the DELETE request
	if err := s.httpClient.DeleteWithAuth(ctx, token, deleteURL, nil); err != nil {
		return errors.NewAPIError(0, "content_delete_failed",
			fmt.Sprintf("Failed to delete content with ID %d", id), err.Error())
	}

	return nil
}

// GetCount retrieves the number of content files on the network.
func (s *contentService) GetCount(ctx context.Context) (*types.ContentFileCount, error) {
	// Ensure we have authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return nil, err
	}

	// Build URL
	countURL := fmt.Sprintf("%s/%s/Content/Count", s.config.BSNBaseURL, s.config.APIVersion)

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Make the request
	var result types.ContentFileCount
	if err := s.httpClient.GetWithAuth(ctx, token, countURL, &result); err != nil {
		return nil, errors.NewAPIError(0, "content_count_failed", "Failed to get content file count", err.Error())
	}

	return &result, nil
}

// GetByID retrieves metadata for the specified content file.
func (s *contentService) GetByID(ctx context.Context, id int) (*types.ContentFile, error) {
	// Ensure we have authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return nil, err
	}

	if id <= 0 {
		return nil, errors.NewValidationError("id", id, "content ID must be greater than 0")
	}

	// Build URL
	getURL := fmt.Sprintf("%s/%s/Content/%d", s.config.BSNBaseURL, s.config.APIVersion, id)

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Make the request
	var result types.ContentFile
	if err := s.httpClient.GetWithAuth(ctx, token, getURL, &result); err != nil {
		return nil, errors.NewAPIError(0, "content_get_failed",
			fmt.Sprintf("Failed to get content file %d", id), err.Error())
	}

	return &result, nil
}

// Download downloads the actual content file data by ID.
func (s *contentService) Download(ctx context.Context, id int) ([]byte, error) {
	// Ensure we have authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return nil, err
	}

	if id <= 0 {
		return nil, errors.NewValidationError("id", id, "content ID must be greater than 0")
	}

	// Build URL for download
	downloadURL := fmt.Sprintf("%s/%s/Content/%d/Download", s.config.BSNBaseURL, s.config.APIVersion, id)

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Make the request - use GetBytesWithAuth to download binary data
	data, err := s.httpClient.GetBytesWithAuth(ctx, token, downloadURL)
	if err != nil {
		return nil, errors.NewAPIError(0, "content_download_failed",
			fmt.Sprintf("Failed to download content file %d", id), err.Error())
	}

	return data, nil
}
