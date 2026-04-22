<template>
  <el-container class="layout-container">
    <!-- 根据路由元信息决定是否显示侧边栏 -->
    <el-aside v-if="!route.meta.hideSidebar" width="220px" class="sidebar">
      <div class="logo">
        <el-icon><Cloudy /></el-icon>
        <span>openCMP</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        router
        :collapse="false"
      >
        <el-sub-menu index="compute">
          <template #title>
            <el-icon><Cpu /></el-icon>
            <span>主机</span>
          </template>
          <el-sub-menu index="compute-host">
            <template #title>
              <span>主机</span>
            </template>
            <el-menu-item index="/compute/vms">虚拟机</el-menu-item>
            <el-menu-item index="/compute/host-templates">主机模版</el-menu-item>
            <el-menu-item index="/compute/autoscaling-groups">弹性伸缩组</el-menu-item>
          </el-sub-menu>
          <el-sub-menu index="compute-images">
            <template #title>
              <span>镜像</span>
            </template>
            <el-menu-item index="/compute/images">系统镜像</el-menu-item>
          </el-sub-menu>
          <el-sub-menu index="compute-storage">
            <template #title>
              <span>存储</span>
            </template>
            <el-menu-item index="/compute/disks">硬盘</el-menu-item>
            <el-menu-item index="/compute/disk-snapshots">硬盘快照</el-menu-item>
            <el-menu-item index="/compute/host-snapshots">主机快照</el-menu-item>
            <el-menu-item index="/compute/snapshot-policies">自动快照策略</el-menu-item>
          </el-sub-menu>
          <el-sub-menu index="compute-network">
            <template #title>
              <span>网络</span>
            </template>
            <el-menu-item index="/compute/security-groups">安全组</el-menu-item>
            <el-menu-item index="/compute/subnets">IP子网</el-menu-item>
            <el-menu-item index="/compute/eips">弹性公网IP</el-menu-item>
          </el-sub-menu>
          <el-sub-menu index="compute-keys">
            <template #title>
              <span>密钥</span>
            </template>
            <el-menu-item index="/compute/keys">密钥</el-menu-item>
          </el-sub-menu>
        </el-sub-menu>

        <el-sub-menu index="network">
          <template #title>
            <el-icon><Connection /></el-icon>
            <span>网络</span>
          </template>

          <!-- 地域子菜单 -->
          <el-sub-menu index="network-geography">
            <template #title>
              <span>地域</span>
            </template>
            <el-menu-item index="/network/geography/regions">区域</el-menu-item>
            <el-menu-item index="/network/geography/zones">可用区</el-menu-item>
          </el-sub-menu>

          <!-- 基础网络子菜单 -->
          <el-sub-menu index="network-basic">
            <template #title>
              <span>基础网络</span>
            </template>
            <el-menu-item index="/network/basic/vpc-interconnect">vpc互联</el-menu-item>
            <el-menu-item index="/network/basic/vpc-peering">vpc对等连接</el-menu-item>
            <el-menu-item index="/network/basic/global-vpc">全局vpc</el-menu-item>
            <el-menu-item index="/network/basic/vpcs">vpc</el-menu-item>
            <el-menu-item index="/network/basic/route-tables">路由表</el-menu-item>
            <el-menu-item index="/network/basic/l2-networks">二层网络</el-menu-item>
            <el-menu-item index="/network/basic/subnets">ip子网</el-menu-item>
          </el-sub-menu>

          <!-- 网络服务子菜单 -->
          <el-sub-menu index="network-services">
            <template #title>
              <span>网络服务</span>
            </template>
            <el-menu-item index="/network/services/eips">弹性公网ip</el-menu-item>
            <el-menu-item index="/network/services/nat-gateways">nat网关</el-menu-item>
            <el-menu-item index="/network/services/dns">dns解析</el-menu-item>
            <el-menu-item index="/network/services/ipv6-gateways">ipv6网关</el-menu-item>
          </el-sub-menu>

          <!-- 网络安全子菜单 -->
          <el-sub-menu index="network-security">
            <template #title>
              <span>网络安全</span>
            </template>
            <el-menu-item index="/network/security/waf-policies">waf策略</el-menu-item>
            <el-menu-item index="/network/security/app-services">应用程序服务</el-menu-item>
          </el-sub-menu>

          <!-- 负载均衡子菜单 -->
          <el-sub-menu index="network-loadbalancer">
            <template #title>
              <span>负载均衡</span>
            </template>
            <el-menu-item index="/network/loadbalancer/instances">实例</el-menu-item>
            <el-menu-item index="/network/loadbalancer/acls">访问控制</el-menu-item>
            <el-menu-item index="/network/loadbalancer/certificates">证书</el-menu-item>
          </el-sub-menu>

          <!-- 内容分发网络子菜单 -->
          <el-sub-menu index="network-cdn">
            <template #title>
              <span>内容分发网络</span>
            </template>
            <el-menu-item index="/network/cdn/domains">cdn域名</el-menu-item>
          </el-sub-menu>
        </el-sub-menu>

        <!-- 存储一级菜单 -->
        <el-sub-menu index="storage">
          <template #title>
            <el-icon><Folder /></el-icon>
            <span>存储</span>
          </template>

          <!-- 块存储子菜单 -->
          <el-sub-menu index="storage-block">
            <template #title>
              <span>块存储</span>
            </template>
            <el-menu-item index="/storage/block/block-storage">块存储</el-menu-item>
          </el-sub-menu>

          <!-- 对象存储子菜单 -->
          <el-sub-menu index="storage-object">
            <template #title>
              <span>对象存储</span>
            </template>
            <el-menu-item index="/storage/object/buckets">存储桶</el-menu-item>
          </el-sub-menu>

          <!-- 表格存储子菜单 -->
          <el-sub-menu index="storage-table">
            <template #title>
              <span>表格存储</span>
            </template>
            <el-menu-item index="/storage/table/table-storage">表格存储</el-menu-item>
          </el-sub-menu>

          <!-- 文件存储子菜单 -->
          <el-sub-menu index="storage-file">
            <template #title>
              <span>文件存储</span>
            </template>
            <el-menu-item index="/storage/file/file-systems">文件系统</el-menu-item>
            <el-menu-item index="/storage/file/nas-groups">NAS权限组</el-menu-item>
          </el-sub-menu>
        </el-sub-menu>

        <!-- 数据库一级菜单 -->
        <el-sub-menu index="database">
          <template #title>
            <el-icon><Coin /></el-icon>
            <span>数据库</span>
          </template>

          <!-- RDS子菜单 -->
          <el-sub-menu index="database-rds">
            <template #title>
              <span>RDS</span>
            </template>
            <el-menu-item index="/database/rds/instances">RDS实例</el-menu-item>
          </el-sub-menu>

          <!-- Redis子菜单 -->
          <el-sub-menu index="database-redis">
            <template #title>
              <span>Redis</span>
            </template>
            <el-menu-item index="/database/redis/instances">Redis实例</el-menu-item>
          </el-sub-menu>

          <!-- MongoDB子菜单 -->
          <el-sub-menu index="database-mongodb">
            <template #title>
              <span>MongoDB</span>
            </template>
            <el-menu-item index="/database/mongodb/instances">MongoDB实例</el-menu-item>
          </el-sub-menu>
        </el-sub-menu>

        <!-- 中间件一级菜单 -->
        <el-sub-menu index="middleware">
          <template #title>
            <el-icon><Promotion /></el-icon>
            <span>中间件</span>
          </template>

          <!-- 消息队列子菜单 -->
          <el-sub-menu index="middleware-message-queue">
            <template #title>
              <span>消息队列</span>
            </template>
            <el-menu-item index="/middleware/message-queue/kafka">Kafka</el-menu-item>
          </el-sub-menu>

          <!-- 数据分析子菜单 -->
          <el-sub-menu index="middleware-data-analysis">
            <template #title>
              <span>数据分析</span>
            </template>
            <el-menu-item index="/middleware/data-analysis/elasticsearch">Elasticsearch</el-menu-item>
          </el-sub-menu>
        </el-sub-menu>

        <!-- 容器一级菜单 -->
        <el-sub-menu index="container">
          <template #title>
            <el-icon><Box /></el-icon>
            <span>容器</span>
          </template>

          <!-- 容器服务子菜单 -->
          <el-sub-menu index="container-service">
            <template #title>
              <span>容器服务</span>
            </template>
            <el-menu-item index="/container/service/kubernetes">Kubernetes</el-menu-item>
            <el-menu-item index="/container/service/image-repositories">镜像仓库</el-menu-item>
          </el-sub-menu>
        </el-sub-menu>

        <!-- 监控一级菜单 -->
        <el-sub-menu index="monitoring">
          <template #title>
            <el-icon><DataLine /></el-icon>
            <span>监控</span>
          </template>

          <!-- 纵览子菜单 -->
          <el-sub-menu index="monitoring-overview">
            <template #title>
              <span>纵览</span>
            </template>
            <el-menu-item index="/monitoring/overview/dashboard">大盘</el-menu-item>
          </el-sub-menu>

          <!-- 资源子菜单 -->
          <el-sub-menu index="monitoring-resources">
            <template #title>
              <span>资源</span>
            </template>
            <el-menu-item index="/monitoring/resources/vms">虚拟机</el-menu-item>
          </el-sub-menu>

          <!-- 监控子菜单 -->
          <el-sub-menu index="monitoring-monitor">
            <template #title>
              <span>监控</span>
            </template>
            <el-menu-item index="/monitoring/monitor/dashboards">监控面板</el-menu-item>
            <el-menu-item index="/monitoring/monitor/query">监控查询</el-menu-item>
          </el-sub-menu>

          <!-- 告警子菜单 -->
          <el-sub-menu index="monitoring-alerts">
            <template #title>
              <span>告警</span>
            </template>
            <el-menu-item index="/monitoring/alerts/policies">告警策略</el-menu-item>
            <el-menu-item index="/monitoring/alerts/resources">告警资源</el-menu-item>
            <el-menu-item index="/monitoring/alerts/history">告警历史</el-menu-item>
            <el-menu-item index="/monitoring/alerts/blocked">屏蔽资源</el-menu-item>
          </el-sub-menu>
        </el-sub-menu>

        <el-sub-menu index="cloud-management">
          <template #title>
            <el-icon><Cloudy /></el-icon>
            <span>多云管理</span>
          </template>
          <el-menu-item index="/cloud-accounts">云账户管理</el-menu-item>
          <el-menu-item index="/cloud-management/sync-policies">同步策略</el-menu-item>
          <el-menu-item index="/cloud-management/cloud-user-groups">云用户组</el-menu-item>
          <el-menu-item index="/cloud-management/proxies">代理</el-menu-item>
          <el-menu-item index="/scheduled-tasks">定时同步任务</el-menu-item>
          <el-menu-item index="/sync-logs">同步日志</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="iam">
          <template #title>
            <el-icon><Lock /></el-icon>
            <span>认证与安全</span>
          </template>
          <el-menu-item index="/iam/auth-sources">认证源</el-menu-item>
          <el-menu-item index="/iam/domains">域</el-menu-item>
          <el-menu-item index="/iam/projects">项目</el-menu-item>
          <el-menu-item index="/iam/users">用户管理</el-menu-item>
          <el-menu-item index="/iam/groups">用户组</el-menu-item>
          <el-menu-item index="/iam/roles">角色</el-menu-item>
          <el-menu-item index="/iam/permissions">权限</el-menu-item>
          <!-- 安全告警选项 -->
          <el-menu-item index="/iam/alerts">安全告警</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="message-center">
          <template #title>
            <el-icon><Bell /></el-icon>
            <span>消息中心</span>
          </template>
          <!-- 根据环境显示不同的消息中心选项 -->
          <el-menu-item v-if="currentEnvironment.type === 'management_console'" index="/message-center/inbox">站内信</el-menu-item>
          <el-menu-item v-else index="/message-center/project-inbox">项目站内信</el-menu-item>
          <el-menu-item v-if="currentEnvironment.type === 'management_console'" index="/message-center/channels">通知渠道设置</el-menu-item>
          <el-menu-item v-else index="/message-center/project-channels">项目通知渠道设置</el-menu-item>
          <el-menu-item v-if="currentEnvironment.type === 'management_console'" index="/message-center/robots">机器人管理</el-menu-item>
          <el-menu-item v-else index="/message-center/project-robots">项目机器人管理</el-menu-item>
          <el-menu-item v-if="currentEnvironment.type === 'management_console'" index="/message-center/receivers">接收人管理</el-menu-item>
          <el-menu-item v-else index="/message-center/project-receivers">项目接收人管理</el-menu-item>
          <el-menu-item v-if="currentEnvironment.type === 'management_console'" index="/message-center/message-types">消息订阅设置</el-menu-item>
          <el-menu-item v-else index="/message-center/project-subscriptions">项目消息订阅</el-menu-item>
        </el-sub-menu>

        <!-- 费用中心 -->
        <el-sub-menu index="finance">
          <template #title>
            <el-icon><Wallet /></el-icon>
            <span>费用中心</span>
          </template>
          <!-- 订单管理 -->
          <el-sub-menu index="finance-orders">
            <template #title>订单管理</template>
            <el-menu-item index="/finance/orders/my-orders">我的订单</el-menu-item>
            <el-menu-item index="/finance/orders/renewals">续费管理</el-menu-item>
          </el-sub-menu>
          <!-- 费用账单 -->
          <el-sub-menu index="finance-bills">
            <template #title>费用账单</template>
            <el-menu-item index="/finance/bills/view">账单查看</el-menu-item>
            <el-menu-item index="/finance/bills/export">账单导出中心</el-menu-item>
          </el-sub-menu>
          <!-- 成本管理 -->
          <el-sub-menu index="finance-cost">
            <template #title>成本管理</template>
            <el-menu-item index="/finance/cost/analysis">成本分析</el-menu-item>
            <el-menu-item index="/finance/cost/reports">成本报告</el-menu-item>
            <el-menu-item index="/finance/cost/budgets">预算管理</el-menu-item>
            <el-menu-item index="/finance/cost/anomaly">异常监测</el-menu-item>
          </el-sub-menu>
        </el-sub-menu>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-left">
          <!-- 环境切换下拉菜单 -->
          <el-dropdown @command="handleEnvironmentChange" class="environment-dropdown">
            <span class="environment-selector">
              <el-icon><Grid /></el-icon>
              <span class="environment-text">{{ currentEnvironment.name }}</span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="management_console" :class="{ 'selected': currentEnvironment.type === 'management_console' }">
                  <el-icon><Monitor /></el-icon>
                  管理控制台
                </el-dropdown-item>
                <el-dropdown-item divided />
                <el-dropdown-item
                  v-for="project in projectList"
                  :key="project.id"
                  :command="`project_${project.id}`"
                  :class="{ 'selected': currentEnvironment.type === 'project' && currentEnvironment.id === project.id }"
                >
                  <el-icon><FolderOpened /></el-icon>
                  {{ project.name }}
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>

          <!-- 原有的面包屑导航（可选保留） -->
          <el-breadcrumb separator="/" v-if="showBreadcrumb">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentRoute">{{ currentRoute }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32" :icon="User" />
              <span class="username">{{ currentUser?.display_name || currentUser?.name || '用户' }}</span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人信息
                </el-dropdown-item>
                <el-dropdown-item command="password">
                  <el-icon><Lock /></el-icon>
                  修改密码
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>

    <!-- 修改密码对话框 -->
    <el-dialog v-model="passwordDialogVisible" title="修改密码" width="400px">
      <el-form :model="passwordForm" label-width="100px" :rules="passwordRules" ref="passwordFormRef">
        <el-form-item label="当前密码" prop="oldPassword">
          <el-input v-model="passwordForm.oldPassword" type="password" placeholder="请输入当前密码" show-password />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="passwordForm.newPassword" type="password" placeholder="请输入新密码（至少6位）" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" placeholder="请再次输入新密码" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="passwordDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handlePasswordSubmit" :loading="passwordLoading">确认修改</el-button>
      </template>
    </el-dialog>

    <!-- 个人信息对话框 -->
    <el-dialog v-model="profileDialogVisible" title="个人信息" width="500px">
      <el-form :model="profileForm" label-width="100px">
        <el-form-item label="用户名">
          <el-input :value="currentUser?.name" disabled />
        </el-form-item>
        <el-form-item label="显示名称">
          <el-input v-model="profileForm.display_name" placeholder="请输入显示名称" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="profileForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="profileForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="profileForm.remark" type="textarea" placeholder="请输入备注" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="profileDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleProfileSubmit" :loading="profileLoading">保存</el-button>
      </template>
    </el-dialog>
  </el-container>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import {
  User, ArrowDown, SwitchButton, Lock, Bell, Cloudy, Cpu, Connection,
  Grid, Monitor, FolderOpened, ArrowDown as IconArrowDown, Folder, Coin, Promotion, Box, DataLine, Wallet
} from '@element-plus/icons-vue'
import { getProjects } from '@/api/iam' // 使用iam中的项目API
import { initializeProjectContext, setProjectContext, getCurrentProjectId, getCurrentProjectName, isInProjectMode, clearProjectContext } from '@/utils/projectContext'
import { changePassword, updateProfile } from '@/api/auth'

const router = useRouter()
const route = useRoute()

// 密码修改对话框
const passwordDialogVisible = ref(false)
const passwordLoading = ref(false)
const passwordFormRef = ref<any>(null)
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})
const passwordRules = {
  oldPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [{ required: true, min: 6, message: '密码长度至少6位', trigger: 'blur' }],
  confirmPassword: [{ required: true, message: '请确认新密码', trigger: 'blur' }]
}

// 个人信息对话框
const profileDialogVisible = ref(false)
const profileLoading = ref(false)
const profileForm = ref({
  display_name: '',
  email: '',
  phone: '',
  remark: ''
})

// 环境状态
const currentEnvironment = ref({
  type: 'management_console', // 'management_console' 或 'project'
  id: null,
  name: '管理控制台'
})

// 项目列表
const projectList = ref([])

// 是否显示面包屑（根据需要决定）
const showBreadcrumb = computed(() => {
  // 如果不是管理控制台环境，可以隐藏面包屑或显示项目相关路径
  return currentEnvironment.value.type === 'management_console' && !route.meta?.hideSidebar
})

// 加载项目列表
const loadProjects = async () => {
  try {
    const response = await getProjects()
    projectList.value = response.items || []
  } catch (error) {
    console.error('Failed to load projects:', error)
    ElMessage.error('加载项目列表失败')
  }
}

// 环境切换处理
const handleEnvironmentChange = (command: string) => {
  if (command === 'management_console') {
    currentEnvironment.value = {
      type: 'management_console',
      id: null,
      name: '管理控制台'
    }
    // 清除项目上下文，返回管理控制台视图
    clearProjectContext()

    // 刷新页面以确保所有组件都更新到管理控制台状态
    router.go(0)
  } else if (command.startsWith('project_')) {
    const projectId = parseInt(command.split('_')[1])
    const project = projectList.value.find((proj: any) => proj.id === projectId)

    if (project) {
      currentEnvironment.value = {
        type: 'project',
        id: projectId,
        name: project.name
      }

      // 存储项目上下文
      setProjectContext(projectId, project.name)

      // 刷新页面以确保所有组件都更新到项目状态
      router.go(0)
    }
  }
}

const currentUser = ref<any>(() => {
  try {
    const userStr = localStorage.getItem('user')
    return userStr ? JSON.parse(userStr) : null
  } catch {
    return null
  }
})

const activeMenu = computed(() => {
  const path = route.path
  if (path.startsWith('/compute')) return 'compute'
  if (path.startsWith('/network')) return 'network'
  if (path.startsWith('/storage')) return 'storage'
  if (path.startsWith('/database')) return 'database'
  if (path.startsWith('/middleware')) return 'middleware'
  if (path.startsWith('/container')) return 'container'
  if (path.startsWith('/monitoring')) return 'monitoring'
  if (path.startsWith('/iam')) {
    // 根据路径判断是否为项目相关页面来返回正确的菜单项
    if (path.includes('/project-alerts')) {
      return 'iam'
    }
    return 'iam'
  }
  if (path.startsWith('/message-center')) {
    // 根据路径判断是否为项目相关页面来返回正确的菜单项
    if (path.includes('/project-')) {
      return 'message-center'
    }
    return 'message-center'
  }
  return path
})

const currentRoute = computed(() => {
  return route.meta.title as string || ''
})

const handleCommand = async (command: string) => {
  switch (command) {
    case 'logout':
      await handleLogout()
      break
    case 'password':
      passwordDialogVisible.value = true
      passwordForm.value = { oldPassword: '', newPassword: '', confirmPassword: '' }
      break
    case 'profile':
      profileDialogVisible.value = true
      profileForm.value = {
        display_name: currentUser.value?.display_name || '',
        email: currentUser.value?.email || '',
        phone: currentUser.value?.phone || '',
        remark: currentUser.value?.remark || ''
      }
      break
  }
}

// 处理密码修改
const handlePasswordSubmit = async () => {
  if (!passwordFormRef.value) return
  try {
    await passwordFormRef.value.validate()
    if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
      ElMessage.error('两次输入的密码不一致')
      return
    }
    passwordLoading.value = true
    await changePassword(passwordForm.value.oldPassword, passwordForm.value.newPassword)
    ElMessage.success('密码修改成功')
    passwordDialogVisible.value = false
  } catch (e: any) {
    if (e?.response?.data?.error) {
      ElMessage.error(e.response.data.error)
    } else if (e !== 'cancel') {
      ElMessage.error('密码修改失败')
    }
  } finally {
    passwordLoading.value = false
  }
}

// 处理个人信息更新
const handleProfileSubmit = async () => {
  profileLoading.value = true
  try {
    const res = await updateProfile(profileForm.value)
    ElMessage.success('个人信息更新成功')
    // 更新本地存储的用户信息
    const userStr = localStorage.getItem('user')
    if (userStr) {
      const user = JSON.parse(userStr)
      user.display_name = profileForm.value.display_name
      user.email = profileForm.value.email
      user.phone = profileForm.value.phone
      user.remark = profileForm.value.remark
      localStorage.setItem('user', JSON.stringify(user))
    }
    profileDialogVisible.value = false
  } catch (e: any) {
    if (e?.response?.data?.error) {
      ElMessage.error(e.response.data.error)
    } else {
      ElMessage.error('更新失败')
    }
  } finally {
    profileLoading.value = false
  }
}

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      type: 'warning'
    })
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    // 登出时也清除项目上下文
    clearProjectContext()
    ElMessage.success('已退出登录')
    router.push('/login')
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

// 初始化时加载项目列表和项目上下文
onMounted(() => {
  loadProjects()
  initializeProjectContext()

  // 根据当前存储的项目上下文设置环境
  const projectId = getCurrentProjectId()
  const projectName = getCurrentProjectName()

  if (projectId && projectName) {
    currentEnvironment.value = {
      type: 'project',
      id: projectId,
      name: projectName
    }
  } else {
    currentEnvironment.value = {
      type: 'management_console',
      id: null,
      name: '管理控制台'
    }
  }
})
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.sidebar {
  background-color: #304156;
  overflow-x: hidden;
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 60px;
  color: #fff;
  font-size: 20px;
  font-weight: bold;
  gap: 10px;
}

.logo .el-icon {
  font-size: 28px;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.environment-dropdown {
  margin-right: 20px;
}

.environment-selector {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  cursor: pointer;
  background-color: #f8f9fa;
  transition: all 0.3s;
}

.environment-selector:hover {
  border-color: #409eff;
  background-color: #ecf5ff;
}

.environment-text {
  font-size: 14px;
  color: #333;
  margin: 0 5px;
}

.selected {
  background-color: #f0f9ff;
  color: #409eff;
  font-weight: 500;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f5f5;
}

.username {
  font-size: 14px;
  color: #333;
}

.main-content {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>