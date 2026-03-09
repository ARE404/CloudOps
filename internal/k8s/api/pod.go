package api

import (
	"github.com/GoSimplicity/CloudOps/internal/k8s/service"
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/gin-gonic/gin"
)

type K8sPodHandler struct {
	podSvc service.PodService
}

func NewK8sPodHandler(podSvc service.PodService) *K8sPodHandler {
	return &K8sPodHandler{podSvc: podSvc}
}

func (h *K8sPodHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/k8s/pods")
	{
		g.GET("", h.ListPods)
		g.GET("/logs", h.GetPodLogs)
		g.DELETE("", h.DeletePod)
	}
}

func (h *K8sPodHandler) ListPods(ctx *gin.Context) {
	var req model.GetPodsReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.podSvc.ListPods(ctx, &req)
	})
}

func (h *K8sPodHandler) GetPodLogs(ctx *gin.Context) {
	var req model.GetPodLogsReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.podSvc.GetPodLogs(ctx, &req)
	})
}

func (h *K8sPodHandler) DeletePod(ctx *gin.Context) {
	var req model.DeletePodReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.podSvc.DeletePod(ctx, &req)
	})
}
