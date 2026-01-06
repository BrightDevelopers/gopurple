package services

import (
	"context"
	"fmt"
	"net/url"

	"github.com/brightdevelopers/gopurple/internal/auth"
	"github.com/brightdevelopers/gopurple/internal/config"
	"github.com/brightdevelopers/gopurple/internal/errors"
	"github.com/brightdevelopers/gopurple/internal/http"
	"github.com/brightdevelopers/gopurple/internal/types"
)

// ScheduleService provides group schedule management operations.
type ScheduleService interface {
	// GetGroupSchedule retrieves all scheduled presentations for a group by name
	GetGroupSchedule(ctx context.Context, groupName string) (*types.ScheduledPresentationList, error)

	// AddScheduledPresentation adds a scheduled presentation to a group
	AddScheduledPresentation(ctx context.Context, groupName string, schedule *types.ScheduledPresentation) (*types.ScheduledPresentation, error)

	// DeleteScheduledPresentation removes a scheduled presentation from a group
	DeleteScheduledPresentation(ctx context.Context, groupName string, scheduleID int) error
}

// scheduleService implements the ScheduleService interface.
type scheduleService struct {
	config      *config.Config
	httpClient  *http.HTTPClient
	authManager *auth.AuthManager
}

// NewScheduleService creates a new schedule service.
func NewScheduleService(cfg *config.Config, httpClient *http.HTTPClient, authManager *auth.AuthManager) ScheduleService {
	return &scheduleService{
		config:      cfg,
		httpClient:  httpClient,
		authManager: authManager,
	}
}

// GetGroupSchedule retrieves all scheduled presentations for a group by name.
func (s *scheduleService) GetGroupSchedule(ctx context.Context, groupName string) (*types.ScheduledPresentationList, error) {
	if groupName == "" {
		return nil, errors.NewValidationError("groupName", groupName, "group name cannot be empty")
	}

	// Ensure we have authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return nil, err
	}

	// Build URL - name needs to be URL encoded
	scheduleURL := fmt.Sprintf("%s/%s/Groups/Regular/%s/schedule/",
		s.config.BSNBaseURL, s.config.APIVersion, url.PathEscape(groupName))

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Make the request
	var result types.ScheduledPresentationList
	if err := s.httpClient.GetWithAuth(ctx, token, scheduleURL, &result); err != nil {
		return nil, errors.NewAPIError(0, "schedule_list_failed",
			fmt.Sprintf("Failed to get schedule for group '%s'", groupName), err.Error())
	}

	return &result, nil
}

// AddScheduledPresentation adds a scheduled presentation to a group.
func (s *scheduleService) AddScheduledPresentation(ctx context.Context, groupName string, schedule *types.ScheduledPresentation) (*types.ScheduledPresentation, error) {
	if groupName == "" {
		return nil, errors.NewValidationError("groupName", groupName, "group name cannot be empty")
	}

	if schedule == nil {
		return nil, errors.NewValidationError("schedule", schedule, "schedule cannot be nil")
	}

	if schedule.PresentationID <= 0 {
		return nil, errors.NewValidationError("presentationId", schedule.PresentationID, "presentation ID must be positive")
	}

	// Ensure we have authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return nil, err
	}

	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return nil, err
	}

	// Build URL
	scheduleURL := fmt.Sprintf("%s/%s/Groups/Regular/%s/schedule/",
		s.config.BSNBaseURL, s.config.APIVersion, url.PathEscape(groupName))

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return nil, err
	}

	// Make the POST request
	var result types.ScheduledPresentation
	if err := s.httpClient.PostWithAuth(ctx, token, scheduleURL, schedule, &result); err != nil {
		return nil, errors.NewAPIError(0, "schedule_add_failed",
			fmt.Sprintf("Failed to add scheduled presentation to group '%s'", groupName), err.Error())
	}

	return &result, nil
}

// DeleteScheduledPresentation removes a scheduled presentation from a group.
func (s *scheduleService) DeleteScheduledPresentation(ctx context.Context, groupName string, scheduleID int) error {
	if groupName == "" {
		return errors.NewValidationError("groupName", groupName, "group name cannot be empty")
	}

	if scheduleID <= 0 {
		return errors.NewValidationError("scheduleID", scheduleID, "schedule ID must be positive")
	}

	// Ensure we have authentication and network context
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return err
	}

	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return err
	}

	// Build URL
	scheduleURL := fmt.Sprintf("%s/%s/Groups/Regular/%s/schedule/%d/",
		s.config.BSNBaseURL, s.config.APIVersion, url.PathEscape(groupName), scheduleID)

	// Get access token
	token, err := s.authManager.GetToken()
	if err != nil {
		return err
	}

	// Make the DELETE request
	if err := s.httpClient.DeleteWithAuth(ctx, token, scheduleURL, nil); err != nil {
		return errors.NewAPIError(0, "schedule_delete_failed",
			fmt.Sprintf("Failed to delete schedule %d from group '%s'", scheduleID, groupName), err.Error())
	}

	return nil
}
