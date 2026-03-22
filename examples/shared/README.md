# Shared Examples Package

This package provides common helper functions for the example programs in this repository.

## Purpose

The `shared` package centralizes common patterns used across example programs, reducing code duplication and providing consistent behavior for:

- Flag parsing with environment variable fallbacks
- Input validation
- Error handling

## Functions

### GetSerialWithFallback

```go
func GetSerialWithFallback(serialFlag string) (string, error)
```

Returns the serial number from the flag value or `BS_SERIAL` environment variable. Returns an error if neither is provided.

**Usage:**
```go
serial, err := shared.GetSerialWithFallback(*serialFlag)
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
    flag.Usage()
    os.Exit(1)
}
```

### GetDeviceIDWithFallback

```go
func GetDeviceIDWithFallback(idFlag int) (int, error)
```

Returns the device ID from the flag value or `BS_DEVICE_ID` environment variable. Returns an error if neither is provided or if the environment variable is invalid.

**Usage:**
```go
deviceID, err := shared.GetDeviceIDWithFallback(*idFlag)
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
    flag.Usage()
    os.Exit(1)
}
```

### GetPresentationIDWithFallback

```go
func GetPresentationIDWithFallback(idFlag int) (int, error)
```

Returns the presentation ID from the flag value or `BS_PRESENTATION_ID` environment variable.

### GetContentIDWithFallback

```go
func GetContentIDWithFallback(idFlag int) (int, error)
```

Returns the content ID from the flag value or `BS_CONTENT_ID` environment variable.

### GetGroupIDWithFallback

```go
func GetGroupIDWithFallback(idFlag int) (int, error)
```

Returns the group ID from the flag value or `BS_GROUP_ID` environment variable.

## Environment Variables

The following environment variables are supported:

- `BS_SERIAL` - Device serial number
- `BS_DEVICE_ID` - Device ID (numeric)
- `BS_PRESENTATION_ID` - Presentation ID (numeric)
- `BS_CONTENT_ID` - Content file ID (numeric)
- `BS_GROUP_ID` - Group ID (numeric)

**Priority:** Flag values always take precedence over environment variables.

## Example Usage

### Device Operations with Serial Number

```go
package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/brightdevelopers/gopurple/examples/shared"
)

func main() {
    serialFlag := flag.String("serial", "", "Device serial number")
    flag.Parse()

    serial, err := shared.GetSerialWithFallback(*serialFlag)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
        flag.Usage()
        os.Exit(1)
    }

    // Use serial for API calls
    fmt.Printf("Operating on device: %s\n", serial)
}
```

### Multiple ID Fallbacks

```go
// Get serial number and device ID from flags or environment variables
serial, serialErr := shared.GetSerialWithFallback(*serialFlag)
deviceID, idErr := shared.GetDeviceIDWithFallback(*idFlag)

// Validate input
if serialErr != nil && idErr != nil {
    fmt.Fprintf(os.Stderr, "Error: Must specify either --serial/BS_SERIAL or --id/BS_DEVICE_ID\n\n")
    flag.Usage()
    os.Exit(1)
}

if serialErr == nil && idErr == nil {
    fmt.Fprintf(os.Stderr, "Error: Cannot specify both --serial and --id\n\n")
    flag.Usage()
    os.Exit(1)
}

// Use whichever was provided
if serialErr == nil {
    // Use serial
} else {
    // Use deviceID
}
```

## Testing

Run tests with:

```bash
go test -v
```

All functions are fully tested with multiple test cases covering:
- Flag precedence over environment variables
- Environment variable fallback
- Error handling for missing values
- Invalid environment variable values
