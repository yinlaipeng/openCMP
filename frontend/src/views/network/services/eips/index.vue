<template>
  <div class="eips-container">
    <div class="page-header">
      <h2>弹性公网IP管理</h2>
      <el-button type="primary" @click="handleCreate">申请弹性IP</el-button>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="loadEips">
        <el-form-item label="云账号">
          <CloudAccountSelector
            v-model:value="filters.account_id"
            placeholder="选择云账号"
            @change="handleAccountChange"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="选择状态" clearable>
            <el-option label="已绑定" value="In-Use" />
            <el-option label="未绑定" value="Available" />
          </el-select>
        </el-form-item>
        <el-form-item label="IP地址">
          <el-input v-model="filters.ip" placeholder="IP地址搜索" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadEips">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table
      :data="eips"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
    >
      <el-table-column label="名称" min-width="150">
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
      <el-table-column prop="address" label="IP地址" width="150" />
      <el-table-column prop="bandwidth" label="带宽" width="100">
        <template #default="{ row }">
          {{ row.bandwidth }} Mbps
        </template>
      </el-table-column>
      <el-table-column prop="resource_id" label="绑定资源" width="150">
        <template #default="{ row }">
          {{ row.resource_id || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="region_id" label="区域" width="120" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="handleBind(row)" v-if="row.status === 'Available'">绑定</el-button>
          <el-button size="small" @click="handleUnbind(row)" v-if="row.status === 'In-Use'">解绑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">释放</el-button>
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
    <el-dialog v-model="detailVisible" title="弹性公网IP详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedEip?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedEip?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedEip?.status }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ selectedEip?.address }}</el-descriptions-item>
        <el-descriptions-item label="带宽">{{ selectedEip?.bandwidth }} Mbps</el-descriptions-item>
        <el-descriptions-item label="绑定资源">{{ selectedEip?.resource_id || '-' }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedEip?.region_id }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedEip?.account_name }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Bind Modal -->
    <el-dialog v-model="bindVisible" title="绑定弹性IP" width="400px">
      <el-form :model="bindForm" label-width="100px">
        <el-form-item label="资源类型">
          <el-select v-model="bindForm.resource_type" placeholder="选择资源类型">
            <el-option label="ECS实例" value="ecs" />
            <el-option label="SLB实例" value="slb" />
          </el-select>
        </el-form-item>
        <el-form-item label="资源ID">
          <el-input v-model="bindForm.resource_id" placeholder="输入资源ID" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="bindVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmBind">确定</el-button>
      </template>
    </el-dialog>

    <!-- Create Modal -->
    <el-dialog v-model="createVisible" title="申请弹性IP" width="400px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="带宽(Mbps)">
          <el-input-number v-model="createForm.bandwidth" :min="1" :max="1000" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCreate">申请</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import { getEIPs, createEIP, bindEIP, unbindEIP, deleteEIP } from '@/api/network'

interface EIP {
  id: string
  name: string
  status: string
  address: string
  bandwidth: number
  resource_id: string
  resource_type: string
  region_id: string
  platform: string
  account_name: string
  cloud_account_id: number
}

const eips = ref<EIP[]>([])
const loading = ref(false)
const detailVisible = ref(false)
const bindVisible = ref(false)
const createVisible = ref(false)
const selectedEip = ref<EIP | null>(null)

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 筛选条件
const filters = reactive({
  account_id: null as number | null,
  status: '',
  ip: ''
})

// 绑定表单
const bindForm = reactive({
  eip_id: '',
  resource_type: 'ecs',
  resource_id: ''
})

// 创建表单
const createForm = reactive({
  account_id: 0,
  bandwidth: 100
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
  if (status === 'Available' || status === 'available') return 'success'
  if (status === 'In-Use' || status === 'in-use') return 'primary'
  if (status === 'Associating' || status === 'Unassociating') return 'warning'
  if (status === 'Error' || status === 'error') return 'danger'
  return 'info'
}

const handleAccountChange = (accountId: number | null) => {
  filters.account_id = accountId
}

const loadEips = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      size: pagination.pageSize
    }
    if (filters.account_id) params.account_id = filters.account_id
    if (filters.status) params.status = filters.status
    if (filters.ip) params.ip = filters.ip

    const res = await getEIPs(params)
    eips.value = res.items || res
    pagination.total = res.total || eips.value.length
  } catch (e) {
    console.error(e)
    ElMessage.error('加载弹性IP列表失败')
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: EIP) => {
  selectedEip.value = row
  detailVisible.value = true
}

const handleCreate = () => {
  if (!filters.account_id) {
    ElMessage.warning('请先选择云账号')
    return
  }
  createForm.account_id = filters.account_id
  createVisible.value = true
}

const confirmCreate = async () => {
  try {
    await createEIP(createForm.account_id, { bandwidth: createForm.bandwidth })
    ElMessage.success('申请成功')
    createVisible.value = false
    loadEips()
  } catch (e) {
    ElMessage.error('申请失败')
  }
}

const handleBind = (row: EIP) => {
  selectedEip.value = row
  bindForm.eip_id = row.id
  bindForm.resource_type = 'ecs'
  bindForm.resource_id = ''
  bindVisible.value = true
}

const confirmBind = async () => {
  try {
    await bindEIP(filters.account_id!, bindForm.eip_id, bindForm.resource_id, bindForm.resource_type)
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
    await unbindEIP(filters.account_id!, row.id)
    ElMessage.success('解绑成功')
    loadEips()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('解绑失败')
  }
}

const handleDelete = async (row: EIP) => {
  try {
    await ElMessageBox.confirm(`确认释放弹性公网IP ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteEIP(filters.account_id!, row.id)
    ElMessage.success('释放成功')
    loadEips()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('释放失败')
  }
}

const resetFilters = () => {
  filters.account_id = null
  filters.status = ''
  filters.ip = ''
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
})
</script>

<style scoped>
.eips-container {
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
</style>