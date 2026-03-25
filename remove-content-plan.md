# Plan: Remove Presentations and Content from gopurple SDK

Remove all code and documentation related to **presentations**, **content**, **schedules**, and **file upload** (which exclusively serves content/presentation workflows). **Device web pages are kept.**

---

## Phase 1: Delete Service Files

Delete entirely — no code in these files should be kept:

- `internal/services/presentations.go`
- `internal/services/content.go`
- `internal/services/schedules.go`
- `internal/services/upload.go`

---

## Phase 2: Delete Example Programs

Delete these directories entirely:

**Presentation examples (8):**
- `examples/main-presentation-list/`
- `examples/main-presentation-create/`
- `examples/main-presentation-info/`
- `examples/main-presentation-info-by-name/`
- `examples/main-presentation-count/`
- `examples/main-presentation-update/`
- `examples/main-presentation-delete/`
- `examples/main-presentation-delete-by-filter/`

**Content examples (4):**
- `examples/main-content-list/`
- `examples/main-content-upload/`
- `examples/main-content-download/`
- `examples/main-content-delete/`

---

## Phase 3: Delete Config Files

- `configs/presentation-example.json`
- `configs/bsn-content.json`

---

## Phase 4: Delete Documentation Files

- `docs/download-presentation.md`

---

## Phase 5: Edit `internal/types/types.go`

Remove the following type definitions:

**Content types (~lines 1480–1592):**
- `ContentFile`
- `ContentFileList`
- `ContentFileCount`
- `ContentDeleteResult`
- `ContentUploadArguments`
- `ContentUploadStatus`
- `UploadSessionResponse`
- `ChunkUploadResponse`
- `UploadRequest`
- `UploadResponse`

**Presentation types (~lines 1598–1712):**
- `Presentation`
- `PresentationCount`
- `PresentationList`
- `PresentationCreateRequest`
- `PresentationAutorun`
- `PresentationDeviceWebPage`
- `PresentationScreenSettings`
- `PresentationDeleteResult`
- `ScheduleSettings`

**Schedule types (~lines 1719+):**
- `ScheduledPresentation`
- `ScheduledPresentationList`

---

## Phase 6: Edit `gopurple.go`

**6a. Remove type aliases (~lines 325–392):**

Remove all exported type aliases for every type deleted in Phase 5:
- `ContentFile`, `ContentFileList`, `ContentFileCount`, `ContentDeleteResult`
- `ContentUploadArguments`, `ContentUploadStatus`, `UploadSessionResponse`, `ChunkUploadResponse`
- `UploadRequest`, `UploadResponse`
- `Presentation`, `PresentationCount`, `PresentationList`, `PresentationCreateRequest`
- `PresentationAutorun`, `PresentationDeviceWebPage`, `PresentationScreenSettings`, `PresentationDeleteResult`
- `ScheduleSettings`, `ScheduledPresentation`, `ScheduledPresentationList`

**6b. Remove fields from `Client` struct (~lines 507–511):**
```go
Content        services.ContentService
Upload         services.UploadService
Presentations  services.PresentationService
Schedules      services.ScheduleService
```

**6c. Remove service initialization in `New()` (~lines 534–575):**

Remove the lines that construct and assign `Content`, `Upload`, `Presentations`, and `Schedules`.

---

## Phase 7: Edit `docs/all-apis.md`

Remove entire sections covering:
- Content endpoints
- Presentation endpoints
- Schedule endpoints
- Upload API endpoints

Update any summary tables or endpoint counts that reference these sections.

---

## Phase 8: Edit `README.md`

Remove all presentation and content references, including:

- Feature bullet points mentioning presentations or content management
- The "Presentation Management" row in the feature table (~line 347)
- Code examples using `client.Content.*` or `client.Presentations.*` (~lines 587–599)
- The `presentation-example.json` entry in the configs section (~line 399)
- The Presentation Download Guide link (~line 757)
- Example binary listings for `main-content-*` and `main-presentation-*` (~lines 364–365, 449–450)
- Any directory tree entries showing removed files (~lines 701–706)

---

## Phase 9: Edit `Makefile`

Remove build targets for deleted example binaries:
- All `main-presentation-*` targets
- All `main-content-*` targets

---

## Phase 10: Edit `examples/README.md` (if it exists)

Remove any presentation/content entries from the examples index.

---

## Verification

After all phases are complete:

1. `go build ./...` — must compile with zero errors
2. `go vet ./...` — must pass clean
3. `make test` — all tests must pass
4. Confirm no stray references remain:
   ```
   grep -r "Presentation\|ContentFile\|ContentService\|ScheduleService\|UploadService" --include="*.go" .
   ```
   Should return zero results.

---

## Order of Operations

1. Delete example directories (Phase 2)
2. Delete config and doc files (Phases 3–4)
3. Delete service files (Phase 1)
4. Edit `internal/types/types.go` (Phase 5)
5. Edit `gopurple.go` (Phase 6)
6. Edit `docs/all-apis.md` (Phase 7)
7. Edit `README.md` (Phase 8)
8. Edit `Makefile` (Phase 9)
9. Edit `examples/README.md` if present (Phase 10)
10. Build and verify
