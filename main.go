package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
)

func main() {
	http.HandleFunc("/", limitRate(limitSize(handler)))

	log.Println("Starting server on :3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers first
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle OPTIONS method for preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get URL from either query parameter or path
	url := r.URL.Query().Get("url")
	if url == "" {
		// Try to extract from path if query parameter is empty
		path := r.URL.Path[1:] // Remove leading slash
		if strings.HasPrefix(path, "url=") {
			url = path[4:] // Remove "url=" prefix
		}
	}

	if url == "" {
		http.Error(w, "URL is required. Use format: /?url=https://example.com", http.StatusBadRequest)
		return
	}

	// Ensure URL has proper protocol
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	// Fix common URL formatting issues
	url = strings.Replace(url, "https:/", "https://", 1)
	url = strings.Replace(url, "http:/", "http://", 1)

	// Proxy the request
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch URL: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Write the status code and response body
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
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
