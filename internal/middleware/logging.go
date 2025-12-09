package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Logger creates a structured logging middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Log format
		log.Printf("[%s] %d | %s | %s | %s %s | %v",
			requestID[:8],
			statusCode,
			clientIP,
			c.Request.Method,
			path,
			query,
			latency,
		)

		// Log errors if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Printf("[%s] ERROR: %s", requestID[:8], err.Error())
			}
		}
	}
}

// GetRequestID extracts request ID from context
func GetRequestID(c *gin.Context) string {
	requestID, exists := c.Get("request_id")
	if !exists {
		return ""
	}
	id, ok := requestID.(string)
	if !ok {
		return ""
	}
	return id
}
