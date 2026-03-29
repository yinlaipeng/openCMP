<template>
  <div class="roles-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">角色</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建角色
          </el-button>
        </div>
      </template>

      <!-- 筛选条件 -->
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="角色名">
          <el-input v-model="filterForm.name" placeholder="请输入角色名" clearable @keyup.enter="loadRoles" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="filterForm.type" placeholder="全部" clearable style="width: 120px">
            <el-option label="系统" value="system" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterForm.enabled" placeholder="全部" clearable style="width: 100px">
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadRoles">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 角色列表 -->
      <el-table :data="roles" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名" min-width="150" show-overflow-tooltip />
        <el-table-column prop="display_name" label="显示名" min-width="150" show-overflow-tooltip />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'system' ? 'warning' : 'success'">
              {{ row.type === 'system' ? '系统' : '自定义' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_public" label="公开" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_public ? 'primary' : 'info'" size="small">
              {{ row.is_public ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="320" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">详情</el-button>
            <el-button size="small" @click="handleEdit(row)" :disabled="row.type === 'system'">编辑</el-button>
            <el-button size="small" @click="handlePermissions(row)">权限</el-button>
            <el-button 
              size="small" 
              :type="row.enabled ? 'warning' : 'success'" 
              @click="handleToggleEnable(row)"
              :disabled="row.type === 'system'"
            >
              {{ row.enabled ? '禁用' : '启用' }}
            </el-button>
            <el-button 
              size="small" 
              type="info" 
              @click="handleMakePublic(row)"
              :disabled="row.type === 'system' || row.is_public"
            >
              公开
            </el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)" :disabled="row.type === 'system'">删除</el-button>
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
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑角色' : '新建角色'" width="700px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="角色名" prop="name">
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
        
        <!-- 策略选择 -->
        <el-form-item label="关联策略">
          <el-select
            v-model="selectedPolicyIds"
            multiple
            filterable
            placeholder="请选择策略"
            style="width: 100%"
          >
            <el-option
              v-for="item in policies"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            >
              <span style="float: left">{{ item.name }}</span>
              <span style="float: right; color: #8492a6; font-size: 13px">
                <el-tag size="small" :type="item.is_system ? 'warning' : 'success'">
                  {{ item.is_system ? '系统' : '自定义' }}
                </el-tag>
              </span>
            </el-option>
          </el-select>
          <div class="form-item-tip">选择策略后，角色将拥有策略中定义的权限</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 权限管理对话框 -->
    <el-dialog v-model="permDialogVisible" title="角色权限管理" width="900px">
      <div class="perm-dialog-header">
        <span>当前角色：<strong>{{ currentRole?.name }}</strong>（{{ currentRole?.display_name }}）</span>
      </div>
      
      <el-tabs v-model="permTab">
        <!-- 权限页签 -->
        <el-tab-pane label="权限" name="permission">
          <el-transfer
            v-model="selectedPermissions"
            :data="allPermissions"
            :titles="['可选权限', '已选权限']"
            filterable
            :filter-method="filterPermission"
            :props="{ key: 'id', label: 'display_name' }"
          >
            <template #default="{ option }">
              <span>{{ option.display_name }}</span>
              <el-tag size="small" style="margin-left: 8px">{{ option.resource }}</el-tag>
              <el-tag size="small" type="success" style="margin-left: 4px">{{ option.action }}</el-tag>
            </template>
          </el-transfer>
        </el-tab-pane>
        
        <!-- 策略页签 -->
        <el-tab-pane label="策略" name="policy">
          <div class="policy-transfer-container">
            <el-transfer
              v-model="selectedPolicyIds"
              :data="allPolicies"
              :titles="['可选策略', '已选策略']"
              filterable
              :filter-method="filterPolicy"
              :props="{ key: 'id', label: 'name' }"
            >
              <template #default="{ option }">
                <span>{{ option.name }}</span>
                <el-tag size="small" style="margin-left: 8px" :type="option.is_system ? 'warning' : 'success'">
                  {{ option.is_system ? '系统' : '自定义' }}
                </el-tag>
              </template>
            </el-transfer>
          </div>
        </el-tab-pane>
      </el-tabs>
      
      <template #footer>
        <el-button @click="permDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSavePermissions" :loading="permSubmitting">保存</el-button>
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
            <el-table-column prop="id" label="ID" width="200" />
            <el-table-column prop="name" label="策略名" />
            <el-table-column prop="is_system" label="类型" width="100">
              <template #default="{ row }">
                <el-tag size="small" :type="row.is_system ? 'warning' : 'success'">
                  {{ row.is_system ? '系统' : '自定义' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>

      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import {
  getRoles,
  createRole,
  updateRole,
  deleteRole,
  getPermissions,
  getRolePermissions,
  assignPermission,
  revokePermission,
  enableRole,
  disableRole,
  makeRolePublic,
  getRoleUsers,
  getRoleGroups,
  getPolicies,
  getRolePolicies,
  assignPolicyToRole,
  revokePolicyFromRole
} from '@/api/iam'

const roles = ref<any[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const permDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const detailTab = ref('users')
const permTab = ref('permission')
const isEdit = ref(false)
const submitting = ref(false)
const permSubmitting = ref(false)
const currentRoleId = ref(0)
const currentRole = ref<any>(null)
const selectedPermissions = ref<number[]>([])
const allPermissions = ref<any[]>([])
const selectedPolicyIds = ref<string[]>([])
const allPolicies = ref<any[]>([])
const roleUsers = ref<any[]>([])
const roleGroups = ref<any[]>([])
const rolePolicies = ref<any[]>([])
const usersLoading = ref(false)
const groupsLoading = ref(false)
const policiesLoading = ref(false)
const formRef = ref<FormInstance>()

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
  name: [{ required: true, message: '请输入角色名', trigger: 'blur' }],
  display_name: [{ required: true, message: '请输入显示名', trigger: 'blur' }]
}

const loadRoles = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit
    }
    if (filterForm.name) params.keyword = filterForm.name
    if (filterForm.type) params.type = filterForm.type
    if (filterForm.enabled !== undefined) params.enabled = filterForm.enabled

    const res = await getRoles(params)
    roles.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载角色列表失败')
  } finally {
    loading.value = false
  }
}

const loadPermissions = async () => {
  try {
    const res = await getPermissions({ limit: 1000 })
    allPermissions.value = (res.items || []).map((p: any) => ({
      id: p.id,
      display_name: `${p.display_name || p.name} (${p.resource}.${p.action})`,
      resource: p.resource,
      action: p.action,
      data: p
    }))
  } catch (e: any) {
    ElMessage.error(e.message || '加载权限列表失败')
  }
}

const loadPolicies = async () => {
  try {
    const res = await getPolicies({ limit: 1000 })
    allPolicies.value = (res.data || []).map((p: any) => ({
      id: p.id,
      name: `${p.name}`,
      is_system: p.is_system,
      data: p
    }))
  } catch (e: any) {
    ElMessage.error(e.message || '加载策略列表失败')
  }
}

const loadRolePermissions = async (roleId: number) => {
  try {
    const res = await getRolePermissions(roleId)
    selectedPermissions.value = (res.items || []).map((p: any) => p.id)
  } catch (e: any) {
    console.error(e)
  }
}

const loadRolePolicies = async (roleId: number) => {
  try {
    const res = await getRolePolicies(roleId)
    selectedPolicyIds.value = (res.items || []).map((p: any) => p.id)
  } catch (e: any) {
    console.error(e)
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

const loadRolePoliciesDetail = async (roleId: number) => {
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

const resetFilter = () => {
  filterForm.name = ''
  filterForm.type = ''
  filterForm.enabled = undefined
  pagination.page = 1
  loadRoles()
}

const handleCreate = async () => {
  isEdit.value = false
  form.name = ''
  form.display_name = ''
  form.description = ''
  form.type = 'custom'
  selectedPolicyIds.value = []
  await loadPolicies()
  dialogVisible.value = true
}

const handleEdit = async (row: any) => {
  isEdit.value = true
  currentRoleId.value = row.id
  form.name = row.name
  form.display_name = row.display_name
  form.description = row.description
  form.type = row.type
  selectedPolicyIds.value = []
  await loadPolicies()
  await loadRolePolicies(row.id)
  dialogVisible.value = true
}

const handleView = async (row: any) => {
  currentRole.value = row
  detailTab.value = 'users'
  await loadRoleUsers(row.id)
  await loadRoleGroups(row.id)
  await loadRolePoliciesDetail(row.id)
  detailDialogVisible.value = true
}

const handlePermissions = async (row: any) => {
  currentRoleId.value = row.id
  currentRole.value = row
  permTab.value = 'permission'
  await loadPermissions()
  await loadPolicies()
  await loadRolePermissions(row.id)
  await loadRolePolicies(row.id)
  permDialogVisible.value = true
}

const filterPermission = (query: string, item: any) => {
  return item.display_name.toLowerCase().includes(query.toLowerCase())
}

const filterPolicy = (query: string, item: any) => {
  return item.name.toLowerCase().includes(query.toLowerCase())
}

const handleSavePermissions = async () => {
  permSubmitting.value = true
  try {
    if (permTab.value === 'permission') {
      // 保存权限
      const currentPermIds = selectedPermissions.value
      const existingPermIds = (await getRolePermissions(currentRoleId.value)).items?.map((p: any) => p.id) || []

      const toAdd = currentPermIds.filter(id => !existingPermIds?.includes(id))
      const toRemove = existingPermIds?.filter(id => !currentPermIds.includes(id)) || []

      for (const permId of toAdd) {
        await assignPermission(currentRoleId.value, permId)
      }
      for (const permId of toRemove) {
        await revokePermission(currentRoleId.value, permId)
      }
    } else {
      // 保存策略
      const currentPolicyIds = selectedPolicyIds.value
      const existingPolicyIds = (await getRolePolicies(currentRoleId.value)).items?.map((p: any) => p.id) || []

      const toAdd = currentPolicyIds.filter(id => !existingPolicyIds?.includes(id))
      const toRemove = existingPolicyIds?.filter(id => !currentPolicyIds.includes(id)) || []

      for (const policyId of toAdd) {
        await assignPolicyToRole(currentRoleId.value, policyId)
      }
      for (const policyId of toRemove) {
        await revokePolicyFromRole(currentRoleId.value, policyId)
      }
    }

    ElMessage.success('保存成功')
    permDialogVisible.value = false
    loadRoles()
  } catch (e: any) {
    ElMessage.error(e.message || '保存失败')
  } finally {
    permSubmitting.value = false
  }
}

const handleToggleEnable = async (row: any) => {
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
      ElMessage.error(e.message || `${action}失败`)
    }
  }
}

const handleMakePublic = async (row: any) => {
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

const handleDelete = async (row: any) => {
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

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (isEdit.value) {
        await updateRole(currentRoleId.value, form)
        // 更新策略关联
        await updateRolePolicies(currentRoleId.value, selectedPolicyIds.value)
        ElMessage.success('更新成功')
      } else {
        const result = await createRole(form)
        // 创建后关联策略
        if (selectedPolicyIds.value.length > 0 && result.id) {
          for (const policyId of selectedPolicyIds.value) {
            await assignPolicyToRole(result.id, policyId)
          }
        }
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

const updateRolePolicies = async (roleId: number, newPolicyIds: string[]) => {
  const existingPolicies = (await getRolePolicies(roleId)).items || []
  const existingPolicyIds = existingPolicies.map((p: any) => p.id)
  
  const toAdd = newPolicyIds.filter(id => !existingPolicyIds.includes(id))
  const toRemove = existingPolicyIds.filter(id => !newPolicyIds.includes(id))
  
  for (const policyId of toAdd) {
    await assignPolicyToRole(roleId, policyId)
  }
  for (const policyId of toRemove) {
    await revokePolicyFromRole(roleId, policyId)
  }
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

.filter-form {
  margin-bottom: 16px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.perm-dialog-header {
  margin-bottom: 16px;
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.form-item-tip {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.policy-transfer-container {
  display: flex;
  justify-content: center;
}
</style>
