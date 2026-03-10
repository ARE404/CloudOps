<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { k8sApi, type K8sService, type K8sCluster } from '../../../api/k8s'

const list = ref<K8sService[]>([])
const clusters = ref<K8sCluster[]>([])
const selectedClusterId = ref<number | null>(null)
const namespace = ref('')
const loading = ref(false)

const columns = [
  { title: '名称', dataIndex: 'name' },
  { title: '命名空间', dataIndex: 'namespace' },
  { title: '类型', dataIndex: 'type' },
  { title: 'ClusterIP', dataIndex: 'clusterIP' },
  { title: '端口', dataIndex: 'ports' },
  { title: '创建时间', dataIndex: 'createdAt' },
]

async function fetchClusters() {
  const res = await k8sApi.listClusters()
  clusters.value = res || []
  if (clusters.value.length > 0) {
    selectedClusterId.value = clusters.value[0]!.id
    fetchList()
  }
}

async function fetchList() {
  if (!selectedClusterId.value) return
  loading.value = true
  try {
    const res = await k8sApi.listServices({ clusterId: selectedClusterId.value!, namespace: namespace.value || undefined })
    list.value = res || []
  } finally {
    loading.value = false
  }
}

onMounted(fetchClusters)
</script>

<template>
  <a-card>
    <div style="margin-bottom: 16px; display: flex; align-items: center; gap: 12px; flex-wrap: wrap;">
      <h3 style="margin: 0;">Service 列表</h3>
      <a-select v-model:value="selectedClusterId" style="width: 180px" placeholder="选择集群" @change="fetchList">
        <a-select-option v-for="c in clusters" :key="c.id" :value="c.id">{{ c.name }}</a-select-option>
      </a-select>
      <a-input v-model:value="namespace" placeholder="命名空间（可选）" style="width: 160px" @press-enter="fetchList" />
      <a-button @click="fetchList">查询</a-button>
    </div>
    <a-table :dataSource="list" :columns="columns" :loading="loading" row-key="name" :pagination="{ pageSize: 20 }" />
  </a-card>
</template>
