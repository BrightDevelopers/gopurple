# Example Program: main-deploy-presentation

## Summary

A Go example program that performs a full end-to-end presentation deployment workflow:
parse an autoplay.json, upload referenced content to BSN.cloud, convert the autoplay to BPFX
via bs-bpfx-builder, create a presentation with the BPFX as its projectFile, publish it, and
optionally assign it to a group or device.

## Why Go Instead of the Shell Script

The existing `examples/scripts/deploy-presentation.sh` ends with 7 "Required Manual Steps"
(lines 370-381) because bash cannot carry uploaded content IDs forward into a structured BPFX
payload. The script creates an empty presentation and then tells the user to open
BrightAuthor:connected manually to add zones.

A Go program solves this because:

- Content ID resolution (check if exists → upload if not → get ID) is straightforward Go, ugly in bash
- The bs-bpfx-builder REST call is a simple `http.Post`
- BPFX post-processing (if content IDs must be injected) is just JSON manipulation
- Error handling is explicit, not `set -e` brittle
- Fits the existing examples pattern exactly

## Flags

| Flag | Description | Required |
|------|-------------|----------|
| `--autoplay` | Path to autoplay.json | Yes |
| `--name` | Presentation name | Yes |
| `--group` | Group name to assign presentation to | No |
| `--serial` | Device serial number to assign presentation to | No |
| `--bpfx-url` | bs-bpfx-builder service URL (default: `http://localhost:3000`) | No |
| `--network` | BSN.cloud network name (overrides BS_NETWORK) | No |
| `--json` | Output as JSON | No |
| `--verbose` | Show detailed progress | No |

## Execution Sequence

1. Parse `autoplay.json` — extract content file references from `zones[].playlist.states[]`
2. For each content file: search BSN.cloud by filename; upload if missing; collect `{filename → contentID}` map
3. POST `autoplay.json` body to `{bpfx-url}/autoplay2bpfx` → receive BPFX JSON blob
4. *(Conditional — see Shape B below)* Walk BPFX JSON and inject BSN content IDs where asset references appear
5. Create presentation with `ProjectFile` set to the BPFX blob
6. Publish presentation
7. If `--group` or `--serial` specified, assign presentation

## Open Question: Content ID Resolution Shape

Before writing the program, one empirical question must be answered: does BSN.cloud resolve
content references in a BPFX by **filename** or does it require **explicit content IDs**?

**Shape A — BSN resolves by filename:**
The BPFX asset references contain filenames. BSN.cloud matches them to uploaded content during
publish. Steps 1-2 are still needed to ensure content is uploaded, but step 4 (ID injection)
is not needed. The autoplay.json can be posted to bs-bpfx-builder as-is.

**Shape B — BPFX must contain BSN content IDs:**
After uploading content and collecting IDs, the generated BPFX JSON must be post-processed to
substitute BSN content IDs into the asset references before the presentation is created.

**How to answer this:** Fetch an existing working presentation via `client.Presentations.GetByID`
and inspect the `projectFile` field. If asset references contain filenames, it is Shape A. If
they contain numeric IDs or content hashes, it is Shape B.

## bs-bpfx-builder Dependency

The program requires bs-bpfx-builder running as an HTTP service. See
`howto-fix-bs-bpfx-builder.md` for the one-line Dockerfile fix required to make it start
in server mode by default.

The `--bpfx-url` flag makes this dependency configurable. If the service is not reachable the
program fails immediately with a clear error before touching BSN.cloud.

## Prerequisites Before Writing

1. Answer the Shape A/B question by inspecting a real presentation's `projectFile` field
2. Apply the bs-bpfx-builder Dockerfile fix (`CMD ["--server"]`)
3. Confirm `PresentationCreateRequest.ProjectFile` (currently `interface{}`) accepts a raw
   JSON object (not a string) when POSTed to the BSN.cloud API
