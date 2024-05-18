package main

import (
    "io"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", handler)
    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
    // Read the 'url' query parameter
    url := r.URL.Query().Get("url")
    if url == "" {
        http.Error(w, "URL is required", http.StatusBadRequest)
        return
    }

    // Set CORS headers
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

    // Handle OPTIONS method for preflight requests
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    // Proxy the request
    resp, err := http.Get(url)
    if err != nil {
        http.Error(w, "Failed to fetch URL", http.StatusInternalServerError)
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

