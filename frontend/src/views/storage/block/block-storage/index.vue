<template>
  <div class="block-storage-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">块存储列表</span>
          <div class="toolbar">
            <el-button @click="handleView">查看</el-button>
            <el-dropdown :disabled="!hasSelection" @command="handleBatchCommand">
              <el-button>
                批量操作 <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="delete">批量删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button :disabled="!hasSelection" @click="handleTags">标签</el-button>
          </div>
        </div>
      </template>

      <!-- Search Filters -->
      <el-form :inline="true" class="search-form">
        <el-form-item label="名称">
          <el-input v-model="searchParams.name" placeholder="请输入名称" clearable @clear="loadBlockStorage" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchParams.status" placeholder="请选择状态" clearable @clear="loadBlockStorage">
            <el-option label="运行中" value="active" />
            <el-option label="离线" value="offline" />
            <el-option label="创建中" value="creating" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="searchParams.platform" placeholder="请选择平台" clearable @clear="loadBlockStorage">
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
            <el-option label="本地IDC" value="idc" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-input v-model="searchParams.region" placeholder="请输入区域" clearable @clear="loadBlockStorage" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadBlockStorage">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>

      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="宿主机" name="host">
          <el-table :data="hostStorage" v-loading="loading" @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55" />
            <el-table-column label="名称" width="180">
              <template #default="{ row }">
                <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="enabled" label="启用状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '已启用' : '未启用' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="physical_capacity" label="物理容量" width="150">
              <template #default="{ row }">
                {{ row.physical_capacity }} GB
              </template>
            </el-table-column>
            <el-table-column prop="virtual_capacity" label="虚拟容量" width="150">
              <template #default="{ row }">
                {{ row.virtual_capacity }} GB
              </template>
            </el-table-column>
            <el-table-column prop="platform" label="平台" width="100" />
            <el-table-column prop="domain" label="所属域" width="150" />
            <el-table-column prop="region" label="区域" width="120" />
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <el-button size="small" @click="handleEdit(row)">编辑</el-button>
                <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            v-model:current-page="searchParams.page"
            v-model:page-size="searchParams.page_size"
            :total="total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            @size-change="loadBlockStorage"
            @current-change="loadBlockStorage"
          />
        </el-tab-pane>
        <el-tab-pane label="物理机" name="physical_host">
          <el-table :data="physicalHostStorage" v-loading="loading" @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55" />
            <el-table-column label="名称" width="180">
              <template #default="{ row }">
                <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="enabled" label="启用状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '已启用' : '未启用' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="physical_capacity" label="物理容量" width="150">
              <template #default="{ row }">
                {{ row.physical_capacity }} GB
              </template>
            </el-table-column>
            <el-table-column prop="virtual_capacity" label="虚拟容量" width="150">
              <template #default="{ row }">
                {{ row.virtual_capacity }} GB
              </template>
            </el-table-column>
            <el-table-column prop="platform" label="平台" width="100" />
            <el-table-column prop="domain" label="所属域" width="150" />
            <el-table-column prop="region" label="区域" width="120" />
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <el-button size="small" @click="handleEdit(row)">编辑</el-button>
                <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            v-model:current-page="searchParams.page"
            v-model:page-size="searchParams.page_size"
            :total="total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            @size-change="loadBlockStorage"
            @current-change="loadBlockStorage"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="块存储详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedStorage?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedStorage?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedStorage?.status }}</el-descriptions-item>
        <el-descriptions-item label="启用状态">{{ selectedStorage?.enabled ? '已启用' : '未启用' }}</el-descriptions-item>
        <el-descriptions-item label="物理容量">{{ selectedStorage?.physical_capacity }} GB</el-descriptions-item>
        <el-descriptions-item label="虚拟容量">{{ selectedStorage?.virtual_capacity }} GB</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedStorage?.platform }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ selectedStorage?.domain }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedStorage?.region }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ selectedStorage?.type }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedStorage?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'

interface BlockStorage {
  id: string
  name: string
  status: string
  enabled: boolean
  physical_capacity: number
  virtual_capacity: number
  platform: string
  domain: string
  region: string
  type: string
  created_at: string
}

const hostStorage = ref<BlockStorage[]>([])
const physicalHostStorage = ref<BlockStorage[]>([])
const loading = ref(false)
const activeTab = ref('host')
const total = ref(0)
const selectedRows = ref<BlockStorage[]>([])
const detailDialogVisible = ref(false)
const selectedStorage = ref<BlockStorage | null>(null)

const searchParams = ref({
  page: 1,
  page_size: 10,
  name: '',
  status: '',
  platform: '',
  region: ''
})

const hasSelection = computed(() => selectedRows.value.length > 0)

const handleSelectionChange = (rows: BlockStorage[]) => {
  selectedRows.value = rows
}

const handleTabChange = () => {
  searchParams.value.page = 1
  loadBlockStorage()
}

const resetSearch = () => {
  searchParams.value = {
    page: 1,
    page_size: 10,
    name: '',
    status: '',
    platform: '',
    region: ''
  }
  loadBlockStorage()
}

const handleView = () => {
  if (selectedRows.value.length === 1) {
    viewDetail(selectedRows.value[0])
  } else {
    ElMessage.warning('请选择一条记录查看')
  }
}

const handleBatchCommand = async (command: string) => {
  if (!hasSelection.value) return
  if (command === 'delete') {
    try {
      await ElMessageBox.confirm(`确认删除选中的 ${selectedRows.value.length} 个块存储？`, '警告', { type: 'warning' })
      hostStorage.value = hostStorage.value.filter(s => !selectedRows.value.some(r => r.id === s.id))
      physicalHostStorage.value = physicalHostStorage.value.filter(s => !selectedRows.value.some(r => r.id === s.id))
      ElMessage.success('删除成功')
      loadBlockStorage()
    } catch (e) {
      console.error(e)
    }
  }
}

const handleTags = () => {
  if (!hasSelection.value) {
    ElMessage.warning('请选择需要设置标签的块存储')
    return
  }
  ElMessage.info('标签功能开发中')
}

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
    case 'online':
      return 'success'
    case 'pending':
    case 'creating':
      return 'warning'
    case 'offline':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const loadBlockStorage = async () => {
  loading.value = true
  try {
    // Mock data - 宿主机
    hostStorage.value = [
      {
        id: 'bs-host-1',
        name: '宿主机块存储 1',
        status: 'Active',
        enabled: true,
        physical_capacity: 1000,
        virtual_capacity: 2000,
        platform: '阿里云',
        domain: 'Default Domain',
        region: 'cn-hangzhou',
        type: '宿主机',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'bs-host-2',
        name: '宿主机块存储 2',
        status: 'Active',
        enabled: true,
        physical_capacity: 2000,
        virtual_capacity: 4000,
        platform: '阿里云',
        domain: 'Default Domain',
        region: 'cn-shanghai',
        type: '宿主机',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'bs-host-3',
        name: '宿主机块存储 3',
        status: 'Offline',
        enabled: false,
        physical_capacity: 500,
        virtual_capacity: 1000,
        platform: '阿里云',
        domain: 'Domain A',
        region: 'cn-beijing',
        type: '宿主机',
        created_at: '2024-01-03 10:00:00'
      }
    ]

    // Mock data - 物理机
    physicalHostStorage.value = [
      {
        id: 'bs-ph-1',
        name: '物理机块存储 1',
        status: 'Active',
        enabled: true,
        physical_capacity: 5000,
        virtual_capacity: 10000,
        platform: '本地IDC',
        domain: 'Default Domain',
        region: '本地机房',
        type: '物理机',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'bs-ph-2',
        name: '物理机块存储 2',
        status: 'Creating',
        enabled: true,
        physical_capacity: 3000,
        virtual_capacity: 6000,
        platform: '本地IDC',
        domain: 'Domain B',
        region: '本地机房',
        type: '物理机',
        created_at: '2024-01-02 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    hostStorage.value = []
    physicalHostStorage.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: BlockStorage) => {
  selectedStorage.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row: BlockStorage) => {
  ElMessage.info(`编辑块存储: ${row.name}`)
}

const handleDelete = async (row: BlockStorage) => {
  try {
    await ElMessageBox.confirm(`确认删除块存储 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    hostStorage.value = hostStorage.value.filter(s => s.id !== row.id)
    physicalHostStorage.value = physicalHostStorage.value.filter(s => s.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadBlockStorage()
})
</script>

<style scoped>
.block-storage-page {
  height: 100%;
}

.page-card {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 18px;
  font-weight: bold;
}

.toolbar {
  display: flex;
  gap: 8px;
}

.search-form {
  margin-bottom: 16px;
}
</style>