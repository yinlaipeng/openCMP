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
        path: '/middleware',
        name: 'Middleware',
        meta: { title: '中间件', icon: 'Promotion' },
        children: [
          {
            path: 'message-queue/kafka',
            name: 'MiddlewareKafka',
            component: () => import('@/views/middleware/message-queue/kafka/index.vue'),
            meta: { title: 'Kafka' }
          },
          {
            path: 'data-analysis/elasticsearch',
            name: 'MiddlewareElasticsearch',
            component: () => import('@/views/middleware/data-analysis/elasticsearch/index.vue'),
            meta: { title: 'Elasticsearch' }
          }
        ]
      },
      {
        path: '/container',
        name: 'Container',
        meta: { title: '容器', icon: 'Box' },
        children: [
          {
            path: 'service/kubernetes',
            name: 'ContainerKubernetes',
            component: () => import('@/views/container/service/kubernetes/index.vue'),
            meta: { title: 'Kubernetes' }
          },
          {
            path: 'service/image-repositories',
            name: 'ContainerImageRepositories',
            component: () => import('@/views/container/service/image-repositories/index.vue'),
            meta: { title: '镜像仓库' }
          }
        ]
      },
      {
        path: '/monitoring',
        name: 'Monitoring',
        meta: { title: '监控', icon: 'DataLine' },
        children: [
          // 纵览
          {
            path: 'overview/dashboard',
            name: 'MonitoringDashboard',
            component: () => import('@/views/monitoring/overview/dashboard/index.vue'),
            meta: { title: '大盘' }
          },
          // 资源
          {
            path: 'resources/vms',
            name: 'MonitoringResourcesVMs',
            component: () => import('@/views/monitoring/resources/vms/index.vue'),
            meta: { title: '虚拟机' }
          },
          // 监控
          {
            path: 'monitor/dashboards',
            name: 'MonitoringDashboards',
            component: () => import('@/views/monitoring/monitor/dashboards/index.vue'),
            meta: { title: '监控面板' }
          },
          {
            path: 'monitor/query',
            name: 'MonitoringQuery',
            component: () => import('@/views/monitoring/monitor/query/index.vue'),
            meta: { title: '监控查询' }
          },
          // 告警
          {
            path: 'alerts/policies',
            name: 'MonitoringAlertPolicies',
            component: () => import('@/views/monitoring/alerts/policies/index.vue'),
            meta: { title: '告警策略' }
          },
          {
            path: 'alerts/resources',
            name: 'MonitoringAlertResources',
            component: () => import('@/views/monitoring/alerts/resources/index.vue'),
            meta: { title: '告警资源' }
          },
          {
            path: 'alerts/history',
            name: 'MonitoringAlertHistory',
            component: () => import('@/views/monitoring/alerts/history/index.vue'),
            meta: { title: '告警历史' }
          },
          {
            path: 'alerts/blocked',
            name: 'MonitoringBlockedResources',
            component: () => import('@/views/monitoring/alerts/blocked/index.vue'),
            meta: { title: '屏蔽资源' }
          }
        ]
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
        component: () => import('@/views/compute/index.vue'), // Parent component that renders children
        meta: { title: '主机', icon: 'Cpu' },
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
          },
          {
            path: 'host-templates',
            name: 'HostTemplates',
            component: () => import('@/views/compute/host-templates/index.vue'),
            meta: { title: '主机模版' }
          },
          {
            path: 'autoscaling-groups',
            name: 'AutoscalingGroups',
            component: () => import('@/views/compute/autoscaling-groups/index.vue'),
            meta: { title: '弹性伸缩组' }
          },
          {
            path: 'images/system-images',
            name: 'SystemImages',
            component: () => import('@/views/compute/images/system-images.vue'),
            meta: { title: '系统镜像' }
          },
          {
            path: 'disks',
            name: 'Disks',
            component: () => import('@/views/compute/storage/disks/index.vue'),
            meta: { title: '硬盘' }
          },
          {
            path: 'disk-snapshots',
            name: 'DiskSnapshots',
            component: () => import('@/views/compute/storage/disk-snapshots/index.vue'),
            meta: { title: '硬盘快照' }
          },
          {
            path: 'host-snapshots',
            name: 'HostSnapshots',
            component: () => import('@/views/compute/storage/host-snapshots/index.vue'),
            meta: { title: '主机快照' }
          },
          {
            path: 'snapshot-policies',
            name: 'SnapshotPolicies',
            component: () => import('@/views/compute/storage/snapshot-policies/index.vue'),
            meta: { title: '自动快照策略' }
          },
          {
            path: 'security-groups',
            name: 'ComputeSecurityGroups',
            component: () => import('@/views/compute/network/security-groups/index.vue'),
            meta: { title: '安全组' }
          },
          {
            path: 'subnets',
            name: 'IPSubnets',
            component: () => import('@/views/compute/network/ip-subnets/index.vue'),
            meta: { title: 'IP子网' }
          },
          {
            path: 'eips',
            name: 'ElasticIPs',
            component: () => import('@/views/compute/network/eips/index.vue'),
            meta: { title: '弹性公网IP' }
          },
          {
            path: 'keys',
            name: 'SSHKeys',
            component: () => import('@/views/compute/keys/keys/index.vue'),
            meta: { title: '密钥' }
          }
        ]
      },
      {
        path: '/network',
        name: 'Network',
        meta: { title: '网络', icon: 'Connection' },
        children: [
          // 地域
          {
            path: 'geography/regions',
            name: 'NetworkGeographyRegions',
            component: () => import('@/views/network/geography/regions/index.vue'),
            meta: { title: '区域' }
          },
          {
            path: 'geography/zones',
            name: 'NetworkGeographyZones',
            component: () => import('@/views/network/geography/zones/index.vue'),
            meta: { title: '可用区' }
          },
          // 基础网络
          {
            path: 'basic/vpc-interconnect',
            name: 'VPCInterconnect',
            component: () => import('@/views/network/vpc-interconnect/index.vue'),
            meta: { title: 'vpc互联' }
          },
          {
            path: 'basic/vpc-peering',
            name: 'VPCPeering',
            component: () => import('@/views/network/vpc-peering/index.vue'),
            meta: { title: 'vpc对等连接' }
          },
          {
            path: 'basic/global-vpc',
            name: 'GlobalVPC',
            component: () => import('@/views/network/global-vpc/index.vue'),
            meta: { title: '全局vpc' }
          },
          {
            path: 'basic/vpcs',
            name: 'VPCs',
            component: () => import('@/views/network/vpcs/index.vue'),
            meta: { title: 'vpc' }
          },
          {
            path: 'basic/route-tables',
            name: 'RouteTables',
            component: () => import('@/views/network/routes/index.vue'),
            meta: { title: '路由表' }
          },
          {
            path: 'basic/l2-networks',
            name: 'L2Networks',
            component: () => import('@/views/network/l2-networks/index.vue'),
            meta: { title: '二层网络' }
          },
          {
            path: 'basic/subnets',
            name: 'Subnets',
            component: () => import('@/views/network/subnets/index.vue'),
            meta: { title: 'ip子网' }
          },
          // 网络服务
          {
            path: 'services/eips',
            name: 'NetworkServicesEIPs',
            component: () => import('@/views/network/services/eips/index.vue'),
            meta: { title: '弹性公网ip' }
          },
          {
            path: 'services/nat-gateways',
            name: 'NetworkServicesNATGateways',
            component: () => import('@/views/network/services/nat-gateways/index.vue'),
            meta: { title: 'nat网关' }
          },
          {
            path: 'services/dns',
            name: 'NetworkServicesDNS',
            component: () => import('@/views/network/services/dns/index.vue'),
            meta: { title: 'dns解析' }
          },
          {
            path: 'services/ipv6-gateways',
            name: 'NetworkServicesIPv6Gateways',
            component: () => import('@/views/network/services/ipv6-gateways/index.vue'),
            meta: { title: 'ipv6网关' }
          },
          // 网络安全
          {
            path: 'security/waf-policies',
            name: 'NetworkSecurityWAFPolicies',
            component: () => import('@/views/network/security/waf-policies/index.vue'),
            meta: { title: 'waf策略' }
          },
          {
            path: 'security/app-services',
            name: 'NetworkSecurityAppServices',
            component: () => import('@/views/network/security/app-services/index.vue'),
            meta: { title: '应用程序服务' }
          },
          // 负载均衡
          {
            path: 'loadbalancer/instances',
            name: 'NetworkLoadBalancerInstances',
            component: () => import('@/views/network/loadbalancer/instances/index.vue'),
            meta: { title: '实例' }
          },
          {
            path: 'loadbalancer/acls',
            name: 'NetworkLoadBalancerACLs',
            component: () => import('@/views/network/loadbalancer/acls/index.vue'),
            meta: { title: '访问控制' }
          },
          {
            path: 'loadbalancer/certificates',
            name: 'NetworkLoadBalancerCertificates',
            component: () => import('@/views/network/loadbalancer/certificates/index.vue'),
            meta: { title: '证书' }
          },
          // 内容分发网络
          {
            path: 'cdn/domains',
            name: 'NetworkCDNDomains',
            component: () => import('@/views/network/cdn/domains/index.vue'),
            meta: { title: 'cdn域名' }
          }
        ]
      },
      {
        path: '/storage',
        name: 'Storage',
        meta: { title: '存储', icon: 'Folder' },
        children: [
          // 块存储
          {
            path: 'block/block-storage',
            name: 'StorageBlockStorage',
            component: () => import('@/views/storage/block/block-storage/index.vue'),
            meta: { title: '块存储' }
          },
          // 对象存储
          {
            path: 'object/buckets',
            name: 'StorageObjectBuckets',
            component: () => import('@/views/storage/object/buckets/index.vue'),
            meta: { title: '存储桶' }
          },
          // 表格存储
          {
            path: 'table/table-storage',
            name: 'StorageTableStorage',
            component: () => import('@/views/storage/table/table-storage/index.vue'),
            meta: { title: '表格存储' }
          },
          // 文件存储
          {
            path: 'file/file-systems',
            name: 'StorageFileSystems',
            component: () => import('@/views/storage/file/file-systems/index.vue'),
            meta: { title: '文件系统' }
          },
          {
            path: 'file/nas-groups',
            name: 'StorageNASGroups',
            component: () => import('@/views/storage/file/nas-groups/index.vue'),
            meta: { title: 'NAS权限组' }
          }
        ]
      },
      {
        path: '/database',
        name: 'Database',
        meta: { title: '数据库', icon: 'Coin' },
        children: [
          // RDS
          {
            path: 'rds/instances',
            name: 'DatabaseRDSInstances',
            component: () => import('@/views/database/rds/instances/index.vue'),
            meta: { title: 'RDS实例' }
          },
          // Redis
          {
            path: 'redis/instances',
            name: 'DatabaseRedisInstances',
            component: () => import('@/views/database/redis/instances/index.vue'),
            meta: { title: 'Redis实例' }
          },
          // MongoDB
          {
            path: 'mongodb/instances',
            name: 'DatabaseMongoDBInstances',
            component: () => import('@/views/database/mongodb/instances/index.vue'),
            meta: { title: 'MongoDB实例' }
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
