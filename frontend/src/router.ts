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
        meta: { title: '云账户管理', icon: 'Cloudy' }
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
            path: 'users',
            name: 'Users',
            component: () => import('@/views/iam/users/index.vue'),
            meta: { title: '用户管理' }
          },
          {
            path: 'roles',
            name: 'Roles',
            component: () => import('@/views/iam/roles/index.vue'),
            meta: { title: '角色权限' }
          },
          {
            path: 'permissions',
            name: 'Permissions',
            component: () => import('@/views/iam/permissions/index.vue'),
            meta: { title: '权限管理' }
          },
          {
            path: 'policies',
            name: 'Policies',
            component: () => import('@/views/iam/policies/index.vue'),
            meta: { title: '策略管理' }
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
