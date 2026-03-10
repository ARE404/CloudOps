<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { k8sApi, type K8sCluster } from '../../../api/k8s'

const list = ref<K8sCluster[]>([])
const loading = ref(false)
const modalVisible = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref()

const form = reactive({ name: '', kubeConfig: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '集群名', dataIndex: 'name' },
  { title: '状态', dataIndex: 'status', customRender: ({ value }: { value: number }) => value === 1 ? '正常' : '异常' },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '操作', key: 'action', width: 150 },
]

async function fetchList() {
  loading.value = true
  try {
    const res = await k8sApi.listClusters()
    list.value = res || []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  Object.assign(form, { name: '', kubeConfig: '' })
  modalVisible.value = true
}

function openEdit(record: K8sCluster) {
  editingId.value = record.id
  Object.assign(form, record)
  modalVisible.value = true
}

async function handleDelete(id: number) {
  try {
    await k8sApi.deleteCluster(id)
    message.success('删除成功')
    fetchList()
  } catch (e) {
    message.error(String(e))
  }
}

async function handleSubmit() {
  try {
    await formRef.value.validate()
    if (editingId.value) {
      await k8sApi.updateCluster(editingId.value, form)
      message.success('更新成功')
    } else {
      await k8sApi.createCluster(form)
      message.success('创建成功')
    }
    modalVisible.value = false
    fetchList()
  } catch (e) {
    if (typeof e === 'string') message.error(e)
  }
}

onMounted(fetchList)
</script>

<template>
  <a-card>
    <div style="margin-bottom: 16px; display: flex; justify-content: space-between; align-items: center;">
      <h3 style="margin: 0;">K8s 集群管理</h3>
      <a-button type="primary" @click="openCreate">添加集群</a-button>
    </div>
    <a-table :dataSource="list" :columns="columns" :loading="loading" row-key="id" :pagination="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a-space>
            <a-button type="link" size="small" @click="openEdit(record as K8sCluster)">编辑</a-button>
            <a-popconfirm title="确认删除？" @confirm="handleDelete((record as K8sCluster).id)">
              <a-button type="link" size="small" danger>删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>
  </a-card>

  <a-modal v-model:open="modalVisible" :title="editingId ? '编辑集群' : '添加集群'" width="700px" @ok="handleSubmit">
    <a-form ref="formRef" :model="form" layout="vertical">
      <a-form-item label="集群名" name="name" :rules="[{ required: true }]">
        <a-input v-model:value="form.name" />
      </a-form-item>
      <a-form-item label="KubeConfig" name="kubeConfig" :rules="[{ required: true }]">
        <a-textarea v-model:value="form.kubeConfig" :rows="12" placeholder="粘贴 kubeconfig 内容" style="font-family: monospace; font-size: 12px;" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>
