package proxy

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"sync/atomic"

	"github.com/FLASH2332/novus-load-balancer/internal/cache"
	"github.com/FLASH2332/novus-load-balancer/internal/loadbalancer"
)

// ReverseProxy struct
type ReverseProxy struct {
	LB    *loadbalancer.LoadBalancer
	Cache *cache.LRUCache
}

// NewReverseProxy initializes the reverse proxy with cache
func NewReverseProxy(targets []string, strategy string, cacheSize int) *ReverseProxy {
	lb := loadbalancer.NewLoadBalancer(targets, strategy)
	return &ReverseProxy{
		LB:    lb,
		Cache: cache.NewLRUCache(cacheSize),
	}
}

// ServeHTTP forwards requests to backend servers with caching
func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Use request URL as cache key
	cacheKey := r.URL.String()

	// Try cache first
	if cached, ok := rp.Cache.Get(cacheKey); ok {
		log.Printf("✅ Cache hit for %s", cacheKey)
		w.Write(cached)
		return
	}

	// If not cached, forward to backend
	backend := rp.LB.GetNextBackend()
	if backend == nil {
		http.Error(w, "No backend available", http.StatusServiceUnavailable)
		return
	}

	// Increase connection count
	atomic.AddUint64(&backend.Connections, 1)
	defer atomic.AddUint64(&backend.Connections, ^uint64(0)) // decrease after request

	log.Printf("➡ Forwarding request to: %s%s", backend.URL, r.URL.Path)

	proxy := httputil.NewSingleHostReverseProxy(backend.URL)

	// Capture response using custom transport
	proxy.ModifyResponse = func(resp *http.Response) error {
		// Read body into memory
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()

		// Save a copy to cache
		rp.Cache.Put(cacheKey, body)

		// Replace response body so client can still read it
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
		return nil
	}

	proxy.ServeHTTP(w, r)
}
