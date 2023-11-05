package core

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/sivaosorg/govm/configx"
	"github.com/sivaosorg/govm/dbx"
	"github.com/sivaosorg/govm/server"
	"github.com/sivaosorg/mysqlconn/mysqlconn"
	"github.com/sivaosorg/postgresconn/postgresconn"
	"github.com/sivaosorg/redisconn/redisconn"
)

// Conf global
// Get & Set key-value
// Based on conf, to create new cluster / instance
var Conf *configx.KeysConfig

type CoreCommand struct {
	psql        *postgresconn.Postgres
	psqlStatus  dbx.Dbx
	msql        *mysqlconn.MySql
	msqlStatus  dbx.Dbx
	redis       *redis.Client
	redisStatus dbx.Dbx
	handlers    *coreHandler
}

func (c *CoreCommand) Name() string {
	return "start_server"
}

func (c *CoreCommand) Description() string {
	return "Start Server HTTP/1.1"
}

func (c *CoreCommand) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("CoreCommand Args is required")
	}
	keys, err := configx.ReadConfig[configx.KeysConfig](args[0])
	if err != nil {
		return err
	}
	Conf = keys
	c.conn()
	c.handler()
	c.run()
	select {} // Keep application running
}

func (c *CoreCommand) conn() {
	Conf.Postgres.SetTimeout(10 * time.Second)
	Conf.MySql.SetTimeout(10 * time.Second)
	Conf.Redis.SetTimeout(10 * time.Second)

	// Instances
	go func() {
		psql, s := postgresconn.NewClient(Conf.Postgres)
		if s.IsConnected {
			defer psql.Close()
		}
		c.psql = psql
		c.psqlStatus = s
	}()
	go func() {
		msql, s := mysqlconn.NewClient(Conf.MySql)
		if s.IsConnected {
			defer msql.Close()
		}
		c.msql = msql
		c.msqlStatus = s
	}()
	go func() {
		redis, s := redisconn.NewClient(Conf.Redis)
		if s.IsConnected {
			defer redis.Close()
		}
		c.redis = redis
		c.redisStatus = s
	}()
}

func (c *CoreCommand) run() {
	gin.SetMode(Conf.Server.Mode)
	core := gin.New()
	core.SetTrustedProxies(nil)

	// base middlewares
	core.Use(gin.Logger())
	core.Use(c.handlers.middlewares.CorsMiddleware())
	core.Use(c.handlers.middlewares.Recovery())

	// set routes
	c.routes(core)

	// start server
	go func() {
		if Conf.Server.SSL.IsEnabled {
			err := core.RunTLS(Conf.Server.CreateAppServer(core.Handler()).Addr,
				Conf.Server.SSL.CertFile, Conf.Server.SSL.KeyFile)
			if err != nil {
				panic(err)
			}
		} else {
			err := core.Run(Conf.Server.CreateAppServer(core.Handler()).Addr)
			if err != nil {
				panic(err)
			}
		}
	}()
	go func() {
		if Conf.Server.SP.IsEnabled {
			debug := Conf.Server.SP.CreateAppServer(core.Handler())
			server.StartServer(debug)
		}
	}()
}
