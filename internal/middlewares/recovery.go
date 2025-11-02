package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery catches panics and returns a clean JSON error response
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Internal Server Error",
				})
			}
		}()
		c.Next()
	}
}
