<template>
  <div class="reports-container">
    <div class="page-header">
      <h2>成本报告</h2>
      <div class="toolbar">
        <el-button type="primary" @click="handleGenerate">生成报告</el-button>
      </div>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" @submit.prevent="loadReports">
        <el-form-item label="云账号">
          <el-select v-model="filterAccountId" placeholder="选择云账号" clearable style="width: 200px">
            <el-option v-for="account in cloudAccounts" :key="account.id" :label="account.name" :value="account.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="报告类型">
          <el-select v-model="filterReportType" placeholder="报告类型" clearable style="width: 150px">
            <el-option label="月度报告" value="monthly" />
            <el-option label="季度报告" value="quarterly" />
            <el-option label="年度报告" value="yearly" />
            <el-option label="自定义报告" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadReports">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table :data="reports" v-loading="loading" style="width: 100%" row-key="id">
      <el-table-column prop="report_id" label="报告ID" width="150" />
      <el-table-column prop="report_type" label="报告类型" width="120">
        <template #default="{ row }">
          <el-tag>{{ getReportTypeText(row.report_type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="period" label="报告周期" width="200">
        <template #default="{ row }">
          {{ row.start_date }} 至 {{ row.end_date }}
        </template>
      </el-table-column>
      <el-table-column prop="total_cost" label="总成本" width="120">
        <template #default="{ row }">
          ¥{{ row.total_cost?.toFixed(2) || '0.00' }}
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="生成时间" width="160">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 'completed' ? 'success' : 'warning'">
            {{ row.status === 'completed' ? '已完成' : '生成中' }}
          </el-tag>
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

    <el-pagination
      v-model:current-page="pagination.currentPage"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      layout="total, sizes, prev, pager, next"
      class="pagination"
    />

    <!-- 生成报告对话框 -->
    <el-dialog v-model="generateDialogVisible" title="生成成本报告" width="500px">
      <el-form :model="generateForm" label-width="100px">
        <el-form-item label="云账号">
          <el-select v-model="generateForm.cloud_account_id" placeholder="选择云账号（可选，留空为全局）" clearable style="width: 100%;">
            <el-option v-for="account in cloudAccounts" :key="account.id" :label="account.name" :value="account.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="报告类型">
          <el-select v-model="generateForm.report_type" placeholder="请选择报告类型" style="width: 100%;">
            <el-option label="月度报告" value="monthly" />
            <el-option label="季度报告" value="quarterly" />
            <el-option label="年度报告" value="yearly" />
            <el-option label="自定义报告" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="generateForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            style="width: 100%;"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="generateDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitGenerate" :loading="generating">生成</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getCostReports, generateCostReport } from '@/api/finance'
import { getCloudAccounts } from '@/api/cloud-account'

const reports = ref<any[]>([])
const loading = ref(false)
const cloudAccounts = ref<any[]>([])
const filterAccountId = ref<number | undefined>()
const filterReportType = ref<string>('')

const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

const generateDialogVisible = ref(false)
const generating = ref(false)
const generateForm = reactive({
  cloud_account_id: undefined as number | undefined,
  report_type: 'monthly',
  dateRange: [] as string[]
})

const getReportTypeText = (type: string) => {
  const map: Record<string, string> = {
    monthly: '月度',
    quarterly: '季度',
    yearly: '年度',
    custom: '自定义'
  }
  return map[type] || type
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

const loadReports = async () => {
  loading.value = true
  try {
    const res = await getCostReports({
      cloud_account_id: filterAccountId.value,
      report_type: filterReportType.value,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    reports.value = res.items || []
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
  filterAccountId.value = undefined
  filterReportType.value = ''
  pagination.currentPage = 1
  loadReports()
}

const handleGenerate = () => {
  generateForm.cloud_account_id = undefined
  generateForm.report_type = 'monthly'
  generateForm.dateRange = []
  generateDialogVisible.value = true
}

const submitGenerate = async () => {
  generating.value = true
  try {
    const params: any = {
      report_type: generateForm.report_type
    }
    if (generateForm.cloud_account_id) {
      params.cloud_account_id = generateForm.cloud_account_id
    }
    if (generateForm.dateRange && generateForm.dateRange.length === 2) {
      params.start_date = generateForm.dateRange[0]
      params.end_date = generateForm.dateRange[1]
    }
    await generateCostReport(params)
    ElMessage.success('报告生成任务已提交')
    generateDialogVisible.value = false
    loadReports()
  } catch (e) {
    ElMessage.error('生成失败')
  } finally {
    generating.value = false
  }
}

const handleDownload = (row: any) => {
  ElMessage.info('下载功能开发中...')
}

watch([filterAccountId, filterReportType, pagination.currentPage, pagination.pageSize], loadReports)

onMounted(() => {
  loadCloudAccounts()
  loadReports()
})
</script>

<style scoped>
.reports-container {
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