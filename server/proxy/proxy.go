package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//Endpoint is a proxy server endpoint definition
type Endpoint struct {
	URL               *url.URL
	URLtext           string
	proxy             *httputil.ReverseProxy
	ActiveConnections int
	CountCh           chan int
}

var proxies = map[int]*Endpoint{
	1: &Endpoint{
		URLtext:           "http://localhost:5000",
		ActiveConnections: 0,
	},
}

//GetProxies returns the list of proxies
func GetProxies() map[int]*Endpoint {
	return proxies
}

//GetProxyByID returns an endpoint by id
func GetProxyByID(id int) *Endpoint {
	return proxies[id]
}

//SetupProxies creates the proxy connections
func SetupProxies() {
	for _, p := range proxies {
		p.SetupProxy()
		//To accomodate for multiple client communications
		p.CountCh = make(chan int, 10)
	}
}

//SetupProxy initilises the proxy server
func (ep *Endpoint) SetupProxy() {
	ep.URL, _ = url.Parse(ep.URLtext)
	ep.proxy = httputil.NewSingleHostReverseProxy(ep.URL)
}

//RedirectClient sends the client to the given proxy
func (ep *Endpoint) RedirectClient(w http.ResponseWriter, r *http.Request) {
	ep.proxy.ServeHTTP(w, r)
	fmt.Printf("%+v", ep.proxy.ErrorLog)
}
