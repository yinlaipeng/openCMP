<template>
  <div class="anomaly-container">
    <div class="page-header">
      <h2>异常监测</h2>
      <div class="toolbar">
        <el-button type="primary" @click="handleRefresh">刷新</el-button>
      </div>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" @submit.prevent="loadAnomalies">
        <el-form-item label="云账号">
          <el-select v-model="selectedAccountId" placeholder="选择云账号" clearable style="width: 200px">
            <el-option v-for="account in cloudAccounts" :key="account.id" :label="account.name" :value="account.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="严重程度">
          <el-select v-model="selectedSeverity" placeholder="严重程度" clearable style="width: 120px">
            <el-option label="高" value="high" />
            <el-option label="中" value="medium" />
            <el-option label="低" value="low" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="selectedStatus" placeholder="状态" clearable style="width: 120px">
            <el-option label="新发现" value="new" />
            <el-option label="已确认" value="confirmed" />
            <el-option label="已解决" value="resolved" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadAnomalies">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="16" style="margin-bottom: 16px;">
      <el-col :span="6">
        <el-statistic title="异常总数" :value="pagination.total" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="高严重" :value="highSeverityCount">
          <template #suffix>
            <el-tag type="danger" size="small">高</el-tag>
          </template>
        </el-statistic>
      </el-col>
      <el-col :span="6">
        <el-statistic title="待处理" :value="pendingCount" />
      </el-col>
    </el-row>

    <!-- 数据表格 -->
    <el-table :data="anomalies" v-loading="loading" style="width: 100%" row-key="id">
      <el-table-column prop="anomaly_type" label="异常类型" width="120">
        <template #default="{ row }">
          <el-tag>{{ getAnomalyTypeText(row.anomaly_type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="period" label="发生周期" width="100" />
      <el-table-column prop="detected_at" label="检测时间" width="160">
        <template #default="{ row }">
          {{ formatDate(row.detected_at) }}
        </template>
      </el-table-column>
      <el-table-column prop="expected_cost" label="预期费用" width="100">
        <template #default="{ row }">
          ¥{{ row.expected_cost?.toFixed(2) || '0.00' }}
        </template>
      </el-table-column>
      <el-table-column prop="actual_cost" label="实际费用" width="100">
        <template #default="{ row }">
          ¥{{ row.actual_cost?.toFixed(2) || '0.00' }}
        </template>
      </el-table-column>
      <el-table-column prop="deviation_rate" label="偏差率" width="100">
        <template #default="{ row }">
          <span :style="{ color: row.deviation_rate >= 50 ? '#F56C6C' : '#E6A23C' }">
            {{ row.deviation_rate?.toFixed(1) }}%
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="severity" label="严重程度" width="80">
        <template #default="{ row }">
          <el-tag :type="getSeverityType(row.severity)">
            {{ getSeverityText(row.severity) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="handleResolve(row)" :disabled="row.status === 'resolved'">
            处理
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="pagination.currentPage"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      layout="total, sizes, prev, pager, next"
      class="pagination"
    />

    <!-- 处理异常对话框 -->
    <el-dialog v-model="resolveDialogVisible" title="处理异常" width="400px">
      <el-form :model="resolveForm" label-width="80px">
        <el-form-item label="异常信息">
          <el-descriptions :column="1" border size="small">
            <el-descriptions-item label="类型">{{ getAnomalyTypeText(currentAnomaly?.anomaly_type) }}</el-descriptions-item>
            <el-descriptions-item label="偏差率">{{ currentAnomaly?.deviation_rate?.toFixed(1) }}%</el-descriptions-item>
            <el-descriptions-item label="实际费用">¥{{ currentAnomaly?.actual_cost?.toFixed(2) }}</el-descriptions-item>
          </el-descriptions>
        </el-form-item>
        <el-form-item label="处理说明" required>
          <el-input v-model="resolveForm.resolution" type="textarea" :rows="3" placeholder="请输入处理说明" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resolveDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitResolve" :loading="submitting">确认处理</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getAnomalies, resolveAnomaly } from '@/api/finance'
import { getCloudAccounts } from '@/api/cloud-account'
import type { CostAnomaly } from '@/types/finance'

const anomalies = ref<CostAnomaly[]>([])
const loading = ref(false)
const cloudAccounts = ref<any[]>([])
const selectedAccountId = ref<number | undefined>()
const selectedSeverity = ref<string>('')
const selectedStatus = ref<string>('')
const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

const resolveDialogVisible = ref(false)
const submitting = ref(false)
const currentAnomaly = ref<CostAnomaly | null>(null)
const resolveForm = reactive({
  resolution: ''
})

const highSeverityCount = computed(() => {
  return anomalies.value.filter(a => a.severity === 'high').length
})

const pendingCount = computed(() => {
  return anomalies.value.filter(a => a.status !== 'resolved').length
})

const getAnomalyTypeText = (type: string) => {
  const map: Record<string, string> = {
    spike: '费用突增',
    drop: '费用骤降',
    unusual_pattern: '异常模式'
  }
  return map[type] || type
}

const getSeverityType = (severity: string) => {
  const map: Record<string, string> = {
    high: 'danger',
    medium: 'warning',
    low: 'info'
  }
  return map[severity] || 'default'
}

const getSeverityText = (severity: string) => {
  const map: Record<string, string> = {
    high: '高',
    medium: '中',
    low: '低'
  }
  return map[severity] || severity
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    new: 'danger',
    confirmed: 'warning',
    resolved: 'success'
  }
  return map[status] || 'default'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    new: '新发现',
    confirmed: '已确认',
    resolved: '已解决'
  }
  return map[status] || status
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const loadAnomalies = async () => {
  loading.value = true
  try {
    const res = await getAnomalies({
      cloud_account_id: selectedAccountId.value,
      severity: selectedSeverity.value,
      status: selectedStatus.value,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    anomalies.value = res.items || []
    pagination.total = res.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const loadCloudAccounts = async () => {
  try {
    const res = await getCloudAccounts({ page: 1, page_size: 100 })
    cloudAccounts.value = res.items || []
  } catch (e) {
    console.error(e)
  }
}

const resetFilters = () => {
  selectedAccountId.value = undefined
  selectedSeverity.value = ''
  selectedStatus.value = ''
  pagination.currentPage = 1
  loadAnomalies()
}

const handleRefresh = () => {
  loadAnomalies()
  ElMessage.success('数据已刷新')
}

const handleResolve = (row: CostAnomaly) => {
  currentAnomaly.value = row
  resolveForm.resolution = ''
  resolveDialogVisible.value = true
}

const submitResolve = async () => {
  if (!resolveForm.resolution) {
    ElMessage.warning('请输入处理说明')
    return
  }
  if (!currentAnomaly.value) return

  submitting.value = true
  try {
    await resolveAnomaly(currentAnomaly.value.id, resolveForm.resolution)
    ElMessage.success('异常已处理')
    resolveDialogVisible.value = false
    loadAnomalies()
  } catch (e) {
    ElMessage.error('处理失败')
  } finally {
    submitting.value = false
  }
}

watch([selectedAccountId, selectedSeverity, selectedStatus, pagination.currentPage, pagination.pageSize], loadAnomalies)

onMounted(() => {
  loadCloudAccounts()
  loadAnomalies()
})
</script>

<style scoped>
.anomaly-container {
  padding: 20px;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.page-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}
.filter-card {
  margin-bottom: 20px;
}
.pagination {
  margin-top: 20px;
  text-align: right;
}
</style>