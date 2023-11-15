package middlewares

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sivaosorg/gocell/pkg/constant"
	"github.com/sivaosorg/gocell/pkg/utils"
	"github.com/sivaosorg/govm/blueprint"
	"github.com/sivaosorg/govm/bot/telegram"
	"github.com/sivaosorg/govm/builder"
	"github.com/sivaosorg/govm/charge"
	"github.com/sivaosorg/govm/logger"
	"github.com/sivaosorg/govm/timex"
	commonX "github.com/sivaosorg/govm/utils"
)

type responseWriterWrapper struct {
	gin.ResponseWriter
	bodyBuffer *bytes.Buffer
}

func (w responseWriterWrapper) Write(b []byte) (int, error) {
	w.bodyBuffer.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w responseWriterWrapper) WriteString(s string) (int, error) {
	w.bodyBuffer.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (m *MiddlewareManager) RequestMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		wrappedWriter := &responseWriterWrapper{
			ResponseWriter: c.Writer,
			bodyBuffer:     bytes.NewBufferString(""),
		}
		c.Writer = wrappedWriter
		// Make the request handling asynchronous
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Debugf("Asynchronous request recover handling: %v", r)
				}
			}()
			c.Next()
		}()
		m.async(c, wrappedWriter)
	}
}

func (m *MiddlewareManager) async(c *gin.Context, response *responseWriterWrapper) {
	tz := timex.With(time.Now()).Format(timex.DateTimeFormYearMonthDayHourMinuteSecond)
	raw, err := utils.SnapshotRequestBodyWith(c)
	go m.console(c, response, tz, raw, err)
	go m.notify(c, response, tz, raw, err)
}

func (m *MiddlewareManager) console(c *gin.Context, response *responseWriterWrapper, timestamp string, raw []byte, err error) {
	url := c.Request.URL.String()
	method := c.Request.Method
	builder := builder.NewMapBuilder()
	builder.
		Add("_tz", timestamp).
		Add("_method", method).
		Add("_endpoint", url).
		Add("_headers", utils.GetHeaders(c)).
		Add("_request", string(raw)).
		Add("_status_code", c.Writer.Status()).
		Add("_response", response.bodyBuffer.String())
	if err != nil {
		builder.Add("_error", err.Error())
	}
	logger.Debugf(builder.Json())
}

func (m *MiddlewareManager) notify(c *gin.Context, response *responseWriterWrapper, timestamp string, raw []byte, err error) {
	conf, err := m.conf.FindTelegramSeeker(constant.TelegramKeyTenant1)
	if err != nil {
		logger.Errorf("Telegram Bot Notify", err)
		return
	}
	if !conf.Config.IsEnabled {
		return
	}
	var builder strings.Builder
	url := c.Request.URL.String()
	method := c.Request.Method
	icon, _ := blueprint.TypeIcons[blueprint.TypeNotification]
	builder.WriteString(fmt.Sprintf("%v %s\n", icon, "Request Notify"))
	builder.WriteString(fmt.Sprintf("Tz: `%s`\n\n", timestamp))
	builder.WriteString(fmt.Sprintf("Method: `%s`\n", method))
	builder.WriteString(fmt.Sprintf("Endpoint: `%s`\n", url))
	builder.WriteString(fmt.Sprintf("Headers: \n`%s`\n", commonX.ToJson(utils.GetHeaders(c))))
	if m.allowRequestBody(c) {
		if charge.IsPostForms(c.Request) {
			builder.WriteString(fmt.Sprintf("Form-Data: \n`%s`\n", string(raw)))
		} else {
			builder.WriteString(fmt.Sprintf("Request-Body: \n`%s`\n", string(raw)))
		}
	}
	builder.WriteString(fmt.Sprintf("HTTP Code: `%d`\n", c.Writer.Status()))
	builder.WriteString(fmt.Sprintf("Response: \n`%s`\n", response.bodyBuffer.String()))
	if err != nil {
		builder.WriteString(fmt.Sprintf("Error(s): `%s`\n", err.Error()))
	}
	svc := telegram.NewTelegramService(conf.Config, conf.Option)
	svc.SendMessage(builder.String())
}

func (m *MiddlewareManager) allowRequestBody(c *gin.Context) bool {
	return (charge.IsPOST(c.Request) || charge.IsPUT(c.Request) || charge.IsPATCH(c.Request)) ||
		(charge.IsPostForms(c.Request) && charge.IsGET(c.Request))
}
