<template>
  <div class="vms-container">
    <div class="page-header">
      <h2>虚拟机管理</h2>
      <div class="toolbar">
        <el-button link type="primary" @click="handleViewMode">
          <el-icon><View /></el-icon>
          查看
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建虚拟机
        </el-button>
        <el-button :disabled="selectedVMs.length === 0" @click="handleBatchStart">
          <el-icon><VideoPlay /></el-icon>
          启动
        </el-button>
        <el-button :disabled="selectedVMs.length === 0" @click="handleBatchStop">
          <el-icon><VideoPause /></el-icon>
          停止
        </el-button>
        <el-button :disabled="selectedVMs.length === 0" @click="handleBatchReboot">
          <el-icon><RefreshRight /></el-icon>
          重启
        </el-button>
        <el-button :disabled="selectedVMs.length === 0" :loading="syncing" @click="syncStatus">
          <el-icon><Refresh /></el-icon>
          同步状态
        </el-button>
        <el-dropdown trigger="click" @command="handleBatchCommand" :disabled="selectedVMs.length === 0">
          <el-button :disabled="selectedVMs.length === 0">
            批量操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="start">批量启动</el-dropdown-item>
              <el-dropdown-item command="stop">批量停止</el-dropdown-item>
              <el-dropdown-item command="reboot">批量重启</el-dropdown-item>
              <el-dropdown-item command="delete" divided>批量删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button @click="handleTags">
          <el-icon><PriceTag /></el-icon>
          标签
        </el-button>
        <el-dropdown trigger="click" @command="handleRemoteControl" :disabled="selectedVMs.length === 0">
          <el-button link type="success" :disabled="selectedVMs.length === 0">
            远程控制<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="vnc">VNC 远程终端</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-dropdown trigger="click" @command="handleMoreCommand">
          <el-button link>
            更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="export">导出列表</el-dropdown-item>
              <el-dropdown-item command="refresh">刷新数据</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="queryForm" @submit.prevent="loadVMs">
        <el-form-item label="云账户">
          <CloudAccountSelector
            v-model:value="queryForm.account_id"
            placeholder="请选择云账户"
            @change="handleAccountChange"
          />
        </el-form-item>
        <el-form-item label="名称/IP">
          <el-input v-model="queryForm.name" placeholder="搜索名称或IP" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryForm.status" placeholder="选择状态" clearable style="width: 140px">
            <el-option label="运行中" value="Running" />
            <el-option label="已停止" value="Stopped" />
            <el-option label="启动中" value="Starting" />
            <el-option label="停止中" value="Stopping" />
            <el-option label="创建中" value="Pending" />
          </el-select>
        </el-form-item>
        <el-form-item label="项目">
          <el-select v-model="queryForm.project_id" placeholder="选择项目" clearable style="width: 140px">
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-select v-model="queryForm.region" placeholder="选择区域" clearable style="width: 140px">
            <el-option label="华东1(杭州)" value="cn-hangzhou" />
            <el-option label="华东2(上海)" value="cn-shanghai" />
            <el-option label="华北2(北京)" value="cn-beijing" />
            <el-option label="华南1(深圳)" value="cn-shenzhen" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="queryForm.platform" placeholder="选择平台" clearable style="width: 140px">
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadVMs">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table
      :data="vms"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column label="名称" min-width="150">
        <template #default="{ row }">
          <el-link type="primary" :underline="false" @click="showDetails(row)">
            {{ row.name }}
          </el-link>
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
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusName(row.status) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="IP" width="200">
        <template #default="{ row }">
          <div v-if="row.private_ip">{{ row.private_ip }}</div>
          <div v-if="row.public_ip">{{ row.public_ip }}</div>
          <span v-if="!row.private_ip && !row.public_ip">-</span>
        </template>
      </el-table-column>

      <el-table-column prop="os_name" label="系统" width="120">
        <template #default="{ row }">
          {{ row.os_name || '未知' }}
        </template>
      </el-table-column>

      <el-table-column prop="billing_method" label="计费" width="100">
        <template #default="{ row }">
          {{ row.billing_method || '按量付费' }}
        </template>
      </el-table-column>

      <el-table-column prop="region_id" label="区域" width="120" />

      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-dropdown trigger="click" @command="(cmd: string) => handleRemoteCommand(cmd, row)">
            <el-button size="small" type="success" link>
              远程控制<el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="vnc">VNC 远程终端</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <el-dropdown trigger="click" @command="(cmd: string) => handleDropdownCommand(cmd, row)">
            <el-button size="small" link>
              更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="start" v-if="row.status === 'Stopped'">启动</el-dropdown-item>
                <el-dropdown-item command="stop" v-if="row.status === 'Running'">停止</el-dropdown-item>
                <el-dropdown-item command="reboot">重启</el-dropdown-item>
                <el-dropdown-item command="details">详情</el-dropdown-item>
                <el-dropdown-item divided command="delete">删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      class="pagination"
    />

    <!-- VM 详情模态框 -->
    <VMModal
      v-model="detailsModalVisible"
      :vm-id="selectedVM?.id || ''"
      :account-id="queryForm.account_id || 0"
      @close="detailsModalVisible = false"
    />

    <!-- VNC 控制台模态框 -->
    <VNCConsole
      v-model="vncModalVisible"
      :vm-id="selectedVM?.id || ''"
      :vm-name="selectedVM?.name || ''"
      :account-id="queryForm.account_id || 0"
      @close="vncModalVisible = false"
    />

    <!-- 创建虚拟机模态框 -->
    <CreateVMModal
      v-model:visible="createModalVisible"
      @success="handleCreateSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, ArrowDown, Refresh, View, VideoPlay, VideoPause, RefreshRight, PriceTag } from '@element-plus/icons-vue'
import { getVMs, vmAction, deleteVM, batchVMAction } from '@/api/compute'
import type { VirtualMachine } from '@/types'
import type { PaginatedResponse } from '@/types/api'
import VMModal from '@/components/vm/VMModal.vue'
import VNCConsole from '@/components/vm/VNCConsole.vue'
import CreateVMModal from '@/components/vm/CreateVMModal.vue'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import {
  getStatusLabel,
  getStatusTagType,
  getProviderLabel,
  getProviderTagType
} from '@/utils/status-mappers'

// 响应式数据
const vms = ref<VirtualMachine[]>([])
const projects = ref<{ id: number; name: string }[]>([])
const loading = ref(false)
const syncing = ref(false)
const detailsModalVisible = ref(false)
const vncModalVisible = ref(false)
const createModalVisible = ref(false)
const tagsModalVisible = ref(false)
const selectedVM = ref<VirtualMachine | null>(null)
const selectedVMs = ref<VirtualMachine[]>([])

// 分页数据
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 查询表单
const queryForm = reactive({
  account_id: null as number | null,
  name: '',
  status: '',
  project_id: null as number | null,
  region: '',
  platform: ''
})

// 使用共享的状态映射函数
const getPlatformLabel = (platform: string): string => getProviderLabel(platform)
const getPlatformType = (platform: string) => getProviderTagType(platform)
const getStatusName = (status: string) => getStatusLabel(status)
const getStatusType = (status: string) => getStatusTagType(status)

const handleAccountChange = (accountId: number | null) => {
  queryForm.account_id = accountId
  loadVMs()
}

const handleSelectionChange = (selection: VirtualMachine[]) => {
  selectedVMs.value = selection
}

const loadVMs = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      size: pageSize.value,
      ...(queryForm.account_id ? { account_id: queryForm.account_id } : {}),
      ...(queryForm.name ? { name: queryForm.name } : {}),
      ...(queryForm.status ? { status: queryForm.status } : {}),
      ...(queryForm.platform ? { platform: queryForm.platform } : {})
    }

    const res = await getVMs(params) as PaginatedResponse<VirtualMachine>
    vms.value = res.items ?? []
    total.value = res.total ?? 0
  } catch (e) {
    console.error(e)
    ElMessage.error('加载虚拟机列表失败')
  } finally {
    loading.value = false
  }
  }

const syncStatus = async () => {
  syncing.value = true
  try {
    await loadVMs()
    ElMessage.success('状态已同步')
  } finally {
    syncing.value = false
  }
}

const handleCreate = () => {
  createModalVisible.value = true
  }

const handleViewMode = () => {
  ElMessage.info('切换查看模式')
  }

const handleTags = () => {
  if (selectedVMs.value.length === 0) {
    ElMessage.warning('请先选择虚拟机')
    return
  }
  tagsModalVisible.value = true
  }

const handleRemoteControl = (command: string) => {
  if (selectedVMs.value.length === 0) return
  if (command === 'vnc' && selectedVMs.value.length === 1) {
    openVNC(selectedVMs.value[0])
  } else if (command === 'vnc') {
    ElMessage.warning('VNC 远程终端只能对单个虚拟机操作')
  }
  }

const handleMoreCommand = (command: string) => {
  if (command === 'export') {
    ElMessage.success('导出功能开发中')
  } else if (command === 'refresh') {
    loadVMs()
    ElMessage.success('数据已刷新')
  }
  }

const handleBatchStart = async () => {
  await handleBatchCommand('start')
  }

const handleBatchStop = async () => {
  await handleBatchCommand('stop')
  }

const handleBatchReboot = async () => {
  await handleBatchCommand('reboot')
  }

// 批量操作
const handleBatchCommand = async (command: string) => {
  if (selectedVMs.value.length === 0) {
    ElMessage.warning('请先选择虚拟机')
    return
  }

  const actionNames: Record<string, string> = {
    start: '启动',
    stop: '停止',
    reboot: '重启',
    delete: '删除'
  }

  try {
    await ElMessageBox.confirm(
      `确定要批量${actionNames[command]}选中的 ${selectedVMs.value.length} 个虚拟机吗？`,
      '批量操作确认',
      { type: 'warning' }
    )

    if (!queryForm.account_id) {
      ElMessage.error('请先选择云账户')
      return
    }

    const vmIds = selectedVMs.value.map(vm => vm.id)
    const res = await batchVMAction({
      vm_ids: vmIds,
      account_id: queryForm.account_id,
      action: command as any
    })

    ElMessage.success(res.message || `批量${actionNames[command]}完成`)
    selectedVMs.value = []
    loadVMs()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '批量操作失败')
    }
  }
}

// 远程控制下拉
const handleRemoteCommand = (command: string, row: VirtualMachine) => {
  if (command === 'vnc') {
    openVNC(row)
  }
}

const handleDropdownCommand = async (command: string, row: VirtualMachine) => {
  if (command === 'details') {
    showDetails(row)
    return
  }

  if (command === 'delete') {
    try {
      await ElMessageBox.confirm('确定要删除该虚拟机吗？', '提示', { type: 'warning' })
      const accountId = row.cloud_account_id || queryForm.account_id
      if (accountId) {
        await deleteVM(row.id, accountId)
        ElMessage.success('删除成功')
        loadVMs()
      }
    } catch (e: any) {
      if (e !== 'cancel') {
        ElMessage.error(`删除失败: ${e.message}`)
      }
    }
  } else {
    try {
      await ElMessageBox.confirm(`确定要${command === 'start' ? '启动' : command === 'stop' ? '停止' : '重启'}该虚拟机吗？`, '提示', { type: 'warning' })
      const accountId = row.cloud_account_id || queryForm.account_id
      if (accountId) {
        await vmAction(row.id, accountId, command as any)
        ElMessage.success('操作成功')
        loadVMs()
      }
    } catch (e: any) {
      if (e !== 'cancel') {
        ElMessage.error(`操作失败: ${e.message}`)
      }
    }
  }
}

const showDetails = (row: VirtualMachine) => {
  selectedVM.value = row
  detailsModalVisible.value = true
}

const openVNC = (row: VirtualMachine) => {
  selectedVM.value = row
  vncModalVisible.value = true
}

const resetQuery = () => {
  queryForm.account_id = null
  queryForm.name = ''
  queryForm.status = ''
  queryForm.platform = ''
  currentPage.value = 1
  loadVMs()
  }

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadVMs()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadVMs()
}

const handleCreateSuccess = (vm: VirtualMachine) => {
  ElMessage.success(`${vm.name} 创建成功`)
  loadVMs()
}

const fetchProjects = async () => {
  // Mock 项目数据，实际应从 API 获取
  projects.value = [
    { id: 1, name: 'system' },
    { id: 2, name: '项目A' },
    { id: 3, name: '项目B' }
  ]
}

onMounted(() => {
  loadVMs()
  fetchProjects()
})
</script>

<style scoped>
.vms-container {
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
  gap: 10px;
  align-items: center;
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