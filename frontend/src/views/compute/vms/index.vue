<template>
  <div class="vms-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">虚拟机管理</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            创建虚拟机
          </el-button>
        </div>
      </template>

      <el-form :inline="true" :model="queryForm" class="query-form">
        <el-form-item label="云账户">
          <el-input v-model="queryForm.account_id" placeholder="请输入云账户 ID" />
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="queryForm.name" placeholder="虚拟机名称" />
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

      <el-table
        :data="vms"
        v-loading="loading"
        style="width: 100%"
        @row-dblclick="showDetails"
      >
        <el-table-column prop="name" label="名称" min-width="150">
          <template #default="{ row }">
            <el-link type="primary" :underline="false" @click="showDetails(row)">
              {{ row.name }}
            </el-link>
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
            <div>
              <div v-if="row.private_ip" class="ip-item">
                <el-tag size="small" type="info">内网</el-tag>
                {{ row.private_ip }}
              </div>
              <div v-if="row.public_ip" class="ip-item">
                <el-tag size="small" type="success">公网</el-tag>
                {{ row.public_ip }}
              </div>
              <div v-if="!row.private_ip && !row.public_ip">-</div>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="os_name" label="系统" width="120">
          <template #default="{ row }">
            {{ row.os_name || '未知' }}
          </template>
        </el-table-column>

        <el-table-column label="密码" width="120">
          <template #default="{ row }">
            <el-popover
              placement="top"
              title="密码信息"
              :width="200"
              trigger="hover"
            >
              <div>密码信息将在安全通道中提供</div>
              <template #reference>
                <el-button link size="small">
                  <el-icon><Hide /></el-icon>
                  <span>点击获取</span>
                </el-button>
              </template>
            </el-popover>
          </template>
        </el-table-column>

        <el-table-column label="安全组" min-width="150">
          <template #default="{ row }">
            <el-tooltip
              v-if="row.security_group_names && row.security_group_names.length > 0"
              effect="dark"
              placement="top"
            >
              <template #content>
                <div v-for="sg in row.security_group_names" :key="sg" class="sg-tooltip-item">
                  {{ sg }}
                </div>
              </template>
              <el-tag
                v-for="sg in row.security_group_names?.slice(0, 2)"
                :key="sg"
                size="small"
                type="info"
                class="sg-tag"
              >
                {{ sg }}
              </el-tag>
              <el-tag
                v-if="row.security_group_names && row.security_group_names.length > 2"
                size="small"
                type="info"
              >
                +{{ row.security_group_names.length - 2 }}
              </el-tag>
            </el-tooltip>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column prop="billing_method" label="计费方式" width="120">
          <template #default="{ row }">
            {{ row.billing_method || '按量付费' }}
          </template>
        </el-table-column>

        <el-table-column prop="platform" label="平台" width="120">
          <template #default="{ row }">
            {{ row.platform || '未知' }}
          </template>
        </el-table-column>

        <el-table-column prop="project_id" label="项目" width="150">
          <template #default="{ row }">
            {{ row.project_id || '未分配' }}
          </template>
        </el-table-column>

        <el-table-column prop="region_id" label="区域" width="150" />

        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <div class="operation-buttons">
              <el-button
                size="small"
                type="primary"
                @click="openVNC(row)"
              >
                <el-icon><Monitor /></el-icon>
                远程控制
              </el-button>

              <VMActionDropdown
                :vm="row"
                :account-id="parseInt(queryForm.account_id)"
                @remote-control="openVNC(row)"
                @refresh="loadVMs"
                @vm-action="handleVmAction"
              />
            </div>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        class="pagination"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </el-card>

    <!-- VM 详情模态框 -->
    <VMModal
      v-model="detailsModalVisible"
      :vm-id="selectedVM?.id || ''"
      :account-id="parseInt(queryForm.account_id)"
      @close="detailsModalVisible = false"
    />

    <!-- VNC 控制台模态框 -->
    <VNCConsole
      v-model="vncModalVisible"
      :vm-id="selectedVM?.id || ''"
      :vm-name="selectedVM?.name || ''"
      :account-id="parseInt(queryForm.account_id)"
      @close="vncModalVisible = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, ElPagination } from 'element-plus'
import {
  Plus,
  Hide,
  Monitor,
  ArrowRight,
  CaretRight
} from '@element-plus/icons-vue'
import { getVMs, vmAction, deleteVM } from '@/api/compute'
import type { VirtualMachine } from '@/types'
import VMModal from '@/components/vm/VMModal.vue'
import VNCConsole from '@/components/vm/VNCConsole.vue'
import VMActionDropdown from '@/components/vm/VMActionDropdown.vue'

// 响应式数据
const vms = ref<VirtualMachine[]>([])
const loading = ref(false)
const detailsModalVisible = ref(false)
const vncModalVisible = ref(false)
const selectedVM = ref<VirtualMachine | null>(null)

// 分页数据
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 查询表单
const queryForm = reactive({
  account_id: '',
  name: '',
  status: ''
})

// 方法
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

const loadVMs = async () => {
  if (!queryForm.account_id) {
    ElMessage.warning('请先输入云账户 ID')
    return
  }

  loading.value = true
  try {
    const params: any = {
      account_id: parseInt(queryForm.account_id),
    }

    if (queryForm.name) params.name = queryForm.name
    if (queryForm.status) params.status = queryForm.status

    // 模拟分页参数
    params.page = currentPage.value
    params.size = pageSize.value

    const res = await getVMs(params)
    vms.value = Array.isArray(res) ? res : res.items || res
    total.value = vms.value.length // 实际使用中应该是后端返回的总数
  } catch (e) {
    console.error(e)
    ElMessage.error('加载虚拟机列表失败')
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  ElMessage.info('创建虚拟机功能开发中')
}

const handleAction = async (row: VirtualMachine, action: string) => {
  try {
    await ElMessageBox.confirm(`确定要${action === 'start' ? '启动' : '停止'}该虚拟机吗？`, '提示', { type: 'warning' })
    await vmAction(row.id, parseInt(queryForm.account_id), action as any)
    ElMessage.success(`${action === 'start' ? '启动' : '停止'}成功`)
    loadVMs()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
      ElMessage.error(`操作失败: ${e.message}`)
    }
  }
}

const handleDelete = async (row: VirtualMachine) => {
  try {
    await ElMessageBox.confirm('确定要删除该虚拟机吗？', '提示', { type: 'warning' })
    await deleteVM(row.id, parseInt(queryForm.account_id))
    ElMessage.success('删除成功')
    loadVMs()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
      ElMessage.error(`删除失败: ${e.message}`)
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
  queryForm.account_id = ''
  queryForm.name = ''
  queryForm.status = ''
  currentPage.value = 1
  loadVMs()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  loadVMs()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadVMs()
}

const handleVmAction = (action: string, data?: any) => {
  console.log('VM action triggered:', action, data)
  // 处理特定的 VM 操作
}

onMounted(() => {
  // 可以从项目上下文获取默认账户ID
  loadVMs()
})

// 双击行查看详情
const onRowDoubleClick = (row: VirtualMachine) => {
  showDetails(row)
}
</script>

<style scoped>
.vms-page {
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

.query-form {
  margin-bottom: 20px;
}

.ip-item {
  margin: 2px 0;
  font-size: 12px;
}

.sg-tag {
  margin-right: 4px;
  margin-bottom: 4px;
}

.sg-tooltip-item {
  padding: 2px 0;
}

.operation-buttons {
  display: flex;
  gap: 6px;
  align-items: center;
}

.operation-buttons > * {
  margin-right: 0 !important;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}

.el-table :deep(.el-table__cell) {
  padding: 8px 0;
}
</style>
