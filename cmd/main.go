package main

import (
	"fmt"
	"log"

	"github.com/FLASH2332/novus-load-balancer/config"
	"github.com/FLASH2332/novus-load-balancer/internal/proxy"
	"net/http"
)

func main() {
	// Load configuration
	config.LoadConfig("config/config.yaml")

	// Create a new reverse proxy
	reverseProxy := proxy.NewReverseProxy(config.Cfg.Proxy.Targets)

	// Start the server
	port := config.Cfg.Server.Port
	fmt.Printf("Reverse Proxy running on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), reverseProxy))

}