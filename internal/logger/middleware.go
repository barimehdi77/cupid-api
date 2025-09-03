package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GinMiddleware returns a Gin middleware that logs HTTP requests using Zap
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get client IP
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()
		userAgent := c.Request.UserAgent()

		if raw != "" {
			path = path + "?" + raw
		}

		// Choose log level based on status code
		var logLevel zapcore.Level
		switch {
		case statusCode >= 500:
			logLevel = zap.ErrorLevel
		case statusCode >= 400:
			logLevel = zap.WarnLevel
		default:
			logLevel = zap.InfoLevel
		}

		// Create log fields
		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
			zap.String("user_agent", userAgent),
			zap.Int("body_size", bodySize),
		}

		// Add error field if there are any errors
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		// Log the request
		switch logLevel {
		case zap.ErrorLevel:
			Logger.Error("HTTP Request", fields...)
		case zap.WarnLevel:
			Logger.Warn("HTTP Request", fields...)
		default:
			Logger.Info("HTTP Request", fields...)
		}
	}
}

// GinRecoveryMiddleware returns a recovery middleware that logs panics using Zap
func GinRecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		Logger.Error("Panic recovered",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.Any("panic", recovered),
		)
		c.AbortWithStatus(500)
	})
}
