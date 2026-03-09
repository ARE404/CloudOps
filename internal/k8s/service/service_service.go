package service

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/k8s/client"
	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceService interface {
	ListServices(ctx context.Context, req *model.GetServicesReq) ([]corev1.Service, error)
}

type serviceService struct {
	k8sClient *client.K8sClient
	logger    *zap.Logger
}

func NewServiceService(k8sClient *client.K8sClient, logger *zap.Logger) ServiceService {
	return &serviceService{k8sClient: k8sClient, logger: logger}
}

func (s *serviceService) ListServices(ctx context.Context, req *model.GetServicesReq) ([]corev1.Service, error) {
	cs, err := s.k8sClient.GetClientSet(req.ClusterID)
	if err != nil {
		return nil, err
	}
	ns := req.Namespace
	if ns == "" {
		ns = metav1.NamespaceAll
	}
	list, err := cs.CoreV1().Services(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
