package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"tailscale.com/client/tailscale"
	"tailscale.com/ipn/ipnstate"
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

	// Create a local client
	localClient := &tailscale.LocalClient{}

	// Start the server with a TailScale TLS certificate
	log.Println("Sleeping for 15 seconds to allow tailscale to start...")
	time.Sleep(15 * time.Second)
	s := &http.Server{
		TLSConfig: &tls.Config{
			GetCertificate: localClient.GetCertificate,
		},
		Handler: proxy,
	}
	log.Printf("Server started on %s\n", addr)
	log.Println("Please note that this server is only accessible via Tailscale VPN")
	log.Fatal(s.ListenAndServeTLS("", ""))

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

func waitForTailscale() {

	// Creatge a local client
	localClient := &tailscale.LocalClient{}

	// Get the local machine status.
	var status *ipnstate.Status
	var err error

	for {

		// Get the local machine status.
		log.Println("Get local machine status...")
		status, err = localClient.Status(context.Background())
		if err != nil {
			log.Printf("Failed to get local machine status: %v", err)
		}

		if err == nil {
			if !status.Self.Online {
				//sleep for 5 seconds
				log.Println("Local machine is not online, sleeping for 5 seconds...")
				time.Sleep(5 * time.Second)
				continue
			} else {
				log.Println("Online:", status.Self.Online)
				break
			}
		}
	}
}
