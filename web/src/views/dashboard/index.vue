<script setup lang="ts">
</script>

<template>
  <div>
    <a-card style="border-radius: 8px;">
      <a-result
        status="success"
        title="欢迎使用 CloudOps"
        sub-title="云原生运维管理平台，支持 K8s 多集群管理、Prometheus 监控告警、工单系统和 AI 助手"
      >
        <template #extra>
          <a-space>
            <a-button type="primary" @click="$router.push('/k8s/cluster')">管理集群</a-button>
            <a-button @click="$router.push('/prometheus/alert-event')">查看告警</a-button>
            <a-button @click="$router.push('/ai/chat')">AI 助手</a-button>
          </a-space>
        </template>
      </a-result>
    </a-card>

    <a-row :gutter="16" style="margin-top: 24px;">
      <a-col :span="6">
        <a-card>
          <a-statistic title="K8s 集群" :value="0" suffix="个" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="告警事件" :value="0" suffix="条" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="工单实例" :value="0" suffix="个" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="系统用户" :value="0" suffix="人" />
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>
