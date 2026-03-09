package model

// K8sCluster K8s集群表
type K8sCluster struct {
	Model
	Name           string `json:"name" gorm:"uniqueIndex;size:64;not null"`
	NameZh         string `json:"name_zh" gorm:"size:64"`
	Description    string `json:"description" gorm:"size:255"`
	KubeConfigData string `json:"kube_config_data" gorm:"type:text"`
	ApiServerAddr  string `json:"api_server_addr" gorm:"size:255"`
	Status         int8   `json:"status" gorm:"default:0;comment:0=unknown,1=active,2=error"`
	Version        string `json:"version" gorm:"size:32"`
	NodeCount      int    `json:"node_count" gorm:"-"`
	CreateUserID   int    `json:"create_user_id"`
	CreateUserName string `json:"create_user_name" gorm:"size:64"`
}

// --- Requests ---

type ListClustersReq struct {
	PageReq
	Name string `json:"name" form:"name"`
}

type GetClusterReq struct {
	ID int `json:"-"`
}

type CreateClusterReq struct {
	Name           string `json:"name" binding:"required"`
	NameZh         string `json:"name_zh"`
	Description    string `json:"description"`
	KubeConfigData string `json:"kube_config_data" binding:"required"`
	CreateUserID   int    `json:"-"`
	CreateUserName string `json:"-"`
}

type UpdateClusterReq struct {
	ID             int    `json:"-"`
	NameZh         string `json:"name_zh"`
	Description    string `json:"description"`
	KubeConfigData string `json:"kube_config_data"`
}

type DeleteClusterReq struct {
	ID int `json:"-"`
}

type RefreshClusterReq struct {
	ID int `json:"-"`
}

// K8s 资源通用请求

type GetNamespacesReq struct {
	ClusterID int    `json:"cluster_id" form:"cluster_id" binding:"required"`
	Name      string `json:"name" form:"name"`
}

type GetDeploymentsReq struct {
	ClusterID int    `json:"cluster_id" form:"cluster_id" binding:"required"`
	Namespace string `json:"namespace" form:"namespace"`
}

type CreateDeploymentReq struct {
	ClusterID   int               `json:"cluster_id" binding:"required"`
	Namespace   string            `json:"namespace" binding:"required"`
	Name        string            `json:"name" binding:"required"`
	Image       string            `json:"image" binding:"required"`
	Replicas    int32             `json:"replicas"`
	Labels      map[string]string `json:"labels"`
	EnvVars     map[string]string `json:"env_vars"`
	CPULimit    string            `json:"cpu_limit"`
	MemoryLimit string            `json:"memory_limit"`
}

type ScaleDeploymentReq struct {
	ClusterID int    `json:"cluster_id" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Replicas  int32  `json:"replicas" binding:"required"`
}

type DeleteDeploymentReq struct {
	ClusterID int    `json:"cluster_id" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
	Name      string `json:"name" binding:"required"`
}

type GetPodsReq struct {
	ClusterID int    `json:"cluster_id" form:"cluster_id" binding:"required"`
	Namespace string `json:"namespace" form:"namespace"`
}

type GetPodLogsReq struct {
	ClusterID int    `json:"cluster_id" form:"cluster_id" binding:"required"`
	Namespace string `json:"namespace" form:"namespace" binding:"required"`
	Name      string `json:"name" form:"name" binding:"required"`
	Container string `json:"container" form:"container"`
	Lines     int64  `json:"lines" form:"lines"`
}

type DeletePodReq struct {
	ClusterID int    `json:"cluster_id" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
	Name      string `json:"name" binding:"required"`
}

type GetServicesReq struct {
	ClusterID int    `json:"cluster_id" form:"cluster_id" binding:"required"`
	Namespace string `json:"namespace" form:"namespace"`
}

type GetConfigMapsReq struct {
	ClusterID int    `json:"cluster_id" form:"cluster_id" binding:"required"`
	Namespace string `json:"namespace" form:"namespace"`
}

type CreateNamespaceReq struct {
	ClusterID int               `json:"cluster_id" binding:"required"`
	Name      string            `json:"name" binding:"required"`
	Labels    map[string]string `json:"labels"`
}

type DeleteNamespaceReq struct {
	ClusterID int    `json:"cluster_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
}
