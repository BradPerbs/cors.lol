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

    // Handles HTTPS connections as well, Cloudflare will send these as HTTP to your server
    proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

    addr := fmt.Sprintf(":%d", tcpPort)
    err := http.ListenAndServe(addr, proxy)
    if err != nil {
        fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
        return
    }
}
