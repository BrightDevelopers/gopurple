package services

import (
	"context"
	"testing"

	"github.com/brightdevelopers/gopurple/internal/auth"
	"github.com/brightdevelopers/gopurple/internal/config"
	"github.com/brightdevelopers/gopurple/internal/http"
	"github.com/brightdevelopers/gopurple/internal/types"
)

func createTestScheduleService() ScheduleService {
	cfg := config.DefaultConfig()
	cfg.ClientID = "test-id"
	cfg.ClientSecret = "test-secret"

	httpClient := http.NewHTTPClient(cfg)
	authManager := auth.NewAuthManager(cfg, httpClient)

	return NewScheduleService(cfg, httpClient, authManager)
}

func TestScheduleService_GetGroupSchedule(t *testing.T) {
	service := createTestScheduleService()
	ctx := context.Background()

	// Test with empty group name
	_, err := service.GetGroupSchedule(ctx, "")
	if err == nil {
		t.Error("Expected error when getting schedule with empty group name")
	}

	// Test without authentication should fail
	_, err = service.GetGroupSchedule(ctx, "TestGroup")
	if err == nil {
		t.Error("Expected error when getting schedule without authentication")
	}
}

func TestScheduleService_GetScheduledPresentation(t *testing.T) {
	service := createTestScheduleService()
	ctx := context.Background()

	// Test with empty group name
	_, err := service.GetScheduledPresentation(ctx, "", 123)
	if err == nil {
		t.Error("Expected error when getting scheduled presentation with empty group name")
	}

	// Test with zero schedule ID
	_, err = service.GetScheduledPresentation(ctx, "TestGroup", 0)
	if err == nil {
		t.Error("Expected error when getting scheduled presentation with zero schedule ID")
	}

	// Test with negative schedule ID
	_, err = service.GetScheduledPresentation(ctx, "TestGroup", -1)
	if err == nil {
		t.Error("Expected error when getting scheduled presentation with negative schedule ID")
	}

	// Test without authentication should fail
	_, err = service.GetScheduledPresentation(ctx, "TestGroup", 123)
	if err == nil {
		t.Error("Expected error when getting scheduled presentation without authentication")
	}
}

func TestScheduleService_AddScheduledPresentation(t *testing.T) {
	service := createTestScheduleService()
	ctx := context.Background()

	// Test with empty group name
	schedule := &types.ScheduledPresentation{
		PresentationID: 123,
	}
	_, err := service.AddScheduledPresentation(ctx, "", schedule)
	if err == nil {
		t.Error("Expected error when adding scheduled presentation with empty group name")
	}

	// Test with nil schedule
	_, err = service.AddScheduledPresentation(ctx, "TestGroup", nil)
	if err == nil {
		t.Error("Expected error when adding nil scheduled presentation")
	}

	// Test with zero presentation ID
	schedule = &types.ScheduledPresentation{
		PresentationID: 0,
	}
	_, err = service.AddScheduledPresentation(ctx, "TestGroup", schedule)
	if err == nil {
		t.Error("Expected error when adding scheduled presentation with zero presentation ID")
	}

	// Test with negative presentation ID
	schedule = &types.ScheduledPresentation{
		PresentationID: -1,
	}
	_, err = service.AddScheduledPresentation(ctx, "TestGroup", schedule)
	if err == nil {
		t.Error("Expected error when adding scheduled presentation with negative presentation ID")
	}

	// Test without authentication should fail
	schedule = &types.ScheduledPresentation{
		PresentationID: 123,
	}
	_, err = service.AddScheduledPresentation(ctx, "TestGroup", schedule)
	if err == nil {
		t.Error("Expected error when adding scheduled presentation without authentication")
	}
}

func TestScheduleService_UpdateScheduledPresentation(t *testing.T) {
	service := createTestScheduleService()
	ctx := context.Background()

	// Test with empty group name
	schedule := &types.ScheduledPresentation{
		PresentationID: 123,
	}
	_, err := service.UpdateScheduledPresentation(ctx, "", 456, schedule)
	if err == nil {
		t.Error("Expected error when updating scheduled presentation with empty group name")
	}

	// Test with zero schedule ID
	_, err = service.UpdateScheduledPresentation(ctx, "TestGroup", 0, schedule)
	if err == nil {
		t.Error("Expected error when updating scheduled presentation with zero schedule ID")
	}

	// Test with negative schedule ID
	_, err = service.UpdateScheduledPresentation(ctx, "TestGroup", -1, schedule)
	if err == nil {
		t.Error("Expected error when updating scheduled presentation with negative schedule ID")
	}

	// Test with nil schedule
	_, err = service.UpdateScheduledPresentation(ctx, "TestGroup", 456, nil)
	if err == nil {
		t.Error("Expected error when updating with nil scheduled presentation")
	}

	// Test with zero presentation ID
	schedule = &types.ScheduledPresentation{
		PresentationID: 0,
	}
	_, err = service.UpdateScheduledPresentation(ctx, "TestGroup", 456, schedule)
	if err == nil {
		t.Error("Expected error when updating scheduled presentation with zero presentation ID")
	}

	// Test with negative presentation ID
	schedule = &types.ScheduledPresentation{
		PresentationID: -1,
	}
	_, err = service.UpdateScheduledPresentation(ctx, "TestGroup", 456, schedule)
	if err == nil {
		t.Error("Expected error when updating scheduled presentation with negative presentation ID")
	}

	// Test without authentication should fail
	schedule = &types.ScheduledPresentation{
		PresentationID: 123,
	}
	_, err = service.UpdateScheduledPresentation(ctx, "TestGroup", 456, schedule)
	if err == nil {
		t.Error("Expected error when updating scheduled presentation without authentication")
	}
}

func TestScheduleService_DeleteScheduledPresentation(t *testing.T) {
	service := createTestScheduleService()
	ctx := context.Background()

	// Test with empty group name
	err := service.DeleteScheduledPresentation(ctx, "", 123)
	if err == nil {
		t.Error("Expected error when deleting scheduled presentation with empty group name")
	}

	// Test with zero schedule ID
	err = service.DeleteScheduledPresentation(ctx, "TestGroup", 0)
	if err == nil {
		t.Error("Expected error when deleting scheduled presentation with zero schedule ID")
	}

	// Test with negative schedule ID
	err = service.DeleteScheduledPresentation(ctx, "TestGroup", -1)
	if err == nil {
		t.Error("Expected error when deleting scheduled presentation with negative schedule ID")
	}

	// Test without authentication should fail
	err = service.DeleteScheduledPresentation(ctx, "TestGroup", 123)
	if err == nil {
		t.Error("Expected error when deleting scheduled presentation without authentication")
	}
}
