package main

import (
	"github.com/reiver/go-cors"

	"fmt"
	"net/http"
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

	var handler http.Handler
	{
		proxy := cors.ProxyHandler{
			LogWriter:os.Stdout,
		}

		handler = &proxy
	}


	{
		var addr string = fmt.Sprintf(":%d", tcpport)

		err := http.ListenAndServe(addr, handler)
		if nil != err {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			return
		}
	}
}