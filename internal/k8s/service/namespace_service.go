package service

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/k8s/client"
	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceService interface {
	ListNamespaces(ctx context.Context, req *model.GetNamespacesReq) ([]corev1.Namespace, error)
	CreateNamespace(ctx context.Context, req *model.CreateNamespaceReq) error
	DeleteNamespace(ctx context.Context, req *model.DeleteNamespaceReq) error
}

type namespaceService struct {
	k8sClient *client.K8sClient
	logger    *zap.Logger
}

func NewNamespaceService(k8sClient *client.K8sClient, logger *zap.Logger) NamespaceService {
	return &namespaceService{k8sClient: k8sClient, logger: logger}
}

func (s *namespaceService) ListNamespaces(ctx context.Context, req *model.GetNamespacesReq) ([]corev1.Namespace, error) {
	cs, err := s.k8sClient.GetClientSet(req.ClusterID)
	if err != nil {
		return nil, err
	}
	list, err := cs.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *namespaceService) CreateNamespace(ctx context.Context, req *model.CreateNamespaceReq) error {
	cs, err := s.k8sClient.GetClientSet(req.ClusterID)
	if err != nil {
		return err
	}
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   req.Name,
			Labels: req.Labels,
		},
	}
	_, err = cs.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	return err
}

func (s *namespaceService) DeleteNamespace(ctx context.Context, req *model.DeleteNamespaceReq) error {
	cs, err := s.k8sClient.GetClientSet(req.ClusterID)
	if err != nil {
		return err
	}
	return cs.CoreV1().Namespaces().Delete(ctx, req.Name, metav1.DeleteOptions{})
}
