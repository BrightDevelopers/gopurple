# Missing Functionality for Complete Workflows

This document identifies functionality that is currently missing from the SDK examples and would be needed to support complete end-to-end workflows.

## Priority 1: Presentation Content Management (CRITICAL)

### Problem

The SDK can create empty presentations but cannot add content zones to them. This makes presentations unusable since they have no content to display.

### Missing Examples

#### 1. main-presentation-add-video-zone

Add a video zone to a presentation.

**Flags:**
```bash
--id <presentation-id>          # Presentation ID (required)
--content-id <content-id>       # Video content ID (required)
--name <zone-name>              # Zone name (default: "Video Zone")
--x <x>                         # X position (default: 0)
--y <y>                         # Y position (default: 0)
--width <width>                 # Width (default: 1920)
--height <height>               # Height (default: 1080)
--layer <layer>                 # Layer/z-index (default: 0)
```

**API Endpoint:**
- `GET /2024-10/REST/presentation/{id}` - Get current presentation
- `PUT /2024-10/REST/presentation/{id}` - Update with new zone

#### 2. main-presentation-add-image-zone

Add an image zone with duration to a presentation.

**Flags:**
```bash
--id <presentation-id>          # Presentation ID (required)
--content-id <content-id>       # Image content ID (required)
--duration <seconds>            # Display duration (required)
--name <zone-name>              # Zone name (default: "Image Zone")
--x <x>                         # X position (default: 0)
--y <y>                         # Y position (default: 0)
--width <width>                 # Width (default: 1920)
--height <height>               # Height (default: 1080)
--layer <layer>                 # Layer/z-index (default: 0)
--transition <type>             # Transition effect (optional)
```

**API Endpoint:**
- `GET /2024-10/REST/presentation/{id}` - Get current presentation
- `PUT /2024-10/REST/presentation/{id}` - Update with new zone

#### 3. main-presentation-configure-zones

Update zone configuration for an existing presentation.

**Flags:**
```bash
--id <presentation-id>          # Presentation ID (required)
--zones-json <file>             # JSON file with zones configuration
```

**API Endpoint:**
- `PUT /2024-10/REST/presentation/{id}` - Update presentation

## Priority 2: Presentation Publishing and Distribution (CRITICAL)

### Problem

Created presentations cannot be published or assigned to devices, making them inaccessible to players.

### Missing Examples

#### 4. main-presentation-publish

Publish a presentation to make it available for distribution.

**Flags:**
```bash
--id <presentation-id>          # Presentation ID (required)
--network <name>                # Network name
--verbose                       # Show detailed information
```

**API Endpoint:**
- `POST /2024-10/REST/presentation/{id}/publish`

**Environment Variables:**
- `BS_PRESENTATION_ID` - Default presentation ID

#### 5. main-group-assign-presentation

Assign a presentation to a group of devices.

**Flags:**
```bash
--group-id <id>                 # Group ID (required)
--presentation-id <id>          # Presentation ID (required)
--network <name>                # Network name
--verbose                       # Show detailed information
```

**API Endpoint:**
- `PUT /2024-10/REST/group/{id}` - Update group with presentation

**Environment Variables:**
- `BS_GROUP_ID` - Default group ID
- `BS_PRESENTATION_ID` - Default presentation ID

#### 6. main-device-assign-presentation

Assign a presentation to a specific device.

**Flags:**
```bash
--serial <serial>               # Device serial (or use --id)
--id <device-id>                # Device ID (alternative)
--presentation-id <id>          # Presentation ID (required)
--network <name>                # Network name
--verbose                       # Show detailed information
```

**API Endpoint:**
- `PUT /2024-10/REST/device/{id}` - Update device with presentation

**Environment Variables:**
- `BS_SERIAL` - Default device serial
- `BS_DEVICE_ID` - Default device ID
- `BS_PRESENTATION_ID` - Default presentation ID

## Priority 3: Playback Verification (HIGH)

### Problem

No way to verify what presentation is currently playing on a device or check distribution status.

### Missing Examples

#### 7. main-device-current-presentation

Get the currently assigned/playing presentation for a device.

**Flags:**
```bash
--serial <serial>               # Device serial (or use --id)
--id <device-id>                # Device ID (alternative)
--network <name>                # Network name
--json                          # Output as JSON
--verbose                       # Show detailed information
```

**API Endpoint:**
- `GET /2024-10/REST/device/{id}` - Get device with presentation info

**Output:**
- Current presentation ID and name
- Presentation status (assigned, downloading, ready, playing)
- Last update time

**Environment Variables:**
- `BS_SERIAL` - Default device serial
- `BS_DEVICE_ID` - Default device ID

#### 8. main-device-distribution-status

Check content distribution status for a device.

**Flags:**
```bash
--serial <serial>               # Device serial (or use --id)
--id <device-id>                # Device ID (alternative)
--network <name>                # Network name
--json                          # Output as JSON
```

**API Endpoint:**
- `GET /2024-10/REST/device/{id}/status` - Get distribution status

**Output:**
- Content download progress
- Files pending download
- Distribution errors

## Priority 4: Schedule Management (MEDIUM)

### Problem

While schedule examples exist, there may be gaps in complete schedule management workflows.

### Review Needed

Check if the existing schedule examples (`main-schedule-*`) provide complete coverage for:
- Creating schedules with presentation assignments
- Updating schedule presentation mappings
- Deleting schedules
- Listing schedules by group/device

## Priority 5: Advanced Content Management (LOW)

### Missing Examples

#### 9. main-content-info

Get detailed information about a content file.

**Flags:**
```bash
--id <content-id>               # Content ID (required)
--network <name>                # Network name
--json                          # Output as JSON
```

**Environment Variables:**
- `BS_CONTENT_ID` - Default content ID

#### 10. main-content-update

Update content file properties (name, virtual path, tags).

**Flags:**
```bash
--id <content-id>               # Content ID (required)
--name <name>                   # New name
--virtual-path <path>           # New virtual path
--tags <tags>                   # Comma-separated tags
--network <name>                # Network name
```

## Implementation Guide

### Step 1: Create SDK Service Methods

Before creating examples, implement the underlying SDK methods in the appropriate service files:

**presentations.go:**
```go
func (s *PresentationService) AddVideoZone(ctx context.Context, presentationID int, zone *VideoZone) error
func (s *PresentationService) AddImageZone(ctx context.Context, presentationID int, zone *ImageZone) error
func (s *PresentationService) Publish(ctx context.Context, presentationID int) error
```

**groups.go:**
```go
func (s *GroupService) AssignPresentation(ctx context.Context, groupID int, presentationID int) error
```

**devices.go:**
```go
func (s *DeviceService) AssignPresentation(ctx context.Context, deviceID int, presentationID int) error
func (s *DeviceService) GetCurrentPresentation(ctx context.Context, deviceID int) (*CurrentPresentation, error)
func (s *DeviceService) GetDistributionStatus(ctx context.Context, deviceID int) (*DistributionStatus, error)
```

### Step 2: Create Example Programs

Follow the pattern established in existing examples:
1. Use `examples/shared` package for flag helpers
2. Support environment variables for IDs
3. Provide `--json` output option
4. Include `--verbose` for detailed information
5. Handle errors gracefully with clear messages
6. Include usage examples in `--help`

### Step 3: Add Tests

Each new example should have:
1. Unit tests for the SDK methods
2. Integration tests (if possible)
3. Documentation in `examples/README.md`

### Step 4: Update Documentation

1. Add examples to `examples/README.md`
2. Update `docs/test-content-examples.md`
3. Create workflow guides showing complete end-to-end examples
4. Update `examples/scripts/deploy-presentation.sh` to use new functionality

## Workarounds

Until these examples are implemented:

### Manual API Calls

Use `curl` with authentication token from `main-auth-info`:

```bash
# Get token
./bin/main-auth-info --show
TOKEN=$(cat .token)

# Add zones to presentation (simplified example)
curl -X PUT "https://api.bsn.cloud/2024-10/REST/presentation/$PRES_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d @zones-config.json

# Publish presentation
curl -X POST "https://api.bsn.cloud/2024-10/REST/presentation/$PRES_ID/publish" \
  -H "Authorization: Bearer $TOKEN"
```

### BrightAuthor:connected

Use the official GUI tool for:
- Adding content zones to presentations
- Publishing presentations
- Assigning presentations to groups/devices
- Monitoring distribution status

## Impact Assessment

**Without Priority 1 & 2 examples:**
- Cannot create usable presentations programmatically
- Cannot automate content deployment workflows
- Must use GUI tools for all presentation management

**With Priority 1 & 2 implemented:**
- Complete automation of presentation deployment
- CI/CD integration possible
- Scripted content updates
- Batch presentation creation

**With Priority 3 added:**
- Closed-loop deployment verification
- Automated testing of deployments
- Monitoring and alerting on distribution issues

## Estimated Implementation Effort

- **Priority 1** (3 examples): 2-3 days per example = 6-9 days
- **Priority 2** (3 examples): 1-2 days per example = 3-6 days
- **Priority 3** (2 examples): 1-2 days per example = 2-4 days

**Total:** 11-19 days for complete implementation

This includes:
- SDK method implementation
- Example program creation
- Testing
- Documentation
