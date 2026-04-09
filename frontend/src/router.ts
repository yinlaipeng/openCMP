import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/layout/index.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/',
    component: Layout,
    redirect: '/cloud-accounts',
    children: [
      {
        path: '/cloud-accounts',
        name: 'CloudAccounts',
        component: () => import('@/views/cloud-accounts/index.vue'),
        meta: { title: '多云管理', icon: 'Cloudy' }
      },
      {
        path: '/cloud-management',
        name: 'CloudManagement',
        meta: { title: '多云管理', icon: 'Cloudy' },
        children: [
          {
            path: 'sync-policies',
            name: 'SyncPolicies',
            component: () => import('@/views/cloud-management/sync-policies/index.vue'),
            meta: { title: '同步策略' }
          }
        ]
      },
      {
        path: '/scheduled-tasks',
        name: 'ScheduledTasks',
        component: () => import('@/views/cloud-accounts/scheduled-tasks.vue'),
        meta: { title: '定时同步任务', icon: 'Timer' }
      },
      {
        path: '/compute',
        name: 'Compute',
        meta: { title: '计算资源', icon: 'Cpu' },
        children: [
          {
            path: 'vms',
            name: 'VMs',
            component: () => import('@/views/compute/vms/index.vue'),
            meta: { title: '虚拟机管理' }
          },
          {
            path: 'images',
            name: 'Images',
            component: () => import('@/views/compute/images/index.vue'),
            meta: { title: '镜像管理' }
          }
        ]
      },
      {
        path: '/network',
        name: 'Network',
        meta: { title: '网络资源', icon: 'Connection' },
        children: [
          {
            path: 'vpcs',
            name: 'VPCs',
            component: () => import('@/views/network/vpcs/index.vue'),
            meta: { title: 'VPC 管理' }
          },
          {
            path: 'subnets',
            name: 'Subnets',
            component: () => import('@/views/network/subnets/index.vue'),
            meta: { title: '子网管理' }
          },
          {
            path: 'security-groups',
            name: 'SecurityGroups',
            component: () => import('@/views/network/security-groups/index.vue'),
            meta: { title: '安全组管理' }
          },
          {
            path: 'eips',
            name: 'EIPs',
            component: () => import('@/views/network/eips/index.vue'),
            meta: { title: '弹性 IP' }
          }
        ]
      },
      {
        path: '/iam',
        name: 'IAM',
        meta: { title: '认证与安全', icon: 'Lock' },
        children: [
          {
            path: 'auth-sources',
            name: 'AuthSources',
            component: () => import('@/views/iam/auth-sources/index.vue'),
            meta: { title: '认证源' }
          },
          {
            path: 'domains',
            name: 'Domains',
            component: () => import('@/views/iam/domains/index.vue'),
            meta: { title: '域' }
          },
          {
            path: 'projects',
            name: 'Projects',
            component: () => import('@/views/iam/projects/index.vue'),
            meta: { title: '项目' }
          },
          {
            path: 'users',
            name: 'Users',
            component: () => import('@/views/iam/users/index.vue'),
            meta: { title: '用户' }
          },
          {
            path: 'groups',
            name: 'Groups',
            component: () => import('@/views/iam/groups/index.vue'),
            meta: { title: '用户组' }
          },
          {
            path: 'roles',
            name: 'Roles',
            component: () => import('@/views/iam/roles/index.vue'),
            meta: { title: '角色' }
          },
          {
            path: 'messages',
            name: 'Messages',
            component: () => import('@/views/iam/messages/index.vue'),
            meta: { title: '消息中心' }
          },
          {
            path: 'alerts',
            name: 'Alerts',
            component: () => import('@/views/iam/alerts/index.vue'),
            meta: { title: '安全告警' }
          },
          {
            path: 'project-alerts',
            name: 'ProjectAlerts',
            component: () => import('@/views/iam/project-alerts/index.vue'),
            meta: { title: '项目安全告警' }
          },
          {
            path: 'permissions',
            name: 'Permissions',
            component: () => import('@/views/iam/permissions/index.vue'),
            meta: { title: '权限' }
          }
        ]
      },
      {
        path: '/message-center',
        name: 'MessageCenter',
        meta: { title: '消息中心', icon: 'Bell' },
        children: [
          {
            path: 'inbox',
            name: 'Inbox',
            component: () => import('@/views/message-center/inbox/index.vue'),
            meta: { title: '站内信' }
          },
          {
            path: 'project-inbox',
            name: 'ProjectInbox',
            component: () => import('@/views/message-center/project-inbox/index.vue'),
            meta: { title: '项目站内信' }
          },
          {
            path: 'channels',
            name: 'NotificationChannels',
            component: () => import('@/views/message-center/channels/index.vue'),
            meta: { title: '通知渠道' }
          },
          {
            path: 'project-channels',
            name: 'ProjectNotificationChannels',
            component: () => import('@/views/message-center/channels/index.vue'), // 暂时使用相同组件
            meta: { title: '项目通知渠道' }
          },
          {
            path: 'robots',
            name: 'Robots',
            component: () => import('@/views/message-center/robots/index.vue'),
            meta: { title: '机器人管理' }
          },
          {
            path: 'project-robots',
            name: 'ProjectRobots',
            component: () => import('@/views/message-center/project-robots/index.vue'),
            meta: { title: '项目机器人管理' }
          },
          {
            path: 'receivers',
            name: 'Receivers',
            component: () => import('@/views/message-center/receivers/index.vue'),
            meta: { title: '接收人管理' }
          },
          {
            path: 'project-receivers',
            name: 'ProjectReceivers',
            component: () => import('@/views/message-center/receivers/index.vue'), // 暂时使用相同组件
            meta: { title: '项目接收人管理' }
          },
          {
            path: 'subscriptions',
            name: 'Subscriptions',
            component: () => import('@/views/message-center/subscriptions/index.vue'),
            meta: { title: '消息订阅' }
          },
          {
            path: 'project-subscriptions',
            name: 'ProjectSubscriptions',
            component: () => import('@/views/message-center/subscriptions/index.vue'), // 暂时使用相同组件
            meta: { title: '项目消息订阅' }
          }
        ]
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  if (to.path === '/login') {
    next()
  } else {
    if (token) {
      next()
    } else {
      next('/login')
    }
  }
})

export default router
