package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"api-wiremock/internal/configuration"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	config  = &configuration.Config
	remote  url.URL
	service Service
)

type (
	proxyTransport struct {
		http.RoundTripper
		request Request
	}

	// Request describe the http request object
	Request struct {
		url    *url.URL
		body   []byte
		method string
	}
)

func main() {
	fmt.Println("Starting Wiremock Test Service")
	InitConfiguration()
	runServer()
}

func (t *proxyTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := t.RoundTripper.RoundTrip(request)

	err = service.CreateStubMapping(t.request, response)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Pass response on to the client
	return response, err
}

// InitConfiguration initialise package instance
func InitConfiguration() {
	configuration.InitConfig()
}

func runServer() {
	remote, err := url.Parse(config.APITarget)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)

	// Declare new Router
	r := mux.NewRouter()

	r.PathPrefix("/").HandlerFunc(handler(proxy))

	// Server default mux settings
	s := &http.Server{
		Addr:           ":8888",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	log.Printf("serving on port 8888")

	// Start Server
	log.Fatal(s.ListenAndServe())
}

func handler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Bad thing happen when I call this in roundTripper
		reqBuffer, _ := service.GetRequestBody(r)

		requestData := Request{
			url:    r.URL,
			method: r.Method,
			body:   reqBuffer,
		}

		var transport *http.Transport

		// Define http transport to use proxy
		if config.ProxyURL != "" {
			proxyURL, _ := url.Parse(config.ProxyURL)
			transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		} else {
			transport = &http.Transport{}
		}
		proxy.Transport = &proxyTransport{transport, requestData}

		r.Host = remote.Host

		proxy.ServeHTTP(w, r)
	}
}
