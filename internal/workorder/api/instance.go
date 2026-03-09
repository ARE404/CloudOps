package api

import (
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/workorder/service"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/GoSimplicity/CloudOps/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type InstanceHandler struct {
	svc service.InstanceService
}

func NewInstanceHandler(svc service.InstanceService) *InstanceHandler {
	return &InstanceHandler{svc: svc}
}

func (h *InstanceHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/workorder/instances")
	{
		g.GET("", h.List)
		g.GET("/:id", h.Get)
		g.POST("", h.Create)
	}
}

func (h *InstanceHandler) List(ctx *gin.Context) {
	var req model.ListInstancesReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.svc.ListInstances(ctx, &req)
	})
}

func (h *InstanceHandler) Get(ctx *gin.Context) {
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return h.svc.GetInstance(ctx, id)
	})
}

func (h *InstanceHandler) Create(ctx *gin.Context) {
	var req model.CreateInstanceReq
	uc := ctx.MustGet("user").(jwt.UserClaims)
	req.ReporterID = uc.Uid
	req.ReporterName = uc.Username
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.svc.CreateInstance(ctx, &req)
	})
}
