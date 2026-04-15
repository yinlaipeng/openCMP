<template>
  <div class="cloud-disks-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="title">云硬盘列表</span>
          <div class="header-actions">
            <el-button type="primary" size="small" @click="showCreateDialog = true">创建云硬盘</el-button>
            <el-button size="small" @click="handleSync" :loading="syncing">同步云硬盘</el-button>
          </div>
        </div>
      </template>

      <!-- 筛选 -->
      <div class="filter-bar" style="margin-bottom: 16px;">
        <el-select v-model="selectedAccountId" placeholder="选择云账号" clearable style="width: 200px;">
          <el-option v-for="account in cloudAccounts" :key="account.id" :label="account.name" :value="account.id" />
        </el-select>
        <el-select v-model="filterStatus" placeholder="状态筛选" clearable style="width: 120px; margin-left: 8px;">
          <el-option label="可用" value="available" />
          <el-option label="已挂载" value="in_use" />
          <el-option label="创建中" value="creating" />
        </el-select>
        <el-button size="small" @click="loadCloudDisks" :loading="loading" style="margin-left: 8px;">刷新</el-button>
      </div>

      <!-- 数据表格 -->
      <el-table :data="cloudDisks" v-loading="loading" style="width: 100%">
        <el-table-column prop="disk_id" label="磁盘ID" width="150" />
        <el-table-column prop="name" label="名称" width="180" />
        <el-table-column prop="size" label="容量" width="100">
          <template #default="{ row }">{{ row.size }} GB</template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag>{{ getDiskTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="vm_id" label="挂载实例" width="150">
          <template #default="{ row }">{{ row.vm_id || '-' }}</template>
        </el-table-column>
        <el-table-column prop="zone_id" label="可用区" width="120" />
        <el-table-column prop="provider_type" label="云平台" width="100" />
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" link type="primary" @click="handleAttach(row)" v-if="row.status === 'available'">挂载</el-button>
            <el-button size="small" link type="primary" @click="handleDetach(row)" v-if="row.status === 'in_use'">卸载</el-button>
            <el-button size="small" link type="primary" @click="handleResize(row)">扩容</el-button>
            <el-button size="small" link type="primary" @click="handleSnapshot(row)">创建快照</el-button>
            <el-button size="small" link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.currentPage"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        style="margin-top: 16px; justify-content: flex-end;"
      />
    </el-card>

    <!-- 创建云硬盘对话框 -->
    <el-dialog v-model="showCreateDialog" title="创建云硬盘" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="云账号" required>
          <el-select v-model="createForm.cloud_account_id" placeholder="请选择云账号">
            <el-option v-for="account in cloudAccounts" :key="account.id" :label="account.name" :value="account.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称" required>
          <el-input v-model="createForm.name" placeholder="请输入云硬盘名称" />
        </el-form-item>
        <el-form-item label="容量" required>
          <el-input-number v-model="createForm.size" :min="10" :max="32000" style="width: 100%;" />
          <span style="margin-left: 8px;">GB</span>
        </el-form-item>
        <el-form-item label="磁盘类型">
          <el-select v-model="createForm.type" placeholder="请选择磁盘类型">
            <el-option label="ESSD PL0" value="cloud_essd" />
            <el-option label="ESSD PL1" value="cloud_essd_pl1" />
            <el-option label="ESSD PL2" value="cloud_essd_pl2" />
            <el-option label="ESSD PL3" value="cloud_essd_pl3" />
            <el-option label="高效云盘" value="cloud_efficiency" />
            <el-option label="普通云盘" value="cloud" />
          </el-select>
        </el-form-item>
        <el-form-item label="可用区" required>
          <el-input v-model="createForm.zone_id" placeholder="如：cn-hangzhou-a" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="creating">创建</el-button>
      </template>
    </el-dialog>

    <!-- 挂载对话框 -->
    <el-dialog v-model="showAttachDialog" title="挂载云硬盘" width="400px">
      <el-form label-width="80px">
        <el-form-item label="云硬盘">
          <span>{{ currentDisk?.name }} ({{ currentDisk?.size }}GB)</span>
        </el-form-item>
        <el-form-item label="目标实例">
          <el-input v-model="attachVMId" placeholder="请输入实例ID" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAttachDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAttachSubmit" :loading="attaching">挂载</el-button>
      </template>
    </el-dialog>

    <!-- 扩容对话框 -->
    <el-dialog v-model="showResizeDialog" title="扩容云硬盘" width="400px">
      <el-form label-width="80px">
        <el-form-item label="云硬盘">
          <span>{{ currentDisk?.name }}</span>
        </el-form-item>
        <el-form-item label="当前容量">
          <span>{{ currentDisk?.size }} GB</span>
        </el-form-item>
        <el-form-item label="新容量">
          <el-input-number v-model="resizeNewSize" :min="currentDisk?.size || 10" :max="32000" style="width: 100%;" />
          <span style="margin-left: 8px;">GB</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showResizeDialog = false">取消</el-button>
        <el-button type="primary" @click="handleResizeSubmit" :loading="resizing">扩容</el-button>
      </template>
    </el-dialog>

    <!-- 创建快照对话框 -->
    <el-dialog v-model="showSnapshotDialog" title="创建快照" width="400px">
      <el-form label-width="80px">
        <el-form-item label="云硬盘">
          <span>{{ currentDisk?.name }}</span>
        </el-form-item>
        <el-form-item label="快照名称">
          <el-input v-model="snapshotName" placeholder="请输入快照名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showSnapshotDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSnapshotSubmit" :loading="creatingSnapshot">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getCloudDisks, createCloudDisk, deleteCloudDisk, attachCloudDisk, detachCloudDisk, resizeCloudDisk, syncCloudDisks, createCloudSnapshot } from '@/api/storage'
import { getCloudAccounts } from '@/api/cloud-account'
import type { CloudDisk } from '@/api/storage'

const cloudDisks = ref<CloudDisk[]>([])
const loading = ref(false)
const syncing = ref(false)
const cloudAccounts = ref<any[]>([])
const selectedAccountId = ref<number | undefined>()
const filterStatus = ref<string>('')

const pagination = reactive({ currentPage: 1, pageSize: 20, total: 0 })

// Dialog states
const showCreateDialog = ref(false)
const showAttachDialog = ref(false)
const showResizeDialog = ref(false)
const showSnapshotDialog = ref(false)
const creating = ref(false)
const attaching = ref(false)
const resizing = ref(false)
const creatingSnapshot = ref(false)
const currentDisk = ref<CloudDisk | null>(null)
const attachVMId = ref('')
const resizeNewSize = ref(100)
const snapshotName = ref('')

const createForm = ref({
  cloud_account_id: undefined as number | undefined,
  name: '',
  size: 100,
  type: 'cloud_essd',
  zone_id: ''
})

onMounted(() => {
  loadCloudAccounts()
  loadCloudDisks()
})

watch([selectedAccountId, filterStatus, pagination.currentPage, pagination.pageSize], loadCloudDisks)

async function loadCloudAccounts() {
  try {
    const res = await getCloudAccounts({ page: 1, page_size: 100 })
    cloudAccounts.value = res.items || []
  } catch {}
}

async function loadCloudDisks() {
  loading.value = true
  try {
    const res = await getCloudDisks({
      cloud_account_id: selectedAccountId.value,
      status: filterStatus.value,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    cloudDisks.value = res.items || []
    pagination.total = res.total || 0
  } catch (error) {
    ElMessage.warning('获取云硬盘列表失败')
  } finally { loading.value = false }
}

async function handleSync() {
  if (!selectedAccountId.value) {
    ElMessage.warning('请先选择云账号')
    return
  }
  syncing.value = true
  try {
    const res = await syncCloudDisks(selectedAccountId.value)
    ElMessage.success(`同步完成，新增 ${res.count} 条，总计 ${res.total} 条`)
    loadCloudDisks()
  } catch (error) {
    ElMessage.error('同步失败')
  } finally { syncing.value = false }
}

async function handleCreate() {
  if (!createForm.value.cloud_account_id || !createForm.value.name || !createForm.value.zone_id) {
    ElMessage.warning('请填写必要信息')
    return
  }
  creating.value = true
  try {
    await createCloudDisk({
      cloud_account_id: createForm.value.cloud_account_id,
      name: createForm.value.name,
      size: createForm.value.size,
      type: createForm.value.type,
      zone_id: createForm.value.zone_id
    })
    ElMessage.success('云硬盘创建请求已提交')
    showCreateDialog.value = false
    createForm.value = { cloud_account_id: undefined, name: '', size: 100, type: 'cloud_essd', zone_id: '' }
    loadCloudDisks()
  } catch (error) {
    ElMessage.error('创建失败')
  } finally { creating.value = false }
}

function handleAttach(row: CloudDisk) {
  currentDisk.value = row
  attachVMId.value = ''
  showAttachDialog.value = true
}

async function handleAttachSubmit() {
  if (!currentDisk.value || !attachVMId.value) return
  attaching.value = true
  try {
    await attachCloudDisk(currentDisk.value.id, attachVMId.value)
    ElMessage.success('挂载请求已提交')
    showAttachDialog.value = false
    loadCloudDisks()
  } catch {
    ElMessage.error('挂载失败')
  } finally { attaching.value = false }
}

function handleDetach(row: CloudDisk) {
  currentDisk.value = row
  ElMessageBox.confirm('确认卸载该云硬盘？', '卸载确认', { type: 'warning' }).then(async () => {
    try {
      await detachCloudDisk(row.id)
      ElMessage.success('卸载请求已提交')
      loadCloudDisks()
    } catch {
      ElMessage.error('卸载失败')
    }
  }).catch(() => {})
}

function handleResize(row: CloudDisk) {
  currentDisk.value = row
  resizeNewSize.value = row.size
  showResizeDialog.value = true
}

async function handleResizeSubmit() {
  if (!currentDisk.value) return
  resizing.value = true
  try {
    await resizeCloudDisk(currentDisk.value.id, resizeNewSize.value)
    ElMessage.success('扩容请求已提交')
    showResizeDialog.value = false
    loadCloudDisks()
  } catch {
    ElMessage.error('扩容失败')
  } finally { resizing.value = false }
}

function handleSnapshot(row: CloudDisk) {
  currentDisk.value = row
  snapshotName.value = `${row.name}-snapshot-${Date.now()}`
  showSnapshotDialog.value = true
}

async function handleSnapshotSubmit() {
  if (!currentDisk.value || !snapshotName.value) return
  creatingSnapshot.value = true
  try {
    await createCloudSnapshot(currentDisk.value.id, snapshotName.value)
    ElMessage.success('快照创建请求已提交')
    showSnapshotDialog.value = false
  } catch {
    ElMessage.error('创建快照失败')
  } finally { creatingSnapshot.value = false }
}

async function handleDelete(row: CloudDisk) {
  try {
    await ElMessageBox.confirm('确认删除该云硬盘？此操作不可恢复。', '删除确认', { type: 'warning' })
    await deleteCloudDisk(row.id)
    ElMessage.success('删除请求已提交')
    loadCloudDisks()
  } catch {
    // canceled
  }
}

function formatTime(time: string): string {
  return time ? new Date(time).toLocaleString('zh-CN') : ''
}

function getDiskTypeText(type: string): string {
  const map: Record<string, string> = {
    'cloud_essd': 'ESSD',
    'cloud_essd_pl1': 'ESSD PL1',
    'cloud_essd_pl2': 'ESSD PL2',
    'cloud_essd_pl3': 'ESSD PL3',
    'cloud_efficiency': '高效云盘',
    'cloud': '普通云盘',
    'cloud_ssd': 'SSD云盘',
    'local_ssd': '本地SSD'
  }
  return map[type] || type
}

function getStatusType(status: string): string {
  const map: Record<string, string> = {
    'available': 'success',
    'in_use': 'primary',
    'creating': 'warning',
    'deleting': 'danger',
    'attaching': 'warning',
    'detaching': 'warning'
  }
  return map[status] || 'info'
}

function getStatusText(status: string): string {
  const map: Record<string, string> = {
    'available': '可用',
    'in_use': '已挂载',
    'creating': '创建中',
    'deleting': '删除中',
    'attaching': '挂载中',
    'detaching': '卸载中'
  }
  return map[status] || status
}
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
.title { font-size: 18px; font-weight: bold; }
.header-actions { display: flex; gap: 8px; }
</style>