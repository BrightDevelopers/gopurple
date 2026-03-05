#!/bin/bash
set -e  # Exit on error

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory and bin directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="${SCRIPT_DIR}/../../bin"
CONFIG_FILE="${1:-presentation-config.json}"
SKIP_CLEANUP="${SKIP_CLEANUP:-false}"

echo -e "${BLUE}=== BSN.cloud Presentation Deployment Script ===${NC}\n"

# Check if config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo -e "${RED}Error: Config file '$CONFIG_FILE' not found${NC}"
    echo "Usage: $0 [config-file]"
    echo ""
    echo "Example config file format (presentation-config.json):"
    cat << 'EOF'
{
  "presentation_name": "Store Display",
  "video_file": "videos/promo.mp4",
  "image_file": "images/logo.png",
  "image_duration": 10,
  "network": "Production",
  "device_serial": "BS123456789"
}
EOF
    exit 1
fi

# Parse config file using jq
if ! command -v jq &> /dev/null; then
    echo -e "${RED}Error: jq is required but not installed${NC}"
    echo "Install with: sudo dnf install jq  (or apt-get, brew, etc.)"
    exit 1
fi

echo -e "${GREEN}Reading configuration from $CONFIG_FILE${NC}"
PRESENTATION_NAME=$(jq -r '.presentation_name' "$CONFIG_FILE")
VIDEO_FILE=$(jq -r '.video_file' "$CONFIG_FILE")
IMAGE_FILE=$(jq -r '.image_file' "$CONFIG_FILE")
IMAGE_DURATION=$(jq -r '.image_duration' "$CONFIG_FILE")
NETWORK=$(jq -r '.network // empty' "$CONFIG_FILE")
DEVICE_SERIAL=$(jq -r '.device_serial // empty' "$CONFIG_FILE")

echo "  Presentation: $PRESENTATION_NAME"
echo "  Video: $VIDEO_FILE"
echo "  Image: $IMAGE_FILE"
echo "  Image Duration: ${IMAGE_DURATION}s"
[ -n "$NETWORK" ] && echo "  Network: $NETWORK"
[ -n "$DEVICE_SERIAL" ] && echo "  Device: $DEVICE_SERIAL"
echo ""

# Set network environment variable if specified
[ -n "$NETWORK" ] && export BS_NETWORK="$NETWORK"
[ -n "$DEVICE_SERIAL" ] && export BS_SERIAL="$DEVICE_SERIAL"

# Verify files exist
if [ ! -f "$VIDEO_FILE" ]; then
    echo -e "${RED}Error: Video file '$VIDEO_FILE' not found${NC}"
    exit 1
fi

if [ ! -f "$IMAGE_FILE" ]; then
    echo -e "${RED}Error: Image file '$IMAGE_FILE' not found${NC}"
    exit 1
fi

# Verify example programs exist
REQUIRED_PROGRAMS=(
    "main-content-list"
    "main-content-delete"
    "main-content-upload"
    "main-presentation-list"
    "main-presentation-delete-by-filter"
    "main-presentation-create"
    "main-presentation-info-by-name"
)

for prog in "${REQUIRED_PROGRAMS[@]}"; do
    if [ ! -f "$BIN_DIR/$prog" ]; then
        echo -e "${RED}Error: Required program '$prog' not found in $BIN_DIR${NC}"
        echo "Run: make build-examples"
        exit 1
    fi
done

echo -e "${BLUE}=== Step 1: Clean up existing content ===${NC}"

# Get base filenames for filtering
VIDEO_BASENAME=$(basename "$VIDEO_FILE")
IMAGE_BASENAME=$(basename "$IMAGE_FILE")

echo "Checking for existing content files..."

# Check if video exists (using 'contains' since 'eq' is not supported)
echo "  Checking for video: $VIDEO_BASENAME"
VIDEO_LIST=$(timeout 10 "$BIN_DIR/main-content-list" --filter "name contains '$VIDEO_BASENAME'" --json 2>&1 || echo '{"error":"timeout or failed"}')
if echo "$VIDEO_LIST" | jq -e . >/dev/null 2>&1; then
    VIDEO_EXISTS=$(echo "$VIDEO_LIST" | jq -r '.items | length' 2>/dev/null || echo "0")
    if [ "$VIDEO_EXISTS" -gt 0 ] 2>/dev/null; then
        echo -e "${YELLOW}  Deleting existing video: $VIDEO_BASENAME${NC}"
        "$BIN_DIR/main-content-delete" --filter "name contains '$VIDEO_BASENAME'" --yes 2>&1 | grep -v "^2026/" || true
    else
        echo "  Video not found in cloud (will upload fresh)"
    fi
else
    echo -e "${YELLOW}  Warning: Could not check for existing video (may not have list permissions)${NC}"
    echo "  Proceeding with upload..."
fi

# Check if image exists (using 'contains' since 'eq' is not supported)
echo "  Checking for image: $IMAGE_BASENAME"
IMAGE_LIST=$(timeout 10 "$BIN_DIR/main-content-list" --filter "name contains '$IMAGE_BASENAME'" --json 2>&1 || echo '{"error":"timeout or failed"}')
if echo "$IMAGE_LIST" | jq -e . >/dev/null 2>&1; then
    IMAGE_EXISTS=$(echo "$IMAGE_LIST" | jq -r '.items | length' 2>/dev/null || echo "0")
    if [ "$IMAGE_EXISTS" -gt 0 ] 2>/dev/null; then
        echo -e "${YELLOW}  Deleting existing image: $IMAGE_BASENAME${NC}"
        "$BIN_DIR/main-content-delete" --filter "name contains '$IMAGE_BASENAME'" --yes 2>&1 | grep -v "^2026/" || true
    else
        echo "  Image not found in cloud (will upload fresh)"
    fi
else
    echo -e "${YELLOW}  Warning: Could not check for existing image (may not have list permissions)${NC}"
    echo "  Proceeding with upload..."
fi

echo ""
echo -e "${BLUE}=== Step 2: Upload content files ===${NC}"

echo "Uploading video: $VIDEO_FILE"
VIDEO_RESULT=$("$BIN_DIR/main-content-upload" --file "$VIDEO_FILE" --json 2>&1)
if echo "$VIDEO_RESULT" | jq -e . >/dev/null 2>&1; then
    VIDEO_ID=$(echo "$VIDEO_RESULT" | jq -r '.contentID')
    if [ "$VIDEO_ID" != "null" ] && [ -n "$VIDEO_ID" ]; then
        echo -e "${GREEN}  ✓ Video uploaded (ID: $VIDEO_ID)${NC}"
    else
        echo -e "${RED}  ✗ Video upload failed (invalid response)${NC}"
        echo "$VIDEO_RESULT"
        exit 1
    fi
else
    echo -e "${RED}  ✗ Video upload failed${NC}"
    echo "$VIDEO_RESULT" | grep -v "^Creating\|^Authenticating\|^Authentication\|^Using\|^Getting\|^Uploading\|^File size\|^Upload process\|^Reading\|^SHA1\|^Initiating" || echo "$VIDEO_RESULT"
    exit 1
fi

echo "Uploading image: $IMAGE_FILE"
IMAGE_RESULT=$("$BIN_DIR/main-content-upload" --file "$IMAGE_FILE" --json 2>&1)
if echo "$IMAGE_RESULT" | jq -e . >/dev/null 2>&1; then
    IMAGE_ID=$(echo "$IMAGE_RESULT" | jq -r '.contentID')
    if [ "$IMAGE_ID" != "null" ] && [ -n "$IMAGE_ID" ]; then
        echo -e "${GREEN}  ✓ Image uploaded (ID: $IMAGE_ID)${NC}"
    else
        echo -e "${RED}  ✗ Image upload failed (invalid response)${NC}"
        echo "$IMAGE_RESULT"
        exit 1
    fi
else
    echo -e "${RED}  ✗ Image upload failed${NC}"
    echo "$IMAGE_RESULT" | grep -v "^Creating\|^Authenticating\|^Authentication\|^Using\|^Getting\|^Uploading\|^File size\|^Upload process\|^Reading\|^SHA1\|^Initiating" || echo "$IMAGE_RESULT"
    exit 1
fi

echo ""
echo -e "${BLUE}=== Step 3: Clean up existing presentation ===${NC}"

# Check if presentation exists (using 'contains' since 'eq' is not supported)
echo "  Checking for presentation: $PRESENTATION_NAME"
PRES_LIST=$(timeout 10 "$BIN_DIR/main-presentation-list" --filter "name contains '$PRESENTATION_NAME'" --json 2>&1 || echo '{"error":"timeout or failed"}')
if echo "$PRES_LIST" | jq -e . >/dev/null 2>&1; then
    PRES_EXISTS=$(echo "$PRES_LIST" | jq -r '.items | length' 2>/dev/null || echo "0")
    if [ "$PRES_EXISTS" -gt 0 ] 2>/dev/null; then
        echo -e "${YELLOW}  Deleting existing presentation: $PRESENTATION_NAME${NC}"
        "$BIN_DIR/main-presentation-delete-by-filter" --filter "name contains '$PRESENTATION_NAME'" --yes 2>&1 | grep -v "^2026/" || true
    else
        echo "  Presentation not found (will create fresh)"
    fi
else
    echo -e "${YELLOW}  Warning: Could not check for existing presentation${NC}"
    echo "  Proceeding with creation..."
fi

echo ""
echo -e "${BLUE}=== Step 4: Create presentation ===${NC}"

echo "Creating presentation: $PRESENTATION_NAME"
PRES_RESULT=$("$BIN_DIR/main-presentation-create" --name "$PRESENTATION_NAME" --json)
PRES_ID=$(echo "$PRES_RESULT" | jq -r '.id')
echo -e "${GREEN}  ✓ Presentation created (ID: $PRES_ID)${NC}"

# Verify creation
"$BIN_DIR/main-presentation-info-by-name" --name "$PRESENTATION_NAME"

echo ""
echo -e "${RED}=== MISSING FUNCTIONALITY ===${NC}"
echo -e "${YELLOW}"
cat << 'EOF'
⚠️  The following steps CANNOT be completed with existing example programs:

MISSING: Add content to presentation
  - Need: Example program to add video zone to presentation
  - Need: Example program to add image zone with duration
  - API: PUT /2024-10/presentation/{id} with zones configuration

MISSING: Publish/distribute presentation
  - Need: Example program to publish presentation
  - Need: Example program to assign presentation to device/group
  - API: POST /2024-10/presentation/{id}/publish
  - API: PUT /2024-10/group/{id} with contentDistribution

MISSING: Verify playback
  - Need: Example program to check current presentation on device
  - API: GET /2024-10/device/{id}/status with presentation info

To complete this workflow, you need to implement these additional examples:
  1. examples/main-presentation-add-content
  2. examples/main-presentation-publish
  3. examples/main-group-assign-presentation (or device-assign-presentation)
  4. examples/main-device-current-presentation
EOF
echo -e "${NC}"

echo ""
echo -e "${BLUE}=== What WAS accomplished ===${NC}"
echo -e "${GREEN}✓ Cleaned up old content files${NC}"
echo -e "${GREEN}✓ Uploaded video (ID: $VIDEO_ID)${NC}"
echo -e "${GREEN}✓ Uploaded image (ID: $IMAGE_ID)${NC}"
echo -e "${GREEN}✓ Deleted old presentation${NC}"
echo -e "${GREEN}✓ Created new empty presentation (ID: $PRES_ID)${NC}"

echo ""
echo -e "${YELLOW}Next steps (manual):${NC}"
echo "  1. Use BrightAuthor:connected or API to add zones with content"
echo "  2. Publish the presentation"
echo "  3. Assign to device/group"
echo "  4. Verify playback on device"

echo ""
echo -e "${BLUE}Content IDs for reference:${NC}"
echo "  Video ID: $VIDEO_ID"
echo "  Image ID: $IMAGE_ID"
echo "  Presentation ID: $PRES_ID"
