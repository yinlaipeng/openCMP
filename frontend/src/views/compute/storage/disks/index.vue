<template>
  <div class="disks-container">
    <div class="page-header">
      <h2>硬盘</h2>
      <div class="toolbar">
        <el-button link type="primary" @click="handleView">
          <el-icon><View /></el-icon>
          查看
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建
        </el-button>
        <el-dropdown trigger="click" @command="handleBatchCommand" :disabled="selectedDisks.length === 0">
          <el-button :disabled="selectedDisks.length === 0">
            批量操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="attach">批量挂载</el-dropdown-item>
              <el-dropdown-item command="detach">批量卸载</el-dropdown-item>
              <el-dropdown-item command="delete" divided>批量删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button @click="handleTags">
          <el-icon><PriceTag /></el-icon>
          标签
        </el-button>
        <el-button link :disabled="selectedDisks.length !== 1" @click="handleResize">
          <el-icon><FullScreen /></el-icon>
          扩容
        </el-button>
        <el-dropdown trigger="click" @command="handleMoreCommand">
          <el-button link>
            更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="sync">同步状态</el-dropdown-item>
              <el-dropdown-item command="refresh">刷新数据</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- 公有云/私有云 Tabs -->
    <el-tabs v-model="activeTab" class="disk-tabs" @tab-change="handleTabChange">
      <el-tab-pane label="公有云" name="public">
        <el-card class="filter-card">
          <el-form :inline="true" :model="filters" @submit.prevent="fetchData">
            <el-form-item label="名称">
              <el-input v-model="filters.name" placeholder="搜索硬盘名称" clearable style="width: 180px" />
            </el-form-item>
            <el-form-item label="状态">
              <el-select v-model="filters.status" placeholder="选择状态" clearable style="width: 120px">
                <el-option label="可用" value="available" />
                <el-option label="使用中" value="in_use" />
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
            <el-form-item label="类型">
              <el-select v-model="filters.type" placeholder="选择类型" clearable style="width: 120px">
                <el-option label="系统盘" value="system" />
                <el-option label="数据盘" value="data" />
              </el-select>
            </el-form-item>
            <el-form-item label="区域">
              <el-select v-model="filters.region" placeholder="选择区域" clearable style="width: 140px">
                <el-option label="华东1(杭州)" value="cn-hangzhou" />
                <el-option label="华东2(上海)" value="cn-shanghai" />
                <el-option label="华北2(北京)" value="cn-beijing" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="fetchData">查询</el-button>
              <el-button @click="resetFilters">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-table
          :data="disks"
          v-loading="loading"
          style="width: 100%"
          row-key="id"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column label="名称" min-width="150">
            <template #default="{ row }">
              <el-link type="primary" :underline="false" @click="handleDetails(row)">
                {{ row.name }}
              </el-link>
            </template>
          </el-table-column>
          <el-table-column label="标签" width="100">
            <template #default="{ row }">
              <el-tag v-for="tag in (row.tags || [])" :key="tag.key" size="small" class="tag-item">
                {{ tag.key }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="容量" width="100">
            <template #default="{ row }">
              {{ row.size }} GB
            </template>
          </el-table-column>
          <el-table-column prop="max_iops" label="最大IOPS" width="100" />
          <el-table-column prop="disk_format" label="格式" width="80" />
          <el-table-column prop="storage_type" label="存储类型" width="100" />
          <el-table-column prop="type" label="磁盘类型" width="100">
            <template #default="{ row }">
              <el-tag size="small" :type="row.type === 'system' ? 'primary' : 'info'">
                {{ row.type === 'system' ? '系统盘' : '数据盘' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="挂载状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.vm_id ? 'success' : 'info'" size="small">
                {{ row.vm_id ? '已挂载' : '未挂载' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="vm_name" label="虚拟机" width="150" />
          <el-table-column prop="device_name" label="设备名" width="100" />
          <el-table-column prop="primary_storage" label="主存储" width="100" />
          <el-table-column prop="billing_type" label="计费方式" width="100" />
          <el-table-column prop="medium_type" label="媒体类型" width="100" />
          <el-table-column label="关机自动重置" width="100">
            <template #default="{ row }">
              <el-tag :type="row.shutdown_reset ? 'warning' : 'info'" size="small">
                {{ row.shutdown_reset ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="平台/云账号" width="150">
            <template #default="{ row }">
              <div class="platform-cell">
                <el-tag size="small" :type="getPlatformType(row.provider_type)">
                  {{ getPlatformLabel(row.provider_type) }}
                </el-tag>
                <span class="account-name">{{ row.account_name || '-' }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="region_id" label="区域" width="120" />
          <el-table-column prop="zone_id" label="可用区" width="120" />
          <el-table-column prop="project_name" label="项目" width="120" />
          <el-table-column prop="created_at" label="创建时间" width="150" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="160" fixed="right">
            <template #default="{ row }">
              <el-dropdown trigger="click" @command="(cmd: string) => handleActionCommand(cmd, row)">
                <el-button size="small" link type="primary">
                  操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="details">查看详情</el-dropdown-item>
                    <el-dropdown-item command="attach" :disabled="row.vm_id">挂载</el-dropdown-item>
                    <el-dropdown-item command="detach" :disabled="!row.vm_id">卸载</el-dropdown-item>
                    <el-dropdown-item command="resize">扩容</el-dropdown-item>
                    <el-dropdown-item command="snapshot">创建快照</el-dropdown-item>
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
            <template #title>私有云硬盘管理功能开发中</template>
          </el-alert>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 创建硬盘弹窗 -->
    <el-dialog
      title="创建硬盘"
      v-model="createDialogVisible"
      width="600px"
    >
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="120px">
        <el-form-item label="云账号" prop="cloud_account_id">
          <el-select v-model="createForm.cloud_account_id" placeholder="选择云账号" style="width: 100%" @change="handleAccountChange">
            <el-option v-for="acc in cloudAccounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入硬盘名称" />
        </el-form-item>
        <el-form-item label="容量(GB)" prop="size">
          <el-input-number v-model="createForm.size" :min="10" :max="32768" style="width: 100%" />
        </el-form-item>
        <el-form-item label="硬盘类型" prop="disk_type">
          <el-select v-model="createForm.disk_type" placeholder="选择硬盘类型" style="width: 100%">
            <el-option label="SSD云盘" value="ssd" />
            <el-option label="高效云盘" value="efficiency" />
            <el-option label="普通云盘" value="normal" />
            <el-option label="ESSD云盘" value="essd" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域" prop="region_id">
          <el-select v-model="createForm.region_id" placeholder="选择区域" style="width: 100%" @change="handleRegionChange">
            <el-option v-for="r in regions" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="可用区" prop="zone_id">
          <el-select v-model="createForm.zone_id" placeholder="选择可用区" style="width: 100%">
            <el-option v-for="z in zones" :key="z.id" :label="z.name" :value="z.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="项目">
          <el-select v-model="createForm.project_id" placeholder="选择项目" clearable style="width: 100%">
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="createForm.description" type="textarea" :rows="2" placeholder="请输入描述（可选）" />
        </el-form-item>
        <el-form-item label="标签">
          <div class="tag-editor">
            <div v-for="(tag, index) in createForm.tags" :key="index" class="tag-row">
              <el-input v-model="tag.key" placeholder="标签键" style="width: 150px" />
              <el-input v-model="tag.value" placeholder="标签值" style="width: 150px" />
              <el-button type="danger" size="small" @click="removeCreateTag(index)">删除</el-button>
            </div>
            <el-button type="primary" size="small" @click="addCreateTag">添加标签</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCreate" :loading="creating">创建</el-button>
      </template>
    </el-dialog>

    <!-- 详情弹窗 -->
    <el-dialog
      title="硬盘详情"
      v-model="detailDialogVisible"
      width="800px"
    >
      <el-tabs v-model="detailTab" v-if="selectedDisk">
        <el-tab-pane label="基础信息" name="basic">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">{{ selectedDisk.disk_id || selectedDisk.id }}</el-descriptions-item>
            <el-descriptions-item label="名称">{{ selectedDisk.name }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(selectedDisk.status)">{{ getStatusLabel(selectedDisk.status) }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="容量">{{ selectedDisk.size }} GB</el-descriptions-item>
            <el-descriptions-item label="类型">{{ selectedDisk.type === 'system' ? '系统盘' : '数据盘' }}</el-descriptions-item>
            <el-descriptions-item label="存储类型">{{ selectedDisk.storage_type || '-' }}</el-descriptions-item>
            <el-descriptions-item label="挂载状态">{{ selectedDisk.vm_id ? '已挂载' : '未挂载' }}</el-descriptions-item>
            <el-descriptions-item label="虚拟机">{{ selectedDisk.vm_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="设备名">{{ selectedDisk.device_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="平台">
              <el-tag size="small" :type="getPlatformType(selectedDisk.provider_type)">
                {{ getPlatformLabel(selectedDisk.provider_type) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="云账号">{{ selectedDisk.account_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="可用区">{{ selectedDisk.zone_id || '-' }}</el-descriptions-item>
            <el-descriptions-item label="区域">{{ selectedDisk.region_id || '-' }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ selectedDisk.created_at }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="标签" name="tags">
          <el-table :data="selectedDiskTags" style="width: 100%">
            <el-table-column prop="key" label="标签键" width="200" />
            <el-table-column prop="value" label="标签值" />
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button size="small" link type="danger" @click="removeTag(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-button type="primary" size="small" style="margin-top: 10px" @click="addTag">添加标签</el-button>
        </el-tab-pane>
        <el-tab-pane label="操作日志" name="logs">
          <el-table :data="operationLogs" style="width: 100%">
            <el-table-column prop="time" label="时间" width="180" />
            <el-table-column prop="action" label="操作" width="120" />
            <el-table-column prop="operator" label="操作人" width="120" />
            <el-table-column prop="result" label="结果">
              <template #default="{ row }">
                <el-tag size="small" :type="row.result === '成功' ? 'success' : 'danger'">{{ row.result }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>

    <!-- 扩容弹窗 -->
    <el-dialog
      title="硬盘扩容"
      v-model="resizeDialogVisible"
      width="400px"
    >
      <el-form :model="resizeForm" ref="resizeFormRef" label-width="100px">
        <el-form-item label="当前容量">
          <span>{{ resizeForm.current_size }} GB</span>
        </el-form-item>
        <el-form-item label="新容量(GB)" prop="new_size">
          <el-input-number v-model="resizeForm.new_size" :min="resizeForm.current_size + 1" :max="32768" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resizeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmResize" :loading="resizing">确认扩容</el-button>
      </template>
    </el-dialog>

    <!-- 挂载弹窗 -->
    <el-dialog
      title="挂载硬盘"
      v-model="attachDialogVisible"
      width="400px"
    >
      <el-form :model="attachForm" ref="attachFormRef" label-width="100px">
        <el-form-item label="硬盘">{{ attachForm.disk_name }}</el-form-item>
        <el-form-item label="虚拟机" prop="vm_id">
          <el-select v-model="attachForm.vm_id" placeholder="选择虚拟机" style="width: 100%">
            <el-option v-for="vm in vms" :key="vm.id" :label="vm.name" :value="vm.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="attachDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAttach" :loading="attaching">确认挂载</el-button>
      </template>
    </el-dialog>

    <!-- 创建快照弹窗 -->
    <el-dialog
      title="创建快照"
      v-model="snapshotDialogVisible"
      width="400px"
    >
      <el-form :model="snapshotForm" ref="snapshotFormRef" label-width="100px">
        <el-form-item label="硬盘">{{ snapshotForm.disk_name }}</el-form-item>
        <el-form-item label="快照名称" prop="name">
          <el-input v-model="snapshotForm.name" placeholder="请输入快照名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="snapshotDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmSnapshot" :loading="creatingSnapshot">创建快照</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Plus, ArrowDown, PriceTag, FullScreen } from '@element-plus/icons-vue'
import {
  getCloudDisks,
  createCloudDisk,
  attachCloudDisk,
  detachCloudDisk,
  resizeCloudDisk,
  deleteCloudDisk,
  batchDeleteCloudDisks,
  createCloudDiskSnapshot,
  type CloudDisk,
  type DiskListParams,
  type CreateDiskParams
} from '@/api/storage'

// Data
const loading = ref(false)
const creating = ref(false)
const resizing = ref(false)
const attaching = ref(false)
const creatingSnapshot = ref(false)
const disks = ref<CloudDisk[]>([])
const selectedDisks = ref<CloudDisk[]>([])
const cloudAccounts = ref<any[]>([])
const projects = ref<any[]>([])
const vms = ref<any[]>([])
const zones = ref<any[]>([])
const regions = ref<any[]>([])

const activeTab = ref('public')
const detailTab = ref('basic')
const createDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const resizeDialogVisible = ref(false)
const attachDialogVisible = ref(false)
const snapshotDialogVisible = ref(false)

const selectedDisk = ref<CloudDisk | null>(null)
const selectedDiskTags = ref<{ key: string; value: string }[]>([])
const operationLogs = ref<{ time: string; action: string; operator: string; result: string }[]>([])

const filters = reactive({
  name: '',
  status: '',
  platform: '',
  type: '',
  region: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const createForm = reactive({
  cloud_account_id: '',
  name: '',
  size: 100,
  disk_type: 'ssd',
  region_id: '',
  zone_id: '',
  project_id: '',
  description: '',
  tags: [] as { key: string; value: string }[]
})

const resizeForm = reactive({
  disk_id: '',
  disk_name: '',
  current_size: 0,
  new_size: 0
})

const attachForm = reactive({
  disk_id: '',
  disk_name: '',
  vm_id: ''
})

const snapshotForm = reactive({
  disk_id: '',
  disk_name: '',
  name: ''
})

const createRules = {
  cloud_account_id: [{ required: true, message: '请选择云账号', trigger: 'change' }],
  name: [{ required: true, message: '请输入硬盘名称', trigger: 'blur' }],
  size: [{ required: true, message: '请输入容量', trigger: 'blur' }],
  region_id: [{ required: true, message: '请选择区域', trigger: 'change' }],
  zone_id: [{ required: true, message: '请选择可用区', trigger: 'change' }]
}

// Methods
const getStatusType = (status: string) => {
  switch (status) {
    case 'available': return 'success'
    case 'in_use': return 'primary'
    case 'creating': return 'warning'
    case 'error': return 'danger'
    default: return 'info'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'available': return '可用'
    case 'in_use': return '使用中'
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
    aws: 'success',
    azure: 'info'
  }
  return types[platform] || 'info'
}

const handleSelectionChange = (selection: CloudDisk[]) => {
  selectedDisks.value = selection
}

const handleTabChange = () => {
  pagination.page = 1
  fetchData()
}

const resetFilters = () => {
  filters.name = ''
  filters.status = ''
  filters.platform = ''
  filters.type = ''
  filters.region = ''
  pagination.page = 1
  fetchData()
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: DiskListParams = {
      page: pagination.page,
      page_size: pagination.pageSize,
      name: filters.name || undefined,
      status: filters.status || undefined,
      platform: filters.platform || undefined,
      type: filters.type || undefined,
      region: filters.region || undefined
    }
    const res = await getCloudDisks(params)
    disks.value = res.items
    pagination.total = res.total
  } catch (error) {
    console.error('Failed to fetch disks:', error)
    ElMessage.error('获取硬盘列表失败')
  } finally {
    loading.value = false
  }
}

const fetchCloudAccounts = async () => {
  cloudAccounts.value = [{ id: 1, name: '阿里云账号1' }, { id: 2, name: '腾讯云账号1' }]
}

const fetchProjects = async () => {
  projects.value = [{ id: 1, name: '项目A' }, { id: 2, name: '项目B' }]
}

const fetchVMs = async () => {
  vms.value = [{ id: 'vm-001', name: 'web-server-01' }, { id: 'vm-002', name: 'db-server-01' }]
}

const fetchRegions = async () => {
  regions.value = [
    { id: 'cn-hangzhou', name: '华东1(杭州)' },
    { id: 'cn-shanghai', name: '华东2(上海)' },
    { id: 'cn-beijing', name: '华北2(北京)' },
    { id: 'cn-shenzhen', name: '华南1(深圳)' }
  ]
}

const handleAccountChange = () => {
  // 根据云账号加载区域
  createForm.region_id = ''
  createForm.zone_id = ''
  fetchRegions()
}

const handleRegionChange = () => {
  // 根据区域加载可用区
  createForm.zone_id = ''
  zones.value = [
    { id: 'zone-a', name: '可用区A' },
    { id: 'zone-b', name: '可用区B' },
    { id: 'zone-c', name: '可用区C' }
  ]
}

const addCreateTag = () => {
  createForm.tags.push({ key: '', value: '' })
}

const removeCreateTag = (index: number) => {
  createForm.tags.splice(index, 1)
}

const handleView = () => {
  ElMessage.info('视图切换功能开发中')
}

const handleTags = () => {
  ElMessage.info('标签管理功能开发中')
}

const handleCreate = () => {
  Object.assign(createForm, {
    cloud_account_id: '',
    name: '',
    size: 100,
    disk_type: 'ssd',
    region_id: '',
    zone_id: '',
    project_id: '',
    description: '',
    tags: []
  })
  createDialogVisible.value = true
}

const confirmCreate = async () => {
  creating.value = true
  try {
    const params: CreateDiskParams = {
      cloud_account_id: Number(createForm.cloud_account_id),
      name: createForm.name,
      size: createForm.size,
      type: createForm.disk_type,
      zone_id: createForm.zone_id,
      project_id: createForm.project_id ? Number(createForm.project_id) : undefined
    }
    await createCloudDisk(params)
    ElMessage.success('硬盘创建任务已提交')
    createDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('创建失败')
  } finally {
    creating.value = false
  }
}

const handleBatchCommand = async (command: string) => {
  if (selectedDisks.value.length === 0) return
  const actionNames = { attach: '挂载', detach: '卸载', delete: '删除' }
  try {
    await ElMessageBox.confirm(
      `确定要批量${actionNames[command]}选中的 ${selectedDisks.value.length} 个硬盘吗？`,
      '批量操作确认',
      { type: 'warning' }
    )
    if (command === 'delete') {
      const ids = selectedDisks.value.map(d => d.id)
      await batchDeleteCloudDisks(ids)
      ElMessage.success(`批量删除完成`)
    } else {
      ElMessage.info(`批量${actionNames[command]}功能开发中`)
    }
    selectedDisks.value = []
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('批量操作失败')
  }
}

const handleMoreCommand = (command: string) => {
  if (command === 'sync') {
    ElMessage.info('同步状态功能开发中')
  } else if (command === 'refresh') {
    fetchData()
  }
}

const handleResize = () => {
  if (selectedDisks.value.length !== 1) return
  const disk = selectedDisks.value[0]
  Object.assign(resizeForm, {
    disk_id: disk.disk_id,
    disk_name: disk.name,
    current_size: disk.size,
    new_size: disk.size + 10
  })
  resizeDialogVisible.value = true
}

const confirmResize = async () => {
  resizing.value = true
  try {
    await resizeCloudDisk(resizeForm.disk_id, resizeForm.new_size)
    ElMessage.success('扩容任务已提交')
    resizeDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('扩容失败')
  } finally {
    resizing.value = false
  }
}

const handleActionCommand = (command: string, row: CloudDisk) => {
  switch (command) {
    case 'details': handleDetails(row); break
    case 'attach': handleAttach(row); break
    case 'detach': handleDetach(row); break
    case 'resize': handleResizeSingle(row); break
    case 'snapshot': handleSnapshot(row); break
    case 'delete': handleDelete(row); break
  }
}

const handleDetails = (row: CloudDisk) => {
  selectedDisk.value = row
  detailTab.value = 'basic'
  selectedDiskTags.value = [
    { key: 'environment', value: 'production' },
    { key: 'owner', value: 'team-dev' }
  ]
  operationLogs.value = [
    { time: row.created_at || '2024-01-01 10:00:00', action: '创建硬盘', operator: 'admin', result: '成功' },
    { time: '2024-01-02 14:30:00', action: '挂载', operator: 'admin', result: '成功' }
  ]
  detailDialogVisible.value = true
}

const handleAttach = (row: CloudDisk) => {
  Object.assign(attachForm, {
    disk_id: row.disk_id,
    disk_name: row.name,
    vm_id: ''
  })
  attachDialogVisible.value = true
}

const confirmAttach = async () => {
  attaching.value = true
  try {
    await attachCloudDisk(attachForm.disk_id, attachForm.vm_id)
    ElMessage.success('挂载任务已提交')
    attachDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('挂载失败')
  } finally {
    attaching.value = false
  }
}

const handleDetach = async (row: CloudDisk) => {
  try {
    await ElMessageBox.confirm(`确定要卸载硬盘 "${row.name}" 吗？`, '卸载确认', { type: 'warning' })
    await detachCloudDisk(row.id)
    ElMessage.success('卸载任务已提交')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('卸载失败')
  }
}

const handleResizeSingle = (row: CloudDisk) => {
  Object.assign(resizeForm, {
    disk_id: row.disk_id,
    disk_name: row.name,
    current_size: row.size,
    new_size: row.size + 10
  })
  resizeDialogVisible.value = true
}

const handleSnapshot = (row: CloudDisk) => {
  Object.assign(snapshotForm, {
    disk_id: row.id,
    disk_name: row.name,
    name: `snap-${row.name}-${new Date().toISOString().slice(0, 10)}`
  })
  snapshotDialogVisible.value = true
}

const confirmSnapshot = async () => {
  creatingSnapshot.value = true
  try {
    await createCloudDiskSnapshot({ disk_id: snapshotForm.disk_id, name: snapshotForm.name })
    ElMessage.success('快照创建任务已提交')
    snapshotDialogVisible.value = false
  } catch (error) {
    ElMessage.error('创建快照失败')
  } finally {
    creatingSnapshot.value = false
  }
}

const handleDelete = async (row: CloudDisk) => {
  try {
    await ElMessageBox.confirm(`确定要删除硬盘 "${row.name}" 吗？`, '删除警告', { type: 'warning' })
    await deleteCloudDisk(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

const addTag = () => {
  ElMessage.info('添加标签功能开发中')
}

const removeTag = (tag: { key: string; value: string }) => {
  const index = selectedDiskTags.value.findIndex(t => t.key === tag.key)
  if (index > -1) {
    selectedDiskTags.value.splice(index, 1)
    ElMessage.success('标签已删除')
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
  fetchCloudAccounts()
  fetchProjects()
  fetchVMs()
})
</script>

<style scoped>
.disks-container { padding: 20px; }

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
.disk-tabs { margin-bottom: 20px; }

.platform-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.account-name { font-size: 12px; color: var(--el-text-color-secondary); }

.pagination { margin-top: 20px; justify-content: flex-end; }
.tag-item { margin-right: 4px; margin-bottom: 2px; }

.tag-editor {
  width: 100%;
}

.tag-row {
  display: flex;
  gap: 10px;
  margin-bottom: 8px;
  align-items: center;
}
</style>