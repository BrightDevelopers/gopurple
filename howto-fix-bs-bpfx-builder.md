# How to Fix bs-bpfx-builder as a REST Service

## Summary

`bs-bpfx-builder` already implements a REST service — `src/index.ts` runs an HTTP server on
port 3000 with a `POST /autoplay2bpfx` endpoint. No application code changes are needed.
The Dockerfile has a small bug that makes `--server` mode non-default.

## The Bug

```dockerfile
ENTRYPOINT ["/sbin/tini", "--", "node", "dist/index.js"]
CMD []
```

`CMD []` is empty, so running the container without arguments starts in CLI mode and immediately
errors because no file paths are provided. The `--server` flag must be passed explicitly at
`docker run` time (as the Makefile's `docker-run` target does).

## Fix

Change `CMD []` to `CMD ["--server"]` so the container starts the HTTP server by default:

```dockerfile
ENTRYPOINT ["/sbin/tini", "--", "node", "dist/index.js"]
CMD ["--server"]
```

With this change, `docker run bs-bpfx-builder:latest` starts the REST service on port 3000
without any additional arguments.

## Optional Improvements

1. **Update the base image** — `node:16.16.0-alpine` is EOL. Use `node:20-alpine` or `node:22-alpine`.

2. **Add a `.dockerignore`** — `COPY . .` includes `.git/`, `test/`, etc. A `.dockerignore` reduces image size:
   ```
   .git
   test
   node_modules
   dist
   *.md
   ```

## Verified Workflow (unchanged)

```bash
# Build
make docker-build

# Run REST service on port 3000
make docker-run

# Test
curl -X POST -d @test/abstract1.json http://localhost:3000/autoplay2bpfx
```

## Endpoint

| Method | Path | Body | Response |
|--------|------|------|----------|
| POST | `/autoplay2bpfx` | autoplay JSON (`abstract.json`) | BPFX JSON (`BsBpfxState`) |
