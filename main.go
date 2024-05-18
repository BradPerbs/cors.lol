package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
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

	url := strings.TrimPrefix(r.URL.Path, "/proxy/")
	if url == "" {
		http.Error(w, "Missing URL to proxy", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	req.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error proxying request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, value := range resp.Header {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	w.WriteHeader(resp.StatusCode)
	io.CopyN(w, resp.Body, requestLimit)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/proxy/{url:.*}", proxyHandler).Methods("GET", "POST", "PUT", "DELETE", "OPTIONS")

	log.Println("Proxy server running on :32000")
	if err := http.ListenAndServe(":32000", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
