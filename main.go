package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/poixeai/proxify/infra/logger"
	"github.com/poixeai/proxify/infra/watcher"
	"github.com/poixeai/proxify/router"
	"github.com/poixeai/proxify/util"
)

func main() {
	// load .env
	_ = godotenv.Load()

	// init logger
	logger.InitLogger()

	// init routes watcher
	if err := watcher.InitRoutesWatcher(); err != nil {
		logger.Errorf("Failed to load routes config: %v", err)
		return
	}

	// init gin
	r := gin.New()
	r.SetTrustedProxies(nil)

	// setup routes
	router.SetRoutes(r)

	// start server
	port := util.GetEnvPort()
	if err := r.Run(":" + port); err != nil {
		logger.Errorf("Failed to start server: %v", err)
		return
	}

	logger.Infof("Server running on port %s", port)
}
