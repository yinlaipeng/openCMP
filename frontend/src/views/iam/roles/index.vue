<template>
  <div class="roles-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">角色</span>
          <!-- 工具栏 -->
          <div class="toolbar">
            <el-button @click="handleRefresh" :loading="loading">
              <el-icon><Refresh /></el-icon>
            </el-button>
            <el-button type="primary" @click="handleCreate">
              <el-icon><Plus /></el-icon>
              新建
            </el-button>
            <el-button :disabled="selectedRoles.length === 0" @click="handleBatchDelete">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
            <div class="toolbar-icons">
              <el-tooltip content="下载" placement="top">
                <el-button @click="handleDownload">
                  <el-icon><Download /></el-icon>
                </el-button>
              </el-tooltip>
              <el-tooltip content="设置" placement="top">
                <el-button @click="handleSettings">
                  <el-icon><Setting /></el-icon>
                </el-button>
              </el-tooltip>
            </div>
          </div>
        </div>
      </template>

      <!-- 轻量化搜索栏 -->
      <div class="search-bar">
        <el-dropdown trigger="click" @command="handleFieldChange">
          <el-button>
            {{ currentFieldLabel }} <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="name">名称</el-dropdown-item>
              <el-dropdown-item command="id">ID</el-dropdown-item>
              <el-dropdown-item command="type">类型</el-dropdown-item>
              <el-dropdown-item command="enabled">状态</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-input
          v-model="searchKeyword"
          :placeholder="searchPlaceholder"
          clearable
          style="width: 300px"
          @keyup.enter="loadRoles"
        >
          <template #suffix>
            <el-icon class="search-icon" @click="loadRoles"><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="loadRoles">查询</el-button>
        <el-button @click="handleResetSearch">重置</el-button>
      </div>

      <!-- 角色列表 -->
      <el-table
        :data="roles"
        v-loading="loading"
        border
        stripe
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" sortable />
        <el-table-column label="名称" min-width="180" show-overflow-tooltip sortable>
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)" class="name-link">
              <span>{{ row.name }}</span>
            </el-button>
            <br>
            <small style="color: #999;">{{ row.display_name || row.description || '-' }}</small>
          </template>
        </el-table-column>
        <el-table-column label="策略" min-width="150" show-overflow-tooltip>
          <template #default="{ row }">
            <el-tag v-if="row.policies_count > 0" type="info" size="small">
              {{ row.policies_count }} 个策略
            </el-tag>
            <span v-else style="color: #999;">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="100" sortable>
          <template #default="{ row }">
            <el-tag :type="row.type === 'system' ? 'warning' : 'success'" size="small">
              {{ row.type === 'system' ? '系统' : '自定义' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="handleSetPolicies(row)">
              设置策略
            </el-button>
            <el-dropdown trigger="click" @command="(cmd: string) => handleMoreCommand(cmd, row)">
              <el-button size="small" type="primary" link>
                更多 <el-icon><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="view">详情</el-dropdown-item>
                  <el-dropdown-item command="edit" :disabled="row.type === 'system'">编辑</el-dropdown-item>
                  <el-dropdown-item command="enable" v-if="!row.enabled" :disabled="row.type === 'system'">启用</el-dropdown-item>
                  <el-dropdown-item command="disable" v-if="row.enabled" :disabled="row.type === 'system'">禁用</el-dropdown-item>
                  <el-dropdown-item command="public" :disabled="row.type === 'system' || row.is_public">公开</el-dropdown-item>
                  <el-dropdown-item command="delete" :disabled="row.type === 'system'" divided>删除</el-dropdown-item>
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
        @size-change="loadRoles"
        @current-change="loadRoles"
        class="pagination"
      />
    </el-card>

    <!-- 添加/编辑角色对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑角色' : '新建角色'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入角色名（英文）" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="显示名" prop="display_name">
          <el-input v-model="form.display_name" placeholder="请输入显示名（中文）" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入角色描述" />
        </el-form-item>
        <el-form-item label="类型" prop="type" v-if="!isEdit">
          <el-radio-group v-model="form.type">
            <el-radio label="custom">自定义</el-radio>
            <el-radio label="system">系统</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 策略分配对话框 -->
    <el-dialog v-model="policyDialogVisible" title="设置策略" width="800px">
      <div v-if="currentRoleForPolicy" class="policy-dialog-header">
        <span>角色: <strong>{{ currentRoleForPolicy.name }}</strong></span>
      </div>

      <!-- 已分配的策略 -->
      <div class="assigned-policies">
        <h4>已分配策略</h4>
        <el-table :data="assignedPolicies" v-loading="policiesLoading" border stripe max-height="200">
          <el-table-column prop="id" label="ID" width="200" show-overflow-tooltip />
          <el-table-column prop="name" label="名称" min-width="150" show-overflow-tooltip />
          <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-button size="small" type="danger" link @click="handleRevokePolicy(row)">
                移除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 可分配的策略 -->
      <div class="available-policies" style="margin-top: 20px;">
        <h4>可分配策略</h4>
        <el-table :data="availablePolicies" v-loading="allPoliciesLoading" border stripe max-height="300">
          <el-table-column prop="id" label="ID" width="200" show-overflow-tooltip />
          <el-table-column prop="name" label="名称" min-width="150" show-overflow-tooltip />
          <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-button
                size="small"
                type="primary"
                link
                @click="handleAssignPolicy(row)"
                :disabled="isPolicyAssigned(row.id)"
              >
                添加
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <template #footer>
        <el-button @click="policyDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 角色详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="角色详情" width="700px">
      <el-descriptions :column="2" border v-if="currentRole">
        <el-descriptions-item label="ID">{{ currentRole.id }}</el-descriptions-item>
        <el-descriptions-item label="角色名">{{ currentRole.name }}</el-descriptions-item>
        <el-descriptions-item label="显示名">{{ currentRole.display_name }}</el-descriptions-item>
        <el-descriptions-item label="类型">
          <el-tag :type="currentRole.type === 'system' ? 'warning' : 'success'" size="small">
            {{ currentRole.type === 'system' ? '系统' : '自定义' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentRole.enabled ? 'success' : 'info'" size="small">
            {{ currentRole.enabled ? '启用' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="公开">
          <el-tag :type="currentRole.is_public ? 'primary' : 'info'" size="small">
            {{ currentRole.is_public ? '是' : '否' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentRole.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentRole.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentRole.updated_at) }}</el-descriptions-item>
      </el-descriptions>

      <el-tabs v-model="detailTab" style="margin-top: 20px">
        <el-tab-pane label="关联用户" name="users">
          <el-table :data="roleUsers" v-loading="usersLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="用户名" />
            <el-table-column prop="display_name" label="显示名" />
            <el-table-column prop="email" label="邮箱" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="关联用户组" name="groups">
          <el-table :data="roleGroups" v-loading="groupsLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="用户组名" />
            <el-table-column prop="description" label="描述" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="关联策略" name="policies">
          <el-table :data="rolePolicies" v-loading="policiesLoading">
            <el-table-column prop="id" label="ID" width="200" show-overflow-tooltip />
            <el-table-column prop="name" label="名称" />
            <el-table-column prop="description" label="描述" show-overflow-tooltip />
          </el-table>
        </el-tab-pane>
      </el-tabs>

      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 批量删除确认弹窗 -->
    <el-dialog v-model="batchDeleteVisible" title="批量删除确认" width="500px">
      <p>确定要删除以下 {{ selectedRoles.length }} 个角色吗？</p>
      <el-table :data="selectedRoles" border stripe max-height="300">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'system' ? 'warning' : 'success'" size="small">
              {{ row.type === 'system' ? '系统' : '自定义' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="batchDeleteVisible = false">取消</el-button>
        <el-button type="danger" @click="handleBatchDeleteConfirm" :loading="batchDeleting">删除</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Refresh, Plus, Delete, Download, Setting, ArrowDown, Search } from '@element-plus/icons-vue'
import { Role } from '@/types/iam'
import {
  getRoles,
  createRole,
  updateRole,
  deleteRole,
  enableRole,
  disableRole,
  makeRolePublic,
  getRoleUsers,
  getRoleGroups,
  getRolePolicies,
  assignPolicyToRole,
  revokePolicyFromRole,
  batchDeleteRoles,
  getPolicies
} from '@/api/iam'

interface Policy {
  id: string
  name: string
  description?: string
}

interface RoleWithPolicies extends Role {
  policies_count?: number
}

const roles = ref<RoleWithPolicies[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const policyDialogVisible = ref(false)
const batchDeleteVisible = ref(false)
const detailTab = ref('users')
const isEdit = ref(false)
const submitting = ref(false)
const batchDeleting = ref(false)
const currentRoleId = ref(0)
const currentRole = ref<Role | null>(null)
const currentRoleForPolicy = ref<Role | null>(null)
const roleUsers = ref<any[]>([])
const roleGroups = ref<any[]>([])
const rolePolicies = ref<Policy[]>([])
const assignedPolicies = ref<Policy[]>([])
const availablePolicies = ref<Policy[]>([])
const selectedRoles = ref<Role[]>([])
const usersLoading = ref(false)
const groupsLoading = ref(false)
const policiesLoading = ref(false)
const allPoliciesLoading = ref(false)
const formRef = ref<FormInstance>()

// 搜索相关
const searchField = ref('name')
const searchKeyword = ref('')

const currentFieldLabel = computed(() => {
  const labels: Record<string, string> = {
    name: '名称',
    id: 'ID',
    type: '类型',
    enabled: '状态'
  }
  return labels[searchField.value] || '名称'
})

const searchPlaceholder = computed(() => {
  const placeholders: Record<string, string> = {
    name: '请输入角色名称...',
    id: '请输入角色ID...',
    type: '请选择类型：system/custom...',
    enabled: '请输入状态：true/false...'
  }
  return placeholders[searchField.value] || '请输入搜索关键词...'
})

const filterForm = reactive({
  name: '',
  type: '',
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
  description: '',
  type: 'custom'
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入角色名', trigger: 'blur' }]
}

const handleFieldChange = (field: string) => {
  searchField.value = field
  searchKeyword.value = ''
}

const loadRoles = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit
    }

    // 根据搜索字段添加参数
    if (searchKeyword.value) {
      if (searchField.value === 'name') {
        params.keyword = searchKeyword.value
      } else if (searchField.value === 'type') {
        params.type = searchKeyword.value
      } else if (searchField.value === 'enabled') {
        params.enabled = searchKeyword.value === 'true'
      } else if (searchField.value === 'id') {
        // ID 搜索暂不支持
        params.keyword = searchKeyword.value
      }
    }

    const res = await getRoles(params)
    roles.value = res.items || []
    pagination.total = res.total || 0

    // 加载每个角色的策略数量
    for (const role of roles.value) {
      try {
        const policiesRes = await getRolePolicies(role.id)
        role.policies_count = policiesRes.total || 0
      } catch {
        role.policies_count = 0
      }
    }
  } catch (e: any) {
    ElMessage.error(e.message || '加载角色列表失败')
  } finally {
    loading.value = false
  }
}

const loadRoleUsers = async (roleId: number) => {
  usersLoading.value = true
  try {
    const res = await getRoleUsers(roleId, { limit: 100 })
    roleUsers.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    usersLoading.value = false
  }
}

const loadRoleGroups = async (roleId: number) => {
  groupsLoading.value = true
  try {
    const res = await getRoleGroups(roleId, { limit: 100 })
    roleGroups.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    groupsLoading.value = false
  }
}

const loadRolePolicies = async (roleId: number) => {
  policiesLoading.value = true
  try {
    const res = await getRolePolicies(roleId)
    rolePolicies.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    policiesLoading.value = false
  }
}

const loadAllPolicies = async () => {
  allPoliciesLoading.value = true
  try {
    const res = await getPolicies({ limit: 100 })
    availablePolicies.value = res.data || []
  } catch (e: any) {
    console.error(e)
  } finally {
    allPoliciesLoading.value = false
  }
}

const isPolicyAssigned = (policyId: string) => {
  return assignedPolicies.value.some(p => p.id === policyId)
}

const handleRefresh = () => {
  loadRoles()
}

const handleResetSearch = () => {
  searchKeyword.value = ''
  searchField.value = 'name'
  pagination.page = 1
  loadRoles()
}

const handleSelectionChange = (selection: Role[]) => {
  selectedRoles.value = selection
}

const handleCreate = async () => {
  isEdit.value = false
  form.name = ''
  form.display_name = ''
  form.description = ''
  form.type = 'custom'
  dialogVisible.value = true
}

const handleEdit = (row: Role) => {
  isEdit.value = true
  currentRoleId.value = row.id
  form.name = row.name
  form.display_name = row.display_name || ''
  form.description = row.description || ''
  form.type = row.type
  dialogVisible.value = true
}

const handleView = async (row: Role) => {
  currentRole.value = row
  detailTab.value = 'users'
  await loadRoleUsers(row.id)
  await loadRoleGroups(row.id)
  await loadRolePolicies(row.id)
  detailDialogVisible.value = true
}

const handleSetPolicies = async (row: Role) => {
  currentRoleForPolicy.value = row
  await loadRolePolicies(row.id)
  assignedPolicies.value = [...rolePolicies.value]
  await loadAllPolicies()
  policyDialogVisible.value = true
}

const handleAssignPolicy = async (policy: Policy) => {
  if (!currentRoleForPolicy.value) return

  try {
    await assignPolicyToRole(currentRoleForPolicy.value.id, policy.id)
    ElMessage.success('策略分配成功')
    assignedPolicies.value.push(policy)
    // 更新角色的策略数量
    const role = roles.value.find(r => r.id === currentRoleForPolicy.value!.id)
    if (role) {
      role.policies_count = (role.policies_count || 0) + 1
    }
  } catch (e: any) {
    ElMessage.error(e.message || '策略分配失败')
  }
}

const handleRevokePolicy = async (policy: Policy) => {
  if (!currentRoleForPolicy.value) return

  try {
    await revokePolicyFromRole(currentRoleForPolicy.value.id, policy.id)
    ElMessage.success('策略移除成功')
    assignedPolicies.value = assignedPolicies.value.filter(p => p.id !== policy.id)
    // 更新角色的策略数量
    const role = roles.value.find(r => r.id === currentRoleForPolicy.value!.id)
    if (role) {
      role.policies_count = Math.max(0, (role.policies_count || 1) - 1)
    }
  } catch (e: any) {
    ElMessage.error(e.message || '策略移除失败')
  }
}

const handleMoreCommand = (command: string, row: Role) => {
  switch (command) {
    case 'view':
      handleView(row)
      break
    case 'edit':
      handleEdit(row)
      break
    case 'enable':
      handleToggleEnable(row)
      break
    case 'disable':
      handleToggleEnable(row)
      break
    case 'public':
      handleMakePublic(row)
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

const handleToggleEnable = async (row: Role) => {
  try {
    const action = row.enabled ? '禁用' : '启用'
    await ElMessageBox.confirm(`确定要${action}该角色吗？`, '提示', { type: 'warning' })

    if (row.enabled) {
      await disableRole(row.id)
    } else {
      await enableRole(row.id)
    }

    ElMessage.success(`${action}成功`)
    loadRoles()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || `${row.enabled ? '禁用' : '启用'}失败`)
    }
  }
}

const handleMakePublic = async (row: Role) => {
  try {
    await ElMessageBox.confirm('确定要公开该角色吗？公开后其他域的用户也可以使用此角色。', '提示', { type: 'warning' })
    await makeRolePublic(row.id)
    ElMessage.success('公开成功')
    loadRoles()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '公开失败')
    }
  }
}

const handleDelete = async (row: Role) => {
  try {
    await ElMessageBox.confirm('确定要删除该角色吗？', '提示', { type: 'warning' })
    await deleteRole(row.id)
    ElMessage.success('删除成功')
    loadRoles()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const handleBatchDelete = () => {
  if (selectedRoles.value.length === 0) {
    ElMessage.warning('请先选择要删除的角色')
    return
  }
  // 检查是否包含系统角色
  const systemRoles = selectedRoles.value.filter(r => r.type === 'system')
  if (systemRoles.length > 0) {
    ElMessage.warning('不能删除系统角色，请重新选择')
    return
  }
  batchDeleteVisible.value = true
}

const handleBatchDeleteConfirm = async () => {
  batchDeleting.value = true
  try {
    const roleIds = selectedRoles.value.map(r => r.id)
    await batchDeleteRoles(roleIds)
    ElMessage.success(`成功删除 ${selectedRoles.value.length} 个角色`)
    batchDeleteVisible.value = false
    selectedRoles.value = []
    loadRoles()
  } catch (e: any) {
    ElMessage.error(e.message || '批量删除失败')
  } finally {
    batchDeleting.value = false
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (isEdit.value) {
        await updateRole(currentRoleId.value, form)
        ElMessage.success('更新成功')
      } else {
        await createRole(form)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadRoles()
    } catch (e: any) {
      ElMessage.error(e.message || (isEdit.value ? '更新失败' : '创建失败'))
    } finally {
      submitting.value = false
    }
  })
}

const handleDownload = () => {
  ElMessage.info('导出功能待实现')
}

const handleSettings = () => {
  ElMessage.info('设置功能待实现')
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  loadRoles()
})
</script>

<style scoped>
.roles-page {
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

.toolbar {
  display: flex;
  gap: 8px;
  align-items: center;
}

.toolbar-icons {
  display: flex;
  gap: 4px;
}

.search-bar {
  display: flex;
  gap: 8px;
  padding: 12px 0;
  align-items: center;
}

.search-icon {
  cursor: pointer;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.policy-dialog-header {
  margin-bottom: 16px;
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.assigned-policies h4,
.available-policies h4 {
  margin-bottom: 10px;
  font-size: 14px;
  color: #606266;
}
</style>