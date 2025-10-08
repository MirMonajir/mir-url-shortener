package main

import (
	"github.com/MirMonajir/mir-url-shortener/internal_logic/application"
	"github.com/MirMonajir/mir-url-shortener/internal_logic/infrastructure"
	"github.com/gin-gonic/gin"
)

// setupRouter initializes the server routes and dependencies but DOES NOT start it.
func setupRouter() *gin.Engine {
	inMemoryStore := infrastructure.NewInMemoryStore()
	shortenerService := application.NewShortenerService(inMemoryStore)

	r := gin.Default()
	h := application.NewHTTPHandler(shortenerService)

	r.POST("/shortenurl", h.ShortenURL)
	r.GET("/:shortenedurl", h.Redirect)
	r.GET("/appmetrics", h.Metrics)

	return r
}
