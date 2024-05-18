package main

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/juju/ratelimit"
)

const (
	maxRequestSize = 10 << 20 // 10 MB
	rateLimit      = 5        // requests per minute
)

var bucket *ratelimit.Bucket

func init() {
	// Initialize the rate limiter
	bucket = ratelimit.NewBucket(time.Minute/time.Duration(rateLimit), int64(rateLimit))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{url:.*}", proxyHandler)
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Rate limiting
	if bucket.TakeAvailable(1) == 0 {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// Check request size
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
	defer r.Body.Close()

	// Get the target URL from the request
	vars := mux.Vars(r)
	targetURL := vars["url"]
	if !strings.HasPrefix(targetURL, "http") {
		targetURL = "http://" + targetURL
	}

	// Create the request to the target URL
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Copy the headers
	for k, v := range r.Header {
		req.Header[k] = v
	}

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy the response headers and status code
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)

	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Copy the response body
	io.Copy(w, resp.Body)
}

