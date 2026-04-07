# Testing Modified bdeploy Tools

All commands assume `BS_CLIENT_ID`, `BS_SECRET`, and `BS_NETWORK` are already set. Replace `"some-package"` and `SOME_SERIAL` with real values from your network.

## 1. Verify New Flags Appear in Help

```bash
./bin/bdeploy-add-setup --help 2>&1 | grep -E "package-name|device-name|device-description|template"
./bin/bdeploy-list-setups --help 2>&1 | grep -E "package-name|package"
./bin/bdeploy-get-setup --help 2>&1 | grep "package-name"
./bin/bdeploy-update-setup --help 2>&1 | grep "package-name"
./bin/bdeploy-delete-setup --help 2>&1 | grep "package-name"
./bin/bdeploy-associate --help 2>&1 | grep -E "package-name|device-name|device-description"
```

## 2. Test Environment Variable Resolution

### `bdeploy-list-setups` -- filter by `BS_PACKAGE_NAME`

```bash
# Via flag (existing behavior, renamed)
./bin/bdeploy-list-setups --package-name "some-package"

# Via old alias (backward compat)
./bin/bdeploy-list-setups --package "some-package"

# Via env var (new)
export BS_PACKAGE_NAME="some-package"
./bin/bdeploy-list-setups
unset BS_PACKAGE_NAME
```

All three should produce the same filtered output.

### `bdeploy-get-setup` -- lookup by `BS_PACKAGE_NAME`

```bash
# Via new alias
./bin/bdeploy-get-setup --package-name "some-package"

# Via old flag (backward compat)
./bin/bdeploy-get-setup --setup-name "some-package"

# Via env var
BS_PACKAGE_NAME="some-package" ./bin/bdeploy-get-setup
```

### `bdeploy-delete-setup` -- lookup by `BS_PACKAGE_NAME`

```bash
# Dry run: just confirm it finds the setup (will prompt for confirmation)
./bin/bdeploy-delete-setup --package-name "some-package"

# Via env var
BS_PACKAGE_NAME="some-package" ./bin/bdeploy-delete-setup
```

## 3. Test `bdeploy-add-setup` Template Mode

### Template mode with flags

```bash
./bin/bdeploy-add-setup \
  --package-name "test-template-v1" \
  --device-name "Test Device" \
  --device-description "Created via template mode" \
  --group "Default" \
  --verbose
```

### Template mode with env vars

```bash
export BS_PACKAGE_NAME="test-env-v1"
export BS_DEVICE_NAME="Env Device"
export BS_DEVICE_DESCRIPTION="Created via env vars"
./bin/bdeploy-add-setup --verbose
unset BS_PACKAGE_NAME BS_DEVICE_NAME BS_DEVICE_DESCRIPTION
```

### Flag overrides env var

```bash
export BS_PACKAGE_NAME="should-be-overridden"
./bin/bdeploy-add-setup --package-name "flag-wins" --verbose
unset BS_PACKAGE_NAME
```

### Config file mode (backward compat)

```bash
./bin/bdeploy-add-setup examples/bdeploy-add-setup/config.json --verbose
```

### Error case: no package name and no config file

```bash
./bin/bdeploy-add-setup
# Should print error and usage
```

### Custom template path

```bash
./bin/bdeploy-add-setup \
  --template setups/setup-template.json \
  --package-name "custom-template-test" \
  --verbose
```

## 4. Test `bdeploy-update-setup` with `--package-name` Override

```bash
# Flag overrides config file's packageName
./bin/bdeploy-update-setup \
  --setup-name "existing-setup" \
  --package-name "renamed-package" \
  examples/bdeploy-update-setup/example-config.json \
  --verbose

# Env var overrides config file's packageName
BS_PACKAGE_NAME="env-renamed" ./bin/bdeploy-update-setup \
  --setup-name "existing-setup" \
  examples/bdeploy-update-setup/example-config.json
```

## 5. Test `bdeploy-associate` with Aliases

```bash
# New alias flags
./bin/bdeploy-associate \
  --serial SOME_SERIAL \
  --package-name "some-package" \
  --device-name "Lobby Display" \
  --device-description "Main entrance" \
  --create

# Env vars
export BS_PACKAGE_NAME="some-package"
export BS_DEVICE_NAME="Lobby Display"
export BS_DEVICE_DESCRIPTION="Main entrance"
./bin/bdeploy-associate --serial SOME_SERIAL --create
unset BS_PACKAGE_NAME BS_DEVICE_NAME BS_DEVICE_DESCRIPTION

# Old flags still work
./bin/bdeploy-associate \
  --serial SOME_SERIAL \
  --setup-name "some-package" \
  --name "Lobby Display" \
  --description "Main entrance" \
  --create
```

## 6. Cleanup Test Records

After testing, remove any setup records created during testing:

```bash
./bin/bdeploy-delete-setup --package-name "test-template-v1" --force
./bin/bdeploy-delete-setup --package-name "test-env-v1" --force
./bin/bdeploy-delete-setup --package-name "flag-wins" --force
./bin/bdeploy-delete-setup --package-name "custom-template-test" --force
```
