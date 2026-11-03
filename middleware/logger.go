package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Logger 请求日志中间件
// 记录: 方法、路径、状态码、耗时、客户端 IP
func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := strings.TrimSpace(c.GetHeader("X-Request-Id"))
		if requestID == "" || len(requestID) > 128 {
			requestID = uuid.NewString()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-Id", requestID)

		c.Next()

		status := c.Writer.Status()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}

		attrs := []any{
			slog.String("request_id", requestID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("route", route),
			slog.Int("status", status),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
			slog.String("client_ip", c.ClientIP()),
		}

		if userID, exists := c.Get("user_id"); exists {
			attrs = append(attrs, slog.Any("user_id", userID))
		}

		level := slog.LevelInfo

		switch {
		case status >= http.StatusInternalServerError:
			level = slog.LevelError
		case status == http.StatusTooManyRequests:
			level = slog.LevelWarn
		}

		logger.Log(
			c.Request.Context(),
			level,
			"http request",
			attrs...,
		)

	}
}
