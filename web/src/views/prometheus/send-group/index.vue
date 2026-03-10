<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { prometheusApi, type SendGroup } from '../../../api/prometheus'

const list = ref<SendGroup[]>([])
const total = ref(0)
const loading = ref(false)
const modalVisible = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref()
const params = reactive({ page: 1, pageSize: 10 })
const form = reactive({ name: '', type: 'webhook', webhook: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '名称', dataIndex: 'name' },
  { title: '类型', dataIndex: 'type' },
  { title: 'Webhook', dataIndex: 'webhook', ellipsis: true },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '操作', key: 'action', width: 150 },
]

async function fetchList() {
  loading.value = true
  try {
    const res = await prometheusApi.listSendGroups(params)
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
  Object.assign(form, { name: '', type: 'webhook', webhook: '' })
  modalVisible.value = true
}

function openEdit(record: SendGroup) {
  editingId.value = record.id
  Object.assign(form, record)
  modalVisible.value = true
}

async function handleDelete(id: number) {
  try {
    await prometheusApi.deleteSendGroup(id)
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
      await prometheusApi.updateSendGroup(editingId.value, form)
      message.success('更新成功')
    } else {
      await prometheusApi.createSendGroup(form)
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
      <h3 style="margin: 0;">发送组</h3>
      <a-button type="primary" @click="openCreate">新建发送组</a-button>
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
            <a-button type="link" size="small" @click="openEdit(record as SendGroup)">编辑</a-button>
            <a-popconfirm title="确认删除？" @confirm="handleDelete((record as SendGroup).id)">
              <a-button type="link" size="small" danger>删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>
  </a-card>

  <a-modal v-model:open="modalVisible" :title="editingId ? '编辑发送组' : '新建发送组'" @ok="handleSubmit">
    <a-form ref="formRef" :model="form" layout="vertical">
      <a-form-item label="名称" name="name" :rules="[{ required: true }]">
        <a-input v-model:value="form.name" />
      </a-form-item>
      <a-form-item label="类型" name="type">
        <a-select v-model:value="form.type">
          <a-select-option value="webhook">Webhook</a-select-option>
          <a-select-option value="email">Email</a-select-option>
          <a-select-option value="feishu">飞书</a-select-option>
          <a-select-option value="dingtalk">钉钉</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="Webhook URL" name="webhook">
        <a-input v-model:value="form.webhook" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>
