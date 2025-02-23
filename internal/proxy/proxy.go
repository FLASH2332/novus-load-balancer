package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"sync/atomic"

	"github.com/FLASH2332/novus-load-balancer/internal/loadbalancer"
)

// ReverseProxy struct
type ReverseProxy struct {
	LB *loadbalancer.LoadBalancer
}

// NewReverseProxy initializes the reverse proxy
func NewReverseProxy(targets []string, strategy string) *ReverseProxy {
	lb := loadbalancer.NewLoadBalancer(targets, strategy)
	return &ReverseProxy{LB: lb}
}

// ServeHTTP forwards requests to backend servers
func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := rp.LB.GetNextBackend()
	if backend == nil {
		http.Error(w, "No backend available", http.StatusServiceUnavailable)
		return
	}

	// Increase connection count
	atomic.AddUint64(&backend.Connections, 1)
	defer atomic.AddUint64(&backend.Connections, ^uint64(0)) // Decrease after request

	log.Printf("Forwarding request to: %s%s (Connections: %d)", backend.URL, r.URL.Path, backend.Connections)
	proxy := httputil.NewSingleHostReverseProxy(backend.URL)
	proxy.ServeHTTP(w, r)
}