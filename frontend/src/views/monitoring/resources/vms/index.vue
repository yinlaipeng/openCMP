<template>
  <div class="monitoring-vms-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">虚拟机监控列表</span>
        </div>
      </template>

      <el-table :data="vmList" v-loading="loading">
        <el-table-column label="名称" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP" width="140" />
        <el-table-column prop="monitor_status" label="监控状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getMonitorStatusType(row.monitor_status)">{{ row.monitor_status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column prop="owner" label="归属" width="120" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="handleManagePolicies(row)">管理告警策略</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- VM Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="虚拟机详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="名称">{{ selectedVM?.name }}</el-descriptions-item>
        <el-descriptions-item label="IP">{{ selectedVM?.ip }}</el-descriptions-item>
        <el-descriptions-item label="监控状态">{{ selectedVM?.monitor_status }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedVM?.status }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedVM?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedVM?.account }}</el-descriptions-item>
        <el-descriptions-item label="归属">{{ selectedVM?.owner }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedVM?.region }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Manage Alert Policies Modal -->
    <el-dialog v-model="policiesDialogVisible" title="管理告警策略" width="800px">
      <el-table :data="vmPolicies" v-loading="policiesLoading">
        <el-table-column prop="name" label="名称" width="180" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getPolicyStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="启用状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '已启用' : '已禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_type" label="资源类型" width="100" />
        <el-table-column prop="detail" label="策略详情" width="150">
          <template #default="{ row }">
            <span>{{ row.detail }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="level" label="告警级别" width="100">
          <template #default="{ row }">
            <el-tag :type="getAlertLevelType(row.level)" size="small">{{ row.level }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="owner" label="策略归属" width="120" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="handleEditPolicy(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDeletePolicy(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="policiesDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="handleAddPolicy">新增策略</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

interface VMResource {
  id: string
  name: string
  ip: string
  monitor_status: string
  status: string
  platform: string
  account: string
  owner: string
  region: string
}

interface VMPolicy {
  id: string
  name: string
  status: string
  enabled: boolean
  resource_type: string
  detail: string
  level: string
  owner: string
}

const vmList = ref<VMResource[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedVM = ref<VMResource | null>(null)

const policiesDialogVisible = ref(false)
const vmPolicies = ref<VMPolicy[]>([])
const policiesLoading = ref(false)

const getMonitorStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case '正常':
    case 'healthy':
      return 'success'
    case '告警':
    case 'warning':
      return 'warning'
    case '异常':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'running':
    case 'active':
      return 'success'
    case 'stopped':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const getPolicyStatusType = (status: string) => {
  switch (status) {
    case '正常':
      return 'success'
    case '异常':
      return 'danger'
    default:
      return 'info'
  }
}

const getAlertLevelType = (level: string) => {
  switch (level) {
    case '严重':
      return 'danger'
    case '警告':
      return 'warning'
    case '信息':
      return 'info'
    default:
      return ''
  }
}

const loadVMList = async () => {
  loading.value = true
  try {
    vmList.value = [
      { id: 'vm-1', name: 'prod-web-01', ip: '192.168.1.10', monitor_status: '正常', status: 'Running', platform: '阿里云', account: 'Aliyun Account 1', owner: 'Project A', region: 'cn-hangzhou' },
      { id: 'vm-2', name: 'prod-web-02', ip: '192.168.1.11', monitor_status: '告警', status: 'Running', platform: '阿里云', account: 'Aliyun Account 1', owner: 'Project A', region: 'cn-hangzhou' },
      { id: 'vm-3', name: 'dev-api-01', ip: '192.168.2.10', monitor_status: '正常', status: 'Running', platform: '阿里云', account: 'Aliyun Account 1', owner: 'Project B', region: 'cn-shanghai' },
      { id: 'vm-4', name: 'dev-db-01', ip: '192.168.2.20', monitor_status: '正常', status: 'Stopped', platform: '阿里云', account: 'Aliyun Account 1', owner: 'Project B', region: 'cn-shanghai' },
      { id: 'vm-5', name: 'test-app-01', ip: '192.168.3.10', monitor_status: '异常', status: 'Running', platform: '阿里云', account: 'Aliyun Account 1', owner: 'Project A', region: 'cn-beijing' }
    ]
  } catch (e) {
    console.error(e)
    vmList.value = []
  } finally {
    loading.value = false
  }
}

const loadVMPolicies = async (vmId: string) => {
  policiesLoading.value = true
  try {
    vmPolicies.value = [
      { id: 'policy-1', name: 'CPU告警策略', status: '正常', enabled: true, resource_type: '虚拟机', detail: 'CPU>80%持续5分钟', level: '警告', owner: '系统' },
      { id: 'policy-2', name: '内存告警策略', status: '正常', enabled: true, resource_type: '虚拟机', detail: '内存>90%持续3分钟', level: '严重', owner: '系统' },
      { id: 'policy-3', name: '磁盘告警策略', status: '正常', enabled: false, resource_type: '虚拟机', detail: '磁盘>85%持续10分钟', level: '信息', owner: '自定义' }
    ]
  } catch (e) {
    console.error(e)
    vmPolicies.value = []
  } finally {
    policiesLoading.value = false
  }
}

const viewDetail = (row: VMResource) => {
  selectedVM.value = row
  detailDialogVisible.value = true
}

const handleManagePolicies = (row: VMResource) => {
  selectedVM.value = row
  loadVMPolicies(row.id)
  policiesDialogVisible.value = true
}

const handleEditPolicy = (row: VMPolicy) => {
  ElMessage.info(`编辑策略: ${row.name}`)
}

const handleDeletePolicy = async (row: VMPolicy) => {
  try {
    await ElMessageBox.confirm(`确认删除策略 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    vmPolicies.value = vmPolicies.value.filter(p => p.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

const handleAddPolicy = () => {
  ElMessage.info('新增策略功能开发中')
}

onMounted(() => {
  loadVMList()
})
</script>

<style scoped>
.monitoring-vms-page {
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
</style>