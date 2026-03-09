package api

import (
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/workorder/service"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/GoSimplicity/CloudOps/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type FlowHandler struct {
	svc service.FlowService
}

func NewFlowHandler(svc service.FlowService) *FlowHandler {
	return &FlowHandler{svc: svc}
}

func (h *FlowHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/workorder/instances/:id/flows")
	{
		g.GET("", h.GetFlows)
		g.POST("/action", h.ExecuteAction)
	}
}

func (h *FlowHandler) GetFlows(ctx *gin.Context) {
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return h.svc.GetFlows(ctx, id)
	})
}

func (h *FlowHandler) ExecuteAction(ctx *gin.Context) {
	var req model.FlowActionReq
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	req.InstanceID = id
	uc := ctx.MustGet("user").(jwt.UserClaims)
	req.OperatorID = uc.Uid
	req.OperatorName = uc.Username
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.svc.ExecuteAction(ctx, &req)
	})
}
