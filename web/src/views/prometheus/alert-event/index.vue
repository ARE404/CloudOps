<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { prometheusApi, type AlertEvent } from '../../../api/prometheus'

const list = ref<AlertEvent[]>([])
const total = ref(0)
const loading = ref(false)
const params = reactive({ page: 1, pageSize: 10 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '告警名', dataIndex: 'alertName' },
  { title: '级别', dataIndex: 'severity' },
  { title: '状态', dataIndex: 'status' },
  { title: '摘要', dataIndex: 'summary', ellipsis: true },
  { title: '开始时间', dataIndex: 'startAt' },
]

async function fetchList() {
  loading.value = true
  try {
    const res = await prometheusApi.listAlertEvents(params)
    list.value = (res as any).list || []
    total.value = (res as any).total || 0
  } catch (e) {
    message.error(String(e))
  } finally {
    loading.value = false
  }
}

onMounted(fetchList)
</script>

<template>
  <a-card>
    <div style="margin-bottom: 16px;">
      <h3 style="margin: 0;">告警事件</h3>
    </div>
    <a-table
      :dataSource="list"
      :columns="columns"
      :loading="loading"
      :pagination="{ current: params.page, pageSize: params.pageSize, total, onChange: (p: number) => { params.page = p; fetchList() } }"
      row-key="id"
    />
  </a-card>
</template>
