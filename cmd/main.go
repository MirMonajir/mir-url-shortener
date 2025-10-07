package main

import (
    "github.com/gin-gonic/gin"
    "github.com/MirMonajir/mir-url-shortener/internal_logic/application"
    "github.com/MirMonajir/mir-url-shortener/internal_logic/infrastructure"
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
    r.Run(":8080")
}
