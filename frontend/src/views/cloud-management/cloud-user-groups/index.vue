<template>
  <div class="cloud-user-groups-container">
    <div class="page-header">
      <h2>云用户组</h2>
      <div class="toolbar">
        <el-button @click="loadGroups" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          新建
        </el-button>
        <el-dropdown trigger="click" @command="handleBatchCommand" :disabled="selectedIds.length === 0">
          <el-button>
            <el-icon><Operation /></el-icon>
            批量操作
            <el-badge v-if="selectedIds.length > 0" :value="selectedIds.length" type="primary" />
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="batchDelete" :icon="Delete">
                <span style="color: var(--el-color-danger)">批量删除</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- 筛选区 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="queryForm" @submit.prevent="loadGroups">
        <el-form-item label="名称">
          <el-input v-model="queryForm.name" placeholder="支持名称搜索" clearable style="width: 200px">
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryForm.status" placeholder="全部" clearable style="width: 140px">
            <el-option label="正常" value="active" />
            <el-option label="异常" value="inactive" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="queryForm.platform" placeholder="全部" clearable style="width: 140px">
            <el-option label="阿里云" value="alibaba" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item label="所属域">
          <el-select v-model="queryForm.domain_id" placeholder="全部" clearable style="width: 140px">
            <el-option v-for="domain in domains" :key="domain.id" :label="domain.name" :value="domain.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadGroups">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-table
      ref="tableRef"
      :data="groups"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="50" />
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" min-width="150">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleView(row)" class="name-link">
            {{ row.name }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="permissions" label="权限" width="120">
        <template #default="{ row }">
          {{ row.permissions || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="platform" label="平台" width="100">
        <template #default="{ row }">
          <el-tag :type="getPlatformType(row.platform)" size="small">
            {{ getPlatformText(row.platform) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="cloud_accounts" label="云账号" min-width="150">
        <template #default="{ row }">
          {{ row.cloud_accounts || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="shared_scope" label="共享范围" width="120">
        <template #default="{ row }">
          {{ getSharedScopeText(row.shared_scope) }}
        </template>
      </el-table-column>
      <el-table-column prop="owner_domain" label="所属域" width="120">
        <template #default="{ row }">
          {{ getDomainName(row.domain_id) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="140" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="handleView(row)">查看</el-button>
          <el-dropdown trigger="click" @command="(cmd: string) => handleDropdownCommand(cmd, row)" style="margin-left: 8px">
            <el-button size="small" link>
              更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="edit" :icon="EditPen">编辑</el-dropdown-item>
                <el-dropdown-item command="delete" :icon="Delete" divided>
                  <span style="color: var(--el-color-danger)">删除</span>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      class="pagination"
    />

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="showDialog"
      :title="isEdit ? '编辑云用户组' : '新建云用户组'"
      width="600px"
      destroy-on-close
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入云用户组名称" />
        </el-form-item>

        <el-form-item label="所属域" prop="domain_id">
          <el-select v-model="form.domain_id" placeholder="请选择所属域" style="width: 100%">
            <el-option v-for="domain in domains" :key="domain.id" :label="domain.name" :value="domain.id" />
          </el-select>
        </el-form-item>

        <el-form-item label="平台">
          <el-select v-model="form.platform" placeholder="请选择平台" style="width: 100%" clearable>
            <el-option label="阿里云" value="alibaba" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>

        <el-form-item label="权限">
          <el-input v-model="form.permissions" type="textarea" :rows="2" placeholder="请输入权限描述" />
        </el-form-item>

        <el-form-item label="共享范围">
          <el-select v-model="form.shared_scope" placeholder="请选择共享范围" style="width: 100%" clearable>
            <el-option label="私有" value="private" />
            <el-option label="域内共享" value="domain" />
            <el-option label="全局共享" value="global" />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 详情抽屉 -->
    <el-drawer
      v-model="showDetailDrawer"
      title="云用户组详情"
      size="50%"
      direction="rtl"
    >
      <div class="drawer-header" v-if="currentGroup">
        <div class="group-icon">
          <el-avatar :size="48" :style="{ backgroundColor: '#409EFF' }">
            <el-icon :size="32"><UserFilled /></el-icon>
          </el-avatar>
        </div>
        <div class="group-info">
          <h3>{{ currentGroup.name }}</h3>
          <div class="group-tags">
            <el-tag size="small" :type="getStatusType(currentGroup.status)">
              {{ getStatusText(currentGroup.status) }}
            </el-tag>
            <el-tag size="small" :type="getPlatformType(currentGroup.platform)">
              {{ getPlatformText(currentGroup.platform) }}
            </el-tag>
          </div>
        </div>
        <div class="quick-actions">
          <el-button size="small" @click="handleEdit(currentGroup)">
            <el-icon><EditPen /></el-icon>
            编辑
          </el-button>
          <el-button size="small" type="danger" @click="handleDelete(currentGroup)">
            <el-icon><Delete /></el-icon>
            删除
          </el-button>
        </div>
      </div>

      <!-- 基本信息 -->
      <el-descriptions :column="2" border style="margin-top: 20px" title="基本信息">
        <el-descriptions-item label="ID">{{ currentGroup?.id }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ getStatusText(currentGroup?.status) }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ getPlatformText(currentGroup?.platform) }}</el-descriptions-item>
        <el-descriptions-item label="权限">{{ currentGroup?.permissions || '-' }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ currentGroup?.cloud_accounts || '-' }}</el-descriptions-item>
        <el-descriptions-item label="共享范围">{{ getSharedScopeText(currentGroup?.shared_scope) }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ getDomainName(currentGroup?.domain_id) }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentGroup?.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentGroup?.updated_at) }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { Plus, Delete, UserFilled, ArrowDown, EditPen, Refresh, Operation, Search } from '@element-plus/icons-vue'
import { getCloudUserGroups, getCloudUserGroup, createCloudUserGroup, updateCloudUserGroup, deleteCloudUserGroup, batchDeleteCloudUserGroups } from '@/api/cloud-user-group'
import { getDomains } from '@/api/iam'
import type { CloudUserGroup, CreateCloudUserGroupRequest } from '@/api/cloud-user-group'
import EmptyState from '@/components/common/EmptyState.vue'

const groups = ref<CloudUserGroup[]>([])
const domains = ref<{ id: number; name: string }[]>([])
const loading = ref(false)
const showDialog = ref(false)
const showDetailDrawer = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const tableRef = ref()
const currentGroup = ref<CloudUserGroup | null>(null)
const selectedIds = ref<number[]>([])

// 分页数据
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 查询表单
const queryForm = reactive({
  name: '',
  status: '',
  platform: '',
  domain_id: undefined as number | undefined
})

const form = reactive<CreateCloudUserGroupRequest>({
  name: '',
  domain_id: 1,
  permissions: '',
  platform: '',
  shared_scope: ''
})

const rules = {
  name: [{ required: true, message: '请输入云用户组名称', trigger: 'blur' }],
  domain_id: [{ required: true, message: '请选择所属域', trigger: 'change' }]
}

const loadGroups = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (queryForm.name) params.name = queryForm.name
    if (queryForm.status) params.status = queryForm.status
    if (queryForm.platform) params.platform = queryForm.platform
    if (queryForm.domain_id) params.domain_id = queryForm.domain_id

    const res = await getCloudUserGroups(params)
    groups.value = res.items || []
    total.value = res.total || 0
  } catch (e) {
    console.error(e)
    ElMessage.error('加载云用户组列表失败')
  } finally {
    loading.value = false
  }
}

const loadDomains = async () => {
  try {
    const res = await getDomains()
    domains.value = (res.items || []).map(d => ({ id: d.id, name: d.name }))
  } catch (e) {
    console.error(e)
  }
}

const handleSelectionChange = (selection: CloudUserGroup[]) => {
  selectedIds.value = selection.map(g => g.id)
}

const handleBatchCommand = async (command: string) => {
  if (selectedIds.value.length === 0) {
    ElMessage.warning('请先选择云用户组')
    return
  }

  try {
    if (command === 'batchDelete') {
      await ElMessageBox.confirm(`确定要批量删除 ${selectedIds.value.length} 个云用户组吗？此操作不可恢复。`, '提示', { type: 'warning' })
      await batchDeleteCloudUserGroups(selectedIds.value)
      ElMessage.success('批量删除成功')
    }
    tableRef.value?.clearSelection()
    loadGroups()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error('批量操作失败')
    }
  }
}

const getStatusType = (status?: string) => {
  const map: Record<string, string> = {
    active: 'success',
    inactive: 'danger',
    pending: 'warning'
  }
  return map[status || ''] || 'info'
}

const getStatusText = (status?: string) => {
  const map: Record<string, string> = {
    active: '正常',
    inactive: '异常',
    pending: '待确认'
  }
  return map[status || ''] || status || '-'
}

const getPlatformType = (platform?: string) => {
  const map: Record<string, string> = {
    alibaba: 'warning',
    tencent: 'primary',
    aws: 'danger',
    azure: 'success'
  }
  return map[platform || ''] || 'info'
}

const getPlatformText = (platform?: string) => {
  const map: Record<string, string> = {
    alibaba: '阿里云',
    tencent: '腾讯云',
    aws: 'AWS',
    azure: 'Azure'
  }
  return map[platform || ''] || platform || '-'
}

const getSharedScopeText = (scope?: string) => {
  const map: Record<string, string> = {
    private: '私有',
    domain: '域内共享',
    global: '全局共享'
  }
  return map[scope || ''] || scope || '-'
}

const getDomainName = (domainId?: number) => {
  if (!domainId) return '-'
  const domain = domains.value.find(d => d.id === domainId)
  return domain?.name || `域${domainId}`
}

const formatDate = (date?: string) => {
  if (!date) return ''
  return new Date(date).toLocaleString('zh-CN')
}

const showCreateDialog = () => {
  isEdit.value = false
  resetForm()
  showDialog.value = true
}

const resetForm = () => {
  form.name = ''
  form.domain_id = domains.value[0]?.id || 1
  form.permissions = ''
  form.platform = ''
  form.shared_scope = ''
}

const handleView = async (row: CloudUserGroup) => {
  try {
    const detail = await getCloudUserGroup(row.id)
    currentGroup.value = detail
    showDetailDrawer.value = true
  } catch (e) {
    console.error(e)
    ElMessage.error('获取详情失败')
  }
}

const handleEdit = async (row: CloudUserGroup) => {
  try {
    const detail = await getCloudUserGroup(row.id)
    isEdit.value = true
    form.name = detail.name
    form.domain_id = detail.domain_id || 1
    form.permissions = detail.permissions || ''
    form.platform = detail.platform || ''
    form.shared_scope = detail.shared_scope || ''
    currentGroup.value = detail
    showDialog.value = true
  } catch (e) {
    console.error(e)
    ElMessage.error('获取云用户组详情失败')
  }
}

const handleDelete = async (row: CloudUserGroup) => {
  try {
    await ElMessageBox.confirm(`确定要删除云用户组 "${row.name}" 吗？此操作不可恢复。`, '提示', { type: 'warning' })
    await deleteCloudUserGroup(row.id)
    ElMessage.success('删除成功')
    showDetailDrawer.value = false
    loadGroups()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
      ElMessage.error('删除失败')
    }
  }
}

const handleDropdownCommand = (command: string, row: CloudUserGroup) => {
  switch (command) {
    case 'edit':
      handleEdit(row)
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      if (isEdit.value && currentGroup.value) {
        await updateCloudUserGroup(currentGroup.value.id, form)
        ElMessage.success('更新成功')
      } else {
        await createCloudUserGroup(form)
        ElMessage.success('创建成功')
      }
      showDialog.value = false
      loadGroups()
    } catch (e) {
      console.error(e)
      ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
    } finally {
      submitting.value = false
    }
  })
}

const resetQuery = () => {
  queryForm.name = ''
  queryForm.status = ''
  queryForm.platform = ''
  queryForm.domain_id = undefined
  currentPage.value = 1
  loadGroups()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadGroups()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadGroups()
}

onMounted(() => {
  loadGroups()
  loadDomains()
})
</script>

<style scoped>
.cloud-user-groups-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.toolbar {
  display: flex;
  gap: 8px;
  align-items: center;
}

.filter-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}

.name-link {
  font-weight: 500;
}

/* 抽屉顶部区域 */
.drawer-header {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--el-border-color);
}

.group-icon {
  margin-right: 16px;
}

.group-info h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.group-tags {
  display: flex;
  gap: 8px;
  margin-top: 8px;
  align-items: center;
}

.quick-actions {
  margin-left: auto;
  display: flex;
  gap: 8px;
}
</style>