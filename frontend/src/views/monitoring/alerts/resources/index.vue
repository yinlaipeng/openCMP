<template>
  <div class="alert-resources-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">告警资源</span>
        </div>
      </template>

      <el-table :data="alertResources" v-loading="loading">
        <el-table-column label="资源名称" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.resource_name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="trigger_time" label="触发时间" width="160" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_type" label="资源类型" width="100" />
        <el-table-column prop="trigger_condition" label="触发条件" width="150" />
        <el-table-column prop="level" label="告警级别" width="100">
          <template #default="{ row }">
            <el-tag :type="getAlertLevelType(row.level)" size="small">{{ row.level }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP" width="140" />
        <el-table-column prop="policy_name" label="策略名称" width="150" />
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="trigger_value" label="触发值" width="100" />
        <el-table-column prop="message_status" label="消息状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getMessageStatusType(row.message_status)">{{ row.message_status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="handleViewMonitor(row)">查看监控</el-button>
            <el-button size="small" type="warning" @click="handleBlock(row)">屏蔽</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="告警资源详情" width="800px">
      <el-tabs v-model="detailTab">
        <el-tab-pane label="详情" name="detail">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="资源名称">{{ selectedResource?.resource_name }}</el-descriptions-item>
            <el-descriptions-item label="资源类型">{{ selectedResource?.resource_type }}</el-descriptions-item>
            <el-descriptions-item label="IP">{{ selectedResource?.ip }}</el-descriptions-item>
            <el-descriptions-item label="平台">{{ selectedResource?.platform }}</el-descriptions-item>
            <el-descriptions-item label="触发时间">{{ selectedResource?.trigger_time }}</el-descriptions-item>
            <el-descriptions-item label="告警级别">{{ selectedResource?.level }}</el-descriptions-item>
            <el-descriptions-item label="触发条件">{{ selectedResource?.trigger_condition }}</el-descriptions-item>
            <el-descriptions-item label="触发值">{{ selectedResource?.trigger_value }}</el-descriptions-item>
            <el-descriptions-item label="策略名称">{{ selectedResource?.policy_name }}</el-descriptions-item>
            <el-descriptions-item label="消息状态">{{ selectedResource?.message_status }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="告警记录" name="records">
          <el-table :data="alertRecords" size="small">
            <el-table-column prop="trigger_time" label="触发时间" width="160" />
            <el-table-column prop="level" label="告警级别" width="100" />
            <el-table-column prop="trigger_value" label="触发值" width="100" />
            <el-table-column prop="status" label="处理状态" width="100" />
            <el-table-column prop="handler" label="处理人" width="120" />
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
import { ElMessage, ElMessageBox } from 'element-plus'

interface AlertResource {
  id: string
  resource_name: string
  trigger_time: string
  status: string
  resource_type: string
  trigger_condition: string
  level: string
  ip: string
  policy_name: string
  platform: string
  trigger_value: string
  message_status: string
}

interface AlertRecord {
  trigger_time: string
  level: string
  trigger_value: string
  status: string
  handler: string
}

const alertResources = ref<AlertResource[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const detailTab = ref('detail')
const selectedResource = ref<AlertResource | null>(null)
const alertRecords = ref<AlertRecord[]>([])

const getStatusType = (status: string) => {
  switch (status) {
    case '未处理':
      return 'danger'
    case '已处理':
      return 'success'
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

const loadAlertResources = async () => {
  loading.value = true
  try {
    alertResources.value = [
      { id: 'alert-1', resource_name: 'prod-web-01', trigger_time: '2024-01-15 10:30:00', status: '未处理', resource_type: '虚拟机', trigger_condition: 'CPU>80%', level: '严重', ip: '192.168.1.10', policy_name: 'CPU使用率告警', platform: '阿里云', trigger_value: '85%', message_status: '已发送' },
      { id: 'alert-2', resource_name: 'prod-db-01', trigger_time: '2024-01-15 09:20:00', status: '已处理', resource_type: '数据库', trigger_condition: '内存>90%', level: '警告', ip: '192.168.2.10', policy_name: '内存使用率告警', platform: '阿里云', trigger_value: '92%', message_status: '已发送' },
      { id: 'alert-3', resource_name: 'dev-api-02', trigger_time: '2024-01-15 08:15:00', status: '已屏蔽', resource_type: '虚拟机', trigger_condition: '磁盘>85%', level: '信息', ip: '192.168.3.10', policy_name: '磁盘使用率告警', platform: '阿里云', trigger_value: '88%', message_status: '未发送' },
      { id: 'alert-4', resource_name: 'prod-lb-01', trigger_time: '2024-01-15 07:30:00', status: '未处理', resource_type: '负载均衡', trigger_condition: '连接数>500', level: '警告', ip: '192.168.4.10', policy_name: '连接数告警', platform: '阿里云', trigger_value: '520', message_status: '已发送' }
    ]
  } catch (e) {
    console.error(e)
    alertResources.value = []
  } finally {
    loading.value = false
  }
}

const loadAlertRecords = async (resourceId: string) => {
  alertRecords.value = [
    { trigger_time: '2024-01-15 10:30:00', level: '严重', trigger_value: '85%', status: '未处理', handler: '-' },
    { trigger_time: '2024-01-15 10:00:00', level: '警告', trigger_value: '82%', status: '已处理', handler: 'admin' },
    { trigger_time: '2024-01-15 09:30:00', level: '信息', trigger_value: '78%', status: '已处理', handler: 'admin' }
  ]
}

const viewDetail = (row: AlertResource) => {
  selectedResource.value = row
  loadAlertRecords(row.id)
  detailTab.value = 'detail'
  detailDialogVisible.value = true
}

const handleViewMonitor = (row: AlertResource) => {
  ElMessage.info(`查看监控: ${row.resource_name}`)
}

const handleBlock = async (row: AlertResource) => {
  try {
    await ElMessageBox.confirm(`确认屏蔽资源 ${row.resource_name} 的告警？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    row.status = '已屏蔽'
    ElMessage.success('屏蔽成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadAlertResources()
})
</script>

<style scoped>
.alert-resources-page {
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
</style>