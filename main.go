package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

//APROXY_TARGET_URL="https://go.dev"

func main() {
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

	fmt.Println("Target URL: ", origin)

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			target := fmt.Sprintf("%s://%s%s", origin.Scheme, origin.Host, r.In.URL.String())
			targetUrl, _ := url.Parse(target)
			r.SetURL(targetUrl)
			httpLog(r, origin)
		},
	}

	http.Handle("/", proxy)
	log.Println("Server started!")
	addr := fmt.Sprintf(":%s", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func httpLog(r *httputil.ProxyRequest, origin *url.URL) {
	x := fmt.Sprintf("%s://%s%s", origin.Scheme, origin.Host, r.In.URL.String())
	log.Printf("Fetching %s\n", x)
}
