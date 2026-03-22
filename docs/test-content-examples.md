# Content and Presentation API Test Programs

This document explains how to use the test programs for working with the BSN.cloud Content and Presentation APIs.

## Overview

### Content Management Programs

There are **four** content management test programs in the `examples/` directory:

1. **main-content-upload** - Upload content files to BSN.cloud
2. **main-content-list** - List content files on BSN.cloud
3. **main-content-download** - Download content files from BSN.cloud
4. **main-content-delete** - Delete content files from BSN.cloud

### Presentation Management Programs

There are **eight** presentation management test programs in the `examples/` directory:

1. **main-presentation-create** - Create presentations on BSN.cloud
2. **main-presentation-list** - List presentations on BSN.cloud
3. **main-presentation-info** - Get presentation details by ID
4. **main-presentation-info-by-name** - Get presentation details by name
5. **main-presentation-update** - Update presentation properties
6. **main-presentation-delete** - Delete presentations by ID
7. **main-presentation-delete-by-filter** - Delete presentations using filters
8. **main-presentation-count** - Get count of presentations

## Common Setup

### Required Environment Variables

All examples require authentication credentials:
```bash
export BS_CLIENT_ID=your_client_id
export BS_SECRET=your_client_secret
```

### Optional Environment Variables

These environment variables provide default values for command-line flags:
```bash
export BS_NETWORK=your_network_name          # Default network name
export BS_SERIAL=device_serial                # Default device serial number
export BS_DEVICE_ID=device_id                 # Default device ID
export BS_PRESENTATION_ID=presentation_id     # Default presentation ID
export BS_CONTENT_ID=content_id               # Default content ID
export BS_GROUP_ID=group_id                   # Default group ID
```

**Note:** Command-line flags always take precedence over environment variables.

### Building Examples

```bash
# Build all examples
make build-examples

# Or build individually - Content programs
go build -o bin/main-content-upload ./examples/main-content-upload
go build -o bin/main-content-list ./examples/main-content-list
go build -o bin/main-content-download ./examples/main-content-download
go build -o bin/main-content-delete ./examples/main-content-delete

# Presentation programs
go build -o bin/main-presentation-create ./examples/main-presentation-create
go build -o bin/main-presentation-list ./examples/main-presentation-list
go build -o bin/main-presentation-info ./examples/main-presentation-info
go build -o bin/main-presentation-info-by-name ./examples/main-presentation-info-by-name
go build -o bin/main-presentation-update ./examples/main-presentation-update
go build -o bin/main-presentation-delete ./examples/main-presentation-delete
go build -o bin/main-presentation-delete-by-filter ./examples/main-presentation-delete-by-filter
go build -o bin/main-presentation-count ./examples/main-presentation-count
```

### Network Selection

All programs support flexible network selection:
1. Use `--network` / `-n` flag
2. Falls back to `BS_NETWORK` environment variable
3. If multiple networks available, prompts interactively
4. If only one network, uses it automatically

### Environment Variable Fallbacks

Examples support environment variable fallbacks for common parameters, making it easier to script repeated operations on the same resources:

**Device operations:**
```bash
# Set default device
export BS_SERIAL=BS123456789

# Now you can omit --serial flag
./bin/main-device-info                    # Uses BS_SERIAL
./bin/rdws-reboot --type normal           # Uses BS_SERIAL
./bin/main-device-info --serial BS999999  # Flag overrides BS_SERIAL
```

**Presentation operations:**
```bash
# Set default presentation
export BS_PRESENTATION_ID=12345

# Omit --id flag
./bin/main-presentation-info              # Uses BS_PRESENTATION_ID
./bin/main-presentation-update --name "New Name"
```

**Content operations:**
```bash
# Set default content
export BS_CONTENT_ID=67890

./bin/main-content-download               # Uses BS_CONTENT_ID
./bin/main-content-download --output /tmp/file.mp4
```

**Priority:** Command-line flags always take precedence over environment variables.

## Test Programs

### 1. main-content-upload

Uploads content files to BSN.cloud.

**Location:** `examples/main-content-upload/main.go`

**Usage:**
```bash
# Upload a file
./bin/main-content-upload --file video.mp4

# Upload with virtual path
./bin/main-content-upload --file image.jpg --virtual-path /videos/

# Upload and output JSON
./bin/main-content-upload --file document.pdf --json

# Upload with debug logging
./bin/main-content-upload --file test.mp4 --debug --verbose
```

**Key Flags:**
- `--file <path>`: Path to file to upload (required)
- `--virtual-path <path>`: Virtual path where file should be placed (e.g., /videos/)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--debug`: Enable debug logging (shows HTTP requests/responses)
- `--timeout 30`: Request timeout in seconds

**Output:**
- Content ID
- File name
- File size
- Virtual path
- Media type
- Upload date
- Upload completion status

### 2. main-content-list

Lists content files on BSN.cloud.

**Location:** `examples/main-content-list/main.go`

**Usage:**
```bash
# List content files (first page)
./bin/main-content-list

# List all content files (paginate through all)
./bin/main-content-list --all

# Filter by media type
./bin/main-content-list --filter "mediaType eq 'Video'"

# Sort by name
./bin/main-content-list --sort "name asc" --all

# Output as JSON
./bin/main-content-list --json
```

**Key Flags:**
- `--page-size 100`: Number of items per page
- `--filter <expr>`: Filter expression (e.g., "mediaType eq 'Video'")
- `--sort <expr>`: Sort expression (e.g., "name asc")
- `--all`: Retrieve all content files (auto-paginate)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information

**Filter Expression Examples:**
```bash
# Exact match
--filter "name eq 'video.mp4'"

# Contains
--filter "name contains 'retail'"

# Media type
--filter "mediaType eq 'Video'"

# File size
--filter "fileSize lt 1000000"

# Multiple conditions (AND)
--filter "mediaType eq 'Video' and fileSize lt 1000000"

# Starts with
--filter "name startsWith 'temp_'"
```

### 3. main-content-download

Downloads content files from BSN.cloud.

**Location:** `examples/main-content-download/main.go`

**Usage:**
```bash
# Download by content ID
./bin/main-content-download --id 12345

# Download to specific file
./bin/main-content-download --id 12345 --output /tmp/video.mp4

# Show file info only (no download)
./bin/main-content-download --id 12345 --info-only
```

**Key Flags:**
- `--id <id>`: Content file ID to download (required)
- `--output <path>`: Output file path (default: use original filename)
- `--info-only`: Show file info without downloading
- `--network <name>` / `-n`: Network name
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

### 4. main-content-delete

Deletes content files from BSN.cloud.

**Location:** `examples/main-content-delete/main.go`

**Usage:**
```bash
# Preview files to delete (dry-run)
./bin/main-content-delete --filter "name contains 'old'" --dry-run

# Delete files matching filter
./bin/main-content-delete --filter "name startsWith 'temp_'" --yes

# Delete with confirmation prompt
./bin/main-content-delete --filter "name contains 'test'"
```

**Key Flags:**
- `--filter <expr>`: Filter expression (required)
- `--network <name>` / `-n`: Network name
- `--dry-run`: Preview files without deleting
- `--yes`: Skip confirmation prompt
- `--json`: Output as JSON
- `--timeout 30`: Request timeout in seconds

**Filter Examples:**
- `"name contains 'old'"` - Files with "old" in name
- `"mediaType eq 'Video' and fileSize lt 1000000"` - Videos smaller than 1MB
- `"name startsWith 'temp_'"` - Files starting with "temp_"

## Example Workflows

### Complete Content Management Workflow

```bash
# 1. Upload content
./bin/main-content-upload --file video.mp4 --virtual-path /videos/
# Output: Content ID (e.g., 12345)

# 2. List all content to verify
./bin/main-content-list --filter "name contains 'video'"

# 3. Download content to verify
./bin/main-content-download --id 12345 --output /tmp/verify.mp4

# 4. Delete old content
./bin/main-content-delete --filter "name contains 'old'" --dry-run
./bin/main-content-delete --filter "name contains 'old'" --yes
```

### Bulk Upload Workflow

```bash
# Upload multiple files with virtual paths
for file in videos/*.mp4; do
  ./bin/main-content-upload --file "$file" --virtual-path /videos/ --json
done
```

### Content Cleanup Workflow

```bash
# List old content
./bin/main-content-list --filter "name contains 'test'" --all

# Preview deletion (dry-run)
./bin/main-content-delete --filter "name contains 'test'" --dry-run

# Delete after review
./bin/main-content-delete --filter "name contains 'test'" --yes
```

### Content Inventory Workflow

```bash
# Get all videos
./bin/main-content-list --filter "mediaType eq 'Video'" --all --json > videos.json

# Get all images
./bin/main-content-list --filter "mediaType eq 'Image'" --all --json > images.json

# Get large files (over 10MB)
./bin/main-content-list --filter "fileSize gt 10485760" --all
```

---

# Presentation API Test Programs

## Presentation Test Programs

### 1. main-presentation-create

Creates presentations on BSN.cloud.

**Location:** `examples/main-presentation-create/main.go`

**Usage:**
```bash
# Create a basic presentation
./bin/main-presentation-create --name "My Presentation"

# Create with description and tags
./bin/main-presentation-create --name "Retail Display" --description "Store display" --tags "retail,store"

# Create for specific player model
./bin/main-presentation-create --name "4K Content" --model XT1144

# Create and output JSON
./bin/main-presentation-create --name "Test" --json
```

**Key Flags:**
- `--name <name>`: Presentation name (required)
- `--description <desc>`: Presentation description
- `--tags <tags>`: Comma-separated list of tags
- `--model <model>`: Player model (e.g., HD1024, XT1144, XD1034)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Output:**
- Presentation ID
- Name
- Type
- Publish state
- Creation date
- Tags (if specified)

### 2. main-presentation-list

Lists presentations on BSN.cloud.

**Location:** `examples/main-presentation-list/main.go`

**Usage:**
```bash
# List presentations (first page)
./bin/main-presentation-list

# List all presentations (paginate through all)
./bin/main-presentation-list --all

# Filter by name
./bin/main-presentation-list --filter "name contains 'Retail'"

# Sort by name
./bin/main-presentation-list --sort "name asc" --all

# Output as JSON
./bin/main-presentation-list --json
```

**Key Flags:**
- `--page-size 100`: Number of items per page
- `--filter <expr>`: Filter expression
- `--sort <expr>`: Sort expression (e.g., "name asc")
- `--all`: Retrieve all presentations (auto-paginate)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information (zones, content items, tags)

**Filter Expression Examples:**
```bash
# Name contains
--filter "name contains 'Retail'"

# Exact match
--filter "name eq 'My Presentation'"

# Publish state
--filter "publishState eq 'Published'"

# Multiple conditions
--filter "name contains 'Store' and publishState eq 'Published'"
```

### 3. main-presentation-info

Retrieves presentation details by ID.

**Location:** `examples/main-presentation-info/main.go`

**Usage:**
```bash
# Get presentation by ID
./bin/main-presentation-info --id 12345

# Get presentation and output JSON
./bin/main-presentation-info --id 12345 --json

# Get with verbose output
./bin/main-presentation-info --id 12345 --verbose
```

**Key Flags:**
- `--id <id>`: Presentation ID (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

### 4. main-presentation-info-by-name

Retrieves presentation details by name.

**Location:** `examples/main-presentation-info-by-name/main.go`

**Usage:**
```bash
# Get presentation by name
./bin/main-presentation-info-by-name --name "Retail Display"

# Get and output JSON
./bin/main-presentation-info-by-name --name "My Presentation" --json
```

**Key Flags:**
- `--name <name>`: Presentation name (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Note:** If multiple presentations have the same name, this returns the first match.

### 5. main-presentation-update

Updates presentation properties.

**Location:** `examples/main-presentation-update/main.go`

**Usage:**
```bash
# Update presentation name
./bin/main-presentation-update --id 12345 --name "Updated Display"

# Update multiple properties
./bin/main-presentation-update --id 12345 --name "New Name" --description "New description"

# Update and output JSON
./bin/main-presentation-update --id 12345 --name "Updated" --json
```

**Key Flags:**
- `--id <id>`: Presentation ID (required)
- `--name <name>`: New presentation name
- `--description <desc>`: New presentation description
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

### 6. main-presentation-delete

Deletes presentations from BSN.cloud.

**Location:** `examples/main-presentation-delete/main.go`

**Usage:**
```bash
# Delete presentation (with confirmation)
./bin/main-presentation-delete --id 12345

# Delete without confirmation
./bin/main-presentation-delete --id 12345 --force

# Delete and output JSON
./bin/main-presentation-delete --id 12345 --force --json
```

**Key Flags:**
- `--id <id>`: Presentation ID to delete (required)
- `--network <name>` / `-n`: Network name
- `--force`: Skip confirmation prompt
- `--json`: Output as JSON
- `--timeout 30`: Request timeout in seconds

### 7. main-presentation-delete-by-filter

Deletes presentations using filter expressions.

**Location:** `examples/main-presentation-delete-by-filter/main.go`

**Usage:**
```bash
# Preview presentations to delete (dry-run)
./bin/main-presentation-delete-by-filter --filter "name contains 'test'" --dry-run

# Delete presentations matching filter
./bin/main-presentation-delete-by-filter --filter "name startsWith 'old_'" --yes

# Delete with confirmation prompt
./bin/main-presentation-delete-by-filter --filter "name contains 'temp'"
```

**Key Flags:**
- `--filter <expr>`: Filter expression (required)
- `--network <name>` / `-n`: Network name
- `--dry-run`: Preview presentations without deleting
- `--yes`: Skip confirmation prompt
- `--json`: Output as JSON
- `--timeout 30`: Request timeout in seconds

**Filter Examples:**
- `"name contains 'test'"` - Presentations with "test" in name
- `"name startsWith 'old_'"` - Presentations starting with "old_"
- `"publishState eq 'Draft'"` - All draft presentations

### 8. main-presentation-count

Gets count of presentations on network.

**Location:** `examples/main-presentation-count/main.go`

**Usage:**
```bash
# Count all presentations
./bin/main-presentation-count

# Count with filter
./bin/main-presentation-count --filter "name contains 'Retail'"

# Count published presentations
./bin/main-presentation-count --filter "publishState eq 'Published'"
```

**Key Flags:**
- `--network <name>` / `-n`: Network name
- `--filter <expr>`: Filter expression
- `--timeout 30`: Request timeout in seconds

## Presentation Workflows

### Complete Presentation Lifecycle

```bash
# 1. Create presentation
./bin/main-presentation-create --name "Store Display" --model HD1024
# Output: Presentation ID (e.g., 12345)

# 2. Verify creation
./bin/main-presentation-info --id 12345

# 3. Update presentation
./bin/main-presentation-update --id 12345 --description "Updated display content"

# 4. List all presentations to verify
./bin/main-presentation-list --filter "name contains 'Store'"

# 5. Delete when no longer needed
./bin/main-presentation-delete --id 12345
```

### Bulk Presentation Management

```bash
# Create multiple presentations
for name in "Store-A" "Store-B" "Store-C"; do
  ./bin/main-presentation-create --name "$name" --tags "retail,store" --json
done

# List all store presentations
./bin/main-presentation-list --filter "name contains 'Store'" --all

# Count presentations by tag
./bin/main-presentation-count --filter "tags contains 'retail'"
```

### Presentation Cleanup

```bash
# List test presentations
./bin/main-presentation-list --filter "name contains 'test'" --all

# Preview deletion (dry-run)
./bin/main-presentation-delete-by-filter --filter "name contains 'test'" --dry-run

# Delete after review
./bin/main-presentation-delete-by-filter --filter "name contains 'test'" --yes
```

### Presentation Inventory

```bash
# Export all presentations to JSON
./bin/main-presentation-list --all --json > presentations.json

# Count by publish state
./bin/main-presentation-count --filter "publishState eq 'Published'"
./bin/main-presentation-count --filter "publishState eq 'Draft'"

# Find presentations by name pattern
./bin/main-presentation-list --filter "name startsWith 'Retail'" --all
```

### Content and Presentation Integration

```bash
# 1. Upload content files
./bin/main-content-upload --file video.mp4 --virtual-path /videos/
./bin/main-content-upload --file image.jpg --virtual-path /images/

# 2. Create presentation
./bin/main-presentation-create --name "Digital Signage" --model XT1144

# 3. Get presentation ID
PRES_ID=$(./bin/main-presentation-info-by-name --name "Digital Signage" --json | jq -r '.id')

# 4. Verify presentation
./bin/main-presentation-info --id $PRES_ID --verbose

# 5. List content for verification
./bin/main-content-list --filter "virtualPath contains 'videos'"
```

## Common Patterns

### Output Formats

- **Default**: Human-readable formatted output
- `--json`: Structured JSON for scripting and automation
- `--verbose`: Additional details (timestamps, IDs, debug info)
- `--debug`: HTTP request/response logging (upload only)

### Confirmation Prompts

Destructive operations (delete) require confirmation:
- Interactive prompt asks for verification
- Use `--yes` flag to skip confirmation
- Use `--dry-run` to preview without executing

### Authentication

All examples use OAuth2 Client Credentials flow:
- Credentials from environment: `BS_CLIENT_ID`, `BS_SECRET`
- Token obtained automatically by SDK
- Token cached and reused until expiration
- No manual token management required

### Pagination

List operations support flexible querying:
- `--page-size N`: Control results per page
- `--all`: Auto-paginate through all results
- `--filter <expr>`: OData-like filter expressions
- `--sort <expr>`: Sort by field(s)

## Exit Codes

All examples use standard exit codes:
- **0**: Success
- **1**: Error (invalid arguments, authentication failed, operation failed)

Use in scripts:
```bash
if ./bin/main-content-upload --file video.mp4; then
  echo "Upload successful"
else
  echo "Upload failed"
  exit 1
fi
```

## Troubleshooting

### Authentication Issues

```bash
# Verify credentials are set
echo $BS_CLIENT_ID
echo $BS_SECRET

# Test authentication
./bin/main-content-list --verbose
```

### Network Selection Issues

```bash
# List available networks
./bin/main-content-list --verbose

# Specify network explicitly
./bin/main-content-list --network "Production"

# Set default network
export BS_NETWORK="Production"
./bin/main-content-list
```

### Upload Issues

```bash
# Enable debug logging
./bin/main-content-upload --file video.mp4 --debug --verbose

# Check file exists and is readable
ls -lh video.mp4

# Verify network connectivity
./bin/main-content-list
```

### Download Issues

```bash
# Verify content ID exists
./bin/main-content-list --all | grep 12345

# Check file info without downloading
./bin/main-content-download --id 12345 --info-only

# Ensure output directory exists
mkdir -p /tmp/downloads
./bin/main-content-download --id 12345 --output /tmp/downloads/video.mp4
```

### Presentation Issues

```bash
# Verify presentation exists
./bin/main-presentation-list --filter "name eq 'My Presentation'"

# Get presentation by name if ID unknown
./bin/main-presentation-info-by-name --name "My Presentation" --verbose

# Check presentation count
./bin/main-presentation-count

# List all presentations to debug
./bin/main-presentation-list --all --verbose
```

## Help Information

All examples include `--help` for detailed usage information:

### Content Programs
```bash
./bin/main-content-upload --help
./bin/main-content-list --help
./bin/main-content-download --help
./bin/main-content-delete --help
```

### Presentation Programs
```bash
./bin/main-presentation-create --help
./bin/main-presentation-list --help
./bin/main-presentation-info --help
./bin/main-presentation-info-by-name --help
./bin/main-presentation-update --help
./bin/main-presentation-delete --help
./bin/main-presentation-delete-by-filter --help
./bin/main-presentation-count --help
```

## Related Documentation

- `examples/README.md` - Complete examples documentation covering all SDK features
- `docs/all-apis.md` - Complete API reference documentation
