package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger logs basic request information (method, path, duration)
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start)

		log.Printf("[%d] %s %s (%v)", status, method, path, duration)
	}
}
