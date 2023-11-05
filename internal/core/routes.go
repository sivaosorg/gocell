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
}
