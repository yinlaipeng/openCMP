<template>
  <div class="mongodb-container">
    <div class="page-header">
      <h2>MongoDB实例列表</h2>
      <div class="toolbar">
        <el-button @click="loadMongoDBInstances" :icon="Refresh" circle />
        <el-button :disabled="selectedRows.length === 0" @click="handleSyncStatus">同步状态</el-button>
        <el-dropdown :disabled="selectedRows.length === 0" trigger="click">
          <el-button :disabled="selectedRows.length === 0">
            批量操作<el-icon class="el-icon--right"><arrow-down /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleBatchAction('delete')">批量删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button>标签</el-button>
      </div>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="loadMongoDBInstances">
        <el-form-item label="名称">
          <el-input v-model="filters.name" placeholder="搜索名称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="状态" clearable>
            <el-option label="运行中" value="Running" />
            <el-option label="已停止" value="Stopped" />
            <el-option label="创建中" value="Creating" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadMongoDBInstances">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table
      :data="mongodbInstances"
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
      <el-table-column label="标签" width="150">
        <template #default="{ row }">
          <div v-if="row.tags && Object.keys(row.tags).length > 0">
            <el-tag v-for="(value, key) in row.tags" :key="key" size="small" style="margin-right: 4px;">
              {{ key }}: {{ value }}
            </el-tag>
          </div>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="configuration" label="配置" width="120">
        <template #default="{ row }">
          {{ row.configuration || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="address" label="地址" width="200" />
      <el-table-column prop="network_address" label="网络地址" width="200">
        <template #default="{ row }">
          {{ row.network_address || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="engine_version" label="引擎版本" width="100" />
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
      <el-table-column prop="region" label="区域" width="120" />
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-dropdown trigger="click">
            <el-button size="small">操作<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="viewDetail(row)">详情</el-dropdown-item>
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
    <el-dialog v-model="detailDialogVisible" title="MongoDB实例详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedMongoDB?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedMongoDB?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedMongoDB?.status }}</el-descriptions-item>
        <el-descriptions-item label="配置">{{ selectedMongoDB?.configuration }}</el-descriptions-item>
        <el-descriptions-item label="地址">{{ selectedMongoDB?.address }}</el-descriptions-item>
        <el-descriptions-item label="网络地址">{{ selectedMongoDB?.network_address }}</el-descriptions-item>
        <el-descriptions-item label="引擎版本">{{ selectedMongoDB?.engine_version }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ getPlatformLabel(selectedMongoDB?.platform) }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedMongoDB?.account_name }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedMongoDB?.project }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedMongoDB?.region }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedMongoDB?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown, Refresh } from '@element-plus/icons-vue'
import { listMongoDB, deleteMongoDB, type MongoDBInstance } from '@/api/database'

interface MongoDB extends MongoDBInstance {
  platform?: string
  account_name?: string
  project?: string
}

const mongodbInstances = ref<MongoDB[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedMongoDB = ref<MongoDB | null>(null)
const selectedRows = ref<MongoDB[]>([])

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 筛选条件
const filters = reactive({
  name: '',
  status: ''
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

const getPlatformLabel = (platform?: string): string => {
  if (!platform) return '未知'
  return platformLabels[platform] || platform || '未知'
}

const getPlatformType = (platform?: string): 'primary' | 'warning' | 'success' | 'info' => {
  if (!platform) return 'info'
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

const handleSelectionChange = (rows: MongoDB[]) => {
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

    // MongoDB删除需要account_id，这里使用默认账号ID
    const defaultAccountId = 1

    for (const row of selectedRows.value) {
      if (action === 'delete') {
        await deleteMongoDB(defaultAccountId, row.id)
      }
    }
    ElMessage.success(`批量 ${action} 执行成功`)
    loadMongoDBInstances()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '操作失败')
    }
  }
}

const loadMongoDBInstances = async () => {
  loading.value = true
  try {
    const res = await listMongoDB({
      name: filters.name || undefined,
      status: filters.status || undefined
    })
    mongodbInstances.value = res || []
    pagination.total = mongodbInstances.value.length
  } catch (e: any) {
    console.error(e)
    ElMessage.error(e.message || '获取MongoDB实例列表失败')
    mongodbInstances.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: MongoDB) => {
  selectedMongoDB.value = row
  detailDialogVisible.value = true
}

const handleDelete = async (row: MongoDB) => {
  try {
    await ElMessageBox.confirm(`确认删除MongoDB实例 ${row.name}？此操作不可恢复！`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    // MongoDB删除需要account_id，这里使用默认账号ID
    const defaultAccountId = 1
    await deleteMongoDB(defaultAccountId, row.id)
    ElMessage.success('删除成功')
    loadMongoDBInstances()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const resetFilters = () => {
  filters.name = ''
  filters.status = ''
  pagination.page = 1
  loadMongoDBInstances()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadMongoDBInstances()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadMongoDBInstances()
}

onMounted(() => {
  loadMongoDBInstances()
})
</script>

<style scoped>
.mongodb-container {
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
</style>