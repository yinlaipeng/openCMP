<template>
  <div class="host-template-container">
    <div class="page-header">
      <h2>主机模版</h2>
      <div class="toolbar">
        <el-button link type="primary" @click="handleView">
          <el-icon><View /></el-icon>
          查看
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          新建模版
        </el-button>
        <el-button :disabled="selectedTemplates.length === 0" @click="handleBatchDelete">
          <el-icon><Delete /></el-icon>
          删除
        </el-button>
      </div>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="fetchHostTemplates">
        <el-form-item label="项目">
          <el-select v-model="filters.project_id" placeholder="请选择项目" clearable style="width: 160px">
            <el-option
              v-for="project in projects"
              :key="project.id"
              :label="project.name"
              :value="project.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="filters.name" placeholder="搜索模版名称" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="filters.platform" placeholder="选择平台" clearable style="width: 120px">
            <el-option label="阿里云" value="alibaba" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="选择状态" clearable style="width: 120px">
            <el-option label="激活" value="Active" />
            <el-option label="未激活" value="Inactive" />
            <el-option label="草稿" value="Draft" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchHostTemplates">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table
      :data="hostTemplates"
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

      <el-table-column label="配置" width="180">
        <template #default="{ row }">
          <div class="config-cell">
            <span>{{ row.cpu_count }}核</span>
            <span>{{ row.memory_size }}MB</span>
            <span>{{ row.disk_size }}GB</span>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="OS镜像" width="150">
        <template #default="{ row }">
          <div>{{ row.os_name }} {{ row.os_version }}</div>
        </template>
      </el-table-column>

      <el-table-column prop="billing_method" label="计费方式" width="100" />

      <el-table-column prop="region_id" label="区域" width="120" />

      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button size="small" link type="primary" @click="handleEdit(row)">编辑</el-button>
          <el-dropdown trigger="click" @command="(cmd: string) => handleActionCommand(cmd, row)">
            <el-button size="small" link>
              更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="deploy">基于此模版部署</el-dropdown-item>
                <el-dropdown-item command="duplicate">复制模版</el-dropdown-item>
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
      :total="pagination.total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      class="pagination"
    />

    <!-- 编辑对话框 -->
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
        @submit.prevent
      >
        <el-form-item label="模版名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入模版名称" />
        </el-form-item>

        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            placeholder="请输入模版描述"
            :rows="3"
          />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status" placeholder="请选择状态">
            <el-option label="激活" value="Active" />
            <el-option label="未激活" value="Inactive" />
            <el-option label="草稿" value="Draft" />
          </el-select>
        </el-form-item>

        <el-form-item label="实例规格" prop="instance_type">
          <el-input v-model="form.instance_type" placeholder="例如：ecs.c6.large" />
        </el-form-item>

        <el-form-item label="CPU架构" prop="cpu_arch">
          <el-select v-model="form.cpu_arch" placeholder="请选择CPU架构">
            <el-option label="x86_64" value="x86_64" />
            <el-option label="arm64" value="arm64" />
            <el-option label="aarch64" value="aarch64" />
          </el-select>
        </el-form-item>

        <el-form-item label="CPU核心数" prop="cpu_count">
          <el-input-number v-model="form.cpu_count" :min="1" :max="256" />
        </el-form-item>

        <el-form-item label="内存大小(MB)" prop="memory_size">
          <el-input-number v-model="form.memory_size" :min="1024" :step="1024" />
        </el-form-item>

        <el-form-item label="磁盘大小(GB)" prop="disk_size">
          <el-input-number v-model="form.disk_size" :min="10" :max="10000" />
        </el-form-item>

        <el-form-item label="镜像ID" prop="image_id">
          <el-input v-model="form.image_id" placeholder="请输入镜像ID" />
        </el-form-item>

        <el-form-item label="操作系统" prop="os_name">
          <el-input v-model="form.os_name" placeholder="例如：Ubuntu, CentOS" />
        </el-form-item>

        <el-form-item label="操作系统版本" prop="os_version">
          <el-input v-model="form.os_version" placeholder="例如：20.04, 7.8" />
        </el-form-item>

        <el-form-item label="VPC ID" prop="vpc_id">
          <el-input v-model="form.vpc_id" placeholder="请输入VPC ID" />
        </el-form-item>

        <el-form-item label="子网ID" prop="subnet_id">
          <el-input v-model="form.subnet_id" placeholder="请输入子网ID" />
        </el-form-item>

        <el-form-item label="计费方式" prop="billing_method">
          <el-select v-model="form.billing_method" placeholder="请选择计费方式">
            <el-option label="按量付费" value="Pay-As-You-Go" />
            <el-option label="包年包月" value="Subscription" />
            <el-option label="预留实例" value="Reserved" />
          </el-select>
        </el-form-item>

        <el-form-item label="平台" prop="platform">
          <el-select v-model="form.platform" placeholder="请选择平台">
            <el-option label="阿里云" value="alibaba" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>

        <el-form-item label="项目" prop="project_id">
          <el-select v-model="form.project_id" placeholder="请选择项目" clearable>
            <el-option
              v-for="project in projects"
              :key="project.id"
              :label="project.name"
              :value="project.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="区域ID" prop="region_id">
          <el-input v-model="form.region_id" placeholder="请输入区域ID" />
        </el-form-item>

        <el-form-item label="可用区ID" prop="zone_id">
          <el-input v-model="form.zone_id" placeholder="请输入可用区ID" />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="handleDialogClose">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus'
import { ArrowDown, View, Plus, Delete } from '@element-plus/icons-vue'
import {
  getHostTemplates,
  createHostTemplate,
  updateHostTemplate,
  deleteHostTemplate
} from '@/api/compute'
import { getProjects } from '@/api/iam'
import type { HostTemplate } from '@/types'

// 响应式数据
const hostTemplates = ref<HostTemplate[]>([])
const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogType = ref<'create' | 'edit'>('create')
const selectedTemplates = ref<HostTemplate[]>([])

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 筛选条件
const filters = reactive({
  project_id: '',
  name: '',
  platform: '',
  status: ''
})

// 表单数据
const form = reactive<Partial<HostTemplate>>({
  name: '',
  description: '',
  status: 'Draft',
  instance_type: '',
  cpu_arch: 'x86_64',
  cpu_count: 2,
  memory_size: 4096,
  disk_size: 50,
  image_id: '',
  os_name: '',
  os_version: '',
  vpc_id: '',
  subnet_id: '',
  billing_method: 'Pay-As-You-Go',
  platform: 'alibaba',
  project_id: '',
  region_id: '',
  zone_id: ''
})

// 项目列表
const projects = ref<any[]>([])

// 表单验证规则
const rules: FormRules = {
  name: [
    { required: true, message: '请输入模版名称', trigger: 'blur' },
    { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  instance_type: [
    { required: true, message: '请输入实例规格', trigger: 'blur' }
  ],
  image_id: [
    { required: true, message: '请输入镜像ID', trigger: 'blur' }
  ],
  os_name: [
    { required: true, message: '请输入操作系统名称', trigger: 'blur' }
  ],
  vpc_id: [
    { required: true, message: '请输入VPC ID', trigger: 'blur' }
  ],
  platform: [
    { required: true, message: '请选择平台', trigger: 'change' }
  ],
  project_id: [
    { required: true, message: '请选择项目', trigger: 'change' }
  ]
}

// 表单引用
const formRef = ref<FormInstance>()

// 计算属性
const dialogTitle = computed(() => {
  return dialogType.value === 'create' ? '创建主机模版' : '编辑主机模版'
})

// 获取主机模版列表
const fetchHostTemplates = async () => {
  loading.value = true
  try {
    const response = await getHostTemplates({
      ...filters,
      page: pagination.page,
      page_size: pagination.pageSize
    })

    hostTemplates.value = response.items || []
    pagination.total = response.pagination?.total || 0
  } catch (error) {
    console.error('获取主机模版列表失败:', error)
    ElMessage.error('获取主机模版列表失败')
  } finally {
    loading.value = false
  }
}

// 获取项目列表
const fetchProjects = async () => {
  try {
    const response = await getProjects()
    projects.value = response.items || []
  } catch (error) {
    console.error('获取项目列表失败:', error)
    ElMessage.error('获取项目列表失败')
  }
}

// 获取状态类型
const getStatusType = (status: string) => {
  switch (status) {
    case 'Active':
      return 'success'
    case 'Inactive':
      return 'info'
    case 'Draft':
      return 'warning'
    default:
      return 'info'
  }
}

// 获取平台标签
const getPlatformLabel = (platform: string) => {
  const labels: Record<string, string> = {
    alibaba: '阿里云',
    tencent: '腾讯云',
    aws: 'AWS',
    azure: 'Azure'
  }
  return labels[platform] || platform
}

// 处理选择变化
const handleSelectionChange = (selection: HostTemplate[]) => {
  selectedTemplates.value = selection
}

// 处理查看
const handleView = () => {
  ElMessage.info('视图切换功能开发中')
  }

// 处理批量删除
const handleBatchDelete = async () => {
  if (selectedTemplates.value.length === 0) {
    ElMessage.warning('请先选择要删除的模版')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确认删除选中的 ${selectedTemplates.value.length} 个模版吗？`,
      '批量删除',
      { type: 'warning' }
    )

    for (const template of selectedTemplates.value) {
      await deleteHostTemplate(template.id)
    }
    ElMessage.success('批量删除成功')
    selectedTemplates.value = []
    fetchHostTemplates()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
  }

// 处理操作命令
const handleActionCommand = (command: string, row: HostTemplate) => {
  switch (command) {
    case 'deploy':
      deployFromTemplate(row)
      break
    case 'duplicate':
      duplicateTemplate(row)
      break
    case 'details':
      handleEdit(row)
      break
    case 'delete':
      handleDelete(row)
      break
  }
  }

// 重置筛选条件
const resetFilters = () => {
  filters.project_id = ''
  filters.name = ''
  filters.platform = ''
  filters.status = ''
  pagination.page = 1
  fetchHostTemplates()
  }

// 处理新建
const handleCreate = () => {
  dialogType.value = 'create'
  Object.assign(form, {
    name: '',
    description: '',
    status: 'Draft',
    instance_type: '',
    cpu_arch: 'x86_64',
    cpu_count: 2,
    memory_size: 4096,
    disk_size: 50,
    image_id: '',
    os_name: '',
    os_version: '',
    vpc_id: '',
    subnet_id: '',
    billing_method: 'Pay-As-You-Go',
    platform: 'alibaba',
    project_id: '',
    region_id: '',
    zone_id: ''
  })
  dialogVisible.value = true
}

// 处理编辑
const handleEdit = (row: HostTemplate) => {
  dialogType.value = 'edit'
  Object.assign(form, { ...row })
  dialogVisible.value = true
}

// 处理删除
const handleDelete = async (row: HostTemplate) => {
  try {
    await ElMessageBox.confirm(
      `确认删除主机模版 "${row.name}" 吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await deleteHostTemplate(row.id)
    ElMessage.success('删除成功')
    fetchHostTemplates()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除主机模版失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 基于此模版部署
const deployFromTemplate = (row: HostTemplate) => {
  ElMessage.info(`基于模版 "${row.name}" 部署虚拟机功能开发中`)
  // 这里可以跳转到部署页面或打开部署对话框
}

// 复制模版
const duplicateTemplate = (row: HostTemplate) => {
  dialogType.value = 'create'
  const duplicatedForm = { ...row }
  duplicatedForm.name = `${row.name}-副本`
  duplicatedForm.id = undefined as any
  Object.assign(form, duplicatedForm)
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitLoading.value = true
  try {
    if (dialogType.value === 'create') {
      await createHostTemplate(form as Partial<HostTemplate>)
      ElMessage.success('创建成功')
    } else {
      if (!form.id) return
      await updateHostTemplate(form.id, form as Partial<HostTemplate>)
      ElMessage.success('更新成功')
    }

    dialogVisible.value = false
    fetchHostTemplates()
  } catch (error) {
    console.error('提交主机模版失败:', error)
    ElMessage.error(dialogType.value === 'create' ? '创建失败' : '更新失败')
  } finally {
    submitLoading.value = false
  }
}

// 关闭对话框
const handleDialogClose = () => {
  dialogVisible.value = false
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

// 处理分页大小变化
const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchHostTemplates()
}

// 处理当前页变化
const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchHostTemplates()
}

// 初始化
onMounted(() => {
  fetchHostTemplates()
  fetchProjects()
})
</script>

<style scoped>
.host-template-container {
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

.config-cell {
  display: flex;
  gap: 8px;
  font-size: 13px;
}

.config-cell span {
  color: var(--text-color-secondary);
}

.pagination {
  margin-top: 20px;
  text-align: right;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>