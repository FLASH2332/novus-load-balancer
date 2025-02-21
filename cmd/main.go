package main

import (
	"fmt"
	"log"

	"github.com/FLASH2332/novus-load-balancer/config"
)

func main() {
	// Load configuration
	config.LoadConfig("config/config.yaml")

	// Print loaded config values
	fmt.Println("Server Port:", config.Cfg.Server.Port)
	fmt.Println("Proxy Targets:", config.Cfg.Proxy.Targets)
	fmt.Println("Cache Max Size:", config.Cfg.Cache.MaxSize)
	fmt.Println("Load Balancer Strategy:", config.Cfg.LoadBalancer.Strategy)

	// Keep the application running
	log.Println("Configuration loaded successfully!")
}
