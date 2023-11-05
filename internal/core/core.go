package core

import (
	"fmt"
	"time"

	"github.com/sivaosorg/govm/configx"
	"github.com/sivaosorg/mysqlconn/mysqlconn"
	"github.com/sivaosorg/postgresconn/postgresconn"
	"github.com/sivaosorg/redisconn/redisconn"
)

// Conf global
// Get & Set key-value
// Based on conf, to create new cluster / instance
var Conf *configx.KeysConfig

type CoreCommand struct {
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
	// Set Timeout deadline
	Conf.Postgres.SetTimeout(10 * time.Second)
	Conf.MySql.SetTimeout(10 * time.Second)
	Conf.Redis.SetTimeout(10 * time.Second)

	// Instances
	psql, s := postgresconn.NewClient(Conf.Postgres)
	if s.IsConnected {
		defer psql.Close()
	}
	msql, s := mysqlconn.NewClient(Conf.MySql)
	if s.IsConnected {
		defer msql.Close()
	}
	redis, s := redisconn.NewClient(Conf.Redis)
	if s.IsConnected {
		defer redis.Close()
	}
	return nil
}
