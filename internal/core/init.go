package core

import (
	"github.com/sivaosorg/gocell/internal/handlers"
	"github.com/sivaosorg/gocell/internal/middlewares"
)

func NewCoreHandler() *coreHandler {
	return &coreHandler{}
}

func (c *coreHandler) setMiddlewares(value *middlewares.MiddlewareManager) *coreHandler {
	c.middlewares = value
	return c
}

func (c *coreHandler) setCommonHandler(value *handlers.CommonHandler) *coreHandler {
	c.commonHandler = value
	return c
}
