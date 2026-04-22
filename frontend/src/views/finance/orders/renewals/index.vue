<template>
  <div class="renewals-container">
    <div class="page-header">
      <h2>续费管理</h2>
      <div class="toolbar">
        <el-button type="primary" @click="handleSync">同步资源</el-button>
      </div>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="loadRenewals">
        <el-form-item label="云账号">
          <el-select v-model="filters.cloud_account_id" placeholder="选择云账号" clearable style="width: 200px">
            <el-option v-for="account in cloudAccounts" :key="account.id" :label="account.name" :value="account.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="到期天数">
          <el-select v-model="filters.days_threshold" placeholder="到期天数" clearable style="width: 120px">
            <el-option label="7天内" :value="7" />
            <el-option label="30天内" :value="30" />
            <el-option label="60天内" :value="60" />
            <el-option label="90天内" :value="90" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadRenewals">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="16" style="margin-bottom: 16px;">
      <el-col :span="6">
        <el-statistic title="待续费数量" :value="pagination.total" />
      </el-col>
      <el-col :span="6">
        <el-statistic title="预计续费费用" :value="totalRenewalCost" suffix="元" />
      </el-col>
    </el-row>

    <el-table :data="renewals" v-loading="loading" style="width: 100%" row-key="id">
      <el-table-column prop="instance_id" label="实例ID" width="130" />
      <el-table-column prop="instance_name" label="实例名称" min-width="180" />
      <el-table-column prop="product_type" label="产品类型" width="100">
        <template #default="{ row }">
          <el-tag>{{ row.product_type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="expire_time" label="到期时间" width="160">
        <template #default="{ row }">
          {{ formatDate(row.expire_time) }}
        </template>
      </el-table-column>
      <el-table-column prop="days_remaining" label="剩余天数" width="100">
        <template #default="{ row }">
          <el-tag :type="getDaysTagType(row.days_remaining)">
            {{ row.days_remaining }}天
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="renewal_price" label="续费价格" width="100">
        <template #default="{ row }">
          ¥{{ row.renewal_price?.toFixed(2) || '0.00' }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 'normal' ? 'success' : 'warning'">
            {{ row.status === 'normal' ? '正常' : '即将到期' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary">续费</el-button>
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
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getRenewalResources, syncRenewals } from '@/api/finance'
import { getCloudAccounts } from '@/api/cloud-account'
import type { RenewalResource } from '@/types/finance'

const renewals = ref<RenewalResource[]>([])
const loading = ref(false)
const cloudAccounts = ref<any[]>([])

const filters = reactive({
  cloud_account_id: undefined as number | undefined,
  days_threshold: 30
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

const totalRenewalCost = computed(() => {
  return renewals.value.reduce((sum, item) => sum + (item.renewal_price || 0), 0)
})

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

const getDaysTagType = (days: number) => {
  if (days <= 7) return 'danger'
  if (days <= 30) return 'warning'
  return 'success'
}

const loadRenewals = async () => {
  loading.value = true
  try {
    const res = await getRenewalResources({
      cloud_account_id: filters.cloud_account_id,
      days_threshold: filters.days_threshold,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    renewals.value = res.items || []
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
  filters.days_threshold = 30
  pagination.currentPage = 1
  loadRenewals()
}

const handleSync = async () => {
  if (!filters.cloud_account_id) {
    ElMessage.warning('请先选择云账号')
    return
  }
  try {
    loading.value = true
    const res = await syncRenewals(filters.cloud_account_id, filters.days_threshold)
    ElMessage.success(`同步完成，新增 ${res.count} 条数据`)
    loadRenewals()
  } catch (e) {
    ElMessage.error('同步失败')
  } finally {
    loading.value = false
  }
}

watch([filters.cloud_account_id, filters.days_threshold, pagination.currentPage, pagination.pageSize], loadRenewals)

onMounted(() => {
  loadCloudAccounts()
  loadRenewals()
})
</script>

<style scoped>
.renewals-container {
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