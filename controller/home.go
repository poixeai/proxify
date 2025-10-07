package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/poixeai/proxify/infra/response"
)

// HomeHandler returns a welcome message
func HomeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to Proxify!",
	})
}

// HealthCheckHandler checks the health of the service
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// ShowPathHandler returns the current request path
func ShowPathHandler(c *gin.Context) {
	path := c.Request.URL.Path
	c.JSON(http.StatusOK, gin.H{
		"path": path,
	})
}

// 404
func NotFoundHandler(c *gin.Context) {
	response.RespondSystemRouteNotFoundError(c)
}

// panic
func PanicHandler(c *gin.Context) {
	panic("Intentional panic for testing")
}
