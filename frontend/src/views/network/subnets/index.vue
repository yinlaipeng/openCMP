<template>
  <div class="ip-subnets-container">
    <div class="page-header">
      <h2>IP子网</h2>
      <div class="toolbar">
        <el-button link type="primary" @click="handleView">
          <el-icon><View /></el-icon>
          查看
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建
        </el-button>
        <el-dropdown trigger="click" @command="handleBatchCommand" :disabled="selectedSubnets.length === 0">
          <el-button :disabled="selectedSubnets.length === 0">
            批量操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="tags">设置标签</el-dropdown-item>
              <el-dropdown-item command="schedtag">调整调度标签</el-dropdown-item>
              <el-dropdown-item command="delete" divided>批量删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button @click="handleTags">
          <el-icon><PriceTag /></el-icon>
          标签
        </el-button>
      </div>
    </div>

    <!-- 搜索表单 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="fetchData">
        <el-form-item label="名称">
          <el-input v-model="filters.name" placeholder="搜索子网名称" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="选择状态" clearable style="width: 120px">
            <el-option label="可用" value="available" />
            <el-option label="创建中" value="creating" />
            <el-option label="错误" value="error" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="filters.platform" placeholder="选择平台" clearable style="width: 120px">
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item label="项目">
          <el-select v-model="filters.project_id" placeholder="选择项目" clearable style="width: 150px">
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-input v-model="filters.region" placeholder="区域" clearable style="width: 120px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <el-table
        :data="ipSubnets"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <!-- Cloudpods 表头顺序: Name → Status → Type → Auto scheduling → IP Address → IPv6 Address → Usage → Scheduler Tag → Platform → Project → Region -->
        <el-table-column prop="name" label="名称" width="180" fixed="left">
          <template #default="{ row }">
            <el-link type="primary" :underline="false" @click="handleDetails(row)">
              {{ row.name }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_classic" label="类型" width="80">
          <template #default="{ row }">
            <el-tag size="small" :type="row.is_classic ? 'info' : 'primary'">
              {{ row.is_classic ? '经典' : 'VPC' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_auto_alloc" label="自动调度" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_auto_alloc ? 'success' : 'info'" size="small">
              {{ row.is_auto_alloc ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="IP地址段" width="180">
          <template #default="{ row }">
            {{ getIPRange(row) }}
          </template>
        </el-table-column>
        <el-table-column label="IPv6地址" width="150">
          <template #default="{ row }">
            {{ row.guest_ip6_mask || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="使用情况" width="120">
          <template #default="{ row }">
            {{ getUsagePercent(row) }}
          </template>
        </el-table-column>
        <el-table-column prop="schedtag" label="调度标签" width="120">
          <template #default="{ row }">
            {{ row.schedtag || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="provider_type" label="平台" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="getPlatformType(row.provider_type)">
              {{ getPlatformLabel(row.provider_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="project_name" label="项目" width="150" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleAdjustSchedtag(row)">调整调度标签</el-button>
            <el-dropdown trigger="click">
              <el-button size="small">
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleEdit(row)">编辑</el-dropdown-item>
                  <el-dropdown-item @click="handleAllocation(row)">分配</el-dropdown-item>
                  <el-dropdown-item divided @click="handleDelete(row)">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        class="pagination"
      />
    </el-card>

    <!-- 创建IP子网弹窗 -->
    <el-dialog
      title="创建IP子网"
      v-model="createDialogVisible"
      width="500px"
    >
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="100px">
        <el-form-item label="云账户" prop="cloud_account_id">
          <el-select v-model="createForm.cloud_account_id" placeholder="选择云账户" style="width: 100%" @change="handleAccountChange">
            <el-option v-for="acc in cloudAccounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入子网名称" />
        </el-form-item>
        <el-form-item label="CIDR" prop="cidr">
          <el-input v-model="createForm.cidr" placeholder="例如: 192.168.1.0/24" />
        </el-form-item>
        <el-form-item label="VPC">
          <el-select v-model="createForm.vpc_id" placeholder="选择VPC" style="width: 100%">
            <el-option v-for="vpc in vpcs" :key="vpc.id" :label="vpc.name" :value="vpc.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="自动分配">
          <el-switch v-model="createForm.is_auto_alloc" />
        </el-form-item>
        <el-form-item label="项目">
          <el-select v-model="createForm.project_id" placeholder="选择项目" clearable style="width: 100%">
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCreate" :loading="creating">创建</el-button>
      </template>
    </el-dialog>

    <!-- 调度标签弹窗 -->
    <el-dialog
      title="调整调度标签"
      v-model="schedtagDialogVisible"
      width="400px"
    >
      <el-form :model="schedtagForm" label-width="100px">
        <el-form-item label="子网">
          {{ schedtagForm.subnet_name }}
        </el-form-item>
        <el-form-item label="调度标签">
          <el-select v-model="schedtagForm.schedtag" placeholder="选择调度标签" style="width: 100%">
            <el-option label="无" value="" />
            <el-option label="高可用" value="ha" />
            <el-option label="性能优先" value="performance" />
            <el-option label="成本优先" value="cost" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="schedtagDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmSchedtag">确定</el-button>
      </template>
    </el-dialog>

    <!-- 详情弹窗 -->
    <el-dialog
      title="IP子网详情"
      v-model="detailDialogVisible"
      width="700px"
    >
      <el-descriptions :column="2" border v-if="selectedSubnet">
        <el-descriptions-item label="ID">{{ selectedSubnet.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedSubnet.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(selectedSubnet.status)">{{ getStatusLabel(selectedSubnet.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="类型">{{ selectedSubnet.is_classic ? '经典网络' : 'VPC网络' }}</el-descriptions-item>
        <el-descriptions-item label="CIDR">{{ selectedSubnet.cidr }}</el-descriptions-item>
        <el-descriptions-item label="IP地址段">{{ getIPRange(selectedSubnet) }}</el-descriptions-item>
        <el-descriptions-item label="IPv6地址">{{ selectedSubnet.guest_ip6_mask || '-' }}</el-descriptions-item>
        <el-descriptions-item label="网关">{{ selectedSubnet.guest_gateway || '-' }}</el-descriptions-item>
        <el-descriptions-item label="使用情况">{{ getUsagePercent(selectedSubnet) }}</el-descriptions-item>
        <el-descriptions-item label="自动调度">{{ selectedSubnet.is_auto_alloc ? '是' : '否' }}</el-descriptions-item>
        <el-descriptions-item label="调度标签">{{ selectedSubnet.schedtag || '-' }}</el-descriptions-item>
        <el-descriptions-item label="平台">
          <el-tag size="small" :type="getPlatformType(selectedSubnet.provider_type)">
            {{ getPlatformLabel(selectedSubnet.provider_type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedSubnet.project_name }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedSubnet.region }}</el-descriptions-item>
        <el-descriptions-item label="DNS">{{ selectedSubnet.dns || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Plus, ArrowDown, PriceTag } from '@element-plus/icons-vue'
import { getSubnets, createSubnet, deleteSubnet, batchDeleteSubnets, Subnet } from '@/api/networkSync'
import { getCloudAccounts } from '@/api/cloud-account'
import type { CloudAccount } from '@/types'

// Data
const loading = ref(false)
const creating = ref(false)
const ipSubnets = ref<Subnet[]>([])
const selectedSubnets = ref<Subnet[]>([])
const cloudAccounts = ref<CloudAccount[]>([])
const projects = ref<any[]>([])
const vpcs = ref<any[]>([])

const createDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const schedtagDialogVisible = ref(false)
const selectedSubnet = ref<Subnet | null>(null)

const filters = reactive({
  name: '',
  status: '',
  platform: '',
  project_id: null as number | null,
  region: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const createForm = reactive({
  cloud_account_id: '',
  name: '',
  cidr: '',
  vpc_id: '',
  is_auto_alloc: true,
  project_id: ''
})

const schedtagForm = reactive({
  subnet_id: '',
  subnet_name: '',
  schedtag: ''
})

const createRules = {
  cloud_account_id: [{ required: true, message: '请选择云账户', trigger: 'change' }],
  name: [{ required: true, message: '请输入子网名称', trigger: 'blur' }],
  cidr: [{ required: true, message: '请输入CIDR', trigger: 'blur' }]
}

// Methods
const getStatusType = (status: string) => {
  switch (status) {
    case 'available': return 'success'
    case 'creating': return 'warning'
    case 'error': return 'danger'
    default: return 'info'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'available': return '可用'
    case 'creating': return '创建中'
    case 'error': return '错误'
    default: return status
  }
}

const getPlatformLabel = (platform: string) => {
  const labels: Record<string, string> = {
    aliyun: '阿里云',
    alibaba: '阿里云',
    tencent: '腾讯云',
    Qcloud: '腾讯云',
    aws: 'AWS',
    azure: 'Azure'
  }
  return labels[platform] || platform || '未知'
}

const getPlatformType = (platform: string) => {
  const types: Record<string, 'primary' | 'warning' | 'success' | 'info'> = {
    aliyun: 'primary',
    alibaba: 'primary',
    tencent: 'warning',
    Qcloud: 'warning',
    aws: 'success',
    azure: 'info'
  }
  return types[platform] || 'info'
}

// 获取IP地址段显示
const getIPRange = (row: Subnet) => {
  if (row.guest_ip_start && row.guest_ip_end) {
    return `${row.guest_ip_start} - ${row.guest_ip_end} /${row.guest_ip_mask}`
  }
  return row.cidr || '-'
}

// 获取使用率
const getUsagePercent = (row: Subnet) => {
  if (row.ports > 0) {
    return `${row.ports_used}/${row.ports} (${Math.round((row.ports_used / row.ports) * 100)}%)`
  }
  return '-'
}

const handleSelectionChange = (selection: Subnet[]) => {
  selectedSubnets.value = selection
}

const resetFilters = () => {
  filters.name = ''
  filters.status = ''
  filters.platform = ''
  filters.project_id = null
  filters.region = ''
  pagination.page = 1
  fetchData()
}

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      ...filters
    }
    const res = await getSubnets(params)
    ipSubnets.value = res.items || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('Failed to fetch IP subnets:', error)
    ElMessage.error('获取IP子网列表失败')
  } finally {
    loading.value = false
  }
}

const fetchCloudAccounts = async () => {
  try {
    const res = await getCloudAccounts({ page: 1, page_size: 100 })
    cloudAccounts.value = res.items || []
  } catch (error) {
    console.error('Failed to fetch cloud accounts:', error)
  }
}

const handleAccountChange = () => {
  vpcs.value = []
}

const handleView = () => {
  ElMessage.info('视图切换功能开发中')
}

const handleTags = () => {
  ElMessage.info('标签管理功能开发中')
}

const handleCreate = () => {
  Object.assign(createForm, {
    cloud_account_id: '',
    name: '',
    cidr: '',
    vpc_id: '',
    is_auto_alloc: true,
    project_id: ''
  })
  createDialogVisible.value = true
}

const confirmCreate = async () => {
  creating.value = true
  try {
    await createSubnet({
      cloud_account_id: Number(createForm.cloud_account_id),
      name: createForm.name,
      cidr: createForm.cidr,
      vpc_id: createForm.vpc_id,
      is_auto_alloc: createForm.is_auto_alloc,
      project_id: createForm.project_id ? Number(createForm.project_id) : undefined
    })
    ElMessage.success('IP子网创建成功')
    createDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('创建失败')
  } finally {
    creating.value = false
  }
}

const handleBatchCommand = async (command: string) => {
  if (selectedSubnets.value.length === 0) return
  const actionNames = { tags: '设置标签', schedtag: '调整调度标签', delete: '删除' }
  try {
    await ElMessageBox.confirm(
      `确定要批量${actionNames[command]}选中的 ${selectedSubnets.value.length} 个子网吗？`,
      '批量操作确认',
      { type: 'warning' }
    )
    if (command === 'delete') {
      const ids = selectedSubnets.value.map(s => s.id)
      await batchDeleteSubnets(ids)
      ElMessage.success('批量删除成功')
    } else {
      ElMessage.info(`批量${actionNames[command]}功能开发中`)
    }
    selectedSubnets.value = []
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('批量操作失败')
  }
}

const handleEdit = (row: Subnet) => {
  ElMessage.info('编辑功能开发中')
}

const handleAllocation = (row: Subnet) => {
  ElMessage.info('IP分配功能开发中')
}

const handleAdjustSchedtag = (row: Subnet) => {
  schedtagForm.subnet_id = row.id
  schedtagForm.subnet_name = row.name
  schedtagForm.schedtag = row.schedtag || ''
  schedtagDialogVisible.value = true
}

const confirmSchedtag = () => {
  ElMessage.success('调度标签调整成功')
  schedtagDialogVisible.value = false
  fetchData()
}

const handleDetails = (row: Subnet) => {
  selectedSubnet.value = row
  detailDialogVisible.value = true
}

const handleDelete = async (row: Subnet) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除IP子网 "${row.name}" 吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await deleteSubnet(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    console.error('Failed to delete IP subnet:', error)
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleSizeChange = (val: number) => {
  pagination.pageSize = val
  pagination.page = 1
  fetchData()
}

const handleCurrentChange = (val: number) => {
  pagination.page = val
  fetchData()
}

// Lifecycle
onMounted(() => {
  fetchData()
  fetchCloudAccounts()
})
</script>

<style scoped>
.ip-subnets-container {
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
  gap: 10px;
  align-items: center;
}

.filter-card {
  margin-bottom: 20px;
}

.table-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
}

.tag-item {
  margin-right: 4px;
  margin-bottom: 2px;
}
</style>