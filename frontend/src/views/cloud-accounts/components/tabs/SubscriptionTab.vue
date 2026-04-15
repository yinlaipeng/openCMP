<template>
  <div class="subscription-tab">
    <div class="toolbar">
      <el-button type="primary" size="small" @click="showCreateDialog = true">新建订阅</el-button>
      <el-button size="small" @click="loadSubscriptions" :loading="loading">刷新</el-button>
    </div>
    <el-table :data="subscriptions" v-loading="loading" style="width: 100%">
      <el-table-column prop="name" label="名称" width="150" />
      <el-table-column prop="subscription_id" label="订阅ID" width="180" />
      <el-table-column prop="enabled" label="启用状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '启用' : '禁用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100" />
      <el-table-column prop="sync_time" label="同步时间" width="160" />
      <el-table-column prop="sync_duration" label="耗时" width="80">
        <template #default="{ row }">{{ row.sync_duration }}s</template>
      </el-table-column>
      <el-table-column prop="sync_status" label="同步状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getSyncStatusType(row.sync_status)">{{ getSyncStatusText(row.sync_status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="domain_name" label="所属域" width="100" />
      <el-table-column prop="default_project_name" label="默认项目" width="120" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="handleChangeProject(row)">更改项目</el-button>
          <el-dropdown trigger="click">
            <el-button size="small">更多<el-icon class="el-icon--right"><ArrowDown /></el-icon></el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleSyncPolicy(row)">同步策略设置</el-dropdown-item>
                <el-dropdown-item @click="handleToggle(row)">{{ row.enabled ? '禁用' : '启用' }}</el-dropdown-item>
                <el-dropdown-item @click="handleDelete(row)" divided style="color: #f56c6c">删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新建订阅对话框 -->
    <el-dialog v-model="showCreateDialog" title="新建订阅" width="500px">
      <el-form :model="createForm" label-width="120px">
        <el-form-item label="订阅名称">
          <el-input v-model="createForm.name" placeholder="请输入订阅名称" />
        </el-form-item>
        <el-form-item label="订阅ID">
          <el-input v-model="createForm.subscription_id" placeholder="请输入订阅ID" />
        </el-form-item>
        <el-form-item label="所属域">
          <el-select v-model="createForm.domain_id" placeholder="请选择所属域">
            <el-option label="默认域" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item label="默认项目">
          <el-select v-model="createForm.default_project_id" placeholder="请选择默认项目" clearable>
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>

    <!-- 更改项目对话框 -->
    <el-dialog v-model="showProjectDialog" title="更改项目" width="400px">
      <el-form label-width="80px">
        <el-form-item label="当前订阅">
          <span>{{ currentSubscription?.name }}</span>
        </el-form-item>
        <el-form-item label="选择项目">
          <el-select v-model="selectedProjectId" placeholder="请选择项目" clearable>
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showProjectDialog = false">取消</el-button>
        <el-button type="primary" @click="handleProjectChangeSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import { getSubscriptions, createSubscription, updateSubscriptionProject, toggleSubscription, syncSubscription, deleteSubscription } from '@/api/cloud-account'
import { getProjects } from '@/api/project'

interface Props {
  accountId: number
}

const props = defineProps<Props>()

const loading = ref(false)
const showCreateDialog = ref(false)
const showProjectDialog = ref(false)
const subscriptions = ref<any[]>([])
const projects = ref<any[]>([])
const currentSubscription = ref<any>(null)

const createForm = ref({
  name: '',
  subscription_id: '',
  domain_id: 1,
  default_project_id: null as number | null
})

const selectedProjectId = ref<number | null>(null)

onMounted(() => {
  loadSubscriptions()
  loadProjects()
})

async function loadSubscriptions() {
  loading.value = true
  try {
    const res = await getSubscriptions(props.accountId)
    subscriptions.value = res.items || []
  } catch (error) {
    ElMessage.warning('获取订阅列表失败')
    // 使用模拟数据作为后备
    subscriptions.value = [
      { id: 1, name: '订阅1', subscription_id: 'sub-001', enabled: true, status: '正常', sync_time: '2026-04-14', sync_duration: 12, sync_status: 'completed', domain_id: 1, default_project_id: 1 }
    ]
  } finally {
    loading.value = false
  }
}

async function loadProjects() {
  try {
    const res = await getProjects()
    projects.value = res.items || []
  } catch {}
}

function getSyncStatusType(status: string): string {
  const types: Record<string, string> = { completed: 'success', failed: 'danger', in_progress: 'warning' }
  return types[status] || 'info'
}

function getSyncStatusText(status: string): string {
  const texts: Record<string, string> = { completed: '同步完成', failed: '同步失败', in_progress: '同步中' }
  return texts[status] || status
}

async function handleCreate() {
  if (!createForm.value.name || !createForm.value.subscription_id) {
    ElMessage.warning('请填写订阅名称和订阅ID')
    return
  }
  try {
    await createSubscription(props.accountId, {
      name: createForm.value.name,
      subscription_id: createForm.value.subscription_id,
      domain_id: createForm.value.domain_id,
      default_project_id: createForm.value.default_project_id
    })
    ElMessage.success('订阅创建成功')
    showCreateDialog.value = false
    createForm.value = { name: '', subscription_id: '', domain_id: 1, default_project_id: null }
    loadSubscriptions()
  } catch (error) {
    ElMessage.error('创建订阅失败')
  }
}

function handleChangeProject(row: any) {
  currentSubscription.value = row
  selectedProjectId.value = row.default_project_id
  showProjectDialog.value = true
}

async function handleProjectChangeSubmit() {
  if (!currentSubscription.value) return
  try {
    await updateSubscriptionProject(props.accountId, currentSubscription.value.id, selectedProjectId.value!)
    ElMessage.success('项目更改成功')
    showProjectDialog.value = false
    loadSubscriptions()
  } catch (error) {
    ElMessage.error('更改项目失败')
  }
}

async function handleSyncPolicy(row: any) {
  try {
    await syncSubscription(props.accountId, row.id)
    ElMessage.success('同步任务已触发')
    loadSubscriptions()
  } catch (error) {
    ElMessage.error('同步失败')
  }
}

async function handleToggle(row: any) {
  try {
    const newEnabled = !row.enabled
    await toggleSubscription(props.accountId, row.id, newEnabled)
    ElMessage.success(newEnabled ? '已启用' : '已禁用')
    loadSubscriptions()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

async function handleDelete(row: any) {
  try {
    await ElMessageBox.confirm('确认删除该订阅？', '删除确认', { type: 'warning' })
    await deleteSubscription(props.accountId, row.id)
    ElMessage.success('删除成功')
    loadSubscriptions()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}
</script>

<style scoped>
.subscription-tab { padding: 10px; }
.toolbar { margin-bottom: 16px; display: flex; gap: 8px; }
</style>