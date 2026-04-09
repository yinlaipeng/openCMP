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
            <span>计算资源</span>
          </template>
          <el-menu-item index="/compute/vms">虚拟机管理</el-menu-item>
          <el-menu-item index="/compute/images">镜像管理</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="network">
          <template #title>
            <el-icon><Connection /></el-icon>
            <span>网络资源</span>
          </template>
          <el-menu-item index="/network/vpcs">VPC 管理</el-menu-item>
          <el-menu-item index="/network/subnets">子网管理</el-menu-item>
          <el-menu-item index="/network/security-groups">安全组管理</el-menu-item>
          <el-menu-item index="/network/eips">弹性 IP</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="cloud-management">
          <template #title>
            <el-icon><Cloudy /></el-icon>
            <span>多云管理</span>
          </template>
          <el-menu-item index="/cloud-accounts">云账户管理</el-menu-item>
          <el-menu-item index="/cloud-management/sync-policies">同步策略</el-menu-item>
          <el-menu-item index="/scheduled-tasks">定时同步任务</el-menu-item>
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
          <!-- 根据环境显示不同的安全告警选项 -->
          <el-menu-item v-if="currentEnvironment.type === 'management_console'" index="/iam/alerts">安全告警</el-menu-item>
          <el-menu-item v-else index="/iam/project-alerts">项目安全告警</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="message-center">
          <template #title>
            <el-icon><Bell /></el-icon>
            <span>消息中心</span>
          </template>
          <!-- 根据环境显示不同的消息中心选项 -->
          <el-menu-item v-if="currentEnvironment.type === 'management_console'" index="/message-center/inbox">站内信</el-menu-item>
          <el-menu-item v-else index="/message-center/project-inbox">项目站内信</el-menu-item>
          <el-menu-item v-if="currentEnvironment.type === 'management_console'" index="/message-center/channels">通知渠道</el-menu-item>
          <el-menu-item v-else index="/message-center/project-channels">项目通知渠道</el-menu-item>
          <el-menu-item v-if="currentEnvironment.type === 'management_console'" index="/message-center/robots">机器人管理</el-menu-item>
          <el-menu-item v-else index="/message-center/project-robots">项目机器人管理</el-menu-item>
          <el-menu-item v-if="currentEnvironment.type === 'management_console'" index="/message-center/receivers">接收人管理</el-menu-item>
          <el-menu-item v-else index="/message-center/project-receivers">项目接收人管理</el-menu-item>
          <el-menu-item v-if="currentEnvironment.type === 'management_console'" index="/message-center/subscriptions">消息订阅</el-menu-item>
          <el-menu-item v-else index="/message-center/project-subscriptions">项目消息订阅</el-menu-item>
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
  </el-container>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import {
  User, ArrowDown, SwitchButton, Lock, Bell, Cloudy, Cpu, Connection,
  Grid, Monitor, FolderOpened, ArrowDown as IconArrowDown
} from '@element-plus/icons-vue'
import { getProjects } from '@/api/iam' // 使用iam中的项目API
import { initializeProjectContext, setProjectContext, getCurrentProjectId, getCurrentProjectName, isInProjectMode, clearProjectContext } from '@/utils/projectContext'

const router = useRouter()
const route = useRoute()

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
      ElMessage.info('修改密码功能开发中')
      break
    case 'profile':
      ElMessage.info('个人信息功能开发中')
      break
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