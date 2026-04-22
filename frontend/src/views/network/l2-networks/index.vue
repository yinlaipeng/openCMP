<template>
  <div class="l2-networks-container">
    <div class="page-header">
      <h2>二层网络</h2>
      <div class="toolbar">
        <el-button link type="primary" @click="handleView">
          <el-icon><View /></el-icon>
          查看
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建
        </el-button>
        <el-dropdown trigger="click" @command="handleBatchCommand" :disabled="selectedNetworks.length === 0">
          <el-button :disabled="selectedNetworks.length === 0">
            批量操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="tags">设置标签</el-dropdown-item>
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
          <el-input v-model="filters.name" placeholder="搜索二层网络名称" clearable style="width: 180px" />
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
        :data="l2Networks"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <!-- Cloudpods 表头顺序: Name → Status → Tags → Bandwidth → VPC → Networks → Platform → Owner Domain → Region -->
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
        <el-table-column label="标签" width="100">
          <template #default="{ row }">
            <el-tag v-for="tag in (row.tags || []).slice(0, 2)" :key="tag.key" size="small" class="tag-item">
              {{ tag.key }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="bandwidth" label="带宽" width="100">
          <template #default="{ row }">
            {{ row.bandwidth }} Mbps
          </template>
        </el-table-column>
        <el-table-column prop="vpc_name" label="VPC" width="150" />
        <el-table-column prop="network_count" label="网络数" width="100" />
        <el-table-column prop="provider_type" label="平台" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="getPlatformType(row.provider_type)">
              {{ getPlatformLabel(row.provider_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="owner_domain" label="所属域" width="150" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-dropdown trigger="click">
              <el-button size="small">
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleChangeDomain(row)">更改域</el-dropdown-item>
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

    <!-- 创建二层网络弹窗 -->
    <el-dialog
      title="创建二层网络"
      v-model="createDialogVisible"
      width="500px"
    >
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入二层网络名称" />
        </el-form-item>
        <el-form-item label="VPC" prop="vpc_id">
          <el-select v-model="createForm.vpc_id" placeholder="选择VPC" style="width: 100%">
            <el-option v-for="vpc in vpcs" :key="vpc.id" :label="vpc.name" :value="vpc.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="带宽">
          <el-input-number v-model="createForm.bandwidth" :min="1" :max="10000" style="width: 100%" />
        </el-form-item>
        <el-form-item label="MTU">
          <el-input-number v-model="createForm.mtu" :min="64" :max="9000" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCreate" :loading="creating">创建</el-button>
      </template>
    </el-dialog>

    <!-- 编辑弹窗 -->
    <el-dialog
      title="编辑二层网络"
      v-model="editDialogVisible"
      width="400px"
    >
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="带宽">
          <el-input-number v-model="editForm.bandwidth" :min="1" :max="10000" style="width: 100%" />
        </el-form-item>
        <el-form-item label="MTU">
          <el-input-number v-model="editForm.mtu" :min="64" :max="9000" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmEdit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 详情弹窗 -->
    <el-dialog
      title="二层网络详情"
      v-model="detailDialogVisible"
      width="700px"
    >
      <el-descriptions :column="2" border v-if="selectedNetwork">
        <el-descriptions-item label="ID">{{ selectedNetwork.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedNetwork.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(selectedNetwork.status)">{{ getStatusLabel(selectedNetwork.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="带宽">{{ selectedNetwork.bandwidth }} Mbps</el-descriptions-item>
        <el-descriptions-item label="VPC">{{ selectedNetwork.vpc_name }}</el-descriptions-item>
        <el-descriptions-item label="网络数">{{ selectedNetwork.network_count }}</el-descriptions-item>
        <el-descriptions-item label="VLAN ID">{{ selectedNetwork.vlan_id || '-' }}</el-descriptions-item>
        <el-descriptions-item label="MTU">{{ selectedNetwork.mtu || 1500 }}</el-descriptions-item>
        <el-descriptions-item label="平台">
          <el-tag size="small" :type="getPlatformType(selectedNetwork.provider_type)">
            {{ getPlatformLabel(selectedNetwork.provider_type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="所属域">{{ selectedNetwork.owner_domain }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedNetwork.region }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Plus, ArrowDown, PriceTag } from '@element-plus/icons-vue'
import { getL2Networks, createL2Network, deleteL2Network, batchDeleteL2Networks, L2Network } from '@/api/networkSync'

// Data
const loading = ref(false)
const creating = ref(false)
const l2Networks = ref<L2Network[]>([])
const selectedNetworks = ref<L2Network[]>([])
const vpcs = ref<any[]>([])

const createDialogVisible = ref(false)
const editDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const selectedNetwork = ref<L2Network | null>(null)

const filters = reactive({
  name: '',
  status: '',
  platform: '',
  region: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const createForm = reactive({
  name: '',
  vpc_id: '',
  bandwidth: 1000,
  mtu: 1500
})

const editForm = reactive({
  id: '',
  name: '',
  bandwidth: 1000,
  mtu: 1500
})

const createRules = {
  name: [{ required: true, message: '请输入二层网络名称', trigger: 'blur' }],
  vpc_id: [{ required: true, message: '请选择VPC', trigger: 'change' }]
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

const handleSelectionChange = (selection: L2Network[]) => {
  selectedNetworks.value = selection
}

const resetFilters = () => {
  filters.name = ''
  filters.status = ''
  filters.platform = ''
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
    const res = await getL2Networks(params)
    l2Networks.value = res.items || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('Failed to fetch L2 networks:', error)
    ElMessage.error('获取二层网络列表失败')
  } finally {
    loading.value = false
  }
}

const handleView = () => {
  ElMessage.info('视图切换功能开发中')
}

const handleTags = () => {
  ElMessage.info('标签管理功能开发中')
}

const handleCreate = () => {
  Object.assign(createForm, {
    name: '',
    vpc_id: '',
    bandwidth: 1000,
    mtu: 1500
  })
  createDialogVisible.value = true
}

const confirmCreate = async () => {
  creating.value = true
  try {
    await createL2Network({
      name: createForm.name,
      vpc_id: createForm.vpc_id,
      bandwidth: createForm.bandwidth,
      mtu: createForm.mtu
    })
    ElMessage.success('二层网络创建成功')
    createDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('创建失败')
  } finally {
    creating.value = false
  }
}

const handleBatchCommand = async (command: string) => {
  if (selectedNetworks.value.length === 0) return
  const actionNames = { tags: '设置标签', delete: '删除' }
  try {
    await ElMessageBox.confirm(
      `确定要批量${actionNames[command]}选中的 ${selectedNetworks.value.length} 个二层网络吗？`,
      '批量操作确认',
      { type: 'warning' }
    )
    if (command === 'delete') {
      const ids = selectedNetworks.value.map(n => n.id)
      await batchDeleteL2Networks(ids)
      ElMessage.success('批量删除成功')
    } else {
      ElMessage.info(`批量${actionNames[command]}功能开发中`)
    }
    selectedNetworks.value = []
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('批量操作失败')
  }
}

const handleEdit = (row: L2Network) => {
  editForm.id = row.id
  editForm.name = row.name
  editForm.bandwidth = row.bandwidth
  editForm.mtu = row.mtu || 1500
  editDialogVisible.value = true
}

const confirmEdit = () => {
  ElMessage.success('编辑成功')
  editDialogVisible.value = false
  fetchData()
}

const handleChangeDomain = (row: L2Network) => {
  ElMessage.info('更改域功能开发中')
}

const handleDetails = (row: L2Network) => {
  selectedNetwork.value = row
  detailDialogVisible.value = true
}

const handleDelete = async (row: L2Network) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除二层网络 "${row.name}" 吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await deleteL2Network(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    console.error('Failed to delete L2 network:', error)
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
})
</script>

<style scoped>
.l2-networks-container {
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