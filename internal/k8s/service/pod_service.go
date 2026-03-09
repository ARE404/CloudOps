package service

import (
	"context"
	"io"

	"github.com/GoSimplicity/CloudOps/internal/k8s/client"
	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodService interface {
	ListPods(ctx context.Context, req *model.GetPodsReq) ([]corev1.Pod, error)
	GetPodLogs(ctx context.Context, req *model.GetPodLogsReq) (string, error)
	DeletePod(ctx context.Context, req *model.DeletePodReq) error
}

type podService struct {
	k8sClient *client.K8sClient
	logger    *zap.Logger
}

func NewPodService(k8sClient *client.K8sClient, logger *zap.Logger) PodService {
	return &podService{k8sClient: k8sClient, logger: logger}
}

func (s *podService) ListPods(ctx context.Context, req *model.GetPodsReq) ([]corev1.Pod, error) {
	cs, err := s.k8sClient.GetClientSet(req.ClusterID)
	if err != nil {
		return nil, err
	}
	ns := req.Namespace
	if ns == "" {
		ns = metav1.NamespaceAll
	}
	list, err := cs.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *podService) GetPodLogs(ctx context.Context, req *model.GetPodLogsReq) (string, error) {
	cs, err := s.k8sClient.GetClientSet(req.ClusterID)
	if err != nil {
		return "", err
	}
	lines := req.Lines
	if lines == 0 {
		lines = 100
	}
	opts := &corev1.PodLogOptions{
		Container: req.Container,
		TailLines: &lines,
	}
	stream, err := cs.CoreV1().Pods(req.Namespace).GetLogs(req.Name, opts).Stream(ctx)
	if err != nil {
		return "", err
	}
	defer stream.Close()
	data, err := io.ReadAll(stream)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *podService) DeletePod(ctx context.Context, req *model.DeletePodReq) error {
	cs, err := s.k8sClient.GetClientSet(req.ClusterID)
	if err != nil {
		return err
	}
	return cs.CoreV1().Pods(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
}
