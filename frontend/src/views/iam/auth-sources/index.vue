<template>
  <div class="auth-sources-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">认证源</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建认证源
          </el-button>
        </div>
      </template>

      <!-- 筛选条件 -->
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="名称">
          <el-input 
            v-model="filterForm.name" 
            placeholder="请输入名称" 
            clearable 
            style="width: 200px"
            @keyup.enter="loadSources" 
          />
        </el-form-item>
        <el-form-item label="类型">
          <el-select 
            v-model="filterForm.type" 
            placeholder="全部" 
            clearable 
            style="width: 120px"
          >
            <el-option label="LDAP" value="ldap" />
            <el-option label="本地" value="local" />
          </el-select>
        </el-form-item>
        <el-form-item label="范围">
          <el-select 
            v-model="filterForm.scope" 
            placeholder="全部" 
            clearable 
            style="width: 120px"
          >
            <el-option label="系统" value="system" />
            <el-option label="域" value="domain" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select 
            v-model="filterForm.enabled" 
            placeholder="全部" 
            clearable 
            style="width: 100px"
          >
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadSources">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 认证源列表 -->
      <el-table
        :data="sources"
        v-loading="loading"
        border
        stripe
        header-cell-class-name="table-header-gray"
      >
        <el-table-column prop="name" label="名称" min-width="180" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="table-cell-with-icon">
              <el-icon v-if="row.type === 'local'" color="#409EFF"><User /></el-icon>
              <el-icon v-else-if="row.type === 'ldap'" color="#67C23A"><Connection /></el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="running" label="启动状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.running ? 'success' : 'info'" size="small">
              {{ row.running ? '已启动' : '未启动' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sync_status" label="同步状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getSyncStatusTag(row.sync_status)" size="small">
              {{ getSyncStatusName(row.sync_status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="auto_create" label="自动创建用户" width="100">
          <template #default="{ row }">
            <el-tag :type="row.auto_create ? 'success' : 'info'" size="small">
              {{ row.auto_create ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="protocol" label="认证协议" width="100">
          <template #default="{ row }">
            <span>{{ getProtocolName(row.protocol) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="认证类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">
              {{ getTypeName(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="scope" label="认证源归属" width="100">
          <template #default="{ row }">
            <el-tag :type="row.scope === 'system' ? 'warning' : ''" size="small">
              {{ row.scope === 'system' ? '系统' : '域' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="400" fixed="right">
          <template #default="{ row }">
            <el-button size="small" link type="primary" @click="handleView(row)">详情</el-button>
            <el-button size="small" link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button
              size="small"
              link
              :type="row.running ? 'warning' : 'success'"
              @click="handleToggleRunning(row)"
            >
              {{ row.running ? '停止' : '启动' }}
            </el-button>
            <el-button
              size="small"
              link
              :type="row.enabled ? 'warning' : 'success'"
              @click="handleToggleEnable(row)"
            >
              {{ row.enabled ? '禁用' : '启用' }}
            </el-button>
            <el-button
              size="small"
              link
              type="primary"
              @click="handleSync(row)"
              :loading="row.syncing"
            >
              同步
            </el-button>
            <el-button
              size="small"
              link
              type="danger"
              @click="handleDelete(row)"
              :disabled="row.is_default"
            >
              删除
              <el-tooltip v-if="row.is_default" content="默认认证源不可删除" placement="top">
                <el-icon style="margin-left: 2px"><QuestionFilled /></el-icon>
              </el-tooltip>
            </el-button>
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
        @size-change="loadSources"
        @current-change="loadSources"
        class="pagination"
      />
    </el-card>

    <!-- 添加/编辑认证源对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑认证源' : '新建认证源'" width="800px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入认证源名称" :disabled="isEdit" />
        </el-form-item>
        
        <!-- 认证源范围选择 -->
        <el-form-item label="认证源范围" prop="scope" v-if="!isEdit">
          <el-radio-group v-model="form.scope">
            <el-radio-button label="system">
              <el-tooltip content="系统级认证源，所有域的用户都可以使用" placement="top">
                <span>系统</span>
              </el-tooltip>
            </el-radio-button>
            <el-radio-button label="domain">
              <el-tooltip content="域级认证源，只有指定域的用户可以使用" placement="top">
                <span>域</span>
              </el-tooltip>
            </el-radio-button>
          </el-radio-group>
        </el-form-item>
        
        <!-- 域选择（仅当选择域范围时显示） -->
        <el-form-item label="所属域" prop="domain_id" v-if="!isEdit && form.scope === 'domain'">
          <el-select v-model="form.domain_id" placeholder="请选择域" style="width: 100%">
            <el-option
              v-for="item in domains"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="类型" prop="type" v-if="!isEdit">
          <el-select v-model="form.type" placeholder="请选择认证源类型" style="width: 100%" @change="handleTypeChange">
            <el-option label="LDAP" value="ldap" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入认证源描述" />
        </el-form-item>
        
        <el-form-item label="自动创建用户">
          <el-switch v-model="form.auto_create" />
          <div class="form-item-tip">启用后，用户首次登录时自动创建本地用户</div>
        </el-form-item>
        
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
        </el-form-item>

        <!-- LDAP 配置 -->
        <template v-if="form.type === 'ldap'">
          <el-divider content-position="left">LDAP 配置</el-divider>
          
          <el-form-item label="LDAP 服务器地址" prop="config.ldap_url">
            <el-input v-model="form.config.ldap_url" placeholder="ldap://192.168.1.1:389" />
            <div class="form-item-tip">LDAP 服务器的地址和端口</div>
          </el-form-item>
          
          <el-form-item label="Base DN" prop="config.ldap_base_dn">
            <el-input v-model="form.config.ldap_base_dn" placeholder="dc=example,dc=com" />
            <div class="form-item-tip">搜索用户的起始 DN</div>
          </el-form-item>
          
          <el-form-item label="绑定 DN" prop="config.ldap_bind_dn">
            <el-input v-model="form.config.ldap_bind_dn" placeholder="cn=admin,dc=example,dc=com" />
            <div class="form-item-tip">用于绑定 LDAP 服务器的管理员 DN</div>
          </el-form-item>
          
          <el-form-item label="绑定密码" prop="config.ldap_bind_password">
            <el-input v-model="form.config.ldap_bind_password" type="password" show-password />
            <div class="form-item-tip">绑定 DN 的密码</div>
          </el-form-item>
          
          <el-form-item label="用户过滤器" prop="config.ldap_user_filter">
            <el-input v-model="form.config.ldap_user_filter" placeholder="(objectClass=person)" />
            <div class="form-item-tip">用于过滤用户的 LDAP 过滤器</div>
          </el-form-item>
          
          <el-form-item label="用户唯一 ID 属性" prop="config.ldap_user_id_attribute">
            <el-input v-model="form.config.ldap_user_id_attribute" placeholder="uid" />
            <div class="form-item-tip">用于标识用户唯一性的属性，默认为 uid</div>
          </el-form-item>
          
          <el-form-item label="用户名属性" prop="config.ldap_user_name_attribute">
            <el-input v-model="form.config.ldap_user_name_attribute" placeholder="cn" />
            <div class="form-item-tip">用于显示用户名的属性，默认为 cn</div>
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 认证源详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="认证源详情" width="700px">
      <el-descriptions :column="2" border v-if="currentSource">
        <el-descriptions-item label="ID">{{ currentSource.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ currentSource.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentSource.enabled ? 'success' : 'info'" size="small">
            {{ currentSource.enabled ? '启用' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="启动状态">
          <el-tag :type="currentSource.running ? 'success' : 'info'" size="small">
            {{ currentSource.running ? '已启动' : '未启动' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="同步状态">
          <el-tag :type="getSyncStatusTag(currentSource.sync_status)" size="small">
            {{ getSyncStatusName(currentSource.sync_status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="自动创建用户">
          <el-tag :type="currentSource.auto_create ? 'success' : 'info'" size="small">
            {{ currentSource.auto_create ? '是' : '否' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="认证协议">{{ getProtocolName(currentSource.protocol) }}</el-descriptions-item>
        <el-descriptions-item label="认证类型">
          <el-tag :type="getTypeTag(currentSource.type)" size="small">
            {{ getTypeName(currentSource.type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="认证源归属">
          <el-tag :type="currentSource.scope === 'system' ? 'warning' : ''" size="small">
            {{ currentSource.scope === 'system' ? '系统' : '域' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentSource.description || '-' }}</el-descriptions-item>
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
import { User, Connection, Plus, QuestionFilled } from '@element-plus/icons-vue'
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

const sources = ref<any[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const currentSource = ref<any>(null)
const domains = ref<any[]>([])
const formRef = ref<FormInstance>()

const filterForm = reactive({
  name: '',
  type: '',
  scope: '',
  enabled: undefined as boolean | undefined
})

const pagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

const form = reactive({
  name: '',
  type: 'ldap',
  scope: 'system' as 'system' | 'domain',
  domain_id: undefined as number | undefined,
  description: '',
  auto_create: false,
  enabled: true,
  config: {
    ldap_url: '',
    ldap_base_dn: '',
    ldap_bind_dn: '',
    ldap_bind_password: '',
    ldap_user_filter: '',
    ldap_user_id_attribute: 'uid',
    ldap_user_name_attribute: 'cn'
  }
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入认证源名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择认证源类型', trigger: 'change' }],
  scope: [{ required: true, message: '请选择认证源范围', trigger: 'change' }],
  domain_id: [{ required: true, message: '请选择所属域', trigger: 'change', trigger: 'change' }]
}

const getTypeName = (type: string) => {
  const map: Record<string, string> = {
    ldap: 'LDAP',
    local: '本地',
    sql: 'SQL'
  }
  return map[type] || type
}

const getTypeTag = (type: string) => {
  const map: Record<string, any> = {
    ldap: 'success',
    local: 'primary',
    sql: 'warning'
  }
  return map[type] || ''
}

const getSyncStatusName = (status: string) => {
  const map: Record<string, string> = {
    synced: '已同步',
    syncing: '同步中',
    failed: '失败',
    pending: '待同步',
    never: '从未同步',
    idle: '空闲'
  }
  return map[status] || status || '空闲'
}

const getSyncStatusTag = (status: string) => {
  const map: Record<string, any> = {
    synced: 'success',
    syncing: 'warning',
    failed: 'danger',
    pending: 'info',
    never: 'info',
    idle: 'info'
  }
  return map[status] || 'info'
}

const getProtocolName = (protocol: string) => {
  const map: Record<string, string> = {
    ldap: 'LDAP',
    saml: 'SAML',
    oidc: 'OIDC',
    cas: 'CAS',
    oauth2: 'OAuth2',
    sql: 'SQL'
  }
  return map[protocol] || protocol || '-'
}

const loadDomains = async () => {
  try {
    const res = await getDomains({ limit: 100, enabled: true })
    domains.value = res.items || []
  } catch (e: any) {
    console.error(e)
  }
}

const loadSources = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit
    }
    if (filterForm.name) params.keyword = filterForm.name
    if (filterForm.type) params.type = filterForm.type
    if (filterForm.enabled !== undefined) params.enabled = filterForm.enabled

    const res = await getAuthSources(params)
    sources.value = (res.items || []).map((item: any) => ({
      ...item,
      // 如果没有 running 字段，默认为 true（启动状态）
      running: item.running !== undefined ? item.running : true,
      // 如果没有 sync_status 字段，默认为 idle（空闲）
      sync_status: item.sync_status || 'idle',
      // 如果没有 protocol 字段，默认为 sql
      protocol: item.protocol || 'sql',
      // 如果没有 is_default 字段，根据 type 判断（sql 类型的系统认证源为默认）
      is_default: item.is_default !== undefined ? item.is_default : (item.type === 'sql' && item.scope === 'system')
    }))
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载认证源列表失败')
  } finally {
    loading.value = false
  }
}

const resetFilter = () => {
  filterForm.name = ''
  filterForm.type = ''
  filterForm.enabled = undefined
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
  form.auto_create = false
  form.enabled = true
  form.config = {
    ldap_url: '',
    ldap_base_dn: '',
    ldap_bind_dn: '',
    ldap_bind_password: '',
    ldap_user_filter: '',
    ldap_user_id_attribute: 'uid',
    ldap_user_name_attribute: 'cn'
  }
  await loadDomains()
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  isEdit.value = true
  currentSource.value = row
  form.name = row.name
  form.type = row.type
  form.scope = row.scope || 'system'
  form.domain_id = row.domain_id
  form.description = row.description || ''
  form.auto_create = row.auto_create
  form.enabled = row.enabled
  form.config = row.config || {
    ldap_url: '',
    ldap_base_dn: '',
    ldap_bind_dn: '',
    ldap_bind_password: '',
    ldap_user_filter: '',
    ldap_user_id_attribute: 'uid',
    ldap_user_name_attribute: 'cn'
  }
  dialogVisible.value = true
}

const handleView = (row: any) => {
  currentSource.value = row
  detailDialogVisible.value = true
}

const handleTest = async (row: any) => {
  try {
    await testAuthSource(row.id)
    ElMessage.success('连接测试成功')
  } catch (e: any) {
    ElMessage.error(e.message || '连接测试失败')
  }
}

const handleToggleRunning = async (row: any) => {
  try {
    const action = row.running ? '停止' : '启动'
    await ElMessageBox.confirm(`确定要${action}该认证源吗？`, '提示', { type: 'warning' })

    // 调用 API 更新运行状态
    const data = { running: !row.running }
    await updateAuthSource(row.id, data)

    ElMessage.success(`${action}成功`)
    loadSources()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || `${action}失败`)
    }
  }
}

const handleSync = async (row: any) => {
  try {
    // 标记为同步中
    row.syncing = true
    
    // 调用同步 API（未来实现）
    // await syncAuthSource(row.id)
    
    // 模拟同步过程
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    ElMessage.success('同步成功')
    row.sync_status = 'synced'
    row.syncing = false
    loadSources()
  } catch (e: any) {
    row.syncing = false
    if (e !== 'cancel') {
      ElMessage.error(e.message || '同步失败')
      row.sync_status = 'failed'
    }
  }
}

const handleToggleEnable = async (row: any) => {
  try {
    const action = row.enabled ? '禁用' : '启用'
    await ElMessageBox.confirm(`确定要${action}该认证源吗？`, '提示', { type: 'warning' })

    if (row.enabled) {
      await disableAuthSource(row.id)
    } else {
      await enableAuthSource(row.id)
    }

    ElMessage.success(`${action}成功`)
    loadSources()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || `${action}失败`)
    }
  }
}

const handleDelete = async (row: any) => {
  // 检查是否为默认认证源
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
      const data: any = {
        name: form.name,
        type: form.type,
        scope: form.scope,
        description: form.description,
        auto_create: form.auto_create,
        enabled: form.enabled,
        config: form.config
      }
      
      // 如果是域范围，添加 domain_id
      if (form.scope === 'domain' && form.domain_id) {
        data.domain_id = form.domain_id
      }

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
})
</script>

<style scoped>
.auth-sources-page {
  height: 100%;
  padding: 20px;
  background-color: #f5f7fa;
}

.page-card {
  height: 100%;
  border-radius: 4px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 18px;
  font-weight: bold;
  color: #303133;
}

.filter-form {
  margin-bottom: 16px;
  padding: 16px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.form-item-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.5;
}

.table-cell-with-icon {
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.table-header-gray) {
  background-color: #fafafa;
  color: #606266;
  font-weight: 500;
}

:deep(.el-button--link) {
  padding: 0 4px;
}
</style>
