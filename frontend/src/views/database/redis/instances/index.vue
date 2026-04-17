<template>
  <div class="redis-container">
    <div class="page-header">
      <h2>Redis实例管理</h2>
      <el-button type="primary" @click="showCreateDialog">新建实例</el-button>
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
      <el-table-column prop="instance_type" label="实例规格" width="120" />
      <el-table-column prop="engine" label="引擎类型" width="100" />
      <el-table-column label="版本" width="80">
        <template #default="{ row }">
          {{ row.engine }} {{ row.engine_version }}
        </template>
      </el-table-column>
      <el-table-column prop="endpoint" label="连接地址" width="200" />
      <el-table-column prop="port" label="端口" width="80" />
      <el-table-column prop="vpc_id" label="VPC" width="150" />
      <el-table-column prop="zone_id" label="可用区" width="120" />
      <el-table-column label="操作" width="180" fixed="right">
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
    <el-dialog v-model="createDialogVisible" title="创建Redis实例" width="600px">
      <el-form :model="createForm" label-width="120px">
        <el-form-item label="实例名称" required>
          <el-input v-model="createForm.name" placeholder="请输入实例名称" />
        </el-form-item>
        <el-form-item label="引擎类型" required>
          <el-select v-model="createForm.engine" placeholder="选择引擎类型">
            <el-option label="Redis" value="Redis" />
            <el-option label="Memcached" value="Memcached" />
          </el-select>
        </el-form-item>
        <el-form-item label="引擎版本" required>
          <el-select v-model="createForm.engine_version" placeholder="选择版本">
            <el-option label="Redis 5.0" value="5.0" />
            <el-option label="Redis 6.0" value="6.0" />
            <el-option label="Redis 7.0" value="7.0" />
          </el-select>
        </el-form-item>
        <el-form-item label="实例规格" required>
          <el-input v-model="createForm.instance_type" placeholder="如: redis.standard.small.default" />
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
import { ArrowDown } from '@element-plus/icons-vue'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import { listCache, createCache, deleteCache, cacheAction, resizeCache, createCacheBackup, type CacheInstance, type CacheConfig } from '@/api/database'

interface Redis extends CacheInstance {
  platform?: string
  account_name?: string
}

const cacheInstances = ref<Redis[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedCache = ref<Redis | null>(null)

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
const createForm = ref<CacheConfig>({
  account_id: 0,
  name: '',
  engine: 'Redis',
  engine_version: '6.0',
  instance_type: '',
  vpc_id: '',
  subnet_id: '',
  zone_id: ''
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
  createLoading.value = true
  try {
    const instance = await createCache(createForm.value)
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
  // Wait for account selector to initialize
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