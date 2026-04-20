<template>
  <div class="operation-logs-container">
    <!-- 工具栏 -->
    <div class="toolbar">
      <el-button size="small" @click="handleRefresh" :loading="loading">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-input
        v-model="searchKeyword"
        placeholder="默认为名称搜索，自动匹配IP或ID搜索项，IP或ID多个搜索用英文竖线(|)隔开"
        clearable
        @keyup.enter="handleSearch"
        class="search-input"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-button type="primary" @click="handleSearch">查询</el-button>
    </div>

    <!-- 日志列表 -->
    <el-table
      :data="filteredLogs"
      v-loading="loading"
      row-key="id"
      border
      stripe
      class="log-table"
    >
      <el-table-column prop="id" label="#ID" width="80">
        <template #default="{ row }">
          <span class="mono-text">{{ row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="operation_time" label="操作时间" width="180" sortable>
        <template #default="{ row }">
          <span class="mono-text">{{ formatDate(row.operation_time) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="resource_name" label="资源名称" min-width="150" show-overflow-tooltip />
      <el-table-column prop="resource_type" label="资源类型" width="120" />
      <el-table-column prop="operation_type" label="操作类型" width="100" />
      <el-table-column prop="service_type" label="服务类型" width="100" />
      <el-table-column prop="risk_level" label="风险级别" width="100">
        <template #default="{ row }">
          <el-tag :type="getRiskLevelType(row.risk_level)" size="small">
            {{ row.risk_level }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="event_type" label="事件类型" width="100" />
      <el-table-column prop="result" label="结果" width="80">
        <template #default="{ row }">
          <span class="result-status" :class="row.result === '成功' ? 'success' : 'fail'">
            {{ row.result }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="operator" label="发起人" width="120" />
      <el-table-column prop="project" label="所属项目" width="120" show-overflow-tooltip />
      <el-table-column label="操作" width="80" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="handleViewLog(row)">
            查看
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- Footer -->
    <div class="table-footer" v-if="logs.length > 0">
      <span>没有更多了</span>
    </div>

    <!-- 空状态 -->
    <EmptyState
      v-if="!loading && logs.length === 0"
      icon="Document"
      title="暂无数据"
      description="该用户暂无操作日志记录"
    />

    <!-- 分页 -->
    <el-pagination
      v-if="logs.length > 0"
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.limit"
      :total="pagination.total"
      :page-sizes="[10, 20, 50]"
      layout="total, sizes, prev, pager, next"
      @size-change="loadLogs"
      @current-change="loadLogs"
      class="pagination"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, Search } from '@element-plus/icons-vue'
import type { OperationLog } from '@/types/iam'
import EmptyState from '@/components/common/EmptyState.vue'

// Props
const props = defineProps<{
  userId?: number
}>()

// Emits
const emit = defineEmits<{
  (e: 'view-log', log: OperationLog): void
}>()

// State
const loading = ref(false)
const logs = ref<OperationLog[]>([])
const searchKeyword = ref('')
const pagination = ref({
  page: 1,
  limit: 20,
  total: 0
})

// Computed
const filteredLogs = computed(() => {
  if (!searchKeyword.value) return logs.value
  const keyword = searchKeyword.value.toLowerCase()
  return logs.value.filter(log =>
    log.resource_name?.toLowerCase().includes(keyword) ||
    log.operation_type?.toLowerCase().includes(keyword) ||
    log.operator?.toLowerCase().includes(keyword)
  )
})

// Load data
watch(() => props.userId, (id) => {
  if (id) loadLogs()
}, { immediate: true })

const loadLogs = async () => {
  if (!props.userId) return
  loading.value = true
  try {
    // TODO: Replace with real API when backend implements GET /users/:id/logs
    // const res = await getUserOperationLogs(props.userId, {
    //   limit: pagination.value.limit,
    //   offset: (pagination.value.page - 1) * pagination.value.limit,
    //   keyword: searchKeyword.value
    // })
    // logs.value = res.items || []
    // pagination.value.total = res.total || 0

    // Mock data for now
    logs.value = [
      {
        id: 1,
        operation_time: new Date().toISOString(),
        resource_name: '用户配置',
        resource_type: 'User',
        operation_type: '更新',
        service_type: 'IAM',
        risk_level: '低',
        event_type: 'API',
        result: '成功',
        operator: 'admin',
        project: '默认项目'
      },
      {
        id: 2,
        operation_time: new Date(Date.now() - 3600000).toISOString(),
        resource_name: '角色分配',
        resource_type: 'Role',
        operation_type: '创建',
        service_type: 'IAM',
        risk_level: '中',
        event_type: 'Console',
        result: '成功',
        operator: 'admin',
        project: '项目A'
      }
    ]
    pagination.value.total = 2
  } catch (e: any) {
    ElMessage.error(e.message || '加载日志失败')
  } finally {
    loading.value = false
  }
}

// Handlers
const handleRefresh = () => loadLogs()

const handleSearch = () => {
  pagination.value.page = 1
  loadLogs()
}

const handleViewLog = (log: OperationLog) => {
  emit('view-log', log)
}

// Utils
const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const getRiskLevelType = (level: string) => {
  switch (level) {
    case '高': return 'danger'
    case '中': return 'warning'
    case '低': return 'success'
    default: return 'info'
  }
}
</script>

<style scoped>
.operation-logs-container {
  padding: var(--space-4, 16px);
}

.toolbar {
  display: flex;
  gap: var(--space-2, 8px);
  margin-bottom: var(--space-3, 12px);
}

.search-bar {
  display: flex;
  gap: var(--space-2, 8px);
  margin-bottom: var(--space-3, 12px);
}

.search-input {
  flex: 1;
  max-width: 400px;
}

.log-table {
  margin-bottom: var(--space-3, 12px);
}

.mono-text {
  font-family: var(--font-mono, 'Fira Code', monospace);
}

.result-status.success {
  color: var(--color-success, #22C55E);
}

.result-status.fail {
  color: var(--color-danger, #EF4444);
}

.table-footer {
  text-align: center;
  color: var(--color-muted, #64748B);
  padding: var(--space-3, 12px);
  font-size: var(--font-size-sm, 14px);
}

.pagination {
  margin-top: var(--space-3, 12px);
  display: flex;
  justify-content: flex-end;
}
</style>