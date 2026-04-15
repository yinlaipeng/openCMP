<template>
  <div class="scheduled-task-tab">
    <div class="toolbar">
      <el-button type="primary" size="small" @click="showCreateDialog = true">新建定时任务</el-button>
      <el-button size="small" @click="loadScheduledTasks" :loading="loading">刷新</el-button>
    </div>
    <el-table :data="scheduledTasks" v-loading="loading" style="width: 100%">
      <el-table-column prop="name" label="名称" width="150" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'info'">{{ row.status === 'active' ? '正常' : '停用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="enabled" label="启用状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '启用' : '禁用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="type" label="操作动作" width="120" />
      <el-table-column prop="frequency" label="频次" width="100" />
      <el-table-column prop="trigger_time" label="触发时间" width="100" />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="handleExecute(row)" :loading="row.executing">执行</el-button>
          <el-dropdown trigger="click">
            <el-button size="small">更多<el-icon class="el-icon--right"><ArrowDown /></el-icon></el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleEdit(row)">编辑</el-dropdown-item>
                <el-dropdown-item @click="handleToggle(row)">{{ row.enabled ? '禁用' : '启用' }}</el-dropdown-item>
                <el-dropdown-item @click="handleDelete(row)" divided style="color: #f56c6c">删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新建定时任务对话框 -->
    <el-dialog v-model="showCreateDialog" title="新建定时任务" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="任务名称" required>
          <el-input v-model="createForm.name" placeholder="请输入任务名称" />
        </el-form-item>
        <el-form-item label="任务类型" required>
          <el-select v-model="createForm.type" placeholder="请选择任务类型">
            <el-option label="同步云账号" value="sync_cloud_account" />
            <el-option label="同步资源" value="sync_resources" />
            <el-option label="生成报告" value="generate_report" />
          </el-select>
        </el-form-item>
        <el-form-item label="执行频次" required>
          <el-select v-model="createForm.frequency" placeholder="请选择频次">
            <el-option label="每日" value="daily" />
            <el-option label="每周" value="weekly" />
            <el-option label="每月" value="monthly" />
            <el-option label="每小时" value="hourly" />
          </el-select>
        </el-form-item>
        <el-form-item label="触发时间" required>
          <el-input v-model="createForm.trigger_time" placeholder="如：02:00 或 周一" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>

    <!-- 编辑定时任务对话框 -->
    <el-dialog v-model="showEditDialog" title="编辑定时任务" width="500px">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="任务名称" required>
          <el-input v-model="editForm.name" placeholder="请输入任务名称" />
        </el-form-item>
        <el-form-item label="任务类型" required>
          <el-select v-model="editForm.type" placeholder="请选择任务类型">
            <el-option label="同步云账号" value="sync_cloud_account" />
            <el-option label="同步资源" value="sync_resources" />
            <el-option label="生成报告" value="generate_report" />
          </el-select>
        </el-form-item>
        <el-form-item label="执行频次" required>
          <el-select v-model="editForm.frequency" placeholder="请选择频次">
            <el-option label="每日" value="daily" />
            <el-option label="每周" value="weekly" />
            <el-option label="每月" value="monthly" />
            <el-option label="每小时" value="hourly" />
          </el-select>
        </el-form-item>
        <el-form-item label="触发时间" required>
          <el-input v-model="editForm.trigger_time" placeholder="如：02:00 或 周一" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleEditSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import {
  getScheduledTasks,
  createScheduledTask,
  updateScheduledTask,
  deleteScheduledTask,
  updateScheduledTaskStatus,
  executeScheduledTask
} from '@/api/scheduled-task'

interface Props { accountId: number }
const props = defineProps<Props>()

interface ScheduledTaskRow {
  id: number
  name: string
  status: string
  enabled: boolean
  type: string
  frequency: string
  trigger_time: string
  executing?: boolean
}

const loading = ref(false)
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const scheduledTasks = ref<ScheduledTaskRow[]>([])
const currentTask = ref<ScheduledTaskRow | null>(null)

const createForm = ref({
  name: '',
  type: 'sync_cloud_account',
  frequency: 'daily',
  trigger_time: '02:00'
})

const editForm = ref({
  name: '',
  type: '',
  frequency: '',
  trigger_time: ''
})

onMounted(() => { loadScheduledTasks() })

async function loadScheduledTasks() {
  loading.value = true
  try {
    const res = await getScheduledTasks({ cloud_account_id: props.accountId })
    scheduledTasks.value = (res.items || []).map((item: any) => ({
      ...item,
      enabled: item.status === 'active'
    }))
  } catch (error) {
    ElMessage.warning('获取定时任务失败')
    scheduledTasks.value = [
      { id: 1, name: '每日同步', status: 'active', enabled: true, type: 'sync_resources', frequency: 'daily', trigger_time: '02:00' },
      { id: 2, name: '每周统计', status: 'active', enabled: true, type: 'generate_report', frequency: 'weekly', trigger_time: '周一' }
    ]
  } finally { loading.value = false }
}

async function handleCreate() {
  if (!createForm.value.name) {
    ElMessage.warning('请填写任务名称')
    return
  }
  try {
    await createScheduledTask({
      name: createForm.value.name,
      type: createForm.value.type,
      frequency: createForm.value.frequency,
      trigger_time: createForm.value.trigger_time,
      status: 'active',
      cloud_account_id: props.accountId
    })
    ElMessage.success('定时任务创建成功')
    showCreateDialog.value = false
    createForm.value = { name: '', type: 'sync_cloud_account', frequency: 'daily', trigger_time: '02:00' }
    loadScheduledTasks()
  } catch (error) {
    ElMessage.error('创建定时任务失败')
  }
}

function handleEdit(row: ScheduledTaskRow) {
  currentTask.value = row
  editForm.value = {
    name: row.name,
    type: row.type,
    frequency: row.frequency,
    trigger_time: row.trigger_time
  }
  showEditDialog.value = true
}

async function handleEditSubmit() {
  if (!currentTask.value) return
  try {
    await updateScheduledTask(currentTask.value.id, {
      name: editForm.value.name,
      type: editForm.value.type,
      frequency: editForm.value.frequency,
      trigger_time: editForm.value.trigger_time
    })
    ElMessage.success('定时任务更新成功')
    showEditDialog.value = false
    loadScheduledTasks()
  } catch (error) {
    ElMessage.error('更新定时任务失败')
  }
}

async function handleExecute(row: ScheduledTaskRow) {
  row.executing = true
  try {
    const res = await executeScheduledTask(row.id)
    ElMessage.success(`任务执行成功：${res.message}`)
    loadScheduledTasks()
  } catch (error) {
    ElMessage.error('执行任务失败')
  } finally {
    row.executing = false
  }
}

async function handleToggle(row: ScheduledTaskRow) {
  try {
    const newStatus = row.enabled ? 'inactive' : 'active'
    await updateScheduledTaskStatus(row.id, newStatus as 'active' | 'inactive')
    ElMessage.success(newStatus === 'active' ? '已启用' : '已禁用')
    loadScheduledTasks()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

async function handleDelete(row: ScheduledTaskRow) {
  try {
    await ElMessageBox.confirm('确认删除该定时任务？', '删除确认', { type: 'warning' })
    await deleteScheduledTask(row.id)
    ElMessage.success('删除成功')
    loadScheduledTasks()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}
</script>

<style scoped>
.scheduled-task-tab { padding: 10px; }
.toolbar { margin-bottom: 16px; display: flex; gap: 8px; }
</style>