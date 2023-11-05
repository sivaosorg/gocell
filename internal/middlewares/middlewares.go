package middlewares

import "github.com/sivaosorg/govm/configx"

type MiddlewareManager struct {
	conf *configx.KeysConfig
}

func NewMiddlewareManager(conf *configx.KeysConfig) *MiddlewareManager {
	return &MiddlewareManager{conf: conf}
}
