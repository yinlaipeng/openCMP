<template>
  <div class="subnets-container">
    <div class="page-header">
      <h2>IP子网管理</h2>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        创建子网
      </el-button>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="loadSubnets">
        <el-form-item label="云账号">
          <CloudAccountSelector
            v-model:value="filters.account_id"
            placeholder="选择云账号"
            @change="handleAccountChange"
          />
        </el-form-item>
        <el-form-item label="VPC">
          <el-input v-model="filters.vpc_id" placeholder="VPC ID" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="状态" clearable>
            <el-option label="可用" value="Available" />
            <el-option label="创建中" value="Creating" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadSubnets">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table
      :data="subnets"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
    >
      <el-table-column label="名称" min-width="180">
        <template #default="{ row }">
          <el-link type="primary" @click="viewDetail(row)">{{ row.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column label="平台/云账号" width="150">
        <template #default="{ row }">
          <div class="platform-cell">
            <el-tag size="small" :type="getPlatformType(row.platform)" effect="plain">
              {{ getPlatformLabel(row.platform) }}
            </el-tag>
            <span class="account-name">{{ row.account_name || '-' }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="cidr" label="CIDR" width="180">
        <template #default="{ row }">
          <span class="font-mono">{{ row.cidr }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="zone_id" label="可用区" width="120">
        <template #default="{ row }">
          <span class="font-mono">{{ row.zone_id }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="vpc_id" label="VPC" width="180">
        <template #default="{ row }">
          <span class="font-mono">{{ row.vpc_id }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="region_id" label="区域" width="100">
        <template #default="{ row }">
          <span class="font-mono">{{ row.region_id }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-dropdown trigger="click">
            <el-button size="small">操作<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="syncStatus(row)">同步状态</el-dropdown-item>
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
      :total="pagination.total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      class="pagination"
    />

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="IP子网详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedSubnet?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedSubnet?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedSubnet?.status }}</el-descriptions-item>
        <el-descriptions-item label="CIDR">{{ selectedSubnet?.cidr }}</el-descriptions-item>
        <el-descriptions-item label="可用区">{{ selectedSubnet?.zone_id }}</el-descriptions-item>
        <el-descriptions-item label="VPC">{{ selectedSubnet?.vpc_id }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedSubnet?.region_id }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedSubnet?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 创建子网模态框 -->
    <CreateSubnetModal
      v-model:visible="createModalVisible"
      @success="handleCreateSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown, Plus } from '@element-plus/icons-vue'
import CreateSubnetModal from '@/components/network/CreateSubnetModal.vue'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import { getSubnets, deleteSubnet } from '@/api/network'
import type { Subnet } from '@/types'

interface ExtendedSubnet extends Subnet {
  platform?: string
  account_name?: string
}

const subnets = ref<ExtendedSubnet[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const createModalVisible = ref(false)
const selectedSubnet = ref<ExtendedSubnet | null>(null)

// 筛选条件
const filters = reactive({
  account_id: null as number | null,
  vpc_id: '',
  status: ''
})

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 平台类型映射
const platformLabels: Record<string, string> = {
  alibaba: '阿里云',
  tencent: '腾讯云',
  aws: 'AWS',
  azure: 'Azure'
}

const platformTypes: Record<string, 'primary' | 'warning' | 'success' | 'info'> = {
  alibaba: 'primary',
  tencent: 'warning',
  aws: 'success',
  azure: 'info'
}

const getPlatformLabel = (platform: string): string => {
  return platformLabels[platform] || platform || '未知'
}

const getPlatformType = (platform: string): 'primary' | 'warning' | 'success' | 'info' => {
  return platformTypes[platform] || 'info'
}

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'available':
    case 'active':
      return 'success'
    case 'pending':
    case 'creating':
      return 'warning'
    case 'failed':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const handleAccountChange = (accountId: number | null) => {
  filters.account_id = accountId
  if (accountId) {
    loadSubnets()
  } else {
    subnets.value = []
  }
}

const loadSubnets = async () => {
  if (!filters.account_id) {
    subnets.value = []
    return
  }

  loading.value = true
  try {
    const params: any = {
      account_id: filters.account_id,
      page: pagination.page,
      size: pagination.pageSize
    }
    if (filters.vpc_id) params.vpc_id = filters.vpc_id
    if (filters.status) params.status = filters.status

    const res = await getSubnets(params)
    subnets.value = Array.isArray(res) ? res : res.items || []
    pagination.total = res.total || subnets.value.length
  } catch (error: any) {
    console.error('Failed to load subnets:', error)
    subnets.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: ExtendedSubnet) => {
  selectedSubnet.value = row
  detailDialogVisible.value = true
}

const syncStatus = async (row: ExtendedSubnet) => {
  ElMessage.info(`正在同步子网 ${row.name} 的状态`)
  await loadSubnets()
}

const handleDelete = async (row: ExtendedSubnet) => {
  try {
    await ElMessageBox.confirm(`确认删除IP子网 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    if (filters.account_id) {
      await deleteSubnet(row.id, filters.account_id)
      ElMessage.success('删除成功')
      loadSubnets()
    }
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
      ElMessage.error(`删除失败: ${e.message}`)
    }
  }
}

const handleCreate = () => {
  createModalVisible.value = true
}

const handleCreateSuccess = (subnet: ExtendedSubnet) => {
  ElMessage.success(`${subnet.name} 创建成功`)
  loadSubnets()
}

const resetFilters = () => {
  filters.account_id = null
  filters.vpc_id = ''
  filters.status = ''
  pagination.page = 1
  subnets.value = []
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadSubnets()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadSubnets()
}

onMounted(() => {
  // Wait for account selector to initialize
})
</script>

<style scoped>
.subnets-container {
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

.platform-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.account-name {
  font-size: 12px;
  color: var(--text-color-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.font-mono {
  font-family: monospace;
}
</style>