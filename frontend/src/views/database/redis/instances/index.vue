<template>
  <div class="redis-container">
    <div class="page-header">
      <h2>Redis实例管理</h2>
      <div class="toolbar">
        <el-button @click="loadCacheInstances" :icon="Refresh" circle />
        <el-button type="primary" @click="showCreateDialog">新建</el-button>
        <el-button :disabled="selectedRows.length === 0" @click="handleSyncStatus">同步状态</el-button>
        <el-dropdown :disabled="selectedRows.length === 0" trigger="click">
          <el-button :disabled="selectedRows.length === 0">
            批量操作<el-icon class="el-icon--right"><arrow-down /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleBatchAction('reboot')">批量重启</el-dropdown-item>
              <el-dropdown-item divided @click="handleBatchAction('delete')">批量删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button>标签</el-button>
      </div>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="loadCacheInstances">
        <el-form-item label="云账号">
          <CloudAccountSelector
            v-model:value="filters.account_id"
            placeholder="选择云账号"
            @change="handleAccountChange"
          />
        </el-form-item>
        <el-form-item label="引擎类型">
          <el-select v-model="filters.engine" placeholder="引擎类型" clearable>
            <el-option label="Redis" value="Redis" />
            <el-option label="Memcached" value="Memcached" />
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
          <el-button type="primary" @click="loadCacheInstances">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table
      :data="cacheInstances"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column label="名称" min-width="180">
        <template #default="{ row }">
          <el-link type="primary" @click="viewDetail(row)">{{ row.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="instance_type" label="实例类型" width="150" />
      <el-table-column label="类型版本" width="120">
        <template #default="{ row }">
          {{ row.engine }} {{ row.engine_version }}
        </template>
      </el-table-column>
      <el-table-column label="密码" width="80">
        <template #default>
          <el-tag size="small">已设置</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="endpoint" label="地址" width="200" />
      <el-table-column prop="port" label="端口" width="80" />
      <el-table-column prop="security_group" label="安全组" width="150">
        <template #default="{ row }">
          {{ row.security_group || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="计费类型" width="100">
        <template #default="{ row }">
          {{ row.billing_type || '按量付费' }}
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
      <el-table-column prop="project" label="项目" width="120">
        <template #default="{ row }">
          {{ row.project || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="zone_id" label="区域" width="120" />
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-dropdown trigger="click">
            <el-button size="small">操作<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
            <template #dropdown>
              <el-dropdown-menu>
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
    <el-dialog v-model="detailDialogVisible" title="Redis实例详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedCache?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedCache?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedCache?.status }}</el-descriptions-item>
        <el-descriptions-item label="实例规格">{{ selectedCache?.instance_type }}</el-descriptions-item>
        <el-descriptions-item label="引擎类型">{{ selectedCache?.engine }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ selectedCache?.engine_version }}</el-descriptions-item>
        <el-descriptions-item label="连接地址">{{ selectedCache?.endpoint }}</el-descriptions-item>
        <el-descriptions-item label="端口">{{ selectedCache?.port }}</el-descriptions-item>
        <el-descriptions-item label="VPC">{{ selectedCache?.vpc_id }}</el-descriptions-item>
        <el-descriptions-item label="子网">{{ selectedCache?.subnet_id }}</el-descriptions-item>
        <el-descriptions-item label="可用区">{{ selectedCache?.zone_id }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedCache?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Create Dialog -->
    <el-dialog v-model="createDialogVisible" title="创建Redis实例" width="700px">
      <el-form :model="createForm" label-width="120px">
        <el-form-item label="指定项目">
          <el-select v-model="createForm.project_id" placeholder="选择项目" clearable>
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称" required>
          <el-input v-model="createForm.name" placeholder="请输入实例名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="createForm.description" type="textarea" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="计费类型">
          <el-radio-group v-model="createForm.billing_type">
            <el-radio value="Postpaid">按量付费</el-radio>
            <el-radio value="Prepaid">包年包月</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="过期释放">
          <el-radio-group v-model="createForm.auto_release">
            <el-radio value="unlimited">不限</el-radio>
            <el-radio value="limited">指定时间</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="创建数量">
          <el-input-number v-model="createForm.quantity" :min="1" :max="10" />
        </el-form-item>
        <el-form-item label="区域" required>
          <el-select v-model="createForm.zone_id" placeholder="选择区域" @change="loadCacheSKUs">
            <el-option v-for="z in zones" :key="z.id" :label="z.name" :value="z.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="类型" required>
          <el-select v-model="createForm.engine" placeholder="选择类型" @change="loadCacheSKUs">
            <el-option label="Redis" value="Redis" />
            <el-option label="Memcached" value="Memcached" />
          </el-select>
        </el-form-item>
        <el-form-item label="版本" required>
          <el-select v-model="createForm.engine_version" placeholder="选择版本" @change="loadCacheSKUs">
            <el-option label="Redis 5.0" value="5.0" />
            <el-option label="Redis 6.0" value="6.0" />
            <el-option label="Redis 7.0" value="7.0" />
          </el-select>
        </el-form-item>
        <el-form-item label="节点类型">
          <el-select v-model="createForm.node_type" placeholder="选择节点类型" @change="loadCacheSKUs">
            <el-option label="单节点" value="standalone" />
            <el-option label="主备" value="ha" />
            <el-option label="集群" value="cluster" />
          </el-select>
        </el-form-item>
        <el-form-item label="性能类型">
          <el-select v-model="createForm.performance_type" placeholder="选择性能类型" @change="loadCacheSKUs">
            <el-option label="标准版" value="standard" />
            <el-option label="性能增强版" value="performance" />
          </el-select>
        </el-form-item>
        <el-form-item label="内存大小">
          <el-select v-model="createForm.memory_mb" placeholder="选择内存">
            <el-option label="1 GB" :value="1024" />
            <el-option label="2 GB" :value="2048" />
            <el-option label="4 GB" :value="4096" />
            <el-option label="8 GB" :value="8192" />
            <el-option label="16 GB" :value="16384" />
          </el-select>
        </el-form-item>
        <el-form-item label="实例规格" required>
          <el-select v-model="createForm.instance_type" placeholder="选择实例规格">
            <el-option
              v-for="sku in cacheSkus"
              :key="sku.id"
              :label="`${sku.instance_type} (${sku.memory_mb}MB)`"
              :value="sku.instance_type"
            />
          </el-select>
          <span v-if="cacheSkus.length > 0" class="sku-count">共 {{ cacheSkus.length }} 个可用规格</span>
        </el-form-item>
        <el-form-item label="VPC" required>
          <el-input v-model="createForm.vpc_id" placeholder="VPC ID" />
        </el-form-item>
        <el-form-item label="子网" required>
          <el-input v-model="createForm.subnet_id" placeholder="子网 ID" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="createLoading">创建</el-button>
      </template>
    </el-dialog>

    <!-- Resize Dialog -->
    <el-dialog v-model="resizeDialogVisible" title="调整Redis规格" width="400px">
      <el-form :model="resizeForm" label-width="100px">
        <el-form-item label="新规格">
          <el-input v-model="resizeForm.instance_type" placeholder="如: redis.standard.medium.default" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resizeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleResize" :loading="resizeLoading">确认调整</el-button>
      </template>
    </el-dialog>

    <!-- Backup Dialog -->
    <el-dialog v-model="backupDialogVisible" title="创建备份" width="300px">
      <p>为实例 <strong>{{ backupForm.instance_id }}</strong> 创建备份</p>
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
import { ArrowDown, Refresh } from '@element-plus/icons-vue'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import { listCache, createCache, deleteCache, cacheAction, resizeCache, createCacheBackup, listCacheSKUs, type CacheInstance, type CacheConfig, type CacheInstanceSKU } from '@/api/database'

interface Redis extends CacheInstance {
  platform?: string
  account_name?: string
  project?: string
  billing_type?: string
  security_group?: string
}

const cacheInstances = ref<Redis[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedCache = ref<Redis | null>(null)
const selectedRows = ref<Redis[]>([])

// SKU数据
const cacheSkus = ref<CacheInstanceSKU[]>([])

// 项目和区域数据
const projects = ref<any[]>([])
const zones = ref<any[]>([])

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
const createForm = ref({
  account_id: 0,
  project_id: null as number | null,
  name: '',
  description: '',
  billing_type: 'Postpaid',
  auto_release: 'unlimited',
  quantity: 1,
  engine: 'Redis',
  engine_version: '6.0',
  node_type: 'standalone',
  performance_type: 'standard',
  memory_mb: 1024,
  instance_type: '',
  vpc_id: '',
  subnet_id: '',
  zone_id: '',
  tags: {} as Record<string, string>
})

// Resize dialog
const resizeDialogVisible = ref(false)
const resizeLoading = ref(false)
const resizeForm = ref({
  instance_id: '',
  instance_type: ''
})

// Backup dialog
const backupDialogVisible = ref(false)
const backupLoading = ref(false)
const backupForm = ref({
  instance_id: ''
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
    case 'running':
    case 'active':
      return 'success'
    case 'creating':
    case 'pending':
      return 'warning'
    case 'stopped':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const handleAccountChange = (accountId: number | null) => {
  filters.account_id = accountId
  if (accountId) {
    loadCacheInstances()
  } else {
    cacheInstances.value = []
  }
}

const handleSelectionChange = (rows: Redis[]) => {
  selectedRows.value = rows
}

const handleSyncStatus = () => {
  ElMessage.info('同步状态功能开发中')
}

const handleBatchAction = async (action: string) => {
  if (selectedRows.value.length === 0) return

  try {
    await ElMessageBox.confirm(`确认对 ${selectedRows.value.length} 个实例执行 ${action} 操作？`, '确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    for (const row of selectedRows.value) {
      if (action === 'reboot') {
        await cacheAction(filters.account_id!, row.id, 'reboot')
      } else if (action === 'delete') {
        await deleteCache(filters.account_id!, row.id)
      }
    }
    ElMessage.success(`批量 ${action} 执行成功`)
    loadCacheInstances()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '操作失败')
    }
  }
}

const loadCacheInstances = async () => {
  if (!filters.account_id) {
    cacheInstances.value = []
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
    const res = await listCache(filter)
    cacheInstances.value = res.items || res
    pagination.total = res.total || cacheInstances.value.length
  } catch (e: any) {
    console.error(e)
    ElMessage.error(e.message || '获取Redis实例列表失败')
    cacheInstances.value = []
  } finally {
    loading.value = false
  }
}

const loadCacheSKUs = async () => {
  if (!filters.account_id) return

  try {
    const res = await listCacheSKUs({
      account_id: filters.account_id,
      engine: createForm.value.engine || undefined,
      engine_version: createForm.value.engine_version || undefined,
      node_type: createForm.value.node_type || undefined,
      performance_type: createForm.value.performance_type || undefined,
      region_id: createForm.value.zone_id || undefined
    })
    cacheSkus.value = res || []

    // 根据内存过滤
    if (createForm.value.memory_mb) {
      cacheSkus.value = cacheSkus.value.filter(sku => sku.memory_mb === createForm.value.memory_mb)
    }
  } catch (e: any) {
    console.error(e)
    cacheSkus.value = []
  }
}

const viewDetail = (row: Redis) => {
  selectedCache.value = row
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
  if (!createForm.value.name || !createForm.value.instance_type) {
    ElMessage.warning('请填写必填字段')
    return
  }

  createLoading.value = true
  try {
    const config: CacheConfig = {
      account_id: createForm.value.account_id,
      name: createForm.value.name,
      engine: createForm.value.engine,
      engine_version: createForm.value.engine_version,
      instance_type: createForm.value.instance_type,
      vpc_id: createForm.value.vpc_id,
      subnet_id: createForm.value.subnet_id,
      zone_id: createForm.value.zone_id,
      tags: createForm.value.tags
    }
    const instance = await createCache(config)
    ElMessage.success(`Redis实例 ${instance.name} 创建成功`)
    createDialogVisible.value = false
    loadCacheInstances()
  } catch (e: any) {
    ElMessage.error(e.message || '创建失败')
  } finally {
    createLoading.value = false
  }
}

const handleAction = async (row: Redis, action: string) => {
  try {
    await ElMessageBox.confirm(`确认执行 ${action} 操作？`, '确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await cacheAction(filters.account_id!, row.id, action)
    ElMessage.success(`操作执行成功`)
    loadCacheInstances()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '操作失败')
    }
  }
}

const showResizeDialog = (row: Redis) => {
  resizeForm.value.instance_id = row.id
  resizeForm.value.instance_type = row.instance_type
  resizeDialogVisible.value = true
}

const handleResize = async () => {
  resizeLoading.value = true
  try {
    await resizeCache(filters.account_id!, resizeForm.value.instance_id, resizeForm.value.instance_type)
    ElMessage.success('规格调整成功')
    resizeDialogVisible.value = false
    loadCacheInstances()
  } catch (e: any) {
    ElMessage.error(e.message || '调整失败')
  } finally {
    resizeLoading.value = false
  }
}

const showBackupDialog = (row: Redis) => {
  backupForm.value.instance_id = row.id
  backupDialogVisible.value = true
}

const handleBackup = async () => {
  backupLoading.value = true
  try {
    await createCacheBackup(filters.account_id!, backupForm.value.instance_id)
    ElMessage.success('备份创建成功')
    backupDialogVisible.value = false
  } catch (e: any) {
    ElMessage.error(e.message || '创建备份失败')
  } finally {
    backupLoading.value = false
  }
}

const handleDelete = async (row: Redis) => {
  try {
    await ElMessageBox.confirm(`确认删除Redis实例 ${row.name}？此操作不可恢复！`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await deleteCache(filters.account_id!, row.id)
    ElMessage.success('删除成功')
    loadCacheInstances()
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
  cacheInstances.value = []
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadCacheInstances()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadCacheInstances()
}

onMounted(() => {
  // 初始化项目和区域数据
  projects.value = [{ id: 1, name: 'system' }]
  zones.value = [
    { id: 'cn-hangzhou', name: '杭州' },
    { id: 'cn-shanghai', name: '上海' },
    { id: 'cn-beijing', name: '北京' }
  ]
})
</script>

<style scoped>
.redis-container {
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

.sku-count {
  margin-left: 8px;
  font-size: 12px;
  color: var(--text-color-secondary);
}
</style>