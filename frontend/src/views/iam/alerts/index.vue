<template>
  <div class="alerts-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">安全告警</span>
          <el-button type="danger" @click="handleRefresh">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>
      
      <el-row :gutter="20" class="stats-row">
        <el-col :span="6">
          <el-statistic title="总告警数" :value="total" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="活跃告警" :value="active" value-style="color: #F56C6C" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="高危告警" :value="critical" value-style="color: #E6A23C" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="已处理" :value="handled" value-style="color: #67C23A" />
        </el-col>
      </el-row>
      
      <el-table :data="alerts" v-loading="loading" style="margin-top: 20px">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)">
              {{ getTypeName(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="level" label="级别" width="100">
          <template #default="{ row }">
            <el-tag :type="getLevelTag(row.level)">
              {{ getLevelName(row.level) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'danger' : 'success'">
              {{ row.status === 'active' ? '活跃' : '已处理' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="source_ip" label="来源 IP" width="140" />
        <el-table-column prop="created_at" label="时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">详情</el-button>
            <el-button 
              v-if="row.status === 'active'" 
              size="small" 
              type="success" 
              @click="handleResolve(row)"
            >
              处理
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <el-dialog v-model="detailVisible" title="告警详情" width="600px">
      <div v-if="currentAlert" class="alert-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="告警 ID">{{ currentAlert.id }}</el-descriptions-item>
          <el-descriptions-item label="级别">
            <el-tag :type="getLevelTag(currentAlert.level)">
              {{ getLevelName(currentAlert.level) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="类型">{{ getTypeName(currentAlert.type) }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="currentAlert.status === 'active' ? 'danger' : 'success'">
              {{ currentAlert.status === 'active' ? '活跃' : '已处理' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="标题" :span="2">{{ currentAlert.title }}</el-descriptions-item>
          <el-descriptions-item label="详情" :span="2">{{ currentAlert.message }}</el-descriptions-item>
          <el-descriptions-item label="来源 IP">{{ currentAlert.source_ip }}</el-descriptions-item>
          <el-descriptions-item label="时间">{{ formatTime(currentAlert.created_at) }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const detailVisible = ref(false)
const currentAlert = ref<any>(null)

const total = ref(0)
const active = ref(0)
const critical = ref(0)
const handled = ref(0)

const alerts = ref<any[]>([])

const getTypeName = (type: string) => {
  const map: Record<string, string> = {
    login_failed: '登录失败',
    password_expired: '密码过期',
    mfa_disabled: 'MFA 禁用',
    abnormal_login: '异常登录'
  }
  return map[type] || type
}

const getTypeTag = (type: string) => {
  const map: Record<string, any> = {
    login_failed: 'info',
    password_expired: 'warning',
    mfa_disabled: 'danger',
    abnormal_login: 'danger'
  }
  return map[type] || ''
}

const getLevelName = (level: string) => {
  const map: Record<string, string> = {
    low: '低',
    medium: '中',
    high: '高',
    critical: '严重'
  }
  return map[level] || level
}

const getLevelTag = (level: string) => {
  const map: Record<string, any> = {
    low: 'info',
    medium: 'warning',
    high: 'danger',
    critical: 'danger'
  }
  return map[level] || ''
}

const formatTime = (time: string) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

const loadAlerts = () => {
  loading.value = true
  // 模拟数据
  alerts.value = [
    {
      id: 1,
      type: 'login_failed',
      level: 'low',
      title: '登录失败次数过多',
      message: '用户 admin 连续 5 次登录失败',
      source_ip: '192.168.1.100',
      status: 'active',
      created_at: new Date().toISOString()
    },
    {
      id: 2,
      type: 'abnormal_login',
      level: 'high',
      title: '异地登录提醒',
      message: '检测到账号在非常见地点登录',
      source_ip: '114.114.114.114',
      status: 'active',
      created_at: new Date(Date.now() - 3600000).toISOString()
    },
    {
      id: 3,
      type: 'password_expired',
      level: 'medium',
      title: '密码即将过期',
      message: '您的密码将在 7 天后过期',
      source_ip: '192.168.1.50',
      status: 'handled',
      created_at: new Date(Date.now() - 86400000).toISOString()
    }
  ]
  
  total.value = alerts.value.length
  active.value = alerts.value.filter(a => a.status === 'active').length
  critical.value = alerts.value.filter(a => a.level === 'critical' || a.level === 'high').length
  handled.value = alerts.value.filter(a => a.status !== 'active').length
  loading.value = false
}

const handleView = (row: any) => {
  currentAlert.value = row
  detailVisible.value = true
}

const handleResolve = (row: any) => {
  row.status = 'handled'
  ElMessage.success('告警已处理')
  loadAlerts()
}

const handleRefresh = () => {
  loadAlerts()
  ElMessage.success('刷新成功')
}

onMounted(() => {
  loadAlerts()
})
</script>

<style scoped>
.alerts-page {
  height: 100%;
}

.page-card {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 18px;
  font-weight: bold;
}

.stats-row {
  margin-bottom: 20px;
}

.stats-row .el-statistic {
  text-align: center;
}

.alert-detail {
  padding: 10px 0;
}
</style>
