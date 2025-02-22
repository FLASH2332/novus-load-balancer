package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// Backend struct with connection count
type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	Connections  uint64 // Tracks active connections
	ReverseProxy *httputil.ReverseProxy
}

// ReverseProxy struct with strategy selection
type ReverseProxy struct {
	Backends   []*Backend
	currentIdx uint64
	Strategy   string // "round_robin" or "least_connections"
}

// NewReverseProxy initializes the proxy with target URLs and strategy
func NewReverseProxy(targets []string, strategy string) *ReverseProxy {
	var backends []*Backend
	for _, target := range targets {
		parsedURL, err := url.Parse(target)
		if err != nil {
			log.Fatalf("Invalid target URL: %v", err)
		}
		backends = append(backends, &Backend{
			URL:          parsedURL,
			Alive:        true,
			Connections:  0,
			ReverseProxy: httputil.NewSingleHostReverseProxy(parsedURL),
		})
	}
	rp := &ReverseProxy{Backends: backends, Strategy: strategy}
	go rp.healthCheckLoop() // Start passive health checks
	return rp
}

// getNextBackend selects backend based on strategy
func (rp *ReverseProxy) getNextBackend() *Backend {
	if rp.Strategy == "least_connections" {
		return rp.getLeastConnectionsBackend()
	}
	return rp.getRoundRobinBackend()
}

// Round Robin Selection
func (rp *ReverseProxy) getRoundRobinBackend() *Backend {
	for i := 0; i < len(rp.Backends); i++ {
		idx := atomic.AddUint64(&rp.currentIdx, 1) % uint64(len(rp.Backends))
		backend := rp.Backends[idx]
		backend.mux.RLock()
		alive := backend.Alive
		backend.mux.RUnlock()
		if alive {
			return backend
		}
	}
	return nil
}

// Least Connections Selection
func (rp *ReverseProxy) getLeastConnectionsBackend() *Backend {
	var selected *Backend
	minConnections := uint64(^uint(0)) // Max uint value

	for _, backend := range rp.Backends {
		backend.mux.RLock()
		alive := backend.Alive
		connections := backend.Connections
		backend.mux.RUnlock()

		if alive && connections < minConnections {
			selected = backend
			minConnections = connections
		}
	}
	return selected
}

// ServeHTTP with Active Connection Tracking
func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := rp.getNextBackend()
	if backend == nil {
		http.Error(w, "No backend available", http.StatusServiceUnavailable)
		return
	}

	// Increase connection count
	atomic.AddUint64(&backend.Connections, 1)
	defer atomic.AddUint64(&backend.Connections, ^uint64(0)) // Decrease after request

	log.Printf("Forwarding request to: %s%s (Connections: %d)", backend.URL, r.URL.Path, backend.Connections)
	backend.ReverseProxy.ServeHTTP(w, r)
}

// healthCheckLoop and isBackendAlive remain the same

func (rp *ReverseProxy) healthCheckLoop() {
	for {
		time.Sleep(5 * time.Second) // Check every 5 seconds
		log.Println("Running health check...") // Add this to confirm it's running

		for _, backend := range rp.Backends {
			alive := isBackendAlive(backend.URL)

			backend.mux.Lock()
			wasAlive := backend.Alive
			backend.Alive = alive
			backend.mux.Unlock()

			// Log when backend status changes
			if wasAlive && !alive {
				log.Printf("❌ Backend %s is DOWN!", backend.URL)
			} else if !wasAlive && alive {
				log.Printf("✅ Backend %s is BACK ONLINE!", backend.URL)
			}
		}
	}
}


func isBackendAlive(url *url.URL) bool {
	client := http.Client{
		Timeout: 2 * time.Second, // Prevents hanging if the server is dead
	}
	req, err := http.NewRequest("HEAD", url.String(), nil) // Use HEAD instead of GET
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Health check failed for %s: %v", url, err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode < 500
}