package services

import (
	"context"
	"fmt"
	"net/url"

	"github.com/brightdevelopers/gopurple/internal/auth"
	"github.com/brightdevelopers/gopurple/internal/config"
	"github.com/brightdevelopers/gopurple/internal/errors"
	"github.com/brightdevelopers/gopurple/internal/http"
	"github.com/brightdevelopers/gopurple/internal/types"
)

// BDeployService provides B-Deploy setup record operations.
type BDeployService interface {
	SetNetworkContext(ctx context.Context, networkName string) error
	GetSetupRecords(ctx context.Context, opts ...BDeployListOption) (*types.BDeployRecordList, error)
	GetSetupRecord(ctx context.Context, setupID string) (*types.BDeploySetupRecord, error)
	AddSetupRecord(ctx context.Context, record *types.BDeploySetupRecord) (*types.BDeployCreateResponse, error)
	UpdateSetupRecord(ctx context.Context, setupID string, record *types.BDeploySetupRecord) (*types.BDeploySetupRecord, error)
	DeleteSetupRecord(ctx context.Context, setupID string) (*types.BDeployDeleteResponse, error)
	GetDeviceBySerial(ctx context.Context, serial string) (*types.BDeployDeviceResponse, error)
	GetAllDevices(ctx context.Context, opts ...BDeployDeviceListOption) (*types.BDeployDeviceListResponse, error)
	CreateDevice(ctx context.Context, request *types.BDeployDeviceRequest) (string, error)
	UpdateDevice(ctx context.Context, deviceID string, request *types.BDeployDeviceRequest) (*types.BDeployDevice, error)
	DeleteDevice(ctx context.Context, deviceID string, serial string) error
}

// bDeployService implements the BDeployService interface.
type bDeployService struct {
	config         *config.Config
	httpClient     *http.HTTPClient
	authManager    *auth.AuthManager
	currentNetwork string // Track the current network context for device API calls
}

// NewBDeployService creates a new B-Deploy service.
func NewBDeployService(cfg *config.Config, httpClient *http.HTTPClient, authManager *auth.AuthManager) BDeployService {
	return &bDeployService{
		config:      cfg,
		httpClient:  httpClient,
		authManager: authManager,
	}
}

// SetNetworkContext sets the network context for B-Deploy operations.
func (s *bDeployService) SetNetworkContext(ctx context.Context, networkName string) error {
	if networkName == "" {
		return errors.NewValidationError("networkName", networkName, "network name cannot be empty")
	}

	// Ensure we have authentication
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return err
	}

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return err
	}

	// Build the BSN.cloud network context endpoint
	contextURL := "https://api.bsn.cloud/2022/06/REST/Self/Session/Network"

	// Build request body
	request := &types.NetworkContextRequest{
		Name: networkName,
	}

	// Make the API request - PUT to set network context
	var response interface{} // API returns empty body on success
	err = s.httpClient.PutWithAuth(ctx, token, contextURL, request, &response)
	if err != nil {
		return errors.NewAPIError(0, "network_context_failed",
			fmt.Sprintf("Failed to set network context to '%s'", networkName), err.Error())
	}

	// Store the network name for device API calls
	s.currentNetwork = networkName

	return nil
}

// GetSetupRecords retrieves B-Deploy setup records with optional filtering.
func (s *bDeployService) GetSetupRecords(ctx context.Context, opts ...BDeployListOption) (*types.BDeployRecordList, error) {
	// Ensure we have authentication
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Build query parameters from options
	config := &bDeployListConfig{}
	for _, opt := range opts {
		opt.apply(config)
	}

	params := url.Values{}
	if config.networkName != "" {
		params.Set("NetworkName", config.networkName)
	}
	if config.username != "" {
		params.Set("username", config.username)
	}
	if config.packageName != "" {
		params.Set("packageName", config.packageName)
	}
	// Don't add pagination parameters unless explicitly set
	// The B-Deploy API may not support pagination or may have issues with it

	// Build URL
	baseURL := "https://provision.bsn.cloud/rest-setup/v3/setup"
	if len(params) > 0 {
		baseURL += "?" + params.Encode()
	}

	// Make the API request - B-Deploy API returns a wrapper with error and result fields
	var apiResponse types.BDeployAPIResponse
	err = s.httpClient.GetWithAuth(ctx, token, baseURL, &apiResponse)
	if err != nil {
		return nil, errors.NewAPIError(0, "bdeploy_records_failed", "Failed to get B-Deploy setup records", err.Error())
	}

	// Check for API-level errors
	if apiResponse.Error != nil {
		return nil, errors.NewAPIError(0, "bdeploy_api_error", "B-Deploy API returned an error", fmt.Sprintf("%v", apiResponse.Error))
	}

	// Convert the response to our expected format
	setupRecords := &types.BDeployRecordList{
		Items:      apiResponse.Result,
		TotalCount: len(apiResponse.Result),
	}

	return setupRecords, nil
}

// GetSetupRecord retrieves a single B-Deploy setup record by ID.
func (s *bDeployService) GetSetupRecord(ctx context.Context, setupID string) (*types.BDeploySetupRecord, error) {
	if setupID == "" {
		return nil, errors.NewValidationError("setupID", setupID, "setup ID cannot be empty")
	}

	// Ensure we have authentication
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Build the B-Deploy setup retrieval endpoint
	// Using v3 API with query parameter format - returns array format with full setup records
	getURL := fmt.Sprintf("https://provision.bsn.cloud/rest-setup/v3/setup/?_id=%s", url.QueryEscape(setupID))

	// Make the API request - B-Deploy API returns array format with full setup record structure
	var apiResponse types.BDeployFullRecordAPIResponse
	err = s.httpClient.GetWithAuth(ctx, token, getURL, &apiResponse)
	if err != nil {
		return nil, errors.NewAPIError(0, "bdeploy_get_failed", "Failed to get B-Deploy setup record", err.Error())
	}

	// Check for API-level errors
	if apiResponse.Error != nil {
		return nil, errors.NewAPIError(0, "bdeploy_api_error", "B-Deploy API returned an error", fmt.Sprintf("%v", apiResponse.Error))
	}

	// Check if any records were returned
	if len(apiResponse.Result) == 0 {
		return nil, errors.NewAPIError(404, "bdeploy_not_found", "Setup record not found", fmt.Sprintf("No setup record found with ID: %s", setupID))
	}

	// Return the first (and should be only) record
	return &apiResponse.Result[0], nil
}

// AddSetupRecord creates a new B-Deploy setup record.
func (s *bDeployService) AddSetupRecord(ctx context.Context, record *types.BDeploySetupRecord) (*types.BDeployCreateResponse, error) {
	if record == nil {
		return nil, errors.NewValidationError("record", "nil", "setup record cannot be nil")
	}

	// Ensure we have authentication
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Build the B-Deploy setup creation endpoint
	createURL := "https://provision.bsn.cloud/rest-setup/v3/setup"

	// Make the API request - B-Deploy API returns wrapper format with full record in result
	var apiResponse types.BDeployCreateAPIResponse
	err = s.httpClient.PostWithAuth(ctx, token, createURL, record, &apiResponse)
	if err != nil {
		return nil, errors.NewAPIError(0, "bdeploy_create_failed", "Failed to create B-Deploy setup record", err.Error())
	}

	// Check for API-level errors
	if apiResponse.Error != nil {
		return nil, errors.NewAPIError(0, "bdeploy_api_error", "B-Deploy API returned an error", fmt.Sprintf("%v", apiResponse.Error))
	}

	// Convert to simplified response format with just the ID
	response := &types.BDeployCreateResponse{
		ID:      apiResponse.Result,
		Success: true,
	}

	return response, nil
}

// UpdateSetupRecord updates an existing B-Deploy setup record.
func (s *bDeployService) UpdateSetupRecord(ctx context.Context, setupID string, record *types.BDeploySetupRecord) (*types.BDeploySetupRecord, error) {
	if setupID == "" {
		return nil, errors.NewValidationError("setupID", setupID, "setup ID cannot be empty")
	}
	if record == nil {
		return nil, errors.NewValidationError("record", "nil", "setup record cannot be nil")
	}

	// Ensure we have authentication
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Ensure the record ID matches the setupID parameter
	record.ID = setupID

	// Build the B-Deploy setup update endpoint
	updateURL := "https://provision.bsn.cloud/rest-setup/v3/setup"

	// Make the API request - B-Deploy API returns wrapper format with full record in result
	var apiResponse types.BDeployUpdateAPIResponse
	err = s.httpClient.PutWithAuth(ctx, token, updateURL, record, &apiResponse)
	if err != nil {
		return nil, errors.NewAPIError(0, "bdeploy_update_failed", "Failed to update B-Deploy setup record", err.Error())
	}

	// Check for API-level errors
	if apiResponse.Error != nil {
		return nil, errors.NewAPIError(0, "bdeploy_api_error", "B-Deploy API returned an error", fmt.Sprintf("%v", apiResponse.Error))
	}

	// Return the updated setup record
	return apiResponse.Result, nil
}

// DeleteSetupRecord deletes a B-Deploy setup record by ID.
func (s *bDeployService) DeleteSetupRecord(ctx context.Context, setupID string) (*types.BDeployDeleteResponse, error) {
	if setupID == "" {
		return nil, errors.NewValidationError("setupID", setupID, "setup ID cannot be empty")
	}

	// Ensure we have authentication
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Build the B-Deploy setup deletion endpoint
	// Note: Using v3 API with query parameter format (v3 path parameter format doesn't work)
	deleteURL := fmt.Sprintf("https://provision.bsn.cloud/rest-setup/v3/setup/?_id=%s", url.QueryEscape(setupID))

	// Make the API request
	var response types.BDeployDeleteResponse
	err = s.httpClient.DeleteWithAuth(ctx, token, deleteURL, &response)
	if err != nil {
		return nil, errors.NewAPIError(0, "bdeploy_delete_failed", "Failed to delete B-Deploy setup record", err.Error())
	}

	return &response, nil
}

// GetDeviceBySerial retrieves a B-Deploy device setup record by serial number.
func (s *bDeployService) GetDeviceBySerial(ctx context.Context, serial string) (*types.BDeployDeviceResponse, error) {
	if serial == "" {
		return nil, errors.NewValidationError("serial", serial, "serial number cannot be empty")
	}

	// Ensure we have authentication
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Build the B-Deploy device endpoint with serial query parameter
	// Include NetworkName if we have a current network set (required for full device record with setupId)
	params := url.Values{}
	params.Set("serial", serial)
	if s.currentNetwork != "" {
		params.Set("NetworkName", s.currentNetwork)
	}
	deviceURL := fmt.Sprintf("https://provision.bsn.cloud/rest-device/v2/device/?%s", params.Encode())

	// Try wrapped response format first (like GetAllDevices does)
	var wrappedResponse types.BDeployDeviceResponse
	err = s.httpClient.GetWithAuth(ctx, token, deviceURL, &wrappedResponse)
	if err == nil && wrappedResponse.Result.Players != nil && len(wrappedResponse.Result.Players) > 0 {
		return &wrappedResponse, nil
	}

	// If that fails or returns empty Players, try direct response format
	var directResult types.BDeployDeviceResult
	err = s.httpClient.GetWithAuth(ctx, token, deviceURL, &directResult)
	if err != nil {
		return nil, err
	}

	// Wrap the direct response in the expected format
	return &types.BDeployDeviceResponse{
		Error:  nil,
		Result: directResult,
	}, nil
}

// GetAllDevices retrieves all B-Deploy devices on the network.
// This method uses the network context set via SetNetworkContext.
// The network context must be set before calling this method, or it may return 0 devices.
func (s *bDeployService) GetAllDevices(ctx context.Context, opts ...BDeployDeviceListOption) (*types.BDeployDeviceListResponse, error) {
	// Ensure we have authentication
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Build query parameters from options
	config := &bDeployDeviceListConfig{}
	for _, opt := range opts {
		opt.apply(config)
	}

	params := url.Values{}

	// Add NetworkName query parameter if we have a current network set
	if s.currentNetwork != "" {
		params.Set("NetworkName", s.currentNetwork)
	}

	// Add setupName filter if specified (using query[setupName] format)
	if config.setupName != "" {
		params.Set("query[setupName]", config.setupName)
	}

	// Build the B-Deploy device list endpoint with query parameters
	// The Device API requires an explicit NetworkName query parameter to filter by network,
	// in addition to the session network context.
	deviceListURL := "https://provision.bsn.cloud/rest-device/v2/device/"
	if len(params) > 0 {
		deviceListURL += "?" + params.Encode()
	}

	// First try the wrapped response format (like other B-Deploy APIs)
	var wrappedResponse types.BDeployDeviceListAPIResponse
	err = s.httpClient.GetWithAuth(ctx, token, deviceListURL, &wrappedResponse)
	if err == nil && wrappedResponse.Result != nil {
		return wrappedResponse.Result, nil
	}

	// If that fails or returns nil result, try direct response format
	var directResponse types.BDeployDeviceListResponse
	err = s.httpClient.GetWithAuth(ctx, token, deviceListURL, &directResponse)
	if err != nil {
		return nil, err
	}

	return &directResponse, nil
}

// BDeployListOption represents an option for B-Deploy record listing.
type BDeployListOption interface {
	apply(*bDeployListConfig)
}

// bDeployListConfig holds configuration for B-Deploy record listing.
type bDeployListConfig struct {
	networkName string
	username    string
	packageName string
	pageSize    int
	page        int
}

// optionFunc is a function that implements BDeployListOption.
type bDeployOptionFunc func(*bDeployListConfig)

func (f bDeployOptionFunc) apply(c *bDeployListConfig) {
	f(c)
}

// WithNetworkName sets the network name filter for B-Deploy record listing.
func WithNetworkName(networkName string) BDeployListOption {
	return bDeployOptionFunc(func(c *bDeployListConfig) {
		c.networkName = networkName
	})
}

// WithUsername sets the username filter for B-Deploy record listing.
func WithUsername(username string) BDeployListOption {
	return bDeployOptionFunc(func(c *bDeployListConfig) {
		c.username = username
	})
}

// WithPackageName sets the package name filter for B-Deploy record listing.
func WithPackageName(packageName string) BDeployListOption {
	return bDeployOptionFunc(func(c *bDeployListConfig) {
		c.packageName = packageName
	})
}

// WithBDeployPageSize sets the page size for B-Deploy record listing.
func WithBDeployPageSize(pageSize int) BDeployListOption {
	return bDeployOptionFunc(func(c *bDeployListConfig) {
		c.pageSize = pageSize
	})
}

// WithBDeployPage sets the page number for B-Deploy record listing.
func WithBDeployPage(page int) BDeployListOption {
	return bDeployOptionFunc(func(c *bDeployListConfig) {
		c.page = page
	})
}

// BDeployDeviceListOption represents an option for B-Deploy device listing.
type BDeployDeviceListOption interface {
	apply(*bDeployDeviceListConfig)
}

// bDeployDeviceListConfig holds configuration for B-Deploy device listing.
type bDeployDeviceListConfig struct {
	setupName string
}

// bDeployDeviceOptionFunc is a function that implements BDeployDeviceListOption.
type bDeployDeviceOptionFunc func(*bDeployDeviceListConfig)

func (f bDeployDeviceOptionFunc) apply(c *bDeployDeviceListConfig) {
	f(c)
}

// WithSetupName sets the setup name filter for B-Deploy device listing.
// This uses the query[setupName] API parameter to filter devices by setup name.
func WithSetupName(setupName string) BDeployDeviceListOption {
	return bDeployDeviceOptionFunc(func(c *bDeployDeviceListConfig) {
		c.setupName = setupName
	})
}

// CreateDevice creates a new B-Deploy device record.
// This registers a device serial number with the B-Deploy system.
func (s *bDeployService) CreateDevice(ctx context.Context, request *types.BDeployDeviceRequest) (string, error) {
	if request.Serial == "" {
		return "", errors.NewValidationError("serial", request.Serial, "serial number cannot be empty")
	}
	if request.NetworkName == "" {
		return "", errors.NewValidationError("NetworkName", request.NetworkName, "network name cannot be empty")
	}
	if request.Username == "" {
		return "", errors.NewValidationError("username", request.Username, "username cannot be empty")
	}

	// Ensure we have authentication
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return "", err
	}

	token, err := s.authManager.GetToken()
	if err != nil {
		return "", errors.NewAuthError("failed to get access token", err)
	}

	// POST to /rest-device/v2/device/
	createURL := "https://provision.bsn.cloud/rest-device/v2/device/"

	var response types.BDeployDeviceCreateResponse
	err = s.httpClient.PostWithAuth(ctx, token, createURL, request, &response)
	if err != nil {
		return "", fmt.Errorf("failed to create device: %w", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("API error creating device: %v", response.Error)
	}

	return response.Result, nil
}

// UpdateDevice updates an existing B-Deploy device record.
// This is used to associate a device with a setup ID.
func (s *bDeployService) UpdateDevice(ctx context.Context, deviceID string, request *types.BDeployDeviceRequest) (*types.BDeployDevice, error) {
	if deviceID == "" {
		return nil, errors.NewValidationError("deviceID", deviceID, "device ID cannot be empty")
	}
	if request.Serial == "" {
		return nil, errors.NewValidationError("serial", request.Serial, "serial number cannot be empty")
	}

	// Ensure we have authentication
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, errors.NewAuthError("failed to get access token", err)
	}

	// Ensure the device ID is set in the request
	request.ID = deviceID

	// PUT to /rest-device/v2/device?_id={deviceID}
	updateURL := fmt.Sprintf("https://provision.bsn.cloud/rest-device/v2/device?_id=%s", url.QueryEscape(deviceID))

	var response types.BDeployDeviceUpdateResponse
	err = s.httpClient.PutWithAuth(ctx, token, updateURL, request, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("API error updating device: %v", response.Error)
	}

	// The API may return success without a result body
	// In this case, fetch the device to verify the update
	if response.Result == nil {
		// Fetch the updated device by serial (including NetworkName for full record)
		deviceResponse, err := s.GetDeviceBySerial(ctx, request.Serial)
		if err != nil {
			return nil, fmt.Errorf("device updated but failed to fetch updated record: %w", err)
		}

		if deviceResponse.Result.Matched > 0 && len(deviceResponse.Result.Players) > 0 {
			device := &deviceResponse.Result.Players[0]
			// If setupId is still empty after fetch, use the one we sent in the request
			// This handles cases where the API doesn't return setupId even though it was set
			if device.SetupID == "" && request.SetupID != "" {
				device.SetupID = request.SetupID
			}
			return device, nil
		}

		// If still not found, return a constructed device from the request
		return &types.BDeployDevice{
			ID:          deviceID,
			Serial:      request.Serial,
			Name:        request.Name,
			NetworkName: request.NetworkName,
			Desc:        request.Desc,
			SetupID:     request.SetupID,
			Username:    request.Username,
		}, nil
	}

	// API returned a result - check if setupId is populated
	// If not, fill it from our request since we know we just set it
	if response.Result.SetupID == "" && request.SetupID != "" {
		response.Result.SetupID = request.SetupID
	}

	return response.Result, nil
}

// DeleteDevice removes a device from the B-Deploy system.
// Either deviceID or serial must be provided. If both are provided, deviceID takes precedence.
func (s *bDeployService) DeleteDevice(ctx context.Context, deviceID string, serial string) error {
	// Validate that at least one identifier is provided
	if deviceID == "" && serial == "" {
		return errors.NewValidationError("deviceID/serial", "", "either device ID or serial number must be provided")
	}

	// Ensure we have authentication
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return err
	}

	token, err := s.authManager.GetToken()
	if err != nil {
		return errors.NewAuthError("failed to get access token", err)
	}

	// Build delete URL with either _id or serial parameter
	var deleteURL string
	if deviceID != "" {
		deleteURL = fmt.Sprintf("https://provision.bsn.cloud/rest-device/v2/device?_id=%s", url.QueryEscape(deviceID))
	} else {
		deleteURL = fmt.Sprintf("https://provision.bsn.cloud/rest-device/v2/device?serial=%s", url.QueryEscape(serial))
	}

	// DELETE returns a simple response, we'll use a generic struct
	var response struct {
		Message string `json:"message,omitempty"`
		Error   string `json:"error,omitempty"`
	}

	// B-Deploy device DELETE endpoint requires Content-Type header even without body
	err = s.httpClient.DoWithAuth(ctx, token, &http.Request{
		Method: "DELETE",
		URL:    deleteURL,
		Result: &response,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}

	if response.Error != "" {
		return fmt.Errorf("API error deleting device: %s", response.Error)
	}

	return nil
}
