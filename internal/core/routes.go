package core

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (c *CoreCommand) routes(e *gin.Engine) {
	c.shared(e)
	c.protected(e)
}

// Collection of authenticated endpoints
func (c *CoreCommand) protected(e *gin.Engine) {
	v1 := e.Group("/api/v1")
	v1.GET("/swagger/index.html", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	c.handlers.commonHandler.Router(v1.Group("/common"), c.handlers.middlewares)
}

// Collection of shared/public endpoints
func (c *CoreCommand) shared(e *gin.Engine) {
	v1 := e.Group("/api/v1/shared")

	c.handlers.commonHandler.Router(
		v1.Group("/common",
			c.handlers.middlewares.RequestMiddleWare(),
			c.handlers.middlewares.NetMiddleware(),
			c.handlers.middlewares.RateLimitMiddleWare("psql_rate"),
		),
		c.handlers.middlewares)
}
