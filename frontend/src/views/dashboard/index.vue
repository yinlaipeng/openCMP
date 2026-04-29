<template>
  <div class="dashboard-page">
    <!-- 顶部统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon users">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.userCount }}</div>
              <div class="stat-label">用户总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon domains">
              <el-icon><OfficeBuilding /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.domainCount }}</div>
              <div class="stat-label">域数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon projects">
              <el-icon><Folder /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.projectCount }}</div>
              <div class="stat-label">项目数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon accounts">
              <el-icon><Cloudy /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.accountCount }}</div>
              <div class="stat-label">云账号数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 中间区域 -->
    <el-row :gutter="20" class="middle-row">
      <!-- 快捷入口 -->
      <el-col :span="8">
        <el-card shadow="hover" class="quick-links-card">
          <template #header>
            <div class="card-header">
              <span>快捷入口</span>
            </div>
          </template>
          <div class="quick-links">
            <div class="quick-link-item" @click="goTo('/cloud-accounts')">
              <el-icon><Cloudy /></el-icon>
              <span>云账号管理</span>
            </div>
            <div class="quick-link-item" @click="goTo('/iam/users')">
              <el-icon><User /></el-icon>
              <span>用户管理</span>
            </div>
            <div class="quick-link-item" @click="goTo('/iam/projects')">
              <el-icon><Folder /></el-icon>
              <span>项目管理</span>
            </div>
            <div class="quick-link-item" @click="goTo('/compute/vms')">
              <el-icon><Monitor /></el-icon>
              <span>虚拟机管理</span>
            </div>
            <div class="quick-link-item" @click="goTo('/monitoring/alerts/resources')">
              <el-icon><Bell /></el-icon>
              <span>资源告警</span>
            </div>
            <div class="quick-link-item" @click="goTo('/finance/bills')">
              <el-icon><Money /></el-icon>
              <span>账单查看</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 资源概览 -->
      <el-col :span="8">
        <el-card shadow="hover" class="resource-overview-card">
          <template #header>
            <div class="card-header">
              <span>资源概览</span>
              <el-button type="primary" size="small" text @click="refreshResources">刷新</el-button>
            </div>
          </template>
          <div class="resource-list">
            <div class="resource-item">
              <span class="resource-name">虚拟机</span>
              <span class="resource-count">{{ resourceStats.vmCount }}</span>
            </div>
            <div class="resource-item">
              <span class="resource-name">硬盘</span>
              <span class="resource-count">{{ resourceStats.diskCount }}</span>
            </div>
            <div class="resource-item">
              <span class="resource-name">弹性IP</span>
              <span class="resource-count">{{ resourceStats.eipCount }}</span>
            </div>
            <div class="resource-item">
              <span class="resource-name">安全组</span>
              <span class="resource-count">{{ resourceStats.sgCount }}</span>
            </div>
            <div class="resource-item">
              <span class="resource-name">数据库</span>
              <span class="resource-count">{{ resourceStats.dbCount }}</span>
            </div>
            <div class="resource-item">
              <span class="resource-name">负载均衡</span>
              <span class="resource-count">{{ resourceStats.lbCount }}</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 告警通知 -->
      <el-col :span="8">
        <el-card shadow="hover" class="alerts-card">
          <template #header>
            <div class="card-header">
              <span>告警通知</span>
              <el-button type="primary" size="small" text @click="goTo('/monitoring/alerts/resources')">查看全部</el-button>
            </div>
          </template>
          <div class="alert-list" v-if="alerts.length > 0">
            <div class="alert-item" v-for="alert in alerts" :key="alert.id">
              <el-tag :type="getAlertLevelType(alert.level)" size="small">{{ alert.level }}</el-tag>
              <span class="alert-title">{{ alert.title }}</span>
              <span class="alert-time">{{ alert.time }}</span>
            </div>
          </div>
          <div class="no-alerts" v-else>
            <el-icon><CircleCheck /></el-icon>
            <span>暂无告警</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 底部区域 -->
    <el-row :gutter="20" class="bottom-row">
      <!-- 最近操作 -->
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>最近操作</span>
              <el-button type="primary" size="small" text @click="goTo('/iam/operation-logs')">查看全部</el-button>
            </div>
          </template>
          <el-table :data="recentOperations" size="small" :show-header="false">
            <el-table-column prop="action" width="120">
              <template #default="{ row }">
                <el-tag size="small" :type="row.actionType">{{ row.action }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="resource" />
            <el-table-column prop="user" width="100" />
            <el-table-column prop="time" width="150" />
          </el-table>
        </el-card>
      </el-col>

      <!-- 服务状态 -->
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>云服务状态</span>
              <el-button type="primary" size="small" text @click="refreshServices">刷新</el-button>
            </div>
          </template>
          <div class="service-list">
            <div class="service-item" v-for="service in services" :key="service.provider">
              <span class="service-provider">{{ service.provider }}</span>
              <el-tag :type="service.status === 'normal' ? 'success' : 'danger'" size="small">
                {{ service.status === 'normal' ? '正常' : '异常' }}
              </el-tag>
              <span class="service-region">{{ service.region }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, OfficeBuilding, Folder, Cloudy, Monitor, Bell, Money, CircleCheck } from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()

// 统计数据
const stats = ref({
  userCount: 0,
  domainCount: 0,
  projectCount: 0,
  accountCount: 0
})

// 资源统计
const resourceStats = ref({
  vmCount: 0,
  diskCount: 0,
  eipCount: 0,
  sgCount: 0,
  dbCount: 0,
  lbCount: 0
})

// 告警列表
const alerts = ref([])

// 最近操作
const recentOperations = ref([
  { action: '创建', actionType: 'success', resource: '虚拟机 vm-test-01', user: 'admin', time: '2026-04-22 10:30' },
  { action: '修改', actionType: 'warning', resource: '用户 user-01', user: 'admin', time: '2026-04-22 09:20' },
  { action: '删除', actionType: 'danger', resource: '硬盘 disk-old', user: 'admin', time: '2026-04-22 08:15' }
])

// 服务状态
const services = ref([
  { provider: '阿里云', status: 'normal', region: '华东1(杭州)' },
  { provider: '腾讯云', status: 'normal', region: '广州' },
  { provider: 'AWS', status: 'normal', region: 'us-east-1' }
])

// 加载统计数据
const loadStats = async () => {
  try {
    const res = await request.get('/auth/stats')
    stats.value = {
      userCount: res.user_count || 0,
      domainCount: res.domain_count || 0,
      projectCount: res.project_count || 0,
      accountCount: 0
    }

    // 获取云账号数量
    try {
      const accountsRes = await request.get('/cloud-accounts')
      stats.value.accountCount = accountsRes.items?.length || accountsRes.data?.length || 0
    } catch (e) {
      // ignore
    }
  } catch (e) {
    // ignore
  }
}

// 加载资源统计
const loadResourceStats = async () => {
  try {
    // 获取虚拟机数量
    const vmsRes = await request.get('/compute/vms')
    resourceStats.value.vmCount = vmsRes.items?.length || 0

    // 获取硬盘数量
    const disksRes = await request.get('/storage/disks')
    resourceStats.value.diskCount = disksRes.items?.length || 0

    // 获取EIP数量
    const eipsRes = await request.get('/network/eips')
    resourceStats.value.eipCount = eipsRes.items?.length || 0

    // 获取安全组数量
    const sgsRes = await request.get('/compute/security-groups')
    resourceStats.value.sgCount = sgsRes.items?.length || 0

    // 获取数据库数量
    const dbsRes = await request.get('/database/rds')
    resourceStats.value.dbCount = dbsRes.items?.length || 0

    // 获取负载均衡数量
    const lbsRes = await request.get('/network/load-balancers')
    resourceStats.value.lbCount = lbsRes.items?.length || 0
  } catch (e) {
    // ignore
  }
}

// 刷新资源
const refreshResources = () => {
  loadResourceStats()
}

// 刷新服务
const refreshServices = () => {
  // TODO: 实现服务状态刷新
}

// 获取告警级别类型
const getAlertLevelType = (level: string) => {
  switch (level) {
    case '严重':
      return 'danger'
    case '警告':
      return 'warning'
    case '信息':
      return 'info'
    default:
      return ''
  }
}

// 跳转到页面
const goTo = (path: string) => {
  router.push(path)
}

onMounted(() => {
  loadStats()
  loadResourceStats()
})
</script>

<style scoped>
.dashboard-page {
  padding: 20px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  height: 100px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 28px;
}

.stat-icon.users {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.domains {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.projects {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.accounts {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.middle-row {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.quick-links-card {
  height: 100%;
}

.quick-links {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
}

.quick-link-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  background: #f5f7fa;
}

.quick-link-item:hover {
  background: #ecf5ff;
  color: #409eff;
}

.quick-link-item .el-icon {
  font-size: 20px;
}

.resource-overview-card {
  height: 100%;
}

.resource-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.resource-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 4px;
}

.resource-name {
  color: #606266;
}

.resource-count {
  font-weight: bold;
  color: #303133;
}

.alerts-card {
  height: 100%;
}

.alert-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.alert-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 4px;
}

.alert-title {
  flex: 1;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
}

.alert-time {
  color: #909399;
  font-size: 12px;
}

.no-alerts {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 150px;
  color: #67c23a;
  gap: 10px;
}

.no-alerts .el-icon {
  font-size: 40px;
}

.bottom-row {
  margin-bottom: 20px;
}

.service-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.service-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 4px;
}

.service-provider {
  font-weight: bold;
  color: #303133;
}

.service-region {
  flex: 1;
  color: #909399;
  text-align: right;
}
</style>