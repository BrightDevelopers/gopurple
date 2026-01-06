package auth

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/brightdevelopers/gopurple/internal/config"
	"github.com/brightdevelopers/gopurple/internal/errors"
	"github.com/brightdevelopers/gopurple/internal/http"
	"github.com/brightdevelopers/gopurple/internal/types"
)

// AuthManager handles OAuth2 authentication and network selection for BSN.cloud.
type AuthManager struct {
	config         *config.Config
	httpClient     *http.HTTPClient
	mu             sync.RWMutex
	accessToken    string
	expiresAt      time.Time
	networkSet     bool
	currentNetwork *types.Network
}

// NewAuthManager creates a new authentication manager.
func NewAuthManager(cfg *config.Config, httpClient *http.HTTPClient) *AuthManager {
	return &AuthManager{
		config:     cfg,
		httpClient: httpClient,
	}
}

// Authenticate performs OAuth2 client credentials authentication.
func (a *AuthManager) Authenticate(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if we already have a valid token
	if a.accessToken != "" && time.Until(a.expiresAt) > 30*time.Second {
		return nil
	}

	// Prepare token request
	formData := map[string]string{
		"grant_type": "client_credentials",
	}

	var tokenResp types.TokenResponse
	err := a.httpClient.PostFormWithAuth(
		ctx,
		a.config.ClientID,
		a.config.ClientSecret,
		a.config.TokenEndpoint,
		formData,
		&tokenResp,
	)
	if err != nil {
		return errors.NewAuthError("failed to get access token", err)
	}

	// Update token information
	a.accessToken = tokenResp.AccessToken
	a.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	// Reset network state since we have a new token
	a.networkSet = false
	a.currentNetwork = nil

	return nil
}

// SetNetwork sets the active network for API operations.
func (a *AuthManager) SetNetwork(ctx context.Context, networkName string) error {
	if networkName == "" {
		return errors.NewValidationError("networkName", networkName, "network name cannot be empty")
	}

	// Ensure we have a valid token
	if err := a.EnsureValid(ctx); err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if we already have this network set
	if a.networkSet && a.currentNetwork != nil && a.currentNetwork.Name == networkName {
		return nil
	}

	// Set the network via API
	networkReq := types.NetworkRequest{
		Name: networkName,
	}

	url := fmt.Sprintf("%s/%s/Self/Session/Network", a.config.BSNBaseURL, a.config.APIVersion)
	err := a.httpClient.PutWithAuth(ctx, a.accessToken, url, networkReq, nil)
	if err != nil {
		return errors.NewAuthError(fmt.Sprintf("failed to set network '%s'", networkName), err)
	}

	// Update network state
	a.networkSet = true
	a.currentNetwork = &types.Network{Name: networkName}

	return nil
}

// SetNetworkByID sets the active network by ID.
func (a *AuthManager) SetNetworkByID(ctx context.Context, networkID int) error {
	if networkID <= 0 {
		return errors.NewValidationError("networkID", networkID, "network ID must be positive")
	}

	// Ensure we have a valid token
	if err := a.EnsureValid(ctx); err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if we already have this network set
	if a.networkSet && a.currentNetwork != nil && a.currentNetwork.ID == networkID {
		return nil
	}

	// Set the network via API
	networkReq := types.NetworkRequest{
		ID: networkID,
	}

	url := fmt.Sprintf("%s/%s/Self/Session/Network", a.config.BSNBaseURL, a.config.APIVersion)
	err := a.httpClient.PutWithAuth(ctx, a.accessToken, url, networkReq, nil)
	if err != nil {
		return errors.NewAuthError(fmt.Sprintf("failed to set network ID %d", networkID), err)
	}

	// Update network state
	a.networkSet = true
	a.currentNetwork = &types.Network{ID: networkID}

	return nil
}

// GetNetworks retrieves the list of available networks.
func (a *AuthManager) GetNetworks(ctx context.Context) ([]types.Network, error) {
	// Ensure we have a valid token
	if err := a.EnsureValid(ctx); err != nil {
		return nil, err
	}

	a.mu.RLock()
	token := a.accessToken
	a.mu.RUnlock()

	url := fmt.Sprintf("%s/%s/Self/Networks", a.config.BSNBaseURL, a.config.APIVersion)
	var networks []types.Network
	err := a.httpClient.GetWithAuth(ctx, token, url, &networks)
	if err != nil {
		return nil, errors.NewAuthError("failed to get networks", err)
	}

	return networks, nil
}

// GetToken returns the current access token.
func (a *AuthManager) GetToken() (string, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.accessToken == "" {
		return "", errors.NewAuthError("not authenticated", nil)
	}

	if time.Until(a.expiresAt) <= 30*time.Second {
		return "", errors.NewAuthError("token expired", nil)
	}

	return a.accessToken, nil
}

// IsAuthenticated returns true if the client has a valid access token.
func (a *AuthManager) IsAuthenticated() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.accessToken != "" && time.Until(a.expiresAt) > 30*time.Second
}

// IsNetworkSet returns true if a network context has been established.
func (a *AuthManager) IsNetworkSet() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.networkSet
}

// GetCurrentNetwork returns the currently selected network.
func (a *AuthManager) GetCurrentNetwork() (*types.Network, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.networkSet || a.currentNetwork == nil {
		return nil, errors.NewAuthError("no network selected", nil)
	}

	// Return a copy to prevent external modifications
	network := *a.currentNetwork
	return &network, nil
}

// EnsureValid ensures that we have a valid access token.
func (a *AuthManager) EnsureValid(ctx context.Context) error {
	a.mu.RLock()
	needsRefresh := a.accessToken == "" || time.Until(a.expiresAt) <= 30*time.Second
	a.mu.RUnlock()

	if needsRefresh {
		return a.Authenticate(ctx)
	}

	return nil
}

// EnsureNetworkSet ensures that a network context has been established.
func (a *AuthManager) EnsureNetworkSet(ctx context.Context) error {
	a.mu.RLock()
	networkSet := a.networkSet
	configuredNetwork := a.config.NetworkName
	a.mu.RUnlock()

	if networkSet {
		return nil
	}

	// If we have a network configured, try to set it
	if configuredNetwork != "" {
		return a.SetNetwork(ctx, configuredNetwork)
	}

	return errors.NewAuthError("no network selected", nil)
}

// WithValidToken executes a function with a valid access token.
func (a *AuthManager) WithValidToken(ctx context.Context, fn func(token string) error) error {
	if err := a.EnsureValid(ctx); err != nil {
		return err
	}

	token, err := a.GetToken()
	if err != nil {
		return err
	}

	return fn(token)
}
