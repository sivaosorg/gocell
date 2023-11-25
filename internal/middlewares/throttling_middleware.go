package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	syncconf "github.com/sivaosorg/gocell/internal/syncConf"
	"github.com/sivaosorg/govm/entity"
	"github.com/sivaosorg/govm/logger"
	"github.com/sivaosorg/govm/ratelimitx"
	"golang.org/x/time/rate"
)

// https://blog.logrocket.com/rate-limiting-go-application/
type client struct {
	limiter *rate.Limiter
	last    time.Time
}

var (
	mu                   sync.Mutex
	clients              = make(map[string]*client)
	cleanupTypeWhitelist = time.Minute
	cleanupWhitelist     = 2 * cleanupTypeWhitelist
)

func (m *MiddlewareManager) RateLimitMiddleWare(key string) gin.HandlerFunc {
	rates := syncconf.Params.RateLimits
	clusters := ratelimitx.NewClusterMultiTenantRateLimitConfig().SetClusters(rates)
	go m.cleanupWhitelist()
	return func(c *gin.Context) {
		conf, err := clusters.FindClusterBy(key)
		if err != nil || !conf.Config.IsEnabled {
			return
		}
		ip, err := m.decodeNetwork(c)
		if err != nil {
			response := entity.NewResponseEntity().BadRequest(http.StatusText(http.StatusBadRequest), nil)
			response.SetError(err)
			c.JSON(response.StatusCode, response)
			c.Abort()
			return
		}
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(conf.Config.Rate), conf.Config.MaxBurst)}
		}
		clients[ip].last = time.Now()
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			response := entity.NewResponseEntity().TooManyRequest(http.StatusText(http.StatusTooManyRequests), nil)
			c.JSON(response.StatusCode, response)
			c.Abort()
			return
		}
		mu.Unlock()
		c.Next()
	}
}

func (m *MiddlewareManager) NetMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, err := m.decodeNetwork(c)
		if err != nil {
			m.notification(c, err, http.StatusBadRequest)
		}
		logger.Debugf(fmt.Sprintf("_endpoint: %v, incoming on IP: %v", c.Request.RequestURI, ip))
		c.Next()
	}
}

func (m *MiddlewareManager) decodeNetwork(c *gin.Context) (ip string, err error) {
	ip, _, err = net.SplitHostPort(c.Request.RemoteAddr)
	return ip, err
}

func (m *MiddlewareManager) cleanupWhitelist() {
	for {
		time.Sleep(cleanupTypeWhitelist)
		mu.Lock()
		for ip, client := range clients {
			remain := time.Since(client.last)
			if remain > cleanupWhitelist {
				logger.Debugf(fmt.Sprintf("Cleanup whitelist too many requests for IP: %v and remain duration: %v", ip, remain))
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}
