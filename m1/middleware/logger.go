// m1/middleware/logger.go
package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs incoming requests and their duration
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Log after request is processed
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		log.Printf("[%s] %s %s - %d (%v)",
			method,
			path,
			c.ClientIP(),
			statusCode,
			duration,
		)
	}
}