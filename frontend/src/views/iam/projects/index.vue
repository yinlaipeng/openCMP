<template>
  <div class="projects-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">项目</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建项目
          </el-button>
        </div>
      </template>

      <!-- 筛选条件 -->
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="项目名称">
          <el-input v-model="filterForm.name" placeholder="请输入项目名称" clearable @keyup.enter="loadProjects" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterForm.enabled" placeholder="全部" clearable style="width: 100px">
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadProjects">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 项目列表 -->
      <el-table :data="projects" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" min-width="180" show-overflow-tooltip>
          <template #default="{ row }">
            <el-button link @click="handleView(row)" class="name-link">
              {{ row.name }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="manager" label="项目管理员" width="120">
          <template #default="{ row }">
            <span>{{ getProjectManagerName(row) || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="domain_id" label="所属域" width="120">
          <template #default="{ row }">
            <span>{{ getDomainName(row.domain_id) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleManageUsers(row)">管理用户/组</el-button>
            <el-dropdown trigger="click" @command="(cmd) => handleCommand(cmd, row)">
              <el-button size="small" type="primary" link>
                更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="setManager">设置项目管理员</el-dropdown-item>
                  <el-dropdown-item command="delete" divided>
                    <span style="color: #F56C6C">删除</span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.limit"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadProjects"
        @current-change="loadProjects"
        class="pagination"
      />
    </el-card>

    <!-- 添加/编辑项目对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑项目' : '新建项目'" width="700px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="项目名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入项目名称" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入项目描述" />
        </el-form-item>
        <el-form-item label="所属域" prop="domain_id" v-if="isEdit">
          <el-select v-model="form.domain_id" placeholder="请选择所属域" style="width: 100%" disabled>
            <el-option
              v-for="item in domains"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="所属域" prop="domain_id" v-else>
          <el-select v-model="form.domain_id" placeholder="请选择所属域" style="width: 100%" @change="handleDomainChange">
            <el-option
              v-for="item in domains"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <!-- 折叠面板：向该项目添加用户（可选） -->
        <el-collapse v-model="activeCollapse" v-if="!isEdit">
          <el-collapse-item title="向该项目添加用户（可选）" name="addUser">
            <el-form-item label="选择域">
              <el-select
                v-model="addUserForm.domain_id"
                placeholder="请选择域"
                style="width: 100%"
                @change="handleAddUserDomainChange"
              >
                <el-option
                  v-for="item in domains"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="选择用户">
              <el-select
                v-model="addUserForm.user_ids"
                multiple
                filterable
                placeholder="请选择用户"
                style="width: 100%"
              >
                <el-option
                  v-for="item in domainUsers"
                  :key="item.id"
                  :label="`${item.name} (${item.display_name})`"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="选择角色">
              <el-select
                v-model="addUserForm.role_ids"
                multiple
                filterable
                placeholder="请选择角色"
                style="width: 100%"
              >
                <el-option
                  v-for="item in domainRoles"
                  :key="item.id"
                  :label="item.display_name || item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
            <el-alert
              title="提示"
              description="用户来源于所选域，权限通过角色配置"
              type="info"
              show-icon
              :closable="false"
            />
          </el-collapse-item>
        </el-collapse>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">创建</el-button>
      </template>
    </el-dialog>

    <!-- 项目详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="项目详情" width="900px">
      <el-descriptions :column="2" border v-if="currentProject">
        <el-descriptions-item label="ID">{{ currentProject.id }}</el-descriptions-item>
        <el-descriptions-item label="项目名称">{{ currentProject.name }}</el-descriptions-item>
        <el-descriptions-item label="项目管理员">{{ getProjectManagerName(currentProject) || '-' }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ getDomainName(currentProject.domain_id) }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentProject.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentProject.enabled ? 'success' : 'info'" size="small">
            {{ currentProject.enabled ? '启用' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentProject.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentProject.updated_at) }}</el-descriptions-item>
      </el-descriptions>

      <el-tabs v-model="detailTab" style="margin-top: 20px">
        <el-tab-pane label="已加入用户/组" name="users">
          <div class="tab-toolbar">
            <span>项目成员列表</span>
            <el-button size="small" type="primary" @click="handleAddUser">添加用户</el-button>
          </div>
          <el-table :data="projectUsers" v-loading="usersLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="用户名" />
            <el-table-column prop="display_name" label="显示名" />
            <el-table-column prop="email" label="邮箱" />
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button size="small" type="danger" @click="handleRemoveUser(row)">移除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="角色" name="roles">
          <el-table :data="projectRoles" v-loading="rolesLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="角色名" />
            <el-table-column prop="display_name" label="显示名" />
            <el-table-column prop="type" label="类型" width="100">
              <template #default="{ row }">
                <el-tag size="small" :type="row.type === 'system' ? 'warning' : 'success'">
                  {{ row.type === 'system' ? '系统' : '自定义' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="操作日志" name="operation_logs">
          <el-table :data="projectOperationLogs" v-loading="logsLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="operation" label="操作" />
            <el-table-column prop="user" label="操作人" />
            <el-table-column prop="timestamp" label="时间" />
          </el-table>
        </el-tab-pane>
      </el-tabs>

      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 添加用户对话框 -->
    <el-dialog v-model="addUserDialogVisible" title="添加用户到项目" width="600px">
      <el-form :model="addUserForm" label-width="100px">
        <el-form-item label="选择用户" required>
          <el-select
            v-model="addUserForm.user_id"
            placeholder="请选择用户"
            style="width: 100%"
            filterable
          >
            <el-option
              v-for="item in allUsers"
              :key="item.id"
              :label="`${item.name} (${item.display_name})`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="选择角色" required>
          <el-select
            v-model="addUserForm.role_id"
            placeholder="请选择角色"
            style="width: 100%"
            filterable
          >
            <el-option
              v-for="item in allRoles"
              :key="item.id"
              :label="item.display_name || item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addUserDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAddUserSubmit" :loading="addUserSubmitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 设置项目管理员对话框 -->
    <el-dialog v-model="setManagerDialogVisible" title="设置项目管理员" width="600px">
      <el-form :model="setManagerForm" label-width="100px">
        <el-form-item label="选择管理员" required>
          <el-select
            v-model="setManagerForm.user_id"
            placeholder="请选择项目管理员"
            style="width: 100%"
            filterable
          >
            <el-option
              v-for="item in allUsers"
              :key="item.id"
              :label="`${item.name} (${item.display_name})`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="setManagerDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSetManagerSubmit" :loading="setManagerSubmitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Plus, ArrowDown } from '@element-plus/icons-vue'
import { Project, Domain, User, Role } from '@/types/iam'
import {
  getProjects,
  createProject,
  updateProject,
  deleteProject,
  enableProject,
  disableProject,
  getProjectUsers,
  getProjectRoles,
  joinProject,
  removeUserFromProject,
  getDomains,
  getUsers,
  getRoles,
  setProjectManager
} from '@/api/iam'

const projects = ref<Project[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const detailTab = ref('users')
const isEdit = ref(false)
const submitting = ref(false)
const activeCollapse = ref<string[]>([])
const currentProject = ref<Project | null>(null)
const projectUsers = ref<User[]>([])
const projectRoles = ref<Role[]>([])
const projectOperationLogs = ref<any[]>([])
const domains = ref<Domain[]>([])
const domainUsers = ref<User[]>([])
const domainRoles = ref<Role[]>([])
const usersLoading = ref(false)
const rolesLoading = ref(false)
const logsLoading = ref(false)
const formRef = ref<FormInstance>()

// Additional variables for project managers
const projectManagers = ref<{[key: number]: User | null}>({})

const filterForm = reactive({
  name: '',
  enabled: undefined as boolean | undefined
})

const pagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

const form = reactive({
  name: '',
  description: '',
  domain_id: 1
})

const addUserForm = reactive({
  domain_id: 1,
  user_ids: [] as number[],
  role_ids: [] as number[],
  user_id: 0,
  role_id: 0
})

const setManagerForm = reactive({
  user_id: 0
})

const addUserDialogVisible = ref(false)
const addUserSubmitting = ref(false)
const setManagerDialogVisible = ref(false)
const setManagerSubmitting = ref(false)

const allUsers = ref<User[]>([])
const allRoles = ref<Role[]>([])

const rules: FormRules = {
  name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
  domain_id: [{ required: true, message: '请选择所属域', trigger: 'change' }]
}

const loadProjects = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit
    }
    if (filterForm.name) params.keyword = filterForm.name
    if (filterForm.enabled !== undefined) params.enabled = filterForm.enabled

    const res = await getProjects(params)
    projects.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载项目列表失败')
  } finally {
    loading.value = false
  }
}

const loadDomains = async () => {
  try {
    const res = await getDomains({ limit: 100, enabled: true })
    domains.value = res.items || []
    if (domains.value.length > 0) {
      // 设置默认域
      if (!form.domain_id) {
        form.domain_id = domains.value[0].id
      }
      if (!addUserForm.domain_id) {
        addUserForm.domain_id = domains.value[0].id
      }
    }
  } catch (e: any) {
    console.error(e)
  }
}

const loadDomainUsersAndRoles = async (domainId: number) => {
  try {
    const [usersRes, rolesRes] = await Promise.all([
      getUsers({ domain_id: domainId, limit: 100 }),
      getRoles({ domain_id: domainId, limit: 100 })
    ])
    domainUsers.value = usersRes.items || []
    domainRoles.value = rolesRes.items || []
  } catch (e: any) {
    console.error(e)
  }
}

const handleAddUserDomainChange = (domainId: number) => {
  addUserForm.user_ids = []
  addUserForm.role_ids = []
  loadDomainUsersAndRoles(domainId)
}

const loadProjectUsers = async (projectId: number) => {
  usersLoading.value = true
  try {
    const res = await getProjectUsers(projectId)
    projectUsers.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    usersLoading.value = false
  }
}

const loadProjectRoles = async (projectId: number) => {
  rolesLoading.value = true
  try {
    const res = await getProjectRoles(projectId)
    projectRoles.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    rolesLoading.value = false
  }
}

const resetFilter = () => {
  filterForm.name = ''
  filterForm.enabled = undefined
  pagination.page = 1
  loadProjects()
}

const handleCreate = async () => {
  isEdit.value = false
  form.name = ''
  form.description = ''
  form.domain_id = 1
  await loadDomains()
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (isEdit.value) {
        await updateProject(currentProject.value.id, form)
        ElMessage.success('更新成功')
      } else {
        // 创建项目
        const project = await createProject({
          name: form.name,
          description: form.description,
          domain_id: form.domain_id
        })

        // 如果折叠面板展开且选择了用户和角色，分配用户到项目
        if (activeCollapse.value.includes('addUser') &&
            addUserForm.user_ids.length > 0 &&
            addUserForm.role_ids.length > 0) {
          await joinProject(project.id, addUserForm.user_ids, addUserForm.role_ids)
          ElMessage.success('创建成功并添加用户')
        } else {
          ElMessage.success('创建成功')
        }
      }
      dialogVisible.value = false
      loadProjects()
    } catch (e: any) {
      ElMessage.error(e.message || (isEdit.value ? '更新失败' : '创建失败'))
    } finally {
      submitting.value = false
    }
  })
}

const handleEdit = async (row: Project) => {
  isEdit.value = true
  currentProject.value = row
  form.name = row.name
  form.description = row.description || ''
  form.domain_id = row.domain_id
  await loadDomains()
  dialogVisible.value = true
}

const getDomainName = (domainId: number) => {
  const domain = domains.value.find(d => d.id === domainId);
  return domain ? domain.name : `域#${domainId}`;
}

const handleManageUsers = (row: Project) => {
  // Open the detail view with users tab active
  currentProject.value = row
  detailTab.value = 'users'
  loadProjectUsers(row.id)
  loadProjectRoles(row.id)
  detailDialogVisible.value = true
}

const handleSetManager = async (row: Project) => {
  currentProject.value = row
  await loadAllUsers(row.domain_id)
  setManagerForm.user_id = row.manager_id || 0
  setManagerDialogVisible.value = true
}

const loadAllUsers = async (domainId: number) => {
  try {
    const [usersRes, rolesRes] = await Promise.all([
      getUsers({ domain_id: domainId, limit: 100 }),
      getRoles({ domain_id: domainId, limit: 100 })
    ])
    allUsers.value = usersRes.items || []
    allRoles.value = rolesRes.items || []
  } catch (e: any) {
    console.error(e)
    ElMessage.error(e.message || '加载用户/角色列表失败')
  }
}

const handleSetManagerSubmit = async () => {
  if (!currentProject.value) return

  if (!setManagerForm.user_id) {
    ElMessage.warning('请选择项目管理员')
    return
  }

  setManagerSubmitting.value = true
  try {
    await setProjectManager(currentProject.value.id, setManagerForm.user_id)
    ElMessage.success('设置项目管理员成功')
    setManagerDialogVisible.value = false
    loadProjects() // Refresh the project list to show updated manager
  } catch (e: any) {
    ElMessage.error(e.message || '设置项目管理员失败')
  } finally {
    setManagerSubmitting.value = false
  }
}

const getProjectManagerName = (project: Project) => {
  if (!project.manager_id) return '-'
  const manager = allUsers.value.find(user => user.id === project.manager_id)
  return manager ? `${manager.name} (${manager.display_name})` : `ID: ${project.manager_id}`
}

const handleCommand = (command: string, row: Project) => {
  switch (command) {
    case 'setManager': handleSetManager(row); break
    case 'delete': handleDelete(row); break
  }
}

const handleView = async (row: Project) => {
  currentProject.value = row
  detailTab.value = 'users'
  await Promise.all([
    loadProjectUsers(row.id),
    loadProjectRoles(row.id),
    loadProjectOperationLogs(row.id)
  ])
  detailDialogVisible.value = true
}

const loadProjectOperationLogs = async (projectId: number) => {
  logsLoading.value = true
  try {
    // For now, returning empty array as placeholder - implement later if needed
    projectOperationLogs.value = []
  } catch (e: any) {
    console.error(e)
    projectOperationLogs.value = []
  } finally {
    logsLoading.value = false
  }
}

const handleAddUser = async () => {
  // 打开折叠面板
  activeCollapse.value = ['addUser']
  addUserForm.domain_id = currentProject.value.domain_id
  addUserForm.user_ids = []
  addUserForm.role_ids = []
  await loadDomainUsersAndRoles(addUserForm.domain_id)
}

const handleRemoveUser = async (row: User) => {
  try {
    await ElMessageBox.confirm('确定要移除该用户吗？', '提示', { type: 'warning' })
    await removeUserFromProject(currentProject.value.id, row.id)
    ElMessage.success('移除用户成功')
    await loadProjectUsers(currentProject.value.id)
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '移除用户失败')
    }
  }
}

const handleToggleEnable = async (row: Project) => {
  try {
    const action = row.enabled ? '禁用' : '启用'
    await ElMessageBox.confirm(`确定要${action}该项目吗？`, '提示', { type: 'warning' })

    if (row.enabled) {
      await disableProject(row.id)
    } else {
      await enableProject(row.id)
    }

    ElMessage.success(`${action}成功`)
    loadProjects()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || `${action}失败`)
    }
  }
}

const handleDelete = async (row: Project) => {
  try {
    await ElMessageBox.confirm('确定要删除该项目吗？', '提示', { type: 'warning' })
    await deleteProject(row.id)
    ElMessage.success('删除成功')
    loadProjects()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const handleAddUserSubmit = async () => {
  if (!currentProject.value) return

  if (!addUserForm.user_id || !addUserForm.role_id) {
    ElMessage.warning('请选择用户和角色')
    return
  }

  addUserSubmitting.value = true
  try {
    await joinProject(currentProject.value.id, [addUserForm.user_id], [addUserForm.role_id])
    ElMessage.success('添加用户成功')
    addUserDialogVisible.value = false
    await loadProjectUsers(currentProject.value.id)
  } catch (e: any) {
    ElMessage.error(e.message || '添加用户失败')
  } finally {
    addUserSubmitting.value = false
  }
}

onMounted(() => {
  loadProjects()
})
</script>

<style scoped>
.projects-page {
  height: 100%;
  padding: 20px;
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

.filter-form {
  margin-bottom: 16px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.tab-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.name-link {
  color: #409eff;
  text-decoration: underline;
  cursor: pointer;
}
</style>