# BPFX Gap Analysis and Agent Build Path

## External Tool Reference

### `bs-bpf-converter`

Converts **BPF** (BrightAuthor Presentation File) — BrightSign's legacy XML presentation format — into a normalized JSON representation consumable by the `bsdatamodel` library.

- Parses BPF XML → intermediate JSON, normalizing zones, states, transitions, commands, metadata, user variables, live data feeds, etc.
- Handles version validation, conversion issue tracking (fatal vs. warning severity), and MRSS live data feed fixups
- Deprecated/archived repo; migrated into the BA:Connected monorepo, but still published to NPM

### `bs-bpfx-builder`

Converts a BrightSign **autoplay JSON** file (`abstract.json` / `autoplay` format) into **BPFX** format — the BrightAuthor:connected project file format.

- Spins up a Redux store with `bsdatamodel` reducers
- Parses autoplay JSON, dispatches it to build a `DmBsProjectState`
- Wraps that state with BPFX UI state (canvas, selection, event menu, etc.) to produce a `BsBpfxState`
- Can run as a **CLI** (`autoplay.json` → `output.bpfx`) or as an **HTTP server** (POST to `/autoplay2bpfx`, returns BPFX JSON)

**Summary:** `bs-bpf-converter` is legacy-XML → JSON; `bs-bpfx-builder` is autoplay-JSON → BPFX. Together they represent the format migration path from BrightAuthor's old format to BA:connected's format.

## Gap Analysis

### Gap 1: `projectFile` format — FULLY CLOSED

`projectFile` is a `.bpfx` file — pure UTF-8 JSON. It is the serialized Redux store state of
BA:Connected. Not XML, not binary, not a content ID reference. It gets POSTed as a JSON blob in
the presentation PUT body.

Top-level structure (`BsBpfxState`):
```json
{
  "meta": {
    "brightAuthorVersion": "1.72.1",
    "buildType": "Standard"
  },
  "bsdm": { ... },
  "interactiveCanvas": {
    "statePositionById": {},
    "eventDataById": {},
    "viewTransformByZoneId": {},
    "isLoaderEnabled": false
  },
  "selection": {
    "hovered": null,
    "selectionContainer": {},
    "selectionEntities": {},
    "isEventAdvancedEnabled": false,
    "activeAssetMenuTab": "ASSETS",
    "prevSelectionEntities": {}
  },
  "eventMenu": { "isOpen": true, "showTitles": true },
  "liveText": {},
  "screenLayoutSettings": {
    "selectedLayoutId": null,
    "screenCount": null,
    "selectedBezel": null,
    "isLayoutBoundaryToggled": false,
    "isSnapToCanvasToggled": false
  }
}
```

The `bsdm` block contains: sign properties, serial port config, GPIO, zones, mediaStates,
events, transitions, assetMap, and several other sections (screens, variables, playlists, etc.).

### Gap 2: `files` array fields — PARTIALLY CLOSED (one unknown remains)

The BPFX examples reveal that the `files` question is actually two separate concerns:

**Inside the BPFX `assetMap`** (fully known — confirmed from real downloaded BPFX files):

Each asset entry in `bsdm.assetMap` is keyed by a UUID (the `assetId` referenced from
mediaStates):
```json
{
  "ca754c2e-99b1-4000-ac58-d84a6b333000": {
    "id": "ca754c2e-99b1-4000-ac58-d84a6b333000",
    "name": "myfile.png",
    "path": "/Shared/Incoming/",
    "networkId": 18842293,
    "location": "Bsn",
    "assetType": "Content",
    "scope": "NetworkName",
    "locator": "bsn://Content/18842293",
    "mediaType": "Image",
    "fileSize": 207975,
    "fileHash": "E3225B7D270081D09814DD62E31980AE8ED5E913",
    "lastModifiedDate": "2024-11-06T16:57:02.952Z",
    "refCount": 1
  }
}
```

The `networkId` is the BSN Cloud content ID (the integer returned by the upload API). All of
`name`, `networkId`, `location`, `mediaType`, `fileSize`, `fileHash` come from the BSN upload
response. `path`, `scope`, `locator`, `assetType` are derived from the BSN folder structure and
network name.

**In the BSN Cloud API `files` array** (minimum fields unconfirmed):

The `files` array in the BSN presentation object is the API-side manifest. The minimum required
fields are not yet confirmed from a live example. We know `id` (BSN content ID) and `name` are
present. Whether `mediaType`, `virtualPath`, `fileSize`, or `fileHash` are required is unknown.

This is the one remaining open question. A single `GET /REST/presentation/{id}` on any working
BSN presentation closes it.

### Gap 3: Zone/content reference syntax — FULLY CLOSED

Confirmed from real BPFX examples. The reference chain is:

```
zone.initialMediaStateId
  → mediaState.id
  → mediaState.contentItem.assetId
  → assetMap[assetId].networkId  (= BSN content ID)
```

**Zone structure** (`bsdm.zones.zonesById`):
```json
{
  "cbf2f70a-0356-4000-a9a6-86b4ee9d4000": {
    "id": "cbf2f70a-0356-4000-a9a6-86b4ee9d4000",
    "name": "Zone 1",
    "type": "VideoOrImages",
    "tag": "VI1",
    "nonInteractive": false,
    "initialMediaStateId": "32767447-caf0-4000-a914-9e14c5a24000",
    "position": { "x": 0, "y": 0, "width": 1920, "height": 1080, "pct": false },
    "properties": {
      "viewMode": "Letterboxed and Centered",
      "videoVolume": 100,
      "imageMode": "Scale to Fit",
      "zOrderFront": true,
      "audioOutput": "Analog",
      "audioMode": "Stereo"
    }
  }
}
```

Zones also have a `zoneLayersById` structure (Graphics, Audio, Video layer types) and a
`zoneLayerSequence` that defines rendering order. These are generated automatically by
`bs-bpfx-builder`.

**MediaState structure** (`bsdm.mediaStates.mediaStatesById`):
```json
{
  "32767447-caf0-4000-a914-9e14c5a24000": {
    "id": "32767447-caf0-4000-a914-9e14c5a24000",
    "name": "myfile.png",
    "tag": "2",
    "container": { "id": "cbf2f70a-0356-4000-a9a6-86b4ee9d4000", "type": 0 },
    "contentItem": {
      "name": "myfile.png",
      "type": "Image",
      "assetId": "ca754c2e-99b1-4000-ac58-d84a6b333000",
      "useImageBuffer": false,
      "videoPlayerRequired": false,
      "defaultTransition": "No effect",
      "transitionDuration": 1
    }
  }
}
```

---

## `bs-bpfx-builder`: What It Is and How It Works

### Purpose

`bs-bpfx-builder` is a TypeScript utility that converts a simple "autoplay JSON" input format
into a complete, valid BPFX file. Writing BPFX by hand means constructing a full Redux store
state — potentially 100K+ lines of JSON — including serial port configurations, IR remote
definitions, button panel config, audio routing, zone layer management, and more. The autoplay
input format is ~20 lines.

### Architecture

```
src/
  index.ts            - Entry point: CLI arg parsing + HTTP server (port 3000)
  bpfxGenerator.ts    - Core conversion orchestrator
  signStateBuilder.ts - Business logic for constructing BPFX sign state
  types.ts            - TypeScript interfaces for the autoplay input format
  helpers.ts          - Cross-platform path utilities

Dockerfile            - Alpine Linux, Node.js 16
Makefile              - build, run, docker targets
```

### Input Format (autoplay JSON)

```json
{
  "BrightAuthor": {
    "version": 1,
    "meta": {
      "name": "My Presentation",
      "videoMode": "1920x1080x60p",
      "model": "XT1144",
      "size": { "width": 1920, "height": 1080 }
    },
    "zones": [{
      "id": "zone1",
      "name": "Main",
      "type": "VideoOrImages",
      "absolutePosition": { "x": 0, "y": 0, "width": 1920, "height": 1080 },
      "playlist": {
        "name": "playlist1",
        "timeoutDuration": 10,
        "states": [
          { "type": "video", "name": "intro", "fileName": "intro.mp4" },
          { "type": "image", "name": "logo", "fileName": "logo.png" }
        ]
      }
    }]
  }
}
```

### Processing Pipeline

1. Parse autoplay JSON
2. Create a Redux store with BrightSign's `@brightsign/bsdatamodel` reducer — the same library
   BA:Connected uses internally
3. Dispatch `buildSignState` thunk which:
   - Creates sign properties (validates videoMode/model, defaults to 1920x1080x60p / XT1144)
   - Forces `fullResGraphicsEnabled: true` to avoid coordinate scaling issues across models
   - Creates zones via Redux actions
   - Populates media states per zone:
     - Images: Timer event with `timeoutDuration`, NoEffect transition (1000ms)
     - Videos: no default event, volume 100, display mode 2D, no auto-loop
4. Extract final Redux store state — this *is* the `bsdm` block
5. Wrap in `BsBpfxState` container with UI metadata fields

### The Dummy Asset Problem (Critical Integration Detail)

When `bs-bpfx-builder` runs, media files are not present locally. It handles this by creating
"dummy" asset items with asset type inferred from file extension and placeholder `undefined`
values for metadata (size, hash, networkId, etc.).

This means the generated BPFX's `assetMap` entries will be incomplete. **The assetMap must be
patched after generation** before the BPFX is POSTed to BSN. Using data from the BSN upload API
responses, gopurple must fill in:

```json
{
  "networkId": 18842293,
  "location": "Bsn",
  "assetType": "Content",
  "scope": "NetworkName",
  "locator": "bsn://Content/18842293",
  "mediaType": "Image",
  "fileSize": 207975,
  "fileHash": "E3225B7D270081D09814DD62E31980AE8ED5E913",
  "lastModifiedDate": "2024-11-06T16:57:02.952Z"
}
```

The patching logic: walk `bsdm.assetMap`, match each entry by `name` to an uploaded file, inject
the BSN upload metadata.

### Two Deployment Modes

- **CLI**: `node index.js <input.json> <output.bpfx>` — batch/offline conversion
- **HTTP service**: `POST /autoplay2bpfx` → returns BPFX JSON (port 3000, Dockerized)

---

## End-to-End Architecture

```
[description / screenshot]
         |
    Claude API call
         |
    autoplay JSON
         |
  bs-bpfx-builder service (POST /autoplay2bpfx)
         |
     .bpfx JSON (with dummy assetMap entries)
         |
  patch assetMap entries with BSN upload metadata
  (networkId, fileSize, fileHash, mediaType, locator, path, scope)
         |
  gopurple:
    1. upload content files  →  BSN content IDs + metadata
    2. create presentation
    3. PUT presentation with patched BPFX as projectFile + files array
    4. publish
```

---

## Practical Approaches for Generating autoplay JSON

**Approach A — Templates (lower complexity, faster to implement)**

Pre-built parameterized autoplay JSON for common patterns:
- `fullscreen-video-loop` — single zone, looping video
- `image-slideshow` — single zone, N images with duration
- `video-then-images` — video followed by image slideshow
- `split-screen` — left zone + right zone
- `l-bar` — main content zone + ticker bar zone

A CLI tool takes `--template fullscreen-video-loop --files video.mp4 --name "My Sign"` and
produces autoplay JSON for the builder service.

**Approach B — Agent generation (higher capability, more flexible)**

An LLM (Claude API call) accepts a text description or screenshot and generates the autoplay
JSON. Zone `absolutePosition` coordinates can be derived from layout descriptions or
bounding-box analysis of a screenshot. Then convert via the `bs-bpfx-builder` service.

---

## Remaining Open Questions

1. **BSN API `files` array minimum required fields**: We know `id` and `name` are present.
   Unknown whether `mediaType`, `virtualPath`, `fileSize`, or `fileHash` are required.
   Closed by: `GET /REST/presentation/{id}` on any working presentation.

2. **assetMap `path` field**: In real BPFX files this is `/Shared/Incoming/`. Unclear if this
   must match the actual BSN folder the content was uploaded to, or if it is metadata-only.

---

## Recommended Next Steps

1. **Get a working presentation export** — closes the `files` array question and confirms the
   `path` field behavior. One `GET /REST/presentation/{id}` on any working presentation.

2. **Deploy `bs-bpfx-builder` as a local service** — Dockerfile is available. Runs
   `POST /autoplay2bpfx` on port 3000. Call it from Go to convert autoplay JSON → BPFX.

3. **Build template autoplay JSON files** for common presentation patterns. Simple JSON,
   trivially parameterized.

4. **Implement the assetMap patching step in Go** — after calling the builder service, walk
   `bsdm.assetMap` and inject BSN upload metadata (networkId, fileHash, fileSize, etc.).

5. **Wire the agent path** — once templates exist, an LLM generates autoplay JSON, gopurple
   calls the builder service, patches the assetMap, uploads files, and publishes.
