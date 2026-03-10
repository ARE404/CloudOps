<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { prometheusApi, type AlertRule } from '../../../api/prometheus'

const list = ref<AlertRule[]>([])
const total = ref(0)
const loading = ref(false)
const modalVisible = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref()
const params = reactive({ page: 1, pageSize: 10 })
const form = reactive({ name: '', expr: '', severity: 'warning', summary: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '规则名', dataIndex: 'name' },
  { title: '表达式', dataIndex: 'expr', ellipsis: true },
  { title: '严重级别', dataIndex: 'severity' },
  { title: '摘要', dataIndex: 'summary', ellipsis: true },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '操作', key: 'action', width: 150 },
]

async function fetchList() {
  loading.value = true
  try {
    const res = await prometheusApi.listAlertRules(params)
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
  Object.assign(form, { name: '', expr: '', severity: 'warning', summary: '' })
  modalVisible.value = true
}

function openEdit(record: AlertRule) {
  editingId.value = record.id
  Object.assign(form, record)
  modalVisible.value = true
}

async function handleDelete(id: number) {
  try {
    await prometheusApi.deleteAlertRule(id)
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
      await prometheusApi.updateAlertRule(editingId.value, form)
      message.success('更新成功')
    } else {
      await prometheusApi.createAlertRule(form)
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
      <h3 style="margin: 0;">告警规则</h3>
      <a-button type="primary" @click="openCreate">新建规则</a-button>
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
            <a-button type="link" size="small" @click="openEdit(record as AlertRule)">编辑</a-button>
            <a-popconfirm title="确认删除？" @confirm="handleDelete((record as AlertRule).id)">
              <a-button type="link" size="small" danger>删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>
  </a-card>

  <a-modal v-model:open="modalVisible" :title="editingId ? '编辑规则' : '新建规则'" @ok="handleSubmit" width="600px">
    <a-form ref="formRef" :model="form" layout="vertical">
      <a-form-item label="规则名" name="name" :rules="[{ required: true }]">
        <a-input v-model:value="form.name" />
      </a-form-item>
      <a-form-item label="PromQL 表达式" name="expr" :rules="[{ required: true }]">
        <a-textarea v-model:value="form.expr" :rows="3" />
      </a-form-item>
      <a-form-item label="严重级别" name="severity">
        <a-select v-model:value="form.severity">
          <a-select-option value="info">info</a-select-option>
          <a-select-option value="warning">warning</a-select-option>
          <a-select-option value="critical">critical</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="摘要" name="summary">
        <a-textarea v-model:value="form.summary" :rows="2" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>
