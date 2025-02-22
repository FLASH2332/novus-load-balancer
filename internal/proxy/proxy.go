package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"log"
)

// ReverseProxy struct
type ReverseProxy struct {
	TargetURLs []*url.URL
}

// NewReverseProxy initializes the proxy with target URLs
func NewReverseProxy(targets []string) *ReverseProxy {
	var urls []*url.URL
	for _, target := range targets {
		parsedURL, err := url.Parse(target)
		if err != nil {
			log.Fatalf("Invalid target URL: %v", err)
		}
		urls = append(urls, parsedURL)
	}
	return &ReverseProxy{TargetURLs: urls}
}

// ServeHTTP forwards requests to backend servers
func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target := rp.TargetURLs[0] // For now, always select the first target

	proxy := httputil.NewSingleHostReverseProxy(target)
	r.URL.Host = target.Host
	r.URL.Scheme = target.Scheme
	r.Host = target.Host

	log.Printf("Forwarding request to: %s%s", target, r.URL.Path)
	proxy.ServeHTTP(w, r)
}