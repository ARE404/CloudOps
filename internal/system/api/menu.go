package api

import (
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/system/service"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuSvc service.MenuService
}

func NewMenuHandler(menuSvc service.MenuService) *MenuHandler {
	return &MenuHandler{menuSvc: menuSvc}
}

func (h *MenuHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/menus")
	{
		g.GET("", h.ListMenus)
		g.POST("", h.CreateMenu)
		g.PUT("/:id", h.UpdateMenu)
		g.DELETE("/:id", h.DeleteMenu)
	}
}

func (h *MenuHandler) ListMenus(ctx *gin.Context) {
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return h.menuSvc.ListMenus(ctx)
	})
}

func (h *MenuHandler) CreateMenu(ctx *gin.Context) {
	var req model.CreateMenuReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.menuSvc.CreateMenu(ctx, &req)
	})
}

func (h *MenuHandler) UpdateMenu(ctx *gin.Context) {
	var req model.UpdateMenuReq
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	req.ID = id
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.menuSvc.UpdateMenu(ctx, &req)
	})
}

func (h *MenuHandler) DeleteMenu(ctx *gin.Context) {
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return nil, h.menuSvc.DeleteMenu(ctx, id)
	})
}
