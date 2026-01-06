package services

import (
	"context"
	"fmt"

	"github.com/brightdevelopers/gopurple/internal/auth"
	"github.com/brightdevelopers/gopurple/internal/config"
	"github.com/brightdevelopers/gopurple/internal/errors"
	"github.com/brightdevelopers/gopurple/internal/http"
	"github.com/brightdevelopers/gopurple/internal/types"
)

// ProvisioningService provides player provisioning operations including token generation.
type ProvisioningService interface {
	GenerateDeviceToken(ctx context.Context) (*types.BSNTokenEntity, error)
	ValidateDeviceToken(ctx context.Context, token string) (*types.BSNTokenEntity, error)
}

// provisioningService implements the ProvisioningService interface.
type provisioningService struct {
	config      *config.Config
	httpClient  *http.HTTPClient
	authManager *auth.AuthManager
}

// NewProvisioningService creates a new provisioning service.
func NewProvisioningService(cfg *config.Config, httpClient *http.HTTPClient, authManager *auth.AuthManager) ProvisioningService {
	return &provisioningService{
		config:      cfg,
		httpClient:  httpClient,
		authManager: authManager,
	}
}

// GenerateDeviceToken generates a new device registration token for the current network.
//
// This token allows BrightSign players to register themselves with BSN.cloud.
// The token is valid for 2 years by default and has "cert" scope, allowing
// multiple devices to use the same token.
//
// Required scope: bsn.api.main.devices.setups.token.create
//
// The returned token should be embedded in B-Deploy setup records to enable
// player registration during provisioning.
func (s *provisioningService) GenerateDeviceToken(ctx context.Context) (*types.BSNTokenEntity, error) {
	// Ensure we have authentication
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Build the provisioning token endpoint
	// Using the 2020/10 API version which is documented and stable
	tokenURL := "https://api.bsn.cloud/2020/10/REST/Provisioning/Setups/Tokens/"

	// Make the API request - POST to generate token
	var response types.BSNTokenEntity
	err = s.httpClient.PostWithAuth(ctx, token, tokenURL, nil, &response)
	if err != nil {
		return nil, errors.NewAPIError(0, "token_generation_failed",
			"Failed to generate device registration token", err.Error())
	}

	// Validate response has required fields
	if response.Token == "" {
		return nil, errors.NewAPIError(0, "invalid_token_response",
			"API returned empty token", "token field is empty")
	}

	if response.Scope == "" {
		// Default to "cert" if not provided
		response.Scope = "cert"
	}

	return &response, nil
}

// ValidateDeviceToken validates a device registration token and retrieves its metadata.
//
// Required scope: bsn.api.main.devices.setups.token.validate
func (s *provisioningService) ValidateDeviceToken(ctx context.Context, tokenValue string) (*types.BSNTokenEntity, error) {
	if tokenValue == "" {
		return nil, errors.NewValidationError("token", tokenValue, "token cannot be empty")
	}

	// Ensure we have authentication
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Build the token validation endpoint
	validateURL := fmt.Sprintf("https://api.bsn.cloud/2020/10/REST/Provisioning/Setups/Tokens/%s/", tokenValue)

	// Make the API request
	var response types.BSNTokenEntity
	err = s.httpClient.GetWithAuth(ctx, token, validateURL, &response)
	if err != nil {
		return nil, errors.NewAPIError(0, "token_validation_failed",
			"Failed to validate device registration token", err.Error())
	}

	return &response, nil
}
