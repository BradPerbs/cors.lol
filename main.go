package main

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

// Define rate limiter with 1 request per second and a burst size of 5
var limiter = rate.NewLimiter(1, 5)

func main() {
	// Create a new mux router
	router := mux.NewRouter()

	// Define the proxy handler
	router.HandleFunc("/{url:.*}", func(w http.ResponseWriter, r *http.Request) {
		// Apply the rate limiter
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// Limit request size to 10MB
		r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)

		// Parse the target URL
		targetURL, err := url.Parse("https://example.com") // Change this to your target URL
		if err != nil {
			http.Error(w, "Invalid target URL", http.StatusInternalServerError)
			return
		}

		// Create a reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Set the CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		
		// Handle OPTIONS requests for CORS preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Serve the proxy request
		proxy.ServeHTTP(w, r)
	})

	// Wrap the router with CORS and logging handlers
	loggedRouter := handlers.LoggingHandler(io.Discard, router) // Use io.Discard to disable logging in production

	// Start the server
	http.ListenAndServe(":8080", loggedRouter)
}
