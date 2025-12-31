# Acceptance Tests for Static File Serving with Conditional Proxying

## Test Setup

### Prerequisites
1. Backend server running on `http://localhost:8002`
2. Webroot directory at `/opt/webroot` with test files:
   - `/opt/webroot/index.html`
   - `/opt/webroot/app.js`
   - `/opt/webroot/style.css`

### Start Command
```bash
./http-proxy -p 8001 -t http://localhost:8002 -w /opt/webroot -e backend
```

## Test Case 1: Static File Serving in Combined Mode

### Description
Verify that requests to non-backend paths serve static files from the webroot directory.

### Test Steps
```bash
# Given: Proxy is running with the command above

# When: Request a static file
curl http://localhost:8001/index.html

# Then: Verify response
# - Status: 200 OK
# - Content: Contents of /opt/webroot/index.html
# - NOT proxied to http://localhost:8002/index.html
```

### Expected Behavior
- ✅ Server reads `/opt/webroot/index.html`
- ✅ Returns file content with correct MIME type (`text/html`)
- ✅ Does NOT proxy to backend
- ✅ Backend server does NOT receive the request

### Routing Decision
```
Request: GET /index.html
  ├─ webroot specified? YES (-w /opt/webroot)
  ├─ endpoint specified? YES (-e backend)
  ├─ path starts with "backend"? NO
  └─ Action: Serve static file /opt/webroot/index.html
```

## Test Case 2: Backend Proxying in Combined Mode

### Description
Verify that requests to paths starting with the endpoint prefix are proxied to the backend.

### Test Steps
```bash
# Given: Proxy is running with the command above

# When: Request a backend path
curl http://localhost:8001/backend/index.html

# Then: Verify response
# - Request is proxied to http://localhost:8002/backend/index.html
# - Response comes from backend server
# - NOT served from /opt/webroot/backend/index.html (even if it exists)
```

### Expected Behavior
- ✅ Request is proxied to `http://localhost:8002/backend/index.html`
- ✅ Full path `/backend/index.html` is preserved
- ✅ Response comes from backend server
- ✅ Static file `/opt/webroot/backend/index.html` is NOT served (even if it exists)

### Routing Decision
```
Request: GET /backend/index.html
  ├─ webroot specified? YES (-w /opt/webroot)
  ├─ endpoint specified? YES (-e backend)
  ├─ path starts with "backend"? YES
  └─ Action: Proxy to http://localhost:8002/backend/index.html
```

## Additional Test Cases

### Test Case 3: Root Path
```bash
# Request root
curl http://localhost:8001/

# Expected: Serve /opt/webroot/index.html (if exists)
# Routing: "/" does not start with "backend" → serve static
```

### Test Case 4: Backend API Path
```bash
# Request backend API
curl http://localhost:8001/backend/api/users

# Expected: Proxy to http://localhost:8002/backend/api/users
# Routing: "/backend/api/users" starts with "backend" → proxy
```

### Test Case 5: Other Static Files
```bash
# Request CSS file
curl http://localhost:8001/style.css

# Expected: Serve /opt/webroot/style.css with Content-Type: text/css
# Routing: "/style.css" does not start with "backend" → serve static
```

### Test Case 6: Non-existent Static File
```bash
# Request non-existent file
curl http://localhost:8001/notfound.html

# Expected: 404 Not Found
# Routing: "/notfound.html" does not start with "backend" → serve static (404)
```

### Test Case 7: Case Sensitivity
```bash
# Request with different case
curl http://localhost:8001/Backend/api/users

# Expected: Serve static file (404 if not exists)
# Routing: "/Backend/api/users" does NOT start with "backend" (case-sensitive)
```

## Routing Logic Summary

### Decision Tree
```
Request arrives at http://localhost:8001{path}
    │
    ├─ Does path start with "/backend"?
    │   │
    │   YES ──> Proxy to http://localhost:8002{path}
    │   │       (preserve full path including "/backend")
    │   │
    │   NO ──> Serve static file from /opt/webroot{path}
    │          (or 404 if file not found)
```

### Examples Table

| Request Path | Starts with "backend"? | Action | Result |
|-------------|------------------------|--------|--------|
| `/index.html` | NO | Serve static | `/opt/webroot/index.html` |
| `/app.js` | NO | Serve static | `/opt/webroot/app.js` |
| `/backend/index.html` | YES | Proxy | `http://localhost:8002/backend/index.html` |
| `/backend/api/users` | YES | Proxy | `http://localhost:8002/backend/api/users` |
| `/backend` | YES | Proxy | `http://localhost:8002/backend` |
| `/Backend/api` | NO | Serve static | `/opt/webroot/Backend/api` (404 if not exists) |
| `/api/backend` | NO | Serve static | `/opt/webroot/api/backend` (404 if not exists) |

## Verification Commands

### Setup Test Environment
```bash
# 1. Create webroot directory
mkdir -p /opt/webroot

# 2. Create test files
echo "<html><body>Index Page</body></html>" > /opt/webroot/index.html
echo "console.log('app');" > /opt/webroot/app.js
echo "body { color: red; }" > /opt/webroot/style.css

# 3. Start a simple backend server on port 8002
# (use any web server or create a simple one)
cd /tmp && python3 -m http.server 8002 &

# 4. Start the proxy
./http-proxy -p 8001 -t http://localhost:8002 -w /opt/webroot -e backend
```

### Run Tests
```bash
# Test 1: Static file
curl -v http://localhost:8001/index.html
# Expected: 200, content from /opt/webroot/index.html

# Test 2: Backend proxy
curl -v http://localhost:8001/backend/index.html
# Expected: proxied to backend server

# Test 3: Root
curl -v http://localhost:8001/
# Expected: 200, content from /opt/webroot/index.html

# Test 4: Backend API
curl -v http://localhost:8001/backend/api/test
# Expected: proxied to backend server

# Test 5: CSS file
curl -v http://localhost:8001/style.css
# Expected: 200, Content-Type: text/css

# Test 6: 404
curl -v http://localhost:8001/notfound.html
# Expected: 404 Not Found
```

## Success Criteria

✅ **Test Case 1 passes**: Static files are served from webroot for non-backend paths  
✅ **Test Case 2 passes**: Backend paths are proxied with full path preservation  
✅ **Routing is deterministic**: Prefix match always proxies, non-match always serves static  
✅ **No ambiguity**: Even if `/opt/webroot/backend/index.html` exists, `/backend/index.html` is proxied  
✅ **Case-sensitive**: `/Backend` and `/backend` are treated differently  
✅ **Backward compatible**: Without `-w`, existing proxy behavior is unchanged

