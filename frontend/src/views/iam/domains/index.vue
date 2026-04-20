<template>
  <div class="domains-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">域</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建域
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
              @keyup.enter="loadDomains"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-divider direction="vertical" />
          <el-form-item label="状态">
            <el-select v-model="filterForm.enabled" placeholder="全部" clearable style="width: 100px">
              <el-option label="启用" :value="true" />
              <el-option label="禁用" :value="false" />
            </el-select>
          </el-form-item>
        </el-form>
        <div class="filter-actions">
          <el-button type="primary" @click="loadDomains">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="resetFilter">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </div>
      </div>

      <!-- 域列表 -->
      <el-table :data="domains" v-loading="loading" border stripe>
        <el-table-column prop="name" label="名称" min-width="150" show-overflow-tooltip>
          <template #default="{ row }">
            <el-button link @click="handleView(row)" class="name-link">
              {{ row.name }}
            </el-button>
            <el-tag v-if="row.id === 'default' || row.name.toLowerCase() === 'default'" type="warning" size="small" style="margin-left: 8px;">默认</el-tag>
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
        <el-table-column prop="project_count" label="项目数" width="80">
          <template #default="{ row }">
            {{ row.project_count || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="role_count" label="角色数" width="80">
          <template #default="{ row }">
            {{ row.role_count || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="policy_count" label="策略数" width="80">
          <template #default="{ row }">
            {{ row.policy_count || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="idp_count" label="认证源" width="80">
          <template #default="{ row }">
            {{ row.idp_count || 0 }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="handleView(row)">详情</el-button>
            <el-dropdown trigger="click" @command="(cmd) => handleCommand(cmd, row)">
              <el-button size="small" type="primary" link>
                更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="edit">编辑</el-dropdown-item>
                  <el-dropdown-item command="toggleEnable" :disabled="isDefaultDomain(row)">
                    {{ row.enabled ? '禁用' : '启用' }}
                  </el-dropdown-item>
                  <el-dropdown-item command="delete" :disabled="isDefaultDomain(row)" divided>
                    <span :style="isDefaultDomain(row) ? '' : 'color: #F56C6C'">
                      {{ isDefaultDomain(row) ? '删除（默认域不可删）' : '删除' }}
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
        @size-change="loadDomains"
        @current-change="loadDomains"
        class="pagination"
      />
    </el-card>

    <!-- 添加/编辑域对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑域' : '新建域'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="域名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入域名称" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入域描述" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 域详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="域详情" width="900px">

      <el-tabs v-model="detailTab" style="margin-top: 20px">
        <el-tab-pane label="基本详情" name="basic_info">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">{{ currentDomain?.id }}</el-descriptions-item>
            <el-descriptions-item label="域名称">
              {{ currentDomain?.name }}
              <el-tag v-if="currentDomain?.id === 'default'" type="warning" size="small" style="margin-left: 8px;">默认</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="描述" :span="2">{{ currentDomain?.description || '-' }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="currentDomain?.enabled ? 'success' : 'info'" size="small">
                {{ currentDomain?.enabled ? '启用' : '禁用' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="是否SSO域">
              <el-tag :type="currentDomain?.is_sso ? 'success' : 'info'" size="small">
                {{ currentDomain?.is_sso ? '是' : '否' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="用户数量">{{ currentDomain?.user_count || 0 }}</el-descriptions-item>
            <el-descriptions-item label="组数量">{{ currentDomain?.group_count || 0 }}</el-descriptions-item>
            <el-descriptions-item label="项目数量">{{ currentDomain?.project_count || 0 }}</el-descriptions-item>
            <el-descriptions-item label="角色数量">{{ currentDomain?.role_count || 0 }}</el-descriptions-item>
            <el-descriptions-item label="策略数量">{{ currentDomain?.policy_count || 0 }}</el-descriptions-item>
            <el-descriptions-item label="认证源数量">{{ currentDomain?.idp_count || 0 }}</el-descriptions-item>
            <el-descriptions-item label="可删除">
              <el-tag :type="currentDomain?.can_delete ? 'success' : 'danger'" size="small">
                {{ currentDomain?.can_delete ? '是' : '否' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="可更新">
              <el-tag :type="currentDomain?.can_update ? 'success' : 'danger'" size="small">
                {{ currentDomain?.can_update ? '是' : '否' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatDate(currentDomain?.created_at) }}</el-descriptions-item>
            <el-descriptions-item label="更新时间">{{ formatDate(currentDomain?.updated_at) }}</el-descriptions-item>
          </el-descriptions>

          <!-- 外部资源统计 -->
          <div v-if="currentDomain?.ext_resource" style="margin-top: 20px">
            <h4>外部资源统计</h4>
            <el-descriptions :column="3" border>
              <el-descriptions-item label="云账号">{{ currentDomain?.ext_resource?.cloudaccounts || 0 }}</el-descriptions-item>
              <el-descriptions-item label="云角色">{{ currentDomain?.ext_resource?.cloudroles || 0 }}</el-descriptions-item>
              <el-descriptions-item label="云用户">{{ currentDomain?.ext_resource?.cloudusers || 0 }}</el-descriptions-item>
              <el-descriptions-item label="宿主机">{{ currentDomain?.ext_resource?.hosts || 0 }}</el-descriptions-item>
              <el-descriptions-item label="存储">{{ currentDomain?.ext_resource?.storages || 0 }}</el-descriptions-item>
              <el-descriptions-item label="VPC">{{ currentDomain?.ext_resource?.vpcs || 0 }}</el-descriptions-item>
              <el-descriptions-item label="路由表">{{ currentDomain?.ext_resource?.route_tables || 0 }}</el-descriptions-item>
              <el-descriptions-item label="网线">{{ currentDomain?.ext_resource?.wires || 0 }}</el-descriptions-item>
              <el-descriptions-item label="代理设置">{{ currentDomain?.ext_resource?.proxysettings || 0 }}</el-descriptions-item>
            </el-descriptions>
          </div>
        </el-tab-pane>
        <el-tab-pane label="用户" name="users">
          <el-table :data="domainUsers" v-loading="usersLoading" border stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="用户名" min-width="120" />
            <el-table-column prop="display_name" label="显示名" min-width="120" />
            <el-table-column prop="email" label="邮箱" min-width="180" />
            <el-table-column prop="phone" label="手机号" width="150" />
            <el-table-column prop="enabled" label="启用状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.enabled ? 'success' : 'info'">
                  {{ row.enabled ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="180" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="项目" name="projects">
          <el-table :data="domainProjects" v-loading="projectsLoading" border stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="项目名称" min-width="150" />
            <el-table-column prop="description" label="项目描述" min-width="200" show-overflow-tooltip />
            <el-table-column prop="enabled" label="启用状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.enabled ? 'success' : 'info'">
                  {{ row.enabled ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="180" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="角色" name="roles">
          <el-table :data="domainRoles" v-loading="rolesLoading" border stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="角色名称" min-width="150" />
            <el-table-column prop="display_name" label="显示名" min-width="120" />
            <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
            <el-table-column prop="type" label="角色类型" width="100" />
            <el-table-column prop="enabled" label="启用状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.enabled ? 'success' : 'info'">
                  {{ row.enabled ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="180" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="云账号" name="cloud_accounts">
          <el-table :data="domainCloudAccounts" v-loading="cloudAccountsLoading" border stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="云账号名称" min-width="150" />
            <el-table-column prop="provider_type" label="云服务商" width="120" />
            <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'active' ? 'success' : 'info'">
                  {{ row.status === 'active' ? '活跃' : '非活跃' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="180" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="操作日志" name="operation_logs">
          <el-table :data="domainOperationLogs" v-loading="logsLoading" border stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="resource_name" label="资源名称" min-width="150" />
            <el-table-column prop="resource_type" label="资源类型" width="120" />
            <el-table-column prop="operation_type" label="操作类型" width="120" />
            <el-table-column prop="service_type" label="服务类型" width="120" />
            <el-table-column prop="operator" label="操作人" width="120" />
            <el-table-column prop="result" label="结果" width="100">
              <template #default="{ row }">
                <el-tag :type="row.result === 'success' ? 'success' : 'danger'">
                  {{ row.result === 'success' ? '成功' : '失败' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="180" />
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
import { Plus, ArrowDown, Search, Refresh } from '@element-plus/icons-vue'
import { Domain, User, Project, Role, AuthSource, OperationLog } from '@/types/iam'
import {
  getDomains,
  createDomain,
  updateDomain,
  deleteDomain,
  enableDomain,
  disableDomain,
  getDomainUsers,
  getDomainGroups,
  getDomainProjects,
  getDomainRoles,
  getDomainCloudAccounts,
  getDomainOperationLogs
} from '@/api/iam'

const domains = ref<Domain[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const detailTab = ref('basic_info')
const isEdit = ref(false)
const submitting = ref(false)
const currentDomain = ref<Domain | null>(null)
const domainUsers = ref<User[]>([])
const domainProjects = ref<Project[]>([])
const domainRoles = ref<Role[]>([])
const domainCloudAccounts = ref<AuthSource[]>([])
const domainOperationLogs = ref<OperationLog[]>([])
const usersLoading = ref(false)
const projectsLoading = ref(false)
const rolesLoading = ref(false)
const cloudAccountsLoading = ref(false)
const logsLoading = ref(false)
const formRef = ref<FormInstance>()

const filterForm = reactive({
  searchField: 'name',
  keyword: '',
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
  enabled: true
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入域名称', trigger: 'blur' }]
}

const isDefaultDomain = (domain: Domain) => {
  return domain.name.toLowerCase() === 'default' || domain.name.toLowerCase() === 'system';
};

const handleCommand = (command: string, row: Domain) => {
  switch (command) {
    case 'view': handleView(row); break
    case 'edit': handleEdit(row); break
    case 'toggleEnable': handleToggleEnable(row); break
    case 'delete': handleDelete(row); break
  }
};

const loadDomains = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit,
      details: true,
      show_fail_reason: true
    }
    // 根据搜索字段和关键词进行搜索
    if (filterForm.keyword) {
      if (filterForm.searchField === 'name') {
        params.name = filterForm.keyword
      } else if (filterForm.searchField === 'description') {
        params.description = filterForm.keyword
      }
    }
    if (filterForm.enabled !== undefined) params.enabled = filterForm.enabled

    const res = await getDomains(params)
    // Map response to add statistics fields
    domains.value = (res.items || res.data || []).map(item => ({
      ...item,
      user_count: item.user_count || 0,
      group_count: item.group_count || 0,
      project_count: item.project_count || 0,
      role_count: item.role_count || 0,
      policy_count: item.policy_count || 0,
      idp_count: item.idp_count || 0,
      is_sso: item.is_sso || false,
      ext_resource: item.ext_resource || null,
      can_delete: item.can_delete !== undefined ? item.can_delete : true,
      can_update: item.can_update !== undefined ? item.can_update : true
    }))
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载域列表失败')
  } finally {
    loading.value = false
  }
}

const loadDomainUsers = async (domainId: number) => {
  usersLoading.value = true
  try {
    const res = await getDomainUsers(domainId)
    domainUsers.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    usersLoading.value = false
  }
}

const loadDomainProjects = async (domainId: number) => {
  projectsLoading.value = true
  try {
    const res = await getDomainProjects(domainId)
    domainProjects.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    projectsLoading.value = false
  }
}

const loadDomainRoles = async (domainId: number) => {
  rolesLoading.value = true
  try {
    const res = await getDomainRoles(domainId)
    domainRoles.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    rolesLoading.value = false
  }
}

const loadDomainCloudAccounts = async (domainId: number) => {
  cloudAccountsLoading.value = true
  try {
    const res = await getDomainCloudAccounts(domainId)
    domainCloudAccounts.value = res.items || []
  } catch (e: any) {
    console.warn("Could not load domain cloud accounts:", e.message)
    domainCloudAccounts.value = []
  } finally {
    cloudAccountsLoading.value = false
  }
}

const loadDomainOperationLogs = async (domainId: number) => {
  logsLoading.value = true
  try {
    const res = await getDomainOperationLogs(domainId)
    domainOperationLogs.value = res.items || []
  } catch (e: any) {
    console.warn("Could not load domain operation logs:", e.message)
    domainOperationLogs.value = []
  } finally {
    logsLoading.value = false
  }
}

const resetFilter = () => {
  filterForm.searchField = 'name'
  filterForm.keyword = ''
  filterForm.name = ''
  filterForm.enabled = undefined
  pagination.page = 1
  loadDomains()
}

const handleCreate = () => {
  isEdit.value = false
  form.name = ''
  form.description = ''
  form.enabled = true
  dialogVisible.value = true
}

const handleEdit = (row: Domain) => {
  isEdit.value = true
  currentDomain.value = row
  form.name = row.name
  form.description = row.description || ''
  form.enabled = row.enabled
  dialogVisible.value = true
}

const handleView = async (row: Domain) => {
  currentDomain.value = row
  detailTab.value = 'basic_info' // Start with basic info tab
  await Promise.all([
    loadDomainUsers(row.id),
    loadDomainProjects(row.id),
    loadDomainRoles(row.id),
    loadDomainCloudAccounts(row.id),
    loadDomainOperationLogs(row.id)
  ])
  detailDialogVisible.value = true
}

const handleToggleEnable = async (row: Domain) => {
  try {
    const action = row.enabled ? '禁用' : '启用'
    await ElMessageBox.confirm(`确定要${action}该域吗？`, '提示', { type: 'warning' })

    if (row.enabled) {
      await disableDomain(row.id)
    } else {
      await enableDomain(row.id)
    }

    ElMessage.success(`${action}成功`)
    loadDomains()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || `${action}失败`)
    }
  }
}

const handleDelete = async (row: Domain) => {
  try {
    await ElMessageBox.confirm('确定要删除该域吗？', '提示', { type: 'warning' })
    await deleteDomain(row.id)
    ElMessage.success('删除成功')
    loadDomains()
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
        await updateDomain(currentDomain.value.id, form)
        ElMessage.success('更新成功')
      } else {
        await createDomain(form)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadDomains()
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

onMounted(() => {
  loadDomains()
})
</script>

<style scoped>
.domains-page {
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

.name-link {
  color: #409eff;
  text-decoration: underline;
  cursor: pointer;
}
</style>
