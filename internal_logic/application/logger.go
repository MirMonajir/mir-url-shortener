package application

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// RequestLoggerMiddleware logs HTTP requests with structured output
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		// Process request
		c.Next()

		// Log after response
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		contentLength := c.Writer.Size()

		log.Printf("[%s] %s %s | status: %d | latency: %v | size: %d bytes | ip: %s",
			time.Now().Format("2006-01-02 15:04:05"),
			method,
			path,
			statusCode,
			duration,
			contentLength,
			clientIP,
		)
	}
}

// ErrorLoggerMiddleware logs errors with context
func ErrorLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Log errors that occurred
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Printf("[ERROR] [%s] %s %s | error: %v | type: %v",
					time.Now().Format("2006-01-02 15:04:05"),
					c.Request.Method,
					c.Request.URL.Path,
					err.Error(),
					err.Type,
				)
			}
		}
	}
}
