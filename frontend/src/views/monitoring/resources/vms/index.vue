<template>
  <div class="monitoring-vms-container">
    <div class="page-header">
      <h2>虚拟机监控</h2>
      <el-button type="primary" @click="handleSync">同步监控数据</el-button>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="loadVMList">
        <el-form-item label="云账号">
          <CloudAccountSelector
            v-model:value="filters.account_id"
            placeholder="选择云账号"
            @change="handleAccountChange"
          />
        </el-form-item>
        <el-form-item label="监控状态">
          <el-select v-model="filters.monitor_status" placeholder="监控状态" clearable>
            <el-option label="正常" value="正常" />
            <el-option label="告警" value="告警" />
            <el-option label="异常" value="异常" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadVMList">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table
      :data="vmList"
      v-loading="loading"
      style="width: 100%"
      row-key="resource_id"
    >
        <el-table-column label="名称" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.resource_name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="resource_id" label="资源ID" width="150" />
        <el-table-column prop="monitor_status" label="监控状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getMonitorStatusType(row.monitor_status)">{{ row.monitor_status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column prop="last_sync_at" label="同步时间" width="150">
          <template #default="{ row }">
            {{ formatTime(row.last_sync_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="viewMetrics(row)">查看指标</el-button>
            <el-button size="small" @click="handleManagePolicies(row)">管理告警策略</el-button>
          </template>
        </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      class="pagination"
    />

    <!-- VM Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="虚拟机详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="资源ID">{{ selectedVM?.resource_id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedVM?.resource_name }}</el-descriptions-item>
        <el-descriptions-item label="监控状态">{{ selectedVM?.monitor_status }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedVM?.platform }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedVM?.region }}</el-descriptions-item>
        <el-descriptions-item label="同步时间">{{ formatTime(selectedVM?.last_sync_at) }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Metrics Modal -->
    <el-dialog v-model="metricsDialogVisible" title="监控指标" width="700px">
      <div v-if="metricsLoading" class="metrics-loading">
        <el-skeleton :rows="5" animated />
      </div>
      <div v-else>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-card shadow="hover">
              <template #header>CPU使用率</template>
              <div class="metric-value">{{ currentMetrics?.cpu_usage?.value?.toFixed(1) || '-' }}%</div>
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card shadow="hover">
              <template #header>内存使用率</template>
              <div class="metric-value">{{ currentMetrics?.memory_usage?.value?.toFixed(1) || '-' }}%</div>
            </el-card>
          </el-col>
        </el-row>
        <el-row :gutter="20" style="margin-top: 20px;">
          <el-col :span="12">
            <el-card shadow="hover">
              <template #header>磁盘使用率</template>
              <div class="metric-value">{{ currentMetrics?.disk_usage?.value?.toFixed(1) || '-' }}%</div>
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card shadow="hover">
              <template #header>网络流量</template>
              <div class="metric-value">
                入: {{ currentMetrics?.network_in?.value?.toFixed(1) || '-' }} KB/s
                出: {{ currentMetrics?.network_out?.value?.toFixed(1) || '-' }} KB/s
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>
      <template #footer>
        <el-button @click="metricsDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Manage Alert Policies Modal -->
    <el-dialog v-model="policiesDialogVisible" title="管理告警策略" width="800px">
      <el-table :data="vmPolicies" v-loading="policiesLoading">
        <el-table-column prop="name" label="名称" width="180" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getPolicyStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="启用状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '已启用' : '已禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_type" label="资源类型" width="100">
          <template #default="{ row }">
            {{ getResourceTypeLabel(row.resource_type) }}
          </template>
        </el-table-column>
        <el-table-column label="策略详情" width="150">
          <template #default="{ row }">
            <span>{{ getMetricLabel(row.metric) }} > {{ row.threshold }}%</span>
          </template>
        </el-table-column>
        <el-table-column prop="level" label="告警级别" width="100">
          <template #default="{ row }">
            <el-tag :type="getAlertLevelType(row.level)" size="small">{{ row.level }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="owner" label="策略归属" width="120" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="handleEditPolicy(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDeletePolicy(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="policiesDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="handleAddPolicy">新增策略</el-button>
      </template>
    </el-dialog>

    <!-- 新增/编辑告警策略对话框 -->
    <el-dialog v-model="policyFormDialogVisible" :title="policyFormMode === 'add' ? '新增告警策略' : '编辑告警策略'" width="500px">
      <el-form :model="policyForm" label-width="100px">
        <el-form-item label="策略名称" required>
          <el-input v-model="policyForm.name" placeholder="请输入策略名称" />
        </el-form-item>
        <el-form-item label="监控指标" required>
          <el-select v-model="policyForm.metric" placeholder="请选择监控指标" style="width: 100%;">
            <el-option label="CPU使用率" value="cpu_usage" />
            <el-option label="内存使用率" value="memory_usage" />
            <el-option label="磁盘使用率" value="disk_usage" />
            <el-option label="网络流量" value="network_traffic" />
          </el-select>
        </el-form-item>
        <el-form-item label="阈值" required>
          <el-input-number v-model="policyForm.threshold" :min="0" :max="100" style="width: 150px;" />
          <span style="margin-left: 10px;">%</span>
        </el-form-item>
        <el-form-item label="持续时间" required>
          <el-input-number v-model="policyForm.duration" :min="1" :max="60" style="width: 150px;" />
          <span style="margin-left: 10px;">分钟</span>
        </el-form-item>
        <el-form-item label="告警级别" required>
          <el-select v-model="policyForm.level" placeholder="请选择告警级别" style="width: 100%;">
            <el-option label="信息" value="信息" />
            <el-option label="警告" value="警告" />
            <el-option label="严重" value="严重" />
          </el-select>
        </el-form-item>
        <el-form-item label="启用状态">
          <el-switch v-model="policyForm.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="policyFormDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handlePolicyFormSubmit" :loading="policyFormLoading">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import { listMonitorResources, syncMonitorResources, getResourceMetrics, type MonitorResource, type ResourceMetrics } from '@/api/monitor'
import { listAlertPolicies, createAlertPolicy, updateAlertPolicy, deleteAlertPolicy, type AlertPolicy, type AlertPolicyRequest } from '@/api/monitor'

interface VMPolicy extends AlertPolicy {}

const vmList = ref<MonitorResource[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedVM = ref<MonitorResource | null>(null)

// 筛选条件
const filters = reactive({
  account_id: null as number | null,
  monitor_status: ''
})

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const policiesDialogVisible = ref(false)
const vmPolicies = ref<AlertPolicy[]>([])
const policiesLoading = ref(false)

// Metrics dialog
const metricsDialogVisible = ref(false)
const metricsLoading = ref(false)
const currentMetrics = ref<ResourceMetrics | null>(null)

// 告警策略表单
const policyFormDialogVisible = ref(false)
const policyFormMode = ref<'add' | 'edit'>('add')
const policyFormLoading = ref(false)
const policyForm = ref<AlertPolicyRequest>({
  name: '',
  resource_type: 'vm',
  metric: 'cpu_usage',
  threshold: 80,
  duration: 5,
  level: '警告',
  enabled: true,
  owner: '自定义'
})

const getMonitorStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case '正常':
    case 'healthy':
      return 'success'
    case '告警':
    case 'warning':
      return 'warning'
    case '异常':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

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
    'network_traffic': '网络流量'
  }
  return labels[metric] || metric
}

const formatTime = (time: string | undefined) => {
  if (!time) return '-'
  return time.substring(0, 19).replace('T', ' ')
}

const handleAccountChange = (accountId: number | null) => {
  filters.account_id = accountId
  if (accountId) {
    loadVMList()
  } else {
    vmList.value = []
  }
}

const loadVMList = async () => {
  if (!filters.account_id) {
    vmList.value = []
    return
  }

  loading.value = true
  try {
    const res = await listMonitorResources(filters.account_id, 'vm', filters.monitor_status || undefined)
    vmList.value = res.items || res
    pagination.total = res.total || vmList.value.length
  } catch (e: any) {
    console.error(e)
    ElMessage.error(e.message || '获取监控资源列表失败')
    vmList.value = []
  } finally {
    loading.value = false
  }
}

const handleSync = async () => {
  if (!filters.account_id) {
    ElMessage.warning('请先选择云账号')
    return
  }

  try {
    const result = await syncMonitorResources(filters.account_id)
    ElMessage.success(result.message + `，同步了 ${result.count} 个资源`)
    loadVMList()
  } catch (e: any) {
    ElMessage.error(e.message || '同步失败')
  }
}

const viewDetail = (row: MonitorResource) => {
  selectedVM.value = row
  detailDialogVisible.value = true
}

const viewMetrics = async (row: MonitorResource) => {
  if (!filters.account_id) {
    return
  }

  metricsDialogVisible.value = true
  metricsLoading.value = true
  try {
    currentMetrics.value = await getResourceMetrics(filters.account_id, row.resource_id)
  } catch (e: any) {
    ElMessage.error(e.message || '获取指标失败')
    currentMetrics.value = null
  } finally {
    metricsLoading.value = false
  }
}

const loadVMPolicies = async () => {
  policiesLoading.value = true
  try {
    const policies = await listAlertPolicies({ resource_type: 'vm' })
    vmPolicies.value = policies
  } catch (e: any) {
    console.error(e)
    vmPolicies.value = []
  } finally {
    policiesLoading.value = false
  }
}

const handleManagePolicies = (row: MonitorResource) => {
  selectedVM.value = row
  loadVMPolicies()
  policiesDialogVisible.value = true
}

const handleEditPolicy = (row: AlertPolicy) => {
  policyFormMode.value = 'edit'
  policyForm.value = {
    name: row.name,
    resource_type: row.resource_type,
    metric: row.metric,
    threshold: row.threshold,
    duration: row.duration,
    level: row.level,
    enabled: row.enabled,
    owner: row.owner
  }
  policyFormDialogVisible.value = true
}

const handleDeletePolicy = async (row: AlertPolicy) => {
  try {
    await ElMessageBox.confirm(`确认删除策略 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await deleteAlertPolicy(row.id)
    vmPolicies.value = vmPolicies.value.filter(p => p.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const handleAddPolicy = () => {
  policyFormMode.value = 'add'
  policyForm.value = {
    name: '',
    resource_type: 'vm',
    metric: 'cpu_usage',
    threshold: 80,
    duration: 5,
    level: '警告',
    enabled: true,
    owner: '自定义'
  }
  policyFormDialogVisible.value = true
}

const handlePolicyFormSubmit = async () => {
  if (!policyForm.value.name) {
    ElMessage.warning('请输入策略名称')
    return
  }

  policyFormLoading.value = true
  try {
    if (policyFormMode.value === 'add') {
      const newPolicy = await createAlertPolicy(policyForm.value)
      vmPolicies.value.push(newPolicy)
      ElMessage.success('策略添加成功')
    } else {
      const currentPolicy = vmPolicies.value.find(p => p.name === policyForm.value.name)
      if (currentPolicy) {
        const updated = await updateAlertPolicy(currentPolicy.id, policyForm.value)
        const index = vmPolicies.value.findIndex(p => p.id === currentPolicy.id)
        if (index !== -1) {
          vmPolicies.value[index] = updated
        }
      }
      ElMessage.success('策略更新成功')
    }
    policyFormDialogVisible.value = false
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    policyFormLoading.value = false
  }
}

onMounted(() => {
  // Wait for account selector to initialize
})

const resetFilters = () => {
  filters.account_id = null
  filters.monitor_status = ''
  pagination.page = 1
  vmList.value = []
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadVMList()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadVMList()
}
</script>

<style scoped>
.monitoring-vms-container {
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

.metric-value {
  font-size: 24px;
  font-weight: bold;
  text-align: center;
  padding: 20px 0;
}

.metrics-loading {
  padding: 20px;
}
</style>