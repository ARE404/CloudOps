package api

import (
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/system/service"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleSvc service.RoleService
}

func NewRoleHandler(roleSvc service.RoleService) *RoleHandler {
	return &RoleHandler{roleSvc: roleSvc}
}

func (h *RoleHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/roles")
	{
		g.GET("", h.ListRoles)
		g.POST("", h.CreateRole)
		g.PUT("/:id", h.UpdateRole)
		g.DELETE("/:id", h.DeleteRole)
	}
}

func (h *RoleHandler) ListRoles(ctx *gin.Context) {
	var req model.ListRolesReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.roleSvc.ListRoles(ctx, &req)
	})
}

func (h *RoleHandler) CreateRole(ctx *gin.Context) {
	var req model.CreateRoleReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.roleSvc.CreateRole(ctx, &req)
	})
}

func (h *RoleHandler) UpdateRole(ctx *gin.Context) {
	var req model.UpdateRoleReq
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	req.ID = id
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.roleSvc.UpdateRole(ctx, &req)
	})
}

func (h *RoleHandler) DeleteRole(ctx *gin.Context) {
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return nil, h.roleSvc.DeleteRole(ctx, id)
	})
}
