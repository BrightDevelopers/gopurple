// Package gopurple provides a comprehensive SDK for interacting with BrightSign Cloud (BSN.cloud)
// and Remote Diagnostic Web Service (RDWS) APIs.
//
// The SDK enables developers to build custom automation tools and applications for BrightSign
// device management with a clean, extensible API that abstracts away implementation details.
//
// Basic usage:
//
//	client, err := gopurple.New(
//	    gopurple.WithCredentials(clientID, clientSecret),
//	    gopurple.WithNetwork("My Network"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	ctx := context.Background()
//
//	// Authenticate
//	if err := client.Authenticate(ctx); err != nil {
//	    log.Fatal(err)
//	}
//
//	// List networks
//	networks, err := client.GetNetworks(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("Found %d networks\n", len(networks))
package gopurple

import (
	"context"

	"github.com/brightdevelopers/gopurple/internal/auth"
	"github.com/brightdevelopers/gopurple/internal/config"
	"github.com/brightdevelopers/gopurple/internal/errors"
	"github.com/brightdevelopers/gopurple/internal/http"
	"github.com/brightdevelopers/gopurple/internal/services"
	"github.com/brightdevelopers/gopurple/internal/types"
)

// Re-export public types from internal packages
type (
	// Network represents a BSN.cloud network.
	Network = types.Network

	// NetworkSubscription represents subscription details for a network.
	NetworkSubscription = types.NetworkSubscription

	// NetworkSettings represents network configuration settings.
	NetworkSettings = types.NetworkSettings

	// Device represents a BrightSign device in the network.
	Device = types.Device

	// DeviceSettings represents device configuration settings.
	DeviceSettings = types.DeviceSettings

	// Group represents a device group.
	Group = types.Group

	// GroupList represents a list of groups.
	GroupList = types.GroupList

	// DeviceList represents a paginated list of devices.
	DeviceList = types.DeviceList

	// Region represents a rectangular region for screenshots.
	Region = types.Region

	// DeviceStatus represents the current operational status of a device.
	DeviceStatus = types.DeviceStatus

	// DeviceError represents an error log entry from a device.
	DeviceError = types.DeviceError

	// DeviceErrorList represents a paginated list of device errors.
	DeviceErrorList = types.DeviceErrorList

	// RebootResponse represents the response from a device reboot request.
	RebootResponse = types.RebootResponse

	// RebootType represents the type of reboot to perform.
	RebootType = types.RebootType

	// ReprovisionResponse represents the response from a device re-provision request.
	ReprovisionResponse = types.ReprovisionResponse

	// DWSPasswordRequest represents a request to set the DWS password.
	DWSPasswordRequest = types.DWSPasswordRequest

	// DWSPasswordInfo represents information about the DWS password.
	DWSPasswordInfo = types.DWSPasswordInfo

	// DWSPasswordGetResponse represents the response from getting DWS password info.
	DWSPasswordGetResponse = types.DWSPasswordGetResponse

	// DWSPasswordSetResponse represents the response from setting DWS password.
	DWSPasswordSetResponse = types.DWSPasswordSetResponse

	// BDeployRecord represents a B-Deploy setup record.
	BDeployRecord = types.BDeployRecord

	// BDeployRecordList represents a list of B-Deploy setup records.
	BDeployRecordList = types.BDeployRecordList

	// BDeploySetupRecord represents a complete B-Deploy setup record for creation.
	BDeploySetupRecord = types.BDeploySetupRecord

	// BDeployCreateResponse represents the response from creating a B-Deploy record.
	BDeployCreateResponse = types.BDeployCreateResponse

	// BDeployDeleteResponse represents the response from deleting a B-Deploy record.
	BDeployDeleteResponse = types.BDeployDeleteResponse

	// BDeployInfo represents the B-Deploy section of a setup record.
	BDeployInfo = types.BDeployInfo

	// BSNTokenEntity represents the BSN device registration token.
	BSNTokenEntity = types.BSNTokenEntity

	// NetworkInterface represents a network interface configuration.
	NetworkInterface = types.NetworkInterface

	// NetworkConfig represents network configuration.
	NetworkConfig = types.NetworkConfig

	// NetworkContextRequest represents a request to set the network context.
	NetworkContextRequest = types.NetworkContextRequest

	// BDeployDevice represents a device in the B-Deploy system.
	BDeployDevice = types.BDeployDevice

	// BDeployDeviceRequest represents a request to create or update a B-Deploy device.
	BDeployDeviceRequest = types.BDeployDeviceRequest

	// BDeployDeviceResponse represents the response from B-Deploy device API.
	BDeployDeviceResponse = types.BDeployDeviceResponse

	// BDeployDeviceListResponse represents a list of B-Deploy devices.
	BDeployDeviceListResponse = types.BDeployDeviceListResponse

	// RDWSInfo represents player information from the rDWS /info/ endpoint
	RDWSInfo = types.RDWSInfo

	// RDWSNetworkInterface represents a network interface from the player
	RDWSNetworkInterface = types.RDWSNetworkInterface

	// RDWSIPAddress represents an IP address configuration
	RDWSIPAddress = types.RDWSIPAddress

	// RDWSTimeInfo represents time information from the rDWS /time/ endpoint
	RDWSTimeInfo = types.RDWSTimeInfo

	// RDWSTimeSetRequest represents the request body for rDWS PUT /time/
	RDWSTimeSetRequest = types.RDWSTimeSetRequest

	// RDWSHealthInfo represents health status from the rDWS /health/ endpoint
	RDWSHealthInfo = types.RDWSHealthInfo

	// RDWSFileStat represents file statistics from the fs module
	RDWSFileStat = types.RDWSFileStat

	// RDWSFileInfo represents a file or directory entry
	RDWSFileInfo = types.RDWSFileInfo

	// RDWSStorageStats represents storage device statistics
	RDWSStorageStats = types.RDWSStorageStats

	// RDWSStorageInfo represents storage device information
	RDWSStorageInfo = types.RDWSStorageInfo

	// RDWSFileListResult represents the result of listing files
	RDWSFileListResult = types.RDWSFileListResult

	// RDWSFileListResponse represents the response from listing files
	RDWSFileListResponse = types.RDWSFileListResponse

	// RDWSFileUploadItem represents a single file to upload
	RDWSFileUploadItem = types.RDWSFileUploadItem

	// RDWSFileUploadRequest represents the request to upload files
	RDWSFileUploadRequest = types.RDWSFileUploadRequest

	// RDWSFileUploadResult represents the result of file upload
	RDWSFileUploadResult = types.RDWSFileUploadResult

	// RDWSFileUploadResponse represents the response from file upload
	RDWSFileUploadResponse = types.RDWSFileUploadResponse

	// RDWSFileRenameRequest represents the request to rename a file
	RDWSFileRenameRequest = types.RDWSFileRenameRequest

	// RDWSFileOperationResponse represents a generic success/error response for file operations
	RDWSFileOperationResponse = types.RDWSFileOperationResponse

	// DeviceDownload represents a content download on a device
	DeviceDownload = types.DeviceDownload

	// DeviceDownloadList represents a list of device downloads
	DeviceDownloadList = types.DeviceDownloadList

	// DeviceOperation represents an operation performed on a device
	DeviceOperation = types.DeviceOperation

	// DeviceOperationList represents a list of device operations
	DeviceOperationList = types.DeviceOperationList

	// RDWSLocalDWSInfo represents the current state of local DWS
	RDWSLocalDWSInfo = types.RDWSLocalDWSInfo

	// RDWSDiagnosticsInfo represents network diagnostics information
	RDWSDiagnosticsInfo = types.RDWSDiagnosticsInfo

	// RDWSDNSLookupResult represents DNS lookup result
	RDWSDNSLookupResult = types.RDWSDNSLookupResult

	// RDWSPingResult represents ping diagnostic result
	RDWSPingResult = types.RDWSPingResult

	// RDWSTraceRouteHop represents a single hop in trace route
	RDWSTraceRouteHop = types.RDWSTraceRouteHop

	// RDWSTraceRouteResult represents trace route diagnostic result
	RDWSTraceRouteResult = types.RDWSTraceRouteResult

	// RDWSNetworkConfig represents network interface configuration
	RDWSNetworkConfig = types.RDWSNetworkConfig

	// RDWSNetworkConfigSetRequest represents the request to set network configuration
	RDWSNetworkConfigSetRequest = types.RDWSNetworkConfigSetRequest

	// RDWSNetworkNeighbor represents a device on the network
	RDWSNetworkNeighbor = types.RDWSNetworkNeighbor

	// RDWSNetworkNeighborhoodResult represents network neighborhood information
	RDWSNetworkNeighborhoodResult = types.RDWSNetworkNeighborhoodResult

	// RDWSPacketCaptureStatus represents packet capture status
	RDWSPacketCaptureStatus = types.RDWSPacketCaptureStatus

	// RDWSPacketCaptureStartRequest represents the request to start packet capture
	RDWSPacketCaptureStartRequest = types.RDWSPacketCaptureStartRequest

	// RDWSTelnetInfo represents telnet status information
	RDWSTelnetInfo = types.RDWSTelnetInfo

	// RDWSSSHInfo represents SSH status information
	RDWSSSHInfo = types.RDWSSSHInfo

	// Subscription represents a device subscription in BSN.cloud
	Subscription = types.Subscription

	// SubscriptionList represents a list of subscriptions with pagination support
	SubscriptionList = types.SubscriptionList

	// SubscriptionCount represents the count of subscriptions
	SubscriptionCount = types.SubscriptionCount

	// SubscriptionOperation represents an operation that can be performed on subscriptions
	SubscriptionOperation = types.SubscriptionOperation

	// SubscriptionOperations represents the list of available operations
	SubscriptionOperations = types.SubscriptionOperations

	// RDWSLogFile represents a single log file from the player
	RDWSLogFile = types.RDWSLogFile

	// RDWSLogs represents the collection of log files from a player
	RDWSLogs = types.RDWSLogs

	// RDWSCrashDumpFile represents a single crash dump file from the player
	RDWSCrashDumpFile = types.RDWSCrashDumpFile

	// RDWSCrashDump represents the collection of crash dump files from a player
	RDWSCrashDump = types.RDWSCrashDump

	// RDWSRegistry represents the full player registry
	RDWSRegistry = types.RDWSRegistry

	// RDWSRegistryValue represents a single registry value
	RDWSRegistryValue = types.RDWSRegistryValue

	// RDWSRecoveryURL represents the player's recovery URL setting
	RDWSRecoveryURL = types.RDWSRecoveryURL

	// ContentFile represents a content file in BSN.cloud
	ContentFile = types.ContentFile

	// ContentFileList represents a paginated list of content files
	ContentFileList = types.ContentFileList

	// ContentFileCount represents the count of content files
	ContentFileCount = types.ContentFileCount

	// ContentDeleteResult represents the result of deleting content files
	ContentDeleteResult = types.ContentDeleteResult

	// Presentation represents a presentation in BSN.cloud
	Presentation = types.Presentation

	// PresentationCount represents the count of presentations
	PresentationCount = types.PresentationCount

	// PresentationList represents a paginated list of presentations
	PresentationList = types.PresentationList

	// ScheduleSettings represents scheduling configuration for presentations
	ScheduleSettings = types.ScheduleSettings

	// PresentationCreateRequest represents a request to create a new presentation
	PresentationCreateRequest = types.PresentationCreateRequest

	// PresentationAutorun represents autorun configuration for a presentation
	PresentationAutorun = types.PresentationAutorun

	// PresentationDeviceWebPage represents device web page configuration
	PresentationDeviceWebPage = types.PresentationDeviceWebPage

	// PresentationScreenSettings represents screen configuration for a presentation
	PresentationScreenSettings = types.PresentationScreenSettings

	// PresentationDeleteResult represents the result of deleting presentations by filter
	PresentationDeleteResult = types.PresentationDeleteResult

	// ScheduledPresentation represents a scheduled presentation in a group
	ScheduledPresentation = types.ScheduledPresentation

	// ScheduledPresentationList represents a paginated list of scheduled presentations
	ScheduledPresentationList = types.ScheduledPresentationList

	// DeviceWebPage represents a device web page template
	DeviceWebPage = types.DeviceWebPage

	// DeviceWebPageList represents a paginated list of device web pages
	DeviceWebPageList = types.DeviceWebPageList

	// UploadRequest represents a file upload request (deprecated)
	UploadRequest = types.UploadRequest

	// UploadResponse represents the response from a file upload
	UploadResponse = types.UploadResponse

	// ContentUploadArguments represents arguments for creating a content upload
	ContentUploadArguments = types.ContentUploadArguments

	// UploadSessionResponse represents the response from creating an upload session
	UploadSessionResponse = types.UploadSessionResponse

	// ContentUploadStatus represents the status of a content upload
	ContentUploadStatus = types.ContentUploadStatus

	// ChunkUploadResponse represents the response from uploading a chunk
	ChunkUploadResponse = types.ChunkUploadResponse
)

// Re-export configuration options
type Option = config.Option

var (
	// WithCredentials sets the BSN.cloud OAuth2 credentials.
	WithCredentials = config.WithCredentials

	// WithNetwork sets the default network name for device operations.
	WithNetwork = config.WithNetwork

	// WithTimeout sets the HTTP request timeout for API calls.
	WithTimeout = config.WithTimeout

	// WithRetryCount sets the number of retry attempts for failed API requests.
	WithRetryCount = config.WithRetryCount

	// WithDebug enables debug logging of all HTTP requests and responses.
	WithDebug = config.WithDebug

	// WithEndpoints sets custom API endpoints for BSN.cloud and RDWS.
	WithEndpoints = config.WithEndpoints

	// WithTokenEndpoint sets a custom OAuth2 token endpoint.
	WithTokenEndpoint = config.WithTokenEndpoint

	// WithOIDCURL sets the OIDC base URL and derives the token endpoint from it.
	WithOIDCURL = config.WithOIDCURL

	// WithDeviceSerial sets a default device serial number for single-device operations.
	WithDeviceSerial = config.WithDeviceSerial
)

// Re-export device list options
type ListOption = services.ListOption
type BDeployListOption = services.BDeployListOption
type BDeployDeviceListOption = services.BDeployDeviceListOption

var (
	// WithPageSize sets the page size for device listing.
	WithPageSize = services.WithPageSize

	// WithMarker sets the pagination marker for device listing.
	WithMarker = services.WithMarker

	// WithFilter sets the filter expression for device listing.
	WithFilter = services.WithFilter

	// WithSort sets the sort expression for device listing.
	WithSort = services.WithSort

	// WithNetworkName sets the network name filter for B-Deploy record listing.
	WithNetworkName = services.WithNetworkName

	// WithUsername sets the username filter for B-Deploy record listing.
	WithUsername = services.WithUsername

	// WithPackageName sets the package name filter for B-Deploy record listing.
	WithPackageName = services.WithPackageName

	// WithBDeployPageSize sets the page size for B-Deploy record listing.
	WithBDeployPageSize = services.WithBDeployPageSize

	// WithBDeployPage sets the page number for B-Deploy record listing.
	WithBDeployPage = services.WithBDeployPage

	// WithSetupName sets the setup name filter for B-Deploy device listing.
	// This filters devices by their associated setup package name using the query[setupName] API parameter.
	WithSetupName = services.WithSetupName
)

// Re-export reboot type constants
const (
	// RebootTypeNormal performs a standard reboot
	RebootTypeNormal = types.RebootTypeNormal
	// RebootTypeCrash performs a reboot and saves a crash report
	RebootTypeCrash = types.RebootTypeCrash
	// RebootTypeFactoryReset performs a factory reset reboot
	RebootTypeFactoryReset = types.RebootTypeFactoryReset
	// RebootTypeDisableAutorun disables autorun script and reboots
	RebootTypeDisableAutorun = types.RebootTypeDisableAutorun
)

// Re-export error checking functions
var (
	// IsAuthenticationError checks if an error is authentication-related.
	IsAuthenticationError = errors.IsAuthenticationError

	// IsNetworkError checks if an error is network-related.
	IsNetworkError = errors.IsNetworkError

	// IsConfigurationError checks if an error is configuration-related.
	IsConfigurationError = errors.IsConfigurationError

	// IsRetryableError checks if an error might succeed on retry.
	IsRetryableError = errors.IsRetryableError
)

// Client is the main SDK client that provides access to all BrightSign services.
//
// It manages authentication, network selection, and provides access to various
// service interfaces for device management, firmware operations, and job orchestration.
type Client struct {
	config      *config.Config
	httpClient  *http.HTTPClient
	authManager *auth.AuthManager

	// Services
	Devices        services.DeviceService
	BDeploy        services.BDeployService
	Provisioning   services.ProvisioningService
	RDWS           services.RDWSService
	Subscriptions  services.SubscriptionService
	Content        services.ContentService
	Upload         services.UploadService
	Presentations  services.PresentationService
	Schedules      services.ScheduleService
	DeviceWebPages services.DeviceWebPageService
}

// New creates a new BrightSign SDK client with the given configuration options.
//
// The client must be configured with at least BSN.cloud credentials:
//
//	client, err := gopurple.New(
//	    gopurple.WithCredentials(clientID, clientSecret),
//	)
//
// Additional options can be provided for custom endpoints, timeouts, and network selection:
//
//	client, err := gopurple.New(
//	    gopurple.WithCredentials(clientID, clientSecret),
//	    gopurple.WithNetwork("Production Network"),
//	    gopurple.WithTimeout(60*time.Second),
//	)
//
// By default, the client will load credentials from environment variables:
//   - BS_CLIENT_ID: BSN.cloud API client ID
//   - BS_SECRET: BSN.cloud API client secret
//   - BS_NETWORK: BSN.cloud network name (optional)
func New(opts ...Option) (*Client, error) {
	cfg := config.DefaultConfig()

	// Load configuration from environment variables
	cfg.LoadFromEnv()

	// Apply any provided options (these override environment variables)
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Create HTTP client
	httpClient := http.NewHTTPClient(cfg)

	// Create authentication manager
	authManager := auth.NewAuthManager(cfg, httpClient)

	client := &Client{
		config:         cfg,
		httpClient:     httpClient,
		authManager:    authManager,
		Devices:        services.NewDeviceService(cfg, httpClient, authManager),
		BDeploy:        services.NewBDeployService(cfg, httpClient, authManager),
		Provisioning:   services.NewProvisioningService(cfg, httpClient, authManager),
		RDWS:           services.NewRDWSService(cfg, httpClient, authManager),
		Subscriptions:  services.NewSubscriptionService(cfg, httpClient, authManager),
		Content:        services.NewContentService(cfg, httpClient, authManager),
		Upload:         services.NewUploadService(cfg, httpClient, authManager),
		Presentations:  services.NewPresentationService(cfg, httpClient, authManager),
		Schedules:      services.NewScheduleService(cfg, httpClient, authManager),
		DeviceWebPages: services.NewDeviceWebPageService(cfg, httpClient, authManager),
	}

	return client, nil
}

// Authenticate performs OAuth2 authentication with BSN.cloud using the configured credentials.
//
// This method must be called before making any API requests. The SDK will automatically
// handle token refresh for subsequent requests.
func (c *Client) Authenticate(ctx context.Context) error {
	return c.authManager.Authenticate(ctx)
}

// SetNetwork sets the active network for subsequent device operations by name.
//
// Many device operations require a network context. This method sets the default
// network for the client session.
func (c *Client) SetNetwork(ctx context.Context, networkName string) error {
	return c.authManager.SetNetwork(ctx, networkName)
}

// SetNetworkByID sets the active network for subsequent device operations by ID.
func (c *Client) SetNetworkByID(ctx context.Context, networkID int) error {
	return c.authManager.SetNetworkByID(ctx, networkID)
}

// GetNetworks retrieves all networks accessible to the authenticated user.
//
// This method requires authentication but does not require a network context.
func (c *Client) GetNetworks(ctx context.Context) ([]Network, error) {
	return c.authManager.GetNetworks(ctx)
}

// GetCurrentNetwork returns the currently selected network.
//
// Returns an error if no network has been selected.
func (c *Client) GetCurrentNetwork(ctx context.Context) (*Network, error) {
	return c.authManager.GetCurrentNetwork()
}

// IsAuthenticated returns true if the client has a valid access token.
func (c *Client) IsAuthenticated() bool {
	return c.authManager.IsAuthenticated()
}

// IsNetworkSet returns true if a network context has been established.
func (c *Client) IsNetworkSet() bool {
	return c.authManager.IsNetworkSet()
}

// GetAccessToken returns the current access token.
// Returns an error if not authenticated.
func (c *Client) GetAccessToken() (string, error) {
	return c.authManager.GetToken()
}

// Config returns a copy of the client configuration.
func (c *Client) Config() config.Config {
	return *c.config
}

// EnsureReady ensures the client is authenticated and has a network context.
//
// This is a convenience method that calls Authenticate() and then attempts
// to set a network if one is configured but not yet set.
func (c *Client) EnsureReady(ctx context.Context) error {
	// Ensure authentication
	if err := c.Authenticate(ctx); err != nil {
		return err
	}

	// Ensure network context if configured
	if c.config.NetworkName != "" && !c.IsNetworkSet() {
		if err := c.SetNetwork(ctx, c.config.NetworkName); err != nil {
			return err
		}
	}

	return nil
}

// WithAuthentication executes a function with authentication ensured.
//
// This is a convenience method for operations that require authentication
// but not necessarily a network context.
func (c *Client) WithAuthentication(ctx context.Context, fn func() error) error {
	if err := c.Authenticate(ctx); err != nil {
		return err
	}
	return fn()
}

// WithNetworkContext executes a function with both authentication and network context ensured.
//
// This is a convenience method for operations that require both authentication
// and a network to be selected.
func (c *Client) WithNetworkContext(ctx context.Context, fn func() error) error {
	if err := c.EnsureReady(ctx); err != nil {
		return err
	}

	if !c.IsNetworkSet() {
		return errors.NewAuthError("no network selected - use SetNetwork() or configure BS_NETWORK", nil)
	}

	return fn()
}
