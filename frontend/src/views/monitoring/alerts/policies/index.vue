<template>
  <div class="alert-policies-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">告警策略</span>
        </div>
      </template>

      <!-- Tabs -->
      <el-tabs v-model="activeTab" class="policy-tabs">
        <el-tab-pane label="全部" name="all">
          <PolicyTable :policies="allPolicies" @edit="handleEdit" @delete="handleDelete" />
        </el-tab-pane>
        <el-tab-pane label="自定义策略" name="custom">
          <PolicyTable :policies="customPolicies" @edit="handleEdit" @delete="handleDelete" />
        </el-tab-pane>
        <el-tab-pane label="默认策略" name="default">
          <PolicyTable :policies="defaultPolicies" @edit="handleEdit" @delete="handleDelete" />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PolicyTable from './components/PolicyTable.vue'

interface AlertPolicy {
  id: string
  name: string
  status: string
  enabled: boolean
  resource_type: string
  detail: string
  level: string
  owner: string
}

const activeTab = ref('all')
const policies = ref<AlertPolicy[]>([])

const allPolicies = computed(() => policies.value)
const customPolicies = computed(() => policies.value.filter(p => p.owner === '自定义'))
const defaultPolicies = computed(() => policies.value.filter(p => p.owner === '系统'))

const loadPolicies = async () => {
  policies.value = [
    { id: 'policy-1', name: 'CPU使用率告警', status: '正常', enabled: true, resource_type: '虚拟机', detail: 'CPU>80%持续5分钟', level: '警告', owner: '系统' },
    { id: 'policy-2', name: '内存使用率告警', status: '正常', enabled: true, resource_type: '虚拟机', detail: '内存>90%持续3分钟', level: '严重', owner: '系统' },
    { id: 'policy-3', name: '磁盘使用率告警', status: '正常', enabled: true, resource_type: '虚拟机', detail: '磁盘>85%持续10分钟', level: '信息', owner: '系统' },
    { id: 'policy-4', name: '数据库连接数告警', status: '正常', enabled: true, resource_type: '数据库', detail: '连接数>500持续2分钟', level: '警告', owner: '系统' },
    { id: 'policy-5', name: '网络流量告警', status: '正常', enabled: false, resource_type: '网络', detail: '流量>1GB/s持续1分钟', level: '严重', owner: '自定义' },
    { id: 'policy-6', name: '应用响应时间告警', status: '正常', enabled: true, resource_type: '应用', detail: '响应时间>2s持续5分钟', level: '警告', owner: '自定义' }
  ]
}

const handleEdit = (row: AlertPolicy) => {
  ElMessage.info(`编辑策略: ${row.name}`)
}

const handleDelete = async (row: AlertPolicy) => {
  try {
    await ElMessageBox.confirm(`确认删除策略 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    policies.value = policies.value.filter(p => p.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadPolicies()
})
</script>

<style scoped>
.alert-policies-page {
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

.policy-tabs {
  margin-top: 10px;
}
</style>