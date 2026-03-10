<script setup lang="ts">
import { computed, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useAppStore } from '../stores/app'
import {
  DashboardOutlined,
  UserOutlined,
  TeamOutlined,
  MenuOutlined,
  ClusterOutlined,
  AppstoreOutlined,
  RocketOutlined,
  ContainerOutlined,
  ApiOutlined,
  DatabaseOutlined,
  AlertOutlined,
  BellOutlined,
  SendOutlined,
  FileTextOutlined,
  OrderedListOutlined,
  BranchesOutlined,
  RobotOutlined,
  BookOutlined,
  PoweroffOutlined,
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const app = useAppStore()

const selectedKeys = computed(() => [route.path])
const openKeys = computed(() => {
  const path = route.path
  if (path.startsWith('/system')) return ['system']
  if (path.startsWith('/k8s')) return ['k8s']
  if (path.startsWith('/prometheus')) return ['prometheus']
  if (path.startsWith('/workorder')) return ['workorder']
  if (path.startsWith('/ai')) return ['ai']
  return []
})

function handleLogout() {
  auth.logout()
  router.push('/login')
}

const menuItems = [
  {
    key: '/dashboard',
    icon: () => h(DashboardOutlined),
    label: '仪表盘',
    onClick: () => router.push('/dashboard'),
  },
  {
    key: 'system',
    icon: () => h(UserOutlined),
    label: '系统管理',
    children: [
      { key: '/system/user', icon: () => h(UserOutlined), label: '用户管理', onClick: () => router.push('/system/user') },
      { key: '/system/role', icon: () => h(TeamOutlined), label: '角色管理', onClick: () => router.push('/system/role') },
      { key: '/system/menu', icon: () => h(MenuOutlined), label: '菜单管理', onClick: () => router.push('/system/menu') },
    ],
  },
  {
    key: 'k8s',
    icon: () => h(ClusterOutlined),
    label: 'Kubernetes',
    children: [
      { key: '/k8s/cluster', icon: () => h(ClusterOutlined), label: '集群管理', onClick: () => router.push('/k8s/cluster') },
      { key: '/k8s/namespace', icon: () => h(AppstoreOutlined), label: '命名空间', onClick: () => router.push('/k8s/namespace') },
      { key: '/k8s/deployment', icon: () => h(RocketOutlined), label: 'Deployment', onClick: () => router.push('/k8s/deployment') },
      { key: '/k8s/pod', icon: () => h(ContainerOutlined), label: 'Pod', onClick: () => router.push('/k8s/pod') },
      { key: '/k8s/service', icon: () => h(ApiOutlined), label: 'Service', onClick: () => router.push('/k8s/service') },
      { key: '/k8s/configmap', icon: () => h(DatabaseOutlined), label: 'ConfigMap', onClick: () => router.push('/k8s/configmap') },
    ],
  },
  {
    key: 'prometheus',
    icon: () => h(AlertOutlined),
    label: '监控告警',
    children: [
      { key: '/prometheus/scrape-pool', icon: () => h(DatabaseOutlined), label: '采集池', onClick: () => router.push('/prometheus/scrape-pool') },
      { key: '/prometheus/alert-rule', icon: () => h(AlertOutlined), label: '告警规则', onClick: () => router.push('/prometheus/alert-rule') },
      { key: '/prometheus/alert-event', icon: () => h(BellOutlined), label: '告警事件', onClick: () => router.push('/prometheus/alert-event') },
      { key: '/prometheus/send-group', icon: () => h(SendOutlined), label: '发送组', onClick: () => router.push('/prometheus/send-group') },
    ],
  },
  {
    key: 'workorder',
    icon: () => h(FileTextOutlined),
    label: '工单系统',
    children: [
      { key: '/workorder/template', icon: () => h(FileTextOutlined), label: '工单模板', onClick: () => router.push('/workorder/template') },
      { key: '/workorder/instance', icon: () => h(OrderedListOutlined), label: '工单实例', onClick: () => router.push('/workorder/instance') },
      { key: '/workorder/flow', icon: () => h(BranchesOutlined), label: '工单流转', onClick: () => router.push('/workorder/flow') },
    ],
  },
  {
    key: 'ai',
    icon: () => h(RobotOutlined),
    label: 'AI 助手',
    children: [
      { key: '/ai/chat', icon: () => h(RobotOutlined), label: 'AI 对话', onClick: () => router.push('/ai/chat') },
      { key: '/ai/knowledge', icon: () => h(BookOutlined), label: '知识库', onClick: () => router.push('/ai/knowledge') },
    ],
  },
]
</script>

<template>
  <a-layout style="min-height: 100vh">
    <a-layout-sider
      v-model:collapsed="app.collapsed"
      collapsible
      :width="220"
      style="background: #001529"
    >
      <div style="height: 64px; display: flex; align-items: center; justify-content: center; color: white; font-size: 18px; font-weight: bold; overflow: hidden;">
        <span v-if="!app.collapsed">CloudOps</span>
        <span v-else>CO</span>
      </div>
      <a-menu
        theme="dark"
        mode="inline"
        :selected-keys="selectedKeys"
        :open-keys="openKeys"
        :items="menuItems"
      />
    </a-layout-sider>
    <a-layout>
      <a-layout-header style="background: #fff; padding: 0 24px; display: flex; align-items: center; justify-content: space-between; box-shadow: 0 1px 4px rgba(0,21,41,.08);">
        <span style="font-size: 16px; font-weight: 500;">{{ route.meta.title || 'CloudOps' }}</span>
        <a-space>
          <a-button type="text" danger @click="handleLogout">
            <template #icon><PoweroffOutlined /></template>
            退出登录
          </a-button>
        </a-space>
      </a-layout-header>
      <a-layout-content style="margin: 24px; background: #f0f2f5; min-height: 280px;">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
