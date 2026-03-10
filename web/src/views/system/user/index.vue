<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { userApi, type User, type CreateUserParams } from '../../../api/user'

const list = ref<User[]>([])
const total = ref(0)
const loading = ref(false)
const modalVisible = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref()

const params = reactive({ page: 1, pageSize: 10 })
const form = reactive<CreateUserParams>({ username: '', password: '', realName: '', email: '', phone: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '用户名', dataIndex: 'username' },
  { title: '姓名', dataIndex: 'realName' },
  { title: '邮箱', dataIndex: 'email' },
  { title: '手机', dataIndex: 'phone' },
  { title: '状态', dataIndex: 'status', customRender: ({ value }: { value: number }) => value === 1 ? '正常' : '禁用' },
  { title: '创建时间', dataIndex: 'createdAt' },
  { title: '操作', key: 'action', width: 150 },
]

async function fetchList() {
  loading.value = true
  try {
    const res = await userApi.list(params)
    list.value = res.list || []
    total.value = res.total || 0
  } catch (e) {
    message.error(String(e))
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  Object.assign(form, { username: '', password: '', realName: '', email: '', phone: '' })
  modalVisible.value = true
}

function openEdit(record: User) {
  editingId.value = record.id
  Object.assign(form, record)
  modalVisible.value = true
}

async function handleDelete(id: number) {
  try {
    await userApi.delete(id)
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
      await userApi.update(editingId.value, form)
      message.success('更新成功')
    } else {
      await userApi.create(form)
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
      <h3 style="margin: 0;">用户管理</h3>
      <a-button type="primary" @click="openCreate">新建用户</a-button>
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
            <a-button type="link" size="small" @click="openEdit(record as User)">编辑</a-button>
            <a-popconfirm title="确认删除？" @confirm="handleDelete((record as User).id)">
              <a-button type="link" size="small" danger>删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>
  </a-card>

  <a-modal v-model:open="modalVisible" :title="editingId ? '编辑用户' : '新建用户'" @ok="handleSubmit">
    <a-form ref="formRef" :model="form" layout="vertical">
      <a-form-item label="用户名" name="username" :rules="[{ required: true }]">
        <a-input v-model:value="form.username" />
      </a-form-item>
      <a-form-item v-if="!editingId" label="密码" name="password" :rules="[{ required: true }]">
        <a-input-password v-model:value="form.password" />
      </a-form-item>
      <a-form-item label="姓名" name="realName">
        <a-input v-model:value="form.realName" />
      </a-form-item>
      <a-form-item label="邮箱" name="email">
        <a-input v-model:value="form.email" />
      </a-form-item>
      <a-form-item label="手机" name="phone">
        <a-input v-model:value="form.phone" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>
