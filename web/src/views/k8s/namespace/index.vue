<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { k8sApi, type K8sNamespace, type K8sCluster } from '../../../api/k8s'

const list = ref<K8sNamespace[]>([])
const clusters = ref<K8sCluster[]>([])
const selectedClusterId = ref<number | null>(null)
const loading = ref(false)
const modalVisible = ref(false)
const formRef = ref()

const form = reactive({ name: '' })

const columns = [
  { title: '名称', dataIndex: 'name' },
  { title: '状态', dataIndex: 'status' },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '操作', key: 'action', width: 100 },
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
    const res = await k8sApi.listNamespaces(selectedClusterId.value!)
    list.value = res || []
  } finally {
    loading.value = false
  }
}

async function handleDelete(name: string) {
  if (!selectedClusterId.value) return
  try {
    await k8sApi.deleteNamespace(selectedClusterId.value, name)
    message.success('删除成功')
    fetchList()
  } catch (e) {
    message.error(String(e))
  }
}

async function handleSubmit() {
  if (!selectedClusterId.value) return
  try {
    await formRef.value.validate()
    await k8sApi.createNamespace({ clusterId: selectedClusterId.value, name: form.name })
    message.success('创建成功')
    modalVisible.value = false
    fetchList()
  } catch (e) {
    if (typeof e === 'string') message.error(e)
  }
}

onMounted(fetchClusters)
</script>

<template>
  <a-card>
    <div style="margin-bottom: 16px; display: flex; justify-content: space-between; align-items: center;">
      <a-space>
        <h3 style="margin: 0;">命名空间</h3>
        <a-select v-model:value="selectedClusterId" style="width: 200px" placeholder="选择集群" @change="fetchList">
          <a-select-option v-for="c in clusters" :key="c.id" :value="c.id">{{ c.name }}</a-select-option>
        </a-select>
      </a-space>
      <a-button type="primary" @click="modalVisible = true">创建命名空间</a-button>
    </div>
    <a-table :dataSource="list" :columns="columns" :loading="loading" row-key="name" :pagination="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a-popconfirm title="确认删除？" @confirm="handleDelete((record as K8sNamespace).name)">
            <a-button type="link" size="small" danger>删除</a-button>
          </a-popconfirm>
        </template>
      </template>
    </a-table>
  </a-card>

  <a-modal v-model:open="modalVisible" title="创建命名空间" @ok="handleSubmit">
    <a-form ref="formRef" :model="form" layout="vertical">
      <a-form-item label="名称" name="name" :rules="[{ required: true }]">
        <a-input v-model:value="form.name" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>
