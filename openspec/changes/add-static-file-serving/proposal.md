# Change: Add Static File Serving with Conditional Proxying

## Why
Users need a single HTTP server that can serve static web assets (HTML, CSS, JS) while also proxying specific API paths to a backend service. This eliminates the need to run separate web servers and simplifies deployment for full-stack applications.

## What Changes
- Add `-w, --webroot` flag to specify a directory for serving static files
- Add path prefix matching to `-e, --endpoint` flag (e.g., `/backend` matches `/backend/*`)
- When webroot is specified, serve static files by default
- When a request path matches the endpoint prefix, proxy to the target URL
- Maintain backward compatibility: if `-w` is not provided, behavior remains unchanged (pure proxy mode)

## Impact
- Affected specs: `http-proxy` (new capability)
- Affected code: `http_proxy.go` (main handler logic, flag parsing)
- Backward compatible: No breaking changes to existing functionality

