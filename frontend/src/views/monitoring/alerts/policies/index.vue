<template>
  <div class="policies-container">
    <div class="page-header">
      <h2>告警策略管理</h2>
      <el-button type="primary" @click="showCreateDialog">新建策略</el-button>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="loadPolicies">
        <el-form-item label="资源类型">
          <el-select v-model="filters.resource_type" placeholder="资源类型" clearable>
            <el-option label="虚拟机" value="vm" />
            <el-option label="数据库" value="database" />
            <el-option label="网络" value="network" />
            <el-option label="应用" value="application" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.enabled" placeholder="启用状态" clearable>
            <el-option label="已启用" value="true" />
            <el-option label="已禁用" value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadPolicies">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Tabs -->
    <el-tabs v-model="activeTab" class="policy-tabs" @tab-change="handleTabChange">
        <el-tab-pane label="全部" name="all">
          <PolicyTable :policies="allPolicies" :loading="loading" @edit="handleEdit" @delete="handleDelete" @toggle="handleToggle" />
        </el-tab-pane>
        <el-tab-pane label="自定义策略" name="custom">
          <PolicyTable :policies="customPolicies" :loading="loading" @edit="handleEdit" @delete="handleDelete" @toggle="handleToggle" />
        </el-tab-pane>
        <el-tab-pane label="默认策略" name="default">
          <PolicyTable :policies="defaultPolicies" :loading="loading" @edit="handleEdit" @delete="handleDelete" @toggle="handleToggle" />
        </el-tab-pane>
      </el-tabs>

    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      class="pagination"
    />

    <!-- 新增/编辑告警策略对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogMode === 'add' ? '新建告警策略' : '编辑告警策略'" width="500px">
      <el-form :model="policyForm" label-width="100px">
        <el-form-item label="策略名称" required>
          <el-input v-model="policyForm.name" placeholder="请输入策略名称" />
        </el-form-item>
        <el-form-item label="资源类型" required>
          <el-select v-model="policyForm.resource_type" placeholder="请选择资源类型" style="width: 100%;">
            <el-option label="虚拟机" value="vm" />
            <el-option label="数据库" value="database" />
            <el-option label="网络" value="network" />
            <el-option label="应用" value="application" />
          </el-select>
        </el-form-item>
        <el-form-item label="监控指标" required>
          <el-select v-model="policyForm.metric" placeholder="请选择监控指标" style="width: 100%;">
            <el-option label="CPU使用率" value="cpu_usage" />
            <el-option label="内存使用率" value="memory_usage" />
            <el-option label="磁盘使用率" value="disk_usage" />
            <el-option label="网络流量" value="network_traffic" />
            <el-option label="数据库连接数" value="db_connections" />
          </el-select>
        </el-form-item>
        <el-form-item label="阈值" required>
          <el-input-number v-model="policyForm.threshold" :min="0" :max="100" style="width: 150px;" />
          <span style="margin-left: 10px;">%</span>
        </el-form-item>
        <el-form-item label="持续时间" required>
          <el-input-number v-model="policyForm.duration" :min="1" :max="60" style="width: 150px;" />
          <span style="margin-left: 10px;">分钟</span>
        </el-form-item>
        <el-form-item label="告警级别" required>
          <el-select v-model="policyForm.level" placeholder="请选择告警级别" style="width: 100%;">
            <el-option label="信息" value="信息" />
            <el-option label="警告" value="警告" />
            <el-option label="严重" value="严重" />
          </el-select>
        </el-form-item>
        <el-form-item label="通知渠道">
          <el-select v-model="policyForm.notify_channel" placeholder="请选择通知渠道" style="width: 100%;">
            <el-option label="邮件" value="email" />
            <el-option label="企业微信" value="wechat" />
            <el-option label="钉钉" value="dingtalk" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="policyForm.description" type="textarea" placeholder="策略描述" />
        </el-form-item>
        <el-form-item label="启用状态">
          <el-switch v-model="policyForm.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PolicyTable from './components/PolicyTable.vue'
import { listAlertPolicies, createAlertPolicy, updateAlertPolicy, deleteAlertPolicy, toggleAlertPolicy, type AlertPolicy, type AlertPolicyRequest } from '@/api/monitor'

const activeTab = ref('all')
const policies = ref<AlertPolicy[]>([])
const loading = ref(false)

// 筛选条件
const filters = reactive({
  resource_type: '',
  enabled: ''
})

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const allPolicies = computed(() => policies.value)
const customPolicies = computed(() => policies.value.filter(p => p.owner === '自定义'))
const defaultPolicies = computed(() => policies.value.filter(p => p.owner === '系统'))

// 对话框
const dialogVisible = ref(false)
const dialogMode = ref<'add' | 'edit'>('add')
const submitLoading = ref(false)
const policyForm = ref<AlertPolicyRequest>({
  name: '',
  resource_type: 'vm',
  metric: 'cpu_usage',
  threshold: 80,
  duration: 5,
  level: '警告',
  enabled: true,
  owner: '自定义',
  description: '',
  notify_channel: ''
})

const loadPolicies = async () => {
  loading.value = true
  try {
    policies.value = await listAlertPolicies()
  } catch (e: any) {
    console.error(e)
    ElMessage.error(e.message || '获取告警策略失败')
    policies.value = []
  } finally {
    loading.value = false
  }
}

const handleTabChange = () => {
  // Tab切换时不需要重新加载，因为数据已经全部加载
}

const showCreateDialog = () => {
  dialogMode.value = 'add'
  policyForm.value = {
    name: '',
    resource_type: 'vm',
    metric: 'cpu_usage',
    threshold: 80,
    duration: 5,
    level: '警告',
    enabled: true,
    owner: '自定义',
    description: '',
    notify_channel: ''
  }
  dialogVisible.value = true
}

const handleEdit = (row: AlertPolicy) => {
  dialogMode.value = 'edit'
  policyForm.value = {
    name: row.name,
    resource_type: row.resource_type,
    metric: row.metric,
    threshold: row.threshold,
    duration: row.duration,
    level: row.level,
    enabled: row.enabled,
    owner: row.owner,
    description: row.description,
    notify_channel: row.notify_channel,
    domain_id: row.domain_id,
    project_id: row.project_id
  }
  // 保存编辑的策略ID
  policyForm.value.name = row.name // 临时存储ID用于提交
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!policyForm.value.name) {
    ElMessage.warning('请输入策略名称')
    return
  }

  submitLoading.value = true
  try {
    if (dialogMode.value === 'add') {
      const newPolicy = await createAlertPolicy(policyForm.value)
      policies.value.push(newPolicy)
      ElMessage.success('策略创建成功')
    } else {
      // 找到当前编辑的策略
      const currentPolicy = policies.value.find(p => p.name === policyForm.value.name)
      if (currentPolicy) {
        const updated = await updateAlertPolicy(currentPolicy.id, policyForm.value)
        const index = policies.value.findIndex(p => p.id === currentPolicy.id)
        if (index !== -1) {
          policies.value[index] = updated
        }
        ElMessage.success('策略更新成功')
      }
    }
    dialogVisible.value = false
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    submitLoading.value = false
  }
}

const handleDelete = async (row: AlertPolicy) => {
  try {
    await ElMessageBox.confirm(`确认删除策略 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await deleteAlertPolicy(row.id)
    policies.value = policies.value.filter(p => p.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const handleToggle = async (row: AlertPolicy) => {
  try {
    const updated = await toggleAlertPolicy(row.id)
    const index = policies.value.findIndex(p => p.id === row.id)
    if (index !== -1) {
      policies.value[index] = updated
    }
    ElMessage.success(updated.enabled ? '策略已启用' : '策略已禁用')
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

const resetFilters = () => {
  filters.resource_type = ''
  filters.enabled = ''
  pagination.page = 1
  loadPolicies()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadPolicies()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadPolicies()
}

onMounted(() => {
  loadPolicies()
})
</script>

<style scoped>
.policies-container {
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

.policy-tabs {
  margin-top: 0;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}
</style>