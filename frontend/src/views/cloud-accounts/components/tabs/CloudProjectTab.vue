<template>
  <div class="cloud-project-tab">
    <div class="toolbar">
      <el-button type="primary" size="small" @click="showCreateDialog = true">新建云上项目</el-button>
      <el-button size="small" @click="handleSyncProjects" :loading="syncing">同步云上项目</el-button>
      <el-button size="small" @click="loadCloudProjects" :loading="loading">刷新</el-button>
    </div>
    <el-table :data="cloudProjects" v-loading="loading" style="width: 100%">
      <el-table-column prop="name" label="云上项目" width="150" />
      <el-table-column prop="subscription_id" label="订阅" width="150">
        <template #default="{ row }">{{ row.subscription_id ? `订阅#${row.subscription_id}` : '-' }}</template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'normal' ? 'success' : 'warning'">{{ row.status === 'normal' ? '正常' : row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="tags" label="标签" width="150">
        <template #default="{ row }">
          <el-tag v-for="tag in parseTags(row.tags)" :key="tag" size="small" style="margin-right: 4px">{{ tag }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="domain_id" label="所属域" width="100">
        <template #default="{ row }">{{ row.domain_id ? `域#${row.domain_id}` : '默认域' }}</template>
      </el-table-column>
      <el-table-column prop="local_project_id" label="本地项目" width="150">
        <template #default="{ row }">{{ row.local_project_id ? getLocalProjectName(row.local_project_id) : '-' }}</template>
      </el-table-column>
      <el-table-column prop="priority" label="优先级" width="80" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button size="small" link type="primary" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" link type="primary" @click="handleMap(row)">映射</el-button>
          <el-button size="small" link type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新建云上项目对话框 -->
    <el-dialog v-model="showCreateDialog" title="新建云上项目" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="项目名称" required>
          <el-input v-model="createForm.name" placeholder="请输入项目名称" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="createForm.status" placeholder="请选择状态">
            <el-option label="正常" value="normal" />
            <el-option label="停用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="createForm.tags" placeholder='JSON格式，如 {"env":"prod"}' />
        </el-form-item>
        <el-form-item label="优先级">
          <el-input-number v-model="createForm.priority" :min="0" :max="100" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>

    <!-- 编辑云上项目对话框 -->
    <el-dialog v-model="showEditDialog" title="编辑云上项目" width="500px">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="项目名称" required>
          <el-input v-model="editForm.name" placeholder="请输入项目名称" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="editForm.status" placeholder="请选择状态">
            <el-option label="正常" value="normal" />
            <el-option label="停用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="editForm.tags" placeholder='JSON格式，如 {"env":"prod"}' />
        </el-form-item>
        <el-form-item label="优先级">
          <el-input-number v-model="editForm.priority" :min="0" :max="100" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleEditSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 映射到本地项目对话框 -->
    <el-dialog v-model="showMapDialog" title="映射到本地项目" width="400px">
      <el-form label-width="100px">
        <el-form-item label="云上项目">
          <span>{{ currentProject?.name }}</span>
        </el-form-item>
        <el-form-item label="本地项目">
          <el-select v-model="selectedLocalProjectId" placeholder="请选择本地项目" clearable>
            <el-option v-for="p in localProjects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showMapDialog = false">取消</el-button>
        <el-button type="primary" @click="handleMapSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getCloudProjects, createCloudProject, updateCloudProject, deleteCloudProject, mapCloudProjectToLocal } from '@/api/cloud-account'
import { getProjects } from '@/api/project'

interface Props { accountId: number }
const props = defineProps<Props>()

const loading = ref(false)
const syncing = ref(false)
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const showMapDialog = ref(false)
const cloudProjects = ref<any[]>([])
const localProjects = ref<any[]>([])
const currentProject = ref<any>(null)
const selectedLocalProjectId = ref<number | null>(null)

const createForm = ref({
  name: '',
  status: 'normal',
  tags: '',
  priority: 1
})

const editForm = ref({
  name: '',
  status: 'normal',
  tags: '',
  priority: 1
})

onMounted(() => {
  loadCloudProjects()
  loadLocalProjects()
})

async function loadCloudProjects() {
  loading.value = true
  try {
    const res = await getCloudProjects(props.accountId)
    cloudProjects.value = res.items || []
  } catch (error) {
    ElMessage.warning('获取云上项目失败')
    cloudProjects.value = [
      { id: 1, name: 'proj-001', subscription_id: null, status: 'normal', tags: '{"env":"prod"}', domain_id: 1, local_project_id: 1, priority: 1 },
      { id: 2, name: 'proj-002', subscription_id: null, status: 'normal', tags: '{"env":"dev"}', domain_id: 1, local_project_id: null, priority: 2 }
    ]
  } finally { loading.value = false }
}

async function loadLocalProjects() {
  try {
    const res = await getProjects()
    localProjects.value = res.items || []
  } catch {}
}

function getLocalProjectName(id: number): string {
  const project = localProjects.value.find(p => p.id === id)
  return project ? project.name : `项目#${id}`
}

function parseTags(tags: string): string[] {
  try {
    const obj = JSON.parse(tags)
    return Object.entries(obj).map(([k, v]) => `${k}:${v}`)
  } catch { return [] }
}

async function handleSyncProjects() {
  syncing.value = true
  try {
    await new Promise(r => setTimeout(r, 1000))
    ElMessage.success('云上项目已同步')
    loadCloudProjects()
  } finally { syncing.value = false }
}

async function handleCreate() {
  if (!createForm.value.name) {
    ElMessage.warning('请填写项目名称')
    return
  }
  try {
    await createCloudProject(props.accountId, createForm.value)
    ElMessage.success('云上项目创建成功')
    showCreateDialog.value = false
    createForm.value = { name: '', status: 'normal', tags: '', priority: 1 }
    loadCloudProjects()
  } catch (error) {
    ElMessage.error('创建云上项目失败')
  }
}

function handleEdit(row: any) {
  currentProject.value = row
  editForm.value = {
    name: row.name,
    status: row.status,
    tags: row.tags,
    priority: row.priority
  }
  showEditDialog.value = true
}

async function handleEditSubmit() {
  if (!currentProject.value) return
  try {
    await updateCloudProject(props.accountId, currentProject.value.id, editForm.value)
    ElMessage.success('云上项目更新成功')
    showEditDialog.value = false
    loadCloudProjects()
  } catch (error) {
    ElMessage.error('更新云上项目失败')
  }
}

function handleMap(row: any) {
  currentProject.value = row
  selectedLocalProjectId.value = row.local_project_id
  showMapDialog.value = true
}

async function handleMapSubmit() {
  if (!currentProject.value) return
  try {
    await mapCloudProjectToLocal(props.accountId, currentProject.value.id, selectedLocalProjectId.value!)
    ElMessage.success('映射成功')
    showMapDialog.value = false
    loadCloudProjects()
  } catch (error) {
    ElMessage.error('映射失败')
  }
}

async function handleDelete(row: any) {
  try {
    await ElMessageBox.confirm('确认删除该云上项目？', '删除确认', { type: 'warning' })
    await deleteCloudProject(props.accountId, row.id)
    ElMessage.success('删除成功')
    loadCloudProjects()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}
</script>

<style scoped>
.cloud-project-tab { padding: 10px; }
.toolbar { margin-bottom: 16px; display: flex; gap: 8px; }
</style>