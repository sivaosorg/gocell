package core

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (c *CoreCommand) routes(core *gin.Engine) {
	core.GET("/api/v1/swagger/index.html", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(-1),
	))
	v1 := core.Group("/api/v1")
	{
		v1.GET("/common/psql-status",
			c.handlers.middlewares.RequestMiddleWare(),
			c.handlers.middlewares.NoopMiddleWare(),
			c.handlers.middlewares.RateLimitMiddleWare("psql_rate"),
			c.handlers.middlewares.NetMiddleware(),
			c.handlers.commonHandler.OnPsqlStatus)
		v1.GET("/common/consumer", // endpoint websocket: ws://127.0.0.1:8081/api/v1/common/consumer
			c.handlers.middlewares.RequestMiddleWare(),
			c.handlers.commonHandler.OnSubscribe)
		v1.POST("/common/producer", // endpoint produce message to websocket
			c.handlers.middlewares.RequestMiddleWare(),
			c.handlers.commonHandler.OnProduce)
	}
}
