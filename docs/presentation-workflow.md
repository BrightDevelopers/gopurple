# Presentation Workflow Guide

This guide explains how to use the gopurple SDK for BSN.cloud presentation deployment and management.

## Understanding the Architecture

### BSN.cloud Presentation Model

BSN.cloud presentations consist of two components:

1. **Presentation Structure** (zones, layouts, timing)
   - Created and edited in **BrightAuthor:connected**
   - Stored as `.bpfx` files (proprietary BrightSign format)
   - Contains zone definitions, playlists, transitions

2. **Content Files** (videos, images, audio)
   - Uploaded via API or BrightAuthor:connected
   - Referenced by presentations
   - Distributed to devices separately

### SDK Capabilities

The gopurple SDK provides **deployment and management** capabilities:

✅ **SDK Can Do:**
- Upload/download/manage content files
- Create empty presentations
- Publish presentations
- Assign presentations to devices/groups (via schedules)
- Monitor content distribution progress
- Check presentation status on devices

❌ **SDK Cannot Do:**
- Add zones to presentations (requires BrightAuthor:connected)
- Edit zone properties, playlists, transitions
- Manipulate `.bpfx` presentation structure

**Why?** The .bpfx format is proprietary and undocumented. BrightSign provides BrightAuthor:connected for presentation authoring.

## Workflow Patterns

### Pattern 1: Manual Authoring + Automated Deployment

**Best for:** Small teams, occasional updates

1. **Author** presentation in BrightAuthor:connected
2. **Publish** from BrightAuthor:connected
3. **Assign** via SDK:
   ```bash
   # Assign to entire group
   bin/main-group-assign-presentation \
     --presentation-id 123 \
     --group "Lobby Displays" \
     --network "Production"

   # Assign to specific device
   bin/main-device-assign-presentation \
     --serial BS123456789 \
     --presentation-id 123 \
     --network "Production"
   ```
4. **Monitor** distribution:
   ```bash
   bin/main-device-distribution-status \
     --serial BS123456789 \
     --network "Production"
   ```

### Pattern 2: Template-Based Automation

**Best for:** Frequent content swaps, same layout

1. **Create template** in BrightAuthor:connected (once)
   - Design zones and layout
   - Use placeholder content
   - Save and publish

2. **Upload new content** via SDK:
   ```go
   videoID, err := client.Content.Upload(ctx, "new-promo.mp4", "/videos/")
   imageID, err := client.Content.Upload(ctx, "new-logo.png", "/images/")
   ```

3. **Clone presentation** and update content IDs:
   ```go
   // Get template presentation
   template, err := client.Presentations.GetByName(ctx, "Store Display Template")

   // Clone with new name
   newPres, err := client.Presentations.Create(ctx, &gopurple.Presentation{
       Name: fmt.Sprintf("Store Display %s", time.Now().Format("2006-01-02")),
       Type: template.Type,
       // Note: Requires understanding .bpfx format to swap content IDs
   })
   ```

4. **Publish and assign** via SDK

**Limitation:** Content swapping requires parsing/modifying .bpfx structure (advanced)

### Pattern 3: Fully Automated with `deploy-presentation.sh`

**Best for:** Scripted deployments, CI/CD pipelines

The SDK includes a complete deployment script at `examples/scripts/deploy-presentation.sh`.

**Configuration** (`presentation-config.json`):
```json
{
  "presentation_name": "Store Display March 2026",
  "video_file": "../../test-media/promo.mp4",
  "image_file": "../../test-media/logo.png",
  "image_duration": 10,
  "network": "Production",
  "group": "Lobby Displays",
  "device_serial": "BS123456789"
}
```

**Run the script:**
```bash
cd examples/scripts
./deploy-presentation.sh presentation-config.json
```

**What it does:**
1. Cleans up old content with the same filenames
2. Uploads new video and image files
3. Deletes old presentation if it exists
4. Creates new empty presentation
5. Publishes presentation
6. Assigns to group or device (if specified)
7. Monitors distribution progress (if device specified)

**Limitations:**
- Creates an **empty** presentation
- You must add zones manually in BrightAuthor:connected
- Best used with template workflow

## Assignment and Scheduling

### How Assignment Works

BSN.cloud doesn't have a direct "assign presentation to device" API. Instead, presentations are assigned via **group-based scheduling**.

The SDK provides a convenience method that creates a 24/7 recurring schedule:

```go
// This creates a schedule that runs every day from 00:00 to 24:00
// Effectively a "permanent assignment"
schedule, err := client.Schedules.AssignPresentationPermanently(ctx, groupName, presentationID)
```

**Behind the scenes:**
```go
schedule := &ScheduledPresentation{
    PresentationID:      presentationID,
    IsRecurrent:         true,
    RecurrenceStartDate: today,
    RecurrenceEndDate:   today + 100 years,
    DaysOfWeek:          "EveryDay",
    StartTime:           "00:00:00",
    Duration:            "24:00:00",
}
```

### Advanced Scheduling

For time-based or dayparted schedules, use the scheduling API directly:

```go
// One-time schedule
schedule := &gopurple.ScheduledPresentation{
    PresentationID: presID,
    IsRecurrent:    false,
    EventDate:      &eventDate,
    StartTime:      "09:00:00",
    Duration:       "08:00:00",
}
client.Schedules.AddScheduledPresentation(ctx, groupName, schedule)

// Dayparting (morning vs afternoon content)
morningSchedule := &gopurple.ScheduledPresentation{
    PresentationID:      morningPresID,
    IsRecurrent:         true,
    DaysOfWeek:          "EveryDay",
    StartTime:           "06:00:00",
    Duration:            "06:00:00",
}
client.Schedules.AddScheduledPresentation(ctx, groupName, morningSchedule)
```

See `examples/main-schedule-*` programs for more scheduling examples.

## Distribution Monitoring

After assigning a presentation, monitor how content is downloading to devices:

```bash
# Check distribution status
bin/main-device-distribution-status --serial BS123456789 --network "Production"
```

**Output:**
```
Distribution Status:
  Device: BS123456789
  Progress: 8/10 files (80.0%)

Content Items:
  ✓ promo.mp4 (41003520 bytes)
  ✓ logo.png (47718 bytes)
  ⏳ background.jpg (1024000 bytes)
```

**Programmatic monitoring:**
```go
status, err := client.Devices.GetDistributionStatus(ctx, deviceSerial)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Progress: %.1f%%\n", status.PercentComplete)
fmt.Printf("Files: %d/%d downloaded\n", status.DownloadedFiles, status.TotalFiles)

for _, item := range status.ContentItems {
    fmt.Printf("  %s: %s (%d bytes)\n", item.Status, item.FileName, item.FileSize)
}

if len(status.Errors) > 0 {
    fmt.Println("Errors:")
    for _, err := range status.Errors {
        fmt.Printf("  Content %d: %s\n", err.ContentID, err.ErrorMessage)
    }
}
```

## Example Programs Reference

### Presentation Management
- `main-presentation-create` - Create empty presentation
- `main-presentation-list` - List all presentations
- `main-presentation-info` - Get presentation details by ID
- `main-presentation-info-by-name` - Get presentation by name
- `main-presentation-update` - Update presentation properties
- `main-presentation-delete` - Delete presentation by ID
- `main-presentation-delete-by-filter` - Bulk delete presentations
- `main-presentation-publish` - Publish presentation

### Assignment
- `main-group-assign-presentation` - Assign to group (24/7 schedule)
- `main-device-assign-presentation` - Assign to device's group

### Monitoring
- `main-device-distribution-status` - Check content download progress

### Scheduling (Advanced)
- `main-schedule-add-onetime` - One-time schedule for specific date
- `main-schedule-add-recurring-daily` - Daily recurring schedule
- `main-schedule-add-weekdays` - Weekday-only schedule
- `main-schedule-add-dayparting` - Time-based dayparting
- `main-schedule-list-group` - List all schedules for group
- `main-schedule-update` - Modify existing schedule
- `main-schedule-delete` - Remove schedule

See `examples/README.md` for complete program documentation.

## Best Practices

### 1. Use Naming Conventions

Include dates/versions in presentation names for easy tracking:
```
Store Display 2026-03-06
Holiday Promotion v2
Summer Sale - Week 12
```

### 2. Clean Up Old Content

The SDK provides filtering for bulk operations:
```bash
# Delete old presentations
bin/main-presentation-delete-by-filter --filter "name contains 'Old'" --yes

# Delete unused content
bin/main-content-delete --filter "name contains 'deprecated'" --yes
```

### 3. Test Before Wide Deployment

1. Create presentation in test group
2. Assign to single test device
3. Verify distribution completes successfully
4. Confirm playback on test device
5. Assign to production group

### 4. Monitor Distribution Progress

Don't assume content is ready immediately after assignment:
```bash
# Poll until distribution completes
while true; do
    status=$(bin/main-device-distribution-status --serial BS123456789 --json)
    percent=$(echo "$status" | jq -r '.percentComplete')

    if [ "$percent" == "100.0" ]; then
        echo "Distribution complete!"
        break
    fi

    echo "Progress: $percent%"
    sleep 10
done
```

### 5. Handle Errors Gracefully

```go
status, err := client.Devices.GetDistributionStatus(ctx, deviceSerial)
if err != nil {
    return fmt.Errorf("failed to get distribution status: %w", err)
}

if status.FailedFiles > 0 {
    log.Printf("Warning: %d files failed to download", status.FailedFiles)
    for _, err := range status.Errors {
        log.Printf("  Content %d: %s", err.ContentID, err.ErrorMessage)
    }
    return fmt.Errorf("distribution incomplete: %d files failed", status.FailedFiles)
}
```

## Troubleshooting

### Presentation Not Playing on Device

**Check 1:** Is the presentation assigned?
```bash
bin/main-schedule-list-group --group "Lobby Displays" --network "Production"
```

**Check 2:** Has content finished downloading?
```bash
bin/main-device-distribution-status --serial BS123456789 --network "Production"
```

**Check 3:** Is the presentation published?
```bash
bin/main-presentation-info --id 123 --network "Production"
# Check publishState field
```

**Check 4:** Does the presentation have zones/content?
- Presentations must have zones added in BrightAuthor:connected
- Empty presentations won't play

### Distribution Stalled

**Check device connectivity:**
```bash
bin/main-device-status --serial BS123456789 --network "Production"
# Check lastModifiedDate, networkStatus
```

**Check for errors:**
```bash
bin/main-device-errors --serial BS123456789 --network "Production"
```

**Verify content exists:**
```bash
bin/main-content-list --network "Production" | grep "ContentID: 123"
```

### Zone Management Errors

If you see errors like:
```
Zone management requires BrightAuthor:connected
```

This is expected. The SDK **cannot** add zones to presentations. You must:
1. Use BrightAuthor:connected desktop application
2. Open the presentation
3. Add zones through the GUI
4. Publish from BrightAuthor:connected

## Additional Resources

- [BSN.cloud API Documentation](https://docs.brightsign.biz/display/DOC/BSN.cloud+REST+API)
- [BrightAuthor:connected Download](https://www.brightsign.biz/resources/software-downloads/brightauthor-connected/)
- [Example Programs](../examples/README.md)
- [Deployment Script](../examples/scripts/README.md)
