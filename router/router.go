package router

import (
	"github.com/gin-gonic/gin"
	"github.com/poixeai/proxify/controller"
	"github.com/poixeai/proxify/middleware"
)

func SetRoutes(r *gin.Engine) {
	// basic middleware
	r.Use(middleware.Recover())
	r.Use(middleware.GinRequestLogger())
	r.Use(middleware.Extractor())

	// ==== reserved routes ====
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/", controller.ShowPathHandler)
	}

	// ==== routes.json ====
	r.NoRoute(controller.ProxyHandler)
}
