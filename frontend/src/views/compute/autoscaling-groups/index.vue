<template>
  <div class="autoscaling-group-container">
    <div class="page-header">
      <h2>弹性伸缩组</h2>
      <div class="toolbar">
        <el-button link type="primary" @click="handleView">
          <el-icon><View /></el-icon>
          查看
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建
        </el-button>
        <el-dropdown trigger="click" @command="handleBatchCommand" :disabled="selectedGroups.length === 0">
          <el-button :disabled="selectedGroups.length === 0">
            批量操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="enable">批量启用</el-dropdown-item>
              <el-dropdown-item command="disable">批量禁用</el-dropdown-item>
              <el-dropdown-item command="delete" divided>批量删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="fetchData">
        <el-form-item label="项目">
          <el-select v-model="filters.project_id" placeholder="请选择项目" clearable style="width: 160px">
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="filters.name" placeholder="搜索名称" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="filters.platform" placeholder="选择平台" clearable style="width: 120px">
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="选择状态" clearable style="width: 120px">
            <el-option label="激活" value="Active" />
            <el-option label="未激活" value="Inactive" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table
      :data="autoscalingGroups"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column label="名称" min-width="150">
        <template #default="{ row }">
          <el-link type="primary" :underline="false" @click="handleEdit(row)">
            {{ row.name }}
          </el-link>
        </template>
      </el-table-column>

      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="平台" width="100">
        <template #default="{ row }">
          <el-tag size="small" type="info">{{ getPlatformLabel(row.platform) }}</el-tag>
        </template>
      </el-table-column>

      <el-table-column label="主机模版" width="150">
        <template #default="{ row }">
          {{ getTemplateName(row.host_template_id) }}
        </template>
      </el-table-column>

      <el-table-column label="实例数" width="120">
        <template #default="{ row }">
          <span class="capacity-cell">{{ row.current_capacity }}/{{ row.desired_capacity }}</span>
        </template>
      </el-table-column>

      <el-table-column label="伸缩范围" width="120">
        <template #default="{ row }">
          <span>{{ row.min_size }} - {{ row.max_size }}</span>
        </template>
      </el-table-column>

      <el-table-column prop="region_id" label="区域" width="120" />

      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-dropdown trigger="click" @command="(cmd: string) => handleActionCommand(cmd, row)">
            <el-button size="small" link type="primary">
              操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="edit">编辑</el-dropdown-item>
                <el-dropdown-item command="scaleUp">扩容</el-dropdown-item>
                <el-dropdown-item command="scaleDown">缩容</el-dropdown-item>
                <el-dropdown-item :command="row.status === 'Active' ? 'disable' : 'enable'">
                  {{ row.status === 'Active' ? '禁用' : '启用' }}
                </el-dropdown-item>
                <el-dropdown-item command="details">查看详情</el-dropdown-item>
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

    <!-- Create/Edit Dialog -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="700px"
      :before-close="handleDialogClose"
    >
      <el-form
        :model="form"
        :rules="rules"
        ref="formRef"
        label-width="140px"
      >
        <el-form-item label="项目" prop="project_id">
          <el-select v-model="form.project_id" placeholder="请选择项目" style="width: 100%">
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>

        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="以小写字母开头,包含小写字母、数字或连字符" />
        </el-form-item>

        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>

        <el-form-item label="平台" prop="platform">
          <el-radio-group v-model="form.platform">
            <el-radio value="aliyun">阿里云</el-radio>
            <el-radio value="tencent">腾讯云</el-radio>
            <el-radio value="aws">AWS</el-radio>
            <el-radio value="azure">Azure</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="主机模版" prop="host_template_id">
          <el-select v-model="form.host_template_id" placeholder="请选择主机模版" style="width: 100%">
            <el-option v-for="t in hostTemplates" :key="t.id" :label="t.name" :value="t.id" />
          </el-select>
        </el-form-item>

        <el-form-item label="网络" prop="network_id">
          <el-select v-model="form.network_id" placeholder="请选择网络" style="width: 100%">
            <el-option v-for="n in networks" :key="n.id" :label="n.name" :value="n.id" />
          </el-select>
        </el-form-item>

        <el-form-item label="最大实例数" prop="max_size">
          <el-input-number v-model="form.max_size" :min="1" :max="1000" style="width: 100%" />
        </el-form-item>

        <el-form-item label="期望实例数" prop="desired_capacity">
          <el-input-number v-model="form.desired_capacity" :min="form.min_size" :max="form.max_size" style="width: 100%" />
        </el-form-item>

        <el-form-item label="最小实例数" prop="min_size">
          <el-input-number v-model="form.min_size" :min="0" :max="form.max_size" style="width: 100%" />
        </el-form-item>

        <el-form-item label="实例移出策略" prop="removal_policy">
          <el-select v-model="form.removal_policy" placeholder="请选择移出策略" style="width: 100%">
            <el-option label="最旧实例" value="oldest" />
            <el-option label="最新实例" value="newest" />
            <el-option label="随机" value="random" />
          </el-select>
        </el-form-item>

        <el-form-item label="负载均衡" prop="load_balancing">
          <el-radio-group v-model="form.load_balancing">
            <el-radio value="enabled">启用</el-radio>
            <el-radio value="disabled">不启用</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="健康检查方式" prop="health_check_type">
          <el-select v-model="form.health_check_type" placeholder="请选择健康检查方式" style="width: 100%">
            <el-option label="EC2" value="EC2" />
            <el-option label="ELB" value="ELB" />
          </el-select>
        </el-form-item>

        <el-form-item label="检查周期(秒)" prop="health_check_period">
          <el-input-number v-model="form.health_check_period" :min="10" :max="300" style="width: 100%" />
        </el-form-item>

        <el-form-item label="健康检查宽限期" prop="health_check_grace_period">
          <el-input-number v-model="form.health_check_grace_period" :min="0" :max="600" style="width: 100%" />
        </el-form-item>

        <el-form-item label="标签" prop="tags">
          <el-select v-model="form.tags" placeholder="选择或输入标签" multiple filterable allow-create style="width: 100%">
            <el-option v-for="tag in availableTags" :key="tag" :label="tag" :value="tag" />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="handleDialogClose">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- Scale Dialog -->
    <el-dialog
      title="调整实例数"
      v-model="scaleDialogVisible"
      width="400px"
    >
      <el-form :model="scaleForm" label-width="100px">
        <el-form-item label="当前实例数">
          <span>{{ selectedGroup?.current_capacity }}</span>
        </el-form-item>
        <el-form-item label="新实例数">
          <el-input-number
            v-model="scaleForm.capacity"
            :min="selectedGroup?.min_size"
            :max="selectedGroup?.max_size"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="scaleDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmScale">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- Detail Dialog -->
    <el-dialog
      title="伸缩组详情"
      v-model="detailDialogVisible"
      width="800px"
    >
      <el-descriptions :column="2" border v-if="selectedGroup">
        <el-descriptions-item label="名称">{{ selectedGroup.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(selectedGroup.status)">{{ selectedGroup.status }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="平台">{{ getPlatformLabel(selectedGroup.platform) }}</el-descriptions-item>
        <el-descriptions-item label="主机模版">{{ getTemplateName(selectedGroup.host_template_id) }}</el-descriptions-item>
        <el-descriptions-item label="当前实例数">{{ selectedGroup.current_capacity }}</el-descriptions-item>
        <el-descriptions-item label="期望实例数">{{ selectedGroup.desired_capacity }}</el-descriptions-item>
        <el-descriptions-item label="最小实例数">{{ selectedGroup.min_size }}</el-descriptions-item>
        <el-descriptions-item label="最大实例数">{{ selectedGroup.max_size }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedGroup.region_id }}</el-descriptions-item>
        <el-descriptions-item label="描述">{{ selectedGroup.description || '无' }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Plus, ArrowDown } from '@element-plus/icons-vue'
import {
  getAutoscalingGroups,
  createAutoscalingGroup,
  updateAutoscalingGroup,
  deleteAutoscalingGroup,
  AutoscalingGroup,
  getHostTemplates
} from '@/api/compute'
import { getProjects } from '@/api/iam'

// Data
const loading = ref(false)
const submitting = ref(false)
const autoscalingGroups = ref<AutoscalingGroup[]>([])
const selectedGroups = ref<AutoscalingGroup[]>([])
const dialogVisible = ref(false)
const scaleDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const dialogType = ref<'create' | 'edit'>('create')
const selectedGroup = ref<AutoscalingGroup | null>(null)
const hostTemplates = ref<any[]>([])
const projects = ref<any[]>([])
const networks = ref<any[]>([])

const formRef = ref()

const filters = reactive({
  project_id: '',
  name: '',
  platform: '',
  status: ''
})

const form = reactive({
  id: '',
  name: '',
  description: '',
  status: 'Inactive',
  host_template_id: '',
  network_id: '',
  current_capacity: 0,
  desired_capacity: 1,
  min_size: 0,
  max_size: 10,
  platform: 'aliyun',
  project_id: '',
  region_id: '',
  zone_id: '',
  removal_policy: 'oldest',
  load_balancing: 'disabled',
  health_check_type: 'EC2',
  health_check_period: 60,
  health_check_grace_period: 300,
  tags: [] as string[]
})

const scaleForm = reactive({
  capacity: 0
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const availableTags = ref(['production', 'development', 'test', 'critical'])

const rules = {
  name: [{ required: true, message: '请输入伸缩组名称', trigger: 'blur' }],
  host_template_id: [{ required: true, message: '请选择主机模版', trigger: 'change' }],
  project_id: [{ required: true, message: '请选择项目', trigger: 'change' }],
  platform: [{ required: true, message: '请选择平台', trigger: 'change' }]
}

// Computed
const dialogTitle = computed(() => dialogType.value === 'create' ? '创建弹性伸缩组' : '编辑弹性伸缩组')

// Methods
const getStatusType = (status: string) => {
  switch (status) {
    case 'Active': return 'success'
    case 'Inactive': return 'info'
    case 'Deleting': return 'warning'
    case 'Error': return 'danger'
    default: return 'info'
  }
}

const getPlatformLabel = (platform: string) => {
  const labels: Record<string, string> = {
    aliyun: '阿里云',
    tencent: '腾讯云',
    aws: 'AWS',
    azure: 'Azure'
  }
  return labels[platform] || platform
}

const getTemplateName = (id: string) => {
  const template = hostTemplates.value.find(t => t.id === id)
  return template?.name || id || '-'
}

const handleSelectionChange = (selection: AutoscalingGroup[]) => {
  selectedGroups.value = selection
}

const handleView = () => {
  ElMessage.info('视图切换功能开发中')
}

const resetFilters = () => {
  filters.project_id = ''
  filters.name = ''
  filters.platform = ''
  filters.status = ''
  pagination.page = 1
  fetchData()
}

const fetchData = async () => {
  loading.value = true
  try {
    const response = await getAutoscalingGroups({
      ...filters,
      page: pagination.page,
      page_size: pagination.pageSize
    })
    autoscalingGroups.value = response.data?.items || response.items || []
    pagination.total = response.data?.pagination?.total || response.pagination?.total || 0
  } catch (error) {
    console.error('Failed to fetch:', error)
    ElMessage.error('获取弹性伸缩组列表失败')
  } finally {
    loading.value = false
  }
}

const fetchHostTemplates = async () => {
  try {
    const res = await getHostTemplates()
    hostTemplates.value = res.items || res.data?.items || []
  } catch (e) {
    hostTemplates.value = []
  }
}

const fetchProjects = async () => {
  try {
    const res = await getProjects()
    projects.value = res.items || []
  } catch (e) {
    projects.value = []
  }
}

const handleBatchCommand = async (command: string) => {
  if (selectedGroups.value.length === 0) return

  const actionNames = { enable: '启用', disable: '禁用', delete: '删除' }
  try {
    await ElMessageBox.confirm(
      `确定要批量${actionNames[command]}选中的 ${selectedGroups.value.length} 个伸缩组吗？`,
      '批量操作确认',
      { type: 'warning' }
    )

    for (const group of selectedGroups.value) {
      if (command === 'delete') {
        await deleteAutoscalingGroup(group.id)
      } else {
        const newStatus = command === 'enable' ? 'Active' : 'Inactive'
        await updateAutoscalingGroup(group.id, { ...group, status: newStatus })
      }
    }
    ElMessage.success(`批量${actionNames[command]}完成`)
    selectedGroups.value = []
    fetchData()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('批量操作失败')
  }
}

const handleActionCommand = (command: string, row: AutoscalingGroup) => {
  switch (command) {
    case 'edit': handleEdit(row); break
    case 'scaleUp': handleScaleUp(row); break
    case 'scaleDown': handleScaleDown(row); break
    case 'enable': handleToggleStatus(row, 'Active'); break
    case 'disable': handleToggleStatus(row, 'Inactive'); break
    case 'details': handleDetails(row); break
    case 'delete': handleDelete(row); break
  }
}

const handleCreate = () => {
  dialogType.value = 'create'
  Object.assign(form, {
    id: '', name: '', description: '', status: 'Inactive',
    host_template_id: '', network_id: '', current_capacity: 0,
    desired_capacity: 1, min_size: 0, max_size: 10,
    platform: 'aliyun', project_id: '', region_id: '', zone_id: '',
    removal_policy: 'oldest', load_balancing: 'disabled',
    health_check_type: 'EC2', health_check_period: 60,
    health_check_grace_period: 300, tags: []
  })
  dialogVisible.value = true
}

const handleEdit = (row: AutoscalingGroup) => {
  dialogType.value = 'edit'
  Object.assign(form, { ...row })
  dialogVisible.value = true
}

const handleDetails = (row: AutoscalingGroup) => {
  selectedGroup.value = row
  detailDialogVisible.value = true
}

const handleDelete = async (row: AutoscalingGroup) => {
  try {
    await ElMessageBox.confirm(`确定要删除 "${row.name}" 吗？`, '警告', { type: 'warning' })
    await deleteAutoscalingGroup(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('删除失败')
  }
}

const handleScaleUp = (row: AutoscalingGroup) => {
  selectedGroup.value = row
  scaleForm.capacity = Math.min(row.desired_capacity + 1, row.max_size)
  scaleDialogVisible.value = true
}

const handleScaleDown = (row: AutoscalingGroup) => {
  selectedGroup.value = row
  scaleForm.capacity = Math.max(row.desired_capacity - 1, row.min_size)
  scaleDialogVisible.value = true
}

const handleToggleStatus = async (row: AutoscalingGroup, newStatus: string) => {
  try {
    await ElMessageBox.confirm(`确定要${newStatus === 'Active' ? '启用' : '禁用'} "${row.name}" 吗？`, '提示', { type: 'warning' })
    await updateAutoscalingGroup(row.id, { ...row, status: newStatus })
    ElMessage.success(newStatus === 'Active' ? '启用成功' : '禁用成功')
    fetchData()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('操作失败')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    submitting.value = true
    if (dialogType.value === 'create') {
      await createAutoscalingGroup({ ...form })
      ElMessage.success('创建成功')
    } else {
      await updateAutoscalingGroup(form.id, { ...form })
      ElMessage.success('更新成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (e) {
    ElMessage.error('提交失败')
  } finally {
    submitting.value = false
  }
}

const confirmScale = async () => {
  if (!selectedGroup.value) return
  try {
    await updateAutoscalingGroup(selectedGroup.value.id, {
      ...selectedGroup.value,
      desired_capacity: scaleForm.capacity
    })
    ElMessage.success('调整实例数成功')
    scaleDialogVisible.value = false
    fetchData()
  } catch (e) {
    ElMessage.error('调整实例数失败')
  }
}

const handleDialogClose = () => { dialogVisible.value = false }
const handleSizeChange = (val: number) => { pagination.pageSize = val; pagination.page = 1; fetchData() }
const handleCurrentChange = (val: number) => { pagination.page = val; fetchData() }

onMounted(() => {
  fetchData()
  fetchHostTemplates()
  fetchProjects()
  networks.value = [{ id: 'n-1', name: '默认VPC' }]
})
</script>

<style scoped>
.autoscaling-group-container { padding: 20px; }

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
.capacity-cell { font-weight: 500; }
.pagination { margin-top: 20px; text-align: right; }
.dialog-footer { display: flex; justify-content: flex-end; gap: 10px; }
</style>