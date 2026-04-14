<template>
  <div class="finance-page">
    <el-card>
      <template #header>
        <span class="title">账单导出中心</span>
      </template>

      <el-form :model="exportForm" label-width="100px">
        <el-form-item label="云账号">
          <el-select v-model="exportForm.cloud_account_id" placeholder="选择云账号" clearable>
            <el-option v-for="a in cloudAccounts" :key="a.id" :label="a.name" :value="a.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="账单周期">
          <el-date-picker v-model="exportForm.billing_cycle" type="month" format="YYYY-MM" value-format="YYYY-MM" />
        </el-form-item>
        <el-form-item label="导出格式">
          <el-radio-group v-model="exportForm.format">
            <el-radio label="csv">CSV</el-radio>
            <el-radio label="excel">Excel</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleExport">导出账单</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { exportBills } from '@/api/finance'
import { getCloudAccounts } from '@/api/cloud-account'

const cloudAccounts = ref<any[]>([])
const exportForm = ref({
  cloud_account_id: undefined as number | undefined,
  billing_cycle: '',
  format: 'excel'
})

const handleExport = async () => {
  try {
    const res = await exportBills(exportForm.value)
    ElMessage.success('导出成功：' + res.download_url)
  } catch (e) {
    ElMessage.error('导出失败')
  }
}

onMounted(async () => {
  const res = await getCloudAccounts({ page: 1, page_size: 100 })
  cloudAccounts.value = res.items || []
})
</script>

<style scoped>
.finance-page {
  height: 100%;
}
.title {
  font-size: 18px;
  font-weight: bold;
}
</style>