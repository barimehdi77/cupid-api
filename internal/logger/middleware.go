package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinMiddleware returns a Gin middleware that logs HTTP requests using enhanced Zap logging
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request details
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()
		userAgent := c.Request.UserAgent()

		if raw != "" {
			path = path + "?" + raw
		}

		// Prepare additional fields
		fields := []zap.Field{
			zap.String("ip", clientIP),
			zap.Int("size", bodySize),
		}

		// Add user agent for non-health checks
		if path != "/health" && path != "/ping" {
			fields = append(fields, zap.String("user_agent", userAgent))
		}

		// Add error field if there are any errors
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		// Use the enhanced LogRequest function
		LogRequest(method, path, statusCode, latency, fields...)
	}
}

// GinRecoveryMiddleware returns a recovery middleware that logs panics using enhanced Zap logging
func GinRecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		Logger.Error("💥 Panic recovered - server error",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Any("panic_value", recovered),
			zap.String("recovery_action", "returning 500 status"),
		)
		c.AbortWithStatus(500)
	})
}
