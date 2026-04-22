<template>
  <div class="orders-container">
    <div class="page-header">
      <h2>我的订单</h2>
      <div class="toolbar">
        <el-button type="primary" @click="handleSync">同步数据</el-button>
      </div>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="loadOrders">
        <el-form-item label="云账号">
          <el-select v-model="filters.cloud_account_id" placeholder="选择云账号" clearable style="width: 200px">
            <el-option v-for="account in cloudAccounts" :key="account.id" :label="account.name" :value="account.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select v-model="filters.status" placeholder="订单状态" clearable style="width: 150px">
            <el-option label="待支付" value="pending" />
            <el-option label="已支付" value="paid" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadOrders">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table :data="orders" v-loading="loading" style="width: 100%" row-key="id">
      <el-table-column prop="order_id" label="订单号" width="180" />
      <el-table-column prop="order_type" label="订单类型" width="100" />
      <el-table-column prop="product_name" label="产品名称" min-width="200" />
      <el-table-column prop="amount" label="金额" width="120">
        <template #default="{ row }">
          ¥{{ row.amount?.toFixed(2) || '0.00' }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="effective_time" label="生效时间" width="160" />
      <el-table-column prop="expire_time" label="到期时间" width="160" />
      <el-table-column prop="provider_type" label="云平台" width="100" />
    </el-table>

    <el-pagination
      v-model:current-page="pagination.currentPage"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      :page-sizes="[10, 20, 50]"
      layout="total, sizes, prev, pager, next"
      class="pagination"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getOrders, syncOrders } from '@/api/finance'
import { getCloudAccounts } from '@/api/cloud-account'
import type { Order } from '@/types/finance'

const orders = ref<Order[]>([])
const loading = ref(false)
const cloudAccounts = ref<any[]>([])

const filters = reactive({
  cloud_account_id: undefined as number | undefined,
  status: ''
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    paid: 'success',
    cancelled: 'info'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    cancelled: '已取消'
  }
  return map[status] || status
}

const loadOrders = async () => {
  loading.value = true
  try {
    const res = await getOrders({
      cloud_account_id: filters.cloud_account_id,
      status: filters.status,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    orders.value = res.items || []
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
  filters.status = ''
  pagination.currentPage = 1
  loadOrders()
}

const handleSync = async () => {
  if (!filters.cloud_account_id) {
    ElMessage.warning('请先选择云账号')
    return
  }
  try {
    await syncOrders(filters.cloud_account_id)
    ElMessage.success('订单数据同步成功')
    loadOrders()
  } catch (e) {
    ElMessage.error('同步失败')
  }
}

watch([filters.cloud_account_id, filters.status, pagination.currentPage, pagination.pageSize], loadOrders)

onMounted(() => {
  loadCloudAccounts()
  loadOrders()
})
</script>

<style scoped>
.orders-container {
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