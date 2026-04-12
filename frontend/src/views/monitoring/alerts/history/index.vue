<template>
  <div class="alert-history-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">告警历史</span>
        </div>
      </template>

      <!-- Time Filter -->
      <div class="filter-bar">
        <el-select v-model="timeFilter" placeholder="选择时间范围" style="width: 150px;" @change="handleTimeFilterChange">
          <el-option label="近1小时" value="1h" />
          <el-option label="近6小时" value="6h" />
          <el-option label="近12小时" value="12h" />
          <el-option label="近24小时" value="24h" />
          <el-option label="近1周" value="1w" />
          <el-option label="近1月" value="1m" />
          <el-option label="全部" value="all" />
        </el-select>
        <el-date-picker
          v-if="timeFilter === 'custom'"
          v-model="customTimeRange"
          type="datetimerange"
          range-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          style="width: 350px;"
        />
      </div>

      <el-table :data="alertHistory" v-loading="loading">
        <el-table-column label="策略名称" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.policy_name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="trigger_time" label="触发时间" width="160" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_type" label="资源类型" width="100" />
        <el-table-column prop="detail" label="策略详情" width="150" />
        <el-table-column prop="level" label="告警级别" width="100">
          <template #default="{ row }">
            <el-tag :type="getAlertLevelType(row.level)" size="small">{{ row.level }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_count" label="资源数量" width="100" />
        <el-table-column prop="message_status" label="消息状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getMessageStatusType(row.message_status)">{{ row.message_status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" @click="viewDetail(row)">查看</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="告警历史详情" width="800px">
      <el-tabs v-model="detailTab">
        <el-tab-pane label="详情" name="detail">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="策略名称">{{ selectedHistory?.policy_name }}</el-descriptions-item>
            <el-descriptions-item label="资源类型">{{ selectedHistory?.resource_type }}</el-descriptions-item>
            <el-descriptions-item label="策略详情">{{ selectedHistory?.detail }}</el-descriptions-item>
            <el-descriptions-item label="告警级别">{{ selectedHistory?.level }}</el-descriptions-item>
            <el-descriptions-item label="触发时间">{{ selectedHistory?.trigger_time }}</el-descriptions-item>
            <el-descriptions-item label="资源数量">{{ selectedHistory?.resource_count }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ selectedHistory?.status }}</el-descriptions-item>
            <el-descriptions-item label="消息状态">{{ selectedHistory?.message_status }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="告警记录" name="records">
          <el-table :data="historyRecords" size="small">
            <el-table-column prop="resource_name" label="资源名称" width="180" />
            <el-table-column prop="trigger_time" label="触发时间" width="160" />
            <el-table-column prop="trigger_value" label="触发值" width="100" />
            <el-table-column prop="status" label="状态" width="100" />
            <el-table-column prop="ip" label="IP" width="140" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="操作日志" name="logs">
          <el-table :data="operationLogs" size="small">
            <el-table-column prop="time" label="操作时间" width="160" />
            <el-table-column prop="action" label="操作" width="120" />
            <el-table-column prop="operator" label="操作人" width="120" />
            <el-table-column prop="result" label="结果" width="100" />
          </el-table>
        </el-tab-pane>
      </el-tabs>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

interface AlertHistoryItem {
  id: string
  policy_name: string
  trigger_time: string
  status: string
  resource_type: string
  detail: string
  level: string
  resource_count: number
  message_status: string
}

interface HistoryRecord {
  resource_name: string
  trigger_time: string
  trigger_value: string
  status: string
  ip: string
}

interface OperationLog {
  time: string
  action: string
  operator: string
  result: string
}

const timeFilter = ref('24h')
const customTimeRange = ref([])
const alertHistory = ref<AlertHistoryItem[]>([])
const loading = ref(false)

const detailDialogVisible = ref(false)
const detailTab = ref('detail')
const selectedHistory = ref<AlertHistoryItem | null>(null)
const historyRecords = ref<HistoryRecord[]>([])
const operationLogs = ref<OperationLog[]>([])

const getStatusType = (status: string) => {
  switch (status) {
    case '已恢复':
      return 'success'
    case '进行中':
      return 'warning'
    case '已屏蔽':
      return 'info'
    default:
      return ''
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

const getMessageStatusType = (status: string) => {
  switch (status) {
    case '已发送':
      return 'success'
    case '未发送':
      return 'info'
    case '发送失败':
      return 'danger'
    default:
      return ''
  }
}

const loadAlertHistory = async () => {
  loading.value = true
  try {
    alertHistory.value = [
      { id: 'history-1', policy_name: 'CPU使用率告警', trigger_time: '2024-01-15 10:30:00', status: '进行中', resource_type: '虚拟机', detail: 'CPU>80%', level: '严重', resource_count: 3, message_status: '已发送' },
      { id: 'history-2', policy_name: '内存使用率告警', trigger_time: '2024-01-15 09:20:00', status: '已恢复', resource_type: '数据库', detail: '内存>90%', level: '警告', resource_count: 2, message_status: '已发送' },
      { id: 'history-3', policy_name: '磁盘使用率告警', trigger_time: '2024-01-15 08:15:00', status: '已屏蔽', resource_type: '虚拟机', detail: '磁盘>85%', level: '信息', resource_count: 1, message_status: '未发送' },
      { id: 'history-4', policy_name: '连接数告警', trigger_time: '2024-01-15 07:30:00', status: '已恢复', resource_type: '负载均衡', detail: '连接数>500', level: '警告', resource_count: 1, message_status: '已发送' }
    ]
  } catch (e) {
    console.error(e)
    alertHistory.value = []
  } finally {
    loading.value = false
  }
}

const loadHistoryRecords = async (historyId: string) => {
  historyRecords.value = [
    { resource_name: 'prod-web-01', trigger_time: '2024-01-15 10:30:00', trigger_value: '85%', status: '告警', ip: '192.168.1.10' },
    { resource_name: 'prod-web-02', trigger_time: '2024-01-15 10:25:00', trigger_value: '82%', status: '告警', ip: '192.168.1.11' },
    { resource_name: 'prod-web-03', trigger_time: '2024-01-15 10:20:00', trigger_value: '88%', status: '告警', ip: '192.168.1.12' }
  ]
}

const loadOperationLogs = async (historyId: string) => {
  operationLogs.value = [
    { time: '2024-01-15 10:30:00', action: '告警触发', operator: '系统', result: '成功' },
    { time: '2024-01-15 10:31:00', action: '消息发送', operator: '系统', result: '成功' },
    { time: '2024-01-15 10:35:00', action: '查看详情', operator: 'admin', result: '成功' }
  ]
}

const handleTimeFilterChange = () => {
  if (timeFilter.value === 'custom') {
    return
  }
  ElMessage.success(`查询时间范围: ${timeFilter.value}`)
  loadAlertHistory()
}

const viewDetail = (row: AlertHistoryItem) => {
  selectedHistory.value = row
  loadHistoryRecords(row.id)
  loadOperationLogs(row.id)
  detailTab.value = 'detail'
  detailDialogVisible.value = true
}

onMounted(() => {
  loadAlertHistory()
})
</script>

<style scoped>
.alert-history-page {
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

.filter-bar {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
}
</style>