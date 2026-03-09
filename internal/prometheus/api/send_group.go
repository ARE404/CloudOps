package api

import (
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/prometheus/service"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/gin-gonic/gin"
)

type SendGroupHandler struct {
	svc service.SendGroupService
}

func NewSendGroupHandler(svc service.SendGroupService) *SendGroupHandler {
	return &SendGroupHandler{svc: svc}
}

func (h *SendGroupHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/prometheus/send-groups")
	{
		g.GET("", h.List)
		g.POST("", h.Create)
		g.PUT("/:id", h.Update)
		g.DELETE("/:id", h.Delete)
	}
}

func (h *SendGroupHandler) List(ctx *gin.Context) {
	var req model.ListSendGroupsReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.svc.ListSendGroups(ctx, &req)
	})
}

func (h *SendGroupHandler) Create(ctx *gin.Context) {
	var req model.CreateSendGroupReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.svc.CreateSendGroup(ctx, &req)
	})
}

func (h *SendGroupHandler) Update(ctx *gin.Context) {
	var req model.UpdateSendGroupReq
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	req.ID = id
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.svc.UpdateSendGroup(ctx, &req)
	})
}

func (h *SendGroupHandler) Delete(ctx *gin.Context) {
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return nil, h.svc.DeleteSendGroup(ctx, id)
	})
}
