<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { message } from 'ant-design-vue'
import { SendOutlined } from '@ant-design/icons-vue'
import { aiApi, type ChatMessage } from '../../../api/ai'

const messages = ref<ChatMessage[]>([])
const inputText = ref('')
const loading = ref(false)
const chatContainerRef = ref<HTMLElement>()

async function sendMessage() {
  const text = inputText.value.trim()
  if (!text || loading.value) return

  messages.value.push({ role: 'user', content: text })
  inputText.value = ''
  loading.value = true

  await nextTick()
  scrollToBottom()

  try {
    const res = await aiApi.chat(messages.value)
    messages.value.push({ role: 'assistant', content: res.content })
    await nextTick()
    scrollToBottom()
  } catch (e) {
    message.error(String(e))
    messages.value.pop()
  } finally {
    loading.value = false
  }
}

function scrollToBottom() {
  if (chatContainerRef.value) {
    chatContainerRef.value.scrollTop = chatContainerRef.value.scrollHeight
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    sendMessage()
  }
}
</script>

<template>
  <a-card style="height: calc(100vh - 160px); display: flex; flex-direction: column;">
    <template #title>AI 对话助手</template>
    <div
      ref="chatContainerRef"
      style="flex: 1; overflow-y: auto; padding: 16px; background: #f5f5f5; border-radius: 8px; margin-bottom: 16px; min-height: 400px; max-height: calc(100vh - 320px);"
    >
      <div v-if="messages.length === 0" style="text-align: center; color: #999; margin-top: 80px;">
        <p style="font-size: 16px;">你好！我是 CloudOps AI 助手</p>
        <p>可以向我询问 K8s、监控告警、运维操作相关的问题</p>
      </div>
      <div v-for="(msg, i) in messages" :key="i" style="margin-bottom: 16px;">
        <div v-if="msg.role === 'user'" style="display: flex; justify-content: flex-end;">
          <div style="background: #1677ff; color: white; padding: 8px 16px; border-radius: 18px 18px 4px 18px; max-width: 70%; white-space: pre-wrap; word-break: break-word;">
            {{ msg.content }}
          </div>
        </div>
        <div v-else style="display: flex; justify-content: flex-start;">
          <div style="background: white; padding: 8px 16px; border-radius: 18px 18px 18px 4px; max-width: 70%; white-space: pre-wrap; word-break: break-word; box-shadow: 0 1px 3px rgba(0,0,0,0.1);">
            {{ msg.content }}
          </div>
        </div>
      </div>
      <div v-if="loading" style="display: flex; justify-content: flex-start;">
        <div style="background: white; padding: 8px 16px; border-radius: 18px; box-shadow: 0 1px 3px rgba(0,0,0,0.1);">
          <a-spin size="small" /> 思考中...
        </div>
      </div>
    </div>
    <div style="display: flex; gap: 8px;">
      <a-textarea
        v-model:value="inputText"
        placeholder="输入消息，Enter 发送，Shift+Enter 换行"
        :rows="3"
        style="flex: 1; resize: none;"
        @keydown="handleKeydown"
      />
      <a-button
        type="primary"
        :loading="loading"
        style="height: auto; padding: 8px 20px;"
        @click="sendMessage"
      >
        <template #icon><SendOutlined /></template>
        发送
      </a-button>
    </div>
  </a-card>
</template>
