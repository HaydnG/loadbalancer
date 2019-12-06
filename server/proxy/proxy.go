package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

//Endpoint is a proxy server endpoint definition
type Endpoint struct {
	URL               *url.URL
	URLtext           string
	proxy             *httputil.ReverseProxy
	activeConnections int
	CountCh           chan int
}

var proxies = map[int]*Endpoint{
	1: &Endpoint{
		URLtext:           "http://localhost:5000",
		activeConnections: 0,
	},
	2: &Endpoint{
		URLtext:           "http://localhost:5001",
		activeConnections: 0,
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
		go p.count()
	}
}

//count is a thread safe counter
func (ep *Endpoint) count() {
	for {
		select {
		case a := <-ep.CountCh:
			ep.activeConnections += a
		}
	}
}

//GetLeastPopulated returns the proxy with the least connections
func GetLeastPopulated() (*Endpoint, bool) {

	low := 1000000
	var lowp *Endpoint
	for _, p := range proxies {
		if p.activeConnections < low {

			low = p.activeConnections
			lowp = p
		}

	}
	if lowp == nil {
		return nil, false
	}

	return lowp, true
}

//ActiveConnections returns the number of clients currently on the proxy
func (ep *Endpoint) ActiveConnections() int {
	return ep.activeConnections
}

//SetupProxy initilises the proxy server
func (ep *Endpoint) SetupProxy() {
	ep.URL, _ = url.Parse(ep.URLtext)
	ep.proxy = httputil.NewSingleHostReverseProxy(ep.URL)
}

//RedirectClient sends the client to the given proxy
func (ep *Endpoint) RedirectClient(w http.ResponseWriter, r *http.Request) {
	ep.proxy.ServeHTTP(w, r)
}
