package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/elazarl/goproxy"
)

const (
	tcpPort = 3001
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	// Handle HTTPS connect method
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	addr := fmt.Sprintf(":%d", tcpPort)
	fmt.Println("Starting Proxy on", addr)
	err := http.ListenAndServe(addr, proxy)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return
	}
}
