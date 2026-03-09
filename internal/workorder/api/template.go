package api

import (
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/workorder/service"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	svc service.TemplateService
}

func NewTemplateHandler(svc service.TemplateService) *TemplateHandler {
	return &TemplateHandler{svc: svc}
}

func (h *TemplateHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/workorder/templates")
	{
		g.GET("", h.List)
		g.POST("", h.Create)
		g.PUT("/:id", h.Update)
		g.DELETE("/:id", h.Delete)
	}
}

func (h *TemplateHandler) List(ctx *gin.Context) {
	var req model.ListTemplatesReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.svc.ListTemplates(ctx, &req)
	})
}

func (h *TemplateHandler) Create(ctx *gin.Context) {
	var req model.CreateTemplateReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.svc.CreateTemplate(ctx, &req)
	})
}

func (h *TemplateHandler) Update(ctx *gin.Context) {
	var req model.UpdateTemplateReq
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	req.ID = id
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.svc.UpdateTemplate(ctx, &req)
	})
}

func (h *TemplateHandler) Delete(ctx *gin.Context) {
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return nil, h.svc.DeleteTemplate(ctx, id)
	})
}
