<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { workorderApi, type WorkorderInstance } from '../../../api/workorder'

const list = ref<WorkorderInstance[]>([])
const total = ref(0)
const loading = ref(false)
const modalVisible = ref(false)
const formRef = ref()
const params = reactive({ page: 1, pageSize: 10 })
const form = reactive({ templateId: undefined as number | undefined, title: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '标题', dataIndex: 'title' },
  { title: '模板', dataIndex: 'templateName' },
  { title: '状态', dataIndex: 'status' },
  { title: '创建人', dataIndex: 'creatorName' },
  { title: '创建时间', dataIndex: 'createdAt' },
]

async function fetchList() {
  loading.value = true
  try {
    const res = await workorderApi.listInstances(params)
    list.value = (res as any).list || []
    total.value = (res as any).total || 0
  } catch (e) {
    message.error(String(e))
  } finally {
    loading.value = false
  }
}

function openCreate() {
  Object.assign(form, { templateId: undefined, title: '' })
  modalVisible.value = true
}

async function handleSubmit() {
  try {
    await formRef.value.validate()
    await workorderApi.createInstance(form)
    message.success('创建成功')
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
      <h3 style="margin: 0;">工单实例</h3>
      <a-button type="primary" @click="openCreate">发起工单</a-button>
    </div>
    <a-table
      :dataSource="list"
      :columns="columns"
      :loading="loading"
      :pagination="{ current: params.page, pageSize: params.pageSize, total, onChange: (p: number) => { params.page = p; fetchList() } }"
      row-key="id"
    />
  </a-card>

  <a-modal v-model:open="modalVisible" title="发起工单" @ok="handleSubmit">
    <a-form ref="formRef" :model="form" layout="vertical">
      <a-form-item label="标题" name="title" :rules="[{ required: true }]">
        <a-input v-model:value="form.title" />
      </a-form-item>
      <a-form-item label="模板ID" name="templateId" :rules="[{ required: true }]">
        <a-input-number v-model:value="form.templateId" style="width: 100%" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>
