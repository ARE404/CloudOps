package api

import (
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/prometheus/service"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/gin-gonic/gin"
)

type AlertRuleHandler struct {
	svc service.AlertRuleService
}

func NewAlertRuleHandler(svc service.AlertRuleService) *AlertRuleHandler {
	return &AlertRuleHandler{svc: svc}
}

func (h *AlertRuleHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/prometheus/alert-rules")
	{
		g.GET("", h.List)
		g.POST("", h.Create)
		g.PUT("/:id", h.Update)
		g.DELETE("/:id", h.Delete)
	}
}

func (h *AlertRuleHandler) List(ctx *gin.Context) {
	var req model.ListAlertRulesReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.svc.ListAlertRules(ctx, &req)
	})
}

func (h *AlertRuleHandler) Create(ctx *gin.Context) {
	var req model.CreateAlertRuleReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.svc.CreateAlertRule(ctx, &req)
	})
}

func (h *AlertRuleHandler) Update(ctx *gin.Context) {
	var req model.UpdateAlertRuleReq
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	req.ID = id
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.svc.UpdateAlertRule(ctx, &req)
	})
}

func (h *AlertRuleHandler) Delete(ctx *gin.Context) {
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return nil, h.svc.DeleteAlertRule(ctx, id)
	})
}
