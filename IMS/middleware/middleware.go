package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.Printf("HTTP Request: method=%s path=%s status=%d",
			param.Method,
			param.Path,
			param.StatusCode,
		)
		return ""
	})
}

func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		if latency > 5*time.Second {
			log.Printf("Slow request detected: path=%s method=%s latency=%v",
				c.Request.URL.Path,
				c.Request.Method,
				latency,
			)
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Next()
	}
}
