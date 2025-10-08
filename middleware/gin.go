package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/poixeai/proxify/infra/config"
	"github.com/poixeai/proxify/infra/ctx"
	"github.com/poixeai/proxify/infra/logger"
	"github.com/poixeai/proxify/util"
)

func GinRequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		reqID := util.GenerateRequestID()
		c.Set(ctx.RequestID, reqID)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.GetString(ctx.SubPath)
		clientIP := c.ClientIP()
		targetURL := c.GetString(ctx.TargetURL)

		topRoute := c.GetString(ctx.TopRoute)
		if config.ReservedTopRoutes[topRoute] {
			targetURL = "-"
		}

		logger.Infof(
			"%s | %d | %s | %s -> %s | %v | %s",
			reqID, status, method, path, targetURL, latency, clientIP,
		)
	}
}
