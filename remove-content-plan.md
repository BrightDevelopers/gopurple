# Plan: Remove Content Cloud APIs from gopurple

## Context

The gopurple SDK should only expose device management, diagnostics, and provisioning APIs.
All "content cloud" functionality — content file management, presentation management, scheduling,
device web pages, and file upload sessions — will not be supported. This removes a large surface
area of API that is out of scope for the library's purpose.

---

## What Gets Removed

### Service files (delete entirely)
- `internal/services/content.go` — ContentService
- `internal/services/presentations.go` — PresentationService
- `internal/services/schedules.go` — ScheduleService
- `internal/services/devicewebpages.go` — DeviceWebPageService
- `internal/services/upload.go` — UploadService

### Types in `internal/types/types.go` (delete type blocks)
All types whose names begin with:
- `Content*` — ContentFile, ContentFileList, ContentFileCount, ContentDeleteResult, ContentUploadArguments, ContentUploadStatus
- `Presentation*` — Presentation, PresentationCount, PresentationList, PresentationCreateRequest, PresentationAutorun, PresentationDeviceWebPage, PresentationScreenSettings, PresentationDeleteResult
- `ScheduledPresentation*` — ScheduledPresentation, ScheduledPresentationList
- `ScheduleSettings`
- `DeviceWebPage*` — DeviceWebPage, DeviceWebPageList
- `Upload*` — UploadRequest, UploadResponse, UploadSessionResponse
- `ChunkUploadResponse`

### `gopurple.go`
- Remove type re-exports for all types listed above (lines ~325–392)
- Remove `Content`, `Upload`, `Presentations`, `Schedules`, `DeviceWebPages` fields from `Client` struct (lines 507–511)
- Remove corresponding `services.New*Service(...)` calls in `New()` function (lines 567–571)

### Example programs (delete directories entirely)
Tracked in git:
- `examples/main-content-delete/`
- `examples/main-content-download/`
- `examples/main-content-list/`
- `examples/main-content-upload/`
- `examples/main-presentation-count/`
- `examples/main-presentation-create/`
- `examples/main-presentation-delete/`
- `examples/main-presentation-delete-by-filter/`
- `examples/main-presentation-info/`
- `examples/main-presentation-info-by-name/`
- `examples/main-presentation-list/`
- `examples/main-presentation-update/`

Untracked (new on this branch, also delete):
- `examples/main-device-assign-presentation/`
- `examples/main-device-distribution-status/`
- `examples/main-group-assign-presentation/`
- `examples/main-presentation-publish/`

### `examples/shared/flags.go`
Remove two functions:
- `GetPresentationIDWithFallback()`
- `GetContentIDWithFallback()`

Also remove their corresponding tests from `examples/shared/flags_test.go`.

### `Makefile`
Remove:
- Variable declarations for content/presentation example names (lines ~53–69):
  `EXAMPLE_MAIN_CONTENT_*`, `EXAMPLE_MAIN_PRESENTATION_*`
- All `build-main-content-*` and `build-main-presentation-*` targets
- All `run-main-content-*` and `run-main-presentation-*` targets
- References to removed targets in the `build-examples` aggregate target (line ~80) and
  the long `all` / CI test target (lines ~602–603)

---

## Execution Order

1. Delete the 5 service files
2. Remove content/presentation/schedule/upload/devicewebpage type blocks from `internal/types/types.go`
3. Edit `gopurple.go`: remove type re-exports, Client fields, and New() initializations
4. Delete all affected example directories (16 dirs)
5. Edit `examples/shared/flags.go`: remove `GetPresentationIDWithFallback` and `GetContentIDWithFallback`
6. Edit `examples/shared/flags_test.go`: remove corresponding test cases
7. Edit `Makefile`: remove all variables, targets, and references for removed examples
8. Run `make build` and `go test ./...` to verify everything compiles and tests pass

---

## Verification

```bash
make build
go test ./...
# Confirm no references remain:
grep -r "Presentation\|ContentFile\|ContentService\|ScheduleService\|DeviceWebPage\|UploadService" \
  --include="*.go" internal/ gopurple.go examples/
```
