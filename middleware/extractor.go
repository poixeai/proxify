package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/poixeai/proxify/infra/config"
	"github.com/poixeai/proxify/infra/ctx"
	"github.com/poixeai/proxify/infra/logger"
	"github.com/poixeai/proxify/infra/response"
	"github.com/poixeai/proxify/infra/watcher"
	"github.com/poixeai/proxify/util"
)

func Extractor() gin.HandlerFunc {
	return func(c *gin.Context) {
		top, sub := util.ExtractRoute(c.Request.URL.Path)
		// logger.Debugf("Extracted top route: %s, sub path: %s", top, sub)

		// store top and sub path into context for later use
		c.Set(ctx.TopRoute, top)
		c.Set(ctx.SubPath, sub)

		// system whitelist → allow directly
		if config.ReservedTopRoutes[top] {
			c.Next()
			return
		}

		// check if route exists in routes.json
		cfg := watcher.GetRoutes()
		found := false
		for _, r := range cfg.Routes {
			if r.Path == "/"+top {
				found = true
				c.Set(ctx.TargetEndpoint, r.Target)
				break
			}
		}

		if !found {
			logger.Error("Route not found: " + top)
			response.RespondTopRouteNotFoundError(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
