package router

import (
	"github.com/gin-gonic/gin"
	"github.com/poixeai/proxify/controller"
	"github.com/poixeai/proxify/middleware"
)

func SetRoutes(r *gin.Engine) {
	// basic middleware
	r.Use(middleware.Recover())
	r.Use(middleware.CORS())
	r.Use(middleware.GinRequestLogger())
	r.Use(middleware.Extractor())

	// ==== reserved routes ====
	r.GET("/", controller.HomeHandler)

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/", controller.ShowPathHandler)
		apiGroup.GET("/routes", controller.RoutesHandler)
		apiGroup.GET("/param", controller.ShowParamHandler)
	}

	// ==== routes.json ====
	r.NoRoute(controller.ProxyHandler)
}
