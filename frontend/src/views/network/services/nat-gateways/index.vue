<template>
  <div class="nat-gateways-container">
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
              <el-dropdown-item command="delete" :disabled="selectedRows.length === 0">批量删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button @click="handleTagManage">标签</el-button>
      </div>
      <div class="toolbar-right">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索名称"
          clearable
          style="width: 200px"
          @keyup.enter="loadNATGateways"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
    </div>

    <!-- 过滤栏 -->
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="filters" @submit.prevent="loadNATGateways">
        <el-form-item label="云账号">
          <CloudAccountSelector
            v-model:value="filters.cloud_account_id"
            placeholder="全部云账号"
            @change="handleAccountChange"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 120px">
            <el-option label="可用" value="Available" />
            <el-option label="创建中" value="Creating" />
            <el-option label="删除中" value="Deleting" />
            <el-option label="错误" value="Error" />
          </el-select>
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="filters.nat_type" placeholder="全部类型" clearable style="width: 120px">
            <el-option label="公网NAT" value="public" />
            <el-option label="私网NAT" value="private" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-select v-model="filters.region" placeholder="全部区域" clearable style="width: 150px">
            <el-option v-for="r in regions" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadNATGateways">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格 -->
    <el-table
      :data="natGateways"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="50" />
      <el-table-column label="名称" min-width="150">
        <template #default="{ row }">
          <el-link type="primary" @click="viewDetail(row)">{{ row.name || row.nat_gateway_id }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">
            {{ getStatusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="类型" width="100">
        <template #default="{ row }">
          {{ getNatTypeLabel(row.nat_type) }}
        </template>
      </el-table-column>
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
      <el-table-column label="规则" width="80">
        <template #default="{ row }">
          <el-link type="primary" @click="manageRules(row)">
            {{ row.snat_table_entries + row.dnat_table_entries }} 条
          </el-link>
        </template>
      </el-table-column>
      <el-table-column label="规格" width="120">
        <template #default="{ row }">
          {{ row.specification || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="计费方式" width="100">
        <template #default="{ row }">
          {{ getBillingLabel(row.billing_method) }}
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
      <el-table-column label="所属VPC" width="150">
        <template #default="{ row }">
          <el-link type="primary" v-if="row.vpc_name">{{ row.vpc_name }}</el-link>
          <span v-else>{{ row.vpc_id || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="所属域" width="120">
        <template #default="{ row }">
          {{ row.owner_domain || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="region" label="区域" width="120" />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button size="small" link type="primary" @click="manageRules(row)">管理规则</el-button>
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
    <el-dialog v-model="detailVisible" title="NAT网关详情" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedNatGateway?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedNatGateway?.name }}</el-descriptions-item>
        <el-descriptions-item label="NAT ID">{{ selectedNatGateway?.nat_gateway_id }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(selectedNatGateway?.status)">
            {{ getStatusLabel(selectedNatGateway?.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="类型">{{ getNatTypeLabel(selectedNatGateway?.nat_type) }}</el-descriptions-item>
        <el-descriptions-item label="规格">{{ selectedNatGateway?.specification || '-' }}</el-descriptions-item>
        <el-descriptions-item label="计费方式">{{ getBillingLabel(selectedNatGateway?.billing_method) }}</el-descriptions-item>
        <el-descriptions-item label="SNAT规则数">{{ selectedNatGateway?.snat_table_entries }}</el-descriptions-item>
        <el-descriptions-item label="DNAT规则数">{{ selectedNatGateway?.dnat_table_entries }}</el-descriptions-item>
        <el-descriptions-item label="绑定EIP">{{ selectedNatGateway?.eip_address || '-' }}</el-descriptions-item>
        <el-descriptions-item label="所属VPC">{{ selectedNatGateway?.vpc_name || selectedNatGateway?.vpc_id }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ selectedNatGateway?.owner_domain || '-' }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedNatGateway?.region }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedNatGateway?.account_name }}</el-descriptions-item>
        <el-descriptions-item label="标签" :span="2">
          <template v-if="selectedNatGateway?.tags && selectedNatGateway.tags.length > 0">
            <el-tag v-for="tag in selectedNatGateway.tags" :key="tag.key" class="tag-item">
              {{ tag.key }}: {{ tag.value }}
            </el-tag>
          </template>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ selectedNatGateway?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 规则管理对话框 -->
    <el-dialog v-model="rulesVisible" title="管理NAT规则" width="900px">
      <div class="rules-header">
        <span class="nat-name">NAT网关: {{ selectedNatGateway?.name }}</span>
        <el-button type="primary" size="small" @click="showAddRule">添加规则</el-button>
      </div>
      <el-tabs v-model="ruleTab">
        <el-tab-pane label="SNAT规则" name="snat">
          <el-table :data="snatRules" v-loading="rulesLoading" style="width: 100%">
            <el-table-column prop="name" label="规则名称" width="150" />
            <el-table-column prop="external_ip" label="外部IP" width="150" />
            <el-table-column prop="internal_ip" label="内部IP/网段" width="150" />
            <el-table-column prop="protocol" label="协议" width="80" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button size="small" link type="danger" @click="deleteRule(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="DNAT规则" name="dnat">
          <el-table :data="dnatRules" v-loading="rulesLoading" style="width: 100%">
            <el-table-column prop="name" label="规则名称" width="150" />
            <el-table-column prop="external_ip" label="外部IP" width="120" />
            <el-table-column prop="external_port" label="外部端口" width="100" />
            <el-table-column prop="internal_ip" label="内部IP" width="120" />
            <el-table-column prop="internal_port" label="内部端口" width="100" />
            <el-table-column prop="protocol" label="协议" width="80" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button size="small" link type="danger" @click="deleteRule(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
      <template #footer>
        <el-button @click="rulesVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 新建NAT网关对话框 -->
    <el-dialog v-model="createVisible" title="新建NAT网关" width="600px">
      <el-form :model="createForm" label-width="120px" :rules="createRules" ref="createFormRef">
        <el-form-item label="所属域">
          <el-select v-model="createForm.domain" placeholder="选择所属域" style="width: 100%">
            <el-option label="system" value="system" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="createForm.name" placeholder="输入NAT网关名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="createForm.description" type="textarea" placeholder="输入描述" />
        </el-form-item>
        <el-form-item label="计费方式">
          <el-radio-group v-model="createForm.billing_method">
            <el-radio-button value="Postpaid">按量付费</el-radio-button>
            <el-radio-button value="Prepaid">包年包月</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="区域" prop="region_id">
          <el-select v-model="createForm.region_id" placeholder="选择区域" style="width: 100%">
            <el-option v-for="r in regions" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="规格" prop="specification">
          <el-select v-model="createForm.specification" placeholder="选择规格" style="width: 100%">
            <el-option label="小型" value="small" />
            <el-option label="中型" value="medium" />
            <el-option label="大型" value="large" />
          </el-select>
        </el-form-item>
        <el-form-item label="VPC" prop="vpc_id">
          <el-select v-model="createForm.vpc_id" placeholder="选择VPC" style="width: 100%">
            <el-option v-for="vpc in vpcs" :key="vpc.id" :label="vpc.name" :value="vpc.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="网络">
          <el-select v-model="createForm.subnet_id" placeholder="选择网络" style="width: 100%">
            <el-option v-for="subnet in subnets" :key="subnet.id" :label="subnet.name" :value="subnet.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="绑定EIP">
          <el-select v-model="createForm.eip_id" placeholder="选择EIP（可选）" style="width: 100%" clearable>
            <el-option v-for="eip in availableEips" :key="eip.eip_id" :label="eip.name || eip.address" :value="eip.eip_id" />
          </el-select>
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

    <!-- 添加规则对话框 -->
    <el-dialog v-model="addRuleVisible" title="添加NAT规则" width="500px">
      <el-form :model="addRuleForm" label-width="100px">
        <el-form-item label="规则类型">
          <el-radio-group v-model="addRuleForm.rule_type">
            <el-radio-button value="SNAT">SNAT</el-radio-button>
            <el-radio-button value="DNAT">DNAT</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="规则名称">
          <el-input v-model="addRuleForm.name" placeholder="输入规则名称" />
        </el-form-item>
        <el-form-item label="外部IP" prop="external_ip">
          <el-input v-model="addRuleForm.external_ip" placeholder="外部IP地址" />
        </el-form-item>
        <el-form-item label="外部端口" v-if="addRuleForm.rule_type === 'DNAT'">
          <el-input v-model="addRuleForm.external_port" placeholder="外部端口" />
        </el-form-item>
        <el-form-item label="内部IP" prop="internal_ip">
          <el-input v-model="addRuleForm.internal_ip" placeholder="内部IP地址或网段" />
        </el-form-item>
        <el-form-item label="内部端口" v-if="addRuleForm.rule_type === 'DNAT'">
          <el-input v-model="addRuleForm.internal_port" placeholder="内部端口" />
        </el-form-item>
        <el-form-item label="协议">
          <el-select v-model="addRuleForm.protocol" style="width: 100%">
            <el-option label="ALL" value="ALL" />
            <el-option label="TCP" value="TCP" />
            <el-option label="UDP" value="UDP" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addRuleVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAddRule">确定</el-button>
      </template>
    </el-dialog>

    <!-- 修改NAT网关对话框 -->
    <el-dialog v-model="editVisible" title="修改NAT网关" width="400px">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="editForm.description" type="textarea" />
        </el-form-item>
        <el-form-item label="规格">
          <el-select v-model="editForm.specification" style="width: 100%">
            <el-option label="小型" value="small" />
            <el-option label="中型" value="medium" />
            <el-option label="大型" value="large" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmEdit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown, Search } from '@element-plus/icons-vue'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import {
  getNATGateways,
  getNATGateway,
  createNATGateway,
  updateNATGateway,
  deleteNATGateway,
  batchDeleteNATGateways,
  getNATRules,
  createNATRule,
  deleteNATRule,
  type NATGateway,
  type NATRule
} from '@/api/networkSync'
import { getEIPs, type EIP } from '@/api/networkSync'
import { getVPCs, type VPC } from '@/api/networkSync'
import { getSubnets, type Subnet } from '@/api/networkSync'

const natGateways = ref<NATGateway[]>([])
const loading = ref(false)
const detailVisible = ref(false)
const rulesVisible = ref(false)
const createVisible = ref(false)
const editVisible = ref(false)
const addRuleVisible = ref(false)
const selectedNatGateway = ref<NATGateway | null>(null)
const selectedRows = ref<NATGateway[]>([])
const createLoading = ref(false)
const rulesLoading = ref(false)

// Tabs
const activeTab = ref('all')
const ruleTab = ref('snat')

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
  nat_type: '',
  region: ''
})

// 可用资源
const regions = ref<{ id: string; name: string }[]>([])
const vpcs = ref<VPC[]>([])
const subnets = ref<Subnet[]>([])
const availableEips = ref<EIP[]>([])

// NAT规则
const natRules = ref<NATRule[]>([])
const snatRules = computed(() => natRules.value.filter(r => r.rule_type === 'SNAT'))
const dnatRules = computed(() => natRules.value.filter(r => r.rule_type === 'DNAT'))

// 创建表单
const createForm = reactive({
  domain: 'system',
  name: '',
  description: '',
  billing_method: 'Postpaid',
  region_id: '',
  specification: '',
  vpc_id: '',
  subnet_id: '',
  eip_id: '',
  tags: [] as { key: string; value: string }[]
})

const createRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  region_id: [{ required: true, message: '请选择区域', trigger: 'change' }],
  specification: [{ required: true, message: '请选择规格', trigger: 'change' }],
  vpc_id: [{ required: true, message: '请选择VPC', trigger: 'change' }]
}

const createFormRef = ref()

// 修改表单
const editForm = reactive({
  name: '',
  description: '',
  specification: ''
})

// 添加规则表单
const addRuleForm = reactive({
  rule_type: 'SNAT',
  name: '',
  external_ip: '',
  external_port: '',
  internal_ip: '',
  internal_port: '',
  protocol: 'ALL'
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
  Postpaid: '按量付费',
  Prepaid: '包年包月',
  PayByTraffic: '按流量',
  PayByBandwidth: '按带宽'
}

const getBillingLabel = (method: string): string => {
  return billingLabels[method] || method || '-'
}

// NAT类型映射
const natTypeLabels: Record<string, string> = {
  public: '公网NAT',
  private: '私网NAT'
}

const getNatTypeLabel = (type: string): string => {
  return natTypeLabels[type] || type || '-'
}

// 状态映射
const statusLabels: Record<string, string> = {
  Available: '可用',
  Creating: '创建中',
  Deleting: '删除中',
  Error: '错误',
  available: '可用',
  creating: '创建中'
}

const getStatusLabel = (status: string): string => {
  return statusLabels[status] || status
}

const getStatusType = (status: string): '' | 'success' | 'warning' | 'info' | 'danger' => {
  if (status === 'Available' || status === 'available') return 'success'
  if (status === 'Creating' || status === 'creating') return 'warning'
  if (status === 'Deleting') return 'warning'
  if (status === 'Error' || status === 'error') return 'danger'
  return 'info'
}

const handleTabChange = (tab: string) => {
  activeTab.value = tab
  loadNATGateways()
}

const handleAccountChange = (accountId: number | null) => {
  filters.cloud_account_id = accountId
}

const handleSelectionChange = (rows: NATGateway[]) => {
  selectedRows.value = rows
}

const handleBatchCommand = (command: string) => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请先选择要操作的NAT网关')
    return
  }
  if (command === 'delete') {
    handleBatchDelete()
  }
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(`确认批量删除 ${selectedRows.value.length} 个NAT网关？此操作不可恢复`, '警告', { type: 'warning' })
    const ids = selectedRows.value.map(n => Number(n.id))
    await batchDeleteNATGateways(ids)
    ElMessage.success('批量删除成功')
    loadNATGateways()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('批量删除失败')
  }
}

const handleTagManage = () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请先选择要管理标签的NAT网关')
    return
  }
  ElMessage.info('标签管理功能开发中')
}

const loadNATGateways = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (filters.cloud_account_id) params.cloud_account_id = filters.cloud_account_id
    if (filters.status) params.status = filters.status
    if (filters.nat_type) params.nat_type = filters.nat_type
    if (filters.region) params.region = filters.region
    if (searchKeyword.value) params.name = searchKeyword.value
    if (activeTab.value !== 'all') params.platform = activeTab.value

    const res = await getNATGateways(params)
    natGateways.value = res.items || []
    pagination.total = res.total || 0

    // 提取区域列表
    const regionMap = new Map<string, string>()
    natGateways.value.forEach(nat => {
      if (nat.region_id && nat.region) {
        regionMap.set(nat.region_id, nat.region)
      }
    })
    regions.value = Array.from(regionMap.entries()).map(([id, name]) => ({ id, name }))
  } catch (e) {
    console.error(e)
    ElMessage.error('加载NAT网关列表失败')
  } finally {
    loading.value = false
  }
}

const loadAvailableResources = async () => {
  try {
    // 加载VPC
    const vpcRes = await getVPCs({ page: 1, page_size: 100 })
    vpcs.value = vpcRes.items || []

    // 加载子网
    const subnetRes = await getSubnets({ page: 1, page_size: 100 })
    subnets.value = subnetRes.items || []

    // 加载可用的EIP
    const eipRes = await getEIPs({ page: 1, page_size: 100, status: 'Available' })
    availableEips.value = eipRes.items || []
  } catch (e) {
    console.error('加载可用资源失败', e)
  }
}

const viewDetail = (row: NATGateway) => {
  selectedNatGateway.value = row
  detailVisible.value = true
}

const manageRules = async (row: NATGateway) => {
  selectedNatGateway.value = row
  rulesVisible.value = true
  rulesLoading.value = true
  try {
    const res = await getNATRules(row.id)
    natRules.value = res.items || []
  } catch (e) {
    console.error('加载规则失败', e)
    natRules.value = []
  } finally {
    rulesLoading.value = false
  }
}

const showAddRule = () => {
  addRuleForm.rule_type = 'SNAT'
  addRuleForm.name = ''
  addRuleForm.external_ip = selectedNatGateway.value?.eip_address || ''
  addRuleForm.external_port = ''
  addRuleForm.internal_ip = ''
  addRuleForm.internal_port = ''
  addRuleForm.protocol = 'ALL'
  addRuleVisible.value = true
}

const confirmAddRule = async () => {
  if (!addRuleForm.external_ip || !addRuleForm.internal_ip) {
    ElMessage.warning('请填写外部IP和内部IP')
    return
  }
  try {
    await createNATRule(selectedNatGateway.value!.id, addRuleForm)
    ElMessage.success('添加规则成功')
    addRuleVisible.value = false
    manageRules(selectedNatGateway.value!)
  } catch (e) {
    ElMessage.error('添加规则失败')
  }
}

const deleteRule = async (row: NATRule) => {
  try {
    await ElMessageBox.confirm(`确认删除规则 "${row.name || row.rule_id}"？`, '警告', { type: 'warning' })
    await deleteNATRule(selectedNatGateway.value!.id, row.id)
    ElMessage.success('删除成功')
    manageRules(selectedNatGateway.value!)
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('删除失败')
  }
}

const handleCreate = () => {
  createForm.domain = 'system'
  createForm.name = ''
  createForm.description = ''
  createForm.billing_method = 'Postpaid'
  createForm.region_id = ''
  createForm.specification = ''
  createForm.vpc_id = ''
  createForm.subnet_id = ''
  createForm.eip_id = ''
  createForm.tags = []
  createVisible.value = true
  loadAvailableResources()
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
      name: createForm.name,
      billing_method: createForm.billing_method,
      region_id: createForm.region_id,
      specification: createForm.specification,
      vpc_id: createForm.vpc_id,
      subnet_id: createForm.subnet_id,
      eip_id: createForm.eip_id
    }
    if (createForm.description) data.description = createForm.description
    if (createForm.tags.length > 0) {
      data.tags = createForm.tags.filter(t => t.key && t.value)
    }
    await createNATGateway(data)
    ElMessage.success('创建成功')
    createVisible.value = false
    loadNATGateways()
  } catch (e) {
    if (e !== false) ElMessage.error('创建失败')
  } finally {
    createLoading.value = false
  }
}

const handleEdit = (row: NATGateway) => {
  selectedNatGateway.value = row
  editForm.name = row.name
  editForm.description = row.description
  editForm.specification = row.specification
  editVisible.value = true
}

const confirmEdit = async () => {
  try {
    await updateNATGateway(selectedNatGateway.value!.id, editForm)
    ElMessage.success('修改成功')
    editVisible.value = false
    loadNATGateways()
  } catch (e) {
    ElMessage.error('修改失败')
  }
}

const handleDelete = async (row: NATGateway) => {
  try {
    await ElMessageBox.confirm(`确认删除NAT网关 "${row.name}"？此操作不可恢复`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteNATGateway(row.id)
    ElMessage.success('删除成功')
    loadNATGateways()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('删除失败')
  }
}

const resetFilters = () => {
  filters.cloud_account_id = null
  filters.status = ''
  filters.nat_type = ''
  filters.region = ''
  searchKeyword.value = ''
  pagination.page = 1
  loadNATGateways()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadNATGateways()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadNATGateways()
}

onMounted(() => {
  loadNATGateways()
})
</script>

<style scoped>
.nat-gateways-container {
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

.rules-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.nat-name {
  font-weight: 500;
}
</style>