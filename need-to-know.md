# Detailed Analysis: Missing Functionality for Full Presentation Automation

## Current State

### What Works ✓
1. **Content Upload**: Video and image files successfully upload to BSN.cloud
   - Video ID: 31386966
   - Image ID: 31386967
2. **Presentation Creation**: Empty presentation object created (ID: 871394)
3. **Presentation Retrieval**: Can fetch presentation metadata via API

### What's Broken ✗
**Publish fails with 404** because the presentation has no playable content structure.

---

## The Missing Pieces

### 1. Link Uploaded Content to Presentation

**Current State:**
```json
{
  "files": [
    {
      "id": 28147201,
      "name": "autoplugins.brs"  // Only system file
    }
  ]
}
```

**Required State:**
```json
{
  "files": [
    {
      "id": 28147201,
      "name": "autoplugins.brs"
    },
    {
      "id": 31386966,
      "name": "meet-brightsign.mp4",
      "mediaType": "Video",
      // ... other metadata?
    },
    {
      "id": 31386967,
      "name": "logo.png",
      "mediaType": "Image",
      // ... other metadata?
    }
  ]
}
```

**What I Need to Know:**
- ✓ API endpoint to update: `PUT /2024-10/REST/presentation/{id}` (we have this)
- ✓ How to construct file entries: Use `Presentations.Update()` (we have this)
- ❓ **What fields are required** in each file entry?
  - `id` (we have)
  - `name` (we have)
  - `mediaType`? (Video, Image)
  - `virtualPath`? (`\\Shared\\Videos\\`, `\\Shared\\Images\\`)
  - `size`? (we have from upload response?)
  - `hash`?
  - `creationDate`?
  - `path`? (S3 URL?)
  - `type`? ("Stored"?)

---

### 2. Create Presentation Structure (projectFile)

**Current State:**
```json
{
  "projectFile": null  // ❌ No presentation structure
}
```

**Required State:**
```json
{
  "projectFile": {
    // ❓ What format is this?
    // ❓ XML? JSON? Binary blob?
  }
}
```

**Critical Unknown:** What is `projectFile`?

Based on the docs mentioning `.bpfx` (BrightAuthor Project File XML), I suspect:
- It's an XML structure defining zones, playlists, and playback logic
- It's BrightSign's proprietary presentation format
- It contains the "zones" that determine what plays where and when

---

## What I Need: projectFile Structure

### Questions About projectFile Format

1. **What is the data type?**
   - Is it a JSON object with zone definitions?
   - Is it XML embedded as a string?
   - Is it a base64-encoded binary blob?
   - Is it a reference to an uploaded .bpfx file (content ID)?

2. **How are zones defined?**
   For a simple presentation with:
   - One video zone (full screen, loops)
   - One image zone (full screen, 10 second duration)

   What does the projectFile structure look like?

3. **How do zones reference content?**
   - Do they use content IDs (31386966, 31386967)?
   - Do they use file names ("meet-brightsign.mp4")?
   - Do they use content paths/URLs?

4. **What's the minimal working projectFile?**
   - What fields are absolutely required?
   - What can be omitted or left as defaults?

---

## Specific Information Needed

### A. Example of a Working Presentation

**Request:** Can you provide the JSON output of a working presentation that has:
- At least one video zone
- At least one image zone
- Has been published successfully

**Command to generate:**
```bash
../../bin/main-presentation-info --id <working-presentation-id> --network Test_Rack --json > working-presentation.json
```

This will show me:
- Complete `files` array structure
- The `projectFile` field format
- How zones are defined
- How content is referenced

### B. Documentation on .bpfx Format

**Questions:**
- Is there BSN.cloud API documentation for the projectFile structure?
- Are there example API requests showing projectFile creation?
- Is there a schema or specification for the zone structure?

### C. Alternative: Can BrightAuthor:connected Export Structure?

If creating projectFile from scratch is too complex:
- Can BrightAuthor:connected create a template presentation?
- Can we export its structure via API?
- Can we modify an existing projectFile to swap content IDs?

---

## Implementation Plan (Once I Have the Info)

### Phase 1: Add Files to Presentation
```go
// Update presentation with uploaded content
func addFilesToPresentation(ctx context.Context, client *gopurple.Client, presID int, videoID, imageID int) error {
    pres, err := client.Presentations.GetByID(ctx, presID)

    // Build file entries with proper structure (NEED FORMAT)
    files := []FileEntry{
        {ID: videoID, Name: "meet-brightsign.mp4", MediaType: "Video", /* ... */},
        {ID: imageID, Name: "logo.png", MediaType: "Image", /* ... */},
    }

    updateReq := &types.PresentationCreateRequest{
        Name: pres.Name,
        Files: files,
        // ... preserve other fields
    }

    return client.Presentations.Update(ctx, presID, updateReq)
}
```

### Phase 2: Create ProjectFile with Zones
```go
// Create simple projectFile with video and image zones
func createProjectFile(videoID, imageID int, imageDuration int) interface{} {
    // NEED TO KNOW FORMAT
    // Option 1: JSON structure?
    // Option 2: XML string?
    // Option 3: Reference to uploaded .bpfx file?

    return projectFile
}
```

### Phase 3: Update Presentation with ProjectFile
```go
// Add projectFile to presentation
func addProjectFile(ctx context.Context, client *gopurple.Client, presID int, projectFile interface{}) error {
    pres, err := client.Presentations.GetByID(ctx, presID)

    updateReq := &types.PresentationCreateRequest{
        Name: pres.Name,
        Files: pres.Files,
        ProjectFile: projectFile,  // NEED FORMAT
        // ... preserve other fields
    }

    return client.Presentations.Update(ctx, presID, updateReq)
}
```

---

## Summary

**I can implement full automation once I know:**

1. **File entry structure** - What fields are required when adding content to `files` array?
2. **ProjectFile format** - What does the `projectFile` field contain? (JSON/XML/binary?)
3. **Zone definition syntax** - How are zones structured and how do they reference content?

**Best way to get this info:**
- Export a working presentation's JSON (with zones)
- Or provide BSN.cloud API documentation on presentation structure
- Or show me an example API request that creates a presentation with zones

Once I have this information, I can implement the complete workflow in the deploy script.

---

## Action Items

1. **Provide a working presentation export:**
   ```bash
   # Find a presentation ID that works and has zones
   ../../bin/main-presentation-list --network Test_Rack

   # Export its full structure
   ../../bin/main-presentation-info --id <working-id> --network Test_Rack --json > working-presentation.json
   ```

2. **Or point me to documentation:**
   - BSN.cloud REST API docs for presentation structure
   - Example API calls showing projectFile creation
   - Schema/specification for zone definitions

3. **Or provide access to:**
   - A .bpfx file from BrightAuthor:connected
   - Documentation on the .bpfx XML format
   - Example zone structures

With any of these, I can reverse-engineer the format and implement full automation.
