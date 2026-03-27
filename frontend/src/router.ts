import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/layout/index.vue'

const routes = [
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
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
