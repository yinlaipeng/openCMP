<template>
  <div class="permissions-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">权限</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建权限
          </el-button>
        </div>
      </template>

      <!-- 权限类型切换 -->
      <el-tabs v-model="activeTab" @tab-change="handleTabChange" class="permission-tabs">
        <el-tab-pane label="全部权限" name="all"></el-tab-pane>
        <el-tab-pane label="自定义权限" name="custom"></el-tab-pane>
        <el-tab-pane label="系统权限" name="system"></el-tab-pane>
      </el-tabs>

      <!-- 筛选条件 -->
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="权限名称">
          <el-input
            v-model="filterForm.name"
            placeholder="请输入权限名称"
            clearable
            @keyup.enter="loadPermissions"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterForm.enabled" placeholder="全部" clearable style="width: 100px">
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item label="权限范围">
          <el-select v-model="filterForm.scope" placeholder="全部" clearable style="width: 120px">
            <el-option label="全局" value="global" />
            <el-option label="域" value="domain" />
            <el-option label="项目" value="project" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadPermissions">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 权限列表 -->
      <el-table :data="permissions" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="权限名称" min-width="150" show-overflow-tooltip>
          <template #default="{ row }">
            <el-link type="primary" @click="showPermissionDetails(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="display_name" label="显示名称" min-width="150" show-overflow-tooltip />
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
        <el-table-column prop="scope" label="权限范围" width="120">
          <template #default="{ row }">
            <el-tag
              :type="row.scope === 'global' ? 'primary' :
                     row.scope === 'domain' ? 'warning' :
                     'danger'"
            >
              {{ row.scope === 'global' ? '全局' : row.scope === 'domain' ? '域' : '项目' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="domain_id" label="域" width="100">
          <template #default="{ row }">
            {{ row.domain_id || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="resource" label="资源" width="120" show-overflow-tooltip />
        <el-table-column prop="action" label="操作" width="120" show-overflow-tooltip />
        <el-table-column prop="is_public" label="公开" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_public ? 'primary' : 'info'" size="small">
              {{ row.is_public ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="showPermissionDetails(row)">详情</el-button>
            <el-button size="small" @click="handleEdit(row)" :disabled="row.type === 'system'">编辑</el-button>
            <el-button
              size="small"
              :type="row.enabled ? 'warning' : 'success'"
              @click="handleToggleEnable(row)"
              :disabled="row.type === 'system'"
            >
              {{ row.enabled ? '禁用' : '启用' }}
            </el-button>
            <el-dropdown trigger="click" placement="bottom" size="small">
              <el-button size="small" type="info">
                更多
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleClone(row)">
                    <el-icon><CopyDocument /></el-icon>
                    克隆
                  </el-dropdown-item>
                  <el-dropdown-item
                    @click="handleMakePublic(row)"
                    :disabled="row.type === 'system' || row.is_public"
                  >
                    <el-icon><Promotion /></el-icon>
                    设为公开
                  </el-dropdown-item>
                  <el-dropdown-item
                    @click="handleDelete(row)"
                    :disabled="row.type === 'system'"
                    divided
                  >
                    <el-icon><Delete /></el-icon>
                    删除
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
        @size-change="loadPermissions"
        @current-change="loadPermissions"
        class="pagination"
      />
    </el-card>

    <!-- 添加/编辑权限对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑权限' : '新建权限'" width="800px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input
            v-model="form.name"
            placeholder="请输入权限名称（英文标识符）"
            :disabled="isEdit"
          />
        </el-form-item>
        <el-form-item label="备注" prop="display_name">
          <el-input v-model="form.display_name" placeholder="请输入显示名称（中文）" />
        </el-form-item>
        <el-form-item label="权限范围" prop="scope">
          <el-radio-group v-model="form.scope" @change="onScopeChange">
            <el-radio label="global">管理后台</el-radio>
            <el-radio label="domain">无管理后台</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="编辑模式">
          <el-radio-group v-model="editorMode">
            <el-radio label="form">表单编辑</el-radio>
            <el-radio label="yaml">YAML编辑</el-radio>
          </el-radio-group>
        </el-form-item>
        <template v-if="editorMode === 'yaml'">
          <el-form-item label="权限内容" class="yaml-editor-container">
            <el-input
              v-model="yamlContent"
              type="textarea"
              :rows="12"
              placeholder="请输入YAML格式的权限内容，例如：
policy:
  '*': allow"
              class="yaml-textarea"
            />
          </el-form-item>
        </template>
        <template v-else>
          <el-form-item label="资源类型" prop="resource">
            <el-input v-model="form.resource" placeholder="请输入资源类型，如 vm, vpc, subnet 等" />
          </el-form-item>
          <el-form-item label="操作类型" prop="action">
            <el-input v-model="form.action" placeholder="请输入操作类型，如 read, write, delete 等" />
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入权限描述" />
          </el-form-item>
          <el-form-item label="域" prop="domain_id" v-if="form.scope === 'domain'">
            <el-select v-model="form.domain_id" placeholder="请选择域" style="width: 100%">
              <el-option
                v-for="domain in domains"
                :key="domain.id"
                :label="domain.name"
                :value="domain.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="项目" prop="project_id" v-if="form.scope === 'project'">
            <el-select v-model="form.project_id" placeholder="请选择项目" style="width: 100%">
              <el-option
                v-for="project in projects"
                :key="project.id"
                :label="project.name"
                :value="project.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="类型" prop="type" v-if="!isEdit">
            <el-radio-group v-model="form.type">
              <el-radio label="custom">自定义</el-radio>
              <el-radio label="system">系统</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="启用状态">
            <el-switch v-model="form.enabled" />
          </el-form-item>
          <el-form-item label="设为公开">
            <el-switch v-model="form.is_public" />
            <div class="form-item-tip">
              公开权限可以被其他域使用
            </div>
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 权限详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="权限详情" width="800px">
      <el-tabs v-model="detailActiveTab">
        <el-tab-pane label="详情" name="details">
          <el-descriptions :column="2" border v-if="currentPermission">
            <el-descriptions-item label="ID">{{ currentPermission.id }}</el-descriptions-item>
            <el-descriptions-item label="权限名称">{{ currentPermission.name }}</el-descriptions-item>
            <el-descriptions-item label="显示名称">{{ currentPermission.display_name }}</el-descriptions-item>
            <el-descriptions-item label="资源类型">{{ currentPermission.resource }}</el-descriptions-item>
            <el-descriptions-item label="操作类型">{{ currentPermission.action }}</el-descriptions-item>
            <el-descriptions-item label="类型">
              <el-tag :type="currentPermission.type === 'system' ? 'warning' : 'success'" size="small">
                {{ currentPermission.type === 'system' ? '系统' : '自定义' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="权限范围">
              <el-tag
                :type="currentPermission.scope === 'global' ? 'primary' :
                       currentPermission.scope === 'domain' ? 'warning' :
                       'danger'"
                size="small"
              >
                {{ currentPermission.scope === 'global' ? '全局' : currentPermission.scope === 'domain' ? '域' : '项目' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="域">
              {{ currentPermission.domain_id || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="currentPermission.enabled ? 'success' : 'info'" size="small">
                {{ currentPermission.enabled ? '启用' : '禁用' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="公开">
              <el-tag :type="currentPermission.is_public ? 'primary' : 'info'" size="small">
                {{ currentPermission.is_public ? '是' : '否' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="描述" :span="2">{{ currentPermission.description || '-' }}</el-descriptions-item>
            <el-descriptions-item label="条件" :span="2">
              {{ currentPermission.conditions ? JSON.stringify(currentPermission.conditions) : '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatDate(currentPermission.created_at) }}</el-descriptions-item>
            <el-descriptions-item label="更新时间">{{ formatDate(currentPermission.updated_at) }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="操作日志" name="logs">
          <el-empty description="暂无操作日志" />
          <!-- 在实际实现中，这里应该显示权限的操作日志 -->
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
import { ArrowDown, Plus, CopyDocument, Promotion, Delete } from '@element-plus/icons-vue'
import { Permission } from '@/types/permission'
import {
  getPermissions,
  createPermission,
  updatePermission,
  deletePermission,
  enablePermission,
  disablePermission,
  clonePermission,
  makePermissionPublic,
  getPermission
} from '@/api/iam'
import { getDomains } from '@/api/iam'
import { getProjects } from '@/api/iam'

const permissions = ref<Permission[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const currentPermissionId = ref(0)
const currentPermission = ref<Permission | null>(null)
const formRef = ref<FormInstance>()
const domains = ref<{id: number, name: string}[]>([])
const projects = ref<{id: number, name: string}[]>([])

// 新增的响应式变量
const activeTab = ref('custom') // 默认选中自定义权限
const detailActiveTab = ref('details') // 权限详情标签页
const editorMode = ref<'form' | 'yaml'>('form') // 编辑模式
const yamlContent = ref('') // YAML内容

const filterForm = reactive({
  name: '',
  type: '',
  enabled: undefined as boolean | undefined,
  scope: ''
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
  resource: '',
  action: '',
  scope: 'global' as 'global' | 'domain' | 'project',
  domain_id: undefined as number | undefined,
  project_id: undefined as number | undefined,
  type: 'custom' as 'custom' | 'system',
  enabled: true,
  is_public: false
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入权限名称', trigger: 'blur' }],
  display_name: [{ required: true, message: '请输入显示名称', trigger: 'blur' }],
  resource: [{ required: true, message: '请输入资源类型', trigger: 'blur' }],
  action: [{ required: true, message: '请输入操作类型', trigger: 'blur' }],
  scope: [{ required: true, message: '请选择权限范围', trigger: 'change' }]
}

// 处理标签页切换
const handleTabChange = (tabName: string) => {
  // 重置筛选条件，设置对应的类型过滤
  filterForm.type = tabName === 'all' ? '' : tabName
  pagination.page = 1
  loadPermissions()
}

// 显示权限详情
const showPermissionDetails = async (row: Permission) => {
  currentPermission.value = row
  detailDialogVisible.value = true
  detailActiveTab.value = 'details'
}

const loadPermissions = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit
    }
    if (filterForm.name) params.keyword = filterForm.name
    // 如果没有明确选择类型，则不添加类型筛选
    if (filterForm.type) params.type = filterForm.type
    if (filterForm.enabled !== undefined) params.enabled = filterForm.enabled
    if (filterForm.scope) params.scope = filterForm.scope

    const res = await getPermissions(params)
    permissions.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载权限列表失败')
  } finally {
    loading.value = false
  }
}

const loadDomains = async () => {
  try {
    const res = await getDomains({ limit: 1000 })
    domains.value = (res.items || []).map(d => ({ id: d.id, name: d.name }))
  } catch (e: any) {
    console.error('加载域列表失败:', e)
  }
}

const loadProjects = async () => {
  try {
    const res = await getProjects({ limit: 1000 })
    projects.value = (res.items || []).map(p => ({ id: p.id, name: p.name }))
  } catch (e: any) {
    console.error('加载项目列表失败:', e)
  }
}

const resetFilter = () => {
  filterForm.name = ''
  filterForm.type = activeTab.value === 'all' ? '' : activeTab.value
  filterForm.enabled = undefined
  filterForm.scope = ''
  pagination.page = 1
  loadPermissions()
}

const handleCreate = async () => {
  isEdit.value = false
  form.name = ''
  form.display_name = ''
  form.description = ''
  form.resource = ''
  form.action = ''
  form.scope = 'global'
  form.domain_id = undefined
  form.project_id = undefined
  form.type = 'custom'
  form.enabled = true
  form.is_public = false

  // 重置 YAML 内容
  yamlContent.value = ''
  editorMode.value = 'form'

  dialogVisible.value = true
}

const handleEdit = async (row: Permission) => {
  isEdit.value = true
  currentPermissionId.value = row.id
  form.name = row.name
  form.display_name = row.display_name
  form.description = row.description
  form.resource = row.resource
  form.action = row.action
  form.scope = row.scope
  form.domain_id = row.domain_id
  form.project_id = row.project_id
  form.type = row.type
  form.enabled = row.enabled
  form.is_public = row.is_public

  // 设置为表单模式进行编辑
  editorMode.value = 'form'
  dialogVisible.value = true
}

const handleView = async (row: Permission) => {
  currentPermission.value = row
  detailDialogVisible.value = true
}

const handleToggleEnable = async (row: Permission) => {
  try {
    const action = row.enabled ? '禁用' : '启用'
    await ElMessageBox.confirm(`确定要${action}该权限吗？`, '提示', { type: 'warning' })

    if (row.enabled) {
      await disablePermission(row.id)
    } else {
      await enablePermission(row.id)
    }

    ElMessage.success(`${action}成功`)
    loadPermissions()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || `${action}失败`)
    }
  }
}

const handleClone = async (row: Permission) => {
  try {
    await ElMessageBox.confirm('确定要克隆该权限吗？', '提示', { type: 'warning' })

    // 生成克隆后的默认名称
    const clonedName = `${row.name}-copy`
    const clonedDisplayName = `${row.display_name || row.name} (副本)`

    await clonePermission(row.id, {
      name: clonedName,
      display_name: clonedDisplayName
    })

    ElMessage.success('克隆成功')
    loadPermissions()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '克隆失败')
    }
  }
}

const handleMakePublic = async (row: Permission) => {
  try {
    await ElMessageBox.confirm('确定要将该权限设为公开吗？公开后其他域的用户也可以使用此权限。', '提示', { type: 'warning' })
    await makePermissionPublic(row.id)
    ElMessage.success('设为公开成功')
    loadPermissions()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '设为公开失败')
    }
  }
}

const handleDelete = async (row: Permission) => {
  try {
    await ElMessageBox.confirm('确定要删除该权限吗？', '提示', { type: 'warning' })
    await deletePermission(row.id)
    ElMessage.success('删除成功')
    loadPermissions()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  if (editorMode.value === 'yaml') {
    // 如果是YAML模式，需要解析YAML并填充到表单
    // 这里简单处理，实际上可能需要更复杂的YAML解析
    try {
      // 模拟解析YAML，提取必要字段
      const parsedYaml = yamlContent.value
      if (!parsedYaml.trim()) {
        ElMessage.error('YAML内容不能为空')
        return
      }
      // 在实际应用中，这里需要真正的YAML解析逻辑
    } catch (e) {
      ElMessage.error('YAML格式错误，请检查语法')
      return
    }
  }

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (isEdit.value) {
        await updatePermission(currentPermissionId.value, form)
        ElMessage.success('更新成功')
      } else {
        await createPermission(form)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadPermissions()
    } catch (e: any) {
      ElMessage.error(e.message || (isEdit.value ? '更新失败' : '创建失败'))
    } finally {
      submitting.value = false
    }
  })
}

const onScopeChange = (val: string) => {
  if (val !== 'domain') {
    form.domain_id = undefined
  }
  if (val !== 'project') {
    form.project_id = undefined
  }
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(async () => {
  // 默认加载自定义权限
  filterForm.type = 'custom'
  activeTab.value = 'custom'

  loadPermissions()
  await Promise.all([
    loadDomains(),
    loadProjects()
  ])
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

.permission-tabs {
  margin-bottom: 16px;
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

.yaml-editor-container {
  width: 100%;
}

.yaml-textarea textarea {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  background-color: #2d333b;
  color: #c9d1d9;
  padding: 10px;
  border-radius: 4px;
}
</style>