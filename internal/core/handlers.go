package core

import (
	"github.com/sivaosorg/gocell/internal/handlers"
	"github.com/sivaosorg/gocell/internal/middlewares"
	"github.com/sivaosorg/gocell/internal/repository"
	"github.com/sivaosorg/gocell/internal/service"
	syncconf "github.com/sivaosorg/gocell/internal/syncConf"
)

type coreHandler struct {
	middlewares   *middlewares.MiddlewareManager
	commonHandler *handlers.CommonHandler
}

func (c *CoreCommand) handler() {
	commonRepository := repository.NewCommonRepository(c.resolver)
	commonSvc := service.NewCommonService(commonRepository)
	commonHandler := handlers.NewCommonHandler(commonSvc)

	c.handlers = NewCoreHandler().
		setMiddlewares(middlewares.NewMiddlewareManager(syncconf.Conf)).
		setCommonHandler(commonHandler)
}
