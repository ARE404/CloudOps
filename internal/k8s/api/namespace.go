package api

import (
	"github.com/GoSimplicity/CloudOps/internal/k8s/service"
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/gin-gonic/gin"
)

type K8sNamespaceHandler struct {
	nsSvc service.NamespaceService
}

func NewK8sNamespaceHandler(nsSvc service.NamespaceService) *K8sNamespaceHandler {
	return &K8sNamespaceHandler{nsSvc: nsSvc}
}

func (h *K8sNamespaceHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/k8s/namespaces")
	{
		g.GET("", h.ListNamespaces)
		g.POST("", h.CreateNamespace)
		g.DELETE("/:name", h.DeleteNamespace)
	}
}

func (h *K8sNamespaceHandler) ListNamespaces(ctx *gin.Context) {
	var req model.GetNamespacesReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.nsSvc.ListNamespaces(ctx, &req)
	})
}

func (h *K8sNamespaceHandler) CreateNamespace(ctx *gin.Context) {
	var req model.CreateNamespaceReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.nsSvc.CreateNamespace(ctx, &req)
	})
}

func (h *K8sNamespaceHandler) DeleteNamespace(ctx *gin.Context) {
	var req model.DeleteNamespaceReq
	name, err := base.GetCustomParamID(ctx, "name")
	if err != nil {
		// name is a string, not ID; handle separately
		nameStr := ctx.Param("name")
		clusterID, _ := base.GetCustomParamID(ctx, "cluster_id")
		req.Name = nameStr
		req.ClusterID = clusterID
	} else {
		_ = name
	}
	// Re-read from query
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	req.Name = ctx.Param("name")
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return nil, h.nsSvc.DeleteNamespace(ctx, &req)
	})
}
