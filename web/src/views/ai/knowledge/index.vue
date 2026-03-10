<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { aiApi, type KnowledgeItem } from '../../../api/ai'

const list = ref<KnowledgeItem[]>([])
const loading = ref(false)
const uploading = ref(false)
const fileInputRef = ref<HTMLInputElement>()

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '文件名', dataIndex: 'name' },
  { title: '大小', dataIndex: 'size', customRender: ({ value }: { value: number }) => `${(value / 1024).toFixed(1)} KB` },
  { title: '上传时间', dataIndex: 'createdAt' },
  { title: '操作', key: 'action', width: 100 },
]

async function fetchList() {
  loading.value = true
  try {
    list.value = await aiApi.listKnowledge()
  } catch (e) {
    message.error(String(e))
  } finally {
    loading.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await aiApi.deleteKnowledge(id)
    message.success('删除成功')
    fetchList()
  } catch (e) {
    message.error(String(e))
  }
}

async function handleUpload(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('file', file)
    await aiApi.uploadKnowledge(formData)
    message.success('上传成功')
    fetchList()
  } catch (err) {
    message.error(String(err))
  } finally {
    uploading.value = false
    if (fileInputRef.value) fileInputRef.value.value = ''
  }
}

onMounted(fetchList)
</script>

<template>
  <a-card>
    <div style="margin-bottom: 16px; display: flex; justify-content: space-between; align-items: center;">
      <h3 style="margin: 0;">知识库管理</h3>
      <div>
        <input ref="fileInputRef" type="file" style="display: none" @change="handleUpload" />
        <a-button type="primary" :loading="uploading" @click="fileInputRef?.click()">上传文件</a-button>
      </div>
    </div>
    <a-table
      :dataSource="list"
      :columns="columns"
      :loading="loading"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a-popconfirm title="确认删除？" @confirm="handleDelete((record as KnowledgeItem).id)">
            <a-button type="link" size="small" danger>删除</a-button>
          </a-popconfirm>
        </template>
      </template>
    </a-table>
  </a-card>
</template>
