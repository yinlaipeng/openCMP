<template>
  <div class="host-snapshots-container">
    <div class="page-header">
      <h2>主机快照</h2>
      <div class="toolbar">
        <el-button link type="primary" @click="handleView">
          <el-icon><View /></el-icon>
          查看
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
          <el-table-column label="磁盘快照数" width="100">
            <template #default="{ row }">
              {{ row.disk_snapshots || 0 }}
            </template>
          </el-table-column>
          <el-table-column label="内存快照" width="100">
            <template #default="{ row }">
              <el-tag :type="row.memory_snapshot ? 'success' : 'info'" size="small">
                {{ row.memory_snapshot ? '已包含' : '未包含' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="cpu_arch" label="CPU架构" width="100" />
          <el-table-column label="快照大小" width="100">
            <template #default="{ row }">
              {{ row.size }} GB
            </template>
          </el-table-column>
          <el-table-column prop="instance_name" label="虚拟机" width="150" />
          <el-table-column label="平台" width="100">
            <template #default="{ row }">
              <el-tag size="small" :type="getPlatformType(row.provider_type)">
                {{ getPlatformLabel(row.provider_type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="project_name" label="项目" width="120" />
          <el-table-column label="操作" width="160" fixed="right">
            <template #default="{ row }">
              <el-dropdown trigger="click" @command="(cmd: string) => handleActionCommand(cmd, row)">
                <el-button size="small" link type="primary">
                  操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="details">查看详情</el-dropdown-item>
                    <el-dropdown-item command="rollback">恢复主机</el-dropdown-item>
                    <el-dropdown-item command="create_vm">创建主机</el-dropdown-item>
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
            <template #title>私有云主机快照管理功能开发中</template>
          </el-alert>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 详情弹窗 -->
    <el-dialog
      title="主机快照详情"
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
            <el-descriptions-item label="磁盘快照数">{{ selectedSnapshot.disk_snapshots || 0 }}</el-descriptions-item>
            <el-descriptions-item label="内存快照">
              <el-tag :type="selectedSnapshot.memory_snapshot ? 'success' : 'info'" size="small">
                {{ selectedSnapshot.memory_snapshot ? '已包含' : '未包含' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="CPU架构">{{ selectedSnapshot.cpu_arch || '-' }}</el-descriptions-item>
            <el-descriptions-item label="虚拟机">{{ selectedSnapshot.vm_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="平台">
              <el-tag size="small" :type="getPlatformType(selectedSnapshot.provider_type)">
                {{ getPlatformLabel(selectedSnapshot.provider_type) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="区域">{{ selectedSnapshot.region_id || '-' }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ selectedSnapshot.created_at }}</el-descriptions-item>
            <el-descriptions-item label="描述" :span="2">{{ selectedSnapshot.description || '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="磁盘快照" name="disk_snapshots">
          <el-table :data="selectedDiskSnapshots" style="width: 100%">
            <el-table-column prop="disk_name" label="磁盘名称" width="200" />
            <el-table-column prop="disk_type" label="磁盘类型" width="100" />
            <el-table-column prop="size" label="快照大小" width="100" />
            <el-table-column prop="snapshot_name" label="快照名称" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="标签" name="tags">
          <el-table :data="selectedSnapshotTags" style="width: 100%">
            <el-table-column prop="key" label="标签键" width="200" />
            <el-table-column prop="value" label="标签值" />
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>

    <!-- 恢复主机弹窗 -->
    <el-dialog
      title="恢复主机"
      v-model="rollbackDialogVisible"
      width="500px"
    >
      <el-form :model="rollbackForm" label-width="100px">
        <el-form-item label="快照">{{ rollbackForm.snapshot_name }}</el-form-item>
        <el-form-item label="目标主机">{{ rollbackForm.vm_name }}</el-form-item>
        <el-alert type="warning" :closable="false" show-icon>
          <template #title>恢复操作将覆盖主机当前状态，请确保已做好数据备份</template>
        </el-alert>
      </el-form>
      <template #footer>
        <el-button @click="rollbackDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmRollback" :loading="rollbacking">确认恢复</el-button>
      </template>
    </el-dialog>

    <!-- 创建主机弹窗 -->
    <el-dialog
      title="使用快照创建主机"
      v-model="createVMDialogVisible"
      width="500px"
    >
      <el-form :model="createVMForm" :rules="createVMRules" ref="createVMFormRef" label-width="100px">
        <el-form-item label="快照">{{ createVMForm.snapshot_name }}</el-form-item>
        <el-form-item label="主机名称" prop="name">
          <el-input v-model="createVMForm.name" placeholder="请输入主机名称" />
        </el-form-item>
        <el-form-item label="可用区">
          <el-select v-model="createVMForm.zone_id" placeholder="选择可用区" style="width: 100%">
            <el-option label="可用区A" value="zone-a" />
            <el-option label="可用区B" value="zone-b" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVMDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCreateVM" :loading="creatingVM">创建主机</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, ArrowDown, PriceTag } from '@element-plus/icons-vue'
import {
  getInstanceSnapshots,
  deleteInstanceSnapshot,
  batchDeleteInstanceSnapshots,
  rollbackInstanceSnapshot,
  createVMFromSnapshot,
  type InstanceSnapshot,
  type InstanceSnapshotListParams
} from '@/api/storage'

interface HostSnapshot {
  id: number | string
  snapshot_id: string
  name: string
  status: string
  size: number
  disk_snapshots: number
  memory_snapshot: boolean
  cpu_arch: string
  vm_id: string
  vm_name: string
  provider_type: string
  cloud_account_id: number
  region_id: string
  created_at: string
  description: string
}

// Data
const loading = ref(false)
const rollbacking = ref(false)
const creatingVM = ref(false)
const snapshots = ref<InstanceSnapshot[]>([])
const selectedSnapshots = ref<InstanceSnapshot[]>([])

const activeTab = ref('public')
const detailTab = ref('basic')
const detailDialogVisible = ref(false)
const rollbackDialogVisible = ref(false)
const createVMDialogVisible = ref(false)

const selectedSnapshot = ref<InstanceSnapshot | null>(null)
const selectedSnapshotTags = ref<{ key: string; value: string }[]>([])
const selectedDiskSnapshots = ref<{ disk_name: string; disk_type: string; size: number; snapshot_name: string }[]>([])

const filters = reactive({
  name: '',
  status: '',
  platform: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const rollbackForm = reactive({
  snapshot_id: '',
  snapshot_name: '',
  vm_id: '',
  vm_name: ''
})

const createVMForm = reactive({
  snapshot_id: '',
  snapshot_name: '',
  name: '',
  zone_id: ''
})

const createVMRules = {
  name: [{ required: true, message: '请输入主机名称', trigger: 'blur' }]
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
    case 'rollbacking': return '恢复中'
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

const handleSelectionChange = (selection: InstanceSnapshot[]) => {
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
  pagination.page = 1
  fetchData()
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: InstanceSnapshotListParams = {
      page: pagination.page,
      page_size: pagination.pageSize,
      name: filters.name || undefined,
      status: filters.status || undefined,
      platform: filters.platform || undefined
    }
    const res = await getInstanceSnapshots(params)
    snapshots.value = res.items
    pagination.total = res.total
  } catch (error) {
    console.error('Failed to fetch snapshots:', error)
    ElMessage.error('获取主机快照列表失败')
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

const handleSetTags = () => {
  ElMessage.info('设置标签功能开发中')
}

const handleBatchDelete = async () => {
  if (selectedSnapshots.value.length === 0) return
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedSnapshots.value.length} 个主机快照吗？`,
      '批量删除确认',
      { type: 'warning' }
    )
    const ids = selectedSnapshots.value.map(s => s.id)
    await batchDeleteInstanceSnapshots(ids)
    ElMessage.success('删除成功')
    selectedSnapshots.value = []
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

const handleActionCommand = (command: string, row: InstanceSnapshot) => {
  switch (command) {
    case 'details': handleDetails(row); break
    case 'rollback': handleRollback(row); break
    case 'create_vm': handleCreateVM(row); break
    case 'delete': handleDelete(row); break
  }
}

const handleDetails = (row: InstanceSnapshot) => {
  selectedSnapshot.value = row
  detailTab.value = 'basic'
  selectedSnapshotTags.value = [
    { key: 'snapshot-type', value: '主机快照' },
    { key: 'source-vm', value: row.instance_name }
  ]
  selectedDiskSnapshots.value = [
    { disk_name: 'disk-system', disk_type: '系统盘', size: 5, snapshot_name: `${row.name}-disk-system` },
    { disk_name: 'disk-data-1', disk_type: '数据盘', size: 5, snapshot_name: `${row.name}-disk-data-1` },
    { disk_name: 'disk-data-2', disk_type: '数据盘', size: 5, snapshot_name: `${row.name}-disk-data-2` }
  ]
  detailDialogVisible.value = true
}

const handleRollback = (row: InstanceSnapshot) => {
  Object.assign(rollbackForm, {
    snapshot_id: row.snapshot_id,
    snapshot_name: row.name,
    vm_id: row.instance_id,
    vm_name: row.instance_name
  })
  rollbackDialogVisible.value = true
}

const confirmRollback = async () => {
  rollbacking.value = true
  try {
    await rollbackInstanceSnapshot(selectedSnapshot.value!.id)
    ElMessage.success('恢复任务已提交')
    rollbackDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('恢复失败')
  } finally {
    rollbacking.value = false
  }
}

const handleCreateVM = (row: InstanceSnapshot) => {
  Object.assign(createVMForm, {
    snapshot_id: row.snapshot_id,
    snapshot_name: row.name,
    name: `vm-from-${row.name}`,
    zone_id: ''
  })
  createVMDialogVisible.value = true
}

const confirmCreateVM = async () => {
  creatingVM.value = true
  try {
    await createVMFromSnapshot(selectedSnapshot.value!.id, {
      name: createVMForm.name,
      zone_id: createVMForm.zone_id
    })
    ElMessage.success('创建主机任务已提交')
    createVMDialogVisible.value = false
  } catch (error) {
    ElMessage.error('创建主机失败')
  } finally {
    creatingVM.value = false
  }
}

const handleDelete = async (row: InstanceSnapshot) => {
  try {
    await ElMessageBox.confirm(`确定要删除主机快照 "${row.name}" 吗？`, '删除警告', { type: 'warning' })
    await deleteInstanceSnapshot(row.id)
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
.host-snapshots-container { padding: 20px; }

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

.pagination { margin-top: 20px; justify-content: flex-end; }
.tag-item { margin-right: 4px; margin-bottom: 2px; }
</style>