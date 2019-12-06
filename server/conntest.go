package main

import (
	"fmt"
	"loadbalancer/server/proxy"
	"net/http"
)

type contextKey struct {
	key string
}

func main() {

	proxy.SetupProxies()

	http.HandleFunc("/", defaultHandler)

	server := http.Server{
		Addr: ":8080",
	}
	fmt.Println("Listening")
	server.ListenAndServe()

}

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	endpoint, ok := proxy.GetLeastPopulated()

	if ok {
		fmt.Printf("Client sent to %s active: %d \n", endpoint.URLtext, endpoint.ActiveConnections())

		endpoint.CountCh <- 1
		endpoint.RedirectClient(w, r)

		endpoint.CountCh <- -1

	}

}
