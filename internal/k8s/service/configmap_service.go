package service

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/k8s/client"
	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMapService interface {
	ListConfigMaps(ctx context.Context, req *model.GetConfigMapsReq) ([]corev1.ConfigMap, error)
}

type configMapService struct {
	k8sClient *client.K8sClient
	logger    *zap.Logger
}

func NewConfigMapService(k8sClient *client.K8sClient, logger *zap.Logger) ConfigMapService {
	return &configMapService{k8sClient: k8sClient, logger: logger}
}

func (s *configMapService) ListConfigMaps(ctx context.Context, req *model.GetConfigMapsReq) ([]corev1.ConfigMap, error) {
	cs, err := s.k8sClient.GetClientSet(req.ClusterID)
	if err != nil {
		return nil, err
	}
	ns := req.Namespace
	if ns == "" {
		ns = metav1.NamespaceAll
	}
	list, err := cs.CoreV1().ConfigMaps(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
