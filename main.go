package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

const (
	tcpport = 3001
)

func main() {
	fmt.Println("CORS Proxy")
	fmt.Println("")
	fmt.Println("In your /etc/hosts file add this line:")
	fmt.Println("")
	fmt.Println("127.0.0.1\tproxy.cors")
	fmt.Println("")
	fmt.Println("And then run a request against:")
	fmt.Println("")
	fmt.Printf("http://proxy.cors:%d/http://example.com/user/joeblow.atom\n", tcpport)
	fmt.Print("\n\n\n\n")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Trim the leading '/' from the request URL
		originalURL := r.URL.Path[1:]

		if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
			http.Error(w, "URL must start with http:// or https://", http.StatusBadRequest)
			return
		}

		parsedURL, err := url.Parse(originalURL)
		if err != nil {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(parsedURL)

		// Modify the request to enable CORS
		r.URL = parsedURL
		r.Host = parsedURL.Host

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		proxy.ServeHTTP(w, r)
	})

	addr := fmt.Sprintf(":%d", tcpport)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}
}
