<template>
  <el-table :data="policies" v-loading="loading">
    <el-table-column prop="name" label="名称" width="180" />
    <el-table-column prop="status" label="状态" width="100">
      <template #default="{ row }">
        <el-tag :type="getPolicyStatusType(row.status)">{{ row.status }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="enabled" label="启用状态" width="100">
      <template #default="{ row }">
        <el-switch
          :model-value="row.enabled"
          @change="$emit('toggle', row)"
          size="small"
        />
      </template>
    </el-table-column>
    <el-table-column prop="resource_type" label="资源类型" width="100">
      <template #default="{ row }">
        {{ getResourceTypeLabel(row.resource_type) }}
      </template>
    </el-table-column>
    <el-table-column label="策略详情" width="200">
      <template #default="{ row }">
        <span>{{ getMetricLabel(row.metric) }} > {{ row.threshold }}% 持续 {{ row.duration }}分钟</span>
      </template>
    </el-table-column>
    <el-table-column prop="level" label="告警级别" width="100">
      <template #default="{ row }">
        <el-tag :type="getAlertLevelType(row.level)" size="small">{{ row.level }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="owner" label="策略归属" width="120" />
    <el-table-column prop="notify_channel" label="通知渠道" width="120">
      <template #default="{ row }">
        {{ row.notify_channel || '-' }}
      </template>
    </el-table-column>
    <el-table-column prop="created_at" label="创建时间" width="150">
      <template #default="{ row }">
        {{ formatTime(row.created_at) }}
      </template>
    </el-table-column>
    <el-table-column label="操作" width="150" fixed="right">
      <template #default="{ row }">
        <el-button size="small" @click="$emit('edit', row)">编辑</el-button>
        <el-button size="small" type="danger" @click="$emit('delete', row)">删除</el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
interface AlertPolicy {
  id: number
  name: string
  status: string
  enabled: boolean
  resource_type: string
  metric: string
  threshold: number
  duration: number
  level: string
  owner: string
  notify_channel: string
  created_at: string
}

defineProps<{
  policies: AlertPolicy[]
  loading?: boolean
}>()

defineEmits<{
  (e: 'edit', row: AlertPolicy): void
  (e: 'delete', row: AlertPolicy): void
  (e: 'toggle', row: AlertPolicy): void
}>()

const getPolicyStatusType = (status: string) => {
  switch (status) {
    case '正常':
      return 'success'
    case '异常':
      return 'danger'
    default:
      return 'info'
  }
}

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

const getResourceTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    'vm': '虚拟机',
    'database': '数据库',
    'network': '网络',
    'application': '应用'
  }
  return labels[type] || type
}

const getMetricLabel = (metric: string) => {
  const labels: Record<string, string> = {
    'cpu_usage': 'CPU',
    'memory_usage': '内存',
    'disk_usage': '磁盘',
    'network_traffic': '网络流量',
    'db_connections': '数据库连接数'
  }
  return labels[metric] || metric
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return time.substring(0, 19).replace('T', ' ')
}
</script>