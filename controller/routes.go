package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/poixeai/proxify/infra/watcher"
)

func RoutesHandler(c *gin.Context) {
	cfg := watcher.GetRoutes()
	c.JSON(http.StatusOK, gin.H{
		"data": cfg.Routes,
	})
}
