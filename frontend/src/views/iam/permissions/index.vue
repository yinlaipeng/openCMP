<template>
  <div class="permissions-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">权限</span>
          <!-- 工具栏 -->
          <div class="toolbar">
            <el-button @click="handleRefresh" :loading="loading">
              <el-icon><Refresh /></el-icon>
            </el-button>
            <el-button type="primary" @click="handleCreate">
              <el-icon><Plus /></el-icon>
              新建
            </el-button>
            <el-button :disabled="selectedPolicies.length === 0" @click="handleBatchDisable">
              禁用
            </el-button>
            <el-button :disabled="selectedPolicies.length === 0" @click="handleBatchEnable">
              启用
            </el-button>
            <el-button :disabled="selectedPolicies.length === 0" @click="handleBatchDelete">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </div>
        </div>
      </template>

      <!-- 权限类型切换 Tabs -->
      <el-tabs v-model="activeTab" @tab-change="handleTabChange" class="permission-tabs">
        <el-tab-pane label="全部" name="all"></el-tab-pane>
        <el-tab-pane label="自定义权限" name="custom"></el-tab-pane>
        <el-tab-pane label="系统权限" name="system"></el-tab-pane>
      </el-tabs>

      <!-- 搜索栏 -->
      <div class="search-bar">
        <el-dropdown trigger="click" @command="handleFieldChange">
          <el-button>
            {{ currentFieldLabel }} <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="name">名称</el-dropdown-item>
              <el-dropdown-item command="scope">策略范围</el-dropdown-item>
              <el-dropdown-item command="enabled">启用状态</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-input
          v-model="searchKeyword"
          :placeholder="searchPlaceholder"
          clearable
          style="width: 300px"
          @keyup.enter="loadPolicies"
        >
          <template #suffix>
            <el-icon class="search-icon" @click="loadPolicies"><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="loadPolicies">查询</el-button>
        <el-button @click="handleResetSearch">重置</el-button>
      </div>

      <!-- 策略列表 -->
      <el-table
        :data="policies"
        v-loading="loading"
        border
        stripe
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="名称" min-width="200" show-overflow-tooltip sortable>
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)" class="name-link">
              <span>{{ row.name }}</span>
            </el-button>
            <br>
            <small style="color: #999;">{{ row.description || '-' }}</small>
          </template>
        </el-table-column>
        <el-table-column label="启用状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="策略范围" width="120">
          <template #default="{ row }">
            <el-tag :type="getScopeTagType(row.scope)" size="small">
              {{ getScopeLabel(row.scope) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="所属域" width="100">
          <template #default="{ row }">
            {{ row.domain_id || 'Default' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="handleEdit(row)" :disabled="row.is_system">
              编辑
            </el-button>
            <el-dropdown trigger="click" @command="(cmd: string) => handleMoreCommand(cmd, row)">
              <el-button size="small" type="primary" link>
                更多 <el-icon><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="view">详情</el-dropdown-item>
                  <el-dropdown-item command="enable" v-if="!row.enabled" :disabled="row.is_system">启用</el-dropdown-item>
                  <el-dropdown-item command="disable" v-if="row.enabled" :disabled="row.is_system">禁用</el-dropdown-item>
                  <el-dropdown-item command="roles">关联角色</el-dropdown-item>
                  <el-dropdown-item command="delete" :disabled="row.is_system || !row.can_delete" divided>删除</el-dropdown-item>
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
        @size-change="loadPolicies"
        @current-change="loadPolicies"
        class="pagination"
      />
    </el-card>

    <!-- 创建/编辑策略对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑策略' : '新建策略'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入策略名称" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="2" placeholder="请输入策略描述" />
        </el-form-item>
        <el-form-item label="策略范围" prop="scope">
          <el-select v-model="form.scope" placeholder="选择策略范围" :disabled="isEdit">
            <el-option label="管理后台" value="system" />
            <el-option label="无管理后台" value="domain" />
            <el-option label="项目视图" value="project" />
          </el-select>
        </el-form-item>
        <el-form-item label="策略内容" prop="policyStr">
          <el-input
            v-model="form.policyStr"
            type="textarea"
            :rows="6"
            placeholder="请输入JSON格式的策略内容"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 策略详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="策略详情" width="700px">
      <el-descriptions :column="2" border v-if="currentPolicy">
        <el-descriptions-item label="ID">{{ currentPolicy.id }}</el-descriptions-item>
        <el-descriptions-item label="策略名">{{ currentPolicy.name }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentPolicy.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="策略范围">
          <el-tag :type="getScopeTagType(currentPolicy.scope)" size="small">
            {{ getScopeLabel(currentPolicy.scope) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="启用状态">
          <el-tag :type="currentPolicy.enabled ? 'success' : 'info'" size="small">
            {{ currentPolicy.enabled ? '启用' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="系统策略">
          <el-tag :type="currentPolicy.is_system ? 'warning' : 'success'" size="small">
            {{ currentPolicy.is_system ? '是' : '否' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="所属域">{{ currentPolicy.domain_id || 'Default' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentPolicy.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentPolicy.updated_at) }}</el-descriptions-item>
      </el-descriptions>

      <div style="margin-top: 20px">
        <h4>策略内容</h4>
        <pre class="policy-content">{{ formatPolicy(currentPolicy.policy) }}</pre>
      </div>

      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 关联角色对话框 -->
    <el-dialog v-model="rolesDialogVisible" title="关联角色" width="600px">
      <el-table :data="policyRoles" v-loading="rolesLoading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名" min-width="150" />
        <el-table-column prop="display_name" label="显示名" min-width="150" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
      </el-table>
      <template #footer>
        <el-button @click="rolesDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 批量删除确认弹窗 -->
    <el-dialog v-model="batchDeleteVisible" title="批量删除确认" width="500px">
      <p>确定要删除以下 {{ selectedPolicies.length }} 个策略吗？</p>
      <el-table :data="selectedPolicies" border stripe max-height="300">
        <el-table-column prop="id" label="ID" width="200" show-overflow-tooltip />
        <el-table-column prop="name" label="策略名" />
        <el-table-column label="系统策略" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_system ? 'warning' : 'success'" size="small">
              {{ row.is_system ? '是' : '否' }}
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
import { Refresh, Plus, Delete, ArrowDown, Search } from '@element-plus/icons-vue'
import {
  getPolicies,
  createPolicy,
  updatePolicy,
  deletePolicy,
  enablePolicy,
  disablePolicy,
  getPolicyRoles,
  batchEnablePolicies,
  batchDisablePolicies,
  batchDeletePolicies
} from '@/api/iam'

interface Policy {
  id: string
  name: string
  description?: string
  scope: string
  domain_id?: string
  policy: any
  is_system: boolean
  enabled: boolean
  can_delete?: boolean
  can_update?: boolean
  created_at: string
  updated_at: string
}

const policies = ref<Policy[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const rolesDialogVisible = ref(false)
const batchDeleteVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const batchDeleting = ref(false)
const currentPolicyId = ref('')
const currentPolicy = ref<Policy | null>(null)
const policyRoles = ref<any[]>([])
const selectedPolicies = ref<Policy[]>([])
const rolesLoading = ref(false)
const formRef = ref<FormInstance>()

// Tabs 筛选
const activeTab = ref('all')

// 搜索相关
const searchField = ref('name')
const searchKeyword = ref('')

const currentFieldLabel = computed(() => {
  const labels: Record<string, string> = {
    name: '名称',
    scope: '策略范围',
    enabled: '启用状态'
  }
  return labels[searchField.value] || '名称'
})

const searchPlaceholder = computed(() => {
  const placeholders: Record<string, string> = {
    name: '请输入策略名称...',
    scope: '请选择范围：system/domain/project...',
    enabled: '请输入状态：true/false...'
  }
  return placeholders[searchField.value] || '请输入搜索关键词...'
})

const pagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

const form = reactive({
  name: '',
  description: '',
  scope: 'system',
  policyStr: '{"policy": {"statement": [{"effect": "allow", "action": ["*"], "resource": ["*"]}]}}'
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入策略名称', trigger: 'blur' }],
  scope: [{ required: true, message: '请选择策略范围', trigger: 'change' }],
  policyStr: [{ required: true, message: '请输入策略内容', trigger: 'blur' }]
}

// 处理 Tab 切换
const handleTabChange = (tabName: string) => {
  pagination.page = 1
  loadPolicies()
}

const handleFieldChange = (field: string) => {
  searchField.value = field
  searchKeyword.value = ''
}

// 策略范围显示标签（参考 Cloudpods 定义）
const getScopeLabel = (scope: string) => {
  const labels: Record<string, string> = {
    system: '管理后台',
    domain: '无管理后台',
    project: '项目视图'
  }
  return labels[scope] || scope
}

const getScopeTagType = (scope: string) => {
  const types: Record<string, string> = {
    system: 'danger',
    domain: 'warning',
    project: 'success'
  }
  return types[scope] || 'info'
}

const loadPolicies = async () => {
  loading.value = true
  try {
    const params: any = {
      show_fail_reason: true,
      details: true,
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit
    }

    // 根据 Tab 设置筛选参数（后端期望字符串 "true"/"false")
    if (activeTab.value === 'system') {
      params.is_system = 'true'
    } else if (activeTab.value === 'custom') {
      params.is_system = 'false'
    }
    // 'all' 不添加筛选参数

    if (searchKeyword.value) {
      if (searchField.value === 'name') {
        params.keyword = searchKeyword.value
      } else if (searchField.value === 'scope') {
        params.scope = searchKeyword.value
      } else if (searchField.value === 'enabled') {
        params.enabled = searchKeyword.value === 'true'
      }
    }

    const res = await getPolicies(params)
    policies.value = res.data || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载策略列表失败')
  } finally {
    loading.value = false
  }
}

const loadPolicyRoles = async (policyId: string) => {
  rolesLoading.value = true
  try {
    const res = await getPolicyRoles(policyId)
    policyRoles.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    rolesLoading.value = false
  }
}

const handleRefresh = () => {
  loadPolicies()
}

const handleResetSearch = () => {
  searchKeyword.value = ''
  searchField.value = 'name'
  pagination.page = 1
  loadPolicies()
}

const handleSelectionChange = (selection: Policy[]) => {
  selectedPolicies.value = selection
}

const handleCreate = () => {
  isEdit.value = false
  form.name = ''
  form.description = ''
  form.scope = 'system'
  form.policyStr = '{"policy": {"statement": [{"effect": "allow", "action": ["*"], "resource": ["*"]}]}}'
  dialogVisible.value = true
}

const handleEdit = (row: Policy) => {
  isEdit.value = true
  currentPolicyId.value = row.id
  form.name = row.name
  form.description = row.description || ''
  form.scope = row.scope
  form.policyStr = JSON.stringify(row.policy, null, 2)
  dialogVisible.value = true
}

const handleView = (row: Policy) => {
  currentPolicy.value = row
  detailDialogVisible.value = true
}

const handleMoreCommand = (command: string, row: Policy) => {
  switch (command) {
    case 'view':
      handleView(row)
      break
    case 'enable':
      handleToggleEnable(row)
      break
    case 'disable':
      handleToggleEnable(row)
      break
    case 'roles':
      handleViewRoles(row)
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

const handleViewRoles = async (row: Policy) => {
  currentPolicyId.value = row.id
  await loadPolicyRoles(row.id)
  rolesDialogVisible.value = true
}

const handleToggleEnable = async (row: Policy) => {
  try {
    const action = row.enabled ? '禁用' : '启用'
    await ElMessageBox.confirm(`确定要${action}该策略吗？`, '提示', { type: 'warning' })

    if (row.enabled) {
      await disablePolicy(row.id)
    } else {
      await enablePolicy(row.id)
    }

    ElMessage.success(`${action}成功`)
    loadPolicies()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || `${row.enabled ? '禁用' : '启用'}失败`)
    }
  }
}

const handleDelete = async (row: Policy) => {
  try {
    await ElMessageBox.confirm('确定要删除该策略吗？', '提示', { type: 'warning' })
    await deletePolicy(row.id)
    ElMessage.success('删除成功')
    loadPolicies()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const handleBatchEnable = async () => {
  if (selectedPolicies.value.length === 0) {
    ElMessage.warning('请先选择要启用的策略')
    return
  }
  try {
    await ElMessageBox.confirm(`确定要启用选中的 ${selectedPolicies.value.length} 个策略吗？`, '提示', { type: 'warning' })
    const ids = selectedPolicies.value.map(p => p.id)
    await batchEnablePolicies(ids)
    ElMessage.success('批量启用成功')
    selectedPolicies.value = []
    loadPolicies()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '批量启用失败')
    }
  }
}

const handleBatchDisable = async () => {
  if (selectedPolicies.value.length === 0) {
    ElMessage.warning('请先选择要禁用的策略')
    return
  }
  try {
    await ElMessageBox.confirm(`确定要禁用选中的 ${selectedPolicies.value.length} 个策略吗？`, '提示', { type: 'warning' })
    const ids = selectedPolicies.value.map(p => p.id)
    await batchDisablePolicies(ids)
    ElMessage.success('批量禁用成功')
    selectedPolicies.value = []
    loadPolicies()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '批量禁用失败')
    }
  }
}

const handleBatchDelete = () => {
  if (selectedPolicies.value.length === 0) {
    ElMessage.warning('请先选择要删除的策略')
    return
  }
  const systemPolicies = selectedPolicies.value.filter(p => p.is_system)
  if (systemPolicies.length > 0) {
    ElMessage.warning('不能删除系统策略，请重新选择')
    return
  }
  batchDeleteVisible.value = true
}

const handleBatchDeleteConfirm = async () => {
  batchDeleting.value = true
  try {
    const ids = selectedPolicies.value.map(p => p.id)
    await batchDeletePolicies(ids)
    ElMessage.success(`成功删除 ${selectedPolicies.value.length} 个策略`)
    batchDeleteVisible.value = false
    selectedPolicies.value = []
    loadPolicies()
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
      const policyData = {
        name: form.name,
        description: form.description,
        scope: form.scope,
        policy: JSON.parse(form.policyStr)
      }

      if (isEdit.value) {
        await updatePolicy(currentPolicyId.value, policyData)
        ElMessage.success('更新成功')
      } else {
        await createPolicy(policyData)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadPolicies()
    } catch (e: any) {
      ElMessage.error(e.message || (isEdit.value ? '更新失败' : '创建失败'))
    } finally {
      submitting.value = false
    }
  })
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const formatPolicy = (policy: any) => {
  if (!policy) return '-'
  try {
    return JSON.stringify(policy, null, 2)
  } catch {
    return policy
  }
}

onMounted(() => {
  loadPolicies()
})
</script>

<style scoped>
.permissions-page {
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

.permission-tabs {
  margin-bottom: 16px;
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

.policy-content {
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
  overflow-x: auto;
  font-size: 12px;
  max-height: 300px;
}
</style>