package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sivaosorg/gocell/internal/service"
	"github.com/sivaosorg/govm/entity"
)

type CommonHandler struct {
	commonSvc service.CommonService
}

func NewCommonHandler(commonSvc service.CommonService) *CommonHandler {
	h := &CommonHandler{
		commonSvc: commonSvc,
	}
	return h
}

func (h *CommonHandler) OnPsqlStatus(ctx *gin.Context) {
	response := entity.NewResponseEntity().SetStatusCode(http.StatusOK).SetData(h.commonSvc.GetPsqlStatus())
	ctx.JSON(response.StatusCode, response)
	return
}
