package api

import (
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/prometheus/service"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/gin-gonic/gin"
)

type ScrapePoolHandler struct {
	svc service.ScrapePoolService
}

func NewScrapePoolHandler(svc service.ScrapePoolService) *ScrapePoolHandler {
	return &ScrapePoolHandler{svc: svc}
}

func (h *ScrapePoolHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/prometheus/scrape-pools")
	{
		g.GET("", h.List)
		g.POST("", h.Create)
		g.PUT("/:id", h.Update)
		g.DELETE("/:id", h.Delete)
	}
}

func (h *ScrapePoolHandler) List(ctx *gin.Context) {
	var req model.ListScrapePoolsReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.svc.ListScrapePools(ctx, &req)
	})
}

func (h *ScrapePoolHandler) Create(ctx *gin.Context) {
	var req model.CreateScrapePoolReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.svc.CreateScrapePool(ctx, &req)
	})
}

func (h *ScrapePoolHandler) Update(ctx *gin.Context) {
	var req model.UpdateScrapePoolReq
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	req.ID = id
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.svc.UpdateScrapePool(ctx, &req)
	})
}

func (h *ScrapePoolHandler) Delete(ctx *gin.Context) {
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return nil, h.svc.DeleteScrapePool(ctx, id)
	})
}
