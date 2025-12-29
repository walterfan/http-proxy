package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

var targetBase *url.URL
var allowedEndpoint string

func main() {
	// Define command-line flags
	listenPort := flag.Int("listenPort", 7891, "Port to listen on")
	flag.IntVar(listenPort, "p", 7891, "Port to listen on (short)")

	targetUrl := flag.String("targetUrl", "http://localhost:7890", "Target base URL to proxy to (e.g. http://localhost:8081)")
	flag.StringVar(targetUrl, "t", "http://localhost:7890", "Target base URL to proxy to (short)")

	endpoint := flag.String("endpoint", "", "URL path endpoint to allow (e.g. /api/v1). If specified, only requests to this path are allowed. If not specified, all paths are allowed")
	flag.StringVar(endpoint, "e", "", "URL path endpoint to allow (short)")

	flag.Parse()

	// Store the allowed endpoint globally
	allowedEndpoint = *endpoint

	// Validate target URL
	if *targetUrl == "" {
		log.Fatal("Missing required parameter: --targetUrl")
	}
	var err error
	targetBase, err = url.Parse(*targetUrl)
	if err != nil {
		log.Fatalf("Invalid targetUrl: %v", err)
	}

	// Start HTTP server
	http.HandleFunc("/", handleProxy)
	addr := fmt.Sprintf(":%d", *listenPort)
	log.Printf("Proxy server listening on %s, forwarding to %s", addr, targetBase)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	// Check if endpoint filtering is enabled and validate the path
	if allowedEndpoint != "" && r.URL.Path != allowedEndpoint {
		errMsg := fmt.Sprintf("Forbidden: Access to this %s is not allowed, only %s is allowed", r.URL.Path, allowedEndpoint)
		http.Error(w, errMsg, http.StatusForbidden)
		return
	}

	// Construct full target URL
	targetURL := targetBase.ResolveReference(r.URL)

	// Create a new request to the target
	req, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header = r.Header.Clone()

	// Send request to target
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Request to target failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy headers and status code
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)

	// Copy body
	io.Copy(w, resp.Body)
}
