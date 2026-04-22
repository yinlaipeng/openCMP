<template>
  <div class="alerts-container">
    <div class="page-header">
      <h2>安全告警</h2>
      <div class="toolbar">
        <el-button @click="handleRefresh">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button @click="handleExport">
          <el-icon><Download /></el-icon>
          导出
        </el-button>
      </div>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" @submit.prevent="loadAlerts">
        <el-form-item label="搜索">
          <el-select v-model="searchField" placeholder="选择字段" style="width: 100px">
            <el-option label="标题" value="title" />
            <el-option label="级别" value="level" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input v-model="searchValue" placeholder="输入关键词" clearable style="width: 180px" @keyup.enter="loadAlerts">
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadAlerts">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card>
      <el-table
        :data="alerts"
        v-loading="loading"
        border
        stripe
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <!-- 标题 - 点击打开详情 -->
        <el-table-column label="标题" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">
              <span>{{ row.alert_name || row.title || row.name }}</span>
            </el-button>
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
        <el-table-column label="触发时间" width="180" sortable>
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
      </el-table>
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.limit"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next"
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
import { Refresh, Download, Setting, Search } from '@element-plus/icons-vue'
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

// 搜索
const searchField = ref('title')
const searchValue = ref('')
const searchPlaceholder = computed(() => {
  if (searchField.value === 'title') return '请输入标题搜索...'
  if (searchField.value === 'level') return '请选择严重级别...'
  return '请输入搜索关键词...'
})

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

const handleSearch = () => {
  pagination.page = 1
  loadAlerts()
}

const resetFilter = () => {
  searchField.value = 'title'
  searchValue.value = ''
  pagination.page = 1
  loadAlerts()
}

const handleExport = () => {
  ElMessage.info('导出功能')
}

const handleSettings = () => {
  ElMessage.info('设置功能')
}

const loadAlerts = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit,
      details: 'true'
    }

    // 添加搜索参数
    if (searchValue.value) {
      if (searchField.value === 'title') {
        params.title = searchValue.value
      } else if (searchField.value === 'level') {
        params.level = searchValue.value
      }
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
  ElMessage.success('刷新成功')
}

onMounted(() => {
  loadAlerts()
})
</script>

<style scoped>
.alerts-container { padding: 20px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.page-header h2 { margin: 0; font-size: 18px; font-weight: 600; }
.filter-card { margin-bottom: 20px; }
.pagination { margin-top: 20px; text-align: right; }
.alert-detail { padding: 10px 0; }
</style>