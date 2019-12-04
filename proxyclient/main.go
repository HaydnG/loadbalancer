package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	http.HandleFunc("/", defaultHandler)

	server := http.Server{
		Addr: ":5000",
	}

	server.ListenAndServe()

}

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello you have hit a proxy\n"))
	fmt.Printf("%+v \n", r)
	time.Sleep(1000 * time.Millisecond)
}
