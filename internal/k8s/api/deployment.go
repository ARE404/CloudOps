package api

import (
	"github.com/GoSimplicity/CloudOps/internal/k8s/service"
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/gin-gonic/gin"
)

type K8sDeploymentHandler struct {
	deploySvc service.DeploymentService
}

func NewK8sDeploymentHandler(deploySvc service.DeploymentService) *K8sDeploymentHandler {
	return &K8sDeploymentHandler{deploySvc: deploySvc}
}

func (h *K8sDeploymentHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/k8s/deployments")
	{
		g.GET("", h.ListDeployments)
		g.POST("", h.CreateDeployment)
		g.POST("/scale", h.ScaleDeployment)
		g.DELETE("", h.DeleteDeployment)
	}
}

func (h *K8sDeploymentHandler) ListDeployments(ctx *gin.Context) {
	var req model.GetDeploymentsReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.deploySvc.ListDeployments(ctx, &req)
	})
}

func (h *K8sDeploymentHandler) CreateDeployment(ctx *gin.Context) {
	var req model.CreateDeploymentReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.deploySvc.CreateDeployment(ctx, &req)
	})
}

func (h *K8sDeploymentHandler) ScaleDeployment(ctx *gin.Context) {
	var req model.ScaleDeploymentReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.deploySvc.ScaleDeployment(ctx, &req)
	})
}

func (h *K8sDeploymentHandler) DeleteDeployment(ctx *gin.Context) {
	var req model.DeleteDeploymentReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.deploySvc.DeleteDeployment(ctx, &req)
	})
}
