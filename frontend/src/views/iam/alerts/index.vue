<template>
  <div class="alerts-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">安全告警</span>
          <!-- 工具栏按钮 - 参考 Cloudpods -->
          <div class="toolbar">
            <el-button @click="handleViewMode">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button @click="handleRefresh">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <!-- 统计卡片 - 参考 Cloudpods 顶部统计 -->
      <el-row :gutter="16" class="stats-row">
        <el-col :span="3">
          <div class="stat-card critical" v-if="fatalCount > 0">
            <div class="stat-icon">
              <el-icon><WarningFilled /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ fatalCount }}</div>
              <div class="stat-label">致命</div>
            </div>
          </div>
        </el-col>
        <el-col :span="3">
          <div class="stat-card important">
            <div class="stat-icon">
              <el-icon><InfoFilled /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ importantCount }}</div>
              <div class="stat-label">重要</div>
            </div>
          </div>
        </el-col>
        <el-col :span="3">
          <div class="stat-card normal">
            <div class="stat-icon">
              <el-icon><CircleCheckFilled /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ normalCount }}</div>
              <div class="stat-label">普通</div>
            </div>
          </div>
        </el-col>
        <el-col :span="3">
          <div class="stat-card total">
            <div class="stat-icon">
              <el-icon><BellFilled /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ pagination.total }}</div>
              <div class="stat-label">总计</div>
            </div>
          </div>
        </el-col>
      </el-row>

      <!-- 表格 - 中文表头 -->
      <el-table
        :data="alerts"
        v-loading="loading"
        border
        stripe
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <!-- 标题 - 告警名称 -->
        <el-table-column label="标题" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">
              <span>{{ row.alert_name || row.title || row.name }}</span>
            </el-button>
            <div class="alert-desc" v-if="row.message">
              <small>{{ row.message.substring(0, 50) }}...</small>
            </div>
          </template>
        </el-table-column>
        <!-- 严重级别 -->
        <el-table-column label="严重级别" width="120">
          <template #default="{ row }">
            <el-tag :type="getLevelTag(row.level)" effect="plain">
              {{ getLevelName(row.level) }}
            </el-tag>
          </template>
        </el-table-column>
        <!-- 接收人 -->
        <el-table-column label="接收人" width="150">
          <template #default="{ row }">
            <span>{{ row.recipients || row.user_name || '-' }}</span>
          </template>
        </el-table-column>
        <!-- 触发时间 -->
        <el-table-column label="触发时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.trigger_time || row.created_at) }}
          </template>
        </el-table-column>
        <!-- 告警内容 -->
        <el-table-column label="内容" min-width="250" show-overflow-tooltip>
          <template #default="{ row }">
            <span>{{ row.message || row.content || getAlertContent(row) }}</span>
          </template>
        </el-table-column>
        <!-- 操作 -->
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="handleView(row)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.limit"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadAlerts"
        @current-change="loadAlerts"
        class="pagination"
      />
    </el-card>

    <!-- 详情弹窗 -->
    <el-dialog v-model="detailVisible" title="告警详情" width="700px">
      <div v-if="currentAlert" class="alert-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="告警ID">{{ currentAlert.alert_id || currentAlert.id }}</el-descriptions-item>
          <el-descriptions-item label="告警名称">{{ currentAlert.alert_name || currentAlert.title }}</el-descriptions-item>
          <el-descriptions-item label="严重级别">
            <el-tag :type="getLevelTag(currentAlert.level)">
              {{ getLevelName(currentAlert.level) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="告警状态">
            <el-tag :type="currentAlert.status === 'active' || currentAlert.alert_state === 'alerting' ? 'danger' : 'success'">
              {{ currentAlert.status === 'active' || currentAlert.alert_state === 'alerting' ? '告警中' : '已恢复' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="告警类型">{{ currentAlert.metric || currentAlert.type || '-' }}</el-descriptions-item>
          <el-descriptions-item label="资源类型">{{ currentAlert.res_type || currentAlert.type || '-' }}</el-descriptions-item>
          <el-descriptions-item label="资源名称">{{ currentAlert.res_name || currentAlert.name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="触发时间">{{ formatTime(currentAlert.trigger_time || currentAlert.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="告警内容" :span="2">{{ currentAlert.message || getAlertContent(currentAlert) }}</el-descriptions-item>
          <el-descriptions-item label="标签信息" :span="2">
            <div v-if="currentAlert.tags">
              <el-tag v-for="(value, key) in currentAlert.tags" :key="key" style="margin: 2px">
                {{ key }}: {{ value }}
              </el-tag>
            </div>
            <span v-else>-</span>
          </el-descriptions-item>
          <el-descriptions-item label="阈值信息" :span="2" v-if="currentAlert.threshold">
            {{ currentAlert.threshold }} {{ currentAlert.comparator || '' }} {{ currentAlert.value_str || currentAlert.value }}
          </el-descriptions-item>
        </el-descriptions>

        <!-- 告警规则 -->
        <div v-if="currentAlert.alert_rule" class="alert-rule-section">
          <h4>告警规则</h4>
          <el-table :data="currentAlert.alert_rule" border size="small">
            <el-table-column prop="metric" label="指标" />
            <el-table-column prop="comparator" label="比较符" width="80" />
            <el-table-column prop="threshold" label="阈值" width="100" />
            <el-table-column prop="period" label="周期" width="80" />
          </el-table>
        </div>
      </div>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
        <el-button
          v-if="currentAlert?.status === 'active'"
          type="success"
          @click="handleResolve(currentAlert)"
        >
          处理告警
        </el-button>
        <el-button
          v-if="currentAlert?.can_delete"
          type="danger"
          @click="handleDelete(currentAlert)"
        >
          删除
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Refresh, WarningFilled, InfoFilled, CircleCheckFilled, BellFilled } from '@element-plus/icons-vue'
import { getSecurityAlerts, getSecurityAlertStats, resolveSecurityAlert, deleteSecurityAlert } from '@/api/iam'

const loading = ref(false)
const detailVisible = ref(false)
const currentAlert = ref<any>(null)
const selectedAlerts = ref<any[]>([])

const pagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

const alerts = ref<any[]>([])
const stats = reactive({
  total: 0,
  active: 0,
  critical: 0,
  handled: 0
})

// 计算各级别告警数量 - 参考 Cloudpods 级别设计
const fatalCount = computed(() => alerts.value.filter(a => a.level === 'critical' || a.level === 'fatal').length)
const importantCount = computed(() => alerts.value.filter(a => a.level === 'important' || a.level === 'high').length)
const normalCount = computed(() => alerts.value.filter(a => a.level === 'normal' || a.level === 'medium' || a.level === 'low').length)

// 级别映射 - Cloudpods 级别: fatal, important, normal
const getLevelName = (level: string) => {
  const map: Record<string, string> = {
    fatal: '致命',
    critical: '致命',
    important: '重要',
    high: '重要',
    normal: '普通',
    medium: '普通',
    low: '普通'
  }
  return map[level] || level
}

const getLevelTag = (level: string) => {
  const map: Record<string, any> = {
    fatal: 'danger',
    critical: 'danger',
    important: 'warning',
    high: 'warning',
    normal: 'info',
    medium: 'info',
    low: 'info'
  }
  return map[level] || 'info'
}

// 获取告警内容
const getAlertContent = (alert: any) => {
  if (alert.data?.value_str) {
    return `${alert.metric}: 当前值 ${alert.data.value_str}`
  }
  if (alert.res_name) {
    return `资源 ${alert.res_name} (${alert.res_type}) 触发告警`
  }
  return '-'
}

const formatTime = (time: string) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

const handleSelectionChange = (selection: any[]) => {
  selectedAlerts.value = selection
}

const handleViewMode = () => {
  // 切换视图模式
  ElMessage.info('视图切换功能')
}

const loadAlerts = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit,
      details: 'true'
    }

    const res = await getSecurityAlerts(params)
    alerts.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载安全告警失败')
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const res = await getSecurityAlertStats()
    stats.total = res.total || 0
    stats.active = res.active || 0
    stats.critical = res.critical || 0
    stats.handled = res.handled || 0
  } catch (e: any) {
    console.error('加载统计失败', e)
  }
}

const handleView = (row: any) => {
  currentAlert.value = row
  detailVisible.value = true
}

const handleResolve = async (row: any) => {
  try {
    await resolveSecurityAlert(row.id)
    ElMessage.success('告警已处理')
    detailVisible.value = false
    loadAlerts()
    loadStats()
  } catch (e: any) {
    ElMessage.error(e.message || '处理失败')
  }
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定删除此告警?', '确认', { type: 'warning' })
    await deleteSecurityAlert(row.id)
    ElMessage.success('告警已删除')
    detailVisible.value = false
    loadAlerts()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const handleRefresh = () => {
  loadAlerts()
  loadStats()
  ElMessage.success('刷新成功')
}

onMounted(() => {
  loadStats()
  loadAlerts()
})
</script>

<style scoped>
.alerts-page {
  height: 100%;
  padding: 20px;
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

.toolbar {
  display: flex;
  gap: 8px;
}

/* 统计卡片样式 - 参考 Cloudpods */
.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-radius: 8px;
  background: #f5f7fa;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-card.critical {
  background: linear-gradient(135deg, #fff5f5 0%, #ffe0e0 100%);
  border-left: 4px solid #f56c6c;
}

.stat-card.important {
  background: linear-gradient(135deg, #fffaf0 0%, #ffe8d0 100%);
  border-left: 4px solid #e6a23c;
}

.stat-card.normal {
  background: linear-gradient(135deg, #f0f9ff 0%, #d0e8ff 100%);
  border-left: 4px solid #409eff;
}

.stat-card.total {
  background: linear-gradient(135deg, #f5f7fa 0%, #e5e7ea 100%);
  border-left: 4px solid #909399;
}

.stat-icon {
  font-size: 28px;
  margin-right: 12px;
}

.stat-card.critical .stat-icon {
  color: #f56c6c;
}

.stat-card.important .stat-icon {
  color: #e6a23c;
}

.stat-card.normal .stat-icon {
  color: #409eff;
}

.stat-card.total .stat-icon {
  color: #909399;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  line-height: 1.2;
}

.stat-label {
  font-size: 12px;
  color: #909399;
}

.alert-desc {
  color: #909399;
  margin-top: 4px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.alert-detail {
  padding: 10px 0;
}

.alert-rule-section {
  margin-top: 20px;
}

.alert-rule-section h4 {
  margin-bottom: 10px;
  font-size: 14px;
  color: #606266;
}
</style>