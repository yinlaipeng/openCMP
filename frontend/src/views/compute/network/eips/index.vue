<template>
  <div class="elastic-ips-container">
    <div class="page-header">
      <h2>弹性公网IP</h2>
      <div class="toolbar">
        <el-button link type="primary" @click="handleView">
          <el-icon><View /></el-icon>
          查看
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建
        </el-button>
        <el-dropdown trigger="click" @command="handleBatchCommand" :disabled="selectedEips.length === 0">
          <el-button :disabled="selectedEips.length === 0">
            批量操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="tags">设置标签</el-dropdown-item>
              <el-dropdown-item command="delete" divided>批量释放</el-dropdown-item>
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
          <el-input v-model="filters.name" placeholder="搜索EIP名称" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="选择状态" clearable style="width: 120px">
            <el-option label="可用" value="available" />
            <el-option label="使用中" value="in-use" />
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
        <el-form-item label="云账户">
          <el-select v-model="filters.cloud_account_id" placeholder="选择云账户" clearable style="width: 150px">
            <el-option v-for="acc in cloudAccounts" :key="acc.id" :label="acc.name" :value="acc.id" />
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
        :data="elasticIPs"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="名称" width="200" fixed="left">
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
        <el-table-column label="标签" width="80">
          <template #default="{ row }">
            <el-tag v-for="tag in (row.tags || []).slice(0, 1)" :key="tag.key" size="small" class="tag-item">
              {{ tag.key }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="address" label="IP地址" width="150" />
        <el-table-column prop="bandwidth" label="带宽(Mbps)" width="130" />
        <el-table-column prop="billing_method" label="计费方式" width="120">
          <template #default="{ row }">
            <el-tag size="small" :type="row.billing_method === 'bandwidth' ? 'primary' : 'info'">
              {{ row.billing_method === 'bandwidth' ? '按带宽' : '按流量' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="provider_type" label="平台" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="getPlatformType(row.provider_type)">
              {{ getPlatformLabel(row.provider_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="account_name" label="云账户" width="150" />
        <el-table-column prop="resource_name" label="绑定资源" width="150">
          <template #default="{ row }">
            {{ row.resource_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="project_name" label="项目" width="150" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleBind(row)" :disabled="row.status === 'in-use'">绑定</el-button>
            <el-button size="small" @click="handleUnbind(row)" :disabled="row.status !== 'in-use'">解绑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">释放</el-button>
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

    <!-- 创建弹性公网IP弹窗 -->
    <el-dialog
      title="创建弹性公网IP"
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
          <el-input v-model="createForm.name" placeholder="请输入EIP名称" />
        </el-form-item>
        <el-form-item label="带宽(Mbps)" prop="bandwidth">
          <el-input-number v-model="createForm.bandwidth" :min="1" :max="1000" style="width: 100%" />
        </el-form-item>
        <el-form-item label="计费方式">
          <el-select v-model="createForm.billing_method" placeholder="选择计费方式" style="width: 100%">
            <el-option label="按带宽" value="bandwidth" />
            <el-option label="按流量" value="traffic" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-select v-model="createForm.region_id" placeholder="选择区域" style="width: 100%">
            <el-option v-for="r in regions" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
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

    <!-- 绑定资源弹窗 -->
    <el-dialog
      title="绑定资源"
      v-model="bindDialogVisible"
      width="400px"
    >
      <el-form :model="bindForm" ref="bindFormRef" label-width="100px">
        <el-form-item label="EIP">{{ bindForm.eip_name }}</el-form-item>
        <el-form-item label="资源类型">
          <el-select v-model="bindForm.resource_type" placeholder="选择资源类型" style="width: 100%">
            <el-option label="云服务器" value="vm" />
            <el-option label="负载均衡" value="lb" />
            <el-option label="NAT网关" value="nat" />
          </el-select>
        </el-form-item>
        <el-form-item label="资源" prop="resource_id">
          <el-select v-model="bindForm.resource_id" placeholder="选择资源" style="width: 100%">
            <el-option v-for="r in bindableResources" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="bindDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmBind" :loading="binding">绑定</el-button>
      </template>
    </el-dialog>

    <!-- 详情弹窗 -->
    <el-dialog
      title="弹性公网IP详情"
      v-model="detailDialogVisible"
      width="700px"
    >
      <el-descriptions :column="2" border v-if="selectedEip">
        <el-descriptions-item label="ID">{{ selectedEip.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedEip.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(selectedEip.status)">{{ getStatusLabel(selectedEip.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ selectedEip.address }}</el-descriptions-item>
        <el-descriptions-item label="带宽">{{ selectedEip.bandwidth }} Mbps</el-descriptions-item>
        <el-descriptions-item label="计费方式">{{ selectedEip.billing_method === 'bandwidth' ? '按带宽' : '按流量' }}</el-descriptions-item>
        <el-descriptions-item label="平台">
          <el-tag size="small" :type="getPlatformType(selectedEip.provider_type)">
            {{ getPlatformLabel(selectedEip.provider_type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="云账户">{{ selectedEip.account_name }}</el-descriptions-item>
        <el-descriptions-item label="绑定资源">{{ selectedEip.resource_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedEip.project_name }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedEip.region }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Plus, ArrowDown, PriceTag } from '@element-plus/icons-vue'
import { getEIPs, createEIP, deleteEIP, bindEIP, unbindEIP, EIP } from '@/api/networkSync'
import { getCloudAccounts } from '@/api/cloud-account'
import type { CloudAccount } from '@/types'

// Data
const loading = ref(false)
const creating = ref(false)
const binding = ref(false)
const elasticIPs = ref<EIP[]>([])
const selectedEips = ref<EIP[]>([])
const cloudAccounts = ref<CloudAccount[]>([])
const projects = ref<any[]>([])
const regions = ref<any[]>([])
const bindableResources = ref<any[]>([])

const createDialogVisible = ref(false)
const bindDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const selectedEip = ref<EIP | null>(null)

const filters = reactive({
  name: '',
  status: '',
  platform: '',
  cloud_account_id: null as number | null,
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
  bandwidth: 10,
  billing_method: 'bandwidth',
  region_id: '',
  project_id: ''
})

const bindForm = reactive({
  eip_id: '',
  eip_name: '',
  resource_type: 'vm',
  resource_id: ''
})

const createRules = {
  cloud_account_id: [{ required: true, message: '请选择云账户', trigger: 'change' }],
  name: [{ required: true, message: '请输入EIP名称', trigger: 'blur' }]
}

// Methods
const getStatusType = (status: string) => {
  switch (status) {
    case 'available': return 'success'
    case 'in-use': return 'primary'
    case 'creating': return 'warning'
    case 'error': return 'danger'
    default: return 'info'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'available': return '可用'
    case 'in-use': return '使用中'
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

const handleSelectionChange = (selection: EIP[]) => {
  selectedEips.value = selection
}

const resetFilters = () => {
  filters.name = ''
  filters.status = ''
  filters.platform = ''
  filters.cloud_account_id = null
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
    const res = await getEIPs(params)
    elasticIPs.value = res.items || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('Failed to fetch elastic IPs:', error)
    ElMessage.error('获取弹性公网IP列表失败')
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
  regions.value = []
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
    bandwidth: 10,
    billing_method: 'bandwidth',
    region_id: '',
    project_id: ''
  })
  createDialogVisible.value = true
}

const confirmCreate = async () => {
  creating.value = true
  try {
    await createEIP({
      cloud_account_id: Number(createForm.cloud_account_id),
      name: createForm.name,
      bandwidth: createForm.bandwidth,
      billing_method: createForm.billing_method,
      region_id: createForm.region_id,
      project_id: createForm.project_id ? Number(createForm.project_id) : undefined
    })
    ElMessage.success('弹性公网IP创建成功')
    createDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('创建失败')
  } finally {
    creating.value = false
  }
}

const handleBatchCommand = async (command: string) => {
  if (selectedEips.value.length === 0) return
  const actionNames = { tags: '设置标签', delete: '释放' }
  try {
    await ElMessageBox.confirm(
      `确定要批量${actionNames[command]}选中的 ${selectedEips.value.length} 个EIP吗？`,
      '批量操作确认',
      { type: 'warning' }
    )
    if (command === 'delete') {
      for (const eip of selectedEips.value) {
        await deleteEIP(eip.id)
      }
      ElMessage.success('批量释放成功')
    } else {
      ElMessage.info(`批量${actionNames[command]}功能开发中`)
    }
    selectedEips.value = []
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('批量操作失败')
  }
}

const handleEdit = (row: EIP) => {
  ElMessage.info('编辑功能开发中')
}

const handleBind = (row: EIP) => {
  Object.assign(bindForm, {
    eip_id: row.id,
    eip_name: row.name,
    resource_type: 'vm',
    resource_id: ''
  })
  bindableResources.value = []
  bindDialogVisible.value = true
}

const confirmBind = async () => {
  binding.value = true
  try {
    await bindEIP(bindForm.eip_id, bindForm.resource_id)
    ElMessage.success('绑定成功')
    bindDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('绑定失败')
  } finally {
    binding.value = false
  }
}

const handleUnbind = async (row: EIP) => {
  try {
    await ElMessageBox.confirm(
      `确定要解绑弹性公网IP "${row.name}" 吗？`,
      '确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      }
    )

    await unbindEIP(row.id)
    ElMessage.success('解绑成功')
    fetchData()
  } catch (error) {
    console.error('Failed to unbind EIP:', error)
    if (error !== 'cancel') {
      ElMessage.error('解绑失败')
    }
  }
}

const handleDetails = (row: EIP) => {
  selectedEip.value = row
  detailDialogVisible.value = true
}

const handleDelete = async (row: EIP) => {
  try {
    await ElMessageBox.confirm(
      `确定要释放弹性公网IP "${row.name}" 吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await deleteEIP(row.id)
    ElMessage.success('释放成功')
    fetchData()
  } catch (error) {
    console.error('Failed to release elastic IP:', error)
    if (error !== 'cancel') {
      ElMessage.error('释放失败')
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
.elastic-ips-container {
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