package middleware

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/poixeai/proxify/infra/logger"
	"github.com/poixeai/proxify/infra/response"
)

// catch panic and return JSON error
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("panic: %v\n%s", err, debug.Stack())
				response.RespondInternalError(c)
			}
		}()

		c.Next()
	}
}
