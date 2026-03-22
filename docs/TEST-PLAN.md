# Test Plan: BSN.cloud SDK Missing Functionality

## Test Strategy Overview

Phased testing approach aligned with implementation priorities:
- **Phase 1**: Presentation Content Management - zone addition
- **Phase 2**: Publishing and Distribution - publish and assignment
- **Phase 3**: Playback Verification - status tracking
- **Phase 4**: Content Management - metadata operations

Each phase includes:
- Unit tests for SDK service methods
- Integration tests for example programs
- Flag parsing and environment variable tests
- Mock API response handling
- Error condition testing

## Testing Patterns

From existing codebase:
- **Unit Tests**: Validate parameters, auth, network context, errors (no real API calls)
- **Structure Tests**: Verify type definitions and data integrity
- **Flag Tests**: Table-driven tests for environment variable fallback
- **Interface Tests**: Ensure implementations satisfy interfaces
- **Error Tests**: Validate error classification

## Phase 1: Presentation Content Management Tests

### AddVideoZone SDK Method
**Location**: `/gopurple/internal/services/presentations_test.go`

**Tests**:
```go
TestPresentationService_AddVideoZone
- Validates presentationID > 0
- Validates zone parameter not nil
- Validates content ID > 0
- Requires authentication
- Requires network context
- Tests minimal parameters
- Tests all optional parameters (x, y, width, height, layer)
```

**Error Scenarios**:
- Invalid presentation ID (0, negative)
- Nil zone parameter
- Invalid content ID
- Missing authentication
- Network not set
- 404 presentation not found
- 400 invalid zone configuration

### AddImageZone SDK Method
**Location**: `/gopurple/internal/services/presentations_test.go`

**Tests**:
```go
TestPresentationService_AddImageZone
- Validates presentationID > 0
- Validates zone parameter not nil
- Validates content ID > 0
- Validates duration > 0
- Requires authentication
- Requires network context
- Tests transition effects
- Tests default values
```

**Error Scenarios**:
- Invalid duration (0, negative)
- Invalid transition type
- Presentation already published (immutable)

### main-presentation-add-video-zone Program
**Location**: `/gopurple/examples/main-presentation-add-video-zone/main_test.go`

**Tests**:
```go
TestMainPresentationAddVideoZone_FlagParsing
- Validates --id required
- Validates --content-id required
- Tests default values
- Tests all flags override defaults

TestVideoZoneFlags
- Flag precedence over defaults
- BS_PRESENTATION_ID fallback
- BS_CONTENT_ID fallback
- Invalid values rejected
```

**Execution Tests**:
```bash
# Valid
./main-presentation-add-video-zone --id 123 --content-id 456
./main-presentation-add-video-zone --id 123 --content-id 456 --name "Main" --layer 1

# Errors
./main-presentation-add-video-zone  # Missing flags
./main-presentation-add-video-zone --id 0 --content-id 456  # Invalid ID
```

### main-presentation-add-image-zone Program

**Tests**: Similar to video zone plus duration validation

### main-presentation-configure-zones Program

**Tests**:
```go
TestMainPresentationConfigureZones_FlagParsing
- Validates --id required
- Validates --zones-json required
- Tests file reading and parsing

TestZonesJSONParsing
- Valid zones configuration
- Invalid JSON rejected
- Missing required fields rejected
- Unknown zone types handled
```

**Sample zones.json**:
```json
{
  "zones": [
    {
      "type": "video",
      "name": "Main Video",
      "contentId": 123,
      "x": 0,
      "y": 0,
      "width": 1920,
      "height": 1080,
      "layer": 0
    }
  ]
}
```

## Phase 2: Publishing and Distribution Tests

### Publish SDK Method
**Location**: `/gopurple/internal/services/presentations_test.go`

**Tests**:
```go
TestPresentationService_Publish
- Validates presentationID > 0
- Requires authentication
- Requires network context
- Tests successful publish
- Tests already published presentation
```

**Error Scenarios**:
- Invalid presentation ID
- Presentation not found
- Presentation has no content zones
- Already published

### GroupService.AssignPresentation
**Location**: `/gopurple/internal/services/groups_test.go` (new file)

**Tests**:
```go
TestGroupService_AssignPresentation
- Validates groupID > 0
- Validates presentationID > 0
- Requires authentication
- Requires network context
- Tests successful assignment
- Tests overwriting existing assignment
```

**Error Scenarios**:
- Invalid group ID
- Invalid presentation ID
- Group not found
- Presentation not published
- Incompatible with group devices

### DeviceService.AssignPresentation
**Location**: `/gopurple/internal/services/devices_test.go`

**Tests**:
```go
TestDeviceService_AssignPresentation
- Validates deviceID > 0
- Validates presentationID > 0
- Requires authentication
- Requires network context
- Tests by device ID
- Tests by device serial
```

**Error Scenarios**:
- Invalid device ID
- Device not found
- Presentation not published
- Incompatible with device model

### Example Programs

**main-presentation-publish**:
```go
TestMainPresentationPublish_FlagParsing
- Validates --id required
- Tests BS_PRESENTATION_ID
- Tests --verbose output
```

**main-group-assign-presentation**:
```go
TestMainGroupAssignPresentation_FlagParsing
- Validates --group-id required
- Validates --presentation-id required
- Tests BS_GROUP_ID, BS_PRESENTATION_ID
```

**main-device-assign-presentation**:
```go
TestMainDeviceAssignPresentation_FlagParsing
- Validates --serial OR --id (not both)
- Validates --presentation-id required
- Tests BS_SERIAL, BS_DEVICE_ID, BS_PRESENTATION_ID
```

## Phase 3: Playback Verification Tests

### GetCurrentPresentation SDK Method
**Location**: `/gopurple/internal/services/devices_test.go`

**Tests**:
```go
TestDeviceService_GetCurrentPresentation
- Validates deviceID > 0
- Requires authentication
- Requires network context
- Tests device with assigned presentation
- Tests device with no presentation
```

**Mock Response**:
```go
type CurrentPresentation struct {
    ID               int
    Name             string
    Status           string  // assigned, downloading, ready, playing
    LastUpdateTime   time.Time
    DownloadProgress int
}
```

**Error Scenarios**:
- Device not found
- Device offline (no status)
- Network failure

### GetDistributionStatus SDK Method
**Location**: `/gopurple/internal/services/devices_test.go`

**Tests**:
```go
TestDeviceService_GetDistributionStatus
- Validates deviceID > 0
- Requires authentication
- Requires network context
- Tests active downloads
- Tests completed distribution
- Tests distribution errors
```

**Mock Response**:
```go
type DistributionStatus struct {
    PendingFiles     []PendingFile
    DownloadProgress int
    Errors           []DistributionError
    LastSyncTime     time.Time
}
```

### Example Programs

**main-device-current-presentation**:
```go
TestMainDeviceCurrentPresentation_FlagParsing
- Validates --serial OR --id required
- Tests --json output format
- Tests environment variables

TestCurrentPresentationOutput
- Tests formatted text output
- Tests JSON output structure
- Tests --verbose details
```

**main-device-distribution-status**:
```go
TestMainDeviceDistributionStatus_FlagParsing
- Validates device identifier required
- Tests --json output

TestDistributionStatusOutput
- Tests pending downloads display
- Tests completion percentage
- Tests error display
- Tests "no content" state
```

## Phase 4: Content Management Tests

### ContentService.GetInfo
**Location**: `/gopurple/internal/services/content_test.go`

**Tests**:
```go
TestContentService_GetInfo
- Validates contentID > 0
- Requires authentication
- Requires network context
- Tests retrieving metadata
```

### ContentService.Update
**Location**: `/gopurple/internal/services/content_test.go`

**Tests**:
```go
TestContentService_Update
- Validates contentID > 0
- Requires authentication
- Requires network context
- Tests updating name
- Tests updating virtual path
- Tests updating tags
- Tests partial updates
```

**Error Scenarios**:
- Content not found
- Name conflicts
- Invalid virtual path format
- Content in use (cannot modify)

### Example Programs

**main-content-info**:
```go
TestMainContentInfo_FlagParsing
- Validates --id required
- Tests BS_CONTENT_ID
- Tests --json output
```

**main-content-update**:
```go
TestMainContentUpdate_FlagParsing
- Validates --id required
- Tests at least one update field required
- Tests updating individual fields
- Tests updating multiple fields
```

## Test Fixtures and Mock Data

### Mock Presentation with Zones
```json
{
  "id": 123,
  "name": "Test Presentation",
  "publishState": "Published",
  "zones": [
    {
      "id": 1,
      "type": "video",
      "contentId": 456,
      "name": "Main Video",
      "x": 0,
      "y": 0,
      "width": 1920,
      "height": 1080,
      "layer": 0
    }
  ]
}
```

### Mock Device with Presentation
```json
{
  "id": 789,
  "serial": "ABC123DEF456",
  "model": "XD1033",
  "currentPresentation": {
    "id": 123,
    "name": "Test Presentation",
    "status": "playing"
  }
}
```

## Environment Variable Test Matrix

| Variable | Required | Fallback Behavior |
|----------|----------|-------------------|
| BS_CLIENT_ID | Yes | Error if missing |
| BS_SECRET | Yes | Error if missing |
| BS_NETWORK | No | Interactive selection |
| BS_SERIAL | No | Flag precedence |
| BS_DEVICE_ID | No | Flag precedence |
| BS_PRESENTATION_ID | No | Flag precedence |
| BS_CONTENT_ID | No | Flag precedence |
| BS_GROUP_ID | No | Flag precedence |

## Test Execution Order

**Phase 1**: SDK methods → flag helpers → example parsing → integration

**Phase 2**: Publish method → group service → device assignment → examples

**Phase 3**: Device status methods → examples

**Phase 4**: Content methods → examples

## Coverage Targets

- **SDK Service Methods**: Minimum 80% code coverage
- **Example Programs**: Minimum 70% code coverage
- **Integration Tests**: End-to-end workflows must complete successfully

## Validation Criteria

**Phase 1 Complete**: Can create presentation and add zones via API

**Phase 2 Complete**: Can publish and assign presentations

**Phase 3 Complete**: Can retrieve status correctly

**Phase 4 Complete**: Can view and modify content metadata

## Edge Cases

### Zone Management
- Empty presentation (first zone)
- Zone overlap (should succeed)
- Negative coordinates (should fail)
- Oversized zones (should warn)
- Missing/deleted content references

### Publishing
- No zones defined (should fail)
- Already published (idempotent)
- Draft with validation errors

### Assignment
- Incompatible model (4K to HD)
- Offline device (should queue)
- Mixed model groups
- Overwrite existing

### Status
- Never-online device
- Partial download
- Download failure
- Disk full

## Test Data Requirements

**Minimum**:
- 2 test presentations (simple, complex)
- 3 content items (1 video, 2 images)
- 1 device group with 2+ devices
- 2 test devices (different models)

**Credentials**: Test API credentials with limited permissions, separate test network

## Mock HTTP Server

For integration tests without live API:

```go
type MockBSNServer struct {
    presentations map[int]*types.Presentation
    devices       map[int]*types.Device
    groups        map[int]*types.Group
    content       map[int]*types.Content
}
```
