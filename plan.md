# Plan: Template-Driven Setup Configuration via Command-Line Parameters

## Goal

Replace per-example hardcoded JSON config files with a single shared template (`setups/setup-template.json`) that is rendered at runtime using command-line flags. Each example program gains the flags it needs (per the table below) and renders the template before using the resulting configuration.

## Environment Variables

These three new environment variables follow the existing `BS_` prefix convention used by `BS_CLIENT_ID`, `BS_SECRET`, and `BS_NETWORK`:

| Environment Variable | Corresponding Flag | Template Placeholder | Description |
|---|---|---|---|
| `BS_PACKAGE_NAME` | `--package-name` | `{{.PackageName}}` | Setup package identifier |
| `BS_DEVICE_NAME` | `--device-name` | `{{.DeviceName}}` | Human-readable device name |
| `BS_DEVICE_DESCRIPTION` | `--device-description` | `{{.DeviceDescription}}` | Device description text |

**Resolution order:** Command-line flag takes precedence over environment variable. If neither is set, the value is empty (or an error if the field is required for that example).

Example usage:
```bash
# Set once in the environment
export BS_PACKAGE_NAME="retail-display-v1"
export BS_DEVICE_NAME="Lobby Kiosk"
export BS_DEVICE_DESCRIPTION="Main entrance display unit"

# All examples pick up the values automatically
bdeploy-add-setup --network "Production"
bdeploy-list-setups --network "Production"
bdeploy-associate --serial ABC123456789

# Override per-invocation with a flag
bdeploy-add-setup --package-name "retail-display-v2" --network "Staging"
```

## Flag Assignment Per Example

| Example | `--package-name` / `BS_PACKAGE_NAME` | `--device-name` / `BS_DEVICE_NAME` | `--device-description` / `BS_DEVICE_DESCRIPTION` |
|---|---|---|---|
| `bdeploy-add-setup` | required | optional | optional |
| `bdeploy-list-setups` | optional (filter) | -- | -- |
| `bdeploy-get-setup` | optional (lookup) | -- | -- |
| `bdeploy-update-setup` | optional (update field) | -- | -- |
| `bdeploy-delete-setup` | optional (lookup) | -- | -- |
| `bdeploy-associate` | optional (lookup) | optional | -- |

`--` means the example does not use that field and should not gain the flag or read the env var.

## Shared Template Rendering (Phase 1)

### 1.1 Create a shared template rendering package

Create `examples/internal/setuptemplate/render.go` with:

- A `TemplateVars` struct containing all nine placeholder fields:
  ```go
  type TemplateVars struct {
      Username          string
      NetworkName       string
      PackageName       string
      RegistrationToken string
      TokenValidFrom    string
      TokenValidTo      string
      DeviceName        string
      DeviceDescription string
      GroupName         string
  }
  ```
- A `ResolveVar(flagValue string, envVar string) string` helper that returns the flag value if non-empty, otherwise falls back to `os.Getenv(envVar)`. Each example calls this to resolve `--package-name`/`BS_PACKAGE_NAME`, `--device-name`/`BS_DEVICE_NAME`, and `--device-description`/`BS_DEVICE_DESCRIPTION`.

- A `Render(templatePath string, vars TemplateVars) (*gopurple.BDeploySetupRecord, error)` function that:
  1. Reads the template file from disk.
  2. Executes `text/template` with the provided `TemplateVars`.
  3. Unmarshals the rendered JSON into a `gopurple.BDeploySetupRecord`.
  4. Returns the populated struct.

- A `DefaultTemplatePath() string` helper that returns `"setups/setup-template.json"` (can be overridden by a `--template` flag in each example).

- Constants for the environment variable names:
  ```go
  const (
      EnvPackageName       = "BS_PACKAGE_NAME"
      EnvDeviceName        = "BS_DEVICE_NAME"
      EnvDeviceDescription = "BS_DEVICE_DESCRIPTION"
  )
  ```

### 1.2 Add a common `--template` flag

Every example that renders the template gains a `--template` flag (default: `setups/setup-template.json`) so the template path is configurable.

## Per-Example Changes (Phase 2)

### 2.1 `bdeploy-add-setup`

**Current behavior:** Reads a fully-formed JSON config file as a positional argument.

**Changes:**
- Add flags: `--package-name` (required), `--device-name` (optional), `--device-description` (optional), `--template`.
- Each flag falls back to its environment variable: `--package-name` -> `BS_PACKAGE_NAME`, `--device-name` -> `BS_DEVICE_NAME`, `--device-description` -> `BS_DEVICE_DESCRIPTION`.
- Add flags for the other template variables that are currently read from the config: `--network` (or use existing `BS_NETWORK`), `--username` (or use existing `BS_CLIENT_ID`), `--group` (default: `"Default"`).
- When `--template` is provided (or no positional config file argument is given), render the template using the flag/env values to populate `TemplateVars`, then proceed with the existing workflow.
- When a positional config file argument is given (backward compatibility), use the existing `loadConfig` path unchanged.
- Validation: if template mode, require `--package-name` or `BS_PACKAGE_NAME`. If config-file mode, require `packageName` in JSON (existing behavior).

**Files modified:** `examples/bdeploy-add-setup/main.go`

### 2.2 `bdeploy-list-setups`

**Current behavior:** Has `--package` flag already for filtering.

**Changes:**
- Rename `--package` to `--package-name` for consistency (keep `--package` as an alias via `flag.StringVar`).
- Fall back to `BS_PACKAGE_NAME` if neither flag is set.
- No template rendering needed -- this example only uses `packageName` as an API filter, not as a template variable.

**Files modified:** `examples/bdeploy-list-setups/main.go`

### 2.3 `bdeploy-get-setup`

**Current behavior:** Has `--setup-name` flag that looks up by package name.

**Changes:**
- Add `--package-name` as an alias for `--setup-name` for consistency. Keep `--setup-name` working.
- Fall back to `BS_PACKAGE_NAME` if neither flag is set (still requires one of `--setup-id`, `--setup-name`/`--package-name`, or `BS_PACKAGE_NAME`).
- No template rendering needed -- this example only uses `packageName` as a lookup key.

**Files modified:** `examples/bdeploy-get-setup/main.go`

### 2.4 `bdeploy-update-setup`

**Current behavior:** Reads a JSON config file with update fields. `packageName` in the config updates the setup's package name.

**Changes:**
- Add flag: `--package-name` (optional). When provided, overrides `packageName` from the config file. Falls back to `BS_PACKAGE_NAME`.
- In `applyUpdates`, if `--package-name` flag or `BS_PACKAGE_NAME` was set, use that value instead of `config.PackageName`.
- No template rendering needed -- this example applies partial updates, not full setup creation.

**Files modified:** `examples/bdeploy-update-setup/main.go`

### 2.5 `bdeploy-delete-setup`

**Current behavior:** Has `--setup-name` flag that looks up by package name.

**Changes:**
- Add `--package-name` as an alias for `--setup-name` for consistency. Keep `--setup-name` working.
- Fall back to `BS_PACKAGE_NAME` if neither flag is set.
- No template rendering needed -- this example only uses `packageName` as a lookup key.

**Files modified:** `examples/bdeploy-delete-setup/main.go`

### 2.6 `bdeploy-associate`

**Current behavior:** Has `--setup-name` (package name lookup), `--name` (device name), `--description` (device description) flags already.

**Changes:**
- Add `--package-name` as an alias for `--setup-name` for consistency. Falls back to `BS_PACKAGE_NAME`.
- Add `--device-name` as an alias for `--name` for consistency. Falls back to `BS_DEVICE_NAME`.
- Add `--device-description` as an alias for `--description` for consistency. Falls back to `BS_DEVICE_DESCRIPTION`.
- Keep existing flag names working.
- No template rendering needed -- this example uses these as API parameters, not template variables.

**Files modified:** `examples/bdeploy-associate/main.go`

## Template Rendering Integration in `bdeploy-add-setup` (Phase 3)

This is the only example that needs full template rendering because it creates new setup records from scratch.

### 3.1 Dual-mode operation

```
# Template mode (new):
bdeploy-add-setup --package-name "retail-v1" --network "Production" --group "Lobby"

# Config file mode (existing, backward compatible):
bdeploy-add-setup config.json
```

Detection logic:
- If `flag.NArg() == 1`, use config-file mode (existing behavior).
- If `flag.NArg() == 0` and `--package-name` is set, use template mode.
- If neither, print error and usage.

### 3.2 Template mode flow

1. Populate `TemplateVars` from flags and environment using `ResolveVar`:
   - `Username` from `--username` or `BS_CLIENT_ID`
   - `NetworkName` from `--network` or `BS_NETWORK`
   - `PackageName` from `--package-name` or `BS_PACKAGE_NAME` (required -- error if both empty)
   - `GroupName` from `--group` (default: `"Default"`)
   - `DeviceName` from `--device-name` or `BS_DEVICE_NAME` (default: `""`)
   - `DeviceDescription` from `--device-description` or `BS_DEVICE_DESCRIPTION` (default: `""`)
   - `RegistrationToken`, `TokenValidFrom`, `TokenValidTo` left empty (auto-generated later in the workflow)
2. Call `setuptemplate.Render(templatePath, vars)`.
3. Continue with existing workflow (authenticate, set network, generate token, create record).

### 3.3 Handle empty token placeholders

The template has `{{.RegistrationToken}}`, `{{.TokenValidFrom}}`, `{{.TokenValidTo}}` but the existing `bdeploy-add-setup` workflow auto-generates the token after authentication. Two options:

**Option A (recommended):** Render the template with empty token values. The existing code at line 173 already checks `if setupConfig.BSNDeviceRegistrationTokenEntity == nil` and generates a token. After rendering, if the token field is empty, set `BSNDeviceRegistrationTokenEntity` to `nil` so auto-generation kicks in.

**Option B:** Use a two-pass render -- first render without token, authenticate and generate token, then re-render. Unnecessarily complex.

Go with Option A.

## Documentation Updates (Phase 4)

### 4.1 Update `docs/using-setup-templates.md`

Add a section showing how each example uses template variables via command-line flags and environment variables, with usage examples. Document the `BS_PACKAGE_NAME`, `BS_DEVICE_NAME`, and `BS_DEVICE_DESCRIPTION` environment variables alongside the existing `BS_CLIENT_ID`, `BS_SECRET`, and `BS_NETWORK`.

### 4.2 Update per-example READMEs

Update the README files in each example directory to document the new flags and show template-mode usage.

### 4.3 Update `examples/README.md`

Add a section on template-driven workflows.

## Testing (Phase 5)

### 5.1 Unit test for template rendering

Create `examples/internal/setuptemplate/render_test.go`:
- Test that rendering with all variables produces valid JSON.
- Test that rendering with empty optional variables (`DeviceName`, `DeviceDescription`) produces valid JSON with empty strings.
- Test that missing required variables cause template execution errors (if using `template.Option("missingkey=error")`).

### 5.2 Build verification

Run `make build` to confirm all six examples compile with the new flags.

### 5.3 Integration smoke test

For each modified example, run with `--help` to verify new flags appear in usage output.

## Implementation Order

1. **Phase 1** -- Shared template rendering package (`examples/internal/setuptemplate/`)
2. **Phase 2** -- Add flag aliases to the five simpler examples (list, get, update, delete, associate)
3. **Phase 3** -- Add template rendering to `bdeploy-add-setup` with dual-mode operation
4. **Phase 4** -- Documentation updates
5. **Phase 5** -- Tests and build verification

Each phase is independently testable and deployable. Phases 2 and 3 can be done in parallel since they touch different files.
