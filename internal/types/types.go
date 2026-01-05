package types

import (
	"encoding/json"
	"time"
)

// TokenResponse represents the OAuth2 token response from BSN.cloud.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
}

// Network represents a BSN.cloud network.
type Network struct {
	ID               int                  `json:"id"`
	Name             string               `json:"name"`
	CreationDate     time.Time           `json:"creationDate"`
	LastModifiedDate time.Time           `json:"lastModifiedDate"`
	Subscription     *NetworkSubscription `json:"subscription"`
	Settings         *NetworkSettings     `json:"settings"`
	IsLockedOut      bool                `json:"isLockedOut"`
	LockoutDate      *time.Time          `json:"lockoutDate,omitempty"`
}

// NetworkSubscription represents subscription details for a network.
type NetworkSubscription struct {
	Level     string     `json:"level"`     // "Content" or "Control"
	StartDate time.Time  `json:"startDate"`
	EndDate   *time.Time `json:"endDate,omitempty"`
}

// NetworkSettings represents network configuration settings.
type NetworkSettings struct {
	UserAccessTokenLifetime                   string `json:"userAccessTokenLifetime"`
	DeviceAccessTokenLifetime                 string `json:"deviceAccessTokenLifetime"`
	AutomaticTaggedPlaylistApprovalEnabled    bool   `json:"automaticTaggedPlaylistApprovalEnabled"`
}

// NetworkRequest represents a request to set the active network.
type NetworkRequest struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Device represents a BrightSign device in the network.
type Device struct {
	ID                int              `json:"id"`
	Serial            string           `json:"serial"`
	Model             string           `json:"model"`
	Family            string           `json:"family"`
	RegistrationDate  time.Time        `json:"registrationDate"`
	LastModifiedDate  time.Time        `json:"lastModifiedDate"`
	Settings          *DeviceSettings  `json:"settings,omitempty"`
	Status            *DeviceStatusEmbed `json:"status,omitempty"`
}

// DeviceSettings represents device configuration settings.
type DeviceSettings struct {
	Name                   string  `json:"name"`
	Description            string  `json:"description"`
	ConcatNameAndSerial    bool    `json:"concatNameAndSerial"`
	SetupType              string  `json:"setupType"`
	LastModifiedDate       time.Time `json:"lastModifiedDate"`
	Group                  *Group  `json:"group,omitempty"`
	Timezone               string  `json:"timezone"`
}

// Group represents a device group.
type Group struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Link    string   `json:"link,omitempty"`
	Devices []Device `json:"devices,omitempty"` // Devices in this group (populated when getting group by name/ID)
}

// GroupList represents a list of groups.
type GroupList struct {
	Items []Group `json:"items"`
}

// DeviceList represents a paginated list of devices.
type DeviceList struct {
	Items       []Device `json:"items"`
	IsTruncated bool     `json:"isTruncated"`
	NextMarker  string   `json:"nextMarker,omitempty"`
	TotalCount  int      `json:"totalCount,omitempty"`
}

// DeviceDownload represents a content download on a device.
type DeviceDownload struct {
	ID              int       `json:"id"`
	FileName        string    `json:"fileName"`
	FileSize        int64     `json:"fileSize"`
	DownloadedBytes int64     `json:"downloadedBytes"`
	Status          string    `json:"status"` // pending, downloading, complete, error
	StartTime       time.Time `json:"startTime"`
	EndTime         time.Time `json:"endTime,omitempty"`
	Error           string    `json:"error,omitempty"`
}

// DeviceDownloadList represents a list of device downloads.
type DeviceDownloadList struct {
	Items      []DeviceDownload `json:"items"`
	TotalCount int              `json:"totalCount"`
}

// DeviceOperation represents an operation performed on a device.
type DeviceOperation struct {
	ID            int       `json:"id"`
	OperationType string    `json:"operationType"` // reboot, reprovision, update, etc.
	Status        string    `json:"status"`        // pending, in_progress, completed, failed
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
	StartedAt     time.Time `json:"startedAt,omitempty"`
	CompletedAt   time.Time `json:"completedAt,omitempty"`
	Error         string    `json:"error,omitempty"`
	Progress      int       `json:"progress,omitempty"` // 0-100
}

// DeviceOperationList represents a list of device operations.
type DeviceOperationList struct {
	Items      []DeviceOperation `json:"items"`
	TotalCount int               `json:"totalCount"`
}

// Region represents a rectangular region for screenshots.
type Region struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// DeviceStatusEmbed represents the status embedded in a device object.
type DeviceStatusEmbed struct {
	Group            *Group          `json:"group,omitempty"`
	Timezone         string          `json:"timezone"`
	Synchronization  *SyncSettings   `json:"synchronization,omitempty"`
	Script           *ScriptInfo     `json:"script,omitempty"`
	Firmware         *FirmwareInfo   `json:"firmware,omitempty"`
	Storage          []StorageInfo   `json:"storage,omitempty"`
	Network          *NetworkInfo    `json:"network,omitempty"`
	Uptime           string          `json:"uptime"`
	Health           string          `json:"health"`
	LastModifiedDate time.Time       `json:"lastModifiedDate"`
}

// SyncSettings represents synchronization settings
type SyncSettings struct {
	Status   SyncPeriod `json:"status"`
	Settings SyncPeriod `json:"settings"`
	Schedule SyncPeriod `json:"schedule"`
	Content  SyncPeriod `json:"content"`
}

// SyncPeriod represents a synchronization period
type SyncPeriod struct {
	Period string `json:"period"`
}

// ScriptInfo represents script-related information.
type ScriptInfo struct {
	Type    string   `json:"type"`
	Version string   `json:"version"`
	Plugins []string `json:"plugins,omitempty"`
}

// FirmwareInfo represents firmware-related information.
type FirmwareInfo struct {
	Version string `json:"version"`
}

// StorageInfo represents storage-related information.
type StorageInfo struct {
	Interface string `json:"interface"`
	System    string `json:"system"`
	Total     int64  `json:"total"`
	Used      int64  `json:"used"`
	Free      int64  `json:"free"`
}

// NetworkInfo represents network status information.
type NetworkInfo struct {
	ExternalIP string                `json:"externalIp"`
	Interfaces []NetworkInterfaceBSN `json:"interfaces"`
}

// NetworkInterfaceBSN represents a network interface from BSN.cloud
type NetworkInterfaceBSN struct {
	Enabled bool     `json:"enabled"`
	Proto   string   `json:"proto"`
	IP      []string `json:"ip"`
	Gateway *string  `json:"gateway"`
	DNS     *string  `json:"dns"`
	Name    string   `json:"name"`
	Type    string   `json:"type"`
}

// DeviceStatus represents the current operational status of a device.
// This is constructed from the device details including embedded status.
type DeviceStatus struct {
	DeviceID        string    `json:"deviceId"`
	Serial          string    `json:"serial"`
	Model           string    `json:"model"`
	FirmwareVersion string    `json:"firmwareVersion"`
	IsOnline        bool      `json:"isOnline"`
	LastSeen        time.Time `json:"lastSeen"`
	Status          string    `json:"status"`
	Uptime          int       `json:"uptime"`          // Uptime in seconds
	UptimeDisplay   string    `json:"uptimeDisplay"`   // Human readable uptime
	HealthStatus    string    `json:"healthStatus"`    // "Healthy", "Warning", "Error"
	LastHealthCheck time.Time `json:"lastHealthCheck"`
	IPAddress       string    `json:"ipAddress,omitempty"`
	ConnectionType  string    `json:"connectionType"`  // "ethernet", "wifi"
	SignalStrength  int       `json:"signalStrength,omitempty"` // For wifi connections
}

// DeviceError represents an error log entry from a device.
// The actual API structure may vary - this attempts to be flexible.
type DeviceError struct {
	ID               int       `json:"id,omitempty"`
	DeviceID         string    `json:"deviceId,omitempty"`
	Serial           string    `json:"serial,omitempty"`
	ErrorCode        string    `json:"errorCode,omitempty"`
	ErrorType        string    `json:"errorType,omitempty"`        // "system", "content", "network", etc.
	Severity         string    `json:"severity,omitempty"`         // "info", "warning", "error", "critical"  
	Level            string    `json:"level,omitempty"`            // Alternative severity field
	Message          string    `json:"message,omitempty"`
	Description      string    `json:"description,omitempty"`      // Alternative message field
	Details          string    `json:"details,omitempty"`
	Timestamp        time.Time `json:"timestamp,omitempty"`
	CreationDate     time.Time `json:"creationDate,omitempty"`     // Alternative timestamp field
	LastModifiedDate time.Time `json:"lastModifiedDate,omitempty"` // Alternative timestamp field
	Source           string    `json:"source,omitempty"`           // Component that generated the error
	Component        string    `json:"component,omitempty"`        // Alternative source field
	Resolved         bool      `json:"resolved,omitempty"`
	ResolvedAt       *time.Time `json:"resolvedAt,omitempty"`
	Status           string    `json:"status,omitempty"`           // Alternative resolved field
	
	// Generic fields for unknown API structure
	Type             string    `json:"type,omitempty"`
	Code             string    `json:"code,omitempty"`
	Name             string    `json:"name,omitempty"`
}

// DeviceErrorList represents a paginated list of device errors.
type DeviceErrorList struct {
	Items       []DeviceError `json:"items"`
	IsTruncated bool          `json:"isTruncated"`
	NextMarker  string        `json:"nextMarker,omitempty"`
	TotalCount  int           `json:"totalCount,omitempty"`
}

// RebootType represents the type of reboot to perform.
type RebootType string

const (
	// RebootTypeNormal performs a standard reboot
	RebootTypeNormal RebootType = "normal"
	// RebootTypeCrash performs a reboot and saves a crash report
	RebootTypeCrash RebootType = "crash"
	// RebootTypeFactoryReset performs a factory reset reboot
	RebootTypeFactoryReset RebootType = "factoryreset"
	// RebootTypeDisableAutorun disables autorun script and reboots
	RebootTypeDisableAutorun RebootType = "disableautorun"
)

// RebootResponse represents the response from a device reboot request.
type RebootResponse struct {
	DeviceID     string    `json:"deviceId,omitempty"`
	Serial       string    `json:"serial,omitempty"`
	OperationID  string    `json:"operationId,omitempty"`
	Status       string    `json:"status"`           // "success", "failed", "pending"
	Message      string    `json:"message,omitempty"`
	Timestamp    time.Time `json:"timestamp,omitempty"`
	RebootTime   time.Time `json:"rebootTime,omitempty"`
}

// SnapshotRequest represents a request to take a device screenshot.
type SnapshotRequest struct {
	Format          string  `json:"format,omitempty"`          // "png" or "jpeg"
	Quality         int     `json:"quality,omitempty"`         // 1-100 for JPEG quality
	IncludeMetadata bool    `json:"includeMetadata,omitempty"` // Include metadata in response
	Output          string  `json:"output,omitempty"`          // "base64" or "file"
	Compression     string  `json:"compression,omitempty"`     // "low", "medium", "high"
	DisplayID       int     `json:"displayId,omitempty"`       // Display to capture (multi-display devices)
	Region          *Region `json:"region,omitempty"`          // Optional region to capture
}

// SnapshotResponse represents the response from a device snapshot request.
type SnapshotResponse struct {
	RemoteSnapshotThumbnail string    `json:"remoteSnapshotThumbnail,omitempty"` // Base64 thumbnail
	Filename                string    `json:"filename,omitempty"`                // Full-resolution file path
	Timestamp               string    `json:"timestamp,omitempty"`               // When snapshot was taken
	DeviceName              string    `json:"devicename,omitempty"`              // Device that took snapshot
	Width                   int       `json:"width,omitempty"`                   // Image width in pixels
	Height                  int       `json:"height,omitempty"`                  // Image height in pixels
	Data                    string    `json:"data,omitempty"`                    // Base64 image data if requested
	Format                  string    `json:"format,omitempty"`                  // Image format used
	Size                    int64     `json:"size,omitempty"`                    // File size in bytes
}

// ReprovisionResponse represents the response from a device re-provision request.
type ReprovisionResponse struct {
	Success bool   `json:"success"`         // Whether the re-provision succeeded
	Message string `json:"message,omitempty"` // Optional message from the device
}

// DWSPasswordRequest represents a request to set the DWS password.
type DWSPasswordRequest struct {
	Password         string `json:"password"`          // New password (empty string removes password)
	PreviousPassword string `json:"previous_password"` // Current password (empty if no password set)
}

// DWSPasswordInfo represents information about the DWS password.
type DWSPasswordInfo struct {
	IsResultValid bool `json:"isResultValid"` // Whether the password info is valid
	IsBlank       bool `json:"isBlank"`       // Whether no password is currently set
}

// DWSPasswordGetResponse represents the response from getting DWS password info.
type DWSPasswordGetResponse struct {
	Success  bool             `json:"success"`  // Whether the operation succeeded
	Password *DWSPasswordInfo `json:"password"` // Password information
}

// DWSPasswordSetResponse represents the response from setting DWS password.
type DWSPasswordSetResponse struct {
	Success bool `json:"success"` // Whether the password was set successfully
	Reboot  bool `json:"reboot"`  // Whether the device will reboot
}

// BDeployRecord represents a B-Deploy setup record.
type BDeployRecord struct {
	ID           string `json:"_id"`         // Unique identifier for the setup record
	PackageName  string `json:"packageName"` // Name of the setup package
	SetupType    string `json:"setupType"`   // Type of setup (e.g., "web", "usb", etc.)
	Username     string `json:"username,omitempty"`    // User who created the setup
	NetworkName  string `json:"networkName,omitempty"` // Network name
	BSNGroupName string `json:"bsnGroupName,omitempty"` // BSN group name
	CreatedDate  string `json:"createdDate,omitempty"` // When the setup was created
	UpdatedDate  string `json:"updatedDate,omitempty"` // When the setup was last updated
	IsActive     bool   `json:"isActive,omitempty"`    // Whether the setup is active
}

// BDeployRecordList represents a list of B-Deploy setup records.
type BDeployRecordList struct {
	Items      []BDeployRecord `json:"items,omitempty"`      // Array of setup records
	TotalCount int             `json:"totalCount,omitempty"` // Total number of records
}

// BDeployAPIResponse represents the API response wrapper for B-Deploy records.
type BDeployAPIResponse struct {
	Error  interface{}     `json:"error"`  // Error field (null on success)
	Result []BDeployRecord `json:"result"` // Array of setup records
}

// BDeploySingleRecordResponse represents the API response for a single B-Deploy record.
type BDeploySingleRecordResponse struct {
	Error  interface{}   `json:"error"`  // Error field (null on success)
	Result BDeployRecord `json:"result"` // Single setup record
}

// BDeployFullRecordAPIResponse represents the API response for full B-Deploy setup records.
type BDeployFullRecordAPIResponse struct {
	Error  interface{}           `json:"error"`  // Error field (null on success)
	Result []BDeploySetupRecord  `json:"result"` // Array of full setup records
}

// BDeployInfo represents the B-Deploy section of a setup record.
type BDeployInfo struct {
	Username    string `json:"username"`    // BSN.cloud username
	NetworkName string `json:"networkName"` // BSN.cloud network name
	PackageName string `json:"packageName"` // Setup package name
	Client      string `json:"client,omitempty"` // Client identifier
}

// BSNTokenEntity represents the BSN device registration token.
type BSNTokenEntity struct {
	Token     string `json:"token"`     // Registration token
	Scope     string `json:"scope"`     // Token scope (e.g., "cert")
	ValidFrom string `json:"validFrom"` // Token valid from date
	ValidTo   string `json:"validTo"`   // Token valid to date
}

// NetworkInterface represents a network interface configuration.
type NetworkInterface struct {
	ID                         string `json:"id"`                         // Interface ID
	Name                       string `json:"name"`                       // Interface name
	Type                       string `json:"type"`                       // Interface type
	Proto                      string `json:"proto"`                      // Protocol (e.g., "DHCPv4")
	ContentDownloadEnabled     bool   `json:"contentDownloadEnabled"`     // Enable content downloads
	HealthReportingEnabled     bool   `json:"healthReportingEnabled"`     // Enable health reporting
}

// NetworkConfig represents network configuration.
type NetworkConfig struct {
	TimeServers []string           `json:"timeServers"` // Time server URLs
	Interfaces  []NetworkInterface `json:"interfaces"`  // Network interfaces
	Proxy       string             `json:"proxy,omitempty"`       // Proxy server
	ProxyBypass string             `json:"proxyBypass,omitempty"` // Proxy bypass list
}

// IdleScreenColor represents RGBA color for idle screen.
type IdleScreenColor struct {
	R int `json:"r"` // Red (0-255)
	G int `json:"g"` // Green (0-255)
	B int `json:"b"` // Blue (0-255)
	A int `json:"a"` // Alpha (0-1)
}

// BDeploySetupRecord represents a complete B-Deploy setup record for creation.
// Supports B-Deploy API v2.0.0 and v3.0.0 specifications.
type BDeploySetupRecord struct {
	// Core fields
	Version                          string          `json:"version"`                                    // Setup version (2.0.0 or 3.0.0)
	ID                               string          `json:"_id,omitempty"`                              // Record ID (empty for new records)
	BDeploy                          BDeployInfo     `json:"bDeploy"`                                    // B-Deploy information
	SetupType                        string          `json:"setupType"`                                  // Setup type (bsn, standalone, lfn)
	BSNDeviceRegistrationTokenEntity *BSNTokenEntity `json:"bsnDeviceRegistrationTokenEntity,omitempty"` // BSN device registration token
	TimeZone                         string          `json:"timeZone,omitempty"`                         // IANA timezone (e.g., "America/New_York")
	BSNGroupName                     string          `json:"bsnGroupName,omitempty"`                     // BSN group name (default: "Default")
	Network                          *NetworkConfig  `json:"network,omitempty"`                          // Network configuration

	// Firmware & Debugging
	FirmwareUpdateType        string `json:"firmwareUpdateType,omitempty"`        // Firmware update type (standard, latest, etc.)
	EnableSerialDebugging     bool   `json:"enableSerialDebugging,omitempty"`     // Enable serial port debugging output
	EnableSystemLogDebugging  bool   `json:"enableSystemLogDebugging,omitempty"`  // Enable system log debugging

	// DWS (Diagnostic Web Server)
	DWSEnabled       bool   `json:"dwsEnabled,omitempty"`       // Enable local DWS
	DWSPassword      string `json:"dwsPassword,omitempty"`      // DWS password ("none" for no password)
	RemoteDWSEnabled bool   `json:"remoteDwsEnabled,omitempty"` // Enable remote DWS via BSN.cloud

	// LWS (Local Web Server)
	LWSEnabled                    bool   `json:"lwsEnabled,omitempty"`                    // Enable LWS
	LWSConfig                     string `json:"lwsConfig,omitempty"`                     // LWS configuration mode (status, content, etc.)
	LWSUserName                   string `json:"lwsUserName,omitempty"`                   // LWS username
	LWSPassword                   string `json:"lwsPassword,omitempty"`                   // LWS password
	LWSEnableUpdateNotifications  bool   `json:"lwsEnableUpdateNotifications,omitempty"`  // Enable LWS update notifications

	// Device information
	DeviceName         string `json:"deviceName,omitempty"`         // Device name
	DeviceDescription  string `json:"deviceDescription,omitempty"`  // Device description
	UnitNamingMethod   string `json:"unitNamingMethod,omitempty"`   // Unit naming method (appendUnitIDToUnitName, etc.)

	// BSN connection settings
	TimeBetweenNetConnects int `json:"timeBetweenNetConnects,omitempty"` // Time between network connects (seconds)
	TimeBetweenHeartbeats  int `json:"timeBetweenHeartbeats,omitempty"`  // Time between heartbeats (seconds)

	// SFN (Simple File Networking)
	SFNWebFolderURL               string `json:"sfnWebFolderUrl,omitempty"`               // SFN web folder URL
	SFNUserName                   string `json:"sfnUserName,omitempty"`                   // SFN username
	SFNPassword                   string `json:"sfnPassword,omitempty"`                   // SFN password
	SFNEnableBasicAuthentication  bool   `json:"sfnEnableBasicAuthentication,omitempty"`  // Enable SFN basic auth

	// Logging
	PlaybackLoggingEnabled         bool   `json:"playbackLoggingEnabled,omitempty"`         // Enable playback logging
	EventLoggingEnabled            bool   `json:"eventLoggingEnabled,omitempty"`            // Enable event logging
	DiagnosticLoggingEnabled       bool   `json:"diagnosticLoggingEnabled,omitempty"`       // Enable diagnostic logging
	StateLoggingEnabled            bool   `json:"stateLoggingEnabled,omitempty"`            // Enable state logging
	VariableLoggingEnabled         bool   `json:"variableLoggingEnabled,omitempty"`         // Enable variable logging
	UploadLogFilesAtBoot           bool   `json:"uploadLogFilesAtBoot,omitempty"`           // Upload log files at boot
	UploadLogFilesAtSpecificTime   bool   `json:"uploadLogFilesAtSpecificTime,omitempty"`   // Upload logs at specific time
	UploadLogFilesTime             int    `json:"uploadLogFilesTime,omitempty"`             // Time to upload logs (hour 0-23)
	LogHandlerURL                  string `json:"logHandlerUrl,omitempty"`                  // Custom log handler URL

	// Remote snapshot (screenshots)
	EnableRemoteSnapshot              bool   `json:"enableRemoteSnapshot,omitempty"`              // Enable remote screenshots
	RemoteSnapshotInterval            int    `json:"remoteSnapshotInterval,omitempty"`            // Screenshot interval (minutes)
	RemoteSnapshotMaxImages           int    `json:"remoteSnapshotMaxImages,omitempty"`           // Max screenshots to keep
	RemoteSnapshotJPEGQualityLevel    int    `json:"remoteSnapshotJpegQualityLevel,omitempty"`    // JPEG quality (1-100)
	RemoteSnapshotScreenOrientation   string `json:"remoteSnapshotScreenOrientation,omitempty"`   // Screen orientation (Landscape/Portrait)
	RemoteSnapshotHandlerURL          string `json:"remoteSnapshotHandlerUrl,omitempty"`          // Custom snapshot handler URL

	// Device screenshots (different from remote snapshot)
	DeviceScreenShotsEnabled     bool   `json:"deviceScreenShotsEnabled,omitempty"`     // Enable device screenshots
	DeviceScreenShotsInterval    int    `json:"deviceScreenShotsInterval,omitempty"`    // Screenshot interval (seconds)
	DeviceScreenShotsCountLimit  int    `json:"deviceScreenShotsCountLimit,omitempty"`  // Max screenshots
	DeviceScreenShotsQuality     int    `json:"deviceScreenShotsQuality,omitempty"`     // JPEG quality (1-100)
	DeviceScreenShotsOrientation string `json:"deviceScreenShotsOrientation,omitempty"` // Screen orientation

	// Wireless
	UseWireless bool   `json:"useWireless,omitempty"` // Enable wireless networking
	SSID        string `json:"ssid,omitempty"`        // WiFi SSID
	Passphrase  string `json:"passphrase,omitempty"`  // WiFi password

	// Time server (for backwards compatibility - use Network.TimeServers instead)
	TimeServer string `json:"timeServer,omitempty"` // NTP time server URL

	// Display settings
	IdleScreenColor      *IdleScreenColor `json:"idleScreenColor,omitempty"`      // Idle screen color (RGBA)
	UseCustomSplashScreen bool            `json:"useCustomSplashScreen,omitempty"` // Use custom splash screen

	// Network data types - Wired
	ContentDataTypeEnabledWired     bool `json:"contentDataTypeEnabledWired,omitempty"`     // Enable content downloads on wired
	TextFeedsDataTypeEnabledWired   bool `json:"textFeedsDataTypeEnabledWired,omitempty"`   // Enable text feeds on wired
	HealthDataTypeEnabledWired      bool `json:"healthDataTypeEnabledWired,omitempty"`      // Enable health reporting on wired
	MediaFeedsDataTypeEnabledWired  bool `json:"mediaFeedsDataTypeEnabledWired,omitempty"`  // Enable media feeds on wired
	LogUploadsXfersEnabledWired     bool `json:"logUploadsXfersEnabledWired,omitempty"`     // Enable log uploads on wired

	// Network data types - Wireless
	ContentDataTypeEnabledWireless    bool `json:"contentDataTypeEnabledWireless,omitempty"`    // Enable content downloads on wireless
	TextFeedsDataTypeEnabledWireless  bool `json:"textFeedsDataTypeEnabledWireless,omitempty"`  // Enable text feeds on wireless
	HealthDataTypeEnabledWireless     bool `json:"healthDataTypeEnabledWireless,omitempty"`     // Enable health reporting on wireless
	MediaFeedsDataTypeEnabledWireless bool `json:"mediaFeedsDataTypeEnabledWireless,omitempty"` // Enable media feeds on wireless
	LogUploadsXfersEnabledWireless    bool `json:"logUploadsXfersEnabledWireless,omitempty"`    // Enable log uploads on wireless

	// Rate limiting - Wired (primary interface)
	RateLimitModeOutsideWindow      string `json:"rateLimitModeOutsideWindow,omitempty"`      // Rate limit mode outside window (default, unlimited, limited)
	RateLimitRateOutsideWindow      int    `json:"rateLimitRateOutsideWindow,omitempty"`      // Rate limit outside window (kbps)
	RateLimitModeInWindow           string `json:"rateLimitModeInWindow,omitempty"`           // Rate limit mode in window
	RateLimitRateInWindow           int    `json:"rateLimitRateInWindow,omitempty"`           // Rate limit in window (kbps)
	RateLimitModeInitialDownloads   string `json:"rateLimitModeInitialDownloads,omitempty"`   // Rate limit mode for initial downloads
	RateLimitRateInitialDownloads   int    `json:"rateLimitRateInitialDownloads,omitempty"`   // Rate limit for initial downloads (kbps)

	// Rate limiting - Wireless (secondary interface)
	RateLimitModeOutsideWindow2    string `json:"rateLimitModeOutsideWindow_2,omitempty"`    // Rate limit mode outside window for interface 2
	RateLimitRateOutsideWindow2    int    `json:"rateLimitRateOutsideWindow_2,omitempty"`    // Rate limit outside window for interface 2 (kbps)
	RateLimitModeInWindow2         string `json:"rateLimitModeInWindow_2,omitempty"`         // Rate limit mode in window for interface 2
	RateLimitRateInWindow2         int    `json:"rateLimitRateInWindow_2,omitempty"`         // Rate limit in window for interface 2 (kbps)
	RateLimitModeInitialDownloads2 string `json:"rateLimitModeInitialDownloads_2,omitempty"` // Rate limit mode for initial downloads for interface 2
	RateLimitRateInitialDownloads2 int    `json:"rateLimitRateInitialDownloads_2,omitempty"` // Rate limit for initial downloads for interface 2 (kbps)

	// Network priority and diagnostics
	NetworkConnectionPriority string `json:"networkConnectionPriority,omitempty"` // Network priority (wired, wireless)
	NetworkDiagnosticsEnabled bool   `json:"networkDiagnosticsEnabled,omitempty"` // Enable network diagnostics
	TestEthernetEnabled       bool   `json:"testEthernetEnabled,omitempty"`       // Enable ethernet testing
	TestWirelessEnabled       bool   `json:"testWirelessEnabled,omitempty"`       // Enable wireless testing
	TestInternetEnabled       bool   `json:"testInternetEnabled,omitempty"`       // Enable internet connectivity testing

	// BrightWall
	BrightWallName         string `json:"BrightWallName,omitempty"`         // BrightWall configuration name
	BrightWallScreenNumber string `json:"BrightWallScreenNumber,omitempty"` // BrightWall screen number

	// Hostname
	SpecifyHostname bool   `json:"specifyHostname,omitempty"` // Specify custom hostname
	Hostname        string `json:"hostname,omitempty"`        // Custom hostname

	// Proxy
	UseProxy     bool   `json:"useProxy,omitempty"`     // Use proxy server
	ProxyAddress string `json:"proxyAddress,omitempty"` // Proxy server address
	ProxyPort    int    `json:"proxyPort,omitempty"`    // Proxy server port

	// Network hosts
	NetworkHosts []string `json:"networkHosts,omitempty"` // Custom network hosts

	// Content download windows
	ContentDownloadsRestricted bool `json:"contentDownloadsRestricted,omitempty"` // Restrict content downloads to time window
	ContentDownloadRangeStart  int  `json:"contentDownloadRangeStart,omitempty"`  // Content download start time (minutes from midnight)
	ContentDownloadRangeEnd    int  `json:"contentDownloadRangeEnd,omitempty"`    // Content download end time (minutes from midnight)

	// Heartbeat windows
	HeartbeatsRestricted  bool `json:"heartbeatsRestricted,omitempty"`  // Restrict heartbeats to time window
	HeartbeatsRangeStart  int  `json:"heartbeatsRangeStart,omitempty"`  // Heartbeat start time (minutes from midnight)
	HeartbeatsRangeEnd    int  `json:"heartbeatsRangeEnd,omitempty"`    // Heartbeat end time (minutes from midnight)

	// Static IP configuration - Interface 1 (wired)
	UseDHCP         bool   `json:"useDHCP,omitempty"`         // Use DHCP for interface 1
	StaticIPAddress string `json:"staticIPAddress,omitempty"` // Static IP address for interface 1
	SubnetMask      string `json:"subnetMask,omitempty"`      // Subnet mask for interface 1
	Gateway         string `json:"gateway,omitempty"`         // Gateway for interface 1
	DNS1            string `json:"dns1,omitempty"`            // DNS server 1 for interface 1
	DNS2            string `json:"dns2,omitempty"`            // DNS server 2 for interface 1
	DNS3            string `json:"dns3,omitempty"`            // DNS server 3 for interface 1

	// Static IP configuration - Interface 2 (wireless)
	UseDHCP2         bool   `json:"useDHCP_2,omitempty"`         // Use DHCP for interface 2
	StaticIPAddress2 string `json:"staticIPAddress_2,omitempty"` // Static IP address for interface 2
	SubnetMask2      string `json:"subnetMask_2,omitempty"`      // Subnet mask for interface 2
	Gateway2         string `json:"gateway_2,omitempty"`         // Gateway for interface 2
	DNS1_2           string `json:"dns1_2,omitempty"`            // DNS server 1 for interface 2
	DNS2_2           string `json:"dns2_2,omitempty"`            // DNS server 2 for interface 2
	DNS3_2           string `json:"dns3_2,omitempty"`            // DNS server 3 for interface 2

	// USB updates
	USBUpdatePassword string `json:"usbUpdatePassword,omitempty"` // Password for USB updates
}

// BDeployCreateResponse represents the response from creating a B-Deploy record.
type BDeployCreateResponse struct {
	ID      string `json:"_id,omitempty"`      // Created record ID
	Success bool   `json:"success,omitempty"`  // Whether creation succeeded
	Error   string `json:"error,omitempty"`    // Error message if failed
}

// BDeployCreateAPIResponse represents the API wrapper response for creating a B-Deploy record.
type BDeployCreateAPIResponse struct {
	Error  interface{} `json:"error"`  // Error field (null on success)
	Result string      `json:"result"` // Created setup record ID
}

// BDeployUpdateAPIResponse represents the API wrapper response for updating a B-Deploy record.
type BDeployUpdateAPIResponse struct {
	Error  interface{}          `json:"error"`  // Error field (null on success)
	Result *BDeploySetupRecord  `json:"result"` // Updated setup record
}

// BDeployDeleteResponse represents the response from deleting a B-Deploy record.
type BDeployDeleteResponse struct {
	Success bool   `json:"success,omitempty"`  // Whether deletion succeeded
	Message string `json:"message,omitempty"`  // Response message
	Error   string `json:"error,omitempty"`    // Error message if failed
}

// BDeployDevice represents a device in the B-Deploy device response.
type BDeployDevice struct {
	ID          string    `json:"_id"`
	Client      string    `json:"client"`
	NetworkName string    `json:"NetworkName"`
	Username    string    `json:"username"`
	Serial      string    `json:"serial"`
	Name        string    `json:"name"`
	Model       string    `json:"model"`
	Desc        string    `json:"desc"`
	SetupName   string    `json:"setupName"`   // Used by individual device lookup
	SetupID     string    `json:"setupId"`     // Used by device list 
	URL         string    `json:"url"`         // Used by device list
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Version     int       `json:"__v"`
}

// BDeployDeviceResult represents the result portion of the B-Deploy device response.
type BDeployDeviceResult struct {
	Total   int             `json:"total"`
	Matched int             `json:"matched"`
	Players []BDeployDevice `json:"players"`
	Priv    string          `json:"priv"`
}

// BDeployDeviceResponse represents the full response from the B-Deploy device API.
type BDeployDeviceResponse struct {
	Error  interface{}         `json:"error"`
	Result BDeployDeviceResult `json:"result"`
}

// BDeployDeviceListResponse represents the response from the B-Deploy device list API.
type BDeployDeviceListResponse struct {
	Total   int             `json:"total"`
	Matched int             `json:"matched"`
	Players []BDeployDevice `json:"players"`
}

// BDeployDeviceListAPIResponse represents the API wrapper for the B-Deploy device list response.
type BDeployDeviceListAPIResponse struct {
	Error  interface{}                `json:"error"`
	Result *BDeployDeviceListResponse `json:"result"`
}

// NetworkContextRequest represents a request to set the network context.
type NetworkContextRequest struct {
	Name string `json:"name"` // Network name to set as context
}

// BDeployDeviceRequest represents a request to create or update a B-Deploy device.
type BDeployDeviceRequest struct {
	ID          string `json:"_id,omitempty"`      // Device ID (for updates)
	Username    string `json:"username"`            // BSN.cloud username
	Serial      string `json:"serial"`              // Device serial number
	Name        string `json:"name"`                // Device name
	NetworkName string `json:"NetworkName"`         // Network name
	Model       string `json:"model,omitempty"`     // Device model
	Desc        string `json:"desc,omitempty"`      // Device description
	SetupID     string `json:"setupId,omitempty"`   // Setup ID to associate with
	URL         string `json:"url,omitempty"`       // Direct presentation URL (alternative to setupId)
	UserData    string `json:"userdata,omitempty"`  // Additional header information
}

// BDeployDeviceCreateResponse represents the response from creating a device.
type BDeployDeviceCreateResponse struct {
	Error  interface{} `json:"error"`
	Result string      `json:"result"` // Device ID
}

// BDeployDeviceUpdateResponse represents the response from updating a device.
type BDeployDeviceUpdateResponse struct {
	Error  interface{}     `json:"error"`
	Result *BDeployDevice  `json:"result"`
}

// RDWSInfo represents player information from the /info/ endpoint
type RDWSInfo struct {
	Serial           string                   `json:"serial"`
	Model            string                   `json:"model"`
	FWVersion        string                   `json:"FWVersion"`
	BootVersion      string                   `json:"bootVersion"`
	Family           string                   `json:"family"`
	IsPlayer         bool                     `json:"isPlayer"`
	UpTime           string                   `json:"upTime"`
	UpTimeSeconds    int                      `json:"upTimeSeconds"`
	ConnectionType   string                   `json:"connectionType"`
	Ethernet         []RDWSNetworkInterface   `json:"ethernet,omitempty"`
	Wireless         []RDWSNetworkInterface   `json:"wireless,omitempty"`
	Interfaces       []RDWSNetworkInterface   `json:"interfaces,omitempty"`
	Power            *RDWSInfoSubResult       `json:"power,omitempty"`
	POE              *RDWSInfoSubResult       `json:"poe,omitempty"`
	Extensions       *RDWSInfoSubResult       `json:"extensions,omitempty"`
	Blessings        *RDWSInfoSubResult       `json:"blessings,omitempty"`
	Networking       *RDWSInfoSubResult       `json:"networking,omitempty"`
	BVNPipelines     *RDWSInfoSubResult       `json:"bvnPipelines,omitempty"`
	BVNComponents    *RDWSInfoSubResult       `json:"bvnComponents,omitempty"`
	HardwareFeatures map[string]interface{}   `json:"hardware_features,omitempty"`
	APIFeatures      map[string]interface{}   `json:"api_features,omitempty"`
	ActiveFeatures   map[string]interface{}   `json:"active_features,omitempty"`
	BSNCE            bool                     `json:"bsnce"`
}

// RDWSNetworkInterface represents a network interface from the player
type RDWSNetworkInterface struct {
	InterfaceName string          `json:"interfaceName"`
	InterfaceType string          `json:"interfaceType"`
	IPv4          []RDWSIPAddress `json:"IPv4,omitempty"`
	IPv6          []RDWSIPAddress `json:"IPv6,omitempty"`
}

// RDWSIPAddress represents an IP address configuration
type RDWSIPAddress struct {
	Address  string `json:"address"`
	Netmask  string `json:"netmask"`
	Family   string `json:"family"`
	MAC      string `json:"mac"`
	Internal bool   `json:"internal"`
	CIDR     string `json:"cidr"`
	ScopeID  int    `json:"scopeid,omitempty"`
}

// RDWSInfoSubResult represents nested result objects in the info response
type RDWSInfoSubResult struct {
	Result interface{} `json:"result"`
}

// RDWSInfoResponse represents the full API response wrapper for /info/
type RDWSInfoResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result RDWSInfo `json:"result"`
	} `json:"data"`
}

// RDWSTimeInfo represents time information from the /time/ endpoint
type RDWSTimeInfo struct {
	Time         string `json:"time"`
	TimezoneMin  *int   `json:"timezone_mins"`
	TimezoneName string `json:"timezone_name"`
	TimezoneAbbr string `json:"timezone_abbr"`
	Year         int    `json:"year"`
	Month        int    `json:"month"`
	Date         int    `json:"date"`
	Hour         int    `json:"hour"`
	Minute       int    `json:"minute"`
	Second       int    `json:"second"`
	Millisecond  int    `json:"millisecond"`
}

// RDWSTimeResponse represents the full API response wrapper for GET /time/
type RDWSTimeResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result RDWSTimeInfo `json:"result"`
	} `json:"data"`
}

// RDWSTimeSetRequest represents the request body for PUT /time/
type RDWSTimeSetRequest struct {
	Time          string `json:"time"`          // Format: "hh:mm:ss <timezone>" or "hh:mm:ss"
	Date          string `json:"date"`          // Format: "yyyy-mm-dd"
	ApplyTimezone bool   `json:"applyTimezone"` // Apply with player's timezone (true) or UTC (false)
}

// RDWSTimeSetResponse represents the response from PUT /time/
type RDWSTimeSetResponse struct {
	Data struct {
		Result bool `json:"result"`
	} `json:"data"`
}

// RDWSHealthInfo represents health status from the /health/ endpoint
type RDWSHealthInfo struct {
	Status     string `json:"status"`     // "active"
	StatusTime string `json:"statusTime"` // Format: "yyyy-mm-dd hh:mm:ss <timezone>"
}

// RDWSHealthResponse represents the full API response wrapper for /health/
type RDWSHealthResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result RDWSHealthInfo `json:"result"`
	} `json:"data"`
}

// RDWSFileStat represents file statistics from the fs module
type RDWSFileStat struct {
	Dev         int64  `json:"dev"`
	Mode        int    `json:"mode"`
	Nlink       int    `json:"nlink"`
	UID         int    `json:"uid"`
	GID         int    `json:"gid"`
	Rdev        int    `json:"rdev"`
	BlkSize     int    `json:"blksize"`
	Ino         int64  `json:"ino"`
	Size        int64  `json:"size"`
	Blocks      int64  `json:"blocks"`
	AtimeMs     int64  `json:"atimeMs"`
	MtimeMs     int64  `json:"mtimeMs"`
	CtimeMs     int64  `json:"ctimeMs"`
	BirthtimeMs int64  `json:"birthtimeMs"`
	Atime       string `json:"atime"`
	Mtime       string `json:"mtime"`
	Ctime       string `json:"ctime"`
	Birthtime   string `json:"birthtime"`
}

// RDWSFileInfo represents a file or directory entry
type RDWSFileInfo struct {
	Name       string        `json:"name"`
	Type       string        `json:"type"` // "file" or "dir"
	Path       string        `json:"path"`
	Stat       *RDWSFileStat `json:"stat,omitempty"`
	Mime       string        `json:"mime,omitempty"`
	FileSize   int64         `json:"fileSize,omitempty"`
	Streamable bool          `json:"streamable,omitempty"`
	Children   []RDWSFileInfo `json:"children,omitempty"`
}

// RDWSStorageStats represents storage device statistics
type RDWSStorageStats struct {
	BlockSize  int64 `json:"blockSize"`
	BytesFree  int64 `json:"bytesFree"`
	FilesFree  int64 `json:"filesFree"`
	FilesUsed  int64 `json:"filesUsed"`
	IsReadOnly bool  `json:"isReadOnly"`
	SizeBytes  int64 `json:"sizeBytes"`
}

// RDWSStorageInfo represents storage device information
type RDWSStorageInfo struct {
	FileSystemType string            `json:"fileSystemType"`
	Stats          *RDWSStorageStats `json:"stats,omitempty"`
	MountedOn      string            `json:"mountedOn"`
}

// RDWSFileListResult represents the result of listing files
type RDWSFileListResult struct {
	Name        string           `json:"name,omitempty"`
	Type        string           `json:"type,omitempty"`
	Path        string           `json:"path,omitempty"`
	Stat        *RDWSFileStat    `json:"stat,omitempty"`
	Files       []RDWSFileInfo   `json:"files,omitempty"`
	Contents    []RDWSFileInfo   `json:"contents,omitempty"`
	StorageInfo *RDWSStorageInfo `json:"storageInfo,omitempty"`
}

// RDWSFileListResponse represents the response from listing files
type RDWSFileListResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result RDWSFileListResult `json:"result"`
	} `json:"data"`
}

// RDWSFileUploadItem represents a single file to upload
type RDWSFileUploadItem struct {
	FileName     string `json:"fileName"`
	FileContents string `json:"fileContents"` // Plain text or Data URL (base64)
	FileType     string `json:"fileType"`     // MIME type
}

// RDWSFileUploadRequest represents the request to upload files
type RDWSFileUploadRequest struct {
	Data struct {
		FileUploadPath string               `json:"fileUploadPath"`
		Files          []RDWSFileUploadItem `json:"files"`
	} `json:"data"`
}

// RDWSFileUploadResult represents the result of file upload
type RDWSFileUploadResult struct {
	Success bool     `json:"success"`
	Results []string `json:"results,omitempty"`
}

// RDWSFileUploadResponse represents the response from file upload
type RDWSFileUploadResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result RDWSFileUploadResult `json:"result"`
	} `json:"data"`
}

// RDWSFileRenameRequest represents the request to rename a file
type RDWSFileRenameRequest struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

// RDWSFileOperationResponse represents a generic success/error response for file operations
type RDWSFileOperationResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result struct {
			Success bool   `json:"success,omitempty"`
			Message string `json:"message,omitempty"`
			Error   string `json:"error,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// RDWSLocalDWSInfo represents the current state of local DWS
type RDWSLocalDWSInfo struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port,omitempty"`
}

// RDWSLocalDWSResponse represents the response from GET /control/local-dws/
type RDWSLocalDWSResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result RDWSLocalDWSInfo `json:"result"`
	} `json:"data"`
}

// RDWSLocalDWSSetRequest represents the request to enable/disable local DWS
type RDWSLocalDWSSetRequest struct {
	Data struct {
		Enabled bool `json:"enabled"`
	} `json:"data"`
}

// RDWSLocalDWSSetResponse represents the response from PUT /control/local-dws/
type RDWSLocalDWSSetResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result struct {
			Success bool   `json:"success"`
			Message string `json:"message,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// RDWSDiagnosticsInfo represents network diagnostics information
type RDWSDiagnosticsInfo struct {
	Gateway            string   `json:"gateway,omitempty"`
	DNS                []string `json:"dns,omitempty"`
	ConnectedToRouter  bool     `json:"connectedToRouter,omitempty"`
	ConnectedToInternet bool    `json:"connectedToInternet,omitempty"`
	ExternalIPAddress  string   `json:"externalIpAddress,omitempty"`
}

// RDWSDiagnosticsResponse represents the response from GET /diagnostics/
type RDWSDiagnosticsResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result RDWSDiagnosticsInfo `json:"result"`
	} `json:"data"`
}

// RDWSDNSLookupResult represents DNS lookup result
type RDWSDNSLookupResult struct {
	Success   bool     `json:"success"`
	Domain    string   `json:"domain,omitempty"`
	Addresses []string `json:"addresses,omitempty"`
	Error     string   `json:"error,omitempty"`
}

// RDWSDNSLookupResponse represents the response from GET /diagnostics/dns-lookup/
type RDWSDNSLookupResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result RDWSDNSLookupResult `json:"result"`
	} `json:"data"`
}

// RDWSPingResult represents ping diagnostic result
type RDWSPingResult struct {
	Success    bool    `json:"success"`
	Host       string  `json:"host,omitempty"`
	PacketLoss float64 `json:"packetLoss,omitempty"`
	MinRTT     float64 `json:"minRtt,omitempty"`
	MaxRTT     float64 `json:"maxRtt,omitempty"`
	AvgRTT     float64 `json:"avgRtt,omitempty"`
	Output     string  `json:"output,omitempty"`
	Error      string  `json:"error,omitempty"`
}

// RDWSPingResponse represents the response from GET /diagnostics/ping/
type RDWSPingResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result RDWSPingResult `json:"result"`
	} `json:"data"`
}

// RDWSTraceRouteHop represents a single hop in trace route
type RDWSTraceRouteHop struct {
	Hop     int     `json:"hop"`
	Address string  `json:"address,omitempty"`
	RTT     float64 `json:"rtt,omitempty"`
	Timeout bool    `json:"timeout,omitempty"`
}

// RDWSTraceRouteResult represents trace route diagnostic result
type RDWSTraceRouteResult struct {
	Success bool                `json:"success"`
	Host    string              `json:"host,omitempty"`
	Hops    []RDWSTraceRouteHop `json:"hops,omitempty"`
	Output  string              `json:"output,omitempty"`
	Error   string              `json:"error,omitempty"`
}

// RDWSTraceRouteResponse represents the response from GET /diagnostics/trace-route/
type RDWSTraceRouteResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result RDWSTraceRouteResult `json:"result"`
	} `json:"data"`
}

// RDWSNetworkConfig represents network interface configuration
type RDWSNetworkConfig struct {
	Interface   string   `json:"interface"`
	Type        string   `json:"type"` // dhcp, static
	IPAddress   string   `json:"ipAddress,omitempty"`
	Netmask     string   `json:"netmask,omitempty"`
	Gateway     string   `json:"gateway,omitempty"`
	DNS         []string `json:"dns,omitempty"`
	MACAddress  string   `json:"macAddress,omitempty"`
	LinkStatus  string   `json:"linkStatus,omitempty"`
}

// RDWSNetworkConfigResponse represents the response from GET /diagnostics/network-configuration/
type RDWSNetworkConfigResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result RDWSNetworkConfig `json:"result"`
	} `json:"data"`
}

// RDWSNetworkConfigSetRequest represents the request to set network configuration
type RDWSNetworkConfigSetRequest struct {
	Data struct {
		Type      string   `json:"type"` // dhcp, static
		IPAddress string   `json:"ipAddress,omitempty"`
		Netmask   string   `json:"netmask,omitempty"`
		Gateway   string   `json:"gateway,omitempty"`
		DNS       []string `json:"dns,omitempty"`
	} `json:"data"`
}

// RDWSNetworkConfigSetResponse represents the response from PUT /diagnostics/network-configuration/
type RDWSNetworkConfigSetResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result struct {
			Success bool   `json:"success"`
			Message string `json:"message,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// RDWSNetworkNeighbor represents a device on the network
type RDWSNetworkNeighbor struct {
	IPAddress  string `json:"ipAddress"`
	MACAddress string `json:"macAddress,omitempty"`
	Hostname   string `json:"hostname,omitempty"`
}

// RDWSNetworkNeighborhoodResult represents network neighborhood information
type RDWSNetworkNeighborhoodResult struct {
	Success   bool                  `json:"success"`
	Neighbors []RDWSNetworkNeighbor `json:"neighbors,omitempty"`
}

// RDWSNetworkNeighborhoodResponse represents the response from GET /diagnostics/network-neighborhood/
type RDWSNetworkNeighborhoodResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result RDWSNetworkNeighborhoodResult `json:"result"`
	} `json:"data"`
}

// RDWSPacketCaptureStatus represents packet capture status
type RDWSPacketCaptureStatus struct {
	Running     bool   `json:"running"`
	Interface   string `json:"interface,omitempty"`
	Duration    int    `json:"duration,omitempty"`
	FilePath    string `json:"filePath,omitempty"`
	StartTime   string `json:"startTime,omitempty"`
	ElapsedTime int    `json:"elapsedTime,omitempty"`
}

// RDWSPacketCaptureResponse represents the response from GET /diagnostics/packet-capture/
type RDWSPacketCaptureResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result RDWSPacketCaptureStatus `json:"result"`
	} `json:"data"`
}

// RDWSPacketCaptureStartRequest represents the request to start packet capture
type RDWSPacketCaptureStartRequest struct {
	Data struct {
		Interface string `json:"interface"`         // Network interface (e.g., "eth0")
		Duration  int    `json:"duration,omitempty"` // Duration in seconds
		Filter    string `json:"filter,omitempty"`   // tcpdump filter expression
	} `json:"data"`
}

// RDWSPacketCaptureStartResponse represents the response from POST /diagnostics/packet-capture/
type RDWSPacketCaptureStartResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result struct {
			Success  bool   `json:"success"`
			FilePath string `json:"filePath,omitempty"`
			Message  string `json:"message,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// RDWSPacketCaptureStopResponse represents the response from DELETE /diagnostics/packet-capture/
type RDWSPacketCaptureStopResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result struct {
			Success  bool   `json:"success"`
			FilePath string `json:"filePath,omitempty"`
			Message  string `json:"message,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// RDWSTelnetInfo represents telnet status information
type RDWSTelnetInfo struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port,omitempty"`
}

// RDWSTelnetResponse represents the response from GET /diagnostics/telnet/
type RDWSTelnetResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result RDWSTelnetInfo `json:"result"`
	} `json:"data"`
}

// RDWSTelnetSetRequest represents the request to enable/disable telnet
type RDWSTelnetSetRequest struct {
	Data struct {
		Enabled bool `json:"enabled"`
		Port    int  `json:"port,omitempty"`
	} `json:"data"`
}

// RDWSTelnetSetResponse represents the response from PUT /diagnostics/telnet/
type RDWSTelnetSetResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result struct {
			Success bool   `json:"success"`
			Message string `json:"message,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// RDWSSSHInfo represents SSH status information
type RDWSSSHInfo struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port,omitempty"`
}

// RDWSSSHResponse represents the response from GET /diagnostics/ssh/
type RDWSSSHResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result RDWSSSHInfo `json:"result"`
	} `json:"data"`
}

// RDWSSSHSetRequest represents the request to enable/disable SSH
type RDWSSSHSetRequest struct {
	Data struct {
		Enabled  bool   `json:"enabled"`
		Port     int    `json:"port,omitempty"`
		Password string `json:"password,omitempty"`
	} `json:"data"`
}

// RDWSSSHSetResponse represents the response from PUT /diagnostics/ssh/
type RDWSSSHSetResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result struct {
			Success bool   `json:"success"`
			Message string `json:"message,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// RDWSStorageReformatResponse represents the response from reformatting storage
type RDWSStorageReformatResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result struct {
			Success bool   `json:"success"`
			Message string `json:"message,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// RDWSCustomDataRequest represents a request to send custom data to a player
type RDWSCustomDataRequest struct {
	Data struct {
		Data string `json:"data"`
	} `json:"data"`
}

// RDWSCustomDataResponse represents the response from sending custom data
type RDWSCustomDataResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result struct {
			Success bool   `json:"success"`
			Message string `json:"message,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// RDWSFirmwareDownloadRequest represents a request to download firmware
type RDWSFirmwareDownloadRequest struct {
	Data struct {
		URL            string `json:"url"`
		AutoReboot     *bool  `json:"autoReboot,omitempty"`     // Optional: Auto-reboot after download (default: true)
	} `json:"data"`
}

// RDWSFirmwareDownloadResponse represents the response from firmware download request
type RDWSFirmwareDownloadResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result struct {
			Success bool   `json:"success"`
			Message string `json:"message,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// RDWSRegistry represents the full player registry
type RDWSRegistry struct {
	Sections map[string]map[string]string `json:"sections"`
}

// RDWSRegistryResponse represents the response from getting registry
type RDWSRegistryResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result RDWSRegistry `json:"result"`
	} `json:"data"`
}

// RDWSRegistryValue represents a single registry value
type RDWSRegistryValue struct {
	Section string `json:"section"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

// RDWSRegistryValueResponse represents the response from getting a registry value
type RDWSRegistryValueResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result struct {
			Value string `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// RDWSRegistrySetRequest represents a request to set a registry value
type RDWSRegistrySetRequest struct {
	Data struct {
		Value string `json:"value"`
	} `json:"data"`
}

// RDWSRegistrySetResponse represents the response from setting a registry value
type RDWSRegistrySetResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result struct {
			Success bool   `json:"success"`
			Message string `json:"message,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// RDWSRegistryDeleteResponse represents the response from deleting a registry value
type RDWSRegistryDeleteResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result struct {
			Success bool   `json:"success"`
			Message string `json:"message,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// RDWSRegistryFlushResponse represents the response from flushing registry
type RDWSRegistryFlushResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result struct {
			Success bool   `json:"success"`
			Message string `json:"message,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// RDWSRecoveryURL represents the recovery URL from registry
type RDWSRecoveryURL struct {
	URL string `json:"url"`
}

// RDWSRecoveryURLResponse represents the response from getting recovery URL
type RDWSRecoveryURLResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result struct {
			URL string `json:"url"`
		} `json:"result"`
	} `json:"data"`
}

// RDWSRecoveryURLSetRequest represents a request to set recovery URL
type RDWSRecoveryURLSetRequest struct {
	Data struct {
		URL string `json:"url"`
	} `json:"data"`
}

// RDWSRecoveryURLSetResponse represents the response from setting recovery URL
type RDWSRecoveryURLSetResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result struct {
			Success bool   `json:"success"`
			Message string `json:"message,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// RDWSLogFile represents a single log file
type RDWSLogFile struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Content string `json:"content,omitempty"`
}

// RDWSLogs represents the collection of log files from a player
type RDWSLogs struct {
	Files []RDWSLogFile `json:"files"`
}

// RDWSLogsResponse represents the response from getting logs
type RDWSLogsResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result json.RawMessage `json:"result"`
	} `json:"data"`
}

// RDWSLogsResult represents the successful result from getting logs
type RDWSLogsResult struct {
	Logs []RDWSLogFile `json:"logs"`
}

// RDWSCrashDumpFile represents a single crash dump file
type RDWSCrashDumpFile struct {
	Name      string `json:"name"`
	Timestamp string `json:"timestamp"`
	Size      int64  `json:"size"`
	Content   string `json:"content,omitempty"`
}

// RDWSCrashDump represents the collection of crash dump files from a player
type RDWSCrashDump struct {
	Files []RDWSCrashDumpFile `json:"files"`
}

// RDWSCrashDumpResponse represents the response from getting crash dumps
type RDWSCrashDumpResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Data   struct {
		Result json.RawMessage `json:"result"`
	} `json:"data"`
}

// RDWSCrashDumpResult represents the successful result from getting crash dumps
type RDWSCrashDumpResult struct {
	Dumps []RDWSCrashDumpFile `json:"dumps"`
}

// Subscription represents a device subscription in BSN.cloud
type Subscription struct {
	ID               int        `json:"id"`
	DeviceSerial     string     `json:"deviceSerial,omitempty"`
	DeviceID         int        `json:"deviceId,omitempty"`
	Type             string     `json:"type"`             // e.g., "Content", "Control", "Premium"
	Status           string     `json:"status"`           // e.g., "active", "expired", "pending"
	StartDate        time.Time  `json:"startDate"`
	EndDate          *time.Time `json:"endDate,omitempty"`
	CreationDate     time.Time  `json:"creationDate"`
	LastModifiedDate time.Time  `json:"lastModifiedDate"`
	AutoRenew        bool       `json:"autoRenew,omitempty"`
}

// SubscriptionList represents a list of subscriptions with pagination support
type SubscriptionList struct {
	Items       []Subscription `json:"items"`
	IsTruncated bool           `json:"isTruncated"`
	NextMarker  string         `json:"nextMarker,omitempty"`
	TotalCount  int            `json:"totalCount,omitempty"`
}

// SubscriptionCount represents the count of subscriptions
type SubscriptionCount struct {
	Count int `json:"count"`
}

// SubscriptionOperation represents an operation that can be performed on subscriptions
type SubscriptionOperation struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Allowed     bool   `json:"allowed"`
}

// SubscriptionOperations represents the list of available operations
type SubscriptionOperations struct {
	Operations []SubscriptionOperation `json:"operations"`
}

// ContentFile represents a content file in BSN.cloud
type ContentFile struct {
	ID                    int       `json:"id"`
	Name                  string    `json:"name"`
	Type                  string    `json:"type"`                   // "File" or "Folder"
	MediaType             string    `json:"mediaType,omitempty"`    // e.g., "Video", "Image", "Audio"
	FileSize              int64     `json:"fileSize,omitempty"`
	Hash                  string    `json:"hash,omitempty"`
	Path                  string    `json:"path,omitempty"`
	VirtualPath           string    `json:"virtualPath,omitempty"`
	FileExtension         string    `json:"fileExtension,omitempty"`
	MimeType              string    `json:"mimeType,omitempty"`
	CreationDate          time.Time `json:"creationDate"`
	LastModifiedDate      time.Time `json:"lastModifiedDate"`
	LastUploadedDate      *time.Time `json:"lastUploadedDate,omitempty"`
	UploadComplete        bool      `json:"uploadComplete"`
	ThumbnailURL          string    `json:"thumbnailUrl,omitempty"`
	ParentFolderID        *int      `json:"parentFolderId,omitempty"`
	IsShared              bool      `json:"isShared,omitempty"`
	StorageProvider       string    `json:"storageProvider,omitempty"`
	ExternalID            string    `json:"externalId,omitempty"`
}

// ContentFileList represents a paginated list of content files
type ContentFileList struct {
	Items       []ContentFile `json:"items"`
	IsTruncated bool          `json:"isTruncated"`
	NextMarker  string        `json:"nextMarker,omitempty"`
	TotalCount  int           `json:"totalCount,omitempty"`
}

// ContentFileCount represents the count of content files
type ContentFileCount struct {
	Count int `json:"count"`
}

// ContentDeleteResult represents the result of deleting content files
type ContentDeleteResult struct {
	DeletedCount int      `json:"deletedCount"`
	DeletedIDs   []int    `json:"deletedIds,omitempty"`
	Errors       []string `json:"errors,omitempty"`
}

// UploadRequest represents a file upload request (deprecated - use session-based upload)
type UploadRequest struct {
	FileName    string `json:"fileName"`
	FileSize    int64  `json:"fileSize"`
	VirtualPath string `json:"virtualPath,omitempty"`
	FileData    []byte `json:"-"` // Not serialized to JSON
}

// UploadResponse represents the response from a file upload
type UploadResponse struct {
	ContentID        int       `json:"contentId"`
	FileName         string    `json:"fileName"`
	FileSize         int64     `json:"fileSize"`
	VirtualPath      string    `json:"virtualPath"`
	MediaType        string    `json:"mediaType,omitempty"`
	UploadDate       time.Time `json:"uploadDate"`
	FileHash         string    `json:"fileHash,omitempty"`
	UploadComplete   bool      `json:"uploadComplete"`
}

// ContentUploadArguments represents arguments for creating a content upload
type ContentUploadArguments struct {
	FileName             string            `json:"FileName"`
	VirtualPath          string            `json:"VirtualPath,omitempty"`
	MediaType            string            `json:"MediaType,omitempty"` // Auto, Image, Video, Audio, etc.
	FileSize             int64             `json:"FileSize"`
	FileLastModifiedDate *time.Time        `json:"FileLastModifiedDate,omitempty"`
	SHA1Hash             string            `json:"SHA1Hash,omitempty"`
	Tags                 map[string]string `json:"Tags,omitempty"`
}

// UploadSessionResponse represents the response from creating an upload session
// This matches the WebPageUploadStatus schema from the OpenAPI spec
type UploadSessionResponse struct {
	SessionToken string                `json:"SessionToken"`
	Name         string                `json:"Name,omitempty"`
	Assets       []ContentUploadStatus `json:"Assets,omitempty"` // For WebPage uploads
	Uploads      []ContentUploadStatus `json:"Uploads,omitempty"` // Alternative field name

	// Embedded fields from ContentUploadStatus for standalone uploads
	UploadToken      string    `json:"UploadToken,omitempty"`
	ContentID        int       `json:"ContentId,omitempty"`
	FileName         string    `json:"FileName,omitempty"`
	FileSize         int64     `json:"FileSize,omitempty"`
	SHA1Hash         string    `json:"SHA1Hash,omitempty"`
	State            string    `json:"State,omitempty"`
	StartTime        time.Time `json:"StartTime,omitempty"`
	EndTime          time.Time `json:"EndTime,omitempty"`
	LastActivityTime time.Time `json:"LastActivityTime,omitempty"`
}

// ContentUploadStatus represents the status of a content upload
type ContentUploadStatus struct {
	SessionToken     string    `json:"SessionToken,omitempty"`
	UploadToken      string    `json:"UploadToken"`
	ContentID        int       `json:"ContentId,omitempty"`
	FileName         string    `json:"FileName"`
	FileSize         int64     `json:"FileSize"`
	SHA1Hash         string    `json:"SHA1Hash,omitempty"`
	State            string    `json:"State"` // Queued, Started, Uploading, Uploaded, Verified, Completed, etc.
	StartTime        time.Time `json:"StartTime,omitempty"`
	EndTime          time.Time `json:"EndTime,omitempty"`
	LastActivityTime time.Time `json:"LastActivityTime,omitempty"`
}

// ChunkUploadResponse represents the response from uploading a chunk
type ChunkUploadResponse struct {
	BytesReceived int64  `json:"BytesReceived,omitempty"`
	State         string `json:"State,omitempty"`
}

// =========================================================================
// Presentation Types
// =========================================================================

// Presentation represents a presentation in BSN.cloud
type Presentation struct {
	// Basic metadata
	ID                  int       `json:"id"`
	Name                string    `json:"name"`
	Type                string    `json:"type,omitempty"`                // e.g., "Standard", "HTML", "Complete"
	Description         string    `json:"description,omitempty"`
	LastModifiedDate    time.Time `json:"lastModifiedDate"`
	CreationDate        time.Time `json:"creationDate"`
	LastPublishedDate   *time.Time `json:"lastPublishedDate,omitempty"`
	PublishState        string    `json:"publishState,omitempty"`       // e.g., "Draft", "Published"
	ThumbnailURL        string    `json:"thumbnailUrl,omitempty"`
	ScheduleSettings    *ScheduleSettings `json:"scheduleSettings,omitempty"`
	IsSimplePlaylist    bool      `json:"isSimplePlaylist,omitempty"`
	ZoneCount           int       `json:"zoneCount,omitempty"`
	ContentItemCount    int       `json:"contentItemCount,omitempty"`
	Tags                []string  `json:"tags,omitempty"`

	// Presentation files and configuration
	ProjectFile           interface{}                     `json:"projectFile,omitempty"`
	AutoplayFile          interface{}                     `json:"autoplayFile,omitempty"`
	ResourcesFile         interface{}                     `json:"resourcesFile,omitempty"`
	UserDefinedEventsFile interface{}                     `json:"userDefinedEventsFile,omitempty"`
	ThumbnailFile         interface{}                     `json:"thumbnailFile,omitempty"`

	// Assets and dependencies
	Files                 []interface{}                   `json:"files,omitempty"`
	AutorunPlugins        interface{}                     `json:"autorunPlugins,omitempty"`
	Applications          []interface{}                   `json:"applications,omitempty"`
	Dependencies          []interface{}                   `json:"dependencies,omitempty"`
	Groups                interface{}                     `json:"groups,omitempty"`
	Permissions           []interface{}                   `json:"permissions,omitempty"`

	// Player configuration
	Autorun               *PresentationAutorun            `json:"autorun,omitempty"`
	DeviceWebPage         *PresentationDeviceWebPage      `json:"deviceWebPage,omitempty"`
	DeviceModel           string                          `json:"deviceModel,omitempty"`    // e.g., "HD1024"
	ScreenSettings        *PresentationScreenSettings     `json:"screenSettings,omitempty"`
	Language              string                          `json:"language,omitempty"`       // e.g., "English"
	Status                string                          `json:"status,omitempty"`         // e.g., "Draft", "Published"
}

// ScheduleSettings represents scheduling configuration for presentations
type ScheduleSettings struct {
	StartDate     *time.Time `json:"startDate,omitempty"`
	EndDate       *time.Time `json:"endDate,omitempty"`
	RecurrenceRule string    `json:"recurrenceRule,omitempty"`
}

// PresentationCount represents the count of presentations
type PresentationCount struct {
	Count int `json:"count"`
}

// PresentationList represents a paginated list of presentations
type PresentationList struct {
	Items       []Presentation `json:"items"`
	IsTruncated bool           `json:"isTruncated"`
	NextMarker  string         `json:"nextMarker,omitempty"`
	TotalCount  int            `json:"totalCount,omitempty"`
}

// PresentationCreateRequest represents a request to create a new presentation
type PresentationCreateRequest struct {
	ID                    int                             `json:"id"`                       // Must be 0 for new presentations
	Name                  string                          `json:"name"`
	CreationDate          string                          `json:"creationDate"`             // Use "0001-01-01T00:00:00" for new
	LastModifiedDate      string                          `json:"lastModifiedDate"`         // Use "0001-01-01T00:00:00" for new
	ProjectFile           interface{}                     `json:"projectFile"`              // null for new presentations
	Autorun               *PresentationAutorun            `json:"autorun,omitempty"`
	DeviceWebPage         *PresentationDeviceWebPage      `json:"deviceWebPage,omitempty"`
	DeviceModel           string                          `json:"deviceModel,omitempty"`    // e.g., "HD1024"
	ScreenSettings        *PresentationScreenSettings     `json:"screenSettings,omitempty"`
	Language              string                          `json:"language,omitempty"`       // e.g., "English"
	Status                string                          `json:"status,omitempty"`         // e.g., "Draft", "Published"
	AutoplayFile          interface{}                     `json:"autoplayFile"`             // null for new presentations
	ResourcesFile         interface{}                     `json:"resourcesFile"`            // null for new presentations
	UserDefinedEventsFile interface{}                     `json:"userDefinedEventsFile"`    // null for new presentations
	ThumbnailFile         interface{}                     `json:"thumbnailFile"`            // null for new presentations
	Files                 []interface{}                   `json:"files"`                    // [] for new presentations
	AutorunPlugins        interface{}                     `json:"autorunPlugins"`           // null for new presentations
	Applications          []interface{}                   `json:"applications"`             // [] for new presentations
	Dependencies          []interface{}                   `json:"dependencies"`             // [] for new presentations
	Groups                interface{}                     `json:"groups"`                   // null for new presentations
	Permissions           []interface{}                   `json:"permissions"`              // [] for new presentations
	Tags                  []string                        `json:"tags,omitempty"`
}

// PresentationAutorun represents autorun configuration for a presentation
type PresentationAutorun struct {
	Version  string `json:"version,omitempty"`  // e.g., "10.0.98"
	IsCustom bool   `json:"isCustom,omitempty"`
}

// PresentationDeviceWebPage represents device web page configuration
type PresentationDeviceWebPage struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// PresentationScreenSettings represents screen configuration for a presentation
type PresentationScreenSettings struct {
	VideoMode       string `json:"videoMode,omitempty"`       // e.g., "1920x1080x60p"
	Orientation     string `json:"orientation,omitempty"`     // e.g., "Landscape", "Portrait"
	Connector       string `json:"connector,omitempty"`       // e.g., "HDMI"
	BackgroundColor string `json:"backgroundColor,omitempty"` // e.g., "RGB:000000"
	Overscan        string `json:"overscan,omitempty"`        // e.g., "NoOverscan"
}

// PresentationDeleteResult represents the result of deleting presentations by filter
type PresentationDeleteResult struct {
	DeletedCount int      `json:"deletedCount"`
	DeletedIDs   []int    `json:"deletedIds,omitempty"`
	Errors       []string `json:"errors,omitempty"`
}

// =========================================================================
// Schedule Types
// =========================================================================

// ScheduledPresentation represents a scheduled presentation in a group
type ScheduledPresentation struct {
	ID                  int        `json:"id,omitempty"`
	PresentationID      int        `json:"presentationId"`
	PresentationName    string     `json:"presentationName,omitempty"`
	IsRecurrent         bool       `json:"isRecurrent"`
	EventDate           *time.Time `json:"eventDate,omitempty"`           // For non-recurrent schedules
	StartTime           string     `json:"startTime"`                     // TimeSpan format (HH:MM:SS)
	Duration            string     `json:"duration"`                      // TimeSpan format (HH:MM:SS)
	RecurrenceStartDate *time.Time `json:"recurrenceStartDate,omitempty"` // For recurrent schedules
	RecurrenceEndDate   *time.Time `json:"recurrenceEndDate,omitempty"`   // For recurrent schedules
	DaysOfWeek          any        `json:"daysOfWeek,omitempty"`          // Can be int or string ("EveryDay", "Sunday", etc.)
	InterruptScheduling bool       `json:"interruptScheduling,omitempty"`
	ExpirationDate      *time.Time `json:"expirationDate,omitempty"`
}

// ScheduledPresentationList represents a paginated list of scheduled presentations
type ScheduledPresentationList struct {
	Items       []ScheduledPresentation `json:"items"`
	IsTruncated bool                    `json:"isTruncated"`
	NextMarker  string                  `json:"nextMarker,omitempty"`
	TotalCount  int                     `json:"totalCount,omitempty"`
}

// =========================================================================
// Device Web Page Types
// =========================================================================

// DeviceWebPage represents a device web page template
type DeviceWebPage struct {
	ID               int        `json:"id"`
	Name             string     `json:"name"`
	Path             *string    `json:"path,omitempty"`
	Type             string     `json:"type,omitempty"`
	Size             *int64     `json:"size,omitempty"`
	Hash             *string    `json:"hash,omitempty"`
	CreationDate     *time.Time `json:"creationDate,omitempty"`
	LastModifiedDate *time.Time `json:"lastModifiedDate,omitempty"`
}

// DeviceWebPageList represents a paginated list of device web pages
type DeviceWebPageList struct {
	Items       []DeviceWebPage `json:"items"`
	IsTruncated bool            `json:"isTruncated"`
	NextMarker  string          `json:"nextMarker,omitempty"`
	TotalCount  int             `json:"totalCount,omitempty"`
}