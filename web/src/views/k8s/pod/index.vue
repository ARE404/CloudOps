<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { k8sApi, type K8sPod, type K8sCluster } from '../../../api/k8s'

const list = ref<K8sPod[]>([])
const clusters = ref<K8sCluster[]>([])
const selectedClusterId = ref<number | null>(null)
const namespace = ref('')
const loading = ref(false)
const logsModalVisible = ref(false)
const logsContent = ref('')
const logsTitle = ref('')

const columns = [
  { title: '名称', dataIndex: 'name', ellipsis: true },
  { title: '命名空间', dataIndex: 'namespace' },
  { title: '状态', dataIndex: 'status' },
  { title: '节点', dataIndex: 'nodeName' },
  { title: 'IP', dataIndex: 'ip' },
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
    const res = await k8sApi.listPods({ clusterId: selectedClusterId.value!, namespace: namespace.value || undefined })
    list.value = res || []
  } finally {
    loading.value = false
  }
}

async function viewLogs(record: K8sPod) {
  if (!selectedClusterId.value) return
  try {
    const logs = await k8sApi.getPodLogs({ clusterId: selectedClusterId.value, namespace: record.namespace, name: record.name })
    logsContent.value = logs
    logsTitle.value = `${record.name} 日志`
    logsModalVisible.value = true
  } catch (e) {
    message.error(String(e))
  }
}

async function handleDelete(record: K8sPod) {
  if (!selectedClusterId.value) return
  try {
    await k8sApi.deletePod({ clusterId: selectedClusterId.value, namespace: record.namespace, name: record.name })
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
      <h3 style="margin: 0;">Pod 管理</h3>
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
            <a-button type="link" size="small" @click="viewLogs(record as K8sPod)">查看日志</a-button>
            <a-popconfirm title="确认删除？" @confirm="handleDelete(record as K8sPod)">
              <a-button type="link" size="small" danger>删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>
  </a-card>

  <a-modal v-model:open="logsModalVisible" :title="logsTitle" width="900px" :footer="null">
    <pre style="background: #1e1e1e; color: #d4d4d4; padding: 16px; border-radius: 6px; max-height: 500px; overflow: auto; font-size: 12px; white-space: pre-wrap; word-break: break-all;">{{ logsContent }}</pre>
  </a-modal>
</template>
