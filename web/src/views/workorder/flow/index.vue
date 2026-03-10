<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { workorderApi, type WorkorderFlow } from '../../../api/workorder'

const flows = ref<WorkorderFlow[]>([])
const loading = ref(false)
const actionVisible = ref(false)
const selectedInstanceId = ref<number>(0)
const instanceIdInput = ref<number>(1)
const actionForm = reactive({ action: 'approve', comment: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '操作', dataIndex: 'action' },
  { title: '备注', dataIndex: 'comment' },
  { title: '操作人', dataIndex: 'operatorName' },
  { title: '时间', dataIndex: 'createdAt' },
]

async function fetchFlows() {
  if (!instanceIdInput.value) return
  loading.value = true
  try {
    flows.value = await workorderApi.listFlows(instanceIdInput.value)
  } catch (e) {
    message.error(String(e))
  } finally {
    loading.value = false
  }
}

function openAction(instanceId: number) {
  selectedInstanceId.value = instanceId
  Object.assign(actionForm, { action: 'approve', comment: '' })
  actionVisible.value = true
}

async function handleAction() {
  try {
    await workorderApi.doAction(selectedInstanceId.value, actionForm)
    message.success('操作成功')
    actionVisible.value = false
    fetchFlows()
  } catch (e) {
    message.error(String(e))
  }
}

onMounted(fetchFlows)
</script>

<template>
  <a-card>
    <div style="margin-bottom: 16px; display: flex; align-items: center; gap: 16px;">
      <h3 style="margin: 0;">工单流转</h3>
      <a-input-number v-model:value="instanceIdInput" placeholder="工单实例ID" style="width: 160px" />
      <a-button type="primary" @click="fetchFlows">查询</a-button>
      <a-button @click="openAction(instanceIdInput)">执行操作</a-button>
    </div>
    <a-table
      :dataSource="flows"
      :columns="columns"
      :loading="loading"
      row-key="id"
    />
  </a-card>

  <a-modal v-model:open="actionVisible" title="工单操作" @ok="handleAction">
    <a-form :model="actionForm" layout="vertical">
      <a-form-item label="操作" name="action">
        <a-select v-model:value="actionForm.action">
          <a-select-option value="approve">审批通过</a-select-option>
          <a-select-option value="reject">驳回</a-select-option>
          <a-select-option value="close">关闭</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="备注" name="comment">
        <a-textarea v-model:value="actionForm.comment" :rows="3" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>
