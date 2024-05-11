package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sivaosorg/gocell/internal/middlewares"
	"github.com/sivaosorg/gocell/internal/service"
	"github.com/sivaosorg/govm/entity"
	"github.com/sivaosorg/govm/wsconnx"

	"github.com/sivaosorg/wsconn"
)

type CommonHandler struct {
	commonSvc service.CommonService
	wsSvc     wsconn.WebsocketService
}

func NewCommonHandler(commonSvc service.CommonService) *CommonHandler {
	h := &CommonHandler{
		commonSvc: commonSvc,
		wsSvc:     wsconn.NewWebsocketService(wsconn.NewWebsocket()),
	}
	return h
}

func (c *CommonHandler) Router(r *gin.RouterGroup, middlewares *middlewares.MiddlewareManager) *gin.RouterGroup {
	r.GET("/psql-status", c.OnPsqlStatus) // endpoint: http://127.0.0.1:8081/api/v1/common/psql-status
	r.GET("/consumer", c.OnSubscribe)     // endpoint: ws://127.0.0.1:8081/api/v1/common/consumer
	r.POST("/producer", c.OnProduce)      // endpoint: http://127.0.0.1:8081/api/v1/common/producer
	return r
}

func (h *CommonHandler) OnPsqlStatus(ctx *gin.Context) {
	data := h.commonSvc.GetPsqlStatus()
	response := entity.NewResponseEntity().SetData(data)
	if data.IsConnected {
		response.SetStatusCode(http.StatusOK)
	} else {
		response.SetStatusCode(http.StatusInternalServerError)
	}
	ctx.JSON(response.StatusCode, response)
}

func (h *CommonHandler) OnSubscribe(ctx *gin.Context) {
	h.wsSvc.SubscribeMessage(ctx)
}

func (h *CommonHandler) OnProduce(ctx *gin.Context) {
	response := entity.NewResponseEntity()
	var message wsconnx.WsConnMessagePayload
	message.SetGenesisTimestamp(time.Now())
	if err := ctx.ShouldBindJSON(&message); err != nil {
		response.SetStatusCode(http.StatusBadRequest).SetMessage(err.Error()).SetError(err)
		ctx.JSON(response.StatusCode, response)
		return
	}
	go h.wsSvc.BroadcastMessage(message)
	response.SetStatusCode(http.StatusOK).SetMessage("Message sent successfully").SetData(message)
	ctx.JSON(response.StatusCode, response)
}
