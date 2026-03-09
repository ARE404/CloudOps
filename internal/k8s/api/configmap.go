package api

import (
	"github.com/GoSimplicity/CloudOps/internal/k8s/service"
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/gin-gonic/gin"
)

type K8sConfigMapHandler struct {
	cmSvc service.ConfigMapService
}

func NewK8sConfigMapHandler(cmSvc service.ConfigMapService) *K8sConfigMapHandler {
	return &K8sConfigMapHandler{cmSvc: cmSvc}
}

func (h *K8sConfigMapHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/k8s/configmaps")
	{
		g.GET("", h.ListConfigMaps)
	}
}

func (h *K8sConfigMapHandler) ListConfigMaps(ctx *gin.Context) {
	var req model.GetConfigMapsReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.cmSvc.ListConfigMaps(ctx, &req)
	})
}
