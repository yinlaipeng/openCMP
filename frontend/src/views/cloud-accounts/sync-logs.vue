<template>
  <div class="sync-logs-page">
    <el-card class="header-card">
      <div class="header-content">
        <div class="title-section">
          <h2>同步日志</h2>
          <p class="subtitle">查看云账号资源同步历史记录</p>
        </div>
        <div class="filter-section">
          <el-select
            v-model="selectedAccountId"
            placeholder="选择云账号"
            clearable
            @change="handleAccountChange"
            style="width: 300px"
          >
            <el-option
              v-for="account in cloudAccounts"
              :key="account.id"
              :label="account.name"
              :value="account.id"
            />
          </el-select>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            @change="handleDateChange"
            style="width: 300px"
          />
        </div>
      </div>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row" v-if="syncStatistics">
      <el-col :span="4">
        <el-card class="stat-card">
          <div class="stat-content">
            <span class="stat-value">{{ syncStatistics.total_sync_count }}</span>
            <span class="stat-label">总同步次数</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stat-card success">
          <div class="stat-content">
            <span class="stat-value">{{ syncStatistics.success_count }}</span>
            <span class="stat-label">成功次数</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stat-card danger">
          <div class="stat-content">
            <span class="stat-value">{{ syncStatistics.failure_count }}</span>
            <span class="stat-label">失败次数</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stat-card info">
          <div class="stat-content">
            <span class="stat-value">{{ formatDuration(syncStatistics.avg_duration) }}</span>
            <span class="stat-label">平均耗时</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stat-card warning">
          <div class="stat-content">
            <span class="stat-value">{{ syncStatistics.total_new }}</span>
            <span class="stat-label">新增资源</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stat-card primary">
          <div class="stat-content">
            <span class="stat-value">{{ syncStatistics.total_updated }}</span>
            <span class="stat-label">更新资源</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 同步日志表格 -->
    <el-card class="table-card">
      <el-table
        v-loading="loading"
        :data="syncLogs"
        style="width: 100%"
        @row-click="showLogDetail"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="cloud_account_name" label="云账号" width="150">
          <template #default="{ row }">
            <el-tag>{{ row.cloud_account_name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sync_mode" label="同步模式" width="100">
          <template #default="{ row }">
            <el-tag :type="row.sync_mode === 'full' ? 'warning' : 'success'">
              {{ row.sync_mode === 'full' ? '全量' : '增量' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_type" label="资源类型" width="120">
          <template #default="{ row }">
            {{ getResourceTypeLabel(row.resource_type) }}
          </template>
        </el-table-column>
        <el-table-column prop="sync_start_time" label="开始时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.sync_start_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="sync_duration" label="耗时" width="100">
          <template #default="{ row }">
            {{ formatDuration(row.sync_duration) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="资源统计" width="200">
          <template #default="{ row }">
            <div class="resource-stats">
              <span class="stat-item new">新增: {{ row.new_count }}</span>
              <span class="stat-item updated">更新: {{ row.updated_count }}</span>
              <span class="stat-item deleted">删除: {{ row.deleted_count }}</span>
              <span class="stat-item error" v-if="row.error_count > 0">
                错误: {{ row.error_count }}
              </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="triggered_by" label="触发方式" width="100">
          <template #default="{ row }">
            {{ row.triggered_by === 'manual' ? '手动' : '定时' }}
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="100">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="showLogDetail(row)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="同步日志详情"
      width="700px"
      destroy-on-close
    >
      <el-descriptions :column="2" border v-if="selectedLog">
        <el-descriptions-item label="ID">{{ selectedLog.id }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedLog.cloud_account_name }}</el-descriptions-item>
        <el-descriptions-item label="同步模式">
          <el-tag :type="selectedLog.sync_mode === 'full' ? 'warning' : 'success'">
            {{ selectedLog.sync_mode === 'full' ? '全量同步' : '增量同步' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="资源类型">
          {{ getResourceTypeLabel(selectedLog.resource_type) }}
        </el-descriptions-item>
        <el-descriptions-item label="开始时间">
          {{ formatTime(selectedLog.sync_start_time) }}
        </el-descriptions-item>
        <el-descriptions-item label="结束时间">
          {{ selectedLog.sync_end_time ? formatTime(selectedLog.sync_end_time) : '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="同步耗时">
          {{ formatDuration(selectedLog.sync_duration) }}
        </el-descriptions-item>
        <el-descriptions-item label="触发方式">
          {{ selectedLog.triggered_by === 'manual' ? '手动触发' : '定时触发' }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(selectedLog.status)">
            {{ getStatusLabel(selectedLog.status) }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>

      <el-divider>资源统计</el-divider>

      <el-row :gutter="20">
        <el-col :span="6">
          <el-statistic title="新增资源" :value="selectedLog.new_count" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="更新资源" :value="selectedLog.updated_count" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="删除资源" :value="selectedLog.deleted_count" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="跳过资源" :value="selectedLog.skipped_count" />
        </el-col>
      </el-row>

      <el-divider v-if="selectedLog.error_message">错误信息</el-divider>

      <el-alert
        v-if="selectedLog.error_message"
        type="error"
        :closable="false"
        show-icon
      >
        <template #title>同步过程中发生错误</template>
        <pre class="error-message">{{ selectedLog.error_message }}</pre>
      </el-alert>

      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { getCloudAccounts } from '@/api/cloud-account'
import { getSyncLogs, getSyncStatistics, type SyncLog, type SyncStatistics } from '@/api/sync-log'
import type { CloudAccount } from '@/types'

// 数据
const loading = ref(false)
const cloudAccounts = ref<CloudAccount[]>([])
const selectedAccountId = ref<number | null>(null)
const dateRange = ref<[Date, Date] | null>(null)
const syncLogs = ref<SyncLog[]>([])
const syncStatistics = ref<SyncStatistics | null>(null)
const detailDialogVisible = ref(false)
const selectedLog = ref<SyncLog | null>(null)

// 加载云账号列表
const loadCloudAccounts = async () => {
  try {
    const res = await getCloudAccounts()
    cloudAccounts.value = res.items || []
    if (cloudAccounts.value.length > 0) {
      selectedAccountId.value = cloudAccounts.value[0].id
      loadSyncLogs()
      loadSyncStatistics()
    }
  } catch (error: any) {
    ElMessage.error('加载云账号列表失败: ' + error.message)
  }
}

// 加载同步日志
const loadSyncLogs = async () => {
  if (!selectedAccountId.value) return

  loading.value = true
  try {
    const res = await getSyncLogs(selectedAccountId.value, 50)
    syncLogs.value = res.items || []
  } catch (error: any) {
    ElMessage.error('加载同步日志失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 加载同步统计
const loadSyncStatistics = async () => {
  if (!selectedAccountId.value) return

  try {
    syncStatistics.value = await getSyncStatistics(selectedAccountId.value, 30)
  } catch (error: any) {
    console.error('加载同步统计失败:', error)
  }
}

// 处理账号变更
const handleAccountChange = () => {
  loadSyncLogs()
  loadSyncStatistics()
}

// 处理日期变更
const handleDateChange = () => {
  loadSyncLogs()
}

// 显示日志详情
const showLogDetail = (row: SyncLog) => {
  selectedLog.value = row
  detailDialogVisible.value = true
}

// 格式化时间
const formatTime = (time: string) => {
  if (!time) return '-'
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 格式化耗时
const formatDuration = (seconds: number) => {
  if (!seconds) return '0秒'
  if (seconds < 60) return `${seconds}秒`
  const minutes = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${minutes}分${secs}秒`
}

// 获取状态类型
const getStatusType = (status: string) => {
  switch (status) {
    case 'success':
      return 'success'
    case 'failed':
      return 'danger'
    case 'partial_failure':
      return 'warning'
    case 'running':
      return 'info'
    default:
      return 'info'
  }
}

// 获取状态标签
const getStatusLabel = (status: string) => {
  switch (status) {
    case 'success':
      return '成功'
    case 'failed':
      return '失败'
    case 'partial_failure':
      return '部分失败'
    case 'running':
      return '进行中'
    default:
      return status
  }
}

// 获取资源类型标签
const getResourceTypeLabel = (type: string) => {
  const typeMap: Record<string, string> = {
    'all': '全部',
    'vm': '虚拟机',
    'vpc': 'VPC',
    'subnet': '子网',
    'security_group': '安全组',
    'eip': '弹性公网IP',
    'disk': '云硬盘',
    'snapshot': '快照',
    'image': '镜像',
    'rds': 'RDS数据库',
    'redis': 'Redis缓存'
  }
  return typeMap[type] || type
}

// 初始化
onMounted(() => {
  loadCloudAccounts()
})
</script>

<style scoped>
.sync-logs-page {
  padding: 20px;
}

.header-card {
  margin-bottom: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-section h2 {
  margin: 0;
  font-size: 24px;
  color: #303133;
}

.title-section .subtitle {
  margin: 8px 0 0;
  font-size: 14px;
  color: #909399;
}

.filter-section {
  display: flex;
  gap: 20px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  text-align: center;
}

.stat-card .stat-content {
  padding: 10px;
}

.stat-card .stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
  display: block;
}

.stat-card .stat-label {
  font-size: 14px;
  color: #909399;
  display: block;
  margin-top: 5px;
}

.stat-card.success .stat-value {
  color: #67c23a;
}

.stat-card.danger .stat-value {
  color: #f56c6c;
}

.stat-card.warning .stat-value {
  color: #e6a23c;
}

.stat-card.primary .stat-value {
  color: #409eff;
}

.stat-card.info .stat-value {
  color: #909399;
}

.table-card {
  margin-bottom: 20px;
}

.resource-stats {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.resource-stats .stat-item {
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
}

.resource-stats .stat-item.new {
  background-color: #f0f9eb;
  color: #67c23a;
}

.resource-stats .stat-item.updated {
  background-color: #ecf5ff;
  color: #409eff;
}

.resource-stats .stat-item.deleted {
  background-color: #fef0f0;
  color: #f56c6c;
}

.resource-stats .stat-item.error {
  background-color: #faecd8;
  color: #e6a23c;
}

.error-message {
  margin: 10px 0;
  padding: 10px;
  background-color: #f5f5f5;
  border-radius: 4px;
  white-space: pre-wrap;
  word-break: break-all;
  font-size: 12px;
}
</style>