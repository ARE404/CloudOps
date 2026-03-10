<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { prometheusApi, type ScrapePool } from '../../../api/prometheus'

const list = ref<ScrapePool[]>([])
const total = ref(0)
const loading = ref(false)
const modalVisible = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref()
const params = reactive({ page: 1, pageSize: 10 })
const form = reactive({ name: '', clusterId: 0, labels: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '名称', dataIndex: 'name' },
  { title: '集群ID', dataIndex: 'clusterId' },
  { title: '标签', dataIndex: 'labels' },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '操作', key: 'action', width: 150 },
]

async function fetchList() {
  loading.value = true
  try {
    const res = await prometheusApi.listScrapePools(params)
    list.value = (res as any).list || []
    total.value = (res as any).total || 0
  } catch (e) {
    message.error(String(e))
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  Object.assign(form, { name: '', clusterId: 0, labels: '' })
  modalVisible.value = true
}

function openEdit(record: ScrapePool) {
  editingId.value = record.id
  Object.assign(form, record)
  modalVisible.value = true
}

async function handleDelete(id: number) {
  try {
    await prometheusApi.deleteScrapePool(id)
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
      await prometheusApi.updateScrapePool(editingId.value, form)
      message.success('更新成功')
    } else {
      await prometheusApi.createScrapePool(form)
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
      <h3 style="margin: 0;">采集池管理</h3>
      <a-button type="primary" @click="openCreate">新建采集池</a-button>
    </div>
    <a-table
      :dataSource="list"
      :columns="columns"
      :loading="loading"
      :pagination="{ current: params.page, pageSize: params.pageSize, total, onChange: (p: number) => { params.page = p; fetchList() } }"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a-space>
            <a-button type="link" size="small" @click="openEdit(record as ScrapePool)">编辑</a-button>
            <a-popconfirm title="确认删除？" @confirm="handleDelete((record as ScrapePool).id)">
              <a-button type="link" size="small" danger>删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>
  </a-card>

  <a-modal v-model:open="modalVisible" :title="editingId ? '编辑采集池' : '新建采集池'" @ok="handleSubmit">
    <a-form ref="formRef" :model="form" layout="vertical">
      <a-form-item label="名称" name="name" :rules="[{ required: true }]">
        <a-input v-model:value="form.name" />
      </a-form-item>
      <a-form-item label="集群ID" name="clusterId">
        <a-input-number v-model:value="form.clusterId" style="width: 100%" />
      </a-form-item>
      <a-form-item label="标签" name="labels">
        <a-textarea v-model:value="form.labels" :rows="3" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>
