package main

import (
    "io"
    "log"
    "net/http"
    "sync"
    "time"
)

var (
    // Limit each IP to 100 requests per 30 minutes
    rateLimit = 10
    rateLimitDuration = 30 * time.Minute
    requestCounts = make(map[string]int)
    lastAccess = make(map[string]time.Time)
    countsLock = sync.Mutex{}
    // Max allowed size of the request body is 10MB
    maxBodySize int64 = 10 << 20 // 10 MB
)

func main() {
    http.HandleFunc("/", limitRate(limitSize(handler)))

    log.Println("Starting server on :3001")
    log.Fatal(http.ListenAndServe(":3001", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
    // handler implementation remains the same
}

func limitRate(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        countsLock.Lock()
        defer countsLock.Unlock()

        last, found := lastAccess[ip]
        if found && time.Since(last) > rateLimitDuration {
            // Reset count and last access time after the rate limit duration has passed
            delete(requestCounts, ip)
            delete(lastAccess, ip)
        }

        count, exists := requestCounts[ip]
        if exists {
            if count >= rateLimit {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            requestCounts[ip]++
        } else {
            requestCounts[ip] = 1
            lastAccess[ip] = time.Now()
        }

        next(w, r)
    }
}

func limitSize(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
        next(w, r)
    }
}
