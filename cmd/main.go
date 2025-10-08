package main

import (
	"log"

	"github.com/MirMonajir/mir-url-shortener/internal_logic/application"
	"github.com/MirMonajir/mir-url-shortener/internal_logic/infrastructure"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialising the dependencies
	inMemoryStore := infrastructure.NewInMemoryStore()
	shortenerService := application.NewShortenerService(inMemoryStore)

	// Setup gin router
	r := gin.Default()
	h := application.NewHTTPHandler(shortenerService)

	// Routes
	r.POST("/shortenurl", h.ShortenURL)
	r.GET("/:shortenedurl", h.Redirect)
	r.GET("/appmetrics", h.Metrics)

	// Start the application
	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start the URLShotener server: %v", err)
	}
}
