package middlewares

import (
	"net/http"
	"time"

	"github.com/datpham2001/techblog/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger returns a gin middleware for logging HTTP requests
func Logger(l *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()
		latency := time.Since(start)

		entry := l.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"query":      raw,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"latency":    latency.String(),
			"latency_ms": latency.Milliseconds(),
		})

		// Add request ID if available
		if requestID := c.GetString("request_id"); requestID != "" {
			entry = entry.WithField("request_id", requestID)
		}

		// Add user ID if available (from auth middleware)
		if userID := c.GetString("user_id"); userID != "" {
			entry = entry.WithField("user_id", userID)
		}

		// Log based on status code
		if c.Writer.Status() >= http.StatusInternalServerError {
			entry.Error("HTTP Request")
		} else if c.Writer.Status() >= http.StatusBadRequest {
			entry.Warn("HTTP Request")
		} else {
			entry.Info("HTTP Request")
		}
	}
}
