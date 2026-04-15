<template>
  <div class="operation-log-tab">
    <div class="toolbar">
      <el-select v-model="filterResult" placeholder="结果筛选" size="small" clearable style="width: 120px; margin-right: 8px">
        <el-option label="成功" value="success" />
        <el-option label="失败" value="failed" />
      </el-select>
      <el-date-picker v-model="filterTimeRange" type="daterange" size="small" range-separator="-" start-placeholder="开始时间" end-placeholder="结束时间" style="width: 220px; margin-right: 8px" />
      <el-button size="small" @click="loadOperationLogs" :loading="loading">刷新</el-button>
    </div>
    <el-table :data="operationLogs" v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="#ID" width="60" />
      <el-table-column prop="operation_time" label="操作时间" width="160">
        <template #default="{ row }">{{ formatTime(row.operation_time) }}</template>
      </el-table-column>
      <el-table-column prop="resource_name" label="资源名称" width="150" show-overflow-tooltip />
      <el-table-column prop="resource_type" label="资源类型" width="100" />
      <el-table-column prop="operation_type" label="操作类型" width="100" />
      <el-table-column prop="service_type" label="服务类型" width="100" />
      <el-table-column prop="risk_level" label="风险级别" width="80">
        <template #default="{ row }">
          <el-tag :type="getRiskLevelType(row.risk_level)" size="small">{{ getRiskLevelText(row.risk_level) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="result" label="结果" width="80">
        <template #default="{ row }">
          <el-tag :type="row.result === 'success' ? 'success' : 'danger'" size="small">{{ row.result === 'success' ? '成功' : '失败' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="operator" label="发起人" width="100" />
      <el-table-column prop="project_name" label="所属项目" width="120" />
      <el-table-column label="操作" width="80" fixed="right">
        <template #default="{ row }">
          <el-button size="small" link type="primary" @click="handleViewDetail(row)">查看</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      v-model:current-page="pagination.currentPage"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      :page-sizes="[10, 20, 50]"
      layout="total, sizes, prev, pager, next"
      style="margin-top: 16px"
    />

    <!-- 详情弹窗 -->
    <el-dialog v-model="showDetailDialog" title="操作日志详情" width="600px">
      <el-descriptions :column="1" border v-if="currentLog">
        <el-descriptions-item label="ID">#{{ currentLog.id }}</el-descriptions-item>
        <el-descriptions-item label="操作时间">{{ formatTime(currentLog.operation_time) }}</el-descriptions-item>
        <el-descriptions-item label="资源名称">{{ currentLog.resource_name }}</el-descriptions-item>
        <el-descriptions-item label="资源类型">{{ currentLog.resource_type }}</el-descriptions-item>
        <el-descriptions-item label="操作类型">{{ currentLog.operation_type }}</el-descriptions-item>
        <el-descriptions-item label="服务类型">{{ currentLog.service_type }}</el-descriptions-item>
        <el-descriptions-item label="风险级别">{{ getRiskLevelText(currentLog.risk_level) }}</el-descriptions-item>
        <el-descriptions-item label="事件类型">{{ currentLog.time_type || 'API调用' }}</el-descriptions-item>
        <el-descriptions-item label="结果">{{ currentLog.result === 'success' ? '成功' : '失败' }}</el-descriptions-item>
        <el-descriptions-item label="发起人">{{ currentLog.operator }}</el-descriptions-item>
        <el-descriptions-item label="所属项目">{{ currentLog.project_name || '无' }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getOperationLogsByAccount } from '@/api/cloud-account'

interface Props { accountId: number }
const props = defineProps<Props>()

const loading = ref(false)
const showDetailDialog = ref(false)
const currentLog = ref<any>(null)
const operationLogs = ref<any[]>([])
const filterResult = ref('')
const filterTimeRange = ref<[Date, Date] | null>(null)

const pagination = reactive({ currentPage: 1, pageSize: 20, total: 0 })

onMounted(() => { loadOperationLogs() })

async function loadOperationLogs() {
  loading.value = true
  try {
    const res = await getOperationLogsByAccount(props.accountId, {
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    operationLogs.value = res.items || []
    pagination.total = res.total || 0
  } catch (error) {
    ElMessage.warning('获取操作日志失败')
    // 使用模拟数据作为后备
    operationLogs.value = [
      { id: 1, operation_time: '2026-04-14T10:30:00', resource_name: 'VM-001', resource_type: '虚拟机', operation_type: '创建', service_type: 'compute', risk_level: 'low', result: 'success', operator: 'admin', project_id: 1 },
      { id: 2, operation_time: '2026-04-14T10:31:00', resource_name: 'VPC-001', resource_type: '网络', operation_type: '删除', service_type: 'network', risk_level: 'high', result: 'success', operator: 'ops', project_id: 2 }
    ]
    pagination.total = 2
  } finally { loading.value = false }
}

function formatTime(time: string): string {
  return new Date(time).toLocaleString('zh-CN')
}

function getRiskLevelType(level: string): string {
  const types: Record<string, string> = { low: 'success', medium: 'warning', high: 'danger' }
  return types[level] || 'info'
}

function getRiskLevelText(level: string): string {
  const texts: Record<string, string> = { low: '低', medium: '中', high: '高' }
  return texts[level] || level
}

function handleViewDetail(row: any) {
  currentLog.value = row
  showDetailDialog.value = true
}
</script>

<style scoped>
.operation-log-tab { padding: 10px; }
.toolbar { margin-bottom: 16px; display: flex; align-items: center; }
</style>