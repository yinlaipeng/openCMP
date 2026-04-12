<template>
  <div class="users-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">用户</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建用户
          </el-button>
        </div>
      </template>

      <!-- 筛选条件 -->
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="用户名">
          <el-input v-model="filterForm.name" placeholder="请输入用户名" clearable @keyup.enter="loadUsers" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="filterForm.email" placeholder="请输入邮箱" clearable @keyup.enter="loadUsers" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterForm.enabled" placeholder="全部" clearable style="width: 100px">
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadUsers">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 用户列表 -->
      <el-table :data="users" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="名称" min-width="150" show-overflow-tooltip>
          <template #default="{ row }">
            <div>
              <el-button type="primary" link @click="handleView(row)" class="name-link">
                <span>{{ row.name }}</span>
              </el-button>
              <br>
              <small style="color: #999;">{{ row.remark || '-' }}</small>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="display_name" label="显示名" min-width="120" show-overflow-tooltip />
        <el-table-column prop="enabled" label="启用状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="console_login" label="控制台登陆" width="100">
          <template #default="{ row }">
            <el-tag :type="row.console_login ? 'success' : 'info'">
              {{ row.console_login ? '允许' : '禁止' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="mfa_enabled" label="MFA" width="80">
          <template #default="{ row }">
            <el-tag :type="row.mfa_enabled ? 'success' : 'info'" size="small">
              {{ row.mfa_enabled ? '开' : '关' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="所属域" width="120">
          <template #default="{ row }">
            <span>{{ getDomainName(row.domain_id) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">修改属性</el-button>
            <el-dropdown trigger="click" @command="(cmd) => handleCommand(cmd, row)">
              <el-button size="small" type="primary" link>
                更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="enable" v-if="!row.enabled">启用</el-dropdown-item>
                  <el-dropdown-item command="disable" v-if="row.enabled">禁用</el-dropdown-item>
                  <el-dropdown-item command="resetPassword">重置密码</el-dropdown-item>
                  <el-dropdown-item command="manageProjects">管理项目</el-dropdown-item>
                  <el-dropdown-item command="resetMFA">重置MFA</el-dropdown-item>
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
        @size-change="loadUsers"
        @current-change="loadUsers"
        class="pagination"
      />
    </el-card>

    <!-- 添加/编辑用户对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? `修改用户属性（${form.name}）` : '新建用户'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="用户名" prop="name" v-if="!isEdit">
          <el-input v-model="form.name" placeholder="请输入用户名" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="显示名" prop="display_name">
          <el-input v-model="form.display_name" placeholder="请输入显示名" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="form.remark" type="textarea" :rows="2" placeholder="请输入备注信息" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input v-model="form.password" type="password" placeholder="请输入密码（至少 8 位）" show-password />
        </el-form-item>
        <el-form-item label="所属域" prop="domain_id" v-if="!isEdit">
          <el-select v-model="form.domain_id" placeholder="请选择所属域" style="width: 100%">
            <el-option
              v-for="domain in allDomains"
              :key="domain.id"
              :label="domain.name"
              :value="domain.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="控制台登录" prop="console_login">
          <el-switch v-model="form.console_login" :active-value="true" :inactive-value="false" />
        </el-form-item>
        <el-form-item label="启用MFA" prop="mfa_enabled">
          <el-switch v-model="form.mfa_enabled" :active-value="true" :inactive-value="false" />
        </el-form-item>

        <!-- 折叠面板：向该用户添加项目（可选） -->
        <el-collapse v-model="activeCollapse" v-if="!isEdit">
          <el-collapse-item title="向该用户添加项目（可选）" name="addProject">
            <el-form-item label="选择域">
              <el-select
                v-model="addUserToProjectForm.domain_id"
                placeholder="请选择域"
                style="width: 100%"
                @change="handleAddUserToProjectDomainChange"
              >
                <el-option
                  v-for="item in allDomains"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="选择项目">
              <el-select
                v-model="addUserToProjectForm.project_ids"
                multiple
                filterable
                placeholder="请选择项目"
                style="width: 100%"
              >
                <el-option
                  v-for="item in allProjects"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="选择角色">
              <el-select
                v-model="addUserToProjectForm.role_ids"
                multiple
                filterable
                placeholder="请选择角色"
                style="width: 100%"
              >
                <el-option
                  v-for="item in allProjectRoles"
                  :key="item.id"
                  :label="item.display_name || item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
            <el-button size="small" type="primary" plain @click="handleCreateRole">+ 没有我想要的角色？立即创建</el-button>
            <el-alert
              title="提示"
              description="项目和角色来源于所选域，权限通过角色配置"
              type="info"
              show-icon
              :closable="false"
              style="margin-top: 10px;"
            />
          </el-collapse-item>
        </el-collapse>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 重置密码对话框 -->
    <el-dialog v-model="resetPwdDialogVisible" title="重置密码" width="500px">
      <el-form :model="resetPwdForm" label-width="100px">
        <el-alert
          title="重置用户密码"
          description="重置后用户需要使用新密码重新登录"
          type="warning"
          show-icon
          style="margin-bottom: 20px"
        />
        <el-form-item label="新密码" required>
          <el-input v-model="resetPwdForm.password" type="password" placeholder="请输入新密码（至少 8 位）" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetPwdDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleResetPasswordSubmit" :loading="resetPwdSubmitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 用户详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="用户详情" width="800px">
      <el-descriptions :column="2" border v-if="currentUser">
        <el-descriptions-item label="ID">{{ currentUser.id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ currentUser.name }}</el-descriptions-item>
        <el-descriptions-item label="显示名">{{ currentUser.display_name }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ currentUser.remark || '-' }}</el-descriptions-item>
        <el-descriptions-item label="邮箱">{{ currentUser.email || '-' }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ currentUser.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentUser.enabled ? 'success' : 'info'" size="small">
            {{ currentUser.enabled ? '启用' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="控制台登录">
          <el-tag :type="currentUser.console_login ? 'success' : 'info'" size="small">
            {{ currentUser.console_login ? '允许' : '禁止' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="MFA">
          <el-tag :type="currentUser.mfa_enabled ? 'success' : 'info'" size="small">
            {{ currentUser.mfa_enabled ? '已开启' : '未开启' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentUser.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentUser.updated_at) }}</el-descriptions-item>
      </el-descriptions>

      <el-tabs v-model="detailTab" style="margin-top: 20px">
        <el-tab-pane label="关联角色" name="roles">
          <div class="tab-toolbar">
            <span>用户的角色和权限</span>
            <el-button size="small" type="primary" @click="handleAssignRole">分配角色</el-button>
          </div>
          <el-table :data="userRoles" v-loading="rolesLoading">
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
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button size="small" type="danger" @click="handleRevokeRole(row)">撤销</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="关联用户组" name="groups">
          <div class="tab-toolbar">
            <span>用户所属的用户组</span>
            <el-button size="small" type="primary" @click="handleJoinGroup">加入用户组</el-button>
          </div>
          <el-table :data="userGroups" v-loading="groupsLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="用户组名" />
            <el-table-column prop="description" label="描述" />
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button size="small" type="danger" @click="handleLeaveGroup(row)">离开</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="已加入的项目" name="projects">
          <el-table :data="userProjects" v-loading="projectsLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="项目名" />
            <el-table-column prop="description" label="描述" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="操作日志" name="logs">
          <el-table :data="userLogs" v-loading="logsLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="operation" label="操作" />
            <el-table-column prop="timestamp" label="时间" />
            <el-table-column prop="details" label="详情" />
          </el-table>
        </el-tab-pane>
      </el-tabs>

      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 分配角色对话框 -->
    <el-dialog v-model="assignRoleDialogVisible" title="分配角色" width="600px">
      <el-form :model="assignRoleForm" label-width="100px">
        <el-form-item label="选择角色" required>
          <el-select v-model="assignRoleForm.role_id" placeholder="请选择角色" style="width: 100%">
            <el-option
              v-for="item in allRoles"
              :key="item.id"
              :label="item.display_name || item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="所属域" required>
          <el-select v-model="assignRoleForm.domain_id" placeholder="请选择域" style="width: 100%">
            <el-option label="Default" :value="1" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="assignRoleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAssignRoleSubmit" :loading="assignRoleSubmitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 加入用户组对话框 -->
    <el-dialog v-model="joinGroupDialogVisible" title="加入用户组" width="600px">
      <el-form :model="joinGroupForm" label-width="100px">
        <el-form-item label="选择用户组" required>
          <el-select v-model="joinGroupForm.group_id" placeholder="请选择用户组" style="width: 100%">
            <el-option
              v-for="item in allGroups"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="joinGroupDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleJoinGroupSubmit" :loading="joinGroupSubmitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Role, Group, Domain, Project } from '@/types/iam'
import {
  getUsers,
  createUser,
  updateUser,
  deleteUser,
  enableUser,
  disableUser,
  resetUserPassword,
  getUserRoles,
  assignRoleToUser,
  revokeRoleFromUser,
  getUserGroups,
  joinGroup,
  leaveGroup,
  getRoles,
  getGroups,
  getDomains,
  getProjects,
  assignUserToProject,
  getUserProjects
} from '@/api/iam'

const users = ref<User[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const resetPwdDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const assignRoleDialogVisible = ref(false)
const joinGroupDialogVisible = ref(false)
const detailTab = ref('roles') // 'roles', 'groups', 'projects', 'logs'
const isEdit = ref(false)
const submitting = ref(false)
const resetPwdSubmitting = ref(false)
const assignRoleSubmitting = ref(false)
const joinGroupSubmitting = ref(false)
const currentUserId = ref(0)
const currentUser = ref<User | null>(null)
const userRoles = ref<Role[]>([])
const userGroups = ref<Group[]>([])
const userProjects = ref<Project[]>([])
const userLogs = ref<any[]>([]) // Operation logs
const allRoles = ref<Role[]>([])
const allGroups = ref<Group[]>([])
const allDomains = ref<Domain[]>([])
const allProjects = ref<Project[]>([])
const allProjectRoles = ref<Role[]>([])
const rolesLoading = ref(false)
const groupsLoading = ref(false)
const projectsLoading = ref(false)
const logsLoading = ref(false)
const formRef = ref<FormInstance>()

const filterForm = reactive({
  name: '',
  email: '',
  enabled: undefined as boolean | undefined
})

const pagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

const form = reactive({
  name: '',
  display_name: '',
  remark: '',
  email: '',
  phone: '',
  password: '',
  domain_id: 1,
  console_login: true,
  mfa_enabled: false
})

const resetPwdForm = reactive({
  password: ''
})

const assignRoleForm = reactive({
  role_id: 0,
  domain_id: 1
})

const joinGroupForm = reactive({
  group_id: 0
})

const activeCollapse = ref<string[]>([])

const addUserToProjectForm = reactive({
  domain_id: 1,
  project_ids: [] as number[],
  role_ids: [] as number[]
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  display_name: [{ required: true, message: '请输入显示名', trigger: 'blur' }],
  email: [
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { min: 8, message: '密码长度至少 8 位', trigger: 'blur' }
  ],
  domain_id: [{ required: true, message: '请选择所属域', trigger: 'change' }]
}

const loadUsers = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit
    }
    if (filterForm.name) params.keyword = filterForm.name
    if (filterForm.email) params.email = filterForm.email
    if (filterForm.enabled !== undefined) params.enabled = filterForm.enabled

    const res = await getUsers(params)
    users.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载用户列表失败')
  } finally {
    loading.value = false
  }
}

const loadUserRoles = async (userId: number) => {
  rolesLoading.value = true
  try {
    const res = await getUserRoles(userId)
    userRoles.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    rolesLoading.value = false
  }
}

const loadUserGroups = async (userId: number) => {
  groupsLoading.value = true
  try {
    const res = await getUserGroups(userId)
    userGroups.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    groupsLoading.value = false
  }
}

const loadUserProjects = async (userId: number) => {
  projectsLoading.value = true
  try {
    const res = await getUserProjects(userId)
    userProjects.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    projectsLoading.value = false
  }
}

const loadUserLogs = async (userId: number) => {
  logsLoading.value = true
  try {
    // For now, we'll use mock data since we don't have an API for user logs yet
    // In a real implementation, we would call a proper API endpoint
    userLogs.value = [
      { id: 1, operation: '登录', timestamp: '2024-01-01 10:00:00', details: '成功登录系统' },
      { id: 2, operation: '密码修改', timestamp: '2024-01-01 11:00:00', details: '用户修改了自己的密码' },
      { id: 3, operation: '角色分配', timestamp: '2024-01-01 12:00:00', details: '管理员为用户分配了新角色' }
    ]
  } catch (e: any) {
    console.error(e)
    // Use mock data in case of error
    userLogs.value = []
  } finally {
    logsLoading.value = false
  }
}

const loadAllRoles = async () => {
  try {
    const res = await getRoles({ limit: 100 })
    allRoles.value = res.items || []
  } catch (e: any) {
    console.error(e)
  }
}

const loadAllGroups = async () => {
  try {
    const res = await getGroups({ limit: 100 })
    allGroups.value = res.items || []
  } catch (e: any) {
    console.error(e)
  }
}

const loadAllDomains = async () => {
  try {
    const res = await getDomains({ limit: 100 })
    allDomains.value = res.items || []
  } catch (e: any) {
    console.error(e)
  }
}

const resetFilter = () => {
  filterForm.name = ''
  filterForm.email = ''
  filterForm.enabled = undefined
  pagination.page = 1
  loadUsers()
}

const loadAllProjectsAndRoles = async (domainId: number) => {
  try {
    const [projectsRes, rolesRes] = await Promise.all([
      getProjects({ domain_id: domainId, limit: 100 }),
      getRoles({ domain_id: domainId, limit: 100 })
    ])
    allProjects.value = projectsRes.items || []
    allProjectRoles.value = rolesRes.items || []
  } catch (e: any) {
    console.error(e)
    ElMessage.error(e.message || '加载项目/角色列表失败')
  }
}

const handleAddUserToProjectDomainChange = (domainId: number) => {
  addUserToProjectForm.project_ids = []
  addUserToProjectForm.role_ids = []
  loadAllProjectsAndRoles(domainId)
}

const handleCreate = async () => {
  isEdit.value = false
  form.name = ''
  form.display_name = ''
  form.remark = ''
  form.email = ''
  form.phone = ''
  form.password = ''
  form.domain_id = 1
  form.console_login = true
  form.mfa_enabled = false
  await loadAllDomains()
  await loadAllProjectsAndRoles(1) // Load projects and roles for default domain
  addUserToProjectForm.domain_id = 1
  addUserToProjectForm.project_ids = []
  addUserToProjectForm.role_ids = []
  dialogVisible.value = true
}

const handleEdit = (row: User) => {
  isEdit.value = true
  currentUserId.value = row.id
  form.name = row.name
  form.display_name = row.display_name
  form.remark = row.remark || ''
  form.email = row.email
  form.phone = row.phone || ''
  form.domain_id = row.domain_id
  form.console_login = row.console_login
  form.mfa_enabled = row.mfa_enabled
  dialogVisible.value = true
}

const handleView = async (row: User) => {
  currentUser.value = row
  detailTab.value = 'roles'
  await Promise.all([
    loadUserRoles(row.id),
    loadUserGroups(row.id),
    loadUserProjects(row.id),
    loadUserLogs(row.id)
  ])
  detailDialogVisible.value = true
}

const handleResetPassword = (row: User) => {
  currentUserId.value = row.id
  resetPwdForm.password = ''
  resetPwdDialogVisible.value = true
}

const handleToggleEnable = async (row: User) => {
  try {
    const action = row.enabled ? '禁用' : '启用'
    await ElMessageBox.confirm(`确定要${action}该用户吗？`, '提示', { type: 'warning' })

    if (row.enabled) {
      await disableUser(row.id)
    } else {
      await enableUser(row.id)
    }

    ElMessage.success(`${action}成功`)
    loadUsers()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || `${action}失败`)
    }
  }
}

const handleDelete = async (row: User) => {
  try {
    await ElMessageBox.confirm('确定要删除该用户吗？', '提示', { type: 'warning' })
    await deleteUser(row.id)
    ElMessage.success('删除成功')
    loadUsers()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (isEdit.value) {
        // For edit, only update the fields that are in the update payload
        const updateData = {
          display_name: form.display_name,
          remark: form.remark,
          email: form.email,
          phone: form.phone,
          console_login: form.console_login,
          mfa_enabled: form.mfa_enabled
        };
        await updateUser(currentUserId.value, updateData)
        ElMessage.success('更新成功')
      } else {
        // For create, use the full form data
        const newUser = await createUser(form)

        // 如果展开折叠面板且选择了项目和角色，则将用户添加到项目
        if (activeCollapse.value.includes('addProject') &&
            addUserToProjectForm.project_ids.length > 0 &&
            addUserToProjectForm.role_ids.length > 0) {

          // 为每个项目分配每个角色（可能需要调整策略）
          for (const projectId of addUserToProjectForm.project_ids) {
            for (const roleId of addUserToProjectForm.role_ids) {
              await assignUserToProject(newUser.id, projectId, roleId)
            }
          }
          ElMessage.success('创建成功并添加到项目')
        } else {
          ElMessage.success('创建成功')
        }
      }
      dialogVisible.value = false
      loadUsers()
    } catch (e: any) {
      ElMessage.error(e.message || (isEdit.value ? '更新失败' : '创建失败'))
    } finally {
      submitting.value = false
    }
  })
}

const handleResetPasswordSubmit = async () => {
  if (!resetPwdForm.password || resetPwdForm.password.length < 8) {
    ElMessage.error('密码长度至少 8 位')
    return
  }

  resetPwdSubmitting.value = true
  try {
    await resetUserPassword(currentUserId.value, resetPwdForm.password)
    ElMessage.success('密码重置成功')
    resetPwdDialogVisible.value = false
  } catch (e: any) {
    ElMessage.error(e.message || '密码重置失败')
  } finally {
    resetPwdSubmitting.value = false
  }
}

const handleAssignRole = async () => {
  await loadAllRoles()
  assignRoleForm.role_id = 0
  assignRoleForm.domain_id = 1
  assignRoleDialogVisible.value = true
}

const handleAssignRoleSubmit = async () => {
  if (!assignRoleForm.role_id) {
    ElMessage.error('请选择角色')
    return
  }

  assignRoleSubmitting.value = true
  try {
    await assignRoleToUser(currentUserId.value, assignRoleForm.role_id, assignRoleForm.domain_id)
    ElMessage.success('角色分配成功')
    assignRoleDialogVisible.value = false
    await loadUserRoles(currentUserId.value)
  } catch (e: any) {
    ElMessage.error(e.message || '角色分配失败')
  } finally {
    assignRoleSubmitting.value = false
  }
}

const handleRevokeRole = async (row: Role) => {
  try {
    await ElMessageBox.confirm('确定要撤销该角色吗？', '提示', { type: 'warning' })
    await revokeRoleFromUser(currentUserId.value, row.id, 1)
    ElMessage.success('角色撤销成功')
    await loadUserRoles(currentUserId.value)
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '角色撤销失败')
    }
  }
}

const handleJoinGroup = async () => {
  await loadAllGroups()
  joinGroupForm.group_id = 0
  joinGroupDialogVisible.value = true
}

const handleJoinGroupSubmit = async () => {
  if (!joinGroupForm.group_id) {
    ElMessage.error('请选择用户组')
    return
  }

  joinGroupSubmitting.value = true
  try {
    await joinGroup(currentUserId.value, joinGroupForm.group_id)
    ElMessage.success('加入用户组成功')
    joinGroupDialogVisible.value = false
    await loadUserGroups(currentUserId.value)
  } catch (e: any) {
    ElMessage.error(e.message || '加入用户组失败')
  } finally {
    joinGroupSubmitting.value = false
  }
}

const handleLeaveGroup = async (row: Group) => {
  try {
    await ElMessageBox.confirm('确定要离开该用户组吗？', '提示', { type: 'warning' })
    await leaveGroup(currentUserId.value, row.id)
    ElMessage.success('离开用户组成功')
    await loadUserGroups(currentUserId.value)
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '离开用户组失败')
    }
  }
}

const handleCreateRole = () => {
  // This would navigate to the role creation page
  // For now, we'll just show a notification
  ElMessage.info('角色创建功能将在后续实现')
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const getDomainName = (domainId: number) => {
  const domain = allDomains.value.find(d => d.id === domainId);
  return domain ? domain.name : `域#${domainId}`;
}

const handleCommand = (command: string, row: User) => {
  switch (command) {
    case 'enable': handleToggleEnable({...row, enabled: false}); break
    case 'disable': handleToggleEnable({...row, enabled: true}); break
    case 'resetPassword': handleResetPassword(row); break
    case 'manageProjects': handleManageProjects(row); break
    case 'resetMFA': handleResetMFA(row); break
    case 'delete': handleDelete(row); break
  }
}

const handleManageProjects = async (row: User) => {
  // This will open the project management dialog
  currentUser.value = row
  // Load user's current projects and show in a dialog
  try {
    const res = await getUserProjects(row.id)
    ElMessage.info('功能将在后续实现')
  } catch (e: any) {
    ElMessage.error(e.message || '加载项目失败')
  }
}

const handleResetMFA = async (row: User) => {
  try {
    await ElMessageBox.confirm('确定要重置用户的MFA吗？', '提示', { type: 'warning' })
    ElMessage.info('重置MFA功能将在后续实现')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '重置MFA失败')
    }
  }
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.users-page {
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
  padding: 0;
  border: none;
  background: none;
}
</style>
