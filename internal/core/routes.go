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
			c.handlers.commonHandler.OnPsqlStatus)
	}
}
