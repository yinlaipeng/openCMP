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
          <el-tag size="small" type="success">多云管理平台</el-tag>
        </div>
      </el-header>
      
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const activeMenu = computed(() => {
  const path = route.path
  if (path.startsWith('/compute')) return 'compute'
  if (path.startsWith('/network')) return 'network'
  return path
})

const currentRoute = computed(() => {
  return route.meta.title as string || ''
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
  flex: 1;
}

.main-content {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>
