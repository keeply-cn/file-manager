package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Config struct {
	root     string
	password string
	port     int
	basePath string
}

func main() {
	cfg := &Config{}
	flag.StringVar(&cfg.root, "root", "", "Root directory to serve (required)")
	flag.StringVar(&cfg.password, "password", "", "Password for authentication (required)")
	flag.IntVar(&cfg.port, "port", 8080, "Port to listen on")
	flag.StringVar(&cfg.basePath, "basepath", "/", "Base path for routing")
	flag.Parse()

	if cfg.root == "" || cfg.password == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := os.MkdirAll(cfg.root, 0755); err != nil {
		log.Fatalf("Failed to create root directory: %v", err)
	}

	log.Printf("Starting server on :%d with base path %s", cfg.port, cfg.basePath)
	log.Printf("Serving files from: %s", cfg.root)

	if err := run(cfg); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func run(cfg *Config) error {
	srv := newServer(cfg)
	return srv.Start()
}
