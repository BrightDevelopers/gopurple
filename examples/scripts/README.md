# Example Scripts

This directory contains shell scripts that orchestrate multiple example programs to accomplish complete workflows.

## deploy-presentation.sh

Automated presentation deployment workflow that uploads content and creates presentations.

### What It Does

The script performs the following steps:

1. **Clean up existing content** - Deletes old versions of video/image files if they exist
2. **Upload content** - Uploads video and image files to BSN.cloud
3. **Clean up existing presentation** - Deletes old presentation with the same name
4. **Create presentation** - Creates a new empty presentation

### Usage

```bash
# Make script executable
chmod +x deploy-presentation.sh

# Run with default config (presentation-config.json)
./deploy-presentation.sh

# Run with custom config
./deploy-presentation.sh my-config.json
```

### Configuration File

Create a JSON configuration file with the following format:

```json
{
  "presentation_name": "Store Display Demo",
  "video_file": "videos/promo.mp4",
  "image_file": "images/logo.png",
  "image_duration": 10,
  "network": "Production",
  "device_serial": "BS123456789"
}
```

**Fields:**
- `presentation_name` - Name for the presentation (required)
- `video_file` - Path to video file (required)
- `image_file` - Path to image file (required)
- `image_duration` - Duration in seconds for image display (required)
- `network` - BSN.cloud network name (optional, uses BS_NETWORK if not specified)
- `device_serial` - Target device serial (optional, for future use)

### Prerequisites

**Required:**
- jq (JSON processor) - `sudo dnf install jq` or `brew install jq`
- Built example programs in `../../bin/`
- Valid BSN.cloud credentials in environment:
  ```bash
  export BS_CLIENT_ID=your_client_id
  export BS_SECRET=your_client_secret
  ```

**Build examples:**
```bash
cd /gopurple
make build-examples
```

### Example Output

```
=== BSN.cloud Presentation Deployment Script ===

Reading configuration from presentation-config.json
  Presentation: Store Display Demo
  Video: videos/promo.mp4
  Image: images/logo.png
  Image Duration: 10s
  Network: Production

=== Step 1: Clean up existing content ===
Checking for existing content files...
  Deleting existing video: promo.mp4
  Deleting existing image: logo.png

=== Step 2: Upload content files ===
Uploading video: videos/promo.mp4
  ✓ Video uploaded (ID: 12345)
Uploading image: images/logo.png
  ✓ Image uploaded (ID: 67890)

=== Step 3: Clean up existing presentation ===
  Deleting existing presentation: Store Display Demo

=== Step 4: Create presentation ===
Creating presentation: Store Display Demo
  ✓ Presentation created (ID: 99999)

=== MISSING FUNCTIONALITY ===
⚠️  The following steps CANNOT be completed with existing example programs:

MISSING: Add content to presentation
  - Need: Example program to add video zone to presentation
  - Need: Example program to add image zone with duration
  - API: PUT /2024-10/presentation/{id} with zones configuration

MISSING: Publish/distribute presentation
  - Need: Example program to publish presentation
  - Need: Example program to assign presentation to device/group
  - API: POST /2024-10/presentation/{id}/publish

MISSING: Verify playback
  - Need: Example program to check current presentation on device
  - API: GET /2024-10/device/{id}/status with presentation info
```

## Missing Functionality

The script clearly identifies missing functionality that would be needed to complete the full workflow:

### 1. Add Content to Presentation

**What's needed:**
- Example program to add video zones to a presentation
- Example program to add image zones with duration settings
- Configure zone layout and properties

**API endpoints:**
- `PUT /2024-10/REST/presentation/{id}` - Update presentation with zones

**Suggested example programs:**
- `examples/main-presentation-add-video-zone`
- `examples/main-presentation-add-image-zone`
- `examples/main-presentation-configure-zones`

### 2. Publish and Distribute Presentation

**What's needed:**
- Publish the presentation (make it available for deployment)
- Assign presentation to device or group
- Trigger content distribution

**API endpoints:**
- `POST /2024-10/REST/presentation/{id}/publish`
- `PUT /2024-10/REST/group/{id}` - Update group with presentation assignment
- `PUT /2024-10/REST/device/{id}` - Update device with presentation

**Suggested example programs:**
- `examples/main-presentation-publish`
- `examples/main-group-assign-presentation`
- `examples/main-device-assign-presentation`

### 3. Verify Playback

**What's needed:**
- Check what presentation is currently playing on a device
- Verify content distribution status
- Monitor playback status

**API endpoints:**
- `GET /2024-10/REST/device/{id}` - Get device with current presentation info
- `GET /2024-10/REST/device/{id}/status` - Device status with playback info

**Suggested example programs:**
- `examples/main-device-current-presentation`
- `examples/main-device-playback-status`

## Implementation Roadmap

To complete the full presentation deployment workflow, implement these examples in order:

### Phase 1: Presentation Content (Critical)

1. **main-presentation-add-video-zone**
   ```bash
   ./bin/main-presentation-add-video-zone --id <pres-id> --content-id <video-id> \
     --x 0 --y 0 --width 1920 --height 1080
   ```

2. **main-presentation-add-image-zone**
   ```bash
   ./bin/main-presentation-add-image-zone --id <pres-id> --content-id <image-id> \
     --x 0 --y 0 --width 1920 --height 1080 --duration 10
   ```

### Phase 2: Distribution (Critical)

3. **main-presentation-publish**
   ```bash
   ./bin/main-presentation-publish --id <pres-id>
   ```

4. **main-group-assign-presentation**
   ```bash
   ./bin/main-group-assign-presentation --group-id <group-id> \
     --presentation-id <pres-id>
   ```

### Phase 3: Verification (Nice to have)

5. **main-device-current-presentation**
   ```bash
   ./bin/main-device-current-presentation --serial <serial>
   ```

## Workaround: Manual Steps

Until the missing examples are implemented, complete these steps manually:

### Using BrightAuthor:connected

1. Open BrightAuthor:connected
2. Find the presentation created by the script
3. Add video zone with uploaded video
4. Add image zone with uploaded image, set duration
5. Publish the presentation
6. Assign to device/group
7. Verify playback

### Using REST API directly

```bash
# Get the IDs from script output
VIDEO_ID=12345
IMAGE_ID=67890
PRES_ID=99999

# Add zones (example - actual structure varies)
curl -X PUT "https://api.bsn.cloud/2024-10/REST/presentation/$PRES_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "zones": [
      {
        "type": "VideoOrImages",
        "content": [{"id": '$VIDEO_ID'}],
        "geometry": {"x": 0, "y": 0, "width": 1920, "height": 1080}
      },
      {
        "type": "Image",
        "content": [{"id": '$IMAGE_ID'}],
        "duration": 10,
        "geometry": {"x": 0, "y": 0, "width": 1920, "height": 1080}
      }
    ]
  }'

# Publish
curl -X POST "https://api.bsn.cloud/2024-10/REST/presentation/$PRES_ID/publish" \
  -H "Authorization: Bearer $TOKEN"
```

## Future Enhancements

Once missing examples are implemented, the script can be extended to:

- Complete end-to-end deployment
- Wait for distribution to complete
- Verify playback on device
- Handle multiple devices
- Support playlist creation
- Schedule presentation
- Monitor distribution status
