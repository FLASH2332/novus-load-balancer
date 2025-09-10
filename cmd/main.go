package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/FLASH2332/novus-load-balancer/config"
	"github.com/FLASH2332/novus-load-balancer/internal/proxy"
)

func main() {
	// Load configuration
	config.LoadConfig("config/config.yaml")

	// Choose strategy from config
	strategy := config.Cfg.LoadBalancer.Strategy

	// Create the reverse proxy
	reverseProxy := proxy.NewReverseProxy(config.Cfg.Proxy.Targets, strategy,config.Cfg.Cache.MaxSize)

	// Start the server
	port := config.Cfg.Server.Port
	fmt.Printf("Reverse Proxy running on port %d with %s strategy\n", port, strategy)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), reverseProxy))
}