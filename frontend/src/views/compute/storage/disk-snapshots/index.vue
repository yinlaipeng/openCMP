<template>
  <div class="disk-snapshots-container">
    <div class="page-header">
      <h2>硬盘快照</h2>
      <div class="toolbar">
        <el-button link type="primary" @click="handleView">
          <el-icon><View /></el-icon>
          查看
        </el-button>
        <el-button :disabled="selectedSnapshots.length === 0" @click="handleSyncStatus">
          同步状态
        </el-button>
        <el-button :disabled="selectedSnapshots.length === 0" @click="handleSetTags">
          设置标签
        </el-button>
        <el-button :disabled="selectedSnapshots.length === 0" type="danger" @click="handleBatchDelete">
          删除
        </el-button>
        <el-button @click="handleTags">
          <el-icon><PriceTag /></el-icon>
          标签
        </el-button>
      </div>
    </div>

    <!-- 公有云/私有云 Tabs -->
    <el-tabs v-model="activeTab" class="snapshot-tabs" @tab-change="handleTabChange">
      <el-tab-pane label="公有云" name="public">
        <el-card class="filter-card">
          <el-form :inline="true" :model="filters" @submit.prevent="fetchData">
            <el-form-item label="名称">
              <el-input v-model="filters.name" placeholder="搜索快照名称" clearable style="width: 180px" />
            </el-form-item>
            <el-form-item label="状态">
              <el-select v-model="filters.status" placeholder="选择状态" clearable style="width: 120px">
                <el-option label="可用" value="available" />
                <el-option label="创建中" value="creating" />
                <el-option label="回滚中" value="rollbacking" />
                <el-option label="错误" value="error" />
              </el-select>
            </el-form-item>
            <el-form-item label="平台">
              <el-select v-model="filters.platform" placeholder="选择平台" clearable style="width: 120px">
                <el-option label="阿里云" value="aliyun" />
                <el-option label="腾讯云" value="tencent" />
                <el-option label="AWS" value="aws" />
              </el-select>
            </el-form-item>
            <el-form-item label="磁盘类型">
              <el-select v-model="filters.disk_type" placeholder="选择类型" clearable style="width: 120px">
                <el-option label="系统盘" value="system" />
                <el-option label="数据盘" value="data" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="fetchData">查询</el-button>
              <el-button @click="resetFilters">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-table
          :data="snapshots"
          v-loading="loading"
          style="width: 100%"
          row-key="id"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column label="名称" min-width="200">
            <template #default="{ row }">
              <el-link type="primary" :underline="false" @click="handleDetails(row)">
                {{ row.name }}
              </el-link>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="标签" width="100">
            <template #default="{ row }">
              <el-tag v-for="tag in (row.tags || [])" :key="tag.key" size="small" class="tag-item">
                {{ tag.key }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="快照大小" width="100">
            <template #default="{ row }">
              {{ row.size }} GB
            </template>
          </el-table-column>
          <el-table-column prop="disk_type" label="磁盘类型" width="100">
            <template #default="{ row }">
              <el-tag size="small" :type="row.disk_type === 'system' ? 'primary' : 'info'">
                {{ row.disk_type === 'system' ? '系统盘' : '数据盘' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="disk_name" label="磁盘" width="150" />
          <el-table-column prop="vm_name" label="虚拟机" width="150" />
          <el-table-column label="平台" width="100">
            <template #default="{ row }">
              <el-tag size="small" :type="getPlatformType(row.provider_type)">
                {{ getPlatformLabel(row.provider_type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="account_name" label="云账号" width="120" />
          <el-table-column prop="project_name" label="项目" width="120" />
          <el-table-column prop="region_id" label="区域" width="120" />
          <el-table-column label="操作" width="160" fixed="right">
            <template #default="{ row }">
              <el-dropdown trigger="click" @command="(cmd: string) => handleActionCommand(cmd, row)">
                <el-button size="small" link type="primary">
                  操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="details">查看详情</el-dropdown-item>
                    <el-dropdown-item command="rollback">回滚磁盘</el-dropdown-item>
                    <el-dropdown-item command="create_disk">创建磁盘</el-dropdown-item>
                    <el-dropdown-item divided command="delete">删除</el-dropdown-item>
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
      </el-tab-pane>
      <el-tab-pane label="私有云" name="private">
        <el-card class="filter-card">
          <el-alert type="info" :closable="false" show-icon>
            <template #title>私有云硬盘快照管理功能开发中</template>
          </el-alert>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 详情弹窗 -->
    <el-dialog
      title="硬盘快照详情"
      v-model="detailDialogVisible"
      width="800px"
    >
      <el-tabs v-model="detailTab" v-if="selectedSnapshot">
        <el-tab-pane label="基础信息" name="basic">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">{{ selectedSnapshot.snapshot_id || selectedSnapshot.id }}</el-descriptions-item>
            <el-descriptions-item label="名称">{{ selectedSnapshot.name }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(selectedSnapshot.status)">{{ getStatusLabel(selectedSnapshot.status) }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="快照大小">{{ selectedSnapshot.size }} GB</el-descriptions-item>
            <el-descriptions-item label="磁盘类型">{{ selectedSnapshot.disk_type === 'system' ? '系统盘' : '数据盘' }}</el-descriptions-item>
            <el-descriptions-item label="来源磁盘">{{ selectedSnapshot.disk_name }}</el-descriptions-item>
            <el-descriptions-item label="虚拟机">{{ selectedSnapshot.vm_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="进度">
              <el-progress :percentage="selectedSnapshot.progress || 100" />
            </el-descriptions-item>
            <el-descriptions-item label="平台">
              <el-tag size="small" :type="getPlatformType(selectedSnapshot.provider_type)">
                {{ getPlatformLabel(selectedSnapshot.provider_type) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="云账号">{{ selectedSnapshot.account_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="区域">{{ selectedSnapshot.region_id || '-' }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ selectedSnapshot.created_at }}</el-descriptions-item>
            <el-descriptions-item label="描述">{{ selectedSnapshot.description || '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="标签" name="tags">
          <el-table :data="selectedSnapshotTags" style="width: 100%">
            <el-table-column prop="key" label="标签键" width="200" />
            <el-table-column prop="value" label="标签值" />
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>

    <!-- 回滚磁盘弹窗 -->
    <el-dialog
      title="回滚磁盘"
      v-model="rollbackDialogVisible"
      width="500px"
    >
      <el-form :model="rollbackForm" label-width="100px">
        <el-form-item label="快照">{{ rollbackForm.snapshot_name }}</el-form-item>
        <el-form-item label="目标磁盘">{{ rollbackForm.disk_name }}</el-form-item>
        <el-alert type="warning" :closable="false" show-icon>
          <template #title>回滚操作将覆盖磁盘上的数据，请确保已做好数据备份</template>
        </el-alert>
      </el-form>
      <template #footer>
        <el-button @click="rollbackDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmRollback" :loading="rollbacking">确认回滚</el-button>
      </template>
    </el-dialog>

    <!-- 创建磁盘弹窗 -->
    <el-dialog
      title="使用快照创建磁盘"
      v-model="createDiskDialogVisible"
      width="500px"
    >
      <el-form :model="createDiskForm" :rules="createDiskRules" ref="createDiskFormRef" label-width="100px">
        <el-form-item label="快照">{{ createDiskForm.snapshot_name }}</el-form-item>
        <el-form-item label="磁盘名称" prop="name">
          <el-input v-model="createDiskForm.name" placeholder="请输入磁盘名称" />
        </el-form-item>
        <el-form-item label="磁盘类型">
          <el-select v-model="createDiskForm.disk_type" style="width: 100%">
            <el-option label="SSD云盘" value="ssd" />
            <el-option label="高效云盘" value="efficiency" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDiskDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCreateDisk" :loading="creatingDisk">创建磁盘</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, ArrowDown, PriceTag } from '@element-plus/icons-vue'
import {
  getCloudDiskSnapshots,
  deleteCloudDiskSnapshot,
  batchDeleteCloudDiskSnapshots,
  rollbackCloudDiskSnapshot,
  createDiskFromSnapshot,
  type CloudDiskSnapshot,
  type SnapshotListParams
} from '@/api/storage'

// Data
const loading = ref(false)
const rollbacking = ref(false)
const creatingDisk = ref(false)
const snapshots = ref<CloudDiskSnapshot[]>([])
const selectedSnapshots = ref<CloudDiskSnapshot[]>([])

const activeTab = ref('public')
const detailTab = ref('basic')
const detailDialogVisible = ref(false)
const rollbackDialogVisible = ref(false)
const createDiskDialogVisible = ref(false)

const selectedSnapshot = ref<CloudDiskSnapshot | null>(null)
const selectedSnapshotTags = ref<{ key: string; value: string }[]>([])

const filters = reactive({
  name: '',
  status: '',
  platform: '',
  disk_type: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const rollbackForm = reactive({
  snapshot_id: '',
  snapshot_name: '',
  disk_id: '',
  disk_name: ''
})

const createDiskForm = reactive({
  snapshot_id: '',
  snapshot_name: '',
  name: '',
  disk_type: 'ssd'
})

const createDiskRules = {
  name: [{ required: true, message: '请输入磁盘名称', trigger: 'blur' }]
}

// Methods
const getStatusType = (status: string) => {
  switch (status) {
    case 'available': return 'success'
    case 'creating': return 'warning'
    case 'rollbacking': return 'warning'
    case 'error': return 'danger'
    default: return 'info'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'available': return '可用'
    case 'creating': return '创建中'
    case 'rollbacking': return '回滚中'
    case 'error': return '错误'
    default: return status
  }
}

const getPlatformLabel = (platform: string) => {
  const labels: Record<string, string> = {
    aliyun: '阿里云',
    tencent: '腾讯云',
    aws: 'AWS'
  }
  return labels[platform] || platform || '未知'
}

const getPlatformType = (platform: string) => {
  const types: Record<string, 'primary' | 'warning' | 'success' | 'info'> = {
    aliyun: 'primary',
    tencent: 'warning',
    aws: 'success'
  }
  return types[platform] || 'info'
}

const handleSelectionChange = (selection: CloudDiskSnapshot[]) => {
  selectedSnapshots.value = selection
}

const handleTabChange = () => {
  pagination.page = 1
  fetchData()
}

const resetFilters = () => {
  filters.name = ''
  filters.status = ''
  filters.platform = ''
  filters.disk_type = ''
  pagination.page = 1
  fetchData()
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: SnapshotListParams = {
      page: pagination.page,
      page_size: pagination.pageSize,
      name: filters.name || undefined,
      status: filters.status || undefined,
      platform: filters.platform || undefined,
      disk_type: filters.disk_type || undefined
    }
    const res = await getCloudDiskSnapshots(params)
    snapshots.value = res.items
    pagination.total = res.total
  } catch (error) {
    console.error('Failed to fetch snapshots:', error)
    ElMessage.error('获取硬盘快照列表失败')
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

const handleSyncStatus = () => {
  ElMessage.info('同步状态功能开发中')
}

const handleSetTags = () => {
  ElMessage.info('设置标签功能开发中')
}

const handleBatchDelete = async () => {
  if (selectedSnapshots.value.length === 0) return
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedSnapshots.value.length} 个快照吗？`,
      '批量删除确认',
      { type: 'warning' }
    )
    const ids = selectedSnapshots.value.map(s => s.id)
    await batchDeleteCloudDiskSnapshots(ids)
    ElMessage.success('删除成功')
    selectedSnapshots.value = []
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

const handleActionCommand = (command: string, row: CloudDiskSnapshot) => {
  switch (command) {
    case 'details': handleDetails(row); break
    case 'rollback': handleRollback(row); break
    case 'create_disk': handleCreateDisk(row); break
    case 'delete': handleDelete(row); break
  }
}

const handleDetails = (row: CloudDiskSnapshot) => {
  selectedSnapshot.value = row
  detailTab.value = 'basic'
  selectedSnapshotTags.value = [
    { key: 'snapshot-type', value: row.disk_type === 'system' ? '系统盘快照' : '数据盘快照' },
    { key: 'source-disk', value: row.disk_name }
  ]
  detailDialogVisible.value = true
}

const handleRollback = (row: CloudDiskSnapshot) => {
  Object.assign(rollbackForm, {
    snapshot_id: row.snapshot_id,
    snapshot_name: row.name,
    disk_id: row.disk_id,
    disk_name: row.disk_name
  })
  rollbackDialogVisible.value = true
}

const confirmRollback = async () => {
  rollbacking.value = true
  try {
    await rollbackCloudDiskSnapshot(selectedSnapshot.value!.id)
    ElMessage.success('回滚任务已提交')
    rollbackDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('回滚失败')
  } finally {
    rollbacking.value = false
  }
}

const handleCreateDisk = (row: CloudDiskSnapshot) => {
  Object.assign(createDiskForm, {
    snapshot_id: row.snapshot_id,
    snapshot_name: row.name,
    name: `disk-from-${row.name}`,
    disk_type: 'ssd'
  })
  createDiskDialogVisible.value = true
}

const confirmCreateDisk = async () => {
  creatingDisk.value = true
  try {
    await createDiskFromSnapshot(selectedSnapshot.value!.id, {
      name: createDiskForm.name,
      disk_type: createDiskForm.disk_type
    })
    ElMessage.success('创建磁盘任务已提交')
    createDiskDialogVisible.value = false
  } catch (error) {
    ElMessage.error('创建磁盘失败')
  } finally {
    creatingDisk.value = false
  }
}

const handleDelete = async (row: CloudDiskSnapshot) => {
  try {
    await ElMessageBox.confirm(`确定要删除快照 "${row.name}" 吗？`, '删除警告', { type: 'warning' })
    await deleteCloudDiskSnapshot(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
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

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.disk-snapshots-container { padding: 20px; }

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

.toolbar { display: flex; gap: 10px; align-items: center; }
.filter-card { margin-bottom: 20px; }
.snapshot-tabs { margin-bottom: 20px; }

.platform-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.account-name { font-size: 12px; color: var(--el-text-color-secondary); }

.pagination { margin-top: 20px; justify-content: flex-end; }
.tag-item { margin-right: 4px; margin-bottom: 2px; }
</style>