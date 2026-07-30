package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/brightdevelopers/gopurple/internal/errors"
	"github.com/brightdevelopers/gopurple/internal/types"
)

// This file covers the player routes needed to drive and observe a supervisor
// update over rDWS. They follow the same shape as the rest of rdws.go: validate,
// ensure auth + network context, build the destination-scoped URL, call, unwrap
// data.result.

// TriggerUpdateSync asks the player to check the supervisor update service for a
// new build immediately, instead of waiting for its own timer (which defaults to
// 8 hours).
//
// A false Success with a message is a legitimate answer, not a transport failure:
// the player reports "Service is disabled." when supervisor updates are switched
// off by registry, and refuses while shutting down.
func (s *rdwsService) TriggerUpdateSync(ctx context.Context, serial string) (*types.RDWSUpdateSyncResult, error) {
	if serial == "" {
		return nil, errors.NewValidationError("serial", serial, "device serial cannot be empty")
	}

	token, err := s.readyToken(ctx)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/update/sync/?destinationType=player&destinationName=%s", s.config.RDWSBaseURL, serial)

	var response types.RDWSUpdateSyncResponse
	if err := s.httpClient.PostWithAuth(ctx, token, url, struct{}{}, &response); err != nil {
		return nil, errors.NewAPIError(0, "rdws_update_sync_failed",
			fmt.Sprintf("Failed to trigger a supervisor update check on device with serial '%s'", serial), err.Error())
	}
	return &response.Data.Result, nil
}

// GetStoredSupervisors lists the supervisor builds present in the player's
// read-write supervisor directory.
//
// The firmware-bundled supervisor is not listed: this reports downloaded builds
// only. A caller comparing against what the update service offers should treat an
// absent build as "not yet downloaded".
func (s *rdwsService) GetStoredSupervisors(ctx context.Context, serial string) (*types.RDWSStoredSupervisors, error) {
	if serial == "" {
		return nil, errors.NewValidationError("serial", serial, "device serial cannot be empty")
	}

	token, err := s.readyToken(ctx)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/system/supervisors/?destinationType=player&destinationName=%s", s.config.RDWSBaseURL, serial)

	var response types.RDWSStoredSupervisorsResponse
	if err := s.httpClient.GetWithAuth(ctx, token, url, &response); err != nil {
		return nil, errors.NewAPIError(0, "rdws_stored_supervisors_failed",
			fmt.Sprintf("Failed to list stored supervisors on device with serial '%s'", serial), err.Error())
	}
	return &response.Data.Result, nil
}

// DeleteSupervisors removes downloaded supervisor builds from the player.
//
// Set request.Data.Builds to remove specific builds, or request.Data.Clear to
// remove all of them; the two are mutually exclusive and the player rejects a
// request carrying both. Only the read-write directory is touched, so the
// firmware-bundled supervisor always remains available to boot from.
//
// This is destructive: it changes what the device will run at its next boot.
func (s *rdwsService) DeleteSupervisors(ctx context.Context, serial string, request *types.RDWSDeleteSupervisorsRequest) (bool, error) {
	if serial == "" {
		return false, errors.NewValidationError("serial", serial, "device serial cannot be empty")
	}
	if request == nil {
		return false, errors.NewValidationError("request", request, "delete request cannot be nil")
	}
	if len(request.Data.Builds) == 0 && !request.Data.Clear {
		return false, errors.NewValidationError("request", request,
			"specify either builds to delete or clear")
	}
	// Rejected client-side as well as by the player, so the caller gets a clear
	// message instead of a bare 400 from the device.
	if len(request.Data.Builds) > 0 && request.Data.Clear {
		return false, errors.NewValidationError("request", request,
			"builds and clear are mutually exclusive")
	}

	token, err := s.readyToken(ctx)
	if err != nil {
		return false, err
	}

	url := fmt.Sprintf("%s/system/supervisors/delete/?destinationType=player&destinationName=%s", s.config.RDWSBaseURL, serial)

	var response types.RDWSDeleteSupervisorsResponse
	if err := s.httpClient.PostWithAuth(ctx, token, url, request, &response); err != nil {
		return false, errors.NewAPIError(0, "rdws_delete_supervisors_failed",
			fmt.Sprintf("Failed to delete supervisors on device with serial '%s'", serial), err.Error())
	}
	return response.Data.Result.Success, nil
}

// GetSystemInfo returns the player's system information, including the version of
// the supervisor that is actually RUNNING and which stored build is active.
//
// EXPERIMENTAL. This wraps GET /v1/system, which is an Internal-tier route that
// the player deliberately excludes from its published DWS documentation
// (#swagger.ignore). It is reachable over rDWS only because the relay's
// passthrough issues the request from localhost, which satisfies the Internal
// tier's local-address check. BrightSign makes no compatibility promise about this
// route: callers must tolerate it failing or disappearing, and should prefer a
// documented endpoint where one exists. It is provided because no documented
// endpoint reports the running supervisor version, and the alternative is every
// caller hand-rolling the request.
func (s *rdwsService) GetSystemInfo(ctx context.Context, serial string) (*types.RDWSSystemInfo, error) {
	if serial == "" {
		return nil, errors.NewValidationError("serial", serial, "device serial cannot be empty")
	}

	token, err := s.readyToken(ctx)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/system/?destinationType=player&destinationName=%s", s.config.RDWSBaseURL, serial)

	var response types.RDWSSystemInfoResponse
	if err := s.httpClient.GetWithAuth(ctx, token, url, &response); err != nil {
		return nil, errors.NewAPIError(0, "rdws_system_info_failed",
			fmt.Sprintf("Failed to get system information from device with serial '%s'", serial), err.Error())
	}
	return &response.Data.Result, nil
}

// GetCrashDumpFiles lists the crash dumps the player is currently holding.
//
// Distinct from GetCrashDump, which downloads a dump archive from /crash-dump.
// This is the cheap enumeration used to answer "did anything crash since I last
// looked".
func (s *rdwsService) GetCrashDumpFiles(ctx context.Context, serial string) ([]types.RDWSCrashDumpListEntry, error) {
	if serial == "" {
		return nil, errors.NewValidationError("serial", serial, "device serial cannot be empty")
	}

	token, err := s.readyToken(ctx)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/logs/crash-dumps/?destinationType=player&destinationName=%s", s.config.RDWSBaseURL, serial)

	var response types.RDWSCrashDumpListResponse
	if err := s.httpClient.GetWithAuth(ctx, token, url, &response); err != nil {
		return nil, errors.NewAPIError(0, "rdws_crash_dump_list_failed",
			fmt.Sprintf("Failed to list crash dumps on device with serial '%s'", serial), err.Error())
	}
	return parseCrashDumpList(response.Data.Result)
}

// parseCrashDumpList interprets the data.result payload of GET /logs/crash-dumps.
//
// The route returns an array. A JSON string is the player reporting an error, and
// an absent result means no dumps.
func parseCrashDumpList(raw json.RawMessage) ([]types.RDWSCrashDumpListEntry, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var entries []types.RDWSCrashDumpListEntry
	if err := json.Unmarshal(raw, &entries); err == nil {
		return entries, nil
	}
	var deviceErr string
	if err := json.Unmarshal(raw, &deviceErr); err == nil {
		return nil, errors.NewAPIError(0, "rdws_crash_dump_list_error",
			"Device reported an error listing crash dumps", deviceErr)
	}
	return nil, errors.NewAPIError(0, "rdws_crash_dump_list_parse_failed",
		"Could not parse the crash dump list", string(raw))
}

// readyToken ensures authentication and network context, then returns a token.
// Factored out because every rDWS call repeats it verbatim.
func (s *rdwsService) readyToken(ctx context.Context) (string, error) {
	if err := s.authManager.EnsureValid(ctx); err != nil {
		return "", err
	}
	if err := s.authManager.EnsureNetworkSet(ctx); err != nil {
		return "", err
	}
	return s.authManager.GetToken()
}
