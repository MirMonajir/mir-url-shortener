package application

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// SecurityHeadersMiddleware adds security headers to all responses
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'")

		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Prevent MIME sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Enable browser XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// HSTS (HTTP Strict Transport Security) - only for HTTPS
		if c.Request.URL.Scheme == "https" {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		c.Next()
	}
}

// CORSMiddleware configures CORS with secure defaults
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only allow specified origins - in production, this should be configurable
		allowedOrigins := []string{"http://localhost:3000", "https://example.com"}

		origin := c.GetHeader("Origin")
		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				isAllowed = true
				break
			}
		}

		if isAllowed || origin == "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			c.Header("Access-Control-Max-Age", "3600")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// RecoveryMiddleware handles panics and recovers gracefully
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic but don't expose internals to client
				c.JSON(http.StatusInternalServerError, gin.H{
					"error_type": "internal_error",
					"message":    "An unexpected error occurred",
					"code":       http.StatusInternalServerError,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
