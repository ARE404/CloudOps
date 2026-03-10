<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { workorderApi, type WorkorderTemplate } from '../../../api/workorder'

const list = ref<WorkorderTemplate[]>([])
const total = ref(0)
const loading = ref(false)
const modalVisible = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref()
const params = reactive({ page: 1, pageSize: 10 })
const form = reactive({ name: '', description: '', steps: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '模板名', dataIndex: 'name' },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '操作', key: 'action', width: 150 },
]

async function fetchList() {
  loading.value = true
  try {
    const res = await workorderApi.listTemplates(params)
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
  Object.assign(form, { name: '', description: '', steps: '' })
  modalVisible.value = true
}

function openEdit(record: WorkorderTemplate) {
  editingId.value = record.id
  Object.assign(form, record)
  modalVisible.value = true
}

async function handleDelete(id: number) {
  try {
    await workorderApi.deleteTemplate(id)
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
      await workorderApi.updateTemplate(editingId.value, form)
      message.success('更新成功')
    } else {
      await workorderApi.createTemplate(form)
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
      <h3 style="margin: 0;">工单模板</h3>
      <a-button type="primary" @click="openCreate">新建模板</a-button>
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
            <a-button type="link" size="small" @click="openEdit(record as WorkorderTemplate)">编辑</a-button>
            <a-popconfirm title="确认删除？" @confirm="handleDelete((record as WorkorderTemplate).id)">
              <a-button type="link" size="small" danger>删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>
  </a-card>

  <a-modal v-model:open="modalVisible" :title="editingId ? '编辑模板' : '新建模板'" @ok="handleSubmit" width="600px">
    <a-form ref="formRef" :model="form" layout="vertical">
      <a-form-item label="模板名" name="name" :rules="[{ required: true }]">
        <a-input v-model:value="form.name" />
      </a-form-item>
      <a-form-item label="描述" name="description">
        <a-textarea v-model:value="form.description" :rows="2" />
      </a-form-item>
      <a-form-item label="步骤（JSON）" name="steps">
        <a-textarea v-model:value="form.steps" :rows="5" placeholder='[{"name":"审批","assignee":"admin"}]' />
      </a-form-item>
    </a-form>
  </a-modal>
</template>
