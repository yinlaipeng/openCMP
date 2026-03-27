<template>
  <el-container class="layout-container">
    <el-aside width="220px" class="sidebar">
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
      >
        <el-menu-item index="/cloud-accounts">
          <el-icon><Cloudy /></el-icon>
          <span>云账户管理</span>
        </el-menu-item>
        
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
        
        <el-sub-menu index="iam">
          <template #title>
            <el-icon><Lock /></el-icon>
            <span>认证与安全</span>
          </template>
          <el-menu-item index="/iam/auth-sources">认证源</el-menu-item>
          <el-menu-item index="/iam/users">用户管理</el-menu-item>
          <el-menu-item index="/iam/roles">角色权限</el-menu-item>
          <el-menu-item index="/iam/messages">消息中心</el-menu-item>
          <el-menu-item index="/iam/alerts">安全告警</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    
    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-breadcrumb separator="/">
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
import { computed, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { User, ArrowDown, SwitchButton, Lock } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

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
  if (path.startsWith('/iam')) return 'iam'
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
    ElMessage.success('已退出登录')
    router.push('/login')
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}
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
  flex: 1;
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
