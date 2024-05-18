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
    fmt.Println("CORS Proxy")
    fmt.Println("")
    fmt.Println("In your /etc/hosts file add this line:")
    fmt.Println("")
    fmt.Println("127.0.0.1\tproxy.cors")
    fmt.Println("")
    fmt.Println("And then run a request against:")
    fmt.Println("")
    fmt.Printf("http://proxy.cors:%d/http://example.com/user/joeblow.atom\n", tcpPort)
    fmt.Print("\n\n\n\n")

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
