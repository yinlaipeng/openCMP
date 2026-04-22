<template>
  <div class="bills-view-container">
    <div class="page-header">
      <h2>账单查看</h2>
      <div class="toolbar">
        <el-button type="primary" @click="handleSync">同步账单</el-button>
      </div>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="loadBills">
        <el-form-item label="云账号">
          <el-select v-model="filters.cloud_account_id" placeholder="选择云账号" clearable style="width: 200px">
            <el-option v-for="account in cloudAccounts" :key="account.id" :label="account.name" :value="account.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="账单周期">
          <el-date-picker
            v-model="filters.billing_cycle"
            type="month"
            placeholder="选择账单周期"
            format="YYYY-MM"
            value-format="YYYY-MM"
            style="width: 150px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadBills">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="16" style="margin-bottom: 16px;">
      <el-col :span="6">
        <el-statistic title="本月总费用" :value="totalCost" suffix="元" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="账单数量" :value="pagination.total" />
      </el-col>
    </el-row>

    <el-table :data="bills" v-loading="loading" style="width: 100%" row-key="id">
      <el-table-column prop="billing_cycle" label="账期" width="100" />
      <el-table-column prop="product_type" label="产品类型" width="120" />
      <el-table-column prop="product_name" label="产品名称" min-width="200" />
      <el-table-column prop="instance_id" label="实例ID" width="150" />
      <el-table-column prop="usage_amount" label="用量" width="100" />
      <el-table-column prop="total_cost" label="费用" width="120">
        <template #default="{ row }">
          ¥{{ row.total_cost?.toFixed(2) || '0.00' }}
        </template>
      </el-table-column>
      <el-table-column prop="billing_method" label="计费方式" width="100" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'paid' ? 'success' : 'warning'">
            {{ row.status === 'paid' ? '已支付' : '待支付' }}
          </el-tag>
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { getBills, syncBills } from '@/api/finance'
import { getCloudAccounts } from '@/api/cloud-account'
import type { Bill } from '@/types/finance'

const bills = ref<Bill[]>([])
const loading = ref(false)
const cloudAccounts = ref<any[]>([])

const filters = reactive({
  cloud_account_id: undefined as number | undefined,
  billing_cycle: ''
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

const totalCost = computed(() => {
  return bills.value.reduce((sum, bill) => sum + (bill.total_cost || 0), 0)
})

const loadBills = async () => {
  loading.value = true
  try {
    const res = await getBills({
      cloud_account_id: filters.cloud_account_id,
      billing_cycle: filters.billing_cycle,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    bills.value = res.items || []
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
  filters.cloud_account_id = undefined
  filters.billing_cycle = ''
  pagination.currentPage = 1
  loadBills()
}

const handleSync = async () => {
  if (!filters.cloud_account_id) {
    ElMessage.warning('请先选择云账号')
    return
  }
  try {
    await syncBills(filters.cloud_account_id)
    ElMessage.success('账单数据同步成功')
    loadBills()
  } catch (e) {
    ElMessage.error('同步失败')
  }
}

watch([filters.cloud_account_id, filters.billing_cycle, pagination.currentPage, pagination.pageSize], loadBills)

onMounted(() => {
  loadCloudAccounts()
  loadBills()
})
</script>

<style scoped>
.bills-view-container {
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