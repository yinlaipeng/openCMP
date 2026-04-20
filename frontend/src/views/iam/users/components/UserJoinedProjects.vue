<template>
  <div class="joined-projects-container">
    <!-- 工具栏 -->
    <div class="toolbar">
      <el-button size="small" @click="handleRefresh" :loading="loading">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
      <el-button size="small" type="primary" @click="handleJoinProject">
        <el-icon><Plus /></el-icon>
        加入项目
      </el-button>
      <el-button
        size="small"
        :disabled="selectedProjects.length === 0"
        @click="handleRemove"
      >
        移除
      </el-button>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-input
        v-model="searchKeyword"
        placeholder="默认为名称搜索，自动匹配IP或ID搜索项，IP或ID多个搜索用英文竖线(|)隔开"
        clearable
        @keyup.enter="handleSearch"
        class="search-input"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-button type="primary" @click="handleSearch">查询</el-button>
    </div>

    <!-- 项目列表 -->
    <el-table
      :data="filteredProjects"
      v-loading="loading"
      @selection-change="handleSelectionChange"
      row-key="id"
      border
      stripe
      class="project-table"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="project_name" label="所属项目" min-width="150" show-overflow-tooltip />
      <el-table-column prop="domain_name" label="所属域" width="120" />
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }">
          <el-tag size="small">{{ row.type || '-' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="group" label="组" width="120" show-overflow-tooltip />
      <el-table-column prop="role" label="角色" width="120" show-overflow-tooltip />
      <el-table-column prop="permission" label="权限" min-width="150" show-overflow-tooltip />
      <el-table-column label="操作" width="80" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="danger" link @click="handleRemoveSingle(row)">
            移除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 空状态 -->
    <EmptyState
      v-if="!loading && projects.length === 0"
      icon="Folder"
      title="暂无数据"
      description="该用户尚未加入任何项目"
      action-text="加入项目"
      @action="handleJoinProject"
    />

    <!-- 加入项目弹窗 -->
    <el-dialog v-model="joinDialogVisible" title="加入项目" width="500px">
      <el-form :model="joinForm" label-width="100px">
        <el-form-item label="选择域">
          <el-select v-model="joinForm.domain_id" placeholder="请选择域" @change="handleDomainChange">
            <el-option v-for="d in domains" :key="d.id" :label="d.name" :value="d.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="选择项目">
          <el-select v-model="joinForm.project_ids" multiple placeholder="请选择项目">
            <el-option v-for="p in availableProjects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="选择角色">
          <el-select v-model="joinForm.role_ids" multiple placeholder="请选择角色">
            <el-option v-for="r in roles" :key="r.id" :label="r.display_name || r.name" :value="r.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="joinDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleJoinSubmit" :loading="joinLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Plus, Search } from '@element-plus/icons-vue'
import { getUserProjects, assignUserToProject, removeUserFromProject, getDomains, getProjects, getRoles } from '@/api/iam'
import type { Project, Domain, Role } from '@/types/iam'
import EmptyState from '@/components/common/EmptyState.vue'

const props = defineProps<{
  userId?: number
}>()

const emit = defineEmits<{
  (e: 'refresh'): void
}>()

// State
const loading = ref(false)
const joinLoading = ref(false)
const projects = ref<any[]>([])
const selectedProjects = ref<any[]>([])
const searchKeyword = ref('')
const domains = ref<Domain[]>([])
const availableProjects = ref<Project[]>([])
const roles = ref<Role[]>([])
const joinDialogVisible = ref(false)

const joinForm = ref({
  domain_id: 0,
  project_ids: [] as number[],
  role_ids: [] as number[]
})

// Computed
const filteredProjects = computed(() => {
  if (!searchKeyword.value) return projects.value
  const keyword = searchKeyword.value.toLowerCase()
  return projects.value.filter(p =>
    p.project_name?.toLowerCase().includes(keyword) ||
    p.domain_name?.toLowerCase().includes(keyword)
  )
})

// Load data
watch(() => props.userId, (id) => {
  if (id) loadProjects()
}, { immediate: true })

const loadProjects = async () => {
  if (!props.userId) return
  loading.value = true
  try {
    const res = await getUserProjects(props.userId)
    projects.value = res.items || []
  } catch (e: any) {
    ElMessage.error(e.message || '加载项目列表失败')
  } finally {
    loading.value = false
  }
}

const loadDomains = async () => {
  try {
    const res = await getDomains({ limit: 100 })
    domains.value = res.items || []
  } catch (e) {
    console.error('加载域失败', e)
  }
}

const handleDomainChange = async (domainId: number) => {
  try {
    const [pRes, rRes] = await Promise.all([
      getProjects({ domain_id: domainId, limit: 100 }),
      getRoles({ domain_id: domainId, limit: 100 })
    ])
    availableProjects.value = pRes.items || []
    roles.value = rRes.items || []
  } catch (e) {
    console.error('加载项目/角色失败', e)
  }
}

// Handlers
const handleRefresh = () => loadProjects()

const handleSearch = () => {
  // Filter is computed, search triggers re-filter
}

const handleSelectionChange = (selection: any[]) => {
  selectedProjects.value = selection
}

const handleJoinProject = async () => {
  await loadDomains()
  joinForm.value = { domain_id: 0, project_ids: [], role_ids: [] }
  joinDialogVisible.value = true
}

const handleJoinSubmit = async () => {
  if (!props.userId || !joinForm.value.project_ids.length || !joinForm.value.role_ids.length) {
    ElMessage.error('请选择项目和角色')
    return
  }
  joinLoading.value = true
  try {
    for (const projectId of joinForm.value.project_ids) {
      for (const roleId of joinForm.value.role_ids) {
        await assignUserToProject(props.userId, projectId, roleId)
      }
    }
    ElMessage.success('加入项目成功')
    joinDialogVisible.value = false
    loadProjects()
    emit('refresh')
  } catch (e: any) {
    ElMessage.error(e.message || '加入项目失败')
  } finally {
    joinLoading.value = false
  }
}

const handleRemove = async () => {
  try {
    await ElMessageBox.confirm(`确定要移除选中的 ${selectedProjects.value.length} 个项目吗？`, '提示', { type: 'warning' })
    for (const p of selectedProjects.value) {
      await removeUserFromProject(props.userId!, p.id)
    }
    ElMessage.success('移除成功')
    loadProjects()
    emit('refresh')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '移除失败')
    }
  }
}

const handleRemoveSingle = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要移除该项目吗？', '提示', { type: 'warning' })
    await removeUserFromProject(props.userId!, row.id)
    ElMessage.success('移除成功')
    loadProjects()
    emit('refresh')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '移除失败')
    }
  }
}

onMounted(() => loadDomains())
</script>

<style scoped>
.joined-projects-container {
  padding: var(--space-4, 16px);
}

.toolbar {
  display: flex;
  gap: var(--space-2, 8px);
  margin-bottom: var(--space-3, 12px);
}

.search-bar {
  display: flex;
  gap: var(--space-2, 8px);
  margin-bottom: var(--space-3, 12px);
}

.search-input {
  flex: 1;
  max-width: 400px;
}

.project-table {
  margin-bottom: var(--space-4, 16px);
}
</style>