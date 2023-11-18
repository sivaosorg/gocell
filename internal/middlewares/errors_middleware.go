package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sivaosorg/gocell/pkg/constant"
	"github.com/sivaosorg/govm/blueprint"
	"github.com/sivaosorg/govm/bot/telegram"
	"github.com/sivaosorg/govm/entity"
	"github.com/sivaosorg/govm/logger"
	"github.com/sivaosorg/govm/timex"
)

func (m *MiddlewareManager) ErrorMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.Debugf("Recovered errors from panic:", err)
			m.defaultHandleRecovery(c, err)
			return
		}
		if len(c.Errors) > 0 {
			logger.Debugf("Business logic errors:", c.Errors)
			response := entity.NewResponseEntity().BadRequest(http.StatusText(http.StatusBadRequest), nil)
			response.SetErrors(fmt.Sprint(c.Errors))
			m.notification(c, c.Errors, response.StatusCode)
			c.JSON(response.StatusCode, response)
			return
		}
	}()
	c.Next()
}

func (m *MiddlewareManager) notification(c *gin.Context, err any, status int) {
	if !m.conf.AvailableTelegramSeekers() {
		return
	}
	var builder strings.Builder
	icon, _ := blueprint.TypeIcons[blueprint.TypeError]
	builder.WriteString(fmt.Sprintf("%v %s\n", icon, "Core Application Recovery"))
	builder.WriteString(fmt.Sprintf("Tz: %s\n", time.Now().Format(timex.DateTimeFormYearMonthDayHourMinuteSecond)))
	builder.WriteString(fmt.Sprintf("URL: `%s`\n", c.Request.RequestURI))
	builder.WriteString(fmt.Sprintf("Status Code: %d\n", status))
	builder.WriteString(fmt.Sprintf("Message: %s\n", http.StatusText(status)))
	builder.WriteString(fmt.Sprintf("Error(R): `%s`\n", fmt.Sprint(err)))
	conf, e := m.conf.FindTelegramSeeker(constant.TelegramKeyTenant2)
	if e != nil {
		return
	}
	svc := telegram.NewTelegramService(conf.Config, conf.Option)
	go svc.SendMessage(builder.String())
}
