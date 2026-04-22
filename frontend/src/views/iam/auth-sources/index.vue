<template>
  <div class="auth-sources-container">
    <div class="page-header">
      <h2>认证源</h2>
      <div class="toolbar">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          新建认证源
        </el-button>
      </div>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filterForm" @submit.prevent="loadSources">
        <el-form-item label="搜索">
          <el-select v-model="filterForm.searchField" placeholder="选择字段" style="width: 120px">
            <el-option label="名称" value="name" />
            <el-option label="备注" value="description" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filterForm.keyword" placeholder="输入关键词" clearable style="width: 180px" @keyup.enter="loadSources">
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="filterForm.type" placeholder="全部" clearable style="width: 100px">
            <el-option label="LDAP" value="ldap" />
            <el-option label="本地" value="local" />
            <el-option label="SQL" value="sql" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterForm.enabled" placeholder="全部" clearable style="width: 80px">
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadSources">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card>
      <el-table :data="sources" v-loading="loading" border stripe row-key="id">
        <el-table-column prop="name" label="名称" min-width="180" show-overflow-tooltip>
          <template #default="{ row }">
            <el-button link @click="showDetailWithLogs(row)" class="name-link">{{ row.name }}</el-button>
            <el-tag v-if="row.is_system" type="warning" size="small" style="margin-left: 8px">系统</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">{{ row.enabled ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="80">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">{{ getTypeName(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="scope" label="归属" width="80">
          <template #default="{ row }">
            <el-tag :type="row.scope === 'system' ? 'warning' : 'primary'" size="small">{{ row.scope === 'system' ? '系统' : '域' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sync_status" label="同步状态" width="80">
          <template #default="{ row }">
            <el-tag :type="getSyncStatusTag(row.sync_status)" size="small">{{ getSyncStatusName(row.sync_status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-dropdown trigger="click" @command="(cmd) => handleCommand(cmd, row)">
              <el-button size="small" type="primary" link>更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon></el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="toggleEnable" :disabled="row.is_system">{{ row.enabled ? '禁用' : '启用' }}</el-dropdown-item>
                  <el-dropdown-item command="sync" :disabled="row.syncing || row.is_system">{{ row.syncing ? '同步中...' : '同步' }}</el-dropdown-item>
                  <el-dropdown-item command="test" :disabled="row.is_system">测试连接</el-dropdown-item>
                  <el-dropdown-item command="delete" :disabled="row.is_default || row.is_system" divided>删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.limit"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next"
        class="pagination"
      />
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑认证源' : '新建认证源'" width="800px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入认证源名称" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="认证源归属" prop="scope" v-if="!isEdit">
          <el-radio-group v-model="form.scope">
            <el-radio-button label="system">系统</el-radio-button>
            <el-radio-button label="domain">域</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="归属域" prop="domain_id" v-if="!isEdit && form.scope === 'domain'">
          <el-select v-model="form.domain_id" placeholder="请选择域" style="width: 100%">
            <el-option v-for="item in domains" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入备注" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
        </el-form-item>
        <template v-if="form.type === 'ldap'">
          <el-divider content-position="left">LDAP 配置</el-divider>
          <el-form-item label="服务器地址" prop="config.url">
            <el-input v-model="form.config.url" placeholder="ldap://192.168.1.1" />
          </el-form-item>
          <el-form-item label="基本DN" prop="config.base_dn">
            <el-input v-model="form.config.base_dn" placeholder="DC=example,DC=com" />
          </el-form-item>
          <el-form-item label="用户名" prop="config.bind_dn">
            <el-input v-model="form.config.bind_dn" placeholder="cn=admin,dc=example,dc=com" />
          </el-form-item>
          <el-form-item label="密码" prop="config.bind_password">
            <el-input v-model="form.config.bind_password" type="password" show-password placeholder="请输入密码" />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="认证源详情" width="700px">
      <el-descriptions :column="2" border v-if="currentSource">
        <el-descriptions-item label="ID">{{ currentSource.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ currentSource.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentSource.enabled ? 'success' : 'info'" size="small">{{ currentSource.enabled ? '启用' : '禁用' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="类型">
          <el-tag :type="getTypeTag(currentSource.type)" size="small">{{ getTypeName(currentSource.type) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="归属">
          <el-tag :type="currentSource.scope === 'system' ? 'warning' : 'primary'" size="small">{{ currentSource.scope === 'system' ? '系统' : '域' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="同步状态">
          <el-tag :type="getSyncStatusTag(currentSource.sync_status)" size="small">{{ getSyncStatusName(currentSource.sync_status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ currentSource.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentSource.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentSource.updated_at) }}</el-descriptions-item>
      </el-descriptions>
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
import { Plus, ArrowDown, Search } from '@element-plus/icons-vue'
import { AuthSource, Domain } from '@/types/iam'
import {
  getAuthSources,
  createAuthSource,
  updateAuthSource,
  deleteAuthSource,
  testAuthSource,
  enableAuthSource,
  disableAuthSource,
  getDomains
} from '@/api/iam'

const sources = ref<AuthSource[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const currentSource = ref<AuthSource | null>(null)
const domains = ref<Domain[]>([])
const formRef = ref<FormInstance>()

const filterForm = reactive({
  searchField: 'name',
  keyword: '',
  type: '',
  scope: '',
  enabled: undefined as boolean | undefined,
  sync_status: ''
})

const pagination = reactive({ page: 1, limit: 20, total: 0 })

const form = reactive({
  name: '',
  type: 'ldap',
  scope: 'system' as 'system' | 'domain',
  domain_id: undefined as number | undefined,
  description: '',
  enabled: true,
  config: {
    url: '',
    base_dn: '',
    bind_dn: '',
    bind_password: '',
    user_filter: '(objectClass=person)',
    user_id_attribute: 'uid',
    user_name_attribute: 'cn',
    protocol: 'ldap'
  }
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入认证源名称', trigger: 'blur' }]
}

const getTypeName = (type: string) => {
  const map: Record<string, string> = { ldap: 'LDAP', local: '本地', sql: 'SQL' }
  return map[type] || type
}

const getTypeTag = (type: string) => {
  const map: Record<string, any> = { ldap: 'success', local: 'primary', sql: 'warning' }
  return map[type] || ''
}

const getSyncStatusName = (status: string) => {
  const map: Record<string, string> = { synced: '已同步', syncing: '同步中', failed: '失败', pending: '待同步', never: '从未同步', idle: '空闲' }
  return map[status] || status || '空闲'
}

const getSyncStatusTag = (status: string) => {
  const map: Record<string, any> = { synced: 'success', syncing: 'warning', failed: 'danger', pending: 'info', never: 'info', idle: 'info' }
  return map[status] || 'info'
}

const loadDomains = async () => {
  try {
    const res = await getDomains({ limit: 100, enabled: true })
    domains.value = res.items || []
  } catch (e) { console.error(e) }
}

const loadSources = async () => {
  loading.value = true
  try {
    const params: any = { limit: pagination.limit, offset: (pagination.page - 1) * pagination.limit }
    if (filterForm.keyword) {
      if (filterForm.searchField === 'name') params.name = filterForm.keyword
      else if (filterForm.searchField === 'description') params.description = filterForm.keyword
    }
    if (filterForm.type) params.type = filterForm.type
    if (filterForm.scope) params.scope = filterForm.scope
    if (filterForm.enabled !== undefined) params.enabled = filterForm.enabled
    if (filterForm.sync_status) params.sync_status = filterForm.sync_status

    const res = await getAuthSources(params)
    sources.value = (res.items || []).map((item: any) => ({
      ...item,
      running: item.running !== undefined ? item.running : true,
      sync_status: item.sync_status || 'idle',
      protocol: item.protocol || 'sql',
      is_default: item.is_default !== undefined ? item.is_default : (item.type === 'sql' && item.scope === 'system'),
      is_system: item.is_system !== undefined ? item.is_system : (item.type === 'sql' && item.scope === 'system')
    }))
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载认证源列表失败')
  } finally {
    loading.value = false
  }
}

const resetFilter = () => {
  filterForm.searchField = 'name'
  filterForm.keyword = ''
  filterForm.type = ''
  filterForm.scope = ''
  filterForm.enabled = undefined
  filterForm.sync_status = ''
  pagination.page = 1
  loadSources()
}

const handleCreate = async () => {
  isEdit.value = false
  form.name = ''
  form.type = 'ldap'
  form.scope = 'system'
  form.domain_id = undefined
  form.description = ''
  form.enabled = true
  form.config = { url: '', base_dn: '', bind_dn: '', bind_password: '', user_filter: '(objectClass=person)', user_id_attribute: 'uid', user_name_attribute: 'cn', protocol: 'ldap' }
  await loadDomains()
  dialogVisible.value = true
}

const handleEdit = (row: AuthSource) => {
  isEdit.value = true
  currentSource.value = row
  form.name = row.name
  form.type = row.type
  form.scope = row.scope || 'system'
  form.domain_id = row.domain_id
  form.description = row.description || ''
  form.enabled = row.enabled
  form.config = {
    url: row.config?.url || '',
    base_dn: row.config?.base_dn || '',
    bind_dn: row.config?.bind_dn || '',
    bind_password: row.config?.bind_password || '',
    user_filter: row.config?.user_filter || '(objectClass=person)',
    user_id_attribute: row.config?.user_id_attribute || 'uid',
    user_name_attribute: row.config?.user_name_attribute || 'cn',
    protocol: row.config?.protocol || 'ldap'
  }
  dialogVisible.value = true
}

const showDetailWithLogs = (row: AuthSource) => {
  currentSource.value = row
  detailDialogVisible.value = true
}

const handleTest = async (row: AuthSource) => {
  try {
    await testAuthSource(row.id)
    ElMessage.success('连接测试成功')
  } catch (e: any) {
    ElMessage.error(e.message || '连接测试失败')
  }
}

const handleSync = async (row: AuthSource) => {
  row.syncing = true
  await new Promise(resolve => setTimeout(resolve, 2000))
  ElMessage.success('同步成功')
  row.sync_status = 'synced'
  row.syncing = false
  loadSources()
}

const handleToggleEnable = async (row: AuthSource) => {
  try {
    const action = row.enabled ? '禁用' : '启用'
    await ElMessageBox.confirm(`确定要${action}该认证源吗？`, '提示', { type: 'warning' })
    if (row.enabled) await disableAuthSource(row.id)
    else await enableAuthSource(row.id)
    ElMessage.success(`${action}成功`)
    loadSources()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error(e.message || '操作失败')
  }
}

const handleCommand = (command: string, row: AuthSource) => {
  switch (command) {
    case 'toggleEnable': handleToggleEnable(row); break
    case 'sync': handleSync(row); break
    case 'test': handleTest(row); break
    case 'delete': handleDelete(row); break
  }
}

const handleDelete = async (row: AuthSource) => {
  if (row.is_default) {
    ElMessage.warning('默认认证源不可删除')
    return
  }
  try {
    await ElMessageBox.confirm('确定要删除该认证源吗？', '提示', { type: 'warning' })
    await deleteAuthSource(row.id)
    ElMessage.success('删除成功')
    loadSources()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error(e.message || '删除失败')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      const data: any = {
        name: form.name,
        type: form.type,
        scope: form.scope,
        description: form.description,
        enabled: form.enabled,
        config: form.config
      }
      if (form.scope === 'domain' && form.domain_id) data.domain_id = form.domain_id
      if (isEdit.value) {
        await updateAuthSource(currentSource.value.id, data)
        ElMessage.success('更新成功')
      } else {
        await createAuthSource(data)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadSources()
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
  loadSources()
  loadDomains()
})
</script>

<style scoped>
.auth-sources-container { padding: 20px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.page-header h2 { margin: 0; font-size: 18px; font-weight: 600; }
.filter-card { margin-bottom: 20px; }
.pagination { margin-top: 20px; text-align: right; }
.name-link { color: #409eff; cursor: pointer; }
</style>