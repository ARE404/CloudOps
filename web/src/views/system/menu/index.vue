<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { menuApi, type Menu } from '../../../api/menu'

const list = ref<Menu[]>([])
const loading = ref(false)
const modalVisible = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref()

const form = reactive({ parentId: 0, name: '', path: '', icon: '', sort: 0, type: 1 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '菜单名', dataIndex: 'name' },
  { title: '路径', dataIndex: 'path' },
  { title: '图标', dataIndex: 'icon' },
  { title: '排序', dataIndex: 'sort', width: 80 },
  { title: '类型', dataIndex: 'type', customRender: ({ value }: { value: number }) => value === 1 ? '目录' : value === 2 ? '菜单' : '按钮' },
  { title: '操作', key: 'action', width: 150 },
]

async function fetchList() {
  loading.value = true
  try {
    const res = await menuApi.list()
    list.value = res || []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  Object.assign(form, { parentId: 0, name: '', path: '', icon: '', sort: 0, type: 1 })
  modalVisible.value = true
}

function openEdit(record: Menu) {
  editingId.value = record.id
  Object.assign(form, record)
  modalVisible.value = true
}

async function handleDelete(id: number) {
  try {
    await menuApi.delete(id)
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
      await menuApi.update(editingId.value, form)
      message.success('更新成功')
    } else {
      await menuApi.create(form)
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
      <h3 style="margin: 0;">菜单管理</h3>
      <a-button type="primary" @click="openCreate">新建菜单</a-button>
    </div>
    <a-table
      :dataSource="list"
      :columns="columns"
      :loading="loading"
      row-key="id"
      :pagination="false"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a-space>
            <a-button type="link" size="small" @click="openEdit(record as Menu)">编辑</a-button>
            <a-popconfirm title="确认删除？" @confirm="handleDelete((record as Menu).id)">
              <a-button type="link" size="small" danger>删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>
  </a-card>

  <a-modal v-model:open="modalVisible" :title="editingId ? '编辑菜单' : '新建菜单'" @ok="handleSubmit">
    <a-form ref="formRef" :model="form" layout="vertical">
      <a-form-item label="菜单名" name="name" :rules="[{ required: true }]">
        <a-input v-model:value="form.name" />
      </a-form-item>
      <a-form-item label="路径" name="path">
        <a-input v-model:value="form.path" />
      </a-form-item>
      <a-form-item label="图标" name="icon">
        <a-input v-model:value="form.icon" />
      </a-form-item>
      <a-form-item label="排序" name="sort">
        <a-input-number v-model:value="form.sort" style="width: 100%" />
      </a-form-item>
      <a-form-item label="类型" name="type">
        <a-select v-model:value="form.type">
          <a-select-option :value="1">目录</a-select-option>
          <a-select-option :value="2">菜单</a-select-option>
          <a-select-option :value="3">按钮</a-select-option>
        </a-select>
      </a-form-item>
    </a-form>
  </a-modal>
</template>
