package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sivaosorg/gocell/internal/service"
	"github.com/sivaosorg/govm/entity"
	"github.com/sivaosorg/govm/wsconnx"

	"github.com/sivaosorg/wsconn/wsconn"
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

func (h *CommonHandler) OnPsqlStatus(ctx *gin.Context) {
	response := entity.NewResponseEntity().SetStatusCode(http.StatusOK).SetData(h.commonSvc.GetPsqlStatus())
	ctx.JSON(response.StatusCode, response)
	return
}

func (h *CommonHandler) SubscribeMessage(ctx *gin.Context) {
	h.wsSvc.SubscribeMessage(ctx)
}

func (h *CommonHandler) OnMessage(ctx *gin.Context) {
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
	return
}
