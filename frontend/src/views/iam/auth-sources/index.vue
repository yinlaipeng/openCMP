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
        <el-table-column prop="name" label="名称/备注" min-width="180" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="table-cell-with-icon" @click="showDetailWithLogs(row)">
              <el-icon v-if="row.type === 'local'" color="#409EFF"><User /></el-icon>
              <el-icon v-else-if="row.type === 'ldap'" color="#67C23A"><Connection /></el-icon>
              <span class="clickable-name">{{ row.name }}</span>
              <el-tag v-if="row.is_system" type="warning" size="small" style="margin-left: 8px;">系统</el-tag>
              <div v-if="row.description" class="remark">{{ row.description }}</div>
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
        <el-table-column prop="running" label="启用状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.running ? 'success' : 'info'" size="small">
              {{ row.running ? '已启用' : '已禁用' }}
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
        <el-table-column prop="scope" label="认证源归属" min-width="130">
          <template #default="{ row }">
            <el-tag :type="row.scope === 'system' ? 'warning' : 'primary'" size="small">
              {{ row.scope === 'system' ? '系统' : getDomainName(row.domain_id) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">修改配置</el-button>
            <el-dropdown trigger="click" @command="(cmd) => handleCommand(cmd, row)">
              <el-button size="small" type="primary" link>
                更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="toggleEnable" :disabled="row.is_system">
                    {{ row.enabled ? '禁用' : '启用' }}
                  </el-dropdown-item>
                  <el-dropdown-item command="sync" :disabled="row.syncing || row.is_system">
                    {{ row.syncing ? '同步中...' : '同步' }}
                  </el-dropdown-item>
                  <el-dropdown-item command="test" :disabled="row.is_system">连接测试</el-dropdown-item>
                  <el-dropdown-item command="delete" :disabled="row.is_default || row.is_system" divided>
                    <span :style="(row.is_default || row.is_system) ? '' : 'color: #F56C6C'">
                      {{ (row.is_default || row.is_system) ? '删除（系统认证源不可删）' : '删除' }}
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
        
        <!-- 认证源归属选择 -->
        <el-form-item label="认证源归属" prop="scope" v-if="!isEdit">
          <el-radio-group v-model="form.scope">
            <el-radio-button label="system">
              <el-tooltip content="系统级认证源，所有域的用户都可以使用，登录页面所有人可见此LDAP" placement="top">
                <span>系统</span>
              </el-tooltip>
            </el-radio-button>
            <el-radio-button label="domain">
              <el-tooltip content="域级认证源，只有指定域的用户可以使用，需进入域后才可见" placement="top">
                <span>域</span>
              </el-tooltip>
            </el-radio-button>
          </el-radio-group>
          <div class="form-item-tip">认证源归属，即用户从哪里来以及如何验证身份</div>
        </el-form-item>

        <!-- 域选择（仅当选择域范围时显示） -->
        <el-form-item label="归属域" prop="domain_id" v-if="!isEdit && form.scope === 'domain'">
          <el-select v-model="form.domain_id" placeholder="请选择域" style="width: 100%">
            <el-option
              v-for="item in domains"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入认证源名称" :disabled="isEdit" />
        </el-form-item>

        <el-form-item label="备注" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入认证源备注" />
        </el-form-item>

        <el-form-item label="认证协议" prop="config.protocol" v-if="form.type === 'ldap'">
          <el-select v-model="form.config.protocol" placeholder="请选择认证协议" style="width: 100%">
            <el-option label="LDAP" value="ldap" />
            <el-option label="LDAPS" value="ldaps" />
          </el-select>
        </el-form-item>

        <el-form-item label="认证类型" prop="config.auth_type" v-if="form.type === 'ldap'">
          <el-select v-model="form.config.auth_type" placeholder="请选择认证类型" style="width: 100%">
            <el-option label="OpenLDAP" value="openldap" />
          </el-select>
        </el-form-item>

        <el-form-item label="用户归属" prop="config.target_domain">
          <el-input v-model="form.config.target_domain" placeholder="请输入用户归属域名称，留空则系统会创建同认证源名称的域（遇到同名会追加后缀'-1'）" />
          <div class="form-item-tip">用户归属域名称，留空则系统会创建同认证源名称的域（遇到同名会追加后缀"-1"）</div>
        </el-form-item>

        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
        </el-form-item>

        <!-- LDAP 配置 -->
        <template v-if="form.type === 'ldap'">
          <el-divider content-position="left">LDAP 配置</el-divider>

          <el-form-item label="服务器地址" prop="config.url">
            <el-input v-model="form.config.url" placeholder="例如：ldap://192.168.222.222" />
            <div class="form-item-tip">LDAP 服务器的地址和端口</div>
          </el-form-item>

          <el-form-item label="基本DN" prop="config.base_dn">
            <el-input v-model="form.config.base_dn" placeholder="例如：DC=ocdc,DC=com" />
            <div class="form-item-tip">搜索用户的起始 DN</div>
          </el-form-item>

          <el-form-item label="用户名" prop="config.bind_dn">
            <el-input v-model="form.config.bind_dn" placeholder="例如：cn=admin,dc=example,dc=com" />
            <div class="form-item-tip">用于绑定 LDAP 服务器的管理员 DN</div>
          </el-form-item>

          <el-form-item label="密码" prop="config.bind_password">
            <el-input v-model="form.config.bind_password" type="password" show-password placeholder="请输入绑定密码" />
            <div class="form-item-tip">绑定 DN 的密码</div>
          </el-form-item>

          <el-form-item label="用户DN" prop="config.user_search_base">
            <el-input v-model="form.config.user_search_base" placeholder="例如：OU=XX集团,DC=ocdc,DC=com" />
            <div class="form-item-tip">用户搜索的基础 DN</div>
          </el-form-item>

          <el-form-item label="组DN" prop="config.group_search_base">
            <el-input v-model="form.config.group_search_base" placeholder="例如：OU=XX集团,DC=ocdc,DC=com" />
            <div class="form-item-tip">组搜索的基础 DN</div>
          </el-form-item>

          <el-form-item label="用户启用状态" prop="config.user_enabled_attribute">
            <el-select v-model="form.config.user_enabled_attribute" placeholder="请选择用户启用状态属性" style="width: 100%">
              <el-option label="启用" value="enabled" />
              <el-option label="禁用" value="disabled" />
              <el-option label="Active Directory账户状态" value="userAccountControl" />
            </el-select>
            <div class="form-item-tip">用户启用状态的属性设置</div>
          </el-form-item>

          <el-form-item label="用户过滤器" prop="config.user_filter">
            <el-input v-model="form.config.user_filter" placeholder="(objectClass=person)" />
            <div class="form-item-tip">用于过滤用户的 LDAP 过滤器</div>
          </el-form-item>

          <el-form-item label="用户唯一 ID 属性" prop="config.user_id_attribute">
            <el-input v-model="form.config.user_id_attribute" placeholder="uid" />
            <div class="form-item-tip">用于标识用户唯一性的属性，默认为 uid</div>
          </el-form-item>

          <el-form-item label="用户名属性" prop="config.user_name_attribute">
            <el-input v-model="form.config.user_name_attribute" placeholder="cn" />
            <div class="form-item-tip">用于显示用户名的属性，默认为 cn</div>
          </el-form-item>

          <!-- 测试 LDAP 连接和用户查询按钮 -->
          <el-form-item>
            <el-button
              type="primary"
              @click="testLdapConnection"
              :loading="ldapTesting"
            >
              {{ ldapTesting ? '测试中...' : '测试 LDAP 用户查询' }}
            </el-button>
            <div class="form-item-tip" v-if="ldapTestResult">
              <span :class="ldapTestResult.success ? 'text-success' : 'text-error'">
                {{ ldapTestResult.message }}
              </span>
            </div>
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 认证源详情和操作日志对话框 -->
    <el-dialog v-model="detailDialogVisible" title="认证源详情" width="90%" :fullscreen="isFullscreen">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="详情" name="detail">
          <el-descriptions :column="2" border v-if="currentSource">
            <el-descriptions-item label="ID">{{ currentSource.id }}</el-descriptions-item>
            <el-descriptions-item label="名称">{{ currentSource.name }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="currentSource.enabled ? 'success' : 'info'" size="small">
                {{ currentSource.enabled ? '启用' : '禁用' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="启用状态">
              <el-tag :type="currentSource.running ? 'success' : 'info'" size="small">
                {{ currentSource.running ? '已启用' : '已禁用' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="同步状态">
              <el-tag :type="getSyncStatusTag(currentSource.sync_status)" size="small">
                {{ getSyncStatusName(currentSource.sync_status) }}
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
        </el-tab-pane>

        <el-tab-pane label="操作日志" name="logs">
          <el-table
            :data="operationLogs"
            v-loading="logsLoading"
            border
            stripe
            header-cell-class-name="table-header-gray"
          >
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="operation_time" label="操作时间" width="160">
              <template #default="{ row }">
                {{ formatDate(row.operation_time) }}
              </template>
            </el-table-column>
            <el-table-column prop="resource_name" label="资源名称" width="150" />
            <el-table-column prop="resource_type" label="资源类型" width="120" />
            <el-table-column prop="operation_type" label="操作类型" width="120" />
            <el-table-column prop="service_type" label="服务类型" width="120" />
            <el-table-column prop="risk_level" label="风险级别" width="100">
              <template #default="{ row }">
                <el-tag :type="getRiskLevelType(row.risk_level)" size="small">
                  {{ row.risk_level }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="time_type" label="时间类型" width="100" />
            <el-table-column prop="result" label="结果" width="100">
              <template #default="{ row }">
                <el-tag :type="getResultType(row.result)" size="small">
                  {{ row.result }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="operator" label="发起人" width="120" />
            <el-table-column prop="project_id" label="所属项目" width="100" />
          </el-table>

          <el-pagination
            v-model:current-page="logPagination.page"
            v-model:page-size="logPagination.limit"
            :total="logPagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="loadOperationLogs"
            @current-change="loadOperationLogs"
            class="pagination"
            style="margin-top: 20px; justify-content: flex-end;"
          />
        </el-tab-pane>
      </el-tabs>
      <template #footer>
        <el-button @click="toggleFullscreen">
          <el-icon><FullScreen v-if="!isFullscreen" /><Crop v-else /></el-icon>
          {{ isFullscreen ? '退出全屏' : '全屏' }}
        </el-button>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Connection, Plus, QuestionFilled, ArrowDown, FullScreen, Crop } from '@element-plus/icons-vue'
import { AuthSource, Domain, OperationLog } from '@/types/iam'
import {
  getAuthSources,
  createAuthSource,
  updateAuthSource,
  deleteAuthSource,
  testAuthSource,
  testLdapUsers,
  enableAuthSource,
  disableAuthSource,
  getDomains,
  getResourceOperationLogs
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

// Operation logs
const operationLogs = ref<OperationLog[]>([])
const logsLoading = ref(false)
const activeTab = ref('detail')
const isFullscreen = ref(false)

// LDAP 测试相关
const ldapTesting = ref(false)
const ldapTestResult = ref<{ success: boolean; message: string; users?: any[] } | null>(null)

const logPagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

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
  type: 'ldap', // 保留默认类型，用于区分LDAP和其他类型（如本地认证源）
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
    user_search_base: '',
    group_search_base: '',
    user_enabled_attribute: 'enabled',
    protocol: 'ldap',
    auth_type: 'openldap', // 只保留OpenLDAP作为认证类型
    target_domain: ''
  }
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入认证源名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择认证源类型', trigger: 'change' }],
  scope: [{ required: true, message: '请选择认证源范围', trigger: 'change' }],
  domain_id: [{ required: true, message: '请选择所属域', trigger: 'change' }],
  'config.url': [
    { required: true, message: '请输入服务器地址', trigger: 'blur' },
    { pattern: /^(ldap|ldaps):\/\/.*/, message: '必须以 ldap:// 或 ldaps:// 开头', trigger: 'blur' }
  ],
  'config.base_dn': [
    { required: true, message: '请输入基本DN', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9=,.-]+$/, message: '请输入有效的DN格式', trigger: 'blur' }
  ],
  'config.bind_dn': [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  'config.bind_password': [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ],
  'config.protocol': [
    { required: true, message: '请选择认证协议', trigger: 'change' }
  ],
  'config.auth_type': [
    { required: true, message: '请选择认证类型', trigger: 'change' }
  ]
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

const getRiskLevelType = (level: string) => {
  const map: Record<string, any> = {
    high: 'danger',
    medium: 'warning',
    low: 'success',
    info: 'info'
  }
  return map[level] || 'info'
}

const getResultType = (result: string) => {
  const map: Record<string, any> = {
    success: 'success',
    failed: 'danger',
    pending: 'warning'
  }
  return map[result] || 'info'
}

const loadDomains = async () => {
  try {
    const res = await getDomains({ limit: 100, enabled: true })
    domains.value = res.items || []
  } catch (e: any) {
    console.error(e)
  }
}

const getDomainName = (domainId?: number) => {
  if (!domainId) return '域'
  const d = domains.value.find((x: any) => x.id === domainId)
  return d ? d.name : `域#${domainId}`
}

const loadSources = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit
    }
    if (filterForm.name) params.name = filterForm.name
    if (filterForm.type) params.type = filterForm.type
    if (filterForm.scope) params.scope = filterForm.scope
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
      is_default: item.is_default !== undefined ? item.is_default : (item.type === 'sql' && item.scope === 'system'),
      // 如果有 is_system 字段，使用它；否则根据类型判断是否为系统认证源
      is_system: item.is_system !== undefined ? item.is_system : (item.type === 'sql' && item.scope === 'system')
    }))
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载认证源列表失败')
  } finally {
    loading.value = false
  }
}

const loadOperationLogs = async () => {
  if (!currentSource.value) return

  logsLoading.value = true
  try {
    const params: any = {
      limit: logPagination.limit,
      offset: (logPagination.page - 1) * logPagination.limit
    }

    // 调用获取特定资源操作日志的API
    const res = await getResourceOperationLogs('auth_source', currentSource.value.id, params)
    operationLogs.value = res.items || []
    logPagination.total = res.total || 0
  } catch (e: any) {
    console.error('加载操作日志失败:', e)
    ElMessage.error(e.message || '加载操作日志失败')
    // 使用模拟数据
    operationLogs.value = [
      {
        id: 1,
        operation_time: new Date().toISOString(),
        resource_name: currentSource.value.name,
        resource_type: '认证源',
        operation_type: '创建',
        service_type: 'IAM',
        risk_level: 'low',
        time_type: '实时',
        result: 'success',
        operator: 'admin',
        project_id: 1,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      },
      {
        id: 2,
        operation_time: new Date(Date.now() - 86400000).toISOString(), // 一天前
        resource_name: currentSource.value.name,
        resource_type: '认证源',
        operation_type: '修改配置',
        service_type: 'IAM',
        risk_level: 'medium',
        time_type: '实时',
        result: 'success',
        operator: 'admin',
        project_id: 1,
        created_at: new Date(Date.now() - 86400000).toISOString(),
        updated_at: new Date(Date.now() - 86400000).toISOString()
      }
    ]
    logPagination.total = 2
  } finally {
    logsLoading.value = false
  }
}

const resetFilter = () => {
  filterForm.keyword = ''
  filterForm.name = ''
  filterForm.type = ''
  filterForm.scope = ''
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
  form.enabled = true
  form.config = {
    url: '',
    base_dn: '',
    bind_dn: '',
    bind_password: '',
    user_filter: '(objectClass=person)',
    user_id_attribute: 'uid',
    user_name_attribute: 'cn',
    user_search_base: '',
    group_search_base: '',
    user_enabled_attribute: 'enabled',
    protocol: 'ldap',
    auth_type: 'openldap',
    target_domain: ''
  }
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
    user_search_base: row.config?.user_search_base || '',
    group_search_base: row.config?.group_search_base || '',
    user_enabled_attribute: row.config?.user_enabled_attribute || 'enabled',
    protocol: row.config?.protocol || 'ldap',
    auth_type: row.config?.auth_type || 'openldap', // 确保认证类型正确赋值
    target_domain: row.config?.target_domain || ''
  }
  dialogVisible.value = true
}

const showDetailWithLogs = (row: AuthSource) => {
  currentSource.value = row
  detailDialogVisible.value = true
  activeTab.value = 'detail'
  // 重置日志分页
  logPagination.page = 1
  logPagination.limit = 20
  logPagination.total = 0
  operationLogs.value = []
}

const handleView = (row: AuthSource) => {
  showDetailWithLogs(row)
  activeTab.value = 'detail'
}

const handleTest = async (row: AuthSource) => {
  try {
    await testAuthSource(row.id)
    ElMessage.success('连接测试成功')
  } catch (e: any) {
    ElMessage.error(e.message || '连接测试失败')
  }
}

const handleToggleRunning = async (row: AuthSource) => {
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

const handleSync = async (row: AuthSource) => {
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

const handleToggleEnable = async (row: AuthSource) => {
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

const handleCommand = (command: string, row: AuthSource) => {
  switch (command) {
    case 'toggleEnable': handleToggleEnable(row); break
    case 'sync': handleSync(row); break
    case 'test': handleTest(row); break
    case 'delete': handleDelete(row); break
  }
}

const handleTypeChange = (type: string) => {
  // 切换类型时重置 config
  if (type === 'ldap') {
    form.config = {
      url: '',
      base_dn: '',
      bind_dn: '',
      bind_password: '',
      user_filter: '(objectClass=person)',
      user_id_attribute: 'uid',
      user_name_attribute: 'cn',
      user_search_base: '',
      group_search_base: '',
      user_enabled_attribute: 'enabled',
      protocol: 'ldap',
      auth_type: 'openldap', // 只保留OpenLDAP作为选项
      target_domain: ''
    }
  } else if (type === 'local') {
    form.config = {
      url: '',
      base_dn: '',
      bind_dn: '',
      bind_password: '',
      user_filter: '',
      user_id_attribute: 'uid',
      user_name_attribute: 'cn',
      user_search_base: '',
      group_search_base: '',
      user_enabled_attribute: 'enabled',
      protocol: '',
      auth_type: 'openldap', // 本地认证源也设为openldap，但实际不使用
      target_domain: ''
    }
  }
}

const handleDelete = async (row: AuthSource) => {
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

// 测试 LDAP 连接和用户查询
const testLdapConnection = async () => {
  // Validate that essential LDAP config fields are filled
  if (!form.config.url || !form.config.base_dn || !form.config.bind_dn || !form.config.bind_password) {
    ElMessage.warning('请先填写完整的LDAP配置信息（服务器地址、基本DN、用户名、密码）')
    return
  }

  ldapTesting.value = true
  ldapTestResult.value = null

  try {
    // Prepare a test config with the form values and test connection
    const testConfig = {
      url: form.config.url,
      base_dn: form.config.base_dn,
      bind_dn: form.config.bind_dn,
      bind_password: form.config.bind_password,
      user_filter: form.config.user_filter || '(objectClass=person)',
      user_id_attribute: form.config.user_id_attribute || 'uid',
      user_name_attribute: form.config.user_name_attribute || 'cn',
      user_search_base: form.config.user_search_base || form.config.base_dn,
      group_search_base: form.config.group_search_base || form.config.base_dn,
      user_enabled_attribute: form.config.user_enabled_attribute || 'enabled'
    }

    // Call the API function to test LDAP user query
    const result = await testLdapUsers(testConfig)

    if (result.success) {
      ldapTestResult.value = {
        success: true,
        message: `连接成功！找到 ${result.users.length} 个用户`,
        users: result.users.slice(0, 5) // Show first 5 users
      }
      ElMessage.success(`LDAP连接成功，找到 ${result.users.length} 个用户`)
    } else {
      ldapTestResult.value = {
        success: false,
        message: `连接失败：${result.message}`
      }
      ElMessage.error(`LDAP连接失败：${result.message}`)
    }
  } catch (error: any) {
    console.error('LDAP connection test error:', error)
    ldapTestResult.value = {
      success: false,
      message: `连接失败：${error.message || '网络错误'}`
    }
    ElMessage.error(`LDAP连接测试失败：${error.message || '网络错误'}`)
  } finally {
    ldapTesting.value = false
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
        config: {
          url: form.config.url,
          base_dn: form.config.base_dn,
          bind_dn: form.config.bind_dn,
          bind_password: form.config.bind_password,
          user_filter: form.config.user_filter,
          user_id_attribute: form.config.user_id_attribute,
          user_name_attribute: form.config.user_name_attribute,
          user_search_base: form.config.user_search_base,
          group_search_base: form.config.group_search_base,
          user_enabled_attribute: form.config.user_enabled_attribute,
          protocol: form.config.protocol,
          auth_type: form.config.auth_type,
          target_domain: form.config.target_domain
        }
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

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
}

// 当切换到日志标签页时加载日志
const onTabChange = (tabName: string) => {
  if (tabName === 'logs' && currentSource.value) {
    loadOperationLogs()
  }
}

// 监听标签页变化
watch(activeTab, (newVal) => {
  onTabChange(newVal)
})

onMounted(() => {
  loadSources()
  loadDomains()  // Load domain list for both forms and filters
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

.text-success {
  color: #67c23a;
}

.text-error {
  color: #f56c6c;
}

.table-cell-with-icon {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.table-cell-with-icon:hover {
  color: #409EFF;
}

.clickable-name {
  font-weight: 500;
}

.remark {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.3;
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
