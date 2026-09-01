package main

import (
	"fmt"
	"log"

	"github.com/MirMonajir/mir-url-shortener/internal_logic/infrastructure"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration with validation
	cfg, err := infrastructure.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Setup router with middleware and routes
	r := SetupRouter()

	// Start server with configuration
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("Starting URL Shortener server on %s (mode: %s)", addr, cfg.GinMode)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start the URLShortener server: %v", err)
	}
}
