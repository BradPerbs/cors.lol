package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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

	proxyHandler := func(w http.ResponseWriter, r *http.Request) {
		urlString := r.URL.Path[1:]
		parsedURL, err := url.Parse(urlString)
		if err != nil {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(parsedURL)
		r.URL.Path = parsedURL.Path
		r.URL.RawQuery = parsedURL.RawQuery

		// Modify the request to enable CORS
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = parsedURL.Host

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		proxy.ServeHTTP(w, r)
	}

	http.HandleFunc("/", proxyHandler)

	addr := fmt.Sprintf(":%d", tcpport)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
	}
}
