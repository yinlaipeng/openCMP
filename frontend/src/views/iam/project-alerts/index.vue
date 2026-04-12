<template>
  <el-card class="section-card">
    <template #header>
      <div class="card-header">
        <span>安全告警</span>
      </div>
    </template>

    <el-table :data="alerts" v-loading="alertsLoading" style="width: 100%">
      <el-table-column prop="title" label="标题" width="200" />
      <el-table-column prop="level" label="级别" width="120">
        <template #default="{ row }">
          <el-tag :type="getLevelType(row.level)">{{ getLevelName(row.level) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="user_id" label="接受人" width="150">
        <template #default="{ row }">
          {{ getUserDisplayName(row.user_id) }}
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="触发时间" width="180" />
      <el-table-column prop="message" label="内容" />
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const alerts = ref<any[]>([])

// 加载状态
const alertsLoading = ref(false)

// 获取安全告警 - 使用模拟数据直到后端API完成
const loadAlerts = async () => {
  alertsLoading.value = true
  try {
    // 模拟API延迟
    await new Promise(resolve => setTimeout(resolve, 500))

    // 模拟安全告警数据
    alerts.value = [
      {
        id: 1,
        title: '登录失败次数过多',
        level: 'high',
        message: '用户 admin 在短时间内连续 5 次登录失败，可能存在暴力破解尝试',
        user_id: 1,
        created_at: new Date(Date.now() - 3600000).toISOString()
      },
      {
        id: 2,
        title: '异常登录行为',
        level: 'critical',
        message: '检测到账号在非常见地点登录，请确认是否为本人操作',
        user_id: 2,
        created_at: new Date(Date.now() - 1800000).toISOString()
      },
      {
        id: 3,
        title: '密码即将过期',
        level: 'medium',
        message: '您的账户密码将在 7 天后过期，请及时更新密码',
        user_id: 3,
        created_at: new Date(Date.now() - 86400000).toISOString()
      },
      {
        id: 4,
        title: 'MFA 验证失败',
        level: 'high',
        message: '检测到多次 MFA 验证失败，请检查双因素验证设置',
        user_id: 1,
        created_at: new Date(Date.now() - 120000).toISOString()
      },
      {
        id: 5,
        title: '权限变更警告',
        level: 'low',
        message: '您的账户权限发生变更，请注意安全',
        user_id: 4,
        created_at: new Date(Date.now() - 172800000).toISOString()
      }
    ]
  } catch (error) {
    console.error('Failed to load alerts:', error)
    ElMessage.error('加载安全告警失败')
  } finally {
    alertsLoading.value = false
  }
}

// 根据严重程度返回标签类型
const getLevelType = (level: string) => {
  switch (level.toLowerCase()) {
    case 'critical':
      return 'danger'
    case 'high':
      return 'warning'
    case 'medium':
      return 'info'
    case 'low':
      return 'success'
    default:
      return 'info'
  }
}

// 获取级别名称
const getLevelName = (level: string) => {
  const map: Record<string, string> = {
    low: '低',
    medium: '中',
    high: '高',
    critical: '严重'
  }
  return map[level] || level
}

// 获取用户显示名称
const getUserDisplayName = (userId: number | undefined) => {
  if (!userId) return '系统'
  // 在实际应用中，这里应该通过API获取用户名
  const userMap: Record<number, string> = {
    1: 'admin',
    2: 'user1',
    3: 'user2',
    4: 'user3'
  }
  return userMap[userId] || `用户${userId}`
}

onMounted(() => {
  loadAlerts()
})
</script>

<style scoped>
.project-context-notice {
  margin-bottom: 20px;
}

.project-alert {
  border-radius: 4px;
}

.section-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

:deep(.el-card__header) {
  padding: 12px 20px;
  background-color: #fafafa;
  border-bottom: 1px solid #eee;
}
</style>