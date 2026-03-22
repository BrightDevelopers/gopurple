# Design Document: BSN.cloud SDK Missing Functionality Implementation

## Architecture Overview

The implementation follows the established SDK patterns with three-tier architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                    Example Programs                          │
│  (10 new CLI tools in /gopurple/examples/*)                 │
└────────────────┬────────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────────┐
│                   SDK Service Layer                          │
│  PresentationService: Publish, AddZone, ConfigureZones      │
│  DeviceService: AssignPresentation, GetCurrentPresentation  │
│  GroupService: AssignPresentation                           │
│  ContentService: GetByID, Update                            │
└────────────────┬────────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────────┐
│              Internal Services Layer                         │
│  HTTP client, Auth manager, Error handling, Types           │
└─────────────────────────────────────────────────────────────┘
```

## Current State Analysis

The SDK has:
- Robust service pattern with authentication/network context management
- Type-safe interfaces in `/gopurple/internal/services/`
- Rich type definitions in `/gopurple/internal/types/types.go`
- Consistent example programs using shared helpers in `/gopurple/examples/shared/flags.go`
- Support for presentations (create, list, delete, update) but no zone management
- Support for devices and groups but no presentation assignment
- No publishing workflow

## Missing Functionality

1. **Presentation Zone Management**: Add video/image zones to presentations
2. **Presentation Publishing**: Publish presentations to make them distributable
3. **Presentation Assignment**: Assign presentations to groups and devices
4. **Playback Verification**: Get current presentation and distribution status from devices
5. **Content Metadata Management**: Update content file properties

## New Types Required

Add to `/gopurple/internal/types/types.go`:

```go
// PresentationZone represents a content zone in a presentation
type PresentationZone struct {
    ID          int                     `json:"id,omitempty"`
    Name        string                  `json:"name"`
    Type        string                  `json:"type"` // "VideoOrImages", "Audio", etc.
    X           int                     `json:"x"`
    Y           int                     `json:"y"`
    Width       int                     `json:"width"`
    Height      int                     `json:"height"`
    ZOrder      int                     `json:"zOrder"` // Layer/z-index
    Playlist    *ZonePlaylist           `json:"playlist,omitempty"`
    Properties  map[string]interface{}  `json:"properties,omitempty"`
}

// ZonePlaylist represents content items in a zone
type ZonePlaylist struct {
    Items       []PlaylistItem          `json:"items"`
    Transition  string                  `json:"transition,omitempty"`
}

// PlaylistItem represents a content item in a playlist
type PlaylistItem struct {
    ContentID   int     `json:"contentId"`
    ContentName string  `json:"contentName,omitempty"`
    Duration    int     `json:"duration,omitempty"` // seconds for images
    DisplayTime string  `json:"displayTime,omitempty"` // TimeSpan HH:MM:SS
}

// PresentationPublishRequest for publishing a presentation
type PresentationPublishRequest struct {
    Force bool `json:"force,omitempty"`
}

// PresentationPublishResponse from publishing
type PresentationPublishResponse struct {
    PresentationID   int       `json:"presentationId"`
    PublishState     string    `json:"publishState"`
    PublishedDate    time.Time `json:"publishedDate"`
    Warnings         []string  `json:"warnings,omitempty"`
}

// CurrentPresentationInfo represents current presentation on a device
type CurrentPresentationInfo struct {
    DeviceID         int       `json:"deviceId"`
    Serial           string    `json:"serial"`
    PresentationID   int       `json:"presentationId,omitempty"`
    PresentationName string    `json:"presentationName,omitempty"`
    Status           string    `json:"status"` // assigned, downloading, ready, playing
    AssignedDate     time.Time `json:"assignedDate,omitempty"`
    LastUpdateDate   time.Time `json:"lastUpdateDate,omitempty"`
}

// DistributionStatus for content distribution on a device
type DistributionStatus struct {
    DeviceID            int                     `json:"deviceId"`
    Serial              string                  `json:"serial"`
    PresentationID      int                     `json:"presentationId,omitempty"`
    ContentItems        []DistributionItem      `json:"contentItems"`
    TotalFiles          int                     `json:"totalFiles"`
    DownloadedFiles     int                     `json:"downloadedFiles"`
    PendingFiles        int                     `json:"pendingFiles"`
    FailedFiles         int                     `json:"failedFiles"`
    PercentComplete     float64                 `json:"percentComplete"`
    LastUpdateTime      time.Time               `json:"lastUpdateTime"`
    Errors              []DistributionError     `json:"errors,omitempty"`
}

// DistributionItem single content item distribution status
type DistributionItem struct {
    ContentID       int       `json:"contentId"`
    FileName        string    `json:"fileName"`
    FileSize        int64     `json:"fileSize"`
    BytesDownloaded int64     `json:"bytesDownloaded"`
    Status          string    `json:"status"` // pending, downloading, complete, error
    ErrorMessage    string    `json:"errorMessage,omitempty"`
}

// DistributionError represents a distribution error
type DistributionError struct {
    ErrorCode    string    `json:"errorCode"`
    ErrorMessage string    `json:"errorMessage"`
    ContentID    int       `json:"contentId,omitempty"`
    Timestamp    time.Time `json:"timestamp"`
}

// ContentUpdateRequest to update content file properties
type ContentUpdateRequest struct {
    Name        string            `json:"name,omitempty"`
    VirtualPath string            `json:"virtualPath,omitempty"`
    Tags        map[string]string `json:"tags,omitempty"`
    Description string            `json:"description,omitempty"`
}
```

## SDK Service Method Signatures

### PresentationService Extensions

Add to `/gopurple/internal/services/presentations.go`:

```go
// Publish publishes a presentation to make it available for distribution
Publish(ctx context.Context, presentationID int, force bool) (*types.PresentationPublishResponse, error)

// AddVideoZone adds a video zone to a presentation
AddVideoZone(ctx context.Context, presentationID int, zone *types.PresentationZone) (*types.Presentation, error)

// AddImageZone adds an image zone with duration to a presentation
AddImageZone(ctx context.Context, presentationID int, zone *types.PresentationZone, duration int) (*types.Presentation, error)

// ConfigureZones updates all zones for a presentation (replaces existing)
ConfigureZones(ctx context.Context, presentationID int, zones []types.PresentationZone) (*types.Presentation, error)

// GetZones retrieves all zones for a presentation
GetZones(ctx context.Context, presentationID int) ([]types.PresentationZone, error)
```

**Implementation Strategy**: GET current presentation, modify zones, PUT back to API.

### DeviceService Extensions

Add to `/gopurple/internal/services/devices.go`:

```go
// AssignPresentation assigns a presentation to a specific device
AssignPresentation(ctx context.Context, deviceID int, presentationID int) error

// AssignPresentationBySerial assigns by serial number
AssignPresentationBySerial(ctx context.Context, serial string, presentationID int) error

// GetCurrentPresentation retrieves currently assigned presentation
GetCurrentPresentation(ctx context.Context, deviceID int) (*types.CurrentPresentationInfo, error)

// GetCurrentPresentationBySerial retrieves by serial number
GetCurrentPresentationBySerial(ctx context.Context, serial string) (*types.CurrentPresentationInfo, error)

// GetDistributionStatus retrieves content distribution status
GetDistributionStatus(ctx context.Context, deviceID int) (*types.DistributionStatus, error)

// GetDistributionStatusBySerial retrieves by serial number
GetDistributionStatusBySerial(ctx context.Context, serial string) (*types.DistributionStatus, error)
```

### GroupService Extensions

Create `/gopurple/internal/services/groups.go`:

```go
// GroupService provides group management operations
type GroupService interface {
    // AssignPresentation assigns a presentation to all devices in a group
    AssignPresentation(ctx context.Context, groupID int, presentationID int) error

    // AssignPresentationByName assigns by group name
    AssignPresentationByName(ctx context.Context, groupName string, presentationID int) error
}
```

### ContentService Extensions

Add to `/gopurple/internal/services/content.go`:

```go
// Update modifies content file properties (name, virtual path, tags)
Update(ctx context.Context, contentID int, request *types.ContentUpdateRequest) (*types.ContentFile, error)
```

## Example Program Specifications

All 10 programs follow established patterns:
- Use `/gopurple/examples/shared/flags.go` helpers
- Support `--help`, `--json`, `--verbose`, `--network`
- Support environment variables for IDs
- Clear error messages

### 1. main-presentation-add-video-zone

**Flags**: `--id`, `--content-id`, `--name`, `--x`, `--y`, `--width`, `--height`, `--layer`

Adds video zone to presentation using AddVideoZone SDK method.

### 2. main-presentation-add-image-zone

**Flags**: Same as video zone plus `--duration`, `--transition`

Adds image zone with duration using AddImageZone SDK method.

### 3. main-presentation-configure-zones

**Flags**: `--id`, `--zones-json`

Reads zone configuration from JSON file and applies using ConfigureZones SDK method.

### 4. main-presentation-publish

**Flags**: `--id`, `--force`

Publishes presentation using Publish SDK method.

### 5. main-group-assign-presentation

**Flags**: `--group-id`, `--presentation-id`

Assigns presentation to group using GroupService.AssignPresentation.

### 6. main-device-assign-presentation

**Flags**: `--serial` OR `--id`, `--presentation-id`

Assigns presentation to device using DeviceService.AssignPresentation.

### 7. main-device-current-presentation

**Flags**: `--serial` OR `--id`

Gets current presentation using GetCurrentPresentation SDK method.

### 8. main-device-distribution-status

**Flags**: `--serial` OR `--id`

Gets distribution status using GetDistributionStatus SDK method.

### 9. main-content-info

**Flags**: `--id`

Gets content metadata using existing GetByID SDK method.

### 10. main-content-update

**Flags**: `--id`, `--name`, `--virtual-path`, `--tags`, `--description`

Updates content metadata using Update SDK method.

## Implementation Phases

### Phase 1: Type Definitions and Basic Structure
- Add new types to `/gopurple/internal/types/types.go`
- Update service interfaces
- Write unit tests for type marshaling
- **Validation**: All types compile, tests pass

### Phase 2: Presentation Zone Management
- Implement AddVideoZone, AddImageZone, ConfigureZones
- Create 3 example programs
- Write integration tests
- **Validation**: Can add zones to presentation via API

### Phase 3: Presentation Publishing
- Implement Publish method
- Create main-presentation-publish
- Write integration tests
- **Validation**: Can publish presentations

### Phase 4: Presentation Assignment
- Create GroupService interface
- Implement assignment methods for groups/devices
- Create 2 example programs
- Write integration tests
- **Validation**: Can assign presentations

### Phase 5: Playback Verification
- Implement GetCurrentPresentation, GetDistributionStatus
- Create 2 example programs
- Write integration tests
- **Validation**: Can retrieve status from devices

### Phase 6: Content Metadata Management
- Implement Update method
- Create 2 example programs
- Write integration tests
- **Validation**: Can update content properties

### Phase 7: Documentation and Examples
- Update README.md
- Create workflow guide
- Update examples/README.md
- Create end-to-end deployment script

## Data Flow Examples

### Zone Addition Flow
```
User → main-presentation-add-video-zone
    → PresentationService.AddVideoZone(ctx, presentationID, zone)
        → GetByID to fetch current presentation
        → Append zone to presentation.zones
        → Update to save modified presentation
        → Return updated Presentation
    → Display confirmation
```

### Publishing Flow
```
User → main-presentation-publish
    → PresentationService.Publish(ctx, presentationID, force)
        → POST /2024-10/REST/Presentations/{id}/publish
        → Return PresentationPublishResponse
    → Display publish confirmation
```

### Assignment Flow
```
User → main-device-assign-presentation
    → DeviceService.AssignPresentation(ctx, deviceID, presentationID)
        → GetByID to fetch device
        → Modify device assignment field
        → Update to save modified device
        → Return success
    → Display assignment confirmation
```

## Error Handling Strategy

- **Validation Errors (4xx)**: Clear user-facing messages
- **Authentication Errors (401/403)**: Re-auth or check permissions
- **Not Found (404)**: Verify IDs, list available resources
- **Conflicts (409)**: Explain conflict and resolution
- **Server Errors (5xx)**: Retry with exponential backoff

## Testing Approach

- **Unit Tests**: Type marshaling, service validation, error paths
- **Integration Tests**: Complete workflows with BSN.cloud test account
- **Example Tests**: Flag parsing, JSON output, error messages

**Coverage Targets**: 80% SDK methods, 70% example programs

## Critical Implementation Files

1. `/gopurple/internal/types/types.go` - All new type definitions
2. `/gopurple/internal/services/presentations.go` - Zone and publish methods
3. `/gopurple/internal/services/devices.go` - Assignment and status methods
4. `/gopurple/internal/services/groups.go` - New group service
5. `/gopurple/examples/shared/flags.go` - Shared helper patterns
