# Design: Static File Serving with Conditional Proxying

## Context
The HTTP proxy currently operates in pure proxy mode, forwarding all requests (or filtered by exact endpoint match) to a target backend. Users need a single server that can serve static web assets while also proxying API requests, eliminating the need for separate web servers.

## Goals / Non-Goals

### Goals
- Add static file serving capability without breaking existing proxy functionality
- Support combined mode: static files + backend API proxying
- Simple configuration with two new flags: `-w` (webroot) and enhanced `-e` (endpoint prefix)
- Maintain backward compatibility: no changes to behavior when new flags are not used

### Non-Goals
- Advanced web server features (virtual hosts, SSL termination, caching)
- Template rendering or server-side processing
- Authentication/authorization for static files
- URL rewriting or complex routing rules

## Decisions

### Decision 1: Routing Logic Order
**What**: Check endpoint prefix FIRST, then serve static files as fallback

**Why**: 
- Ensures API requests always go to the backend, even if a static file with the same path exists
- Predictable behavior: prefix match = proxy, no match = static file
- Prevents accidental exposure of backend paths as static files

**Implementation**:
```go
if allowedEndpoint != "" && strings.HasPrefix(r.URL.Path, allowedEndpoint) {
    // Proxy to backend
} else if webroot != "" {
    // Serve static file
} else {
    // Existing behavior (proxy all or 403)
}
```

### Decision 2: Endpoint Flag Behavior Change
**What**: Change `-e` flag from exact match to prefix match when `-w` is specified

**Why**:
- Exact match is too restrictive for API endpoints (e.g., `/backend/api/users`, `/backend/health`)
- Prefix match is more intuitive for routing: "all paths starting with `/backend`"
- Backward compatible: without `-w`, existing exact match behavior is preserved

**Alternatives considered**:
- Add a new flag `-ep` for prefix matching → Rejected: adds complexity, confusing to have two similar flags
- Always use prefix matching → Rejected: breaks backward compatibility
- Use regex patterns → Rejected: overkill for this use case, harder to configure

### Decision 3: Use Go's http.FileServer
**What**: Use standard library `http.FileServer` for static file serving

**Why**:
- Battle-tested, handles MIME types, range requests, directory listings
- No external dependencies
- Automatic index.html serving for directories
- Proper error handling (404, 403)

**Implementation**:
```go
fileServer := http.FileServer(http.Dir(webroot))
fileServer.ServeHTTP(w, r)
```

### Decision 4: Webroot Validation at Startup
**What**: Validate webroot directory exists and is a directory before starting the server

**Why**:
- Fail fast: catch configuration errors immediately
- Better user experience: clear error message vs runtime errors
- Prevents server starting in broken state

## Routing Decision Tree

```
Request arrives
    │
    ├─ Is webroot specified (-w)?
    │   │
    │   NO ──> Use existing proxy logic (exact endpoint match or proxy all)
    │   │
    │   YES ──> Is endpoint specified (-e)?
    │           │
    │           NO ──> Serve all requests as static files
    │           │
    │           YES ──> Does request path start with endpoint prefix?
    │                   │
    │                   YES ──> Proxy to backend (preserve full path)
    │                   │
    │                   NO ──> Serve as static file from webroot
```

## Example Routing Scenarios

### Scenario 1: Static file in combined mode
```
Command: ./http-proxy -p 8001 -t http://localhost:8002 -w /opt/webroot -e backend
Request: GET http://localhost:8001/index.html

Routing decision:
1. webroot specified? YES
2. endpoint specified? YES
3. path starts with "backend"? NO (/index.html does not start with backend)
4. Action: Serve /opt/webroot/index.html
```

### Scenario 2: Backend API in combined mode
```
Command: ./http-proxy -p 8001 -t http://localhost:8002 -w /opt/webroot -e backend
Request: GET http://localhost:8001/backend/index.html

Routing decision:
1. webroot specified? YES
2. endpoint specified? YES
3. path starts with "backend"? YES (/backend/index.html starts with backend)
4. Action: Proxy to http://localhost:8002/backend/index.html
```

### Scenario 3: Pure proxy mode (backward compatibility)
```
Command: ./http-proxy -p 8001 -t http://localhost:8002
Request: GET http://localhost:8001/anything

Routing decision:
1. webroot specified? NO
2. Action: Proxy to http://localhost:8002/anything (existing behavior)
```

## Implementation Notes

### Path Prefix Matching
- Use `strings.HasPrefix(r.URL.Path, allowedEndpoint)` for prefix matching
- Case-sensitive matching (standard HTTP behavior)
- No trailing slash normalization (keep it simple)
- Examples:
  - `-e backend` matches: `/backend`, `/backend/`, `/backend/api/users`
  - `-e backend` does NOT match: `/Backend`, `/api/backend`, `/backends`

### Static File Serving
- Use `http.FileServer(http.Dir(webroot))` directly
- No custom error handling needed (FileServer handles 404, 403)
- MIME types automatically detected by Go's standard library
- Directory index: serves `index.html` if present

### Webroot Validation
```go
if webroot != "" {
    info, err := os.Stat(webroot)
    if err != nil {
        log.Fatalf("Invalid webroot: %v", err)
    }
    if !info.IsDir() {
        log.Fatalf("Webroot must be a directory: %s", webroot)
    }
}
```

## Risks / Trade-offs

### Risk: Path traversal attacks
**Mitigation**: `http.FileServer` handles this automatically, preventing `../` attacks

### Risk: Endpoint prefix conflicts
**Example**: `-e /api` with static file `/api/docs.html`
**Mitigation**: Documented behavior (proxy takes precedence), users should choose non-conflicting prefixes

### Trade-off: Prefix matching changes semantics
**Impact**: `-e /backend` now matches `/backend/*` instead of exact `/backend`
**Mitigation**: Only applies when `-w` is specified, backward compatible without `-w`

## Migration Plan

### For Existing Users
1. No changes required if not using new flags
2. Existing proxy behavior unchanged
3. Existing endpoint filtering (`-e`) unchanged when `-w` is not specified

### For New Combined Mode Users
1. Create webroot directory with static files
2. Add `-w /path/to/webroot` flag
3. Add `-e /api-prefix` flag for backend routes
4. Test both static and proxy paths

### Rollback
If issues arise, users can:
1. Remove `-w` flag to revert to pure proxy mode
2. Use separate web server + proxy (previous approach)
3. No data migration needed (stateless proxy)

## Open Questions

None - design is straightforward and well-defined.

