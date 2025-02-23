package loadbalancer

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// Backend struct to store server state
type Backend struct {
	URL         *url.URL
	Alive       bool
	Connections uint64
	mux         sync.RWMutex
}

// LoadBalancer struct
type LoadBalancer struct {
	Backends   []*Backend
	currentIdx uint64
	Strategy   string // "round_robin" or "least_connections"
}

// NewLoadBalancer initializes the load balancer
func NewLoadBalancer(targets []string, strategy string) *LoadBalancer {
	var backends []*Backend
	for _, target := range targets {
		parsedURL, err := url.Parse(target)
		if err != nil {
			log.Fatalf("Invalid target URL: %v", err)
		}
		backends = append(backends, &Backend{
			URL:         parsedURL,
			Alive:       true,
			Connections: 0,
		})
	}
	lb := &LoadBalancer{Backends: backends, Strategy: strategy}
	go lb.HealthCheckLoop() // Start health checks in the background
	return lb
}

// GetNextBackend selects a backend based on strategy
func (lb *LoadBalancer) GetNextBackend() *Backend {
	if lb.Strategy == "least_connections" {
		return lb.getLeastConnectionsBackend()
	}
	return lb.getRoundRobinBackend()
}

// Round Robin Selection
func (lb *LoadBalancer) getRoundRobinBackend() *Backend {
	for i := 0; i < len(lb.Backends); i++ {
		idx := atomic.AddUint64(&lb.currentIdx, 1) % uint64(len(lb.Backends))
		if lb.Backends[idx].Alive {
			return lb.Backends[idx]
		}
	}
	return nil
}

// Least Connections Selection
func (lb *LoadBalancer) getLeastConnectionsBackend() *Backend {
	var selected *Backend
	minConnections := uint64(^uint(0)) // Max uint value

	for _, backend := range lb.Backends {
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

// HealthCheckLoop periodically pings backends to check availability
func (lb *LoadBalancer) HealthCheckLoop() {
	for {
		time.Sleep(5 * time.Second) // Check every 5 seconds
		log.Println("Running health check...")

		for _, backend := range lb.Backends {
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

// isBackendAlive checks if a backend is alive
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