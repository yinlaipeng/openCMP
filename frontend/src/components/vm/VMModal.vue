<template>
  <el-dialog
    v-model="visible"
    :title="`${vm?.name || '虚拟机'} - 详情`"
    width="80%"
    :destroy-on-close="true"
    :close-on-click-modal="false"
  >
    <el-tabs v-model="activeTab" type="card" v-if="vm">
      <!-- 基本信息标签页 -->
      <el-tab-pane label="详情" name="details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID">{{ vm.id }}</el-descriptions-item>
          <el-descriptions-item label="名称">{{ vm.name }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(vm.status)">{{ getStatusName(vm.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="实例类型">{{ vm.instance_type }}</el-descriptions-item>
          <el-descriptions-item label="操作系统">{{ vm.os_name || '未知' }}</el-descriptions-item>
          <el-descriptions-item label="计费方式">{{ vm.billing_method || '未知' }}</el-descriptions-item>
          <el-descriptions-item label="平台">{{ vm.platform || '未知' }}</el-descriptions-item>
          <el-descriptions-item label="项目">{{ vm.project_id || '未分配' }}</el-descriptions-item>
          <el-descriptions-item label="内网IP">{{ vm.private_ip || '无' }}</el-descriptions-item>
          <el-descriptions-item label="外网IP">{{ vm.public_ip || '无' }}</el-descriptions-item>
          <el-descriptions-item label="区域">{{ vm.region_id }}</el-descriptions-item>
          <el-descriptions-item label="可用区">{{ vm.zone_id }}</el-descriptions-item>
          <el-descriptions-item label="VPC">{{ vm.vpc_id }}</el-descriptions-item>
          <el-descriptions-item label="子网">{{ vm.subnet_id }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDate(vm.created_at) }}</el-descriptions-item>
        </el-descriptions>
      </el-tab-pane>

      <!-- 安全组标签页 -->
      <el-tab-pane label="安全组" name="security-groups">
        <el-table :data="securityGroups" v-loading="sgLoading">
          <el-table-column prop="id" label="ID" width="200" />
          <el-table-column prop="name" label="名称" />
          <el-table-column prop="description" label="描述" />
          <el-table-column prop="vpc_id" label="VPC" width="200" />
          <el-table-column prop="created_at" label="创建时间" :formatter="formatDateColumn" />
        </el-table>
      </el-tab-pane>

      <!-- 网络标签页 -->
      <el-tab-pane label="网络" name="network">
        <div v-if="networkInfo" class="network-info">
          <h4>VPC 信息</h4>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="VPC ID">{{ networkInfo.vpc?.id }}</el-descriptions-item>
            <el-descriptions-item label="VPC 名称">{{ networkInfo.vpc?.name }}</el-descriptions-item>
            <el-descriptions-item label="CIDR">{{ networkInfo.vpc?.cidr }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ networkInfo.vpc?.status }}</el-descriptions-item>
          </el-descriptions>

          <h4>子网信息</h4>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="子网 ID">{{ networkInfo.subnet?.id }}</el-descriptions-item>
            <el-descriptions-item label="子网名称">{{ networkInfo.subnet?.name }}</el-descriptions-item>
            <el-descriptions-item label="CIDR">{{ networkInfo.subnet?.cidr }}</el-descriptions-item>
            <el-descriptions-item label="可用区">{{ networkInfo.subnet?.zone_id }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ networkInfo.subnet?.status }}</el-descriptions-item>
          </el-descriptions>

          <h4>IP 信息</h4>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="内网IP">{{ networkInfo.private_ip }}</el-descriptions-item>
            <el-descriptions-item label="外网IP">{{ networkInfo.public_ip }}</el-descriptions-item>
          </el-descriptions>
        </div>
      </el-tab-pane>

      <!-- 磁盘标签页 -->
      <el-tab-pane label="磁盘" name="disks">
        <el-table :data="disks" v-loading="diskLoading">
          <el-table-column prop="id" label="ID" width="200" />
          <el-table-column prop="name" label="名称" />
          <el-table-column prop="size" label="大小(G)" />
          <el-table-column prop="type" label="类型" />
          <el-table-column prop="status" label="状态">
            <template #default="{ row }">
              <el-tag :type="getDiskStatusType(row.status)">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="zone_id" label="可用区" />
          <el-table-column prop="created_at" label="创建时间" :formatter="formatDateColumn" />
        </el-table>
      </el-tab-pane>

      <!-- 快照标签页 -->
      <el-tab-pane label="快照" name="snapshots">
        <el-table :data="snapshots" v-loading="snapshotLoading">
          <el-table-column prop="id" label="ID" width="200" />
          <el-table-column prop="name" label="名称" />
          <el-table-column prop="disk_id" label="磁盘ID" width="200" />
          <el-table-column prop="size" label="大小(G)" />
          <el-table-column prop="status" label="状态" />
          <el-table-column prop="created_at" label="创建时间" :formatter="formatDateColumn" />
        </el-table>
      </el-tab-pane>

      <!-- 监控标签页 -->
      <el-tab-pane label="监控" name="monitoring">
        <div class="monitoring-container">
          <p>监控数据将在后续版本中实现</p>
        </div>
      </el-tab-pane>

      <!-- 告警标签页 -->
      <el-tab-pane label="告警" name="alarms">
        <div class="alarms-container">
          <p>告警信息将在后续版本中实现</p>
        </div>
      </el-tab-pane>

      <!-- 定时任务标签页 -->
      <el-tab-pane label="定时任务" name="scheduled-tasks">
        <div class="scheduled-tasks-container">
          <p>关联的定时任务将在后续版本中实现</p>
        </div>
      </el-tab-pane>

      <!-- 操作日志标签页 -->
      <el-tab-pane label="操作日志" name="operation-logs">
        <el-table :data="operationLogs" v-loading="logLoading">
          <el-table-column prop="id" label="ID" width="150" />
          <el-table-column prop="operation" label="操作" width="120" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getLogStatusType(row.status)">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="operator" label="操作人" width="120" />
          <el-table-column prop="timestamp" label="时间" width="180" />
          <el-table-column prop="details" label="详情" />
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="close">关闭</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import {
  getVMDetails,
  getVMSecurityGroups,
  getVMNetworkInfo,
  getVMDiskInfo,
  getVMSnapshots,
  getVMOperationLogs
} from '@/api/compute'
import type { VirtualMachine } from '@/types'

interface Props {
  modelValue: boolean
  vmId: string
  accountId: number
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'close'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 响应式数据
const activeTab = ref('details')
const vm = ref<VirtualMachine | null>(null)
const securityGroups = ref<any[]>([])
const networkInfo = ref<any>(null)
const disks = ref<any[]>([])
const snapshots = ref<any[]>([])
const operationLogs = ref<any[]>([])

// Loading状态
const sgLoading = ref(false)
const diskLoading = ref(false)
const snapshotLoading = ref(false)
const logLoading = ref(false)

// 计算属性
const visible = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})

// 方法
const close = () => {
  visible.value = false
  emit('close')
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

const getDiskStatusType = (status: string) => {
  if (status === 'available') return 'success'
  if (status === 'in-use') return 'primary'
  if (status === 'error') return 'danger'
  return 'info'
}

const getLogStatusType = (status: string) => {
  if (status === 'success') return 'success'
  if (status === 'failed') return 'danger'
  return 'info'
}

const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

const formatDateColumn = (_: any, __: any, cellValue: string) => {
  return formatDate(cellValue)
}

// 加载详细数据
const loadDetails = async () => {
  if (!props.vmId || !props.accountId) return

  try {
    // 加载基本信息
    const vmDetail = await getVMDetails(props.vmId, props.accountId)
    vm.value = vmDetail

    // 并行加载其他数据
    await Promise.all([
      loadSecurityGroups(),
      loadNetworkInfo(),
      loadDisks(),
      loadSnapshots(),
      loadOperationLogs()
    ])
  } catch (error) {
    console.error('Failed to load VM details:', error)
    ElMessage.error('加载虚拟机详情失败')
  }
}

const loadSecurityGroups = async () => {
  sgLoading.value = true
  try {
    const response = await getVMSecurityGroups(props.vmId, props.accountId)
    securityGroups.value = response.items || []
  } catch (error) {
    console.error('Failed to load security groups:', error)
    ElMessage.error('加载安全组信息失败')
  } finally {
    sgLoading.value = false
  }
}

const loadNetworkInfo = async () => {
  try {
    const response = await getVMNetworkInfo(props.vmId, props.accountId)
    networkInfo.value = response
  } catch (error) {
    console.error('Failed to load network info:', error)
    ElMessage.error('加载网络信息失败')
  }
}

const loadDisks = async () => {
  diskLoading.value = true
  try {
    const response = await getVMDiskInfo(props.vmId, props.accountId)
    disks.value = response.items || []
  } catch (error) {
    console.error('Failed to load disks:', error)
    ElMessage.error('加载磁盘信息失败')
  } finally {
    diskLoading.value = false
  }
}

const loadSnapshots = async () => {
  snapshotLoading.value = true
  try {
    const response = await getVMSnapshots(props.vmId, props.accountId)
    snapshots.value = response.items || []
  } catch (error) {
    console.error('Failed to load snapshots:', error)
    ElMessage.error('加载快照信息失败')
  } finally {
    snapshotLoading.value = false
  }
}

const loadOperationLogs = async () => {
  logLoading.value = true
  try {
    const response = await getVMOperationLogs(props.vmId, props.accountId)
    operationLogs.value = response.items || []
  } catch (error) {
    console.error('Failed to load operation logs:', error)
    ElMessage.error('加载操作日志失败')
  } finally {
    logLoading.value = false
  }
}

// 监听可见性变化
watch(visible, async (newValue) => {
  if (newValue && props.vmId && props.accountId) {
    await loadDetails()
  } else if (!newValue) {
    // 清空数据
    vm.value = null
    securityGroups.value = []
    networkInfo.value = null
    disks.value = []
    snapshots.value = []
    operationLogs.value = []
    activeTab.value = 'details'
  }
})
</script>

<style scoped>
.network-info h4 {
  margin-top: 20px;
  margin-bottom: 10px;
}

.network-info h4:first-child {
  margin-top: 0;
}

.monitoring-container,
.alarms-container,
.scheduled-tasks-container {
  padding: 20px;
  text-align: center;
  color: #999;
}

.dialog-footer {
  text-align: right;
}
</style>