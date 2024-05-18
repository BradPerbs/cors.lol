package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const (
	requestLimit = 10 * 1024 * 1024 // 10 MB
	rateLimit    = 10               // requests per minute
)

// RateLimiter is a wrapper for rate limiting per IP
type RateLimiter struct {
	ips map[string]*rate.Limiter
	r   *rate.Limiter
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   rate.NewLimiter(r, b),
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	if lim, exists := rl.ips[ip]; exists {
		return lim
	}
	lim := rate.NewLimiter(rl.r.Limit(), rl.r.Burst())
	rl.ips[ip] = lim
	return lim
}

var limiter = NewRateLimiter(rate.Every(time.Minute/rateLimit), rateLimit)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	if ipLimiter := limiter.getLimiter(ip); !ipLimiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// Extract the target URL from the request path
	targetURL := r.URL.String()[1:] // Remove the leading slash

	// Ensure the target URL has the correct scheme (http or https)
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		http.Error(w, "Invalid URL scheme", http.StatusBadRequest)
		return
	}

	// Parse the target URL to ensure it is valid
	parsedURL, err := url.Parse(targetURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Create a new request based on the original request
	req, err := http.NewRequest(r.Method, parsedURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Error creating request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy original headers to the new request
	req.Header = r.Header

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error proxying request: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers to the original response
	for key, value := range resp.Header {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}

	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Write response status code and body
	w.WriteHeader(resp.StatusCode)
	if _, err := io.CopyN(w, resp.Body, requestLimit); err != nil && err != io.EOF {
		http.Error(w, "Error reading response body: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", proxyHandler)

	log.Println("Proxy server running on :3001")
	if err := http.ListenAndServe(":3001", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
