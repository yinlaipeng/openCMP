<template>
  <div class="cloud-accounts-container">
    <div class="page-header">
      <h2>云账户管理</h2>
      <el-button type="primary" @click="showWizard = true">
        <el-icon><Plus /></el-icon>
        添加云账户
      </el-button>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="queryForm" @submit.prevent="loadAccounts">
        <el-form-item label="名称">
          <el-input v-model="queryForm.name" placeholder="云账户名称" clearable />
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="queryForm.provider_type" placeholder="选择平台" clearable>
            <el-option label="阿里云" value="alibaba" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryForm.status" placeholder="连接状态" clearable>
            <el-option label="已连接" value="active" />
            <el-option label="未连接" value="inactive" />
          </el-select>
        </el-form-item>
        <el-form-item label="启用状态">
          <el-select v-model="queryForm.enabled" placeholder="启用状态" clearable>
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadAccounts">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 空状态 -->
    <EmptyState
      v-if="!loading && accounts.length === 0"
      title="暂无云账户"
      description="当前没有任何云账户，点击下方按钮添加云账户"
      :icon="Cloudy"
      createButtonText="添加云账户"
      @create="showWizard = true"
    />

    <!-- 数据表格 -->
    <el-table
      :data="accounts"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
      v-if="accounts.length > 0 || loading"
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" min-width="150" />
      <el-table-column prop="status" label="连接状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="enabled" label="启用状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.enabled ? 'success' : 'info'">
            {{ row.enabled ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="health_status" label="健康状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getHealthStatusType(row.health_status)">
            {{ getHealthStatusText(row.health_status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="provider_type" label="平台" width="100">
        <template #default="{ row }">
          <el-tag :type="getProviderType(row.provider_type)">
            {{ getProviderName(row.provider_type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="account_number" label="账号" width="150" show-overflow-tooltip />
      <el-table-column prop="balance" label="余额" width="100">
        <template #default="{ row }">
          ¥{{ row.balance?.toFixed(2) || '0.00' }}
        </template>
      </el-table-column>
      <el-table-column prop="last_sync" label="上次同步" width="160" show-overflow-tooltip />
      <el-table-column label="同步时间" width="140">
        <template #default="{ row }">
          <div style="display: flex; align-items: center; gap: 5px;">
            <el-icon v-if="isSyncing(row.id)" class="is-loading" :size="16" color="#409EFF">
              <Loading />
            </el-icon>
            <span>{{ getSyncTimeDisplay(row) }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="domain_id" label="所属域" width="100">
        <template #default="{ row }">
          {{ getDomainName(row.domain_id) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="140" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="handleSyncClick(row)">同步</el-button>
          <el-dropdown trigger="click" style="margin-left: 8px" @command="(cmd: string) => handleDropdownCommand(cmd, row)">
            <el-button size="small">
              更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="enable" :disabled="row.enabled">启用</el-dropdown-item>
                <el-dropdown-item command="disable" :disabled="!row.enabled">禁用</el-dropdown-item>
                <el-dropdown-item command="test">连接测试</el-dropdown-item>
                <el-dropdown-item command="update">更新账号</el-dropdown-item>
                <el-dropdown-item divided command="delete">
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

    <!-- 向导式添加云账号对话框 -->
    <el-dialog
      v-model="showWizard"
      :title="'添加云账户'"
      width="80%"
      :fullscreen="isMobile"
    >
      <el-steps :active="wizardStep" finish-status="success" align-center>
        <el-step title="选择云平台" />
        <el-step title="配置基本信息" />
        <el-step title="配置同步区域" />
        <el-step title="定时同步任务" />
      </el-steps>

      <!-- 步骤1: 选择云平台 -->
      <div v-if="wizardStep === 0" class="wizard-step-content">
        <h3>请选择您要添加的云平台</h3>
        <div class="provider-grid">
          <el-card
            v-for="provider in providers"
            :key="provider.id"
            class="provider-card"
            @click="selectProvider(provider.id)"
          >
            <div class="provider-item">
              <el-icon :size="40"><component :is="provider.icon" /></el-icon>
              <h4>{{ provider.displayName }}</h4>
            </div>
          </el-card>
        </div>
      </div>

      <!-- 步骤2: 配置基本信息 -->
      <div v-if="wizardStep === 1" class="wizard-step-content">
        <el-form :model="wizardForm" :rules="wizardRules" ref="wizardFormRef" label-width="150px">
          <el-form-item label="名称" prop="name">
            <el-input v-model="wizardForm.name" placeholder="请输入云账户名称" />
          </el-form-item>

          <el-form-item label="备注" prop="remarks">
            <el-input v-model="wizardForm.remarks" type="textarea" placeholder="请输入备注" />
          </el-form-item>

          <el-form-item label="密钥ID" prop="accessKeyId">
            <el-input v-model="wizardForm.accessKeyId" placeholder="请输入Access Key ID" />
          </el-form-item>

          <el-form-item label="密钥Secret" prop="accessKeySecret">
            <el-input v-model="wizardForm.accessKeySecret" type="password" placeholder="请输入Access Key Secret" />
          </el-form-item>

          <el-form-item label="测试连接">
            <el-button @click="testConnectionInWizard" :loading="wizardForm.testConnectionStatus === 'testing'">
              {{ wizardForm.testConnectionStatus === 'testing' ? '测试中...' : '测试连接' }}
            </el-button>
            <span v-if="wizardForm.testConnectionResult" :class="wizardForm.testConnectionResult.includes('成功') ? 'text-success' : 'text-danger'">
              {{ wizardForm.testConnectionResult }}
            </span>
          </el-form-item>
        </el-form>
      </div>

      <!-- 步骤3: 配置同步区域 -->
      <div v-if="wizardStep === 2" class="wizard-step-content">
        <p class="region-tip">
          请选择要同步资源的区域，默认同步所有区域
        </p>

        <div class="region-actions">
          <el-button @click="selectAllRegions">全选</el-button>
          <el-button @click="clearAllRegions">清空</el-button>
        </div>

        <el-table
          :data="availableRegions"
          style="width: 100%; margin-top: 20px;"
          @selection-change="handleRegionSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="name" label="名称" />
          <el-table-column prop="status" label="状态">
            <template #default="{ row }">
              <el-tag :type="getRegionStatusType(row.status)">
                {{ getRegionStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 步骤4: 定时同步任务设置 -->
      <div v-if="wizardStep === 3" class="wizard-step-content">
        <h3>定时同步任务设置（可选）</h3>
        <el-form :model="scheduleForm" label-width="120px">
          <el-form-item label="名称">
            <el-input v-model="scheduleForm.name" placeholder="请输入任务名称" />
          </el-form-item>

          <el-form-item label="触发频次">
            <el-select v-model="scheduleForm.frequency">
              <el-option label="每天" value="daily" />
              <el-option label="每周" value="weekly" />
              <el-option label="每月" value="monthly" />
            </el-select>
          </el-form-item>

          <el-form-item label="触发时间">
            <el-time-picker
              v-model="scheduleForm.triggerTime"
              placeholder="选择时间"
              format="HH:mm"
              value-format="HH:mm"
            />
          </el-form-item>
        </el-form>
      </div>

      <template #footer>
        <div class="wizard-footer">
          <el-button @click="previousStep" :disabled="wizardStep === 0">上一步</el-button>
          <el-button
            v-if="wizardStep < 3"
            type="primary"
            @click="nextStep"
          >
            下一步
          </el-button>
          <el-button
            v-else
            type="primary"
            @click="submitWizard"
            :loading="submitting"
          >
            提交
          </el-button>
          <el-button @click="showWizard = false">取消</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 同步云账号对话框 -->
    <el-dialog v-model="showSyncDialog" title="同步云账号" width="700px">
      <el-form :model="syncForm" label-width="120px">
        <el-form-item label="云账号名称">
          <el-input v-model="syncForm.name" readonly />
        </el-form-item>
        <el-form-item label="环境">
          <el-input v-model="syncForm.environment" readonly />
        </el-form-item>
        <el-form-item label="连接状态">
          <el-tag :type="syncForm.connectionStatus ? 'success' : 'danger'">
            {{ syncForm.connectionStatus ? '已连接' : '未连接' }}
          </el-tag>
        </el-form-item>
        <el-form-item label="同步模式">
          <el-radio-group v-model="syncForm.syncMode">
            <el-radio value="full">全量同步</el-radio>
            <el-radio value="incremental">增量同步</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="同步资源类型">
          <el-radio-group v-model="syncForm.syncAll" style="margin-bottom: 10px;">
            <el-radio value="all">全部资源类型</el-radio>
            <el-radio value="specific">指定资源类型</el-radio>
          </el-radio-group>
          <div v-if="syncForm.syncAll === 'specific'" style="border: 1px solid #dcdfe6; padding: 10px; border-radius: 4px; max-height: 300px; overflow-y: auto;">
            <el-checkbox-group v-model="syncForm.syncResourceTypes">
              <el-row>
                <el-col :span="8" v-for="rt in supportedResourceTypes" :key="rt.id">
                  <el-checkbox :value="rt.id" :label="rt.id">{{ rt.name }}</el-checkbox>
                </el-col>
              </el-row>
            </el-checkbox-group>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="confirmSync" type="primary" :loading="syncing">确认同步</el-button>
        <el-button @click="showSyncDialog = false">取消</el-button>
      </template>
    </el-dialog>

    <!-- 更新账号对话框 -->
    <EditAccountDialog
      v-model="showUpdateDialog"
      :account-id="editAccountId"
      @saved="loadAccounts"
    />

    <!-- 属性设置弹窗 -->
    <CloudAccountDetailDialog
      v-model="showAccountDetailDialog"
      :account-id="accountDetailId"
      :initial-tab="accountDetailTab"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus, Cloudy, ArrowDown, Loading } from '@element-plus/icons-vue'
import { getCloudAccounts, createCloudAccount, updateCloudAccount, deleteCloudAccount, syncCloudAccount, testConnection as testConnectionAPI, updateCloudAccountStatus, updateCloudAccountAttributes, getSupportedResourceTypes, type ResourceType } from '@/api/cloud-account'
import { getSyncPolicies } from '@/api/sync-policy'
import { getProjects } from '@/api/project'
import { getDomains } from '@/api/iam'
import { getLatestSyncLog } from '@/api/sync-log'
import type { CloudAccount, CreateCloudAccountRequest, Project } from '@/types'
import type { SyncPolicy } from '@/types/sync-policy'
import EmptyState from '@/components/common/EmptyState.vue'
import EditAccountDialog from './components/EditAccountDialog.vue'
import CloudAccountDetailDialog from './components/CloudAccountDetailDialog.vue'

const accounts = ref<CloudAccount[]>([])
const syncPolicies = ref<SyncPolicy[]>([])
const projects = ref<Project[]>([])
const domains = ref<{ id: number; name: string }[]>([])
const loading = ref(false)
const syncStatusMap = ref<Map<number, { status: string; startTime: Date; duration?: number }>>(new Map())
const showWizard = ref(false)
const wizardStep = ref(0)
const submitting = ref(false)
const wizardFormRef = ref<FormInstance>()

// 分页数据
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 查询表单
const queryForm = reactive({
  name: '',
  provider_type: '',
  status: '',
  enabled: undefined as boolean | undefined
})

// 供应商选项
const providers = ref([
  { id: 'alibaba', name: 'alibaba', displayName: '阿里云', icon: Cloudy },
  { id: 'tencent', name: 'tencent', displayName: '腾讯云', icon: Cloudy },
  { id: 'aws', name: 'aws', displayName: 'AWS', icon: Cloudy },
  { id: 'azure', name: 'azure', displayName: 'Azure', icon: Cloudy }
])

// 向导表单数据
const wizardForm = reactive({
  name: '',
  remarks: '',
  accessKeyId: '',
  accessKeySecret: '',
  testConnectionStatus: '',
  testConnectionResult: '',
  selectedRegions: [] as string[]
})

// 定时任务表单
const scheduleForm = reactive({
  name: '',
  frequency: 'daily',
  triggerTime: '02:00'
})

// 可用区域列表
const availableRegions = ref([
  { name: '华北1（青岛）', status: 'available' },
  { name: '华北2（北京）', status: 'available' },
  { name: '华东1（杭州）', status: 'available' },
  { name: '华东2（上海）', status: 'available' },
  { name: '华南1（深圳）', status: 'available' },
  { name: '西南1（成都）', status: 'available' },
  { name: '中国（香港）', status: 'available' },
  { name: '新加坡', status: 'available' }
])

// 验证规则
const wizardRules = {
  name: [{ required: true, message: '请输入云账户名称', trigger: 'blur' }],
  accessKeyId: [{ required: true, message: '请输入密钥ID', trigger: 'blur' }],
  accessKeySecret: [{ required: true, message: '请输入密钥Secret', trigger: 'blur' }]
}

// 检测移动端
const isMobile = computed(() => {
  return window.innerWidth < 768
})

// 同步相关
const showSyncDialog = ref(false)
const showUpdateDialog = ref(false)
const showAccountDetailDialog = ref(false)
const editAccountId = ref<number | null>(null)
const accountDetailId = ref<number | null>(null)
const accountDetailTab = ref('detail')
const syncing = ref(false)
const supportedResourceTypes = ref<ResourceType[]>([])

const syncForm = ref({
  name: '',
  environment: '',
  connectionStatus: false,
  syncMode: 'full',
  syncAll: 'all',
  syncResourceTypes: [] as string[]
})

const currentAccount = ref<any>(null)

// 加载云账户列表
const loadAccounts = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (queryForm.name) params.name = queryForm.name
    if (queryForm.provider_type) params.provider_type = queryForm.provider_type
    if (queryForm.status) params.status = queryForm.status
    if (queryForm.enabled !== undefined) params.enabled = queryForm.enabled

    const res = await getCloudAccounts(params)
    accounts.value = res.items || []
    total.value = res.total || 0
  } catch (e) {
    console.error(e)
    ElMessage.error('加载云账户失败')
  } finally {
    loading.value = false
  }
}

const loadSyncPolicies = async () => {
  try {
    const res = await getSyncPolicies()
    syncPolicies.value = res.items || []
  } catch (e) {
    console.error('Failed to load sync policies:', e)
  }
}

const loadProjects = async () => {
  try {
    const res = await getProjects()
    projects.value = res.items || []
  } catch (e) {
    console.error('Failed to load projects:', e)
  }
}

const loadDomains = async () => {
  try {
    const res = await getDomains()
    domains.value = (res.items || []).map(d => ({ id: d.id, name: d.name }))
  } catch (e) {
    console.error('Failed to load domains:', e)
  }
}

const loadResourceTypes = async () => {
  try {
    const res = await getSupportedResourceTypes()
    supportedResourceTypes.value = res.items || []
  } catch (e) {
    console.error('Failed to load resource types:', e)
  }
}

const resetQuery = () => {
  queryForm.name = ''
  queryForm.provider_type = ''
  queryForm.status = ''
  queryForm.enabled = undefined
  currentPage.value = 1
  loadAccounts()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadAccounts()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadAccounts()
}

const getProviderName = (type: string) => {
  const map: Record<string, string> = {
    alibaba: '阿里云',
    tencent: '腾讯云',
    aws: 'AWS',
    azure: 'Azure'
  }
  return map[type] || type
}

const getProviderType = (type: string) => {
  const map: Record<string, any> = {
    alibaba: 'warning',
    tencent: 'primary',
    aws: 'danger',
    azure: 'success'
  }
  return map[type] || 'info'
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    active: 'success',
    connected: 'success',
    inactive: 'info',
    disconnected: 'info',
    error: 'danger',
    pending: 'warning'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '已连接',
    connected: '已连接',
    inactive: '未连接',
    disconnected: '未连接',
    error: '连接错误',
    pending: '连接中'
  }
  return statusMap[status] || status
}

const getHealthStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    healthy: 'success',
    normal: 'success',
    unhealthy: 'danger',
    abnormal: 'danger',
    warning: 'warning'
  }
  return statusMap[status] || 'info'
}

const getHealthStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    healthy: '正常',
    normal: '正常',
    unhealthy: '异常',
    abnormal: '异常',
    warning: '警告'
  }
  return statusMap[status] || status
}

const getDomainName = (domainId: number) => {
  const domain = domains.value.find(d => d.id === domainId)
  return domain ? domain.name : `域${domainId}`
}

const isSyncing = (accountId: number) => {
  const status = syncStatusMap.value.get(accountId)
  return status?.status === 'running'
}

const getSyncTimeDisplay = (row: CloudAccount) => {
  const status = syncStatusMap.value.get(row.id)
  if (status?.status === 'running') {
    const elapsed = Math.floor((Date.now() - status.startTime.getTime()) / 1000)
    return `已同步 ${elapsed}s`
  }
  if (status?.duration) {
    return `耗时 ${status.duration}s`
  }
  return row.last_sync ? formatDate(row.last_sync) : '-'
}

const formatDate = (dateString: string | undefined) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

const getRegionStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    available: 'success',
    unavailable: 'danger'
  }
  return statusMap[status] || 'info'
}

const getRegionStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    available: '可用',
    unavailable: '不可用'
  }
  return statusMap[status] || status
}

const selectProvider = (providerId: string) => {
  ElMessage.success(`已选择 ${providers.value.find(p => p.id === providerId)?.displayName}`)
  wizardStep.value = 1
}

const nextStep = async () => {
  if (wizardStep.value === 1) {
    if (!wizardFormRef.value) return
    const valid = await wizardFormRef.value.validate().catch(() => false)
    if (!valid) return
  }
  if (wizardStep.value < 3) {
    wizardStep.value++
  }
}

const previousStep = () => {
  if (wizardStep.value > 0) {
    wizardStep.value--
  }
}

const selectAllRegions = () => {
  wizardForm.selectedRegions = availableRegions.value
    .filter(region => region.status === 'available')
    .map(region => region.name)
}

const clearAllRegions = () => {
  wizardForm.selectedRegions = []
}

const handleRegionSelectionChange = (selection: any[]) => {
  wizardForm.selectedRegions = selection.map(item => item.name)
}

const testConnectionInWizard = async () => {
  if (!wizardForm.accessKeyId || !wizardForm.accessKeySecret) {
    ElMessage.warning('请先填写密钥ID和密钥Secret')
    return
  }

  wizardForm.testConnectionStatus = 'testing'
  wizardForm.testConnectionResult = '正在测试连接...'

  try {
    const tempAccountData: CreateCloudAccountRequest = {
      name: wizardForm.name || 'temp-test',
      provider_type: 'alibaba',
      credentials: {
        access_key_id: wizardForm.accessKeyId,
        access_key_secret: wizardForm.accessKeySecret
      },
      description: 'Temporary account for connection testing',
      remarks: 'Temporary account for connection testing',
      enabled: true,
      health_status: 'healthy',
      domain_id: 1,
      resource_assignment_method: 'tag_mapping'
    }

    const tempAccount = await createCloudAccount(tempAccountData)
    const response = await testConnectionAPI(tempAccount.id)
    await deleteCloudAccount(tempAccount.id)

    wizardForm.testConnectionStatus = response.connected ? 'success' : 'error'
    wizardForm.testConnectionResult = response.message || (response.connected ? '连接成功' : '连接失败')
  } catch (error: any) {
    wizardForm.testConnectionStatus = 'error'
    wizardForm.testConnectionResult = error.message || '连接测试失败'
  }
}

const submitWizard = async () => {
  submitting.value = true
  try {
    const formData: CreateCloudAccountRequest = {
      name: wizardForm.name,
      provider_type: 'alibaba',
      credentials: {
        access_key_id: wizardForm.accessKeyId,
        access_key_secret: wizardForm.accessKeySecret
      },
      description: wizardForm.remarks,
      remarks: wizardForm.remarks,
      enabled: true,
      health_status: 'healthy',
      account_number: wizardForm.accessKeyId.substring(0, 10) + '...',
      domain_id: 1,
      resource_assignment_method: 'tag_mapping'
    }

    await createCloudAccount(formData)
    ElMessage.success('云账户添加成功')

    resetWizard()
    showWizard.value = false
    loadAccounts()
  } catch (error) {
    console.error(error)
    ElMessage.error('添加云账户失败')
  } finally {
    submitting.value = false
  }
}

const resetWizard = () => {
  wizardStep.value = 0
  Object.assign(wizardForm, {
    name: '',
    remarks: '',
    accessKeyId: '',
    accessKeySecret: '',
    testConnectionStatus: '',
    testConnectionResult: '',
    selectedRegions: []
  })
  Object.assign(scheduleForm, {
    name: '',
    frequency: 'daily',
    triggerTime: '02:00'
  })
}

const handleSyncClick = (row: any) => {
  currentAccount.value = row
  syncForm.value.name = row.name
  syncForm.value.environment = getProviderName(row.provider_type)
  syncForm.value.connectionStatus = row.status === 'active'
  showSyncDialog.value = true
}

const handleDropdownCommand = async (command: string, row: any) => {
  currentAccount.value = row

  switch (command) {
    case 'enable':
      try {
        await updateCloudAccountStatus(row.id, true)
        ElMessage.success('已启用云账号')
        loadAccounts()
      } catch (error) {
        ElMessage.error('启用失败')
      }
      break
    case 'disable':
      try {
        await updateCloudAccountStatus(row.id, false)
        ElMessage.success('已禁用云账号')
        loadAccounts()
      } catch (error) {
        ElMessage.error('禁用失败')
      }
      break
    case 'test':
      try {
        const response = await testConnectionAPI(row.id)
        if (response.connected) {
          ElMessage.success('连接测试成功')
        } else {
          ElMessage.warning('连接测试失败: ' + (response.message || '无法连接'))
        }
      } catch (error) {
        ElMessage.error('连接测试失败')
      }
      break
    case 'update':
      editAccountId.value = row.id
      showUpdateDialog.value = true
      break
    case 'delete':
      try {
        await ElMessageBox.confirm(`确定要删除云账户 "${row.name}" 吗？`, '提示', { type: 'warning' })
        await deleteCloudAccount(row.id)
        ElMessage.success('删除成功')
        loadAccounts()
      } catch (e: any) {
        if (e !== 'cancel') {
          ElMessage.error('删除失败')
        }
      }
      break
  }
}

const confirmSync = async () => {
  syncing.value = true
  try {
    const resourceTypes = syncForm.value.syncAll === 'all' ? ['all'] : syncForm.value.syncResourceTypes

    await syncCloudAccount(currentAccount.value.id, {
      mode: syncForm.value.syncMode,
      resource_types: resourceTypes
    })

    ElMessage.success(`已启动对 ${syncForm.value.name} 的同步`)
    showSyncDialog.value = false

    pollSyncStatus(currentAccount.value.id)
  } catch (error) {
    ElMessage.error('同步启动失败')
  } finally {
    syncing.value = false
  }
}

const pollSyncStatus = async (accountId: number) => {
  syncStatusMap.value.set(accountId, {
    status: 'running',
    startTime: new Date()
  })

  const pollInterval = setInterval(async () => {
    try {
      const log = await getLatestSyncLog(accountId)
      if (log.status !== 'running') {
        clearInterval(pollInterval)
        syncStatusMap.value.set(accountId, {
          status: log.status,
          startTime: new Date(log.sync_start_time),
          duration: log.sync_duration
        })
        loadAccounts()
      }
    } catch (e) {
      console.error('Failed to poll sync status:', e)
    }
  }, 3000)

  setTimeout(() => {
    clearInterval(pollInterval)
  }, 5 * 60 * 1000)
}

onMounted(() => {
  loadAccounts()
  loadSyncPolicies()
  loadProjects()
  loadDomains()
  loadResourceTypes()
})
</script>

<style scoped>
.cloud-accounts-container {
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

.filter-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}

.wizard-step-content {
  margin: 30px 0;
}

.provider-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.provider-card {
  cursor: pointer;
  transition: all 0.3s;
}

.provider-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
}

.provider-item {
  text-align: center;
  padding: 20px 10px;
}

.provider-item h4 {
  margin-top: 10px;
}

.region-tip {
  background-color: #ecf5ff;
  padding: 12px;
  border-radius: 4px;
  border-left: 4px solid #409eff;
  margin-bottom: 20px;
}

.region-actions {
  margin-bottom: 15px;
}

.wizard-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.text-success {
  color: #67c23a;
}

.text-danger {
  color: #f56c6c;
}
</style>