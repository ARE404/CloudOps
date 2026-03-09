package api

import (
	"github.com/GoSimplicity/CloudOps/internal/k8s/service"
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/GoSimplicity/CloudOps/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type K8sClusterHandler struct {
	clusterSvc service.ClusterService
}

func NewK8sClusterHandler(clusterSvc service.ClusterService) *K8sClusterHandler {
	return &K8sClusterHandler{clusterSvc: clusterSvc}
}

func (h *K8sClusterHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/k8s/clusters")
	{
		g.GET("", h.ListClusters)
		g.GET("/:id", h.GetCluster)
		g.POST("", h.CreateCluster)
		g.PUT("/:id", h.UpdateCluster)
		g.DELETE("/:id", h.DeleteCluster)
		g.POST("/:id/refresh", h.RefreshCluster)
	}
}

// ListClusters godoc
// @Summary 获取集群列表
// @Tags K8s集群管理
// @Accept json
// @Produce json
// @Success 200 {object} base.ApiResponse
// @Router /api/k8s/clusters [get]
func (h *K8sClusterHandler) ListClusters(ctx *gin.Context) {
	var req model.ListClustersReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.clusterSvc.ListClusters(ctx, &req)
	})
}

func (h *K8sClusterHandler) GetCluster(ctx *gin.Context) {
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return h.clusterSvc.GetClusterByID(ctx, &model.GetClusterReq{ID: id})
	})
}

func (h *K8sClusterHandler) CreateCluster(ctx *gin.Context) {
	var req model.CreateClusterReq
	uc := ctx.MustGet("user").(jwt.UserClaims)
	req.CreateUserID = uc.Uid
	req.CreateUserName = uc.Username
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.clusterSvc.CreateCluster(ctx, &req)
	})
}

func (h *K8sClusterHandler) UpdateCluster(ctx *gin.Context) {
	var req model.UpdateClusterReq
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	req.ID = id
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.clusterSvc.UpdateCluster(ctx, &req)
	})
}

func (h *K8sClusterHandler) DeleteCluster(ctx *gin.Context) {
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return nil, h.clusterSvc.DeleteCluster(ctx, &model.DeleteClusterReq{ID: id})
	})
}

func (h *K8sClusterHandler) RefreshCluster(ctx *gin.Context) {
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return nil, h.clusterSvc.RefreshClusterStatus(ctx, &model.RefreshClusterReq{ID: id})
	})
}
