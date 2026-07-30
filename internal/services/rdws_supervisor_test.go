package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/brightdevelopers/gopurple/internal/types"
)

// Interface conformance: the new methods must be reachable through RDWSService,
// not only on the concrete type.
func TestRDWSServiceExposesSupervisorUpdateMethods(t *testing.T) {
	var svc RDWSService = &rdwsService{}
	_ = svc
}

func TestSupervisorMethodsRejectEmptySerial(t *testing.T) {
	svc := createTestRDWSService()
	ctx := context.Background()

	if _, err := svc.TriggerUpdateSync(ctx, ""); err == nil {
		t.Error("TriggerUpdateSync: expected an error for an empty serial")
	}
	if _, err := svc.GetStoredSupervisors(ctx, ""); err == nil {
		t.Error("GetStoredSupervisors: expected an error for an empty serial")
	}
	if _, err := svc.GetSystemInfo(ctx, ""); err == nil {
		t.Error("GetSystemInfo: expected an error for an empty serial")
	}
	if _, err := svc.GetCrashDumpFiles(ctx, ""); err == nil {
		t.Error("GetCrashDumpFiles: expected an error for an empty serial")
	}
	req := &types.RDWSDeleteSupervisorsRequest{}
	req.Data.Clear = true
	if _, err := svc.DeleteSupervisors(ctx, "", req); err == nil {
		t.Error("DeleteSupervisors: expected an error for an empty serial")
	}
}

func TestSupervisorMethodsRequireAuthentication(t *testing.T) {
	svc := createTestRDWSService()
	ctx := context.Background()

	if _, err := svc.TriggerUpdateSync(ctx, "ABC123DEF456"); err == nil {
		t.Error("TriggerUpdateSync: expected an error without authentication")
	}
	if _, err := svc.GetStoredSupervisors(ctx, "ABC123DEF456"); err == nil {
		t.Error("GetStoredSupervisors: expected an error without authentication")
	}
	if _, err := svc.GetSystemInfo(ctx, "ABC123DEF456"); err == nil {
		t.Error("GetSystemInfo: expected an error without authentication")
	}
}

// builds and clear are mutually exclusive; rejecting client-side gives a clear
// message instead of a bare 400 from the device.
func TestDeleteSupervisorsValidatesRequestShape(t *testing.T) {
	svc := createTestRDWSService()
	ctx := context.Background()
	const serial = "ABC123DEF456"

	if _, err := svc.DeleteSupervisors(ctx, serial, nil); err == nil {
		t.Error("expected an error for a nil request")
	}

	empty := &types.RDWSDeleteSupervisorsRequest{}
	if _, err := svc.DeleteSupervisors(ctx, serial, empty); err == nil {
		t.Error("expected an error when neither builds nor clear is set")
	}

	both := &types.RDWSDeleteSupervisorsRequest{}
	both.Data.Builds = []string{"2024-01-09T20-31-11.822Z"}
	both.Data.Clear = true
	if _, err := svc.DeleteSupervisors(ctx, serial, both); err == nil {
		t.Error("expected an error when builds and clear are both set")
	}
}

func TestDeleteSupervisorsRequestSerialisesUnderData(t *testing.T) {
	req := &types.RDWSDeleteSupervisorsRequest{}
	req.Data.Builds = []string{"2024-01-09T20-31-11.822Z"}

	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var round struct {
		Data struct {
			Builds []string `json:"builds"`
			Clear  *bool    `json:"clear"`
		} `json:"data"`
	}
	if err := json.Unmarshal(b, &round); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(round.Data.Builds) != 1 || round.Data.Builds[0] != "2024-01-09T20-31-11.822Z" {
		t.Errorf("builds did not round-trip under data: %s", b)
	}
	// clear must be omitted rather than sent as false, or the player sees both
	// fields and rejects the request.
	if round.Data.Clear != nil {
		t.Errorf("clear should be omitted when unset: %s", b)
	}
}

func TestSystemInfoParsesRunningSupervisorAndActiveBuild(t *testing.T) {
	// Shaped after the player's own reply: supervisors_available[].name is the full
	// path plus the version in parentheses.
	const body = `{
	  "route": "/v1/system", "method": "GET",
	  "data": {"result": {
	    "firmware": {"version": {"major": 9, "minor": 1, "patch": 136, "build": 1}},
	    "bootstrap": {
	      "version": {"major": 1, "minor": 0, "patch": 2},
	      "autorun_drive": "/storage/sd",
	      "supervisors_available": [
	        {"name": "/var/lib/brightsign/bootstrap/supervisors/2024-03-12T23-43-57.806Z/2024-03-12T23-43-57.806Z.js (v 2.0.20.2)", "active": true},
	        {"name": "/usr/local/bin/supervisors/2024-01-09T20-31-11.822Z/2024-01-09T20-31-11.822Z.js (v 2.0.18.1)", "active": false}
	      ]
	    },
	    "supervisor": {"version": {"major": 2, "minor": 0, "patch": 20, "build": 2}, "dir_rw": "/var/lib/brightsign/bootstrap/supervisors"}
	  }}
	}`

	var resp types.RDWSSystemInfoResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	info := resp.Data.Result

	if got := info.Supervisor.Version.String(); got != "2.0.20.2" {
		t.Errorf("running supervisor version: got %q, want 2.0.20.2", got)
	}
	if got := info.Firmware.Version.String(); got != "9.1.136.1" {
		t.Errorf("firmware version: got %q", got)
	}
	if got := len(info.Bootstrap.SupervisorsAvailable); got != 2 {
		t.Fatalf("supervisors_available: got %d entries, want 2", got)
	}

	active, ok := info.ActiveSupervisor()
	if !ok {
		t.Fatal("expected an active supervisor")
	}
	// The path is returned verbatim, and it is what discloses whether the running
	// build came from the read-write directory or the firmware-bundled one.
	if want := "/var/lib/brightsign/bootstrap/supervisors/"; !contains(active.Name, want) {
		t.Errorf("active supervisor name %q should carry its full path", active.Name)
	}
}

func TestActiveSupervisorAbsentWhenNoneMarked(t *testing.T) {
	var info types.RDWSSystemInfo
	info.Bootstrap.SupervisorsAvailable = []types.RDWSSupervisorRef{{Name: "a", Active: false}}
	if _, ok := info.ActiveSupervisor(); ok {
		t.Error("no supervisor is marked active, so none should be reported")
	}
}

func TestComponentVersionStringOmitsZeroBuild(t *testing.T) {
	tests := []struct {
		in   types.RDWSComponentVersion
		want string
	}{
		{types.RDWSComponentVersion{Major: 2, Minor: 0, Patch: 20, Build: 2}, "2.0.20.2"},
		{types.RDWSComponentVersion{Major: 2, Minor: 0, Patch: 20}, "2.0.20"},
		{types.RDWSComponentVersion{}, "0.0.0"},
	}
	for _, tt := range tests {
		if got := tt.in.String(); got != tt.want {
			t.Errorf("%+v: got %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestStoredSupervisorsParsesBuildList(t *testing.T) {
	const body = `{"data":{"result":{"success":true,"builds":["2024-01-09T20-31-11.822Z","2024-03-12T23-43-57.806Z"]}}}`
	var resp types.RDWSStoredSupervisorsResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !resp.Data.Result.Success {
		t.Error("success should be true")
	}
	if got := len(resp.Data.Result.Builds); got != 2 {
		t.Errorf("builds: got %d, want 2", got)
	}
}

// A false success carrying a message is a legitimate answer — the update service
// being disabled by registry, for instance — not a transport failure.
func TestUpdateSyncCarriesTheRefusalReason(t *testing.T) {
	const body = `{"data":{"result":{"success":false,"message":"Service is disabled."}}}`
	var resp types.RDWSUpdateSyncResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Data.Result.Success {
		t.Error("success should be false")
	}
	if resp.Data.Result.Message != "Service is disabled." {
		t.Errorf("message: got %q", resp.Data.Result.Message)
	}
}

func TestParseCrashDumpList(t *testing.T) {
	entries, err := parseCrashDumpList(json.RawMessage(`[{"fileName":"dump-1","ctime":"2026-07-30T00:00:00Z"}]`))
	if err != nil {
		t.Fatalf("array payload: %v", err)
	}
	if len(entries) != 1 || entries[0].FileName != "dump-1" {
		t.Errorf("unexpected entries: %+v", entries)
	}

	if entries, err := parseCrashDumpList(json.RawMessage(`[]`)); err != nil || len(entries) != 0 {
		t.Errorf("empty array should yield no entries and no error: %v / %+v", err, entries)
	}
	if entries, err := parseCrashDumpList(nil); err != nil || entries != nil {
		t.Errorf("absent result should yield no entries and no error: %v / %+v", err, entries)
	}
	if _, err := parseCrashDumpList(json.RawMessage(`"cannot read crash dumps"`)); err == nil {
		t.Error("a string payload is the device reporting an error and must surface as one")
	}
	if _, err := parseCrashDumpList(json.RawMessage(`{"unexpected":true}`)); err == nil {
		t.Error("an unparseable payload must surface as an error")
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (func() bool {
		for i := 0; i+len(needle) <= len(haystack); i++ {
			if haystack[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	})()
}
