package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
)

func main() {
	http.HandleFunc("/", limitRate(limitSize(handler)))

	log.Println("Starting server on :3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
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
	// Set security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "no-referrer")

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle OPTIONS method for preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get URL from query parameter
	targetURL := r.URL.Query().Get("url")
	log.Printf("Raw URL from query: %s", targetURL)

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

	log.Printf("Prepared URL to fetch: %s", preparedURL)

	// Create a new request
	req, err := http.NewRequestWithContext(r.Context(), "GET", preparedURL, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create request: %v", err), http.StatusInternalServerError)
		return
	}

	// Set common headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")

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

	log.Printf("Successfully fetched URL %s with status code %d", preparedURL, resp.StatusCode)

	// Copy important headers from the response
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	// Copy other relevant headers
	for _, header := range []string{"Cache-Control", "Expires", "Last-Modified", "ETag"} {
		if value := resp.Header.Get(header); value != "" {
			w.Header().Set(header, value)
		}
	}

	// Write the status code
	w.WriteHeader(resp.StatusCode)

	// Copy the response body
	written, err := io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error copying response body after %d bytes: %v", written, err)
		return
	}

	log.Printf("Successfully copied %d bytes from %s", written, preparedURL)
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
