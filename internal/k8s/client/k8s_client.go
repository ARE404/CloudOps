package client

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// K8sClient 管理多集群 clientset，以 clusterID 为 key
type K8sClient struct {
	mu      sync.RWMutex
	clients map[int]*kubernetes.Clientset
	logger  *zap.Logger
}

func NewK8sClient(logger *zap.Logger) *K8sClient {
	return &K8sClient{
		clients: make(map[int]*kubernetes.Clientset),
		logger:  logger,
	}
}

// AddOrUpdateCluster 加载或更新 kubeconfig，建立 clientset
func (c *K8sClient) AddOrUpdateCluster(clusterID int, kubeConfigData string) error {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeConfigData))
	if err != nil {
		return fmt.Errorf("解析 kubeconfig 失败: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("创建 k8s clientset 失败: %w", err)
	}

	c.mu.Lock()
	c.clients[clusterID] = clientset
	c.mu.Unlock()

	c.logger.Info("k8s 集群已连接", zap.Int("cluster_id", clusterID))
	return nil
}

// GetClientSet 获取指定集群的 clientset
func (c *K8sClient) GetClientSet(clusterID int) (*kubernetes.Clientset, error) {
	c.mu.RLock()
	clientset, ok := c.clients[clusterID]
	c.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("集群 %d 未连接，请先添加集群", clusterID)
	}
	return clientset, nil
}

// RemoveCluster 移除集群连接
func (c *K8sClient) RemoveCluster(clusterID int) {
	c.mu.Lock()
	delete(c.clients, clusterID)
	c.mu.Unlock()
}
