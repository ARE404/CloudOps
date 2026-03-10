import request from './request'

export interface K8sCluster {
  id: number
  name: string
  kubeConfig: string
  status: number
  createdAt: string
}

export interface K8sNamespace {
  name: string
  status: string
  createdAt: string
}

export interface K8sDeployment {
  name: string
  namespace: string
  replicas: number
  readyReplicas: number
  image: string
  createdAt: string
}

export interface K8sPod {
  name: string
  namespace: string
  status: string
  nodeName: string
  ip: string
  createdAt: string
}

export interface K8sService {
  name: string
  namespace: string
  type: string
  clusterIP: string
  ports: string
  createdAt: string
}

export interface K8sConfigMap {
  name: string
  namespace: string
  dataCount: number
  createdAt: string
}

export const k8sApi = {
  // Clusters
  listClusters: (): Promise<K8sCluster[]> =>
    request.get('/api/k8s/clusters'),
  createCluster: (data: Partial<K8sCluster>): Promise<K8sCluster> =>
    request.post('/api/k8s/clusters', data),
  updateCluster: (id: number, data: Partial<K8sCluster>): Promise<K8sCluster> =>
    request.put(`/api/k8s/clusters/${id}`, data),
  deleteCluster: (id: number): Promise<void> =>
    request.delete(`/api/k8s/clusters/${id}`),

  // Namespaces
  listNamespaces: (clusterId: number): Promise<K8sNamespace[]> =>
    request.get('/api/k8s/namespaces', { params: { clusterId } }),
  createNamespace: (data: { clusterId: number; name: string }): Promise<void> =>
    request.post('/api/k8s/namespaces', data),
  deleteNamespace: (clusterId: number, name: string): Promise<void> =>
    request.delete(`/api/k8s/namespaces/${name}`, { params: { clusterId } }),

  // Deployments
  listDeployments: (params: { clusterId: number; namespace?: string }): Promise<K8sDeployment[]> =>
    request.get('/api/k8s/deployments', { params }),
  scaleDeployment: (data: { clusterId: number; namespace: string; name: string; replicas: number }): Promise<void> =>
    request.post('/api/k8s/deployments/scale', data),
  deleteDeployment: (params: { clusterId: number; namespace: string; name: string }): Promise<void> =>
    request.delete('/api/k8s/deployments', { params }),

  // Pods
  listPods: (params: { clusterId: number; namespace?: string }): Promise<K8sPod[]> =>
    request.get('/api/k8s/pods', { params }),
  getPodLogs: (params: { clusterId: number; namespace: string; name: string }): Promise<string> =>
    request.get('/api/k8s/pods/logs', { params }),
  deletePod: (params: { clusterId: number; namespace: string; name: string }): Promise<void> =>
    request.delete('/api/k8s/pods', { params }),

  // Services
  listServices: (params: { clusterId: number; namespace?: string }): Promise<K8sService[]> =>
    request.get('/api/k8s/services', { params }),

  // ConfigMaps
  listConfigMaps: (params: { clusterId: number; namespace?: string }): Promise<K8sConfigMap[]> =>
    request.get('/api/k8s/configmaps', { params }),
}
