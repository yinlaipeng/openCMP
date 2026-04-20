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
      <div class="filter-bar">
        <el-form :inline="true" :model="filterForm" class="filter-form">
          <el-form-item label="搜索字段">
            <el-select v-model="filterForm.searchField" placeholder="选择搜索字段" style="width: 140px">
              <el-option label="名称" value="name" />
              <el-option label="描述" value="description" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-input
              v-model="filterForm.keyword"
              placeholder="请输入搜索关键词"
              clearable
              style="width: 200px"
              @keyup.enter="loadProjects"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-divider direction="vertical" />
          <el-form-item label="所属域">
            <el-select v-model="filterForm.domain_id" placeholder="全部" clearable style="width: 150px">
              <el-option v-for="item in domains" :key="item.id" :label="item.name" :value="item.id" />
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="filterForm.enabled" placeholder="全部" clearable style="width: 100px">
              <el-option label="启用" :value="true" />
              <el-option label="禁用" :value="false" />
            </el-select>
          </el-form-item>
        </el-form>
        <div class="filter-actions">
          <el-button type="primary" @click="loadProjects">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="resetFilter">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </div>
      </div>

      <!-- 项目列表 -->
      <el-table :data="projects" v-loading="loading" border stripe>
        <el-table-column prop="name" label="名称" min-width="150" show-overflow-tooltip>
          <template #default="{ row }">
            <el-button link @click="handleView(row)" class="name-link">
              {{ row.name }}
            </el-button>
            <el-tag v-if="row.name === 'system' || row.is_system" type="warning" size="small" style="margin-left: 8px;">系统</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="150" show-overflow-tooltip />
        <el-table-column prop="enabled" label="启用状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="admin" label="管理员" width="120">
          <template #default="{ row }">
            <span>{{ row.admin || getProjectManagerName(row) || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="domain_id" label="所属域" width="120">
          <template #default="{ row }">
            <el-tag :type="row.domain_id === 1 || row.domain_id === 'default' ? 'primary' : ''" size="small">
              {{ getDomainName(row.domain_id) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="user_count" label="用户数" width="80">
          <template #default="{ row }">
            {{ row.user_count || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="group_count" label="组数" width="80">
          <template #default="{ row }">
            {{ row.group_count || 0 }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleManageUsersGroups(row)">管理用户/组</el-button>
            <el-dropdown trigger="click" @command="(cmd) => handleCommand(cmd, row)">
              <el-button size="small" type="primary" link>
                更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="view">详情</el-dropdown-item>
                  <el-dropdown-item command="setManager">设置管理员</el-dropdown-item>
                  <el-dropdown-item command="toggleEnable">
                    {{ row.enabled ? '禁用' : '启用' }}
                  </el-dropdown-item>
                  <el-dropdown-item command="delete" :disabled="row.is_system" divided>
                    <span :style="row.is_system ? '' : 'color: #F56C6C'">
                      {{ row.is_system ? '删除（系统项目不可删）' : '删除' }}
                    </span>
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
    <el-dialog v-model="detailDialogVisible" :title="detailDialogTitle" width="900px">
      <el-tabs v-model="detailTab" style="margin-top: 20px">
        <el-tab-pane label="详情" name="details">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">{{ currentProject?.id }}</el-descriptions-item>
            <el-descriptions-item label="项目名称">{{ currentProject?.name }}</el-descriptions-item>
            <el-descriptions-item label="项目管理员">{{ getProjectManagerName(currentProject) || '-' }}</el-descriptions-item>
            <el-descriptions-item label="所属域">{{ getDomainName(currentProject?.domain_id) }}</el-descriptions-item>
            <el-descriptions-item label="描述" :span="2">{{ currentProject?.description || '-' }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="currentProject?.enabled ? 'success' : 'info'" size="small">
                {{ currentProject?.enabled ? '启用' : '禁用' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatDate(currentProject?.created_at) }}</el-descriptions-item>
            <el-descriptions-item label="更新时间">{{ formatDate(currentProject?.updated_at) }}</el-descriptions-item>
          </el-descriptions>

          <div style="margin-top: 20px;">
            <h4>资源列表</h4>
            <el-table :data="projectResources" v-loading="resourcesLoading" border stripe>
              <el-table-column prop="type" label="资源类型" width="120" />
              <el-table-column prop="name" label="资源名称" min-width="150" />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.status === 'active' ? 'success' : 'info'">
                    {{ row.status }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>
        <el-tab-pane label="已加入用户/组" name="users">
          <div class="tab-toolbar">
            <span>项目成员列表</span>
            <el-button size="small" type="primary" @click="handleAddUser">添加用户</el-button>
          </div>
          <el-table :data="projectUsers" v-loading="usersLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="名称" min-width="120" />
            <el-table-column label="类型" width="100">
              <template #default>
                <el-tag type="info">用户</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="role_name" label="角色" width="120" />
            <el-table-column prop="permissions" label="权限" width="150" />
            <el-table-column prop="domain_name" label="所属域" width="120" />
            <el-table-column label="操作" width="150">
              <template #default>
                <el-button size="small" @click="handleModifyRole">修改角色</el-button>
                <el-button size="small" type="danger" @click="handleRemoveUser">移除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="操作日志" name="operation_logs">
          <el-table :data="projectOperationLogs" v-loading="logsLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="operation_time" label="操作时间" width="150" />
            <el-table-column prop="resource_name" label="资源名称" min-width="120" />
            <el-table-column prop="resource_type" label="资源类型" width="120" />
            <el-table-column prop="operation_type" label="操作类型" width="120" />
            <el-table-column prop="service_type" label="服务类型" width="120" />
            <el-table-column prop="risk_level" label="风险级别" width="100" />
            <el-table-column prop="time_type" label="时间类型" width="100" />
            <el-table-column prop="result" label="结果" width="80" />
            <el-table-column prop="operator" label="发起人" width="120" />
            <el-table-column prop="project" label="所属项目" width="120" />
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
      <el-alert
        title="提示"
        :description="`你所选的1个项目将执行设置项目管理员操作，你是否确认操作？`"
        type="warning"
        show-icon
        :closable="false"
        style="margin-bottom: 16px;"
      />

      <div style="margin-bottom: 16px;">
        <strong>先展示当前情况：</strong>
        <div style="margin-top: 8px;">
          <span>名称：{{ currentProject?.name || '-' }}</span>
          <span style="margin-left: 20px;">项目管理员：{{ getProjectManagerName(currentProject) || '-' }}</span>
        </div>
      </div>

      <el-form :model="setManagerForm" label-width="100px">
        <el-form-item label="项目管理员" required>
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
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Plus, ArrowDown, Search, Refresh } from '@element-plus/icons-vue'
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
  setProjectManager,
  getResourceOperationLogs
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
const projectResources = ref<any[]>([])
const domains = ref<Domain[]>([])
const domainUsers = ref<User[]>([])
const domainRoles = ref<Role[]>([])
const usersLoading = ref(false)
const rolesLoading = ref(false)
const logsLoading = ref(false)
const resourcesLoading = ref(false)
const formRef = ref<FormInstance>()

// Additional variables for enhanced features
const isManagingUsers = ref(false) // Track whether the dialog is opened for managing users or viewing details
const detailDialogTitle = computed(() => {
  if (isManagingUsers.value) {
    return '管理用户/组'
  }
  return '项目详情'
})

const filterForm = reactive({
  searchField: 'name',
  keyword: '',
  name: '',
  domain_id: undefined as number | undefined,
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
      offset: (pagination.page - 1) * pagination.limit,
      details: true
    }
    // 根据搜索字段和关键词进行搜索
    if (filterForm.keyword) {
      if (filterForm.searchField === 'name') {
        params.name = filterForm.keyword
      } else if (filterForm.searchField === 'description') {
        params.description = filterForm.keyword
      }
    }
    if (filterForm.domain_id) params.domain_id = filterForm.domain_id
    if (filterForm.enabled !== undefined) params.enabled = filterForm.enabled

    const res = await getProjects(params)
    projects.value = (res.items || res.data || []).map(item => ({
      ...item,
      user_count: item.user_count || 0,
      group_count: item.group_count || 0,
      admin: item.admin || '',
      is_system: item.name === 'system' || item.is_system || false,
      can_delete: item.can_delete !== undefined ? item.can_delete : item.name !== 'system'
    }))
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载项目列表失败')
  } finally {
    loading.value = false
  }
}

const loadDomains = async () => {
  try {
    const res = await getDomains({ limit: 100, enabled: true, details: true })
    domains.value = (res.items || res.data || []).map(item => ({
      ...item,
      id: item.id
    }))
    if (domains.value.length > 0) {
      // 设置默认域
      if (!form.domain_id) {
        form.domain_id = domains.value[0].id as number
      }
      if (!addUserForm.domain_id) {
        addUserForm.domain_id = domains.value[0].id as number
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
  filterForm.searchField = 'name'
  filterForm.keyword = ''
  filterForm.name = ''
  filterForm.domain_id = undefined
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

const handleManageUsersGroups = (row: Project) => {
  // Open the detail view with users tab active for managing users
  currentProject.value = row
  isManagingUsers.value = true
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
    case 'view': handleView(row); break
    case 'setManager': handleSetManager(row); break
    case 'toggleEnable': handleToggleEnable(row); break
    case 'delete': handleDelete(row); break
  }
}

const handleView = async (row: Project) => {
  currentProject.value = row
  isManagingUsers.value = false // Indicate this is a project details view, not managing users
  detailTab.value = 'details'
  await Promise.all([
    loadProjectUsers(row.id),
    loadProjectRoles(row.id),
    loadProjectOperationLogs(row.id),
    loadProjectResources(row.id) // Load project resources for the details tab
  ])
  detailDialogVisible.value = true
}

const loadProjectOperationLogs = async (projectId: number) => {
  logsLoading.value = true
  try {
    // Fetch operation logs for the project
    const res = await getResourceOperationLogs('project', projectId)
    projectOperationLogs.value = res.items || []
  } catch (e: any) {
    console.error(e)
    ElMessage.error(e.message || '加载操作日志失败')
    projectOperationLogs.value = []
  } finally {
    logsLoading.value = false
  }
}

const loadProjectResources = async (projectId: number) => {
  resourcesLoading.value = true
  try {
    // For now, we'll simulate project resources since we don't have a dedicated API
    // In a real implementation, this would call an API to get project resources
    projectResources.value = [
      { type: '虚拟机', name: 'VM-001', status: 'active' },
      { type: '存储', name: 'Disk-001', status: 'active' },
      { type: '网络', name: 'VPC-001', status: 'active' },
      { type: '数据库', name: 'MySQL-001', status: 'inactive' }
    ]
  } catch (e: any) {
    console.error(e)
    projectResources.value = []
  } finally {
    resourcesLoading.value = false
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

const handleRemoveUser = async () => {
  // Placeholder for removing user
  ElMessage.info('移除用户功能待实现')
}

const handleModifyRole = async () => {
  // Placeholder for modifying role
  ElMessage.info('修改角色功能待实现')
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
  loadDomains()  // Load domain list for both forms and filters
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

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
  padding: 16px 20px;
  background-color: #fff;
  border-radius: 4px;
  border: 1px solid #ebeef5;
}

.filter-form {
  flex: 1;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.filter-form .el-form-item {
  margin-bottom: 0;
  margin-right: 8px;
}

.filter-form .el-divider--vertical {
  height: 24px;
  margin: 0 12px;
}

.filter-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: 16px;
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