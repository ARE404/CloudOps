package service

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/k8s/client"
	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

type DeploymentService interface {
	ListDeployments(ctx context.Context, req *model.GetDeploymentsReq) ([]appsv1.Deployment, error)
	CreateDeployment(ctx context.Context, req *model.CreateDeploymentReq) error
	ScaleDeployment(ctx context.Context, req *model.ScaleDeploymentReq) error
	DeleteDeployment(ctx context.Context, req *model.DeleteDeploymentReq) error
}

type deploymentService struct {
	k8sClient *client.K8sClient
	logger    *zap.Logger
}

func NewDeploymentService(k8sClient *client.K8sClient, logger *zap.Logger) DeploymentService {
	return &deploymentService{k8sClient: k8sClient, logger: logger}
}

func (s *deploymentService) ListDeployments(ctx context.Context, req *model.GetDeploymentsReq) ([]appsv1.Deployment, error) {
	cs, err := s.k8sClient.GetClientSet(req.ClusterID)
	if err != nil {
		return nil, err
	}
	ns := req.Namespace
	if ns == "" {
		ns = metav1.NamespaceAll
	}
	list, err := cs.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *deploymentService) CreateDeployment(ctx context.Context, req *model.CreateDeploymentReq) error {
	cs, err := s.k8sClient.GetClientSet(req.ClusterID)
	if err != nil {
		return err
	}

	replicas := req.Replicas
	if replicas == 0 {
		replicas = 1
	}

	// Build env vars
	envVars := make([]corev1.EnvVar, 0, len(req.EnvVars))
	for k, v := range req.EnvVars {
		envVars = append(envVars, corev1.EnvVar{Name: k, Value: v})
	}

	// Build resource limits
	resources := corev1.ResourceRequirements{}
	if req.CPULimit != "" || req.MemoryLimit != "" {
		limits := corev1.ResourceList{}
		if req.CPULimit != "" {
			limits[corev1.ResourceCPU] = resource.MustParse(req.CPULimit)
		}
		if req.MemoryLimit != "" {
			limits[corev1.ResourceMemory] = resource.MustParse(req.MemoryLimit)
		}
		resources.Limits = limits
	}

	labels := req.Labels
	if labels == nil {
		labels = map[string]string{"app": req.Name}
	}

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.To(replicas),
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:      req.Name,
							Image:     req.Image,
							Env:       envVars,
							Resources: resources,
						},
					},
				},
			},
		},
	}

	_, err = cs.AppsV1().Deployments(req.Namespace).Create(ctx, deploy, metav1.CreateOptions{})
	return err
}

func (s *deploymentService) ScaleDeployment(ctx context.Context, req *model.ScaleDeploymentReq) error {
	cs, err := s.k8sClient.GetClientSet(req.ClusterID)
	if err != nil {
		return err
	}
	scale, err := cs.AppsV1().Deployments(req.Namespace).GetScale(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	scale.Spec.Replicas = req.Replicas
	_, err = cs.AppsV1().Deployments(req.Namespace).UpdateScale(ctx, req.Name, scale, metav1.UpdateOptions{})
	return err
}

func (s *deploymentService) DeleteDeployment(ctx context.Context, req *model.DeleteDeploymentReq) error {
	cs, err := s.k8sClient.GetClientSet(req.ClusterID)
	if err != nil {
		return err
	}
	return cs.AppsV1().Deployments(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
}
