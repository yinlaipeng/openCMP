<template>
  <div class="vms-container">
    <div class="page-header">
      <h2>虚拟机管理</h2>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        创建虚拟机
      </el-button>
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
        <el-form-item label="名称">
          <el-input v-model="queryForm.name" placeholder="虚拟机名称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryForm.status" placeholder="选择状态" clearable>
            <el-option label="运行中" value="Running" />
            <el-option label="已停止" value="Stopped" />
            <el-option label="启动中" value="Starting" />
            <el-option label="停止中" value="Stopping" />
            <el-option label="创建中" value="Pending" />
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
    >
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
          <el-button size="small" type="success" @click="openVNC(row)">远程控制</el-button>
          <el-dropdown trigger="click" @command="(cmd: string) => handleDropdownCommand(cmd, row)">
            <el-button size="small">
              操作<el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="start" v-if="row.status === 'Stopped'">启动</el-dropdown-item>
                <el-dropdown-item command="stop" v-if="row.status === 'Running'">停止</el-dropdown-item>
                <el-dropdown-item command="reboot">重启</el-dropdown-item>
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
import { Plus, ArrowDown } from '@element-plus/icons-vue'
import { getVMs, vmAction, deleteVM } from '@/api/compute'
import type { VirtualMachine } from '@/types'
import VMModal from '@/components/vm/VMModal.vue'
import VNCConsole from '@/components/vm/VNCConsole.vue'
import CreateVMModal from '@/components/vm/CreateVMModal.vue'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'

// 响应式数据
const vms = ref<VirtualMachine[]>([])
const loading = ref(false)
const detailsModalVisible = ref(false)
const vncModalVisible = ref(false)
const createModalVisible = ref(false)
const selectedVM = ref<VirtualMachine | null>(null)

// 分页数据
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 查询表单
const queryForm = reactive({
  account_id: null as number | null,
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

const getPlatformLabel = (platform: string): string => {
  return platformLabels[platform] || platform || '未知'
}

const getPlatformType = (platform: string): 'primary' | 'warning' | 'success' | 'info' => {
  return platformTypes[platform] || 'info'
}

const getStatusName = (status: string) => {
  const map: Record<string, string> = {
    Running: '运行中',
    Stopped: '已停止',
    Starting: '启动中',
    Stopping: '停止中',
    Pending: '创建中',
    Error: '错误',
    Deleted: '已删除'
  }
  return map[status] || status
}

const getStatusType = (status: string) => {
  if (status === 'Running') return 'success'
  if (status === 'Stopped') return 'info'
  if (status === 'Pending' || status === 'Starting' || status === 'Stopping') return 'warning'
  if (status === 'Error') return 'danger'
  return ''
}

const handleAccountChange = (accountId: number | null) => {
  queryForm.account_id = accountId
  loadVMs()
}

const loadVMs = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      size: pageSize.value
    }

    if (queryForm.account_id) params.account_id = queryForm.account_id
    if (queryForm.name) params.name = queryForm.name
    if (queryForm.status) params.status = queryForm.status

    const res = await getVMs(params)
    vms.value = Array.isArray(res) ? res : res.items || res
    total.value = res.total || vms.value.length
  } catch (e) {
    console.error(e)
    ElMessage.error('加载虚拟机列表失败')
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  createModalVisible.value = true
}

const handleDropdownCommand = async (command: string, row: VirtualMachine) => {
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

onMounted(() => {
  loadVMs()
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