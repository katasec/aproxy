package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"tailscale.com/client/tailscale"
)

func main() {

	port, origin := checkEnv()

	// Setup proxy
	log.Println("Target URL: ", origin)
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			// Set the request host to the target host
			r.SetURL(origin)

			// Log request
			httpLog(r, origin)
		},
	}

	// Handle all requests with the proxy
	http.Handle("/", proxy)
	addr := fmt.Sprintf("0.0.0.0:%s", port)

	// Start the server

	localClient := &tailscale.LocalClient{}
	s := &http.Server{
		TLSConfig: &tls.Config{
			GetCertificate: localClient.GetCertificate,
		},
		Handler: proxy,
	}

	log.Printf("Server started on %s\n", addr)
	log.Fatal(s.ListenAndServeTLS("", ""))
	//log.Fatal(http.ListenAndServe(addr, nil))
}

func checkEnv() (string, *url.URL) {
	if os.Getenv("APROXY_TARGET_URL") == "" {
		fmt.Printf("Please ensure environment variable APROXY_TARGET_URL is set, exitting.\n")
		os.Exit(1)
	}

	port := os.Getenv("APROXY_TARGET_PORT")
	if port == "" {
		fmt.Printf("Please ensure environment variable APROXY_TARGET_PORT is set, exitting.\n")
		os.Exit(1)
	}

	origin, err := url.Parse(os.Getenv("APROXY_TARGET_URL"))
	if err != nil {
		panic(err)
	}
	return port, origin
}

func httpLog(r *httputil.ProxyRequest, origin *url.URL) {

	x := fmt.Sprintf("%s://%s%s", origin.Scheme, origin.Host, r.In.URL.String())
	log.Printf("Fetching %s\n", x)

	// Log request headers
	// for k, v := range r.Out.Header {
	// 	log.Printf("Header: %s: %s\n", k, v)
	// }

}
