<template>
  <div class="snapshot-policies-container">
    <div class="page-header">
      <h2>自动快照策略</h2>
      <div class="toolbar">
        <el-button link type="primary" @click="handleView">
          <el-icon><View /></el-icon>
          查看
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建
        </el-button>
        <el-button :disabled="selectedPolicies.length === 0" type="danger" @click="handleBatchDelete">
          删除
        </el-button>
      </div>
    </div>

    <!-- 公有云/私有云 Tabs -->
    <el-tabs v-model="activeTab" class="policy-tabs" @tab-change="handleTabChange">
      <el-tab-pane label="公有云" name="public">
        <el-card class="filter-card">
          <el-form :inline="true" :model="filters" @submit.prevent="fetchData">
            <el-form-item label="名称">
              <el-input v-model="filters.name" placeholder="搜索策略名称" clearable style="width: 180px" />
            </el-form-item>
            <el-form-item label="状态">
              <el-select v-model="filters.status" placeholder="选择状态" clearable style="width: 120px">
                <el-option label="启用" value="active" />
                <el-option label="禁用" value="inactive" />
              </el-select>
            </el-form-item>
            <el-form-item label="平台">
              <el-select v-model="filters.platform" placeholder="选择平台" clearable style="width: 120px">
                <el-option label="阿里云" value="aliyun" />
                <el-option label="腾讯云" value="tencent" />
                <el-option label="AWS" value="aws" />
              </el-select>
            </el-form-item>
            <el-form-item label="资源类型">
              <el-select v-model="filters.resource_type" placeholder="选择类型" clearable style="width: 120px">
                <el-option label="磁盘" value="disk" />
                <el-option label="主机" value="instance" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="fetchData">查询</el-button>
              <el-button @click="resetFilters">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-table
          :data="policies"
          v-loading="loading"
          style="width: 100%"
          row-key="id"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column label="名称" min-width="200">
            <template #default="{ row }">
              <el-link type="primary" :underline="false" @click="handleDetails(row)">
                {{ row.name }}
              </el-link>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'active' ? 'success' : 'info'">
                {{ row.status === 'active' ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="标签" width="100">
            <template #default="{ row }">
              <el-tag v-for="tag in (row.tags || [])" :key="tag.key" size="small" class="tag-item">
                {{ tag.key }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="平台" width="100">
            <template #default="{ row }">
              <el-tag size="small" :type="getPlatformType(row.provider_type)">
                {{ getPlatformLabel(row.provider_type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="account_name" label="云账号" width="150" />
          <el-table-column prop="resource_type" label="资源类型" width="100">
            <template #default="{ row }">
              <el-tag size="small" :type="row.resource_type === 'disk' ? 'primary' : 'warning'">
                {{ row.resource_type === 'disk' ? '磁盘' : '主机' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="关联资源数" width="100">
            <template #default="{ row }">
              {{ row.associated_count || 0 }}
            </template>
          </el-table-column>
          <el-table-column label="策略" width="200">
            <template #default="{ row }">
              <div class="policy-info">
                <span class="policy-schedule">{{ row.schedule_text }}</span>
                <span class="policy-retention">保留{{ row.retention_days }}天</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="region_id" label="区域" width="120" />
          <el-table-column prop="project_name" label="项目" width="120" />
          <el-table-column label="操作" width="160" fixed="right">
            <template #default="{ row }">
              <el-dropdown trigger="click" @command="(cmd: string) => handleActionCommand(cmd, row)">
                <el-button size="small" link type="primary">
                  操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="details">查看详情</el-dropdown-item>
                    <el-dropdown-item command="edit">编辑</el-dropdown-item>
                    <el-dropdown-item command="toggle">
                      {{ row.status === 'active' ? '禁用' : '启用' }}
                    </el-dropdown-item>
                    <el-dropdown-item command="associate">关联资源</el-dropdown-item>
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
            <template #title>私有云自动快照策略管理功能开发中</template>
          </el-alert>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 创建策略弹窗 -->
    <el-dialog
      title="创建自动快照策略"
      v-model="createDialogVisible"
      width="600px"
    >
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="120px">
        <el-form-item label="云账号" prop="cloud_account_id">
          <el-select v-model="createForm.cloud_account_id" placeholder="选择云账号" style="width: 100%">
            <el-option v-for="acc in cloudAccounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入策略名称" />
        </el-form-item>
        <el-form-item label="资源类型" prop="resource_type">
          <el-radio-group v-model="createForm.resource_type">
            <el-radio value="disk">磁盘</el-radio>
            <el-radio value="instance">主机</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="执行时间" prop="schedule_type">
          <el-radio-group v-model="createForm.schedule_type">
            <el-radio value="daily">每天</el-radio>
            <el-radio value="weekly">每周</el-radio>
            <el-radio value="monthly">每月</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="执行时间点" prop="execute_time">
          <el-time-picker v-model="createForm.execute_time" format="HH:mm" value-format="HH:mm" style="width: 100%" />
        </el-form-item>
        <el-form-item label="执行日期" v-if="createForm.schedule_type === 'weekly'">
          <el-select v-model="createForm.week_day" placeholder="选择周几" style="width: 100%">
            <el-option label="周一" value="1" />
            <el-option label="周二" value="2" />
            <el-option label="周三" value="3" />
            <el-option label="周四" value="4" />
            <el-option label="周五" value="5" />
            <el-option label="周六" value="6" />
            <el-option label="周日" value="7" />
          </el-select>
        </el-form-item>
        <el-form-item label="执行日期" v-if="createForm.schedule_type === 'monthly'">
          <el-select v-model="createForm.month_day" placeholder="选择日期" style="width: 100%">
            <el-option v-for="d in monthDays" :key="d.value" :label="d.label" :value="d.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="保留天数" prop="retention_days">
          <el-input-number v-model="createForm.retention_days" :min="1" :max="365" style="width: 100%" />
        </el-form-item>
        <el-form-item label="区域">
          <el-select v-model="createForm.region_id" placeholder="选择区域" clearable style="width: 100%">
            <el-option label="华东1(杭州)" value="cn-hangzhou" />
            <el-option label="华东2(上海)" value="cn-shanghai" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCreate" :loading="creating">创建</el-button>
      </template>
    </el-dialog>

    <!-- 详情弹窗 -->
    <el-dialog
      title="快照策略详情"
      v-model="detailDialogVisible"
      width="700px"
    >
      <el-tabs v-model="detailTab" v-if="selectedPolicy">
        <el-tab-pane label="基础信息" name="basic">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">{{ selectedPolicy.policy_id || selectedPolicy.id }}</el-descriptions-item>
            <el-descriptions-item label="名称">{{ selectedPolicy.name }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="selectedPolicy.status === 'active' ? 'success' : 'info'">
                {{ selectedPolicy.status === 'active' ? '启用' : '禁用' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="资源类型">
              <el-tag size="small" :type="selectedPolicy.resource_type === 'disk' ? 'primary' : 'warning'">
                {{ selectedPolicy.resource_type === 'disk' ? '磁盘' : '主机' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="执行策略">{{ selectedPolicy.schedule_text }}</el-descriptions-item>
            <el-descriptions-item label="执行时间">{{ selectedPolicy.execute_time }}</el-descriptions-item>
            <el-descriptions-item label="保留天数">{{ selectedPolicy.retention_days }} 天</el-descriptions-item>
            <el-descriptions-item label="关联资源数">{{ selectedPolicy.associated_count || 0 }}</el-descriptions-item>
            <el-descriptions-item label="平台">
              <el-tag size="small" :type="getPlatformType(selectedPolicy.provider_type)">
                {{ getPlatformLabel(selectedPolicy.provider_type) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="云账号">{{ selectedPolicy.account_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="区域">{{ selectedPolicy.region_id || '-' }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ selectedPolicy.created_at }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="关联资源" name="resources">
          <el-table :data="associatedResources" style="width: 100%">
            <el-table-column prop="name" label="名称" width="200" />
            <el-table-column prop="type" label="类型" width="100" />
            <el-table-column prop="region" label="区域" width="120" />
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button size="small" link type="danger" @click="removeAssociation(row)">移除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-button type="primary" size="small" style="margin-top: 10px" @click="addAssociation">添加关联</el-button>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>

    <!-- 关联资源弹窗 -->
    <el-dialog
      title="关联资源"
      v-model="associateDialogVisible"
      width="500px"
    >
      <el-form :model="associateForm" label-width="100px">
        <el-form-item label="策略">{{ associateForm.policy_name }}</el-form-item>
        <el-form-item label="选择资源">
          <el-table :data="availableResources" style="width: 100%" max-height="300">
            <el-table-column type="selection" width="55" />
            <el-table-column prop="name" label="名称" />
            <el-table-column prop="region" label="区域" width="120" />
          </el-table>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="associateDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAssociate" :loading="associating">确认关联</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Plus, ArrowDown } from '@element-plus/icons-vue'
import {
  getSnapshotPolicies,
  createSnapshotPolicy,
  deleteSnapshotPolicy,
  batchDeleteSnapshotPolicies,
  toggleSnapshotPolicy,
  associateResourcesToPolicy,
  disassociateResourceFromPolicy,
  type SnapshotPolicy,
  type SnapshotPolicyListParams,
  type CreateSnapshotPolicyParams
} from '@/api/storage'

// Data
const loading = ref(false)
const creating = ref(false)
const associating = ref(false)
const policies = ref<SnapshotPolicy[]>([])
const selectedPolicies = ref<SnapshotPolicy[]>([])
const cloudAccounts = ref<any[]>([])

const monthDays = Array.from({ length: 28 }, (_, i) => ({
  label: `每月${i + 1}日`,
  value: (i + 1).toString()
}))

const activeTab = ref('public')
const detailTab = ref('basic')
const createDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const associateDialogVisible = ref(false)

const selectedPolicy = ref<SnapshotPolicy | null>(null)
const associatedResources = ref<{ name: string; type: string; region: string }[]>([])
const availableResources = ref<{ name: string; region: string }[]>([])

const filters = reactive({
  name: '',
  status: '',
  platform: '',
  resource_type: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const createForm = reactive({
  cloud_account_id: '',
  name: '',
  resource_type: 'disk',
  schedule_type: 'daily',
  execute_time: '02:00',
  week_day: '1',
  month_day: '1',
  retention_days: 7,
  region_id: ''
})

const associateForm = reactive({
  policy_id: '',
  policy_name: '',
  resource_type: ''
})

const createRules = {
  cloud_account_id: [{ required: true, message: '请选择云账号', trigger: 'change' }],
  name: [{ required: true, message: '请输入策略名称', trigger: 'blur' }],
  resource_type: [{ required: true, message: '请选择资源类型', trigger: 'change' }],
  execute_time: [{ required: true, message: '请选择执行时间', trigger: 'change' }],
  retention_days: [{ required: true, message: '请输入保留天数', trigger: 'blur' }]
}

// Methods
const getPlatformLabel = (platform: string) => {
  const labels: Record<string, string> = {
    aliyun: '阿里云',
    tencent: '腾讯云',
    aws: 'AWS'
  }
  return labels[platform] || platform || '未知'
}

const getPlatformType = (platform: string) => {
  const types: Record<string, 'primary' | 'warning' | 'success' | 'info'> = {
    aliyun: 'primary',
    tencent: 'warning',
    aws: 'success'
  }
  return types[platform] || 'info'
}

const handleSelectionChange = (selection: SnapshotPolicy[]) => {
  selectedPolicies.value = selection
}

const handleTabChange = () => {
  pagination.page = 1
  fetchData()
}

const resetFilters = () => {
  filters.name = ''
  filters.status = ''
  filters.platform = ''
  filters.resource_type = ''
  pagination.page = 1
  fetchData()
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: SnapshotPolicyListParams = {
      page: pagination.page,
      page_size: pagination.pageSize,
      name: filters.name || undefined,
      status: filters.status || undefined,
      platform: filters.platform || undefined,
      resource_type: filters.resource_type || undefined
    }
    const res = await getSnapshotPolicies(params)
    // Transform data to add schedule_text
    policies.value = res.items.map(p => ({
      ...p,
      schedule_text: getScheduleText(p)
    }))
    pagination.total = res.total
  } catch (error) {
    console.error('Failed to fetch policies:', error)
    ElMessage.error('获取快照策略列表失败')
  } finally {
    loading.value = false
  }
}

const getScheduleText = (policy: SnapshotPolicy) => {
  switch (policy.schedule_type) {
    case 'daily': return `每天${policy.execute_time}执行`
    case 'weekly': return `每周${getWeekDayText(policy.week_day)}${policy.execute_time}执行`
    case 'monthly': return `每月${policy.month_day}日${policy.execute_time}执行`
    default: return policy.execute_time
  }
}

const getWeekDayText = (weekDay: string) => {
  const days = ['一', '二', '三', '四', '五', '六', '日']
  return days[parseInt(weekDay) - 1] || weekDay
}

const fetchCloudAccounts = async () => {
  cloudAccounts.value = [{ id: 1, name: '阿里云账号1' }, { id: 2, name: '腾讯云账号1' }]
}

const handleView = () => {
  ElMessage.info('视图切换功能开发中')
}

const handleCreate = () => {
  Object.assign(createForm, {
    cloud_account_id: '',
    name: '',
    resource_type: 'disk',
    schedule_type: 'daily',
    execute_time: '02:00',
    week_day: '1',
    month_day: '1',
    retention_days: 7,
    region_id: ''
  })
  createDialogVisible.value = true
}

const confirmCreate = async () => {
  creating.value = true
  try {
    const params: CreateSnapshotPolicyParams = {
      cloud_account_id: Number(createForm.cloud_account_id),
      name: createForm.name,
      resource_type: createForm.resource_type,
      schedule_type: createForm.schedule_type,
      execute_time: createForm.execute_time,
      week_day: createForm.schedule_type === 'weekly' ? createForm.week_day : undefined,
      month_day: createForm.schedule_type === 'monthly' ? createForm.month_day : undefined,
      retention_days: createForm.retention_days,
      region_id: createForm.region_id || undefined
    }
    await createSnapshotPolicy(params)
    ElMessage.success('策略创建成功')
    createDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('创建失败')
  } finally {
    creating.value = false
  }
}

const handleBatchDelete = async () => {
  if (selectedPolicies.value.length === 0) return
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedPolicies.value.length} 个策略吗？`,
      '批量删除确认',
      { type: 'warning' }
    )
    const ids = selectedPolicies.value.map(p => p.id)
    await batchDeleteSnapshotPolicies(ids)
    ElMessage.success('删除成功')
    selectedPolicies.value = []
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

const handleActionCommand = (command: string, row: SnapshotPolicy) => {
  switch (command) {
    case 'details': handleDetails(row); break
    case 'edit': handleEdit(row); break
    case 'toggle': handleToggle(row); break
    case 'associate': handleAssociate(row); break
    case 'delete': handleDelete(row); break
  }
}

const handleDetails = (row: SnapshotPolicy) => {
  selectedPolicy.value = row
  detailTab.value = 'basic'
  associatedResources.value = [
    { name: 'disk-data-01', type: '数据盘', region: 'cn-hangzhou' },
    { name: 'disk-data-02', type: '数据盘', region: 'cn-hangzhou' },
    { name: 'disk-data-03', type: '数据盘', region: 'cn-hangzhou' }
  ]
  detailDialogVisible.value = true
}

const handleEdit = (row: SnapshotPolicy) => {
  ElMessage.info('编辑功能开发中')
}

const handleToggle = async (row: SnapshotPolicy) => {
  const newStatus = row.status === 'active' ? false : true
  const actionText = newStatus ? '启用' : '禁用'
  try {
    await ElMessageBox.confirm(`确定要${actionText}策略 "${row.name}" 吗？`, '提示', { type: 'warning' })
    await toggleSnapshotPolicy(row.id, newStatus)
    row.status = newStatus ? 'active' : 'inactive'
    ElMessage.success(`${actionText}成功`)
  } catch (error) {
    if (error !== 'cancel') ElMessage.error(`${actionText}失败`)
  }
}

const handleAssociate = (row: SnapshotPolicy) => {
  Object.assign(associateForm, {
    policy_id: row.id,
    policy_name: row.name,
    resource_type: row.resource_type
  })
  availableResources.value = [
    { name: 'disk-backup-01', region: 'cn-hangzhou' },
    { name: 'disk-backup-02', region: 'cn-hangzhou' },
    { name: 'disk-backup-03', region: 'cn-shanghai' }
  ]
  associateDialogVisible.value = true
}

const confirmAssociate = async () => {
  associating.value = true
  try {
    const resourceIds = ['disk-backup-01', 'disk-backup-02'] // Mock selected resources
    await associateResourcesToPolicy(associateForm.policy_id, resourceIds)
    ElMessage.success('关联成功')
    associateDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('关联失败')
  } finally {
    associating.value = false
  }
}

const removeAssociation = async (row: any) => {
  try {
    await disassociateResourceFromPolicy(selectedPolicy.value!.id, row.name)
    const index = associatedResources.value.findIndex(r => r.name === row.name)
    if (index > -1) {
      associatedResources.value.splice(index, 1)
      ElMessage.success('已移除关联')
    }
  } catch (error) {
    ElMessage.error('移除关联失败')
  }
}

const addAssociation = () => {
  ElMessage.info('添加关联功能开发中')
}

const handleDelete = async (row: SnapshotPolicy) => {
  try {
    await ElMessageBox.confirm(`确定要删除策略 "${row.name}" 吗？`, '删除警告', { type: 'warning' })
    await deleteSnapshotPolicy(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
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
})
</script>

<style scoped>
.snapshot-policies-container { padding: 20px; }

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
.policy-tabs { margin-bottom: 20px; }

.policy-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.policy-schedule { font-size: 13px; }
.policy-retention { font-size: 12px; color: var(--el-text-color-secondary); }

.pagination { margin-top: 20px; justify-content: flex-end; }
.tag-item { margin-right: 4px; margin-bottom: 2px; }
</style>