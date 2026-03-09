package api

import (
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/prometheus/service"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/gin-gonic/gin"
)

type AlertEventHandler struct {
	svc service.AlertEventService
}

func NewAlertEventHandler(svc service.AlertEventService) *AlertEventHandler {
	return &AlertEventHandler{svc: svc}
}

func (h *AlertEventHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/prometheus/alert-events")
	{
		g.GET("", h.List)
		g.POST("/receive", h.Receive)
	}
}

func (h *AlertEventHandler) List(ctx *gin.Context) {
	var req model.ListAlertEventsReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.svc.ListAlertEvents(ctx, &req)
	})
}

func (h *AlertEventHandler) Receive(ctx *gin.Context) {
	var event model.PromAlertEvent
	base.HandleRequest(ctx, &event, func() (interface{}, error) {
		return nil, h.svc.ReceiveAlert(ctx, &event)
	})
}
