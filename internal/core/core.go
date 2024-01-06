package core

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	dbresolver "github.com/sivaosorg/db.resolver"
	syncconf "github.com/sivaosorg/gocell/internal/syncConf"
	"github.com/sivaosorg/gocell/pkg/constant"
	"github.com/sivaosorg/govm/blueprint"
	"github.com/sivaosorg/govm/bot/telegram"
	"github.com/sivaosorg/govm/dbx"
	"github.com/sivaosorg/govm/logger"
	"github.com/sivaosorg/govm/server"
	"github.com/sivaosorg/govm/timex"
	"github.com/sivaosorg/govm/utils"
	"github.com/sivaosorg/msqlconn"
	"github.com/sivaosorg/postgresconn"
	"github.com/sivaosorg/redisconn"
	"github.com/sivaosorg/rmqconn"
)

type CoreCommand struct {
	psql           *postgresconn.Postgres
	psqlStatus     dbx.Dbx
	msql           *msqlconn.MySql
	msqlStatus     dbx.Dbx
	redis          *redisconn.Redis
	redisStatus    dbx.Dbx
	handlers       *coreHandler
	rabbitmq       *rmqconn.RabbitMq
	rabbitmqStatus dbx.Dbx
	resolver       *dbresolver.MultiTenantDBResolver
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
	err := c.snap(args)
	if err != nil {
		return err
	}
	c.seeker()
	c.conn()
	c.notify()
	c.handler()
	c.run()
	select {} // Keep application running
}

func (c *CoreCommand) conn() {
	if syncconf.Conf == nil {
		return
	}
	// Instances async
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		psql, s := postgresconn.NewClient(syncconf.Conf.Postgres)
		c.psql = psql
		c.psqlStatus = s
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		msql, s := msqlconn.NewClient(syncconf.Conf.MySql)
		c.msql = msql
		c.msqlStatus = s
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		redis, s := redisconn.NewClient(syncconf.Conf.Redis)
		c.redis = redis
		c.redisStatus = s
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		rabbitmq, s := rmqconn.NewClient(syncconf.Conf.RabbitMq)
		c.rabbitmq = rabbitmq
		c.rabbitmqStatus = s
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.resolver.
			AddPsqlConnectors(syncconf.Conf.PostgresSeekers...).
			AddMsqlConnectors(syncconf.Conf.MySqlSeekers...).
			SetDefaultConnector(dbresolver.NewPostgresConnector(syncconf.Conf.Postgres))
	}()
	wg.Wait()
}

func (c *CoreCommand) seeker() {
	c.resolver = dbresolver.NewMultiTenantDBResolver()
	timeout := 10 * time.Second
	syncconf.Conf.Postgres.SetTimeout(timeout)
	syncconf.Conf.MySql.SetTimeout(timeout)
	syncconf.Conf.Redis.SetTimeout(timeout)
	syncconf.Conf.RabbitMq.SetTimeout(timeout)

	// updating timeout for context db ping
	for idx := range syncconf.Conf.PostgresSeekers {
		syncconf.Conf.PostgresSeekers[idx].Config.SetTimeout(timeout)
	}
	for idx := range syncconf.Conf.MySqlSeekers {
		syncconf.Conf.MySqlSeekers[idx].Config.SetTimeout(timeout)
	}
}

func (c *CoreCommand) run() {
	gin.SetMode(syncconf.Conf.Server.Mode)
	core := gin.New()
	core.SetTrustedProxies(nil)

	// base middlewares
	core.Use(gin.Logger())
	core.Use(c.handlers.middlewares.CorsMiddleware())
	core.Use(c.handlers.middlewares.ErrorMiddleware)
	core.Use(c.handlers.middlewares.Recovery())

	// set routes
	c.routes(core)

	// start server
	go func() {
		if syncconf.Conf.Server.SSL.IsEnabled {
			err := core.RunTLS(syncconf.Conf.Server.CreateAppServer(core.Handler()).Addr,
				syncconf.Conf.Server.SSL.CertFile, syncconf.Conf.Server.SSL.KeyFile)
			if err != nil {
				panic(err)
			}
		} else {
			err := core.Run(syncconf.Conf.Server.CreateAppServer(core.Handler()).Addr)
			if err != nil {
				panic(err)
			}
		}
	}()
	go func() {
		if syncconf.Conf.Server.SP.IsEnabled {
			debug := syncconf.Conf.Server.SP.CreateAppServer(core.Handler())
			server.StartServer(debug)
		}
	}()
}

func (c *CoreCommand) snap(args []string) error {
	s := syncconf.NewSync()
	keys, err := s.GetClusters(args)
	if err != nil {
		return err
	}
	params, ok, err := s.GetParams(args)
	if ok {
		if err != nil {
			return err
		}
	}
	// sync and share config to variable global
	syncconf.Conf = keys
	syncconf.Params = params
	return nil
}

func (c *CoreCommand) notify() {
	conf, err := syncconf.Conf.FindTelegramSeeker(constant.TelegramKeyTenant1)
	if err != nil {
		logger.Errorf("Telegram Bot Notify", err)
		return
	}
	if !conf.Config.IsEnabled {
		return
	}
	go c.sendNotify("Psql Conn Alert", c.psqlStatus, conf)
	go c.sendNotify("Mysql Conn Alert", c.msqlStatus, conf)
	go c.sendNotify("Redis Conn Alert", c.redisStatus, conf)
	go c.sendNotify("RabbitMQ Conn Alert", c.rabbitmqStatus, conf)
}

func (c *CoreCommand) sendNotify(topic string, status dbx.Dbx, conf telegram.MultiTenantTelegramConfig) {
	var builder strings.Builder
	icon, _ := blueprint.TypeIcons[blueprint.TypeSuccess]
	if !status.IsConnected {
		icon, _ = blueprint.TypeIcons[blueprint.TypeError]
	}
	timestamp := timex.With(time.Now()).Format(timex.DateTimeFormYearMonthDayHourMinuteSecond)
	builder.WriteString(fmt.Sprintf("%v %s\n", icon, topic))
	builder.WriteString(fmt.Sprintf("tz: `%s`\n\n", timestamp))
	builder.WriteString(fmt.Sprintf("connected: `%v`\n", status.IsConnected))
	builder.WriteString(fmt.Sprintf("pid: `%v`\n", status.Pid))
	if utils.IsNotEmpty(status.Message) {
		builder.WriteString(fmt.Sprintf("message: `%v`\n", status.Message))
	}
	if status.Error != nil {
		builder.WriteString(fmt.Sprintf("error: `%v`\n", status.Error.Error()))
	}
	svc := telegram.NewTelegramService(conf.Config, conf.Option)
	svc.SendMessage(builder.String())
}
