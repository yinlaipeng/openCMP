<template>
  <div class="eips-container">
    <!-- Tabs 过滤 -->
    <el-tabs v-model="activeTab" class="filter-tabs" @tab-change="handleTabChange">
      <el-tab-pane label="全部" name="all" />
      <el-tab-pane label="私有云" name="on-premise" />
      <el-tab-pane label="公有云" name="public" />
    </el-tabs>

    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <el-button type="primary" @click="handleCreate">新建</el-button>
        <el-dropdown trigger="click" @command="handleBatchCommand">
          <el-button>
            批量操作 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="bind" :disabled="selectedRows.length === 0">批量绑定</el-dropdown-item>
              <el-dropdown-item command="unbind" :disabled="selectedRows.length === 0">批量解绑</el-dropdown-item>
              <el-dropdown-item command="delete" :disabled="selectedRows.length === 0" divided>批量删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button @click="handleTagManage">标签</el-button>
      </div>
      <div class="toolbar-right">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索名称/IP地址"
          clearable
          style="width: 200px"
          @keyup.enter="loadEips"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
    </div>

    <!-- 过滤栏 -->
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="filters" @submit.prevent="loadEips">
        <el-form-item label="云账号">
          <CloudAccountSelector
            v-model:value="filters.cloud_account_id"
            placeholder="全部云账号"
            @change="handleAccountChange"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 120px">
            <el-option label="已绑定" value="In-Use" />
            <el-option label="未绑定" value="Available" />
            <el-option label="绑定中" value="Associating" />
            <el-option label="解绑中" value="Unassociating" />
          </el-select>
        </el-form-item>
        <el-form-item label="计费方式">
          <el-select v-model="filters.billing_method" placeholder="全部" clearable style="width: 120px">
            <el-option label="按流量" value="traffic" />
            <el-option label="按带宽" value="bandwidth" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-select v-model="filters.region" placeholder="全部区域" clearable style="width: 150px">
            <el-option v-for="r in regions" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadEips">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格 -->
    <el-table
      :data="eips"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="50" />
      <el-table-column label="名称" min-width="150">
        <template #default="{ row }">
          <el-link type="primary" @click="viewDetail(row)">{{ row.name || row.eip_id }}</el-link>
        </template>
      </el-table-column>
      <el-table-column label="平台/云账号" width="180">
        <template #default="{ row }">
          <div class="platform-cell">
            <el-tag size="small" :type="getPlatformType(row.provider_type)" effect="plain">
              {{ getPlatformLabel(row.provider_type) }}
            </el-tag>
            <span class="account-name">{{ row.account_name || '-' }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="计费方式" width="100">
        <template #default="{ row }">
          <span>{{ getBillingLabel(row.billing_method) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">
            {{ getStatusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="address" label="IP地址" width="140" />
      <el-table-column label="带宽" width="100">
        <template #default="{ row }">
          {{ row.bandwidth }} Mbps
        </template>
      </el-table-column>
      <el-table-column label="关联资源" min-width="150">
        <template #default="{ row }">
          <template v-if="row.resource_name">
            <el-link type="primary">{{ row.resource_name }}</el-link>
            <span class="resource-type">({{ getResourceTypeLabel(row.resource_type) }})</span>
          </template>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="region" label="区域" width="120" />
      <el-table-column label="标签" width="120">
        <template #default="{ row }">
          <template v-if="row.tags && row.tags.length > 0">
            <el-tag v-for="tag in row.tags.slice(0, 2)" :key="tag.key" size="small" class="tag-item">
              {{ tag.key }}: {{ tag.value }}
            </el-tag>
            <el-tag v-if="row.tags.length > 2" size="small" type="info">
              +{{ row.tags.length - 2 }}
            </el-tag>
          </template>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="项目" width="120">
        <template #default="{ row }">
          <span>{{ row.project_name || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button size="small" link type="primary" @click="handleBind(row)" v-if="row.status === 'Available'">绑定</el-button>
          <el-button size="small" link type="primary" @click="handleUnbind(row)" v-if="row.status === 'In-Use'">解绑</el-button>
          <el-button size="small" link type="primary" @click="handleEdit(row)">修改</el-button>
          <el-button size="small" link type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      class="pagination"
    />

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="弹性公网IP详情" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedEip?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedEip?.name }}</el-descriptions-item>
        <el-descriptions-item label="EIP ID">{{ selectedEip?.eip_id }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(selectedEip?.status)">
            {{ getStatusLabel(selectedEip?.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ selectedEip?.address }}</el-descriptions-item>
        <el-descriptions-item label="带宽">{{ selectedEip?.bandwidth }} Mbps</el-descriptions-item>
        <el-descriptions-item label="计费方式">{{ getBillingLabel(selectedEip?.billing_method) }}</el-descriptions-item>
        <el-descriptions-item label="关联资源">{{ selectedEip?.resource_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="资源类型">{{ getResourceTypeLabel(selectedEip?.resource_type) }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedEip?.region }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedEip?.account_name }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedEip?.project_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="标签" :span="2">
          <template v-if="selectedEip?.tags && selectedEip.tags.length > 0">
            <el-tag v-for="tag in selectedEip.tags" :key="tag.key" class="tag-item">
              {{ tag.key }}: {{ tag.value }}
            </el-tag>
          </template>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ selectedEip?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 绑定对话框 -->
    <el-dialog v-model="bindVisible" title="绑定弹性IP" width="500px">
      <el-form :model="bindForm" label-width="100px">
        <el-form-item label="资源类型">
          <el-select v-model="bindForm.resource_type" placeholder="选择资源类型" style="width: 100%">
            <el-option label="云主机" value="guest" />
            <el-option label="负载均衡" value="loadbalancer" />
            <el-option label="NAT网关" value="natgateway" />
          </el-select>
        </el-form-item>
        <el-form-item label="选择资源">
          <el-select v-model="bindForm.resource_id" placeholder="选择要绑定的资源" style="width: 100%" filterable>
            <el-option v-for="res in availableResources" :key="res.id" :label="res.name" :value="res.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="bindVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmBind">确定</el-button>
      </template>
    </el-dialog>

    <!-- 新建对话框 -->
    <el-dialog v-model="createVisible" title="新建弹性公网IP" width="500px">
      <el-form :model="createForm" label-width="120px" :rules="createRules" ref="createFormRef">
        <el-form-item label="项目" prop="project_id">
          <el-select v-model="createForm.project_id" placeholder="选择项目" style="width: 100%">
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="云账号" prop="cloud_account_id">
          <CloudAccountSelector
            v-model:value="createForm.cloud_account_id"
            placeholder="选择云账号"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="计费方式" prop="billing_method">
          <el-radio-group v-model="createForm.billing_method">
            <el-radio-button value="traffic">按流量</el-radio-button>
            <el-radio-button value="bandwidth">按带宽</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="带宽峰值" prop="bandwidth">
          <el-input-number v-model="createForm.bandwidth" :min="1" :max="1000" style="width: 200px" />
          <span class="unit-label">Mbps</span>
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="createForm.name" placeholder="可选，输入EIP名称" />
        </el-form-item>
        <el-form-item label="标签">
          <div class="tags-editor">
            <div v-for="(tag, idx) in createForm.tags" :key="idx" class="tag-row">
              <el-input v-model="tag.key" placeholder="键" style="width: 120px" />
              <el-input v-model="tag.value" placeholder="值" style="width: 120px" />
              <el-button size="small" type="danger" @click="removeTag(idx)" link>删除</el-button>
            </div>
            <el-button size="small" @click="addTag">添加标签</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCreate" :loading="createLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 修改带宽对话框 -->
    <el-dialog v-model="editVisible" title="修改弹性IP带宽" width="400px">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="当前带宽">
          <span>{{ selectedEip?.bandwidth }} Mbps</span>
        </el-form-item>
        <el-form-item label="新带宽">
          <el-input-number v-model="editForm.bandwidth" :min="1" :max="1000" style="width: 200px" />
          <span class="unit-label">Mbps</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmEdit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 标签管理对话框 -->
    <el-dialog v-model="tagManageVisible" title="标签管理" width="600px">
      <el-table :data="selectedRows" style="width: 100%" max-height="300">
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column prop="address" label="IP地址" width="120" />
        <el-table-column label="当前标签">
          <template #default="{ row }">
            <el-tag v-for="tag in row.tags" :key="tag.key" size="small" class="tag-item">
              {{ tag.key }}: {{ tag.value }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      <div style="margin-top: 16px">
        <el-form :inline="true">
          <el-form-item label="标签键">
            <el-input v-model="tagForm.key" placeholder="输入标签键" style="width: 150px" />
          </el-form-item>
          <el-form-item label="标签值">
            <el-input v-model="tagForm.value" placeholder="输入标签值" style="width: 150px" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleAddTagToSelected">添加</el-button>
            <el-button @click="handleRemoveTagFromSelected">移除</el-button>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="tagManageVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown, Search } from '@element-plus/icons-vue'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import { getEIPs, getEIP, createEIP, bindEIP, unbindEIP, deleteEIP, batchDeleteEIPs } from '@/api/networkSync'
import { getProjects } from '@/api/iam'
import { getCloudAccounts } from '@/api/cloud'

interface EIP {
  id: number | string
  cloud_account_id: number
  eip_id: string
  name: string
  address: string
  bandwidth: number
  billing_method: string
  resource_id: string
  resource_type: string
  resource_name: string
  region_id: string
  region: string
  status: string
  provider_type: string
  account_name: string
  project_id: number
  project_name: string
  tags: { key: string; value: string }[]
  created_at: string
  updated_at: string
}

interface Project {
  id: number
  name: string
}

interface Resource {
  id: string
  name: string
  type: string
}

const eips = ref<EIP[]>([])
const loading = ref(false)
const detailVisible = ref(false)
const bindVisible = ref(false)
const createVisible = ref(false)
const editVisible = ref(false)
const tagManageVisible = ref(false)
const selectedEip = ref<EIP | null>(null)
const selectedRows = ref<EIP[]>([])
const createLoading = ref(false)

// Tabs
const activeTab = ref('all')

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 搜索
const searchKeyword = ref('')

// 过滤条件
const filters = reactive({
  cloud_account_id: null as number | null,
  status: '',
  billing_method: '',
  region: ''
})

// 项目列表
const projects = ref<Project[]>([])
const regions = ref<{ id: string; name: string }[]>([])
const availableResources = ref<Resource[]>([])

// 绑定表单
const bindForm = reactive({
  eip_id: '',
  resource_type: 'guest',
  resource_id: ''
})

// 创建表单
const createForm = reactive({
  project_id: null as number | null,
  cloud_account_id: null as number | null,
  billing_method: 'bandwidth',
  bandwidth: 100,
  name: '',
  tags: [] as { key: string; value: string }[]
})

const createRules = {
  project_id: [{ required: true, message: '请选择项目', trigger: 'change' }],
  cloud_account_id: [{ required: true, message: '请选择云账号', trigger: 'change' }],
  billing_method: [{ required: true, message: '请选择计费方式', trigger: 'change' }],
  bandwidth: [{ required: true, message: '请设置带宽', trigger: 'blur' }]
}

const createFormRef = ref()

// 修改表单
const editForm = reactive({
  bandwidth: 100
})

// 标签管理表单
const tagForm = reactive({
  key: '',
  value: ''
})

// 平台类型映射
const platformLabels: Record<string, string> = {
  alibaba: '阿里云',
  tencent: '腾讯云',
  aws: 'AWS',
  azure: 'Azure',
  huawei: '华为云',
  google: 'GCP',
  onpremise: '私有云'
}

const platformTypes: Record<string, '' | 'success' | 'warning' | 'info' | 'danger'> = {
  alibaba: 'warning',
  tencent: 'success',
  aws: 'info',
  azure: '',
  huawei: 'danger',
  google: 'success',
  onpremise: 'info'
}

const getPlatformLabel = (platform: string): string => {
  return platformLabels[platform] || platform || '未知'
}

const getPlatformType = (platform: string): '' | 'success' | 'warning' | 'info' | 'danger' => {
  return platformTypes[platform] || 'info'
}

// 计费方式映射
const billingLabels: Record<string, string> = {
  traffic: '按流量',
  bandwidth: '按带宽',
  PayByTraffic: '按流量',
  PayByBandwidth: '按带宽'
}

const getBillingLabel = (method: string): string => {
  return billingLabels[method] || method || '-'
}

// 状态映射
const statusLabels: Record<string, string> = {
  Available: '未绑定',
  'In-Use': '已绑定',
  Associating: '绑定中',
  Unassociating: '解绑中',
  available: '未绑定',
  'in-use': '已绑定'
}

const getStatusLabel = (status: string): string => {
  return statusLabels[status] || status
}

const getStatusType = (status: string): '' | 'success' | 'warning' | 'info' | 'danger' => {
  if (status === 'Available' || status === 'available') return 'success'
  if (status === 'In-Use' || status === 'in-use') return 'primary'
  if (status === 'Associating' || status === 'Unassociating') return 'warning'
  if (status === 'Error' || status === 'error') return 'danger'
  return 'info'
}

// 资源类型映射
const resourceTypeLabels: Record<string, string> = {
  guest: '云主机',
  loadbalancer: '负载均衡',
  natgateway: 'NAT网关',
  server: '云主机',
  ecs: 'ECS实例',
  slb: 'SLB实例'
}

const getResourceTypeLabel = (type: string): string => {
  return resourceTypeLabels[type] || type || '-'
}

const handleTabChange = (tab: string) => {
  activeTab.value = tab
  loadEips()
}

const handleAccountChange = (accountId: number | null) => {
  filters.cloud_account_id = accountId
}

const handleSelectionChange = (rows: EIP[]) => {
  selectedRows.value = rows
}

const handleBatchCommand = (command: string) => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请先选择要操作的EIP')
    return
  }
  switch (command) {
    case 'bind':
      ElMessage.info('批量绑定功能开发中')
      break
    case 'unbind':
      handleBatchUnbind()
      break
    case 'delete':
      handleBatchDelete()
      break
  }
}

const handleBatchUnbind = async () => {
  try {
    await ElMessageBox.confirm(`确认批量解绑 ${selectedRows.value.length} 个EIP？`, '提示', { type: 'warning' })
    for (const eip of selectedRows.value) {
      if (eip.status === 'In-Use') {
        await unbindEIP(eip.id)
      }
    }
    ElMessage.success('批量解绑成功')
    loadEips()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('批量解绑失败')
  }
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(`确认批量删除 ${selectedRows.value.length} 个EIP？此操作不可恢复`, '警告', { type: 'warning' })
    const ids = selectedRows.value.map(e => Number(e.id))
    await batchDeleteEIPs(ids)
    ElMessage.success('批量删除成功')
    loadEips()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('批量删除失败')
  }
}

const handleTagManage = () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请先选择要管理标签的EIP')
    return
  }
  tagManageVisible.value = true
}

const loadEips = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (filters.cloud_account_id) params.cloud_account_id = filters.cloud_account_id
    if (filters.status) params.status = filters.status
    if (filters.billing_method) params.billing_method = filters.billing_method
    if (filters.region) params.region = filters.region
    if (searchKeyword.value) params.name = searchKeyword.value
    if (activeTab.value !== 'all') params.platform = activeTab.value

    const res = await getEIPs(params)
    eips.value = res.items || []
    pagination.total = res.total || 0
  } catch (e) {
    console.error(e)
    ElMessage.error('加载弹性IP列表失败')
  } finally {
    loading.value = false
  }
}

const loadProjects = async () => {
  try {
    const res = await getProjects({ page: 1, page_size: 100 })
    projects.value = res.items || []
  } catch (e) {
    console.error('加载项目列表失败', e)
  }
}

const loadRegions = async () => {
  // 从已加载的 EIP 数据中提取区域列表
  const regionMap = new Map<string, string>()
  eips.value.forEach(eip => {
    if (eip.region_id && eip.region) {
      regionMap.set(eip.region_id, eip.region)
    }
  })
  regions.value = Array.from(regionMap.entries()).map(([id, name]) => ({ id, name }))
}

const viewDetail = (row: EIP) => {
  selectedEip.value = row
  detailVisible.value = true
}

const handleCreate = () => {
  createForm.project_id = null
  createForm.cloud_account_id = filters.cloud_account_id
  createForm.billing_method = 'bandwidth'
  createForm.bandwidth = 100
  createForm.name = ''
  createForm.tags = []
  createVisible.value = true
}

const addTag = () => {
  createForm.tags.push({ key: '', value: '' })
}

const removeTag = (idx: number) => {
  createForm.tags.splice(idx, 1)
}

const confirmCreate = async () => {
  try {
    await createFormRef.value.validate()
    createLoading.value = true
    const data: any = {
      cloud_account_id: createForm.cloud_account_id,
      bandwidth: createForm.bandwidth,
      billing_method: createForm.billing_method
    }
    if (createForm.project_id) data.project_id = createForm.project_id
    if (createForm.name) data.name = createForm.name
    if (createForm.tags.length > 0) {
      data.tags = createForm.tags.filter(t => t.key && t.value)
    }
    await createEIP(data)
    ElMessage.success('申请成功')
    createVisible.value = false
    loadEips()
  } catch (e) {
    if (e !== false) ElMessage.error('申请失败')
  } finally {
    createLoading.value = false
  }
}

const handleBind = async (row: EIP) => {
  selectedEip.value = row
  bindForm.eip_id = String(row.id)
  bindForm.resource_type = 'guest'
  bindForm.resource_id = ''
  // TODO: 根据资源类型加载可用资源列表
  availableResources.value = []
  bindVisible.value = true
}

const confirmBind = async () => {
  if (!bindForm.resource_id) {
    ElMessage.warning('请选择要绑定的资源')
    return
  }
  try {
    await bindEIP(selectedEip.value!.id, bindForm.resource_id)
    ElMessage.success('绑定成功')
    bindVisible.value = false
    loadEips()
  } catch (e) {
    ElMessage.error('绑定失败')
  }
}

const handleUnbind = async (row: EIP) => {
  try {
    await ElMessageBox.confirm('确认解绑此弹性IP?', '提示', { type: 'warning' })
    await unbindEIP(row.id)
    ElMessage.success('解绑成功')
    loadEips()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('解绑失败')
  }
}

const handleEdit = (row: EIP) => {
  selectedEip.value = row
  editForm.bandwidth = row.bandwidth
  editVisible.value = true
}

const confirmEdit = async () => {
  try {
    // TODO: 实现修改带宽 API
    ElMessage.success('修改成功')
    editVisible.value = false
    loadEips()
  } catch (e) {
    ElMessage.error('修改失败')
  }
}

const handleDelete = async (row: EIP) => {
  try {
    await ElMessageBox.confirm(`确认删除弹性公网IP "${row.name || row.address}"？此操作不可恢复`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteEIP(row.id)
    ElMessage.success('删除成功')
    loadEips()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('删除失败')
  }
}

const handleAddTagToSelected = async () => {
  if (!tagForm.key) {
    ElMessage.warning('请输入标签键')
    return
  }
  // TODO: 实现批量添加标签 API
  ElMessage.success(`已为 ${selectedRows.value.length} 个EIP添加标签`)
}

const handleRemoveTagFromSelected = async () => {
  if (!tagForm.key) {
    ElMessage.warning('请输入要移除的标签键')
    return
  }
  // TODO: 实现批量移除标签 API
  ElMessage.success(`已从 ${selectedRows.value.length} 个EIP移除标签`)
}

const resetFilters = () => {
  filters.cloud_account_id = null
  filters.status = ''
  filters.billing_method = ''
  filters.region = ''
  searchKeyword.value = ''
  pagination.page = 1
  loadEips()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadEips()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadEips()
}

onMounted(() => {
  loadEips()
  loadProjects()
})
</script>

<style scoped>
.eips-container {
  padding: 20px;
}

.filter-tabs {
  margin-bottom: 16px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.toolbar-left {
  display: flex;
  gap: 8px;
}

.filter-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
}

.platform-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.account-name {
  font-size: 12px;
  color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.resource-type {
  font-size: 12px;
  color: #909399;
  margin-left: 4px;
}

.tag-item {
  margin-right: 4px;
  margin-bottom: 2px;
}

.tags-editor {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tag-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.unit-label {
  margin-left: 8px;
  color: #606266;
}
</style>