<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { k8sApi, type K8sDeployment, type K8sCluster } from '../../../api/k8s'

const list = ref<K8sDeployment[]>([])
const clusters = ref<K8sCluster[]>([])
const selectedClusterId = ref<number | null>(null)
const namespace = ref('')
const loading = ref(false)
const scaleModalVisible = ref(false)
const scalingRecord = ref<K8sDeployment | null>(null)
const scaleReplicas = ref(1)

const columns = [
  { title: '名称', dataIndex: 'name' },
  { title: '命名空间', dataIndex: 'namespace' },
  { title: '副本数', dataIndex: 'replicas' },
  { title: '就绪副本', dataIndex: 'readyReplicas' },
  { title: '镜像', dataIndex: 'image', ellipsis: true },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '操作', key: 'action', width: 180 },
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
    const res = await k8sApi.listDeployments({ clusterId: selectedClusterId.value!, namespace: namespace.value || undefined })
    list.value = res || []
  } finally {
    loading.value = false
  }
}

function openScale(record: K8sDeployment) {
  scalingRecord.value = record
  scaleReplicas.value = record.replicas
  scaleModalVisible.value = true
}

async function handleScale() {
  if (!selectedClusterId.value || !scalingRecord.value) return
  try {
    await k8sApi.scaleDeployment({
      clusterId: selectedClusterId.value,
      namespace: scalingRecord.value.namespace,
      name: scalingRecord.value.name,
      replicas: scaleReplicas.value,
    })
    message.success('扩缩容成功')
    scaleModalVisible.value = false
    fetchList()
  } catch (e) {
    message.error(String(e))
  }
}

async function handleDelete(record: K8sDeployment) {
  if (!selectedClusterId.value) return
  try {
    await k8sApi.deleteDeployment({ clusterId: selectedClusterId.value, namespace: record.namespace, name: record.name })
    message.success('删除成功')
    fetchList()
  } catch (e) {
    message.error(String(e))
  }
}

onMounted(fetchClusters)
</script>

<template>
  <a-card>
    <div style="margin-bottom: 16px; display: flex; align-items: center; gap: 12px; flex-wrap: wrap;">
      <h3 style="margin: 0;">Deployment 管理</h3>
      <a-select v-model:value="selectedClusterId" style="width: 180px" placeholder="选择集群" @change="fetchList">
        <a-select-option v-for="c in clusters" :key="c.id" :value="c.id">{{ c.name }}</a-select-option>
      </a-select>
      <a-input v-model:value="namespace" placeholder="命名空间（可选）" style="width: 160px" @press-enter="fetchList" />
      <a-button @click="fetchList">查询</a-button>
    </div>
    <a-table :dataSource="list" :columns="columns" :loading="loading" row-key="name" :pagination="{ pageSize: 20 }">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a-space>
            <a-button type="link" size="small" @click="openScale(record as K8sDeployment)">扩缩容</a-button>
            <a-popconfirm title="确认删除？" @confirm="handleDelete(record as K8sDeployment)">
              <a-button type="link" size="small" danger>删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>
  </a-card>

  <a-modal v-model:open="scaleModalVisible" title="扩缩容" @ok="handleScale">
    <a-form layout="vertical">
      <a-form-item label="副本数">
        <a-input-number v-model:value="scaleReplicas" :min="0" :max="100" style="width: 100%" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>
