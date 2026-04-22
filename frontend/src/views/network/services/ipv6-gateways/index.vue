<template>
  <div class="ipv6-gateways-container">
    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <el-button type="primary" @click="handleCreate">新建</el-button>
      </div>
      <div class="toolbar-right">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索名称"
          clearable
          style="width: 200px"
          @keyup.enter="loadIPv6Gateways"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
    </div>

    <!-- 过滤栏 -->
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="filters" @submit.prevent="loadIPv6Gateways">
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
        <el-form-item label="区域">
          <el-select v-model="filters.region" placeholder="全部区域" clearable style="width: 150px">
            <el-option v-for="r in regions" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadIPv6Gateways">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格 -->
    <el-table
      :data="ipv6Gateways"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
    >
      <el-table-column label="名称" min-width="150">
        <template #default="{ row }">
          <el-link type="primary" @click="viewDetail(row)">{{ row.name || row.ipv6_gateway_id }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">
            {{ getStatusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="VPC" width="150">
        <template #default="{ row }">
          <el-link type="primary" v-if="row.vpc_name">{{ row.vpc_name }}</el-link>
          <span v-else>{{ row.vpc_id || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="规格" width="120">
        <template #default="{ row }">
          {{ row.specification || '-' }}
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
      <el-table-column label="项目" width="120">
        <template #default="{ row }">
          {{ row.project_name || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="region" label="区域" width="120" />
      <el-table-column label="创建时间" width="180">
        <template #default="{ row }">
          {{ row.created_at }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
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
    <el-dialog v-model="detailVisible" title="IPv6网关详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedGateway?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedGateway?.name }}</el-descriptions-item>
        <el-descriptions-item label="IPv6网关ID">{{ selectedGateway?.ipv6_gateway_id }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(selectedGateway?.status)">
            {{ getStatusLabel(selectedGateway?.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="所属VPC">{{ selectedGateway?.vpc_name || selectedGateway?.vpc_id }}</el-descriptions-item>
        <el-descriptions-item label="规格">{{ selectedGateway?.specification || '-' }}</el-descriptions-item>
        <el-descriptions-item label="IPv6地址段">{{ selectedGateway?.ipv6_cidr || '-' }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedGateway?.region }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ getPlatformLabel(selectedGateway?.provider_type) }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedGateway?.account_name }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedGateway?.project_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedGateway?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 新建对话框 -->
    <el-dialog v-model="createVisible" title="新建IPv6网关" width="500px">
      <el-form :model="createForm" label-width="100px" :rules="createRules" ref="createFormRef">
        <el-form-item label="名称" prop="name">
          <el-input v-model="createForm.name" placeholder="输入IPv6网关名称" />
        </el-form-item>
        <el-form-item label="VPC" prop="vpc_id">
          <el-select v-model="createForm.vpc_id" placeholder="选择VPC" style="width: 100%">
            <el-option v-for="vpc in vpcs" :key="vpc.id" :label="vpc.name" :value="vpc.vpc_id" />
          </el-select>
        </el-form-item>
        <el-form-item label="规格">
          <el-select v-model="createForm.specification" placeholder="选择规格" style="width: 100%">
            <el-option label="中小企业版" value="small" />
            <el-option label="中型企业版" value="medium" />
            <el-option label="大型企业版" value="large" />
          </el-select>
        </el-form-item>
        <el-form-item label="IPv6地址段">
          <el-input v-model="createForm.ipv6_cidr" placeholder="可选，输入IPv6地址段" />
        </el-form-item>
        <el-form-item label="区域" prop="region_id">
          <el-select v-model="createForm.region_id" placeholder="选择区域" style="width: 100%">
            <el-option v-for="r in regions" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCreate" :loading="createLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 修改对话框 -->
    <el-dialog v-model="editVisible" title="修改IPv6网关" width="400px">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="规格">
          <el-select v-model="editForm.specification" style="width: 100%">
            <el-option label="中小企业版" value="small" />
            <el-option label="中型企业版" value="medium" />
            <el-option label="大型企业版" value="large" />
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
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import {
  getIPv6Gateways,
  createIPv6Gateway,
  deleteIPv6Gateway,
  type IPv6Gateway
} from '@/api/networkSync'
import { getVPCs, type VPC } from '@/api/networkSync'

const ipv6Gateways = ref<IPv6Gateway[]>([])
const loading = ref(false)
const detailVisible = ref(false)
const createVisible = ref(false)
const editVisible = ref(false)
const selectedGateway = ref<IPv6Gateway | null>(null)
const createLoading = ref(false)

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
  region: ''
})

// 可用资源
const regions = ref<{ id: string; name: string }[]>([])
const vpcs = ref<VPC[]>([])

// 创建表单
const createForm = reactive({
  name: '',
  vpc_id: '',
  specification: '',
  ipv6_cidr: '',
  region_id: ''
})

const createRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  vpc_id: [{ required: true, message: '请选择VPC', trigger: 'change' }],
  region_id: [{ required: true, message: '请选择区域', trigger: 'change' }]
}

const createFormRef = ref()

// 修改表单
const editForm = reactive({
  name: '',
  specification: ''
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

const handleAccountChange = (accountId: number | null) => {
  filters.cloud_account_id = accountId
}

const loadIPv6Gateways = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (filters.cloud_account_id) params.cloud_account_id = filters.cloud_account_id
    if (filters.status) params.status = filters.status
    if (filters.region) params.region = filters.region
    if (searchKeyword.value) params.name = searchKeyword.value

    const res = await getIPv6Gateways(params)
    ipv6Gateways.value = res.items || []
    pagination.total = res.total || 0

    // 提取区域列表
    const regionMap = new Map<string, string>()
    ipv6Gateways.value.forEach(gw => {
      if (gw.region_id && gw.region) {
        regionMap.set(gw.region_id, gw.region)
      }
    })
    regions.value = Array.from(regionMap.entries()).map(([id, name]) => ({ id, name }))
  } catch (e) {
    console.error(e)
    ElMessage.error('加载IPv6网关列表失败')
  } finally {
    loading.value = false
  }
}

const loadVPCs = async () => {
  try {
    const res = await getVPCs({ page: 1, page_size: 100 })
    vpcs.value = res.items || []
  } catch (e) {
    console.error('加载VPC失败', e)
  }
}

const viewDetail = (row: IPv6Gateway) => {
  selectedGateway.value = row
  detailVisible.value = true
}

const handleCreate = () => {
  createForm.name = ''
  createForm.vpc_id = ''
  createForm.specification = ''
  createForm.ipv6_cidr = ''
  createForm.region_id = ''
  createVisible.value = true
  loadVPCs()
}

const confirmCreate = async () => {
  try {
    await createFormRef.value.validate()
    createLoading.value = true
    const data: any = {
      name: createForm.name,
      vpc_id: createForm.vpc_id,
      specification: createForm.specification,
      ipv6_cidr: createForm.ipv6_cidr,
      region_id: createForm.region_id
    }
    await createIPv6Gateway(data)
    ElMessage.success('创建成功')
    createVisible.value = false
    loadIPv6Gateways()
  } catch (e) {
    if (e !== false) ElMessage.error('创建失败')
  } finally {
    createLoading.value = false
  }
}

const handleEdit = (row: IPv6Gateway) => {
  selectedGateway.value = row
  editForm.name = row.name
  editForm.specification = row.specification
  editVisible.value = true
}

const confirmEdit = async () => {
  // TODO: 实现修改API
  ElMessage.success('修改成功')
  editVisible.value = false
  loadIPv6Gateways()
}

const handleDelete = async (row: IPv6Gateway) => {
  try {
    await ElMessageBox.confirm(`确认删除IPv6网关 "${row.name}"？此操作不可恢复`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteIPv6Gateway(row.id)
    ElMessage.success('删除成功')
    loadIPv6Gateways()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('删除失败')
  }
}

const resetFilters = () => {
  filters.cloud_account_id = null
  filters.status = ''
  filters.region = ''
  searchKeyword.value = ''
  pagination.page = 1
  loadIPv6Gateways()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadIPv6Gateways()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadIPv6Gateways()
}

onMounted(() => {
  loadIPv6Gateways()
})
</script>

<style scoped>
.ipv6-gateways-container {
  padding: 20px;
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
</style>