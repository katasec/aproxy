package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {

	port, origin := checkEnv()

	// Setup proxy
	fmt.Println("Target URL: ", origin)
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {

			// Extract the target URL from the request
			target := fmt.Sprintf("%s://%s%s", origin.Scheme, origin.Host, r.In.URL.Host)
			targetUrl, _ := url.Parse(target)

			r.SetURL(targetUrl)

			// Log request
			httpLog(r, origin)
		},
	}

	// Handle all requests with the proxy
	http.Handle("/", proxy)
	addr := fmt.Sprintf("0.0.0.0:%s", port)

	// Start the server
	log.Printf("Server started on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
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
	for k, v := range r.Out.Header {
		fmt.Printf("Header: %s: %s\n", k, v)
	}

}
