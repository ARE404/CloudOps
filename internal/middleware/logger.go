package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LogMiddleware struct {
	logger *zap.Logger
}

func NewLogMiddleware(l *zap.Logger) *LogMiddleware {
	return &LogMiddleware{logger: l}
}

func (m *LogMiddleware) Log() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		m.logger.Info("request",
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
			zap.Int("status", ctx.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", ctx.ClientIP()),
		)
	}
}
