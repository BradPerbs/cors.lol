package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	// Limit each IP to 20 requests per 5 minutes
	rateLimit         = 20
	rateLimitDuration = 5 * time.Minute
	requestCounts     = make(map[string]int)
	countsLock        = sync.Mutex{}
	// Max allowed size of the request body is 10MB
	maxBodySize int64 = 10 << 20 // 10 MB
	// HTTP client with timeouts
	client = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives:     true,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
		},
	}
	// URL patterns to rewrite
	urlPattern = regexp.MustCompile(`(src|href|url|srcset)=["']?([^"'\s>]+)["']?`)
	cssPattern = regexp.MustCompile(`url\(['"]?([^'")]+)['"]?\)`)
)

func init() {
	// Register additional MIME types
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".mjs", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".json", "application/json")
	mime.AddExtensionType(".woff", "font/woff")
	mime.AddExtensionType(".woff2", "font/woff2")
	mime.AddExtensionType(".ttf", "font/ttf")
	mime.AddExtensionType(".eot", "application/vnd.ms-fontobject")
	mime.AddExtensionType(".svg", "image/svg+xml")
	mime.AddExtensionType(".ico", "image/x-icon")
}

func main() {
	http.HandleFunc("/", limitRate(limitSize(handler)))

	log.Println("Starting server on :3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}

func resolveURL(baseURL *url.URL, ref string) string {
	// Skip data URLs and anchors
	if strings.HasPrefix(ref, "data:") || strings.HasPrefix(ref, "#") {
		return ref
	}

	// Handle protocol-relative URLs
	if strings.HasPrefix(ref, "//") {
		return baseURL.Scheme + ":" + ref
	}

	// Parse the reference URL
	refURL, err := url.Parse(ref)
	if err != nil {
		return ref
	}

	// If it's already absolute, return it
	if refURL.IsAbs() {
		return ref
	}

	// Resolve relative to base URL
	return baseURL.ResolveReference(refURL).String()
}

func rewriteURLs(content []byte, baseURL *url.URL, proxyURL string) []byte {
	// Rewrite HTML/CSS attribute URLs
	content = urlPattern.ReplaceAllFunc(content, func(match []byte) []byte {
		parts := urlPattern.FindSubmatch(match)
		if len(parts) != 3 {
			return match
		}

		attr := string(parts[1])
		urlStr := string(parts[2])

		// Skip data URLs and anchors
		if strings.HasPrefix(urlStr, "data:") || strings.HasPrefix(urlStr, "#") {
			return match
		}

		resolvedURL := resolveURL(baseURL, urlStr)
		return []byte(fmt.Sprintf(`%s="%s?url=%s"`, attr, proxyURL, url.QueryEscape(resolvedURL)))
	})

	// Rewrite CSS url() references
	content = cssPattern.ReplaceAllFunc(content, func(match []byte) []byte {
		parts := cssPattern.FindSubmatch(match)
		if len(parts) != 2 {
			return match
		}

		urlStr := string(parts[1])
		urlStr = strings.TrimSpace(urlStr)
		urlStr = strings.Trim(urlStr, `'"`)

		// Skip data URLs
		if strings.HasPrefix(urlStr, "data:") {
			return match
		}

		resolvedURL := resolveURL(baseURL, urlStr)
		return []byte(fmt.Sprintf(`url("%s?url=%s")`, proxyURL, url.QueryEscape(resolvedURL)))
	})

	return content
}

func detectContentType(filename string, content []byte) string {
	// First try by file extension
	ext := strings.ToLower(path.Ext(filename))
	if ext != "" {
		switch ext {
		case ".js", ".mjs":
			return "application/javascript"
		case ".css":
			return "text/css"
		case ".json":
			return "application/json"
		case ".svg":
			return "image/svg+xml"
		case ".woff":
			return "font/woff"
		case ".woff2":
			return "font/woff2"
		case ".ttf":
			return "font/ttf"
		case ".eot":
			return "application/vnd.ms-fontobject"
		}

		if mimeType := mime.TypeByExtension(ext); mimeType != "" {
			return mimeType
		}
	}

	// Then try by content sniffing
	if len(content) > 0 {
		return http.DetectContentType(content)
	}

	return "application/octet-stream"
}

func prepareURL(rawURL string) (string, error) {
	// Decode URL in case it's encoded
	decodedURL, err := url.QueryUnescape(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to decode URL: %v", err)
	}

	// Remove any whitespace
	decodedURL = strings.TrimSpace(decodedURL)

	// Add scheme if missing
	if !strings.HasPrefix(decodedURL, "http://") && !strings.HasPrefix(decodedURL, "https://") {
		decodedURL = "https://" + decodedURL
	}

	// Parse the URL to validate and normalize it
	parsedURL, err := url.Parse(decodedURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	// Ensure the URL has a host
	if parsedURL.Host == "" {
		return "", fmt.Errorf("invalid URL: missing host")
	}

	// Return the normalized URL
	return parsedURL.String(), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers first
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")

	// Handle OPTIONS method for preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get URL from query parameter
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		http.Error(w, "URL is required. Use format: /?url=https://example.com", http.StatusBadRequest)
		return
	}

	// Prepare and validate the URL
	preparedURL, err := prepareURL(targetURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Proxying request to: %s", preparedURL)

	// Parse the base URL for resolving relative paths
	baseURL, err := url.Parse(preparedURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid URL: %v", err), http.StatusBadRequest)
		return
	}

	// Create a new request
	req, err := http.NewRequestWithContext(r.Context(), r.Method, preparedURL, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create request: %v", err), http.StatusInternalServerError)
		return
	}

	// Copy headers from the original request
	for key, values := range r.Header {
		// Skip certain headers
		if strings.ToLower(key) == "host" || strings.ToLower(key) == "origin" {
			continue
		}
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Set some default headers if not present
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	}
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "*/*")
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching URL %s: %v", preparedURL, err)
		if strings.Contains(err.Error(), "timeout") {
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		} else if strings.Contains(err.Error(), "no such host") {
			http.Error(w, "Invalid host or DNS resolution failed", http.StatusBadGateway)
		} else {
			http.Error(w, fmt.Sprintf("Failed to fetch URL: %v", err), http.StatusBadGateway)
		}
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response: %v", err), http.StatusInternalServerError)
		return
	}

	// Detect content type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = detectContentType(preparedURL, body)
	}
	w.Header().Set("Content-Type", contentType)

	// Rewrite URLs in HTML and CSS content
	if strings.Contains(contentType, "text/html") || strings.Contains(contentType, "text/css") {
		proxyBaseURL := fmt.Sprintf("http://%s", r.Host)
		body = rewriteURLs(body, baseURL, proxyBaseURL)
	}

	// Copy other headers from the response
	for key, values := range resp.Header {
		if key != "Content-Type" && key != "Content-Length" { // Skip these as we handle them separately
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}

	// Ensure CORS headers are set after copying response headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")

	// Write the status code
	w.WriteHeader(resp.StatusCode)

	// Write the body
	if _, err := w.Write(body); err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}

	log.Printf("Successfully proxied %d bytes from %s with content type %s", len(body), preparedURL, contentType)
}

func limitRate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		countsLock.Lock()
		// Retrieve the current count
		count, exists := requestCounts[ip]

		if !exists {
			// Initialize the count for new IPs and set up a reset after the duration
			requestCounts[ip] = 1
			go func(ip string) {
				time.Sleep(rateLimitDuration)
				countsLock.Lock()
				delete(requestCounts, ip)
				countsLock.Unlock()
			}(ip)
		} else {
			// If IP exists and count is already at the limit, return error
			if count >= rateLimit {
				countsLock.Unlock()
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			// Otherwise, increment the count
			requestCounts[ip]++
		}
		countsLock.Unlock()

		next(w, r)
	}
}

func limitSize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
		next(w, r)
	}
}
