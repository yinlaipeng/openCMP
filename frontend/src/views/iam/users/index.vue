<template>
  <div class="users-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">用户</span>
          <div class="header-icons">
            <el-button size="small" @click="handleDownload">
              <el-icon><Download /></el-icon>
            </el-button>
            <el-button size="small" @click="handleSettings">
              <el-icon><Setting /></el-icon>
            </el-button>
          </div>
        </div>
      </template>

      <!-- 工具栏 -->
      <div class="toolbar">
        <el-button size="small" @click="handleRefresh" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button size="small" type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          新建
        </el-button>
        <el-button size="small" @click="handleImportUsers">
          <el-icon><Upload /></el-icon>
          导入用户
        </el-button>
        <el-dropdown trigger="click" @command="handleBatchCommand" :disabled="selectedUsers.length === 0">
          <el-button size="small">
            批量操作 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="batchEnable">批量启用</el-dropdown-item>
              <el-dropdown-item command="batchDisable">批量禁用</el-dropdown-item>
              <el-dropdown-item command="batchResetPassword">批量重置密码</el-dropdown-item>
              <el-dropdown-item command="batchDelete" divided>
                <span style="color: #F56C6C">批量删除</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button size="small" @click="handleManageTags">
          <el-icon><PriceTag /></el-icon>
          标签
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
        <el-button type="primary" @click="handleSearch">
          <el-icon><Search /></el-icon>
          查询
        </el-button>
      </div>

      <!-- 用户列表 -->
      <el-table
        :data="users"
        v-loading="loading"
        @selection-change="handleSelectionChange"
        row-key="id"
        border
        stripe
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="名称" min-width="150" sortable show-overflow-tooltip>
          <template #default="{ row }">
            <div class="name-cell">
              <el-icon class="user-icon"><User /></el-icon>
              <el-button type="primary" link @click="handleView(row)" class="name-link">
                {{ row.name }}
              </el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="display_name" label="显示名" min-width="120" show-overflow-tooltip />
        <el-table-column label="标签" width="100">
          <template #default="{ row }">
            <el-icon v-if="row.tags?.length"><PriceTag /></el-icon>
            <span v-else class="empty-text">-</span>
          </template>
        </el-table-column>
        <el-table-column label="启用状态" width="100">
          <template #default="{ row }">
            <span class="status-dot" :data-status="row.enabled ? 'enabled' : 'disabled'">
              {{ row.enabled ? '启用' : '禁用' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="控制台登录" width="100">
          <template #default="{ row }">
            <span class="status-dot" :data-status="row.console_login ? 'enabled' : 'disabled'">
              {{ row.console_login ? '允许' : '禁止' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="MFA" width="80">
          <template #default="{ row }">
            <span class="status-dot" :data-status="row.mfa_enabled ? 'enabled' : 'disabled'">
              {{ row.mfa_enabled ? '开' : '关' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="所属域" width="120" sortable>
          <template #default="{ row }">
            <span>{{ getDomainName(row.domain_id) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              修改属性
            </el-button>
            <el-dropdown trigger="click" @command="(cmd) => handleCommand(cmd, row)">
              <el-button size="small" type="primary" link>
                更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="enable" v-if="!row.enabled">启用</el-dropdown-item>
                  <el-dropdown-item command="disable" v-if="row.enabled">禁用</el-dropdown-item>
                  <el-dropdown-item command="resetPassword">重置密码</el-dropdown-item>
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

    <!-- 添加用户对话框 -->
    <el-dialog v-model="dialogVisible" title="新建用户" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="用户名" prop="name">
          <el-input v-model="form.name" placeholder="请输入用户名" />
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
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码（至少 8 位）" show-password />
        </el-form-item>
        <el-form-item label="所属域" prop="domain_id">
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
        <el-collapse v-model="activeCollapse">
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

    <!-- 修改属性弹窗 -->
    <ModifyAttributesModal
      v-model="modifyModalVisible"
      :user="currentUser"
      @success="handleModifySuccess"
    />

    <!-- 重置密码弹窗 -->
    <ResetPasswordModal
      v-model="resetPwdModalVisible"
      :user="currentUser"
      @success="handleResetPwdSuccess"
    />

    <!-- 重置MFA弹窗 -->
    <ResetMFAModal
      v-model="resetMFAModalVisible"
      :user="currentUser"
      @success="handleResetMFASuccess"
    />

    <!-- 用户详情抽屉 -->
    <UserDetailDrawer
      v-model="detailDrawerVisible"
      :user-id="currentUserId"
      @refresh="loadUsers"
      @deleted="loadUsers"
    />

    <!-- 导入用户弹窗 -->
    <ImportUsersModal
      v-model="importModalVisible"
      @success="loadUsers"
    />

    <!-- 批量操作弹窗 -->
    <el-dialog v-model="batchModalVisible" :title="batchModalTitle" width="500px">
      <el-alert
        v-if="batchOperation === 'batchDelete'"
        title="警告"
        description="批量删除操作不可恢复，请谨慎操作"
        type="error"
        show-icon
        :closable="false"
        style="margin-bottom: 20px"
      />
      <p class="batch-confirm-text">
        确定要对选中的 <strong>{{ selectedUsers.length }}</strong> 个用户执行{{ batchModalTitle }}操作吗？
      </p>
      <el-table :data="selectedUsers" border size="small" max-height="200">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="用户名" />
        <el-table-column prop="display_name" label="显示名" />
      </el-table>
      <template #footer>
        <el-button @click="batchModalVisible = false">取消</el-button>
        <el-button
          :type="batchOperation === 'batchDelete' ? 'danger' : 'primary'"
          @click="handleBatchSubmit"
          :loading="batchLoading"
        >
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Refresh, Plus, Upload, ArrowDown, Download, Setting,
  PriceTag, Search, User
} from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import type { User as UserType, Domain, Project, Role } from '@/types/iam'
import {
  getUsers, createUser, deleteUser, enableUser, disableUser,
  getDomains, getProjects, getRoles, assignUserToProject, exportUsers
} from '@/api/iam'

// 导入新组件
import UserDetailDrawer from './components/UserDetailDrawer.vue'
import ModifyAttributesModal from './components/ModifyAttributesModal.vue'
import ResetPasswordModal from './components/ResetPasswordModal.vue'
import ResetMFAModal from './components/ResetMFAModal.vue'
import ImportUsersModal from './components/ImportUsersModal.vue'

// State
const users = ref<UserType[]>([])
const loading = ref(false)
const selectedUsers = ref<UserType[]>([])
const allDomains = ref<Domain[]>([])
const allProjects = ref<Project[]>([])
const allProjectRoles = ref<Role[]>([])

// Dialog visibility
const dialogVisible = ref(false)
const modifyModalVisible = ref(false)
const resetPwdModalVisible = ref(false)
const resetMFAModalVisible = ref(false)
const detailDrawerVisible = ref(false)
const importModalVisible = ref(false)
const batchModalVisible = ref(false)

// Current user
const currentUserId = ref(0)
const currentUser = ref<UserType | null>(null)

// Form
const submitting = ref(false)
const formRef = ref<FormInstance>()
const searchKeyword = ref('')
const activeCollapse = ref<string[]>([])

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

const addUserToProjectForm = reactive({
  domain_id: 1,
  project_ids: [] as number[],
  role_ids: [] as number[]
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  display_name: [{ required: true, message: '请输入显示名', trigger: 'blur' }],
  email: [{ type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }],
  password: [{ min: 8, message: '密码长度至少 8 位', trigger: 'blur' }],
  domain_id: [{ required: true, message: '请选择所属域', trigger: 'change' }]
}

// Batch operations
const batchOperation = ref('')
const batchLoading = ref(false)

const batchModalTitle = computed(() => {
  switch (batchOperation.value) {
    case 'batchEnable': return '批量启用'
    case 'batchDisable': return '批量禁用'
    case 'batchResetPassword': return '批量重置密码'
    case 'batchDelete': return '批量删除'
    default: return ''
  }
})

// Load data
const loadUsers = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit
    }
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }

    const res = await getUsers(params)
    users.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载用户列表失败')
  } finally {
    loading.value = false
  }
}

const loadDomains = async () => {
  try {
    const res = await getDomains({ limit: 100 })
    allDomains.value = res.items || []
  } catch (e: any) {
    console.error(e)
  }
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

// Handlers
const handleRefresh = () => loadUsers()

const handleSearch = () => {
  pagination.page = 1
  loadUsers()
}

const handleSelectionChange = (selection: UserType[]) => {
  selectedUsers.value = selection
}

const handleCreate = async () => {
  await loadDomains()
  await loadAllProjectsAndRoles(1)
  form.name = ''
  form.display_name = ''
  form.remark = ''
  form.email = ''
  form.phone = ''
  form.password = ''
  form.domain_id = 1
  form.console_login = true
  form.mfa_enabled = false
  addUserToProjectForm.domain_id = 1
  addUserToProjectForm.project_ids = []
  addUserToProjectForm.role_ids = []
  activeCollapse.value = []
  dialogVisible.value = true
}

const handleImportUsers = () => {
  importModalVisible.value = true
}

const handleManageTags = () => {
  ElMessage.info('标签管理功能将在后续实现')
}

const handleDownload = async () => {
  try {
    loading.value = true
    const res = await exportUsers()
    const exportData = res.items || []

    // 生成 CSV 内容
    const headers = ['ID', '用户名', '显示名', '邮箱', '手机号', '启用状态', '控制台登录', 'MFA', '所属域ID', '备注', '创建时间']
    const csvRows = [headers.join(',')]

    for (const user of exportData) {
      const row = [
        user.id,
        user.name,
        user.display_name || '',
        user.email || '',
        user.phone || '',
        user.enabled ? '启用' : '禁用',
        user.console_login ? '允许' : '禁止',
        user.mfa_enabled ? '开启' : '关闭',
        user.domain_id,
        user.remark || '',
        new Date(user.created_at).toLocaleString('zh-CN')
      ]
      csvRows.push(row.join(','))
    }

    const csvContent = csvRows.join('\n')

    // 创建下载文件
    const blob = new Blob(['\ufeff' + csvContent], { type: 'text/csv;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `users_export_${new Date().toISOString().slice(0,10)}.csv`
    link.click()
    URL.revokeObjectURL(url)

    ElMessage.success(`成功导出 ${exportData.length} 个用户`)
  } catch (e: any) {
    ElMessage.error(e.message || '导出失败')
  } finally {
    loading.value = false
  }
}

const handleSettings = () => {
  ElMessage.info('设置功能将在后续实现')
}

const handleView = (row: UserType) => {
  currentUserId.value = row.id
  currentUser.value = row
  detailDrawerVisible.value = true
}

const handleEdit = (row: UserType) => {
  currentUser.value = row
  modifyModalVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      const newUser = await createUser(form)

      if (activeCollapse.value.includes('addProject') &&
          addUserToProjectForm.project_ids.length > 0 &&
          addUserToProjectForm.role_ids.length > 0) {
        for (const projectId of addUserToProjectForm.project_ids) {
          for (const roleId of addUserToProjectForm.role_ids) {
            await assignUserToProject(newUser.id, projectId, roleId)
          }
        }
        ElMessage.success('创建成功并添加到项目')
      } else {
        ElMessage.success('创建成功')
      }

      dialogVisible.value = false
      loadUsers()
    } catch (e: any) {
      ElMessage.error(e.message || '创建失败')
    } finally {
      submitting.value = false
    }
  })
}

const handleCommand = async (command: string, row: UserType) => {
  currentUser.value = row
  currentUserId.value = row.id

  switch (command) {
    case 'enable':
      await handleToggleEnable(row, true)
      break
    case 'disable':
      await handleToggleEnable(row, false)
      break
    case 'resetPassword':
      resetPwdModalVisible.value = true
      break
    case 'resetMFA':
      resetMFAModalVisible.value = true
      break
    case 'delete':
      await handleDelete(row)
      break
  }
}

const handleToggleEnable = async (row: UserType, enable: boolean) => {
  try {
    const action = enable ? '启用' : '禁用'
    await ElMessageBox.confirm(`确定要${action}该用户吗？`, '提示', { type: 'warning' })
    if (enable) {
      await enableUser(row.id)
    } else {
      await disableUser(row.id)
    }
    ElMessage.success(`${action}成功`)
    loadUsers()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '操作失败')
    }
  }
}

const handleDelete = async (row: UserType) => {
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

// Batch handlers
const handleBatchCommand = (command: string) => {
  batchOperation.value = command
  batchModalVisible.value = true
}

const handleBatchSubmit = async () => {
  batchLoading.value = true
  try {
    const userIds = selectedUsers.value.map(u => u.id)

    switch (batchOperation.value) {
      case 'batchEnable':
        for (const id of userIds) {
          await enableUser(id)
        }
        ElMessage.success(`成功启用 ${userIds.length} 个用户`)
        break
      case 'batchDisable':
        for (const id of userIds) {
          await disableUser(id)
        }
        ElMessage.success(`成功禁用 ${userIds.length} 个用户`)
        break
      case 'batchResetPassword':
        ElMessage.info('批量重置密码功能需要后端API支持')
        break
      case 'batchDelete':
        for (const id of userIds) {
          await deleteUser(id)
        }
        ElMessage.success(`成功删除 ${userIds.length} 个用户`)
        break
    }

    batchModalVisible.value = false
    loadUsers()
  } catch (e: any) {
    ElMessage.error(e.message || '批量操作失败')
  } finally {
    batchLoading.value = false
  }
}

// Helper handlers
const handleAddUserToProjectDomainChange = (domainId: number) => {
  addUserToProjectForm.project_ids = []
  addUserToProjectForm.role_ids = []
  loadAllProjectsAndRoles(domainId)
}

const handleCreateRole = () => {
  ElMessage.info('角色创建功能将在角色管理页面实现')
}

const handleModifySuccess = () => {
  loadUsers()
}

const handleResetPwdSuccess = () => {
  ElMessage.success('密码重置成功')
}

const handleResetMFASuccess = () => {
  loadUsers()
}

const getDomainName = (domainId: number) => {
  const domain = allDomains.value.find(d => d.id === domainId)
  return domain ? domain.name : `域#${domainId}`
}

onMounted(() => {
  loadUsers()
  loadDomains()
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

.header-icons {
  display: flex;
  gap: 8px;
}

.toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.search-bar {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.search-input {
  width: 400px;
  max-width: 100%;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-icon {
  color: var(--color-muted, #64748B);
}

.name-link {
  padding: 0;
  border: none;
  background: none;
}

.status-dot {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.status-dot::before {
  content: '';
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-dot[data-status="enabled"]::before {
  background: var(--color-success, #22C55E);
}

.status-dot[data-status="disabled"]::before {
  background: var(--color-muted, #64748B);
}

.empty-text {
  color: var(--color-muted, #64748B);
}

.batch-confirm-text {
  margin-bottom: 16px;
}

.batch-confirm-text strong {
  color: var(--color-primary, #0F172A);
}
</style>