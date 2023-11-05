package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (m *MiddlewareManager) CorsMiddleware() gin.HandlerFunc {
	if !m.conf.Cors.IsEnabled {
		return m.NoopMiddleWare()
	}
	return cors.New(cors.Config{
		AllowOrigins:     m.conf.Cors.AllowedOrigins,
		AllowMethods:     m.conf.Cors.AllowedMethods,
		AllowHeaders:     m.conf.Cors.AllowedHeaders,
		ExposeHeaders:    m.conf.Cors.ExposedHeaders,
		AllowCredentials: m.conf.Cors.AllowCredentials,
		MaxAge:           time.Duration(m.conf.Cors.MaxAge),
	})
}
