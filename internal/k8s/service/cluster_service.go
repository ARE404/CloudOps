package service

import (
	"context"
	"fmt"

	"github.com/GoSimplicity/CloudOps/internal/k8s/client"
	"github.com/GoSimplicity/CloudOps/internal/k8s/dao"
	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
)

type ClusterService interface {
	ListClusters(ctx context.Context, req *model.ListClustersReq) (*model.PageResp, error)
	GetClusterByID(ctx context.Context, req *model.GetClusterReq) (*model.K8sCluster, error)
	CreateCluster(ctx context.Context, req *model.CreateClusterReq) error
	UpdateCluster(ctx context.Context, req *model.UpdateClusterReq) error
	DeleteCluster(ctx context.Context, req *model.DeleteClusterReq) error
	RefreshClusterStatus(ctx context.Context, req *model.RefreshClusterReq) error
}

type clusterService struct {
	dao       dao.ClusterDAO
	k8sClient *client.K8sClient
	logger    *zap.Logger
}

func NewClusterService(dao dao.ClusterDAO, k8sClient *client.K8sClient, logger *zap.Logger) ClusterService {
	return &clusterService{dao: dao, k8sClient: k8sClient, logger: logger}
}

func (s *clusterService) ListClusters(ctx context.Context, req *model.ListClustersReq) (*model.PageResp, error) {
	list, total, err := s.dao.ListClusters(ctx, req)
	if err != nil {
		return nil, err
	}
	return &model.PageResp{Total: total, List: list}, nil
}

func (s *clusterService) GetClusterByID(ctx context.Context, req *model.GetClusterReq) (*model.K8sCluster, error) {
	return s.dao.GetClusterByID(ctx, req.ID)
}

func (s *clusterService) CreateCluster(ctx context.Context, req *model.CreateClusterReq) error {
	cluster := &model.K8sCluster{
		Name:           req.Name,
		NameZh:         req.NameZh,
		Description:    req.Description,
		KubeConfigData: req.KubeConfigData,
		CreateUserID:   req.CreateUserID,
		CreateUserName: req.CreateUserName,
	}
	if err := s.dao.CreateCluster(ctx, cluster); err != nil {
		return err
	}
	// 异步建立 k8s 连接
	go func() {
		if err := s.k8sClient.AddOrUpdateCluster(cluster.ID, req.KubeConfigData); err != nil {
			s.logger.Error("连接 k8s 集群失败", zap.Int("id", cluster.ID), zap.Error(err))
		}
	}()
	return nil
}

func (s *clusterService) UpdateCluster(ctx context.Context, req *model.UpdateClusterReq) error {
	cluster, err := s.dao.GetClusterByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("集群不存在: %w", err)
	}
	if req.NameZh != "" {
		cluster.NameZh = req.NameZh
	}
	if req.Description != "" {
		cluster.Description = req.Description
	}
	if req.KubeConfigData != "" {
		cluster.KubeConfigData = req.KubeConfigData
		go func() {
			if err := s.k8sClient.AddOrUpdateCluster(cluster.ID, req.KubeConfigData); err != nil {
				s.logger.Error("更新 k8s 集群连接失败", zap.Int("id", cluster.ID), zap.Error(err))
			}
		}()
	}
	return s.dao.UpdateCluster(ctx, cluster)
}

func (s *clusterService) DeleteCluster(ctx context.Context, req *model.DeleteClusterReq) error {
	if err := s.dao.DeleteCluster(ctx, req.ID); err != nil {
		return err
	}
	s.k8sClient.RemoveCluster(req.ID)
	return nil
}

func (s *clusterService) RefreshClusterStatus(ctx context.Context, req *model.RefreshClusterReq) error {
	cluster, err := s.dao.GetClusterByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if err := s.k8sClient.AddOrUpdateCluster(cluster.ID, cluster.KubeConfigData); err != nil {
		cluster.Status = 2
		_ = s.dao.UpdateCluster(ctx, cluster)
		return err
	}
	cluster.Status = 1
	return s.dao.UpdateCluster(ctx, cluster)
}
