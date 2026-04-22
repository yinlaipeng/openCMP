<template>
  <div class="budgets-container">
    <div class="page-header">
      <h2>预算管理</h2>
      <div class="toolbar">
        <el-button type="primary" @click="handleCreate">新建预算</el-button>
      </div>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" @submit.prevent="loadBudgets">
        <el-form-item label="云账号">
          <el-select v-model="selectedAccountId" placeholder="选择云账号" clearable style="width: 200px">
            <el-option v-for="account in cloudAccounts" :key="account.id" :label="account.name" :value="account.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="预算类型">
          <el-select v-model="selectedBudgetType" placeholder="预算类型" clearable style="width: 120px">
            <el-option label="月度预算" value="monthly" />
            <el-option label="季度预算" value="quarterly" />
            <el-option label="年度预算" value="yearly" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="selectedStatus" placeholder="状态" clearable style="width: 100px">
            <el-option label="生效中" value="active" />
            <el-option label="已停用" value="inactive" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadBudgets">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 预算列表表格 -->
    <el-table :data="budgets" v-loading="loading" style="width: 100%" row-key="id">
      <el-table-column prop="name" label="预算名称" width="180" />
      <el-table-column prop="type" label="预算类型" width="100">
        <template #default="{ row }">
          <el-tag>{{ getBudgetTypeText(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="预算金额" width="120">
        <template #default="{ row }">
          ¥{{ row.amount?.toFixed(2) || '0.00' }}
        </template>
      </el-table-column>
      <el-table-column prop="alert_threshold" label="预警阈值" width="100">
        <template #default="{ row }">
          {{ row.alert_threshold }}%
        </template>
      </el-table-column>
      <el-table-column prop="current_usage" label="当前使用" width="120">
        <template #default="{ row }">
          <div>
            ¥{{ row.current_usage?.toFixed(2) || '0.00' }}
            <div style="font-size: 12px; color: #909399;">
              {{ getUsagePercent(row) }}%
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="使用进度" width="150">
        <template #default="{ row }">
          <el-progress
            :percentage="getUsagePercent(row)"
            :status="getProgressStatus(row)"
            :stroke-width="10"
          />
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'info'">
            {{ row.status === 'active' ? '生效中' : '已停用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新建/编辑预算对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑预算' : '新建预算'"
      width="500px"
      destroy-on-close
    >
      <el-form :model="formData" label-width="100px" :rules="formRules" ref="formRef">
        <el-form-item label="预算名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入预算名称" />
        </el-form-item>
        <el-form-item label="关联云账号">
          <el-select v-model="formData.cloud_account_id" placeholder="选择云账号（可选，留空为全局预算）" clearable style="width: 100%;">
            <el-option v-for="account in cloudAccounts" :key="account.id" :label="account.name" :value="account.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="预算类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择预算类型" style="width: 100%;">
            <el-option label="月度预算" value="monthly" />
            <el-option label="季度预算" value="quarterly" />
            <el-option label="年度预算" value="yearly" />
          </el-select>
        </el-form-item>
        <el-form-item label="预算金额" prop="amount">
          <el-input-number v-model="formData.amount" :min="0" :precision="2" style="width: 100%;" />
        </el-form-item>
        <el-form-item label="预警阈值" prop="alert_threshold">
          <el-slider v-model="formData.alert_threshold" :min="0" :max="100" show-input :show-input-controls="false" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="formData.status">
            <el-radio label="active">生效</el-radio>
            <el-radio label="inactive">停用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { getBudgets, createBudget, updateBudget, deleteBudget } from '@/api/finance'
import { getCloudAccounts } from '@/api/cloud-account'
import type { Budget } from '@/types/finance'

const budgets = ref<Budget[]>([])
const loading = ref(false)
const cloudAccounts = ref<any[]>([])
const selectedAccountId = ref<number | undefined>()
const selectedBudgetType = ref<string>('')
const selectedStatus = ref<string>('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive<Partial<Budget>>({
  name: '',
  cloud_account_id: undefined,
  type: 'monthly',
  amount: 1000,
  alert_threshold: 80,
  status: 'active'
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入预算名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择预算类型', trigger: 'change' }],
  amount: [{ required: true, message: '请输入预算金额', trigger: 'blur' }]
}

const getBudgetTypeText = (type: string) => {
  const map: Record<string, string> = {
    monthly: '月度',
    quarterly: '季度',
    yearly: '年度'
  }
  return map[type] || type
}

const getUsagePercent = (budget: Budget) => {
  if (!budget.amount || budget.amount === 0) return 0
  const percent = (budget.current_usage / budget.amount) * 100
  return Math.min(Math.round(percent), 100)
}

const getProgressStatus = (budget: Budget) => {
  const percent = getUsagePercent(budget)
  if (percent >= 100) return 'exception'
  if (percent >= budget.alert_threshold) return 'warning'
  return 'success'
}

const loadBudgets = async () => {
  loading.value = true
  try {
    const params: any = {}
    if (selectedAccountId.value) params.cloud_account_id = selectedAccountId.value
    if (selectedBudgetType.value) params.type = selectedBudgetType.value
    if (selectedStatus.value) params.status = selectedStatus.value
    budgets.value = await getBudgets(params)
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
  selectedAccountId.value = undefined
  selectedBudgetType.value = ''
  selectedStatus.value = ''
  loadBudgets()
}

const handleCreate = () => {
  isEdit.value = false
  Object.assign(formData, {
    name: '',
    cloud_account_id: selectedAccountId.value || undefined,
    type: 'monthly',
    amount: 1000,
    alert_threshold: 80,
    status: 'active'
  })
  dialogVisible.value = true
}

const handleEdit = (row: Budget) => {
  isEdit.value = true
  Object.assign(formData, {
    id: row.id,
    name: row.name,
    cloud_account_id: row.cloud_account_id,
    type: row.type,
    amount: row.amount,
    alert_threshold: row.alert_threshold,
    current_usage: row.current_usage,
    status: row.status
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()

  submitting.value = true
  try {
    if (isEdit.value && formData.id) {
      await updateBudget(formData.id, formData)
      ElMessage.success('预算更新成功')
    } else {
      await createBudget(formData)
      ElMessage.success('预算创建成功')
    }
    dialogVisible.value = false
    loadBudgets()
  } catch (e) {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row: Budget) => {
  try {
    await ElMessageBox.confirm('确定要删除此预算吗？', '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteBudget(row.id)
    ElMessage.success('预算删除成功')
    loadBudgets()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

watch([selectedAccountId, selectedBudgetType, selectedStatus], loadBudgets)

onMounted(() => {
  loadCloudAccounts()
  loadBudgets()
})
</script>

<style scoped>
.budgets-container {
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
</style>