package main

import (
	"fmt"
	"loadbalancer/server/proxy"
	"net"
	"net/http"
)

type contextKey struct {
	key string
}

func main() {

	proxy.SetupProxies()

	http.HandleFunc("/", defaultHandler)

	server := http.Server{
		Addr:      ":8080",
		ConnState: connStateHook,
	}

	server.ListenAndServe()

}

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Connection recieved")
	endpoint := proxy.GetProxyByID(1)
	fmt.Printf("Connection redirected to: %s\n", endpoint.URLtext)
	endpoint.RedirectClient(w, r)
	fmt.Printf("Connection finished\n")

}

func connStateHook(connection net.Conn, state http.ConnState) {
	fmt.Printf("connState: %+v \n", state)
}
