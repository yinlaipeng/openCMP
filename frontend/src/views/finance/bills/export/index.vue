<template>
  <div class="bills-export-container">
    <div class="page-header">
      <h2>账单导出中心</h2>
    </div>

    <el-card>
      <el-tabs v-model="activeTab">
        <!-- 导出任务创建 -->
        <el-tab-pane label="创建导出" name="create">
          <el-form :model="exportForm" label-width="100px" style="max-width: 500px;">
            <el-form-item label="云账号">
              <el-select v-model="exportForm.cloud_account_id" placeholder="选择云账号" clearable style="width: 100%;">
                <el-option v-for="a in cloudAccounts" :key="a.id" :label="a.name" :value="a.id" />
              </el-select>
            </el-form-item>
            <el-form-item label="账单周期">
              <el-date-picker
                v-model="exportForm.billing_cycle"
                type="month"
                format="YYYY-MM"
                value-format="YYYY-MM"
                placeholder="选择月份"
                style="width: 100%;"
              />
            </el-form-item>
            <el-form-item label="导出格式">
              <el-radio-group v-model="exportForm.format">
                <el-radio label="csv">CSV</el-radio>
                <el-radio label="excel">Excel</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleExport" :loading="exporting">导出账单</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 导出历史 -->
        <el-tab-pane label="导出历史" name="history">
          <el-table :data="exportHistory" v-loading="loadingHistory" style="width: 100%" row-key="task_id">
            <el-table-column prop="task_id" label="任务ID" width="150" />
            <el-table-column prop="cloud_account_name" label="云账号" width="150" />
            <el-table-column prop="billing_cycle" label="账单周期" width="100" />
            <el-table-column prop="format" label="格式" width="80">
              <template #default="{ row }">
                <el-tag size="small">{{ row.format.toUpperCase() }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="160">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="file_size" label="文件大小" width="100">
              <template #default="{ row }">
                {{ row.file_size ? formatFileSize(row.file_size) : '-' }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="handleDownload(row)" :disabled="row.status !== 'completed'">
                  下载
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { exportBills } from '@/api/finance'
import { getCloudAccounts } from '@/api/cloud-account'

const activeTab = ref('create')
const cloudAccounts = ref<any[]>([])
const exporting = ref(false)
const loadingHistory = ref(false)
const exportHistory = ref<any[]>([])

const exportForm = reactive({
  cloud_account_id: undefined as number | undefined,
  billing_cycle: '',
  format: 'excel'
})

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    processing: 'warning',
    completed: 'success',
    failed: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待处理',
    processing: '处理中',
    completed: '已完成',
    failed: '失败'
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

const formatFileSize = (size: number) => {
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(2) + ' KB'
  return (size / 1024 / 1024).toFixed(2) + ' MB'
}

const loadCloudAccounts = async () => {
  try {
    const res = await getCloudAccounts({ page: 1, page_size: 100 })
    cloudAccounts.value = res.items || []
  } catch (e) {
    console.error(e)
  }
}

const handleExport = async () => {
  exporting.value = true
  try {
    const res = await exportBills({
      cloud_account_id: exportForm.cloud_account_id,
      billing_cycle: exportForm.billing_cycle,
      format: exportForm.format
    })
    ElMessage.success('导出任务已创建')
    // 添加到历史记录
    exportHistory.value.unshift({
      task_id: 'task-' + Date.now(),
      cloud_account_name: cloudAccounts.value.find(a => a.id === exportForm.cloud_account_id)?.name || '全部',
      billing_cycle: exportForm.billing_cycle,
      format: exportForm.format,
      created_at: new Date().toISOString(),
      status: 'pending',
      file_size: null
    })
    activeTab.value = 'history'
  } catch (e) {
    ElMessage.error('导出失败')
  } finally {
    exporting.value = false
  }
}

const handleDownload = (row: any) => {
  ElMessage.info('下载功能开发中...')
}

onMounted(() => {
  loadCloudAccounts()
})
</script>

<style scoped>
.bills-export-container {
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
</style>