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
        <el-table-column prop="name" label="用户名" min-width="150" show-overflow-tooltip />
        <el-table-column prop="display_name" label="显示名" min-width="120" show-overflow-tooltip />
        <el-table-column prop="email" label="邮箱" min-width="180" show-overflow-tooltip />
        <el-table-column prop="phone" label="手机号" width="120" />
        <el-table-column prop="enabled" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="mfa_enabled" label="MFA" width="60">
          <template #default="{ row }">
            <el-tag :type="row.mfa_enabled ? 'success' : 'info'" size="small">
              {{ row.mfa_enabled ? '开' : '关' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="320" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">详情</el-button>
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleResetPassword(row)">重置密码</el-button>
            <el-button 
              size="small" 
              :type="row.enabled ? 'warning' : 'success'" 
              @click="handleToggleEnable(row)"
            >
              {{ row.enabled ? '禁用' : '启用' }}
            </el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
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
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑用户' : '新建用户'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="用户名" prop="name">
          <el-input v-model="form.name" placeholder="请输入用户名" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="显示名" prop="display_name">
          <el-input v-model="form.display_name" placeholder="请输入显示名" />
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
            <el-option label="Default" :value="1" />
          </el-select>
        </el-form-item>
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
        <el-descriptions-item label="邮箱">{{ currentUser.email || '-' }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ currentUser.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentUser.enabled ? 'success' : 'info'" size="small">
            {{ currentUser.enabled ? '启用' : '禁用' }}
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
import { User, Role, Group, Domain } from '@/types/iam'
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
  getDomains
} from '@/api/iam'

const users = ref<User[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const resetPwdDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const assignRoleDialogVisible = ref(false)
const joinGroupDialogVisible = ref(false)
const detailTab = ref('roles')
const isEdit = ref(false)
const submitting = ref(false)
const resetPwdSubmitting = ref(false)
const assignRoleSubmitting = ref(false)
const joinGroupSubmitting = ref(false)
const currentUserId = ref(0)
const currentUser = ref<User | null>(null)
const userRoles = ref<Role[]>([])
const userGroups = ref<Group[]>([])
const allRoles = ref<Role[]>([])
const allGroups = ref<Group[]>([])
const allDomains = ref<Domain[]>([])
const rolesLoading = ref(false)
const groupsLoading = ref(false)
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
  email: '',
  phone: '',
  password: '',
  domain_id: 1
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

const handleCreate = async () => {
  isEdit.value = false
  form.name = ''
  form.display_name = ''
  form.email = ''
  form.phone = ''
  form.password = ''
  form.domain_id = 1
  await loadAllDomains()
  dialogVisible.value = true
}

const handleEdit = (row: User) => {
  isEdit.value = true
  currentUserId.value = row.id
  form.name = row.name
  form.display_name = row.display_name
  form.email = row.email
  form.phone = row.phone || ''
  form.password = ''
  form.domain_id = row.domain_id
  dialogVisible.value = true
}

const handleView = async (row: User) => {
  currentUser.value = row
  detailTab.value = 'roles'
  await loadUserRoles(row.id)
  await loadUserGroups(row.id)
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
        await updateUser(currentUserId.value, form)
        ElMessage.success('更新成功')
      } else {
        await createUser(form)
        ElMessage.success('创建成功')
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

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
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
</style>
