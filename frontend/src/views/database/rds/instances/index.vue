<template>
  <div class="rds-container">
    <div class="page-header">
      <h2>RDS实例</h2>
      <div class="toolbar">
        <el-button @click="handleRefresh" :icon="Refresh" circle />
        <el-button type="primary" @click="showCreateDialog">新建</el-button>
        <el-button @click="handleSyncStatus" :disabled="selectedRows.length === 0">同步状态</el-button>
        <el-dropdown @command="handleBatchCommand" :disabled="selectedRows.length === 0">
          <el-button>批量操作<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="start">批量启动</el-dropdown-item>
              <el-dropdown-item command="stop">批量停止</el-dropdown-item>
              <el-dropdown-item command="delete">批量删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button @click="handleTags" :icon="PriceTag">标签</el-button>
      </div>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="loadRDSInstances">
        <el-form-item label="云账号">
          <CloudAccountSelector
            v-model:value="filters.account_id"
            placeholder="选择云账号"
            @change="handleAccountChange"
          />
        </el-form-item>
        <el-form-item label="数据库引擎">
          <el-select v-model="filters.engine" placeholder="数据库引擎" clearable>
            <el-option label="MySQL" value="MySQL" />
            <el-option label="PostgreSQL" value="PostgreSQL" />
            <el-option label="SQLServer" value="SQLServer" />
            <el-option label="MariaDB" value="MariaDB" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="状态" clearable>
            <el-option label="运行中" value="Running" />
            <el-option label="已停止" value="Stopped" />
            <el-option label="创建中" value="Creating" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadRDSInstances">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table
      :data="rdsInstances"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column label="名称" min-width="150">
        <template #default="{ row }">
          <el-link type="primary" @click="viewDetail(row)">{{ row.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="instance_type" label="类型" width="120" />
      <el-table-column label="引擎" width="100">
        <template #default="{ row }">
          {{ row.engine }}
        </template>
      </el-table-column>
      <el-table-column prop="endpoint" label="地址" width="180" />
      <el-table-column prop="port" label="端口" width="80" />
      <el-table-column prop="storage_type" label="存储类型" width="100" />
      <el-table-column prop="security_group" label="安全组" width="120" />
      <el-table-column prop="billing_type" label="计费类型" width="100">
        <template #default="{ row }">
          {{ row.billing_type || '按量付费' }}
        </template>
      </el-table-column>
      <el-table-column label="平台" width="100">
        <template #default="{ row }">
          <el-tag size="small" :type="getPlatformType(row.platform)" effect="plain">
            {{ getPlatformLabel(row.platform) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="project" label="项目" width="100" />
      <el-table-column prop="region" label="区域" width="100" />
      <el-table-column label="操作" width="140" fixed="right">
        <template #default="{ row }">
          <el-dropdown trigger="click">
            <el-button size="small">操作<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleAction(row, 'start')">启动</el-dropdown-item>
                <el-dropdown-item @click="handleAction(row, 'stop')">停止</el-dropdown-item>
                <el-dropdown-item @click="handleAction(row, 'reboot')">重启</el-dropdown-item>
                <el-dropdown-item divided @click="showResizeDialog(row)">调整规格</el-dropdown-item>
                <el-dropdown-item @click="showBackupDialog(row)">创建备份</el-dropdown-item>
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
    <el-dialog v-model="detailDialogVisible" title="RDS实例详情" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedRDS?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedRDS?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedRDS?.status }}</el-descriptions-item>
        <el-descriptions-item label="规格">{{ selectedRDS?.instance_type }}</el-descriptions-item>
        <el-descriptions-item label="数据库引擎">{{ selectedRDS?.engine }} {{ selectedRDS?.engine_version }}</el-descriptions-item>
        <el-descriptions-item label="连接地址">{{ selectedRDS?.endpoint }}</el-descriptions-item>
        <el-descriptions-item label="端口">{{ selectedRDS?.port }}</el-descriptions-item>
        <el-descriptions-item label="存储类型">{{ selectedRDS?.storage_type }}</el-descriptions-item>
        <el-descriptions-item label="存储大小">{{ selectedRDS?.storage_size }} GB</el-descriptions-item>
        <el-descriptions-item label="安全组">{{ selectedRDS?.security_group }}</el-descriptions-item>
        <el-descriptions-item label="计费类型">{{ selectedRDS?.billing_type }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedRDS?.project }}</el-descriptions-item>
        <el-descriptions-item label="VPC">{{ selectedRDS?.vpc_id }}</el-descriptions-item>
        <el-descriptions-item label="子网">{{ selectedRDS?.subnet_id }}</el-descriptions-item>
        <el-descriptions-item label="可用区">{{ selectedRDS?.zone_id }}</el-descriptions-item>
        <el-descriptions-item label="主账号">{{ selectedRDS?.master_username }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedRDS?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Create Dialog -->
    <el-dialog v-model="createDialogVisible" title="创建RDS实例" width="700px">
      <el-form :model="createForm" label-width="120px">
        <el-form-item label="指定项目">
          <el-select v-model="createForm.project_id" placeholder="选择项目" clearable>
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="实例名称" required>
          <el-input v-model="createForm.name" placeholder="请输入实例名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="createForm.description" type="textarea" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="计费类型">
          <el-radio-group v-model="createForm.billing_type">
            <el-radio label="postpay">按量付费</el-radio>
            <el-radio label="prepay">包年包月</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="过期释放">
          <el-radio-group v-model="createForm.auto_release">
            <el-radio label="false">不自动释放</el-radio>
            <el-radio label="true">自动释放</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="创建数量">
          <el-input-number v-model="createForm.quantity" :min="1" :max="10" />
        </el-form-item>
        <el-form-item label="区域">
          <el-select v-model="createForm.region_id" placeholder="选择区域" @change="loadRDSSKUs">
            <el-option label="华东1(杭州)" value="cn-hangzhou" />
            <el-option label="华东2(上海)" value="cn-shanghai" />
            <el-option label="华北2(北京)" value="cn-beijing" />
            <el-option label="华南1(深圳)" value="cn-shenzhen" />
          </el-select>
        </el-form-item>
        <el-form-item label="数据库引擎" required>
          <el-select v-model="createForm.engine" placeholder="选择数据库引擎" @change="loadRDSSKUs">
            <el-option label="MySQL" value="MySQL" />
            <el-option label="PostgreSQL" value="PostgreSQL" />
            <el-option label="SQLServer" value="SQLServer" />
            <el-option label="MariaDB" value="MariaDB" />
          </el-select>
        </el-form-item>
        <el-form-item label="版本" required>
          <el-select v-model="createForm.engine_version" placeholder="选择版本" @change="loadRDSSKUs">
            <el-option label="MySQL 5.7" value="5.7" v-if="createForm.engine === 'MySQL'" />
            <el-option label="MySQL 8.0" value="8.0" v-if="createForm.engine === 'MySQL'" />
            <el-option label="PostgreSQL 12" value="12" v-if="createForm.engine === 'PostgreSQL'" />
            <el-option label="PostgreSQL 14" value="14" v-if="createForm.engine === 'PostgreSQL'" />
          </el-select>
        </el-form-item>
        <el-form-item label="实例类型">
          <el-radio-group v-model="createForm.category">
            <el-radio label="ha">高可用版</el-radio>
            <el-radio label="basic">基础版</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="存储类型">
          <el-select v-model="createForm.storage_type" placeholder="选择存储类型" @change="loadRDSSKUs">
            <el-option label="SSD云盘" value="cloud_ssd" />
            <el-option label="ESSD云盘" value="cloud_essd" />
          </el-select>
        </el-form-item>
        <el-form-item label="CPU">
          <el-select v-model="createForm.cpu" placeholder="选择CPU核数" @change="filterSKUsByCPU">
            <el-option label="1核" :value="1" />
            <el-option label="2核" :value="2" />
            <el-option label="4核" :value="4" />
            <el-option label="8核" :value="8" />
          </el-select>
        </el-form-item>
        <el-form-item label="内存">
          <el-select v-model="createForm.memory_mb" placeholder="选择内存大小">
            <el-option label="1GB" :value="1024" />
            <el-option label="2GB" :value="2048" />
            <el-option label="4GB" :value="4096" />
            <el-option label="8GB" :value="8192" />
            <el-option label="16GB" :value="16384" />
          </el-select>
        </el-form-item>
        <el-form-item label="实例规格" required>
          <el-select v-model="createForm.instance_type" placeholder="选择实例规格">
            <el-option v-for="sku in availableSkus" :key="sku.id"
              :label="`${sku.instance_type} (${sku.cpu}核${sku.memory_mb/1024}GB)`"
              :value="sku.instance_type" />
          </el-select>
        </el-form-item>
        <el-form-item label="存储大小(GB)" required>
          <el-input-number v-model="createForm.storage_size" :min="10" :max="2000" />
        </el-form-item>
        <el-form-item label="VPC" required>
          <el-input v-model="createForm.vpc_id" placeholder="VPC ID" />
        </el-form-item>
        <el-form-item label="子网" required>
          <el-input v-model="createForm.subnet_id" placeholder="子网 ID" />
        </el-form-item>
        <el-form-item label="可用区">
          <el-input v-model="createForm.zone_id" placeholder="可用区 ID" />
        </el-form-item>
        <el-form-item label="主账号用户名">
          <el-input v-model="createForm.master_username" placeholder="如: root, admin" />
        </el-form-item>
        <el-form-item label="主账号密码">
          <el-input v-model="createForm.master_password" type="password" placeholder="设置密码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="createLoading">创建</el-button>
      </template>
    </el-dialog>

    <!-- Resize Dialog -->
    <el-dialog v-model="resizeDialogVisible" title="调整RDS规格" width="400px">
      <el-form :model="resizeForm" label-width="100px">
        <el-form-item label="新规格">
          <el-input v-model="resizeForm.instance_type" placeholder="如: mysql.n4.medium.1c" />
        </el-form-item>
        <el-form-item label="新存储大小">
          <el-input-number v-model="resizeForm.storage_size" :min="10" :max="2000" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resizeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleResize" :loading="resizeLoading">确认调整</el-button>
      </template>
    </el-dialog>

    <!-- Backup Dialog -->
    <el-dialog v-model="backupDialogVisible" title="创建备份" width="300px">
      <el-form :model="backupForm" label-width="80px">
        <el-form-item label="备份名称">
          <el-input v-model="backupForm.name" placeholder="可选" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="backupDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleBackup" :loading="backupLoading">创建备份</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown, Refresh, PriceTag } from '@element-plus/icons-vue'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import {
  listRDS, createRDS, deleteRDS, rdsAction, resizeRDS, createRDSBackup,
  listRDSSKUs, type RDSInstance, type RDSConfig, type RDSInstanceSKU
} from '@/api/database'

interface RDS extends RDSInstance {
  platform?: string
  account_name?: string
  billing_type?: string
  security_group?: string
  project?: string
  region?: string
}

interface Project {
  id: number
  name: string
}

const rdsInstances = ref<RDS[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedRDS = ref<RDS | null>(null)
const selectedRows = ref<RDS[]>([])

// Projects list (should be loaded from API)
const projects = ref<Project[]>([
  { id: 1, name: '默认项目' },
  { id: 2, name: '生产环境' },
  { id: 3, name: '测试环境' },
])

// SKU data
const availableSkus = ref<RDSInstanceSKU[]>([])

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 筛选条件
const filters = reactive({
  account_id: null as number | null,
  engine: '',
  status: ''
})

// Create dialog
const createDialogVisible = ref(false)
const createLoading = ref(false)
const createForm = ref<RDSConfig & {
  project_id?: number
  description?: string
  billing_type?: string
  auto_release?: string
  quantity?: number
  region_id?: string
  category?: string
  cpu?: number
  memory_mb?: number
}>({
  account_id: 0,
  name: '',
  engine: 'MySQL',
  engine_version: '5.7',
  instance_type: '',
  storage_size: 100,
  storage_type: 'cloud_essd',
  vpc_id: '',
  subnet_id: '',
  zone_id: '',
  master_username: 'root',
  master_password: '',
  billing_type: 'postpay',
  auto_release: 'false',
  quantity: 1,
  category: 'ha',
  cpu: 2,
  memory_mb: 4096,
})

// Resize dialog
const resizeDialogVisible = ref(false)
const resizeLoading = ref(false)
const resizeForm = ref({
  instance_id: '',
  instance_type: '',
  storage_size: 0
})

// Backup dialog
const backupDialogVisible = ref(false)
const backupLoading = ref(false)
const backupForm = ref({
  instance_id: '',
  name: ''
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

const statusLabels: Record<string, string> = {
  Running: '运行中',
  Stopped: '已停止',
  Creating: '创建中',
  Starting: '启动中',
  Stopping: '停止中'
}

const getPlatformLabel = (platform: string): string => {
  return platformLabels[platform] || platform || '未知'
}

const getPlatformType = (platform: string): 'primary' | 'warning' | 'success' | 'info' => {
  return platformTypes[platform] || 'info'
}

const getStatusLabel = (status: string): string => {
  return statusLabels[status] || status
}

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'running':
    case 'active':
      return 'success'
    case 'creating':
    case 'pending':
    case 'starting':
    case 'stopping':
      return 'warning'
    case 'stopped':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const handleSelectionChange = (rows: RDS[]) => {
  selectedRows.value = rows
}

const handleAccountChange = (accountId: number | null) => {
  filters.account_id = accountId
  if (accountId) {
    loadRDSInstances()
  } else {
    rdsInstances.value = []
  }
}

const handleRefresh = () => {
  loadRDSInstances()
}

const handleSyncStatus = async () => {
  if (selectedRows.value.length === 0) return
  ElMessage.info(`正在同步 ${selectedRows.value.length} 个实例的状态...`)
  await loadRDSInstances()
  ElMessage.success('状态同步完成')
}

const handleBatchCommand = async (command: string) => {
  if (selectedRows.value.length === 0) return

  try {
    await ElMessageBox.confirm(`确认对 ${selectedRows.value.length} 个实例执行 ${command} 操作？`, '确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    for (const row of selectedRows.value) {
      if (command === 'delete') {
        await deleteRDS(filters.account_id!, row.id)
      } else {
        await rdsAction(filters.account_id!, row.id, command)
      }
    }
    ElMessage.success('批量操作执行成功')
    selectedRows.value = []
    loadRDSInstances()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '操作失败')
    }
  }
}

const handleTags = () => {
  ElMessage.info('标签管理功能开发中')
}

const loadRDSSKUs = async () => {
  if (!filters.account_id) return

  try {
    const skus = await listRDSSKUs({
      account_id: filters.account_id,
      engine: createForm.value.engine,
      engine_version: createForm.value.engine_version,
      category: createForm.value.category,
      storage_type: createForm.value.storage_type,
      region_id: createForm.value.region_id,
    })
    availableSkus.value = skus
  } catch (e) {
    console.error('Failed to load SKUs', e)
  }
}

const filterSKUsByCPU = () => {
  // Filter available SKUs based on CPU selection
  if (createForm.value.cpu) {
    availableSkus.value = availableSkus.value.filter(s => s.cpu === createForm.value.cpu)
  }
}

const loadRDSInstances = async () => {
  if (!filters.account_id) {
    rdsInstances.value = []
    return
  }

  loading.value = true
  try {
    const filter = {
      account_id: filters.account_id,
      engine: filters.engine || undefined,
      status: filters.status || undefined,
      page: pagination.page,
      size: pagination.pageSize
    }
    const res = await listRDS(filter)
    rdsInstances.value = res.items || res
    pagination.total = res.total || rdsInstances.value.length
  } catch (e: any) {
    console.error(e)
    ElMessage.error(e.message || '获取RDS实例列表失败')
    rdsInstances.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: RDS) => {
  selectedRDS.value = row
  detailDialogVisible.value = true
}

const showCreateDialog = () => {
  if (!filters.account_id) {
    ElMessage.warning('请先选择云账号')
    return
  }
  createForm.value.account_id = filters.account_id
  createDialogVisible.value = true
}

const handleCreate = async () => {
  createLoading.value = true
  try {
    const instance = await createRDS(createForm.value)
    ElMessage.success(`RDS实例 ${instance.name} 创建成功`)
    createDialogVisible.value = false
    loadRDSInstances()
  } catch (e: any) {
    ElMessage.error(e.message || '创建失败')
  } finally {
    createLoading.value = false
  }
}

const handleAction = async (row: RDS, action: string) => {
  try {
    await ElMessageBox.confirm(`确认执行 ${action} 操作？`, '确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await rdsAction(filters.account_id!, row.id, action)
    ElMessage.success(`操作执行成功`)
    loadRDSInstances()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '操作失败')
    }
  }
}

const showResizeDialog = (row: RDS) => {
  resizeForm.value.instance_id = row.id
  resizeForm.value.instance_type = row.instance_type
  resizeForm.value.storage_size = row.storage_size
  resizeDialogVisible.value = true
}

const handleResize = async () => {
  resizeLoading.value = true
  try {
    await resizeRDS(filters.account_id!, resizeForm.value.instance_id, resizeForm.value.instance_type, resizeForm.value.storage_size)
    ElMessage.success('规格调整成功')
    resizeDialogVisible.value = false
    loadRDSInstances()
  } catch (e: any) {
    ElMessage.error(e.message || '调整失败')
  } finally {
    resizeLoading.value = false
  }
}

const showBackupDialog = (row: RDS) => {
  backupForm.value.instance_id = row.id
  backupForm.value.name = ''
  backupDialogVisible.value = true
}

const handleBackup = async () => {
  backupLoading.value = true
  try {
    await createRDSBackup(filters.account_id!, backupForm.value.instance_id, backupForm.value.name)
    ElMessage.success('备份创建成功')
    backupDialogVisible.value = false
  } catch (e: any) {
    ElMessage.error(e.message || '创建备份失败')
  } finally {
    backupLoading.value = false
  }
}

const handleDelete = async (row: RDS) => {
  try {
    await ElMessageBox.confirm(`确认删除RDS实例 ${row.name}？此操作不可恢复！`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await deleteRDS(filters.account_id!, row.id)
    ElMessage.success('删除成功')
    loadRDSInstances()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const resetFilters = () => {
  filters.account_id = null
  filters.engine = ''
  filters.status = ''
  pagination.page = 1
  rdsInstances.value = []
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadRDSInstances()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadRDSInstances()
}

onMounted(() => {
  // Wait for account selector to initialize
})
</script>

<style scoped>
.rds-container {
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
  gap: 8px;
  align-items: center;
}

.filter-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}
</style>