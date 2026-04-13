<template>
  <div class="scheduled-tasks-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">定时同步任务</span>
          <el-button type="primary" @click="showDialog = true">
            <el-icon><Plus /></el-icon>
            添加任务
          </el-button>
        </div>
      </template>

      <el-table :data="tasks" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column prop="type" label="类型" width="150">
          <template #default="{ row }">
            <el-tag>{{ row.type === 'sync_cloud_account' ? '同步云账号' : row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="frequency" label="触发频次" width="120">
          <template #default="{ row }">
            <el-tag :type="getFrequencyType(row.frequency)">
              {{ getFrequencyText(row.frequency) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="triggerTime" label="触发时间" width="120" />
        <el-table-column prop="validRange" label="有效时间" width="200">
          <template #default="{ row }">
            {{ row.validFrom || '长期有效' }} 至 {{ row.validUntil || '长期有效' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">
              {{ row.status === 'active' ? '启用' : '暂停' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cloud_account_id" label="关联云账号" width="150">
          <template #default="{ row }">
            {{ getCloudAccountName(row.cloud_account_id) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="success" @click="handleExecute(row)">执行</el-button>
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleToggle(row)">
              {{ row.status === 'active' ? '暂停' : '启用' }}
            </el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.currentPage"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        class="pagination"
      />
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="showDialog"
      :title="isEdit ? '编辑任务' : '添加任务'"
      width="600px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入任务名称" />
        </el-form-item>

        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择任务类型" style="width: 100%">
            <el-option label="同步云账号" value="sync_cloud_account" />
          </el-select>
        </el-form-item>

        <el-form-item label="关联云账号" prop="cloud_account_id" v-if="form.type === 'sync_cloud_account'">
          <el-select v-model="form.cloud_account_id" placeholder="请选择要同步的云账号" style="width: 100%">
            <el-option v-for="account in cloudAccounts" :key="account.id" :label="account.name" :value="account.id" />
          </el-select>
        </el-form-item>

        <el-form-item label="触发频次" prop="frequency">
          <el-select v-model="form.frequency" placeholder="请选择触发频次" style="width: 100%">
            <el-option label="单次" value="once" />
            <el-option label="每天" value="daily" />
            <el-option label="每周" value="weekly" />
            <el-option label="每月" value="monthly" />
            <el-option label="周期" value="custom" />
          </el-select>
        </el-form-item>

        <el-form-item label="触发时间" prop="triggerTime">
          <el-time-picker
            v-model="form.triggerTime"
            placeholder="选择时间"
            format="HH:mm"
            value-format="HH:mm"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="有效时间" prop="validRange">
          <el-date-picker
            v-model="form.validRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-switch
            v-model="form.status"
            :active-value="'active'"
            :inactive-value="'inactive'"
            active-text="启用"
            inactive-text="暂停"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getScheduledTasks, createScheduledTask, updateScheduledTask, deleteScheduledTask, updateScheduledTaskStatus, executeScheduledTask } from '@/api/scheduled-task'
import { getCloudAccounts } from '@/api/cloud-account'
import type { ScheduledTask, CreateScheduledTaskRequest } from '@/types'

const tasks = ref<ScheduledTask[]>([])
const cloudAccounts = ref<any[]>([])
const loading = ref(false)
const showDialog = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const currentId = ref<number | null>(null)
const formRef = ref<FormInstance>()

const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

const form = reactive<CreateScheduledTaskRequest & { validRange?: string[], cloud_account_id?: number }>({
  name: '',
  type: 'sync_cloud_account',
  frequency: 'daily',
  triggerTime: '02:00',
  validFrom: undefined,
  validUntil: undefined,
  status: 'active',
  cloud_account_id: undefined
})

const rules = {
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择任务类型', trigger: 'change' }],
  frequency: [{ required: true, message: '请选择触发频次', trigger: 'change' }]
}

const loadTasks = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.currentPage,
      page_size: pagination.pageSize
    }
    const res = await getScheduledTasks(params)
    tasks.value = res.items || []
    pagination.total = res.total || 0
  } catch (e) {
    console.error(e)
    ElMessage.error('加载任务失败')
  } finally {
    loading.value = false
  }
}

const getFrequencyText = (freq: string) => {
  const map: Record<string, string> = {
    'once': '单次',
    'daily': '每天',
    'weekly': '每周',
    'monthly': '每月',
    'custom': '周期'
  }
  return map[freq] || freq
}

const getFrequencyType = (freq: string) => {
  const map: Record<string, string> = {
    'once': 'info',
    'daily': 'success',
    'weekly': 'warning',
    'monthly': 'primary',
    'custom': 'danger'
  }
  return map[freq] || 'info'
}

const getCloudAccountName = (accountId: number | null) => {
  if (!accountId) return '未关联'
  const account = cloudAccounts.value.find(a => a.id === accountId)
  return account?.name || `账号${accountId}`
}

const handleEdit = (row: ScheduledTask) => {
  isEdit.value = true
  currentId.value = row.id
  Object.assign(form, {
    name: row.name,
    type: row.type,
    frequency: row.frequency,
    triggerTime: row.triggerTime,
    validRange: row.validFrom && row.validUntil ? [row.validFrom, row.validUntil] : undefined,
    status: row.status as 'active' | 'inactive',
    cloud_account_id: row.cloud_account_id
  })
  showDialog.value = true
}

const handleToggle = async (row: ScheduledTask) => {
  const newStatus = row.status === 'active' ? 'inactive' : 'active'
  try {
    await updateScheduledTaskStatus(row.id, newStatus as 'active' | 'inactive')
    row.status = newStatus
    ElMessage.success(`${newStatus === 'active' ? '启用' : '暂停'}成功`)
  } catch (e) {
    console.error(e)
    ElMessage.error('更新状态失败')
  }
}

const handleDelete = async (row: ScheduledTask) => {
  try {
    await ElMessageBox.confirm(`确定要删除任务 "${row.name}" 吗？`, '提示', {
      type: 'warning'
    })
    await deleteScheduledTask(row.id)
    ElMessage.success('删除成功')
    loadTasks() // 重新加载数据
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
      ElMessage.error('删除失败')
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      // 处理有效时间范围
      const submitData = {
        name: form.name,
        type: form.type,
        frequency: form.frequency,
        trigger_time: form.triggerTime,
        valid_from: form.validRange?.[0] || undefined,
        valid_until: form.validRange?.[1] || undefined,
        status: form.status,
        cloud_account_id: form.cloud_account_id
      }

      if (isEdit.value && currentId.value) {
        // 编辑现有任务
        await updateScheduledTask(currentId.value, submitData)
        ElMessage.success('更新成功')
      } else {
        // 创建新任务
        await createScheduledTask(submitData)
        ElMessage.success('创建成功')
      }
      showDialog.value = false
      loadTasks() // 重新加载数据
    } catch (e) {
      console.error(e)
      ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
    } finally {
      submitting.value = false
    }
  })
}

const handleExecute = async (row: ScheduledTask) => {
  try {
    const result = await executeScheduledTask(row.id)
    ElMessage.success(`任务执行成功: 同步了 ${Object.entries(result.statistics).map(([k, v]) => `${v}个${k}`).join(', ')}`)
    loadTasks() // 重新加载
  } catch (e) {
    console.error(e)
    ElMessage.error('执行任务失败')
  }
}

const loadCloudAccounts = async () => {
  try {
    const res = await getCloudAccounts({ page_size: 100 })
    cloudAccounts.value = res.items || []
  } catch (e) {
    console.error(e)
  }
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  loadTasks()
}

const handleCurrentChange = (page: number) => {
  pagination.currentPage = page
  loadTasks()
}

onMounted(() => {
  loadTasks()
  loadCloudAccounts()
})
</script>

<style scoped>
.scheduled-tasks-page {
  height: 100%;
}

.page-card {
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

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
}
</style>