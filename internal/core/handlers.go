package core

import "github.com/sivaosorg/gocell/internal/middlewares"

type coreHandler struct {
	middlewares *middlewares.MiddlewareManager
}

func NewCoreHandler() *coreHandler {
	return &coreHandler{}
}

func (c *coreHandler) setMiddlewares(value *middlewares.MiddlewareManager) *coreHandler {
	c.middlewares = value
	return c
}

func (c *CoreCommand) handler() {
	c.handlers = NewCoreHandler().
		setMiddlewares(middlewares.NewMiddlewareManager(Conf))
}
