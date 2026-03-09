package api

import (
	"github.com/GoSimplicity/CloudOps/internal/k8s/service"
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/gin-gonic/gin"
)

type K8sSvcHandler struct {
	svcService service.ServiceService
}

func NewK8sSvcHandler(svcService service.ServiceService) *K8sSvcHandler {
	return &K8sSvcHandler{svcService: svcService}
}

func (h *K8sSvcHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/k8s/services")
	{
		g.GET("", h.ListServices)
	}
}

func (h *K8sSvcHandler) ListServices(ctx *gin.Context) {
	var req model.GetServicesReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.svcService.ListServices(ctx, &req)
	})
}
