<template>
  <div class="autoscaling-group-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">弹性伸缩组</span>
          <el-button type="primary" @click="handleCreate">新建</el-button>
        </div>
      </template>

      <el-table
        :data="autoscalingGroups"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
      >
        <el-table-column prop="name" label="名称" width="200" fixed="left" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="host_template_id" label="主机模版" width="150" />
        <el-table-column prop="current_capacity" label="当前实例数" width="120" />
        <el-table-column prop="desired_capacity" label="期望实例数" width="120" />
        <el-table-column prop="min_size" label="最小实例数" width="120" />
        <el-table-column prop="max_size" label="最大实例数" width="120" />
        <el-table-column prop="platform" label="平台" width="120" />
        <el-table-column prop="project_id" label="项目" width="150" />
        <el-table-column label="操作" fixed="right" width="300">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleScaleUp(row)">扩容</el-button>
            <el-button size="small" @click="handleScaleDown(row)">缩容</el-button>
            <el-button size="small" @click="handleToggleStatus(row)" :type="row.status === 'Active' ? 'warning' : 'success'">
              {{ row.status === 'Active' ? '禁用' : '启用' }}
            </el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
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
    </el-card>

    <!-- Create/Edit Dialog -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="600px"
      :before-close="handleDialogClose"
    >
      <el-form
        :model="form"
        :rules="rules"
        ref="formRef"
        label-width="120px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入伸缩组名称" />
        </el-form-item>

        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            placeholder="请输入描述"
          />
        </el-form-item>

        <el-form-item label="主机模版" prop="host_template_id">
          <el-select v-model="form.host_template_id" placeholder="请选择主机模版" style="width: 100%">
            <el-option
              v-for="item in hostTemplates"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="最小实例数" prop="min_size">
          <el-input-number
            v-model="form.min_size"
            :min="0"
            :max="form.max_size"
            placeholder="请输入最小实例数"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="最大实例数" prop="max_size">
          <el-input-number
            v-model="form.max_size"
            :min="form.min_size"
            :max="100"
            placeholder="请输入最大实例数"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="期望实例数" prop="desired_capacity">
          <el-input-number
            v-model="form.desired_capacity"
            :min="form.min_size"
            :max="form.max_size"
            placeholder="请输入期望实例数"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="平台" prop="platform">
          <el-input v-model="form.platform" placeholder="请输入平台" />
        </el-form-item>

        <el-form-item label="项目ID" prop="project_id">
          <el-input v-model="form.project_id" placeholder="请输入项目ID" />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="handleDialogClose">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
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
        <el-form-item label="新实例数">
          <el-input-number
            v-model="scaleForm.capacity"
            :min="selectedAutoscalingGroup?.min_size"
            :max="selectedAutoscalingGroup?.max_size"
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getAutoscalingGroups,
  createAutoscalingGroup,
  updateAutoscalingGroup,
  deleteAutoscalingGroup,
  AutoscalingGroup
} from '@/api/compute'

// Data
const loading = ref(false)
const autoscalingGroups = ref<AutoscalingGroup[]>([])
const dialogVisible = ref(false)
const scaleDialogVisible = ref(false)
const dialogType = ref<'create' | 'edit'>('create')
const selectedAutoscalingGroup = ref<AutoscalingGroup | null>(null)
const hostTemplates = ref<any[]>([]) // Would be fetched from API

const formRef = ref()

const form = reactive({
  id: '',
  name: '',
  description: '',
  status: '',
  host_template_id: '',
  current_capacity: 0,
  desired_capacity: 1,
  min_size: 0,
  max_size: 10,
  platform: '',
  project_id: '',
  region_id: '',
  zone_id: ''
})

const scaleForm = reactive({
  capacity: 0
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const rules = {
  name: [
    { required: true, message: '请输入伸缩组名称', trigger: 'blur' }
  ],
  host_template_id: [
    { required: true, message: '请选择主机模版', trigger: 'change' }
  ],
  min_size: [
    { required: true, message: '请输入最小实例数', trigger: 'blur' }
  ],
  max_size: [
    { required: true, message: '请输入最大实例数', trigger: 'blur' }
  ],
  desired_capacity: [
    { required: true, message: '请输入期望实例数', trigger: 'blur' }
  ],
  platform: [
    { required: true, message: '请输入平台', trigger: 'blur' }
  ],
  project_id: [
    { required: true, message: '请输入项目ID', trigger: 'blur' }
  ]
}

// Computed
const dialogTitle = computed(() => {
  return dialogType.value === 'create' ? '创建弹性伸缩组' : '编辑弹性伸缩组'
})

// Methods
const getStatusType = (status: string) => {
  switch (status) {
    case 'Active':
      return 'success'
    case 'Inactive':
      return 'info'
    case 'Deleting':
      return 'warning'
    case 'Error':
      return 'danger'
    default:
      return 'info'
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const response = await getAutoscalingGroups({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    autoscalingGroups.value = response.data.items
    pagination.total = response.data.pagination.total
  } catch (error) {
    console.error('Failed to fetch autoscaling groups:', error)
    ElMessage.error('获取弹性伸缩组列表失败')
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  dialogType.value = 'create'
  Object.assign(form, {
    id: '',
    name: '',
    description: '',
    status: 'Inactive',
    host_template_id: '',
    current_capacity: 0,
    desired_capacity: 1,
    min_size: 0,
    max_size: 10,
    platform: '',
    project_id: '',
    region_id: '',
    zone_id: ''
  })
  dialogVisible.value = true
}

const handleEdit = (row: AutoscalingGroup) => {
  dialogType.value = 'edit'
  Object.assign(form, { ...row })
  dialogVisible.value = true
}

const handleDelete = async (row: AutoscalingGroup) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除弹性伸缩组 "${row.name}" 吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await deleteAutoscalingGroup(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    console.error('Failed to delete autoscaling group:', error)
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleScaleUp = (row: AutoscalingGroup) => {
  selectedAutoscalingGroup.value = row
  scaleForm.capacity = Math.min(row.desired_capacity + 1, row.max_size)
  scaleDialogVisible.value = true
}

const handleScaleDown = (row: AutoscalingGroup) => {
  selectedAutoscalingGroup.value = row
  scaleForm.capacity = Math.max(row.desired_capacity - 1, row.min_size)
  scaleDialogVisible.value = true
}

const handleToggleStatus = async (row: AutoscalingGroup) => {
  try {
    const newStatus = row.status === 'Active' ? 'Inactive' : 'Active'
    const actionText = newStatus === 'Active' ? '启用' : '禁用'

    await ElMessageBox.confirm(
      `确定要${actionText}弹性伸缩组 "${row.name}" 吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const updatedRow = { ...row, status: newStatus }
    await updateAutoscalingGroup(row.id, updatedRow)
    ElMessage.success(`${actionText}成功`)
    fetchData()
  } catch (error) {
    console.error(`Failed to ${row.status === 'Active' ? 'disable' : 'enable'} autoscaling group:`, error)
    if (error !== 'cancel') {
      ElMessage.error(`${row.status === 'Active' ? '禁用' : '启用'}失败`)
    }
  }
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()

    if (dialogType.value === 'create') {
      await createAutoscalingGroup({ ...form })
      ElMessage.success('创建成功')
    } else {
      await updateAutoscalingGroup(form.id, { ...form })
      ElMessage.success('更新成功')
    }

    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error('Submit failed:', error)
    ElMessage.error('提交失败')
  }
}

const confirmScale = async () => {
  if (!selectedAutoscalingGroup.value) return

  try {
    // This would be a specific scale API call if it existed
    const updatedRow = {
      ...selectedAutoscalingGroup.value,
      desired_capacity: scaleForm.capacity
    }
    await updateAutoscalingGroup(selectedAutoscalingGroup.value.id, updatedRow)
    ElMessage.success('调整实例数成功')
    scaleDialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error('Scale failed:', error)
    ElMessage.error('调整实例数失败')
  }
}

const handleDialogClose = () => {
  dialogVisible.value = false
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

// Lifecycle
onMounted(() => {
  fetchData()

  // Mock host templates - in real app this would come from API
  hostTemplates.value = [
    { id: 'ht-1', name: '通用型主机模版' },
    { id: 'ht-2', name: '计算型主机模版' },
    { id: 'ht-3', name: '内存型主机模版' }
  ]
})
</script>

<style scoped>
.autoscaling-group-container {
  padding: 20px;
}

.page-card {
  min-height: calc(100vh - 120px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 18px;
  font-weight: bold;
}

.pagination {
  margin-top: 20px;
  text-align: center;
}

.dialog-footer {
  text-align: right;
}
</style>