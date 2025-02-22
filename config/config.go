package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config structure
type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Proxy struct {
		Targets []string `yaml:"targets"`
	} `yaml:"proxy"`
	Cache struct {
		MaxSize int `yaml:"max_size"`
	} `yaml:"cache"`
	LoadBalancer struct {
		Strategy string `yaml:"strategy"`
	} `yaml:"loadbalancer"`
}

// Cfg is the global config instance
var Cfg Config

// LoadConfig reads YAML config from file
func LoadConfig(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	if err := yaml.Unmarshal(file, &Cfg); err != nil {
		log.Fatalf("Error parsing YAML: %v", err)
	}
}