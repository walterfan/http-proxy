# http-proxy

a simple http proxy in golang

## Usage

```bash
# Basic usage: proxy all requests from port 7000 to localhost:7890
./http-proxy -p 7000 -t http://localhost:7890

# Serve static files only
./http-proxy -p 8080 -w /opt/webroot

# Combined mode: serve static files + proxy backend API
./http-proxy -p 8001 -t http://localhost:8002 -w /opt/webroot -e /backend

# Or use the start script
./start.sh
```

## Command-line Options

- `-p, --listenPort`: Port to listen on (default: 7891)
- `-t, --targetUrl`: Target base URL to proxy to (default: http://localhost:7890)
- `-e, --endpoint`: Optional URL path endpoint to allow. When used with `-w`, requests starting with this path are proxied. Without `-w`, only exact matches are allowed.
- `-w, --webroot`: Optional directory path for serving static files. When specified, serves static files by default.

## Examples

### Pure Proxy Mode

```bash
# Proxy all requests from port 80 to localhost:8080
./http-proxy -p 80 -t http://localhost:8080

# Only allow /metrics endpoint (exact match)
./http-proxy -p 9090 -t http://localhost:8080 -e /metrics

# Only allow /healthcheck endpoint (exact match)
./http-proxy -p 80 -t http://localhost:8080 -e /healthcheck
```

### Static File Serving Mode

```bash
# Serve static files from /opt/webroot on port 8080
./http-proxy -p 8080 -w /opt/webroot
```

### Combined Mode (Static Files + Backend Proxy)

```bash
# Serve static files, proxy /backend/* to backend server
./http-proxy -p 8001 -t http://localhost:8002 -w /opt/webroot -e /backend

# Result:
# - http://localhost:8001/index.html        → serves /opt/webroot/index.html
# - http://localhost:8001/app.js            → serves /opt/webroot/app.js
# - http://localhost:8001/backend/api/users → proxies to http://localhost:8002/backend/api/users
# - http://localhost:8001/backend/health    → proxies to http://localhost:8002/backend/health
```

## Routing Logic

When both `-w` (webroot) and `-e` (endpoint) are specified:
1. If request path starts with the endpoint prefix → proxy to backend
2. Otherwise → serve static file from webroot

This allows a single server to handle both frontend assets and backend API calls.

## License

Apache 2.0

## Author

Walter Fan

## Version

1.0.0