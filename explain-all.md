# Everything Needed to Create and Upload Presentations to BSN.cloud

## The Core Mental Model

A BSN.cloud presentation is a two-layer object:

1. **The presentation entity** — the outer envelope the REST API knows about (metadata, file references, status)
2. **The BPFX state** — the actual presentation content (zones, media, events, transitions, assets) embedded inside the entity as `projectFile.body`

The BPFX state is the serialized Redux store from BrightAuthor:Connected. Its structure is fully defined in `/external/bs-bpfx-specification/bpfx-specification-common.md`.

---

## The API Endpoint

One endpoint handles everything:

```
{apiBase}/2022/06/REST/Presentations/{id}/
```

| Operation | Method | URL | Body |
|-----------|--------|-----|------|
| List | GET | `Presentations/` | — |
| Get count | GET | `Presentations/count/` | — |
| Get by ID | GET | `Presentations/{id}/` | — |
| Create (draft) | POST | `Presentations/` | Full entity with `projectFile.body` inline |
| Update (draft) | PUT | `Presentations/{id}/` | Full entity with `projectFile.body` inline |
| Publish | PUT | `Presentations/{id}/` | Same as update + `autoplayFile.body` + `status: "Published"` |
| Delete | DELETE | `Presentations/{id}/` | — |

There is no separate `/publish` endpoint. Publish is the same PUT with extra fields set.

---

## The Read Path

GET returns the presentation **metadata** entity. The `projectFile` field does NOT contain the BPFX body inline on reads — it contains a CDN URL:

```json
{
  "id": 12345,
  "name": "My Presentation",
  "projectFile": {
    "path": "https://cdn.bsn.cloud/some/path/presentation.bpfx",
    ...
  }
}
```

To get the BPFX content, you must make a second unauthenticated GET to `projectFile.path`. That URL returns the full BPFX JSON blob (the `BsBpfxState`).

---

## The Write Path

POST/PUT sends the full entity with `projectFile.body` inline — the entire BPFX state as a nested JSON object:

```json
{
  "id": 0,
  "name": "My Presentation",
  "status": "Draft",
  "projectFile": {
    "id": 0,
    "type": "New",
    "name": "My Presentation.bpfx",
    "body": { <full BsBpfxState here> },
    "transferEncoding": "none"
  },
  "autoplayFile": null,
  ...
}
```

For publish, add `autoplayFile.body` and set `status: "Published"`.

---

## The BPFX Structure

The top-level BPFX object (`BsBpfxState`) has two categories of data:

- `bsdm` — the actual presentation data model consumed by the player
- Everything else (`eventMenu`, `interactiveCanvas`, `liveText`, `selection`, `screenLayoutSettings`, `meta`) — editor UI state, ignored by the player

For programmatic creation, you must supply valid values for ALL top-level keys. The player ignores the editor state, but the API/editor will not accept a file missing them.

### The `bsdm` Section (BsDmState)

This is what the player uses. Minimum required fields:

```
bsdm
  sign
    properties          (version, name, model, videoMode, orientation, etc.)
    serialPortConfigurations  (exactly 8 elements)
    gpio                (exactly 8 elements: "input" or "output")
    buttonPanels        (exactly 8 keys: bp900a-d, bp200a-d)
    irRemote
    audioSignPropertyMap
    wssDeviceSpec       ({})
    lastModifiedTime    (ISO 8601)
  zones
    zonesById           (UUID -> Zone)
    allZones            (UUID[])
    zoneLayerSequence   (UUID[])
    zoneLayersById      (UUID -> ZoneLayer)
    zoneTagIndices      (UUID -> int)
  mediaStates
    mediaStatesById     (UUID -> MediaState)
    sequencesByParentId (UUID -> UUID[])
  events                (flat UUID -> Event map, NO eventsById wrapper)
  transitions
    transitionsById     (UUID -> Transition)
    sequencesByEventId  (UUID -> UUID[])
  commands
    commandsById        (UUID -> Command)
    sequencesById       (UUID -> CommandSequence)
  dataFeeds
    feedsById           (UUID -> DataFeed)
    sourcesById         (UUID -> DataFeedSource)
  userVariables
    variablesById       (UUID -> UserVariable)  // NOTE: "variablesById" not "userVariablesById"
  htmlSites             (flat UUID -> HtmlSite map)
  liveText
    itemsById / canvasesById / layersByCanvasId / dataFeedsByGroupId
  assetMap              (flat UUID -> Asset map)
  auxiliaryFiles        (flat UUID -> AuxiliaryFile map)
  scriptPlugins         (UUID -> ScriptPlugin)
  parserPlugins         ({})
  videoModePlugins      ({})
  nodeApps              (UUID -> NodeApp)
  deviceWebPages        ({})
  linkedPresentations
  partnerProducts
  customAutorun         (0 or null)
  thumbnail             (base64 JPEG or null)
  screens               (v1.3.10 only)
  userDefinedEvents     (v1.3.10 only)
  bmapSpec              (0 or null, v1.3.10 only)
```

---

## How Content Files Become Assets in BPFX

Content files uploaded to BSN.cloud (via the Content API) must be referenced in the BPFX `assetMap` before they can appear in zones. The mapping is:

1. Upload file via Content API → get back a content entity with an integer `id` and a BSN URL
2. Create an `Asset` entry in `bsdm.assetMap` with:
   - `id`: a new UUID v4 (generated locally)
   - `name`: the filename
   - `assetType`: `"Content"`
   - `mediaType`: `"Video"` | `"Image"` | `"Audio"` etc.
   - `locator`: `"pool@"` (for BSN-hosted content)
   - `location`: `"Bsn"`
   - `networkId`: the BSN network UUID string
   - `scope`: the network scope string
3. Reference that asset UUID from content items in `mediaStatesById`

The BPFX asset UUID is local to the BPFX file — it does not need to match any ID in the BSN.cloud Content API response.

---

## Zone Structure

A zone is defined in `bsdm.zones.zonesById`. Each zone needs:

```go
Zone {
    ID                string       // UUID v4
    Name              string
    Type              ZoneType     // "VideoOrImages", "Images", "AudioOnly", "Ticker", "Control"
    InitialMediaStateID string     // UUID of the first MediaState in this zone
    NonInteractive    bool
    Position          ZonePosition { X, Y, Width, Height int; Pct bool }
    Properties        ZoneProperties  // type-specific (see below)
    Tag               string       // can be empty string
}
```

Zone types and their `Properties`:
- `VideoOrImages`: `{ viewMode, videoVolume, maxContentResolution, audioOutput, audioMode, audioMapping, audioOutputAssignments, audioMixMode, audioVolume, minimumVolume, maximumVolume, imageMode, brightness, contrast, saturation, hue, zOrderFront }`
- `Images`: `{ imageMode }`
- `AudioOnly`: `{ audioOutput, audioMode, audioMapping, audioOutputAssignments, audioMixMode, audioVolume, minimumVolume, maximumVolume }`
- `Control`: `{}` (empty)

### Zone Layers

Zones are grouped into layers in `bsdm.zones.zoneLayersById`. A layer has:
```go
ZoneLayer {
    ID             string    // UUID v4
    Type           string    // "Video", "Audio", "Invisible", "Graphics"
    ZoneSequence   []string  // ordered zone UUIDs in this layer
    // Video layers also have: ZoneLayerSpecificProperties { Type "FourK"|"HD", Index int, ... }
}
```

---

## Media States: How Content Goes in a Zone

Content items live in `bsdm.mediaStates.mediaStatesById`. Each `MediaState` holds one `ContentItem`. Playlist ordering is in `sequencesByParentId`.

```go
MediaState {
    ID          string   // UUID v4
    Name        string
    Tag         string   // can be empty
    Container   { ID string; Type int }
    // Type: 0=zone top-level, 1=playlist, 2=SuperState child, 3=MediaList item
    ContentItem ContentItem  // the actual content
}
```

`sequencesByParentId[zoneID]` gives the ordered list of top-level MediaState UUIDs in that zone. This is the "playlist".

### Content Item Types

**Video:**
```go
{ Name string; Type "Video"; AssetID UUID; Volume int; VideoDisplayMode "2D"; AutomaticallyLoop bool }
```

**Image:**
```go
{ Name string; Type "Image"; AssetID UUID; UseImageBuffer bool; VideoPlayerRequired bool; DefaultTransition string; TransitionDuration float64 }
```

**Audio:**
```go
{ Name string; Type "Audio"; AssetID UUID; Volume int }
```

---

## Events and Transitions: Minimum Viable

For a simple non-interactive looping playlist, you need one event per media state: a `MediaEnd` event that advances to the next state (or loops back).

```go
Event {
    ID           string   // UUID v4
    Name         string
    Type         "MediaEnd"
    MediaStateID string   // the MediaState this event belongs to
    Disabled     bool     // false
    Data         nil
    Action       "SeqFwd" // advance to next in sequence
}

Transition {
    ID                UUID
    Name              string
    Type              "No effect"    // always this value
    EventID           UUID           // the Event above
    TargetMediaStateID UUID          // the next MediaState to play
    Duration          0              // always 0
}
```

`bsdm.transitions.sequencesByEventId[eventID]` = `[transitionID]`

For the last item in a playlist that should loop, `TargetMediaStateID` points back to the first MediaState in the zone.

---

## Publishing vs Saving as Draft

Same PUT endpoint, different payload:

**Draft (save):**
```json
{
  "status": "Draft",
  "projectFile": { "body": <BsBpfxState>, "type": "New", ... },
  "autoplayFile": null
}
```

**Publish:**
```json
{
  "status": "Published",
  "projectFile": { "body": <BsBpfxState>, "type": "New", ... },
  "autoplayFile": { "body": <autoplay JSON>, "type": "New", "name": "autoplay.json", ... }
}
```

The `autoplayFile.body` is a separate JSON structure (not BPFX) that describes the autorun configuration for the player. Its structure needs to be determined separately — it is NOT defined in the BPFX spec.

---

## What Is Wrong in the Current SDK

### 1. `Presentation.ProjectFile` is `interface{}`

It must be a typed struct. The read shape is:
```go
type BsnPresentationFile struct {
    ID               int         `json:"id"`
    Type             string      `json:"type"`             // "New", "Existing", etc.
    Name             string      `json:"name"`             // "presentation.bpfx"
    Body             interface{} `json:"body,omitempty"`   // BsBpfxState on write; absent on read
    Path             string      `json:"path,omitempty"`   // CDN URL on read; absent on write
    TransferEncoding string      `json:"transferEncoding"` // "none"
}
```

Write uses `Body`, read uses `Path`. These are mutually exclusive.

### 2. `Publish()` uses the wrong endpoint

Current code: `POST {BSNBaseURL}/2024-10/REST/presentation/{id}/publish`

Correct: `PUT {APIVersion}/Presentations/{id}/` with `status: "Published"` and `autoplayFile.body` set.

### 3. Zone methods return `"not_supported"`

`AddVideoZone`, `AddImageZone`, `ConfigureZones`, `GetZones` all return a `"not_supported"` error. This is wrong — zone management IS possible via the API, but requires the BPFX read-modify-write cycle:

1. GET the presentation entity
2. Fetch `projectFile.path` (CDN URL, no auth) to get the full `BsBpfxState`
3. Unmarshal into typed Go structs
4. Add/modify zones in `bsdm.zones`, add media states in `bsdm.mediaStates`, add assets in `bsdm.assetMap`, wire up events and transitions
5. PUT the entity back with the modified `projectFile.body`

### 4. `PresentationZone`, `ZonePlaylist`, `PlaylistItem` types are wrong

These invented types don't match the actual BPFX structure. They should be replaced by Go structs matching the BPFX spec (see sections above).

### 5. `PresentationPublishResponse` is invented

The publish response is just the updated `Presentation` entity returned from the PUT. There is no distinct publish response type.

---

## Minimum Viable Go Types Needed

The following new types must be added to `internal/types/types.go` (or a new `bpfx.go` file):

```
BsBpfxState             top-level BPFX file
BsDmState               bsdm section
BsSignObject            bsdm.sign
BsSignProperties        bsdm.sign.properties
BsSerialPortConfig      (8 of these)
BsButtonPanels
BsZoneContainer         bsdm.zones
BsZone
BsZonePosition
BsZoneLayer
BsMediaStateContainer   bsdm.mediaStates
BsMediaState
BsContentItem           (interface{} or sum type: Video/Image/Audio/Html/etc.)
BsVideoContentItem
BsImageContentItem
BsAudioContentItem
BsEventHandler          bsdm.events value type
BsTransitionContainer   bsdm.transitions
BsTransition
BsCommandContainer      bsdm.commands (can be empty)
BsDataFeedContainer     bsdm.dataFeeds (can be empty)
BsUserVariableContainer bsdm.userVariables (can be empty)
BsAssetMap              bsdm.assetMap (flat UUID -> BsAsset)
BsAsset
BsLiveTextContainer     bsdm.liveText
BsLinkedPresentations
BsThumbnail
BsScreensContainer      (v1.3.10)
BsnPresentationFile     projectFile / autoplayFile shape
BsRgbaColor             {A, R, G, B int}
BsParameterizedValue    {Params []BsParameterizedParam}
// Editor state (required for valid BPFX but not player-relevant)
BsEventMenuState
BsInteractiveCanvasState
BsSelectionState
BsScreenLayoutSettings  (v1.3.10)
BsMeta                  (v1.3.10)
```

---

## The Full Workflow for "Create and Upload a Presentation"

```
1. Upload each media file to BSN.cloud Content API
   POST /Content/Upload → get content entity with integer ID and BSN URL

2. Build BsDmState:
   a. Generate UUID v4 for each zone, media state, event, transition, asset
   b. Populate bsdm.sign.properties (model, videoMode, name, version "1.3.10", etc.)
   c. For each content file:
      - Add BsAsset to bsdm.assetMap (locator "pool@", location "Bsn", networkId, scope)
   d. For each zone:
      - Add BsZone to bsdm.zones.zonesById
      - Add its UUID to bsdm.zones.allZones and appropriate BsZoneLayer.zoneSequence
   e. For each content item in each zone:
      - Add BsMediaState to bsdm.mediaStates.mediaStatesById
      - Add UUID to bsdm.mediaStates.sequencesByParentId[zoneUUID]
      - Add BsEvent (MediaEnd, action SeqFwd) to bsdm.events
      - Add BsTransition pointing to next MediaState (loop last -> first)
      - Add transition UUID to bsdm.transitions.sequencesByEventId[eventUUID]
   f. Set bsdm.zones.zoneLayersById, zoneLayerSequence, zoneTagIndices

3. Wrap in BsBpfxState:
   - bsdm = the BsDmState from step 2
   - Populate editor state fields with valid defaults (eventMenu, selection, interactiveCanvas,
     liveText {}, screenLayoutSettings, meta {brightAuthorVersion, buildType})

4. POST /Presentations/ with:
   {
     id: 0,
     name: "...",
     status: "Draft",
     projectFile: {
       id: 0, type: "New", name: "name.bpfx",
       body: <BsBpfxState from step 3>,
       transferEncoding: "none"
     },
     autoplayFile: null,
     ...other null/empty fields
   }
   → Response is the created Presentation entity with integer ID

5. To publish, PUT /Presentations/{id}/ with same body but:
   - status: "Published"
   - autoplayFile.body: <autoplay JSON structure>
```

---

## Unknown: Autoplay File Structure

The `autoplayFile.body` structure required for publishing is not defined in the BPFX spec and has not yet been reverse-engineered. It must be determined by:

1. Downloading a published presentation via the API and inspecting `autoplayFile.path` (CDN URL)
2. Or finding the structure in the BAConnected source (`bs-bpf-converter` or `autoplay` package)

This is a blocker for implementing `Publish()` correctly.
