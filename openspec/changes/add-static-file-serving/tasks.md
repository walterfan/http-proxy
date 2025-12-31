# Implementation Tasks

## 1. Command-line Flag Setup
- [x] 1.1 Add `-w, --webroot` flag for static file directory path
- [x] 1.2 Update flag descriptions to clarify new behavior
- [x] 1.3 Add validation for webroot directory existence

## 2. Request Routing Logic
- [x] 2.1 Implement path prefix matching for endpoint flag
- [x] 2.2 Add routing decision logic (static vs proxy)
- [x] 2.3 Create static file handler using `http.FileServer`
- [x] 2.4 Update proxy handler to work with prefix matching

## 3. Handler Integration
- [x] 3.1 Refactor `handleProxy` to support dual mode (static + proxy)
- [x] 3.2 Add proper error handling for file not found cases
- [x] 3.3 Ensure correct MIME types for static files

## 4. Testing & Validation
- [x] 4.1 Test pure proxy mode (no webroot) - backward compatibility
- [x] 4.2 Test static file serving with various file types
- [x] 4.3 Test proxying with path prefix matching
- [x] 4.4 Test combined mode (static + proxy)
- [x] 4.5 Test edge cases (missing files, invalid paths)
- [x] 4.6 **Acceptance Test 1**: Start with `-p 8001 -t http://localhost:8002 -w /tmp/test-webroot -e /backend`, verify `http://localhost:8001/index.html` serves static file
- [x] 4.7 **Acceptance Test 2**: Start with `-p 8001 -t http://localhost:8002 -w /tmp/test-webroot -e /backend`, verify `http://localhost:8001/backend/index.html` proxies to backend

## 5. Documentation
- [x] 5.1 Update README.md with new usage examples
- [x] 5.2 Add examples for combined static + proxy mode
- [x] 5.3 Document backward compatibility guarantees

