<template>
  <div class="finance-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="title">续费管理</span>
        </div>
      </template>

      <el-table :data="renewals" v-loading="loading">
        <el-table-column prop="instance_name" label="实例名称" min-width="200" />
        <el-table-column prop="product_type" label="产品类型" width="120" />
        <el-table-column prop="expire_time" label="到期时间" width="160" />
        <el-table-column prop="days_remaining" label="剩余天数" width="100">
          <template #default="{ row }">
            <el-tag :type="row.days_remaining <= 7 ? 'danger' : row.days_remaining <= 30 ? 'warning' : 'success'">
              {{ row.days_remaining }}天
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="renewal_price" label="续费价格" width="120">
          <template #default="{ row }">¥{{ row.renewal_price?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="primary">续费</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getRenewalResources } from '@/api/finance'
import type { RenewalResource } from '@/types/finance'

const renewals = ref<RenewalResource[]>([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const res = await getRenewalResources({ days_threshold: 30 })
    renewals.value = res.items || []
  } finally {
    loading.value = false
  }
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
</style>