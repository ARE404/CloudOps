import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/login/index.vue'),
    },
    {
      path: '/',
      component: () => import('../layouts/BasicLayout.vue'),
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('../views/dashboard/index.vue'),
          meta: { title: '仪表盘' },
        },
        // System
        { path: 'system/user', component: () => import('../views/system/user/index.vue'), meta: { title: '用户管理' } },
        { path: 'system/role', component: () => import('../views/system/role/index.vue'), meta: { title: '角色管理' } },
        { path: 'system/menu', component: () => import('../views/system/menu/index.vue'), meta: { title: '菜单管理' } },
        // K8s
        { path: 'k8s/cluster', component: () => import('../views/k8s/cluster/index.vue'), meta: { title: 'K8s集群' } },
        { path: 'k8s/namespace', component: () => import('../views/k8s/namespace/index.vue'), meta: { title: '命名空间' } },
        { path: 'k8s/deployment', component: () => import('../views/k8s/deployment/index.vue'), meta: { title: 'Deployment' } },
        { path: 'k8s/pod', component: () => import('../views/k8s/pod/index.vue'), meta: { title: 'Pod' } },
        { path: 'k8s/service', component: () => import('../views/k8s/service/index.vue'), meta: { title: 'Service' } },
        { path: 'k8s/configmap', component: () => import('../views/k8s/configmap/index.vue'), meta: { title: 'ConfigMap' } },
        // Prometheus
        { path: 'prometheus/scrape-pool', component: () => import('../views/prometheus/scrape-pool/index.vue'), meta: { title: '采集池' } },
        { path: 'prometheus/alert-rule', component: () => import('../views/prometheus/alert-rule/index.vue'), meta: { title: '告警规则' } },
        { path: 'prometheus/alert-event', component: () => import('../views/prometheus/alert-event/index.vue'), meta: { title: '告警事件' } },
        { path: 'prometheus/send-group', component: () => import('../views/prometheus/send-group/index.vue'), meta: { title: '发送组' } },
        // Workorder
        { path: 'workorder/template', component: () => import('../views/workorder/template/index.vue'), meta: { title: '工单模板' } },
        { path: 'workorder/instance', component: () => import('../views/workorder/instance/index.vue'), meta: { title: '工单实例' } },
        { path: 'workorder/flow', component: () => import('../views/workorder/flow/index.vue'), meta: { title: '工单流转' } },
        // AI
        { path: 'ai/chat', component: () => import('../views/ai/chat/index.vue'), meta: { title: 'AI 对话' } },
        { path: 'ai/knowledge', component: () => import('../views/ai/knowledge/index.vue'), meta: { title: '知识库' } },
      ],
    },
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
})

router.beforeEach(to => {
  const auth = useAuthStore()
  if (to.path !== '/login' && !auth.token) return '/login'
})

export default router
