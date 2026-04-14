<template>
  <div class="finance-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="title">我的订单</span>
          <el-button type="primary" @click="handleSync">同步数据</el-button>
        </div>
      </template>

      <!-- 云账号筛选 -->
      <div class="filter-bar" style="margin-bottom: 16px;">
        <el-select v-model="selectedAccountId" placeholder="选择云账号" clearable style="width: 200px;">
          <el-option v-for="account in cloudAccounts" :key="account.id" :label="account.name" :value="account.id" />
        </el-select>
        <el-select v-model="selectedStatus" placeholder="订单状态" clearable style="width: 150px; margin-left: 8px;">
          <el-option label="待支付" value="pending" />
          <el-option label="已支付" value="paid" />
          <el-option label="已取消" value="cancelled" />
        </el-select>
      </div>

      <!-- 数据表格 -->
      <el-table :data="orders" v-loading="loading" style="width: 100%">
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

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.currentPage"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        style="margin-top: 16px; justify-content: flex-end;"
      />
    </el-card>
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
const selectedAccountId = ref<number | undefined>()
const selectedStatus = ref<string>('')

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
      cloud_account_id: selectedAccountId.value,
      status: selectedStatus.value,
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

const handleSync = async () => {
  if (!selectedAccountId.value) {
    ElMessage.warning('请先选择云账号')
    return
  }
  try {
    await syncOrders(selectedAccountId.value)
    ElMessage.success('订单数据同步成功')
    loadOrders()
  } catch (e) {
    ElMessage.error('同步失败')
  }
}

watch([selectedAccountId, selectedStatus, pagination.currentPage, pagination.pageSize], loadOrders)

onMounted(() => {
  loadCloudAccounts()
  loadOrders()
})
</script>

<style scoped>
.finance-page {
  height: 100%;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.title {
  font-size: 18px;
  font-weight: bold;
}
.filter-bar {
  display: flex;
  align-items: center;
}
</style>