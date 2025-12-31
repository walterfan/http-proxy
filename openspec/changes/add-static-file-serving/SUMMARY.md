# OpenSpec Proposal Summary

## Change ID
`add-static-file-serving`

## Status
✅ **Validated** - Ready for review and approval

## Quick Overview
This proposal adds static file serving capability to the HTTP proxy, allowing it to serve web assets (HTML, CSS, JS) while proxying specific API paths to a backend service.

## Example Usage
```bash
# Combined mode: serve static files + proxy backend API
./http-proxy -p 8001 -t http://localhost:8002 -w /opt/webroot -e backend

# Result:
# - http://localhost:8001/index.html        → serves /opt/webroot/index.html
# - http://localhost:8001/app.js            → serves /opt/webroot/app.js
# - http://localhost:8001/backend/*         → proxies to http://localhost:8002/backend/*
# - http://localhost:8001/backend/index.html → proxies to http://localhost:8002/backend/index.html
```

## Key Features
1. **Static File Serving** (`-w` flag)
   - Serve HTML, CSS, JS, and other static assets
   - Automatic MIME type detection
   - Directory index support (index.html)

2. **Path Prefix Proxying** (enhanced `-e` flag)
   - Match path prefixes (e.g., `/backend` matches `/backend/*`)
   - Preserve full path when proxying
   - Backward compatible with exact path matching

3. **Combined Mode**
   - Serve static files by default
   - Proxy only paths matching the endpoint prefix
   - Single server for full-stack applications

4. **Backward Compatibility**
   - No breaking changes
   - Pure proxy mode unchanged when `-w` is not specified
   - Existing endpoint filtering behavior preserved

## Acceptance Tests

### Test Case 1: Static file serving in combined mode
```bash
# Given: Start proxy
./http-proxy -p 8001 -t http://localhost:8002 -w /opt/webroot -e backend

# When: Request static file
curl http://localhost:8001/index.html

# Then: Server reads /opt/webroot/index.html and returns its content
# (NOT proxied to backend)
```

### Test Case 2: Backend proxying in combined mode
```bash
# Given: Start proxy (same as above)
./http-proxy -p 8001 -t http://localhost:8002 -w /opt/webroot -e backend

# When: Request backend path
curl http://localhost:8001/backend/index.html

# Then: Request is proxied to http://localhost:8002/backend/index.html
# (NOT served from /opt/webroot/backend/index.html even if it exists)
```

## Files Created
- `openspec/changes/add-static-file-serving/proposal.md` - Rationale and impact
- `openspec/changes/add-static-file-serving/tasks.md` - Implementation checklist (20 tasks)
- `openspec/changes/add-static-file-serving/specs/http-proxy/spec.md` - Requirements and scenarios
- `openspec/changes/add-static-file-serving/design.md` - Routing logic and technical decisions

## Validation
```bash
$ openspec validate add-static-file-serving --strict
Change 'add-static-file-serving' is valid ✓
```

## Next Steps
1. **Review** - Review the proposal documents
2. **Approve** - Approve the change for implementation
3. **Implement** - Follow the tasks in `tasks.md`
4. **Archive** - After deployment, run `openspec archive add-static-file-serving`

## View Full Details
```bash
# View proposal
openspec show add-static-file-serving

# View spec deltas
openspec show add-static-file-serving --json --deltas-only

# View tasks
cat openspec/changes/add-static-file-serving/tasks.md
```

