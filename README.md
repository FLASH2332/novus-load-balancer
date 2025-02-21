````reverse-proxy-lb/
│── cmd/
│   ├── main.go          # Entry point of the application
│
│── config/
│   ├── config.go        # Loads configuration from file/env
│   ├── config.yaml      # Configuration file (optional)
│
│── internal/
│   ├── proxy/
│   │   ├── proxy.go     # Reverse proxy logic
│   ├── loadbalancer/
│   │   ├── lb.go        # Load balancing logic
│   ├── cache/
│   │   ├── lru.go       # LRU cache implementation
│
│── pkg/
│   ├── utils.go         # Helper functions
│
│── go.mod               # Go module file
│── go.sum               # Dependencies checksum
│── README.md            # Project documentation
````
Planning to implement with the following structure
