# http-proxy

a simple http proxy in golang

## Usage

```bash
# Basic usage: proxy all requests from port 7000 to localhost:7890
./http-proxy -p 7000 -t http://localhost:7890

# With endpoint filtering: only allow requests to /healthcheck
./http-proxy -p 80 -t http://localhost:8080 -e /healthcheck

# Or use the start script
./start.sh
```

## Command-line Options

- `-p, --listenPort`: Port to listen on (default: 7891)
- `-t, --targetUrl`: Target base URL to proxy to (default: http://localhost:7890)
- `-e, --endpoint`: Optional URL path endpoint to allow. If specified, only requests to this path are allowed. If not specified, all paths are forwarded.

## Examples

```bash
# Proxy all requests from port 80 to localhost:8080
./http-proxy -p 80 -t http://localhost:8080

# Only allow /metrics endpoint
./http-proxy -p 9090 -t http://localhost:8080 -e /metrics

# Only allow /healthcheck endpoint
./http-proxy -p 80 -t http://localhost:8080 -e /healthcheck
```

## License

Apache 2.0

## Author

Walter Fan

## Version

1.0.0