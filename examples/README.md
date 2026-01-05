# Examples Documentation

This directory contains 75 example programs demonstrating all SDK features.

## Quick Start

```bash
# Set credentials
export BS_CLIENT_ID=your_client_id
export BS_SECRET=your_client_secret
export BS_NETWORK=your_network_name  # Optional

# Build all examples
make build-examples

# Run any example with --help
./bin/list-devices --help
```

## Common Environment Variables

All examples require:
- `BS_CLIENT_ID`: BSN.cloud API client ID (required)
- `BS_SECRET`: BSN.cloud API client secret (required)
- `BS_NETWORK`: Default network name (optional)

## Common Flags

Most examples support:
- `--help`: Display usage information
- `--network` / `-n`: Network name selection
- `--verbose`: Show detailed output
- `--json`: Output as JSON for scripting
- `--timeout 30`: Request timeout in seconds
- `--force` / `-y`: Skip confirmation prompts (destructive operations)

---

## Authentication Examples (2)

### main-auth-info
Authenticate with BSN.cloud and display curl command examples.

**Flags:**
- `--help`: Display usage information
- `--timeout 30`: Request timeout in seconds
- `--show`: Show current token from .token file without re-authenticating

**Output:** Creates `.token` file with access token, provides curl examples

**Usage:**
```bash
./bin/main-auth-info
./bin/main-auth-info --show
```

### main-token-test
Test authentication and token validity.

**Usage:**
```bash
./bin/main-token-test
```

---

## B-Deploy Setup Management (10)

### bdeploy-add-setup
Create a B-Deploy setup record using JSON configuration.

**Flags:**
- `--help`: Display usage information
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds (overrides config file)

**Required Input:** JSON configuration file (see [Configuration Files](#configuration-files))

**Output:** Setup ID for device association

**Usage:**
```bash
./bin/bdeploy-add-setup config.json
./bin/bdeploy-add-setup --verbose examples/bdeploy-add-setup/config.json
```

### bdeploy-associate
Associate or dissociate a device with a B-Deploy setup.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--setup-id <id>`: Setup ID to associate (required unless --dissociate or --setup-name)
- `--setup-name <name>`: Setup package name to associate (alternative to --setup-id)
- `--dissociate`: Remove setup association
- `--network <name>`: Network name
- `--username <user>`: BSN.cloud username
- `--description <desc>`: Device description
- `--name <name>`: Device name
- `--create`: Create device if it doesn't exist
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/bdeploy-associate --serial BS123456789 --setup-id setup-abc123
./bin/bdeploy-associate --serial BS123456789 --setup-name "my-retail-setup"
./bin/bdeploy-associate --serial BS123456789 --dissociate
```

### bdeploy-delete-setup
Delete B-Deploy setup records by ID or package name.

**Flags:**
- `--setup-id <id>`: Setup ID to delete (use this OR --setup-name)
- `--setup-name <name>`: Package name to delete (use this OR --setup-id)
- `--network <name>` / `-n`: Network name
- `--force`: Skip confirmation prompt
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
# Delete by setup ID
./bin/bdeploy-delete-setup --setup-id setup-abc123

# Delete by package name (searches for matching setup)
./bin/bdeploy-delete-setup --setup-name retail-display-v1

# Force delete without confirmation
./bin/bdeploy-delete-setup --setup-id setup-abc123 --force
```

### bdeploy-delete-device
Permanently remove a device from B-Deploy.

**Flags:**
- `--serial <serial>`: Device serial number to delete
- `--device-id <id>`: Device ID to delete (alternative)
- `--network <name>` / `-n`: Network name
- `--force`: Skip confirmation prompt
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/bdeploy-delete-device --serial BS123456789
./bin/bdeploy-delete-device --device-id 12345 --force
```

### bdeploy-get-device
Get B-Deploy device setup record by serial number.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--network <name>` / `-n`: Network name
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/bdeploy-get-device --serial BS123456789
```

### bdeploy-get-records
Get B-Deploy setup records from network.

**Flags:**
- `--network <name>` / `-n`: Network name
- `--username <user>`: Filter by username
- `--package <name>`: Filter by package name
- `--page-size 50`: Number of records per page
- `--page 1`: Page number to retrieve
- `--summary`: Show only summary count
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/bdeploy-get-records
./bin/bdeploy-get-records --package retail-display-v1
```

### bdeploy-get-setup
Retrieve a specific B-Deploy setup record.

**Flags:**
- `--setup-id <id>`: Setup ID to retrieve (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output raw JSON instead of formatted structure
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/bdeploy-get-setup --setup-id setup-abc123
./bin/bdeploy-get-setup --setup-id setup-abc123 --json
```

### bdeploy-list-devices
List all B-Deploy devices on a network, with optional filtering by setup.

**Flags:**
- `--network <name>` / `-n`: Network name
- `--setup-id <id>`: Filter devices by setup ID (optional)
- `--setup-name <name>`: Filter devices by setup package name (optional)
- `--summary`: Show only summary count
- `--debug`: Show raw API response
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/bdeploy-list-devices
./bin/bdeploy-list-devices --summary
./bin/bdeploy-list-devices --setup-id "658f1dbef1d46c829f60a14f"
./bin/bdeploy-list-devices --setup-name "my-setup"
```

### bdeploy-list-setups
List all B-Deploy setup records for a network.

**Flags:**
- `--network <name>`: Network name (required)
- `--username <user>`: Filter by username
- `--package <name>`: Filter by package name
- `--verbose`: Show detailed information

**Output:** Table format with Setup ID, Package Name, Setup Type, Group, Network

**Usage:**
```bash
./bin/bdeploy-list-setups --network Production
./bin/bdeploy-list-setups --network Production --package retail
```

### bdeploy-update-setup
Update an existing B-Deploy setup record.

**Flags:**
- `--setup-id <id>`: Setup ID to update (required)
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Required Input:** JSON configuration file (same structure as add-setup)

**Usage:**
```bash
./bin/bdeploy-update-setup --setup-id setup-abc123 updated-config.json
```

---

## Device Management (9)

### main-device-info
Get real-time BrightSign device information via rDWS.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output raw JSON response
- `--verbose`: Show detailed information
- `--debug`: Enable debug mode (show HTTP requests/responses)
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-device-info --serial BS123456789
./bin/main-device-info --serial BS123456789 --json
```

### main-device-status
Get device status information.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-device-status --serial BS123456789
```

### main-device-errors
Retrieve device error logs.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--limit 100`: Number of error entries to retrieve
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-device-errors --serial BS123456789
./bin/main-device-errors --serial BS123456789 --limit 50
```

### main-device-operations
View operations on devices.

**Flags:**
- `--serial <serial>`: Device serial number
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-device-operations --serial BS123456789
```

### main-device-downloads
View content downloads on devices.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-device-downloads --serial BS123456789
```

### main-device-delete
Delete devices from BSN.cloud.

**Flags:**
- `--serial <serial>`: Device serial number to delete
- `--id <id>`: Device ID to delete (alternative)
- `--network <name>` / `-n`: Network name
- `--force`: Skip confirmation
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-device-delete --serial BS123456789
./bin/main-device-delete --id 12345 --force
```

### main-device-change-group
Change a device's group assignment.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--group <id>`: Target group ID (required)
- `--network <name>` / `-n`: Network name
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-device-change-group --serial BS123456789 --group 42
```

### main-device-local-dws
Manage local Diagnostic Web Server (DWS) on devices.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--enable`: Enable local DWS
- `--disable`: Disable local DWS
- `--password <pass>`: Set DWS password
- `--network <name>` / `-n`: Network name
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-device-local-dws --serial BS123456789 --enable
./bin/main-device-local-dws --serial BS123456789 --password mypass
```

### main--devices-list
List and query devices with filtering and pagination.

**Flags:**
- `--network <name>` / `-n`: Network name
- `--page-size 20`: Number of devices per page
- `--filter <expr>`: Filter expression
- `--sort <expr>`: Sort expression
- `--serial <serial>`: Get specific device by serial
- `--id <id>`: Get specific device by ID
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--quiet`: Suppress non-essential output
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main--devices-list
./bin/main--devices-list --serial BS123456789
./bin/main--devices-list --filter 'model=XD1033' --sort 'registrationDate'
```

---

## Content Management (4)

### main-content-upload
Upload content files to BSN.cloud.

**Flags:**
- `--file <path>`: Path to file to upload (required)
- `--virtual-path <path>`: Virtual path where file should be placed (e.g., /videos/)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--debug`: Enable debug logging (shows HTTP requests/responses)
- `--timeout 30`: Request timeout in seconds

**Output:** Content ID, file name, size, upload date, media type

**Usage:**
```bash
./bin/main-content-upload --file video.mp4
./bin/main-content-upload --file video.mp4 --virtual-path /videos/
```

### main-content-download
Download content files from BSN.cloud.

**Flags:**
- `--id <id>`: Content file ID to download (required)
- `--output <path>`: Output file path (default: use original filename)
- `--info-only`: Show file info without downloading
- `--network <name>` / `-n`: Network name
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-content-download --id 12345
./bin/main-content-download --id 12345 --output /tmp/video.mp4
./bin/main-content-download --id 12345 --info-only
```

### main-content-list
List content files on BSN.cloud.

**Flags:**
- `--network <name>` / `-n`: Network name
- `--page-size 100`: Number of items per page
- `--filter <expr>`: Filter expression
- `--sort <expr>`: Sort expression
- `--all`: Retrieve all content files (paginate through all results)
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-content-list
./bin/main-content-list --filter "mediaType eq 'Video'"
./bin/main-content-list --all
```

### main-content-delete
Delete content files from BSN.cloud.

**Flags:**
- `--filter <expr>`: Filter expression (required)
- `--network <name>` / `-n`: Network name
- `--dry-run`: Preview files without deleting
- `--yes`: Skip confirmation prompt
- `--json`: Output as JSON
- `--timeout 30`: Request timeout in seconds

**Filter Examples:**
- `"name contains 'old'"`
- `"mediaType eq 'Video' and fileSize lt 1000000"`
- `"name startsWith 'temp_'"`

**Usage:**
```bash
./bin/main-content-delete --filter "name contains 'old'" --dry-run
./bin/main-content-delete --filter "name startsWith 'temp_'" --yes
```

---

## Presentation Management (7)

### main-presentation-create
Create presentations on BSN.cloud.

**Flags:**
- `--name <name>`: Presentation name (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-presentation-create --name "Retail Display"
```

### main-presentation-list
List presentations on BSN.cloud.

**Flags:**
- `--network <name>` / `-n`: Network name
- `--page-size 100`: Number of items per page
- `--filter <expr>`: Filter expression
- `--sort <expr>`: Sort expression
- `--all`: Retrieve all presentations
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-presentation-list
./bin/main-presentation-list --filter "name contains 'Retail'"
```

### main-presentation-info
Retrieve presentation details by ID.

**Flags:**
- `--id <id>`: Presentation ID (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-presentation-info --id 12345
```

### main-presentation-info-by-name
Retrieve presentation details by name.

**Flags:**
- `--name <name>`: Presentation name (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-presentation-info-by-name --name "Retail Display"
```

### main-presentation-update
Update presentations on BSN.cloud.

**Flags:**
- `--id <id>`: Presentation ID (required)
- `--name <name>`: New presentation name
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-presentation-update --id 12345 --name "Updated Display"
```

### main-presentation-delete
Delete presentations from BSN.cloud.

**Flags:**
- `--id <id>`: Presentation ID to delete (required)
- `--network <name>` / `-n`: Network name
- `--force`: Skip confirmation
- `--json`: Output as JSON
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-presentation-delete --id 12345
./bin/main-presentation-delete --id 12345 --force
```

### main-presentation-delete-by-filter
Delete presentations using filters.

**Flags:**
- `--filter <expr>`: Filter expression (required)
- `--network <name>` / `-n`: Network name
- `--dry-run`: Preview presentations without deleting
- `--yes`: Skip confirmation
- `--json`: Output as JSON
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-presentation-delete-by-filter --filter "name contains 'old'" --dry-run
./bin/main-presentation-delete-by-filter --filter "name startsWith 'test_'" --yes
```

### main-presentation-count
Get count of presentations on network.

**Flags:**
- `--network <name>` / `-n`: Network name
- `--filter <expr>`: Filter expression
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-presentation-count
./bin/main-presentation-count --filter "name contains 'Retail'"
```

---

## Group Management (3)

### main-group-info
Retrieve BrightSign group information.

**Flags:**
- `--id <id>`: Group ID to retrieve (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-group-info --id 42
```

### main-group-update
Update BrightSign group information.

**Flags:**
- `--id <id>`: Group ID to update (required)
- `--name <name>`: New group name
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-group-update --id 42 --name "Retail Stores"
```

### main-group-delete
Delete BrightSign groups.

**Flags:**
- `--id <id>`: Group ID to delete (required)
- `--network <name>` / `-n`: Network name
- `--force`: Skip confirmation
- `--json`: Output as JSON
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-group-delete --id 42
./bin/main-group-delete --id 42 --force
```

---

## Subscription Management (3)

### main-subscription-operations
List subscription operations and permissions.

**Flags:**
- `--network <name>` / `-n`: Network name
- `--device-id <id>`: Filter by device ID
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-subscription-operations
./bin/main-subscription-operations --device-id 12345
```

### main-subscription-count
Get count of device subscriptions on network.

**Flags:**
- `--network <name>` / `-n`: Network name
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-subscription-count
```

### main-subscriptions-list
List device subscriptions on network.

**Flags:**
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/main-subscriptions-list
```

---

## Remote DWS Operations (34)

Remote Diagnostic Web Server (rDWS) operations allow you to manage and troubleshoot BrightSign devices remotely.

### rdws-info
Retrieve player information via rDWS.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output raw JSON response
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-info --serial BS123456789
```

### rdws-reboot
Remotely reboot BrightSign devices.

**Flags:**
- `--serial <serial>`: Device serial number
- `--id <id>`: Device ID
- `--type <type>`: Reboot type (normal, crash, factoryreset, disableautorun)
- `-t`: Alias for --type
- `--network <name>` / `-n`: Network name
- `-y`: Skip confirmation prompt
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-reboot --serial BS123456789 --type normal
./bin/rdws-reboot --serial BS123456789 -t factoryreset -y
```

### rdws-reprovision
Re-provision BrightSign devices.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--network <name>` / `-n`: Network name
- `-y`: Skip confirmation
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-reprovision --serial BS123456789
./bin/rdws-reprovision --serial BS123456789 -y
```

### rdws-snapshot
Capture screenshots from BrightSign devices.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--output <path>`: Output file path (default: screenshot.png)
- `--network <name>` / `-n`: Network name
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-snapshot --serial BS123456789
./bin/rdws-snapshot --serial BS123456789 --output capture.png
```

### rdws-time
Manage player time via rDWS.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-time --serial BS123456789
```

### rdws-time-set
Set player time via rDWS.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--time <time>`: Time to set (RFC3339 format or Unix timestamp)
- `--network <name>` / `-n`: Network name
- `-y`: Skip confirmation
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-time-set --serial BS123456789 --time "2025-01-15T10:00:00Z"
./bin/rdws-time-set --serial BS123456789 --time 1736935200
```

### rdws-dws-password
Get and set DWS passwords on BrightSign devices.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--set <password>`: Set DWS password
- `--network <name>` / `-n`: Network name
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-dws-password --serial BS123456789
./bin/rdws-dws-password --serial BS123456789 --set newpassword
```

### rdws-health
Check player health status.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-health --serial BS123456789
```

### rdws-local-dws
Manage local DWS on players.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--enable`: Enable local DWS
- `--disable`: Disable local DWS
- `--password <pass>`: Set DWS password
- `--network <name>` / `-n`: Network name
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-local-dws --serial BS123456789 --enable
./bin/rdws-local-dws --serial BS123456789 --password mypass
```

### File Operations

#### rdws-files-list
List files and directories on player storage.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--path <path>`: Directory path on device
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-files-list --serial BS123456789
./bin/rdws-files-list --serial BS123456789 --path /storage/sd/
```

#### rdws-files-upload
Upload a file to player storage.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--file <path>`: Local file path to upload (required)
- `--device-path <path>`: Destination path on device
- `--network <name>` / `-n`: Network name
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-files-upload --serial BS123456789 --file config.xml
./bin/rdws-files-upload --serial BS123456789 --file config.xml --device-path /storage/sd/
```

#### rdws-files-delete
Delete a file from player storage.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--path <path>`: File path on device (required)
- `--network <name>` / `-n`: Network name
- `-y`: Skip confirmation
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-files-delete --serial BS123456789 --path /storage/sd/old.xml
./bin/rdws-files-delete --serial BS123456789 --path /storage/sd/old.xml -y
```

#### rdws-files-rename
Rename a file on player storage.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--path <path>`: File path on device (required)
- `--new-name <name>`: New file name (required)
- `--network <name>` / `-n`: Network name
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-files-rename --serial BS123456789 --path /storage/sd/old.xml --new-name new.xml
```

#### rdws-files-create-folder
Create folders on player storage.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--path <path>`: Folder path to create (required)
- `--network <name>` / `-n`: Network name
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-files-create-folder --serial BS123456789 --path /storage/sd/configs/
```

### Network Operations

#### rdws-network-config
Manage network configuration on players.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-network-config --serial BS123456789
```

#### rdws-ping
Ping hosts from players.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--host <host>`: Host to ping (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-ping --serial BS123456789 --host google.com
./bin/rdws-ping --serial BS123456789 --host 8.8.8.8
```

#### rdws-dns-lookup
Perform DNS lookups on players.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--host <host>`: Host to lookup (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-dns-lookup --serial BS123456789 --host brightsignnetwork.com
```

#### rdws-diagnostics
Run network diagnostics on players.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-diagnostics --serial BS123456789
```

#### rdws-traceroute
Perform traceroute from players.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--host <host>`: Target host (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-traceroute --serial BS123456789 --host 8.8.8.8
```

#### rdws-packet-capture
Packet capture on players.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--network <name>` / `-n`: Network name
- `--output <path>`: Output file path
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-packet-capture --serial BS123456789 --output capture.pcap
```

#### rdws-network-neighborhood
Discover network neighborhood from players.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-network-neighborhood --serial BS123456789
```

### System Operations

#### rdws-registry-get
Read BrightSign player registry values.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--key <key>`: Registry key to read (required)
- `--network <name>` / `-n`: Network name
- `--json`: Output as JSON
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-registry-get --serial BS123456789 --key "networking"
```

#### rdws-registry-set
Write BrightSign player registry values.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--key <key>`: Registry key to set (required)
- `--value <value>`: Registry value (required)
- `--type <type>`: Value type (string, int, binary, etc.)
- `--network <name>` / `-n`: Network name
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-registry-set --serial BS123456789 --key "mykey" --value "myvalue"
```

#### rdws-logs-get
Retrieve BrightSign player log files.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--output <path>`: Output directory for logs
- `--network <name>` / `-n`: Network name
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-logs-get --serial BS123456789 --output /tmp/logs/
```

#### rdws-crashdump-get
Retrieve BrightSign player crash dumps.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--output <path>`: Output directory for dumps
- `--network <name>` / `-n`: Network name
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-crashdump-get --serial BS123456789 --output /tmp/crashdumps/
```

#### rdws-firmware-download
Remotely download and apply firmware updates.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--url <url>`: Firmware download URL (required)
- `--network <name>` / `-n`: Network name
- `-y`: Skip confirmation
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-firmware-download --serial BS123456789 --url https://example.com/firmware.bsfw
```

#### rdws-custom-data
Send custom data to BrightSign players via UDP.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--data <data>`: Data to send (required)
- `--network <name>` / `-n`: Network name
- `--port 5000`: UDP port
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-custom-data --serial BS123456789 --data "command:refresh"
./bin/rdws-custom-data --serial BS123456789 --data "alert" --port 6000
```

#### rdws-ssh
Manage SSH access on players, including enable/disable and password management.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--get`: Get current SSH status
- `--enable`: Enable SSH
- `--disable`: Disable SSH
- `--password <password>`: Set SSH password (optional, if omitted existing password unchanged)
- `--port <port>`: SSH port (default: 22)
- `--network <name>` / `-n`: Network name (uses BS_NETWORK env var if not specified)
- `--json`: Output raw JSON response
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
# Get SSH status
./bin/rdws-ssh --serial BS123456789 --get

# Enable SSH with password
./bin/rdws-ssh --serial BS123456789 --enable --password mypassword

# Enable SSH on custom port with password
./bin/rdws-ssh --serial BS123456789 --enable --port 2222 --password secure123

# Enable SSH without changing password
./bin/rdws-ssh --serial BS123456789 --enable

# Change password on already-enabled SSH
./bin/rdws-ssh --serial BS123456789 --enable --password newpassword

# Disable SSH
./bin/rdws-ssh --serial BS123456789 --disable
```

#### rdws-telnet
Manage telnet access on players.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--enable`: Enable telnet
- `--disable`: Disable telnet
- `--network <name>` / `-n`: Network name
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-telnet --serial BS123456789 --enable
./bin/rdws-telnet --serial BS123456789 --disable
```

#### rdws-reformat-storage
Remotely reformat storage devices on players.

**Flags:**
- `--serial <serial>`: Device serial number (required)
- `--device <device>`: Storage device to reformat (sd, ssd, usb)
- `--network <name>` / `-n`: Network name (uses BS_NETWORK env var if not specified)
- `-y`: Skip confirmation
- `--verbose`: Show detailed information
- `--timeout 30`: Request timeout in seconds

**Usage:**
```bash
./bin/rdws-reformat-storage --serial BS123456789 --device sd
export BS_NETWORK="Production"
./bin/rdws-reformat-storage --serial BS123456789 --device ssd -y
```

---

## Configuration Files

### B-Deploy Setup Configuration (JSON)

Used by `bdeploy-add-setup` and `bdeploy-update-setup`:

```json
{
  "networkName": "Production Network",
  "username": "admin@example.com",
  "packageName": "retail-display-v1",
  "setupType": "standalone",
  "timeZone": "America/New_York",
  "bsnGroupName": "Default",
  "network": {
    "timeServers": ["http://time.brightsignnetwork.com"],
    "interfaces": [
      {
        "id": "wired_eth0",
        "name": "eth0",
        "type": "Ethernet",
        "proto": "DHCPv4",
        "contentDownloadEnabled": true,
        "healthReportingEnabled": true
      }
    ]
  },
  "timeout": 30
}
```

**Required Fields:**
- `networkName`: BSN.cloud network name (must match exactly)
- `username`: Your BSN.cloud username/email
- `packageName`: Unique identifier for this setup configuration
- `setupType`: Setup type: `standalone` or `bsn`
- `network.interfaces`: At least one network interface configuration

**Optional Fields:**
- `timeZone`: Player timezone (default: `America/New_York`)
- `bsnGroupName`: BSN group name (default: `Default`)
- `timeout`: HTTP request timeout in seconds (default: `30`)

**Example configurations** can be found in:
- `examples/bdeploy-add-setup/config.json` - Basic Ethernet setup
- `examples/bdeploy-add-setup/config-annotated.json` - Detailed with comments
- `examples/bdeploy-add-setup/config-wifi-example.json` - WiFi configuration

---

## Common Patterns

### Network Selection
Most examples support flexible network selection:
- Use `--network` / `-n` flag to specify network name
- Falls back to `BS_NETWORK` environment variable
- If multiple networks available, prompts for interactive selection
- If only one network, uses it automatically

### Output Formats
- Default: Human-readable formatted output
- `--json`: Structured JSON for scripting and automation
- `--verbose`: Additional details (timestamps, IDs, debug info)
- `--debug`: HTTP request/response logging (some examples)

### Confirmation Prompts
Destructive operations (delete, reformat, factory reset) require confirmation:
- Interactive prompt asks for verification
- Use `--force` or `-y` flag to skip confirmation
- Useful for automation and scripting

### Authentication
All examples use OAuth2 Client Credentials flow:
- Credentials from environment: `BS_CLIENT_ID`, `BS_SECRET`
- Token obtained automatically by SDK
- Token cached and reused until expiration
- No manual token management required

### Pagination & Filtering
List operations support flexible querying:
- `--page-size N`: Control results per page
- `--filter <expr>`: OData-like filter expressions
- `--sort <expr>`: Sort by field(s)
- `--all`: Auto-paginate through all results

**Filter Expression Examples:**
```bash
# Exact match
--filter "model=XD1033"

# Contains
--filter "name contains 'retail'"

# Multiple conditions (AND)
--filter "mediaType eq 'Video' and fileSize lt 1000000"

# Comparison operators
--filter "registrationDate gt '2025-01-01'"
```

---

## Typical Workflows

### 1. Device Provisioning
Complete workflow for adding new devices:

```bash
# Step 1: Create B-Deploy setup record
./bin/bdeploy-add-setup examples/bdeploy-add-setup/config.json
# Output: setup-id (e.g., setup-abc123)

# Step 2: Associate device with setup
./bin/bdeploy-associate --serial BS123456789 --setup-id setup-abc123

# Step 3: Device boots and receives setup automatically
# Monitor with:
./bin/main-device-status --serial BS123456789
```

### 2. Content Upload & Display
Upload content and create presentations:

```bash
# Step 1: Upload content
./bin/main-content-upload --file video.mp4 --virtual-path /videos/
# Output: Content ID

# Step 2: Create presentation
./bin/main-presentation-create --name "Retail Display"
# Output: Presentation ID

# Step 3: Assign to devices via group management
./bin/main-group-update --id 42 --name "Retail Stores"
./bin/main-device-change-group --serial BS123456789 --group 42
```

### 3. Remote Device Management
Monitor and manage devices remotely:

```bash
# Get device info
./bin/rdws-info --serial BS123456789

# Check health status
./bin/rdws-health --serial BS123456789

# Capture screenshot
./bin/rdws-snapshot --serial BS123456789 --output current.png

# Reboot if needed
./bin/rdws-reboot --serial BS123456789 --type normal

# View logs for troubleshooting
./bin/rdws-logs-get --serial BS123456789 --output /tmp/logs/
```

### 4. Bulk Operations
Automate operations across multiple devices:

```bash
# List all devices and save to file
./bin/main--devices-list --json > devices.json

# Process each device (example script)
for serial in $(jq -r '.items[].serial' devices.json); do
  echo "Processing $serial"
  ./bin/rdws-health --serial $serial --json
done

# Delete old content in bulk
./bin/main-content-delete --filter "name contains 'old'" --dry-run
./bin/main-content-delete --filter "name contains 'old'" --yes
```

---

## Exit Codes

All examples use standard exit codes:
- **0**: Success
- **1**: Error (invalid arguments, authentication failed, operation failed)

Use in scripts:
```bash
if ./bin/main-device-info --serial BS123456789; then
  echo "Device is accessible"
else
  echo "Device error or offline"
fi
```

---

## Building Examples

```bash
# Build all examples
make build-examples

# Build specific example
go build -o bin/list-devices ./examples/main--devices-list

# Clean build artifacts
make clean
```

---

## Adding New Examples

When creating new examples:

1. Create directory under `examples/`
2. Follow naming convention: `main-<feature>` or `<service>-<operation>`
3. Include `--help` flag with clear usage information
4. Support common flags (`--network`, `--verbose`, `--json`, `--timeout`)
5. Use consistent error handling and exit codes
6. Add entry to this README
7. Test thoroughly with different scenarios

Example structure:
```
examples/
└── main-new-feature/
    ├── main.go
    └── README.md  # Optional for complex examples
```
