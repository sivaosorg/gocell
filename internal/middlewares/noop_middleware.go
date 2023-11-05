package middlewares

import "github.com/gin-gonic/gin"

func (m *MiddlewareManager) NoopMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
