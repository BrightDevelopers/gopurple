package shared

import (
	"fmt"
	"os"
	"strconv"
)

// GetSerialWithFallback returns the serial number from flag or BS_SERIAL environment variable.
// Returns an error if neither is provided.
func GetSerialWithFallback(serialFlag string) (string, error) {
	if serialFlag != "" {
		return serialFlag, nil
	}

	if envSerial := os.Getenv("BS_SERIAL"); envSerial != "" {
		return envSerial, nil
	}

	return "", fmt.Errorf("serial number required: use --serial flag or set BS_SERIAL environment variable")
}

// GetDeviceIDWithFallback returns the device ID from flag or BS_DEVICE_ID environment variable.
// Returns an error if neither is provided or if the environment variable is invalid.
func GetDeviceIDWithFallback(idFlag int) (int, error) {
	if idFlag != 0 {
		return idFlag, nil
	}

	if envID := os.Getenv("BS_DEVICE_ID"); envID != "" {
		id, err := strconv.Atoi(envID)
		if err != nil {
			return 0, fmt.Errorf("invalid BS_DEVICE_ID: %w", err)
		}
		return id, nil
	}

	return 0, fmt.Errorf("device ID required: use --id flag or set BS_DEVICE_ID environment variable")
}

// GetGroupIDWithFallback returns the group ID from flag or BS_GROUP_ID environment variable.
// Returns an error if neither is provided or if the environment variable is invalid.
func GetGroupIDWithFallback(idFlag int) (int, error) {
	if idFlag != 0 {
		return idFlag, nil
	}

	if envID := os.Getenv("BS_GROUP_ID"); envID != "" {
		id, err := strconv.Atoi(envID)
		if err != nil {
			return 0, fmt.Errorf("invalid BS_GROUP_ID: %w", err)
		}
		return id, nil
	}

	return 0, fmt.Errorf("group ID required: use --id flag or set BS_GROUP_ID environment variable")
}
