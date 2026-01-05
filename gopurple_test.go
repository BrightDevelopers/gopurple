package gopurple

import (
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	// Test with valid credentials
	client, err := New(
		WithCredentials("test-id", "test-secret"),
		WithNetwork("test-network"),
		WithTimeout(60*time.Second),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	if client == nil {
		t.Error("Expected client to be created, got nil")
	}
	
	// Test configuration
	config := client.Config()
	if config.ClientID != "test-id" {
		t.Errorf("Expected client ID 'test-id', got '%s'", config.ClientID)
	}
	
	if config.ClientSecret != "test-secret" {
		t.Errorf("Expected client secret 'test-secret', got '%s'", config.ClientSecret)
	}
	
	if config.NetworkName != "test-network" {
		t.Errorf("Expected network name 'test-network', got '%s'", config.NetworkName)
	}
	
	if config.Timeout != 60*time.Second {
		t.Errorf("Expected timeout 60s, got %v", config.Timeout)
	}
}

func TestNewWithoutCredentials(t *testing.T) {
	// Test without credentials should fail validation
	_, err := New()
	if err == nil {
		t.Error("Expected error when creating client without credentials")
	}
	
	if !IsConfigurationError(err) {
		t.Errorf("Expected configuration error, got: %v", err)
	}
}

func TestClientMethods(t *testing.T) {
	client, err := New(
		WithCredentials("test-id", "test-secret"),
		WithNetwork("test-network"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	// Test initial state
	if client.IsAuthenticated() {
		t.Error("Expected client to not be authenticated initially")
	}
	
	if client.IsNetworkSet() {
		t.Error("Expected network to not be set initially")
	}
	
	// Test that methods requiring authentication fail appropriately
	ctx := context.Background()
	
	_, err = client.GetNetworks(ctx)
	if err == nil {
		t.Error("Expected error when getting networks without authentication")
	}
	
	if !IsAuthenticationError(err) {
		t.Errorf("Expected authentication error, got: %v", err)
	}
}

func TestErrorTypes(t *testing.T) {
	// Test that error checking functions are available
	testErr := &testError{"test error"}
	
	// These should all return false for a generic error
	if IsAuthenticationError(testErr) {
		t.Error("IsAuthenticationError should return false for generic error")
	}
	
	if IsNetworkError(testErr) {
		t.Error("IsNetworkError should return false for generic error")
	}
	
	if IsConfigurationError(testErr) {
		t.Error("IsConfigurationError should return false for generic error")
	}
	
	if IsRetryableError(testErr) {
		t.Error("IsRetryableError should return false for generic error")
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestConvenienceMethods(t *testing.T) {
	client, err := New(
		WithCredentials("test-id", "test-secret"),
		WithNetwork("test-network"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	ctx := context.Background()
	
	// Test WithAuthentication - should fail since we don't have real creds
	err = client.WithAuthentication(ctx, func() error {
		return nil
	})
	if err == nil {
		t.Error("Expected error in WithAuthentication without real credentials")
	}
	
	// Test WithNetworkContext - should fail since we don't have real creds  
	err = client.WithNetworkContext(ctx, func() error {
		return nil
	})
	if err == nil {
		t.Error("Expected error in WithNetworkContext without real credentials")
	}
}

func TestDeviceService(t *testing.T) {
	client, err := New(
		WithCredentials("test-id", "test-secret"),
		WithNetwork("test-network"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	// Test that device service is available
	if client.Devices == nil {
		t.Error("Expected Devices service to be available")
	}
	
	ctx := context.Background()
	
	// Test device listing without auth should fail
	_, err = client.Devices.List(ctx)
	if err == nil {
		t.Error("Expected error when listing devices without authentication")
	}
	
	// Test device get without auth should fail
	_, err = client.Devices.Get(ctx, "test-serial")
	if err == nil {
		t.Error("Expected error when getting device without authentication")
	}
	
	// Test device get by ID without auth should fail
	_, err = client.Devices.GetByID(ctx, 123)
	if err == nil {
		t.Error("Expected error when getting device by ID without authentication")
	}
}

func TestListOptions(t *testing.T) {
	// Test that list options are available at the top level
	opts := []ListOption{
		WithPageSize(10),
		WithMarker("test-marker"),
		WithFilter("model=XD1033"),
		WithSort("registrationDate"),
	}
	
	if len(opts) != 4 {
		t.Errorf("Expected 4 list options, got %d", len(opts))
	}
}

func TestDeviceStatusAndErrorMethods(t *testing.T) {
	client, err := New(
		WithCredentials("test-id", "test-secret"),
		WithNetwork("test-network"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	ctx := context.Background()
	
	// Test device status methods without auth should fail
	_, err = client.Devices.GetStatus(ctx, 123)
	if err == nil {
		t.Error("Expected error when getting device status without authentication")
	}
	
	_, err = client.Devices.GetStatusBySerial(ctx, "test-serial")
	if err == nil {
		t.Error("Expected error when getting device status by serial without authentication")
	}
	
	// Test device error methods without auth should fail
	_, err = client.Devices.GetErrors(ctx, 123)
	if err == nil {
		t.Error("Expected error when getting device errors without authentication")
	}
	
	_, err = client.Devices.GetErrorsBySerial(ctx, "test-serial")
	if err == nil {
		t.Error("Expected error when getting device errors by serial without authentication")
	}
}

func TestPublicTypeExports(t *testing.T) {
	// Test that all public types are available
	var _ *DeviceStatus
	var _ *DeviceError
	var _ *DeviceErrorList
	
	// Test that we can create instances
	status := &DeviceStatus{
		DeviceID: "test",
		Serial:   "test-serial",
	}
	if status.DeviceID != "test" {
		t.Error("DeviceStatus type export failed")
	}
	
	deviceErr := &DeviceError{
		ID:       1,
		Message:  "test error",
		Severity: "warning",
	}
	if deviceErr.Message != "test error" {
		t.Error("DeviceError type export failed")
	}
	
	errorList := &DeviceErrorList{
		Items: []DeviceError{*deviceErr},
	}
	if len(errorList.Items) != 1 {
		t.Error("DeviceErrorList type export failed")
	}
}