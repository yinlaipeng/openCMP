<template>
  <div class="cost-analysis-container">
    <div class="page-header">
      <h2>成本分析</h2>
      <div class="toolbar">
        <el-button type="primary" @click="handleRefresh">刷新数据</el-button>
      </div>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="loadAnalysis">
        <el-form-item label="云账号">
          <el-select v-model="filters.cloud_account_id" placeholder="选择云账号" clearable style="width: 200px">
            <el-option v-for="account in cloudAccounts" :key="account.id" :label="account.name" :value="account.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="filters.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadAnalysis">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="16" style="margin-bottom: 16px;">
      <el-col :span="6">
        <el-statistic title="总成本" :value="totalCost" suffix="元" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="日均成本" :value="avgDailyCost" suffix="元" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="成本趋势" :value="costTrend">
          <template #suffix>
            <span :style="{ color: costTrend >= 0 ? '#F56C6C' : '#67C23A' }">
              {{ costTrend >= 0 ? '↑' : '↓' }}
            </span>
          </template>
        </el-statistic>
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="16">
      <el-col :span="16">
        <el-card shadow="never">
          <template #header>
            <span>成本趋势</span>
          </template>
          <div v-if="analysisData.length > 0" style="height: 300px;">
            <div class="chart-container">
              <div class="bar-chart">
                <div v-for="(item, index) in analysisData" :key="index" class="bar-item">
                  <div class="bar-label">{{ item.period }}</div>
                  <div class="bar-value" :style="{ width: getBarWidth(item.total_cost) + '%' }"></div>
                  <div class="bar-amount">¥{{ item.total_cost?.toFixed(2) || '0' }}</div>
                </div>
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无数据" />
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="never">
          <template #header>
            <span>产品分布</span>
          </template>
          <div v-if="productDistribution.length > 0" style="height: 300px;">
            <div class="distribution-list">
              <div v-for="(item, index) in productDistribution" :key="index" class="distribution-item">
                <div class="product-name">{{ item.name }}</div>
                <div class="product-bar">
                  <div class="product-bar-fill" :style="{ width: item.percent + '%' }"></div>
                </div>
                <div class="product-percent">{{ item.percent.toFixed(1) }}%</div>
                <div class="product-amount">¥{{ item.amount.toFixed(2) }}</div>
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无数据" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getCostAnalysis } from '@/api/finance'
import { getCloudAccounts } from '@/api/cloud-account'
import type { CostAnalysisData } from '@/types/finance'

const analysisData = ref<CostAnalysisData[]>([])
const loading = ref(false)
const cloudAccounts = ref<any[]>([])

const filters = reactive({
  cloud_account_id: undefined as number | undefined,
  dateRange: [] as string[]
})

const totalCost = computed(() => {
  return analysisData.value.reduce((sum, item) => sum + (item.total_cost || 0), 0)
})

const avgDailyCost = computed(() => {
  if (analysisData.value.length === 0) return 0
  return totalCost.value / analysisData.value.length
})

const costTrend = computed(() => {
  if (analysisData.value.length < 2) return 0
  const last = analysisData.value[analysisData.value.length - 1]?.total_cost || 0
  const prev = analysisData.value[analysisData.value.length - 2]?.total_cost || 0
  if (prev === 0) return 0
  return Math.round(((last - prev) / prev) * 100)
})

const maxCost = computed(() => {
  return Math.max(...analysisData.value.map(item => item.total_cost || 0), 1)
})

const getBarWidth = (cost: number) => {
  return (cost / maxCost.value) * 100
}

const productDistribution = computed(() => {
  const distribution: Record<string, number> = {}
  analysisData.value.forEach(item => {
    if (item.product_costs) {
      Object.entries(item.product_costs).forEach(([name, cost]) => {
        distribution[name] = (distribution[name] || 0) + cost
      })
    }
  })
  const total = Object.values(distribution).reduce((sum, v) => sum + v, 0)
  return Object.entries(distribution)
    .map(([name, amount]) => ({
      name,
      amount,
      percent: total > 0 ? (amount / total) * 100 : 0
    }))
    .sort((a, b) => b.amount - a.amount)
    .slice(0, 5)
})

const loadAnalysis = async () => {
  loading.value = true
  try {
    const params: any = {}
    if (filters.cloud_account_id) {
      params.cloud_account_id = filters.cloud_account_id
    }
    if (filters.dateRange && filters.dateRange.length === 2) {
      params.start_date = filters.dateRange[0]
      params.end_date = filters.dateRange[1]
    }
    analysisData.value = await getCostAnalysis(params)
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
  filters.cloud_account_id = undefined
  filters.dateRange = []
  loadAnalysis()
}

const handleRefresh = () => {
  loadAnalysis()
  ElMessage.success('数据已刷新')
}

watch([filters.cloud_account_id, filters.dateRange], loadAnalysis)

onMounted(() => {
  loadCloudAccounts()
  loadAnalysis()
})
</script>

<style scoped>
.cost-analysis-container {
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

.chart-container {
  padding: 16px;
}
.bar-chart {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.bar-item {
  display: flex;
  align-items: center;
  gap: 8px;
}
.bar-label {
  width: 80px;
  font-size: 12px;
  color: #606266;
}
.bar-value {
  height: 20px;
  background: linear-gradient(90deg, #409EFF, #67C23A);
  border-radius: 4px;
  transition: width 0.3s;
}
.bar-amount {
  font-size: 12px;
  color: #303133;
  min-width: 80px;
}

.distribution-list {
  padding: 16px;
}
.distribution-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}
.product-name {
  width: 80px;
  font-size: 12px;
  color: #606266;
  white-space: nowrap;
  overflow: hidden;
}
.product-bar {
  width: 80px;
  height: 16px;
  background: #EBEEF5;
  border-radius: 4px;
  position: relative;
}
.product-bar-fill {
  height: 100%;
  background: #409EFF;
  border-radius: 4px;
}
.product-percent {
  width: 50px;
  font-size: 12px;
  color: #606266;
}
.product-amount {
  font-size: 12px;
  color: #303133;
}
</style>