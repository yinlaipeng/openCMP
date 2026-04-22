<template>
  <div class="global-vpc-container">
    <div class="page-header">
      <h2>全局VPC</h2>
      <div class="toolbar">
        <el-button link type="primary" @click="handleView">
          <el-icon><View /></el-icon>
          查看
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建
        </el-button>
        <el-dropdown trigger="click" @command="handleBatchCommand" :disabled="selectedGlobalVpcs.length === 0">
          <el-button :disabled="selectedGlobalVpcs.length === 0">
            批量操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="tags">设置标签</el-dropdown-item>
              <el-dropdown-item command="delete" divided>批量删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button @click="handleTags">
          <el-icon><PriceTag /></el-icon>
          标签
        </el-button>
      </div>
    </div>

    <!-- 搜索表单 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="fetchData">
        <el-form-item label="名称">
          <el-input v-model="filters.name" placeholder="搜索全局VPC名称" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="选择状态" clearable style="width: 120px">
            <el-option label="可用" value="available" />
            <el-option label="创建中" value="creating" />
            <el-option label="错误" value="error" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="filters.platform" placeholder="选择平台" clearable style="width: 120px">
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item label="云账户">
          <el-select v-model="filters.cloud_account_id" placeholder="选择云账户" clearable style="width: 150px">
            <el-option v-for="acc in cloudAccounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <el-table
        :data="globalVpcs"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <!-- Cloudpods 表头顺序: Name → Status → Tags → Platform → VPCs → Cloud account → Owner Domain -->
        <el-table-column prop="name" label="名称" width="180" fixed="left">
          <template #default="{ row }">
            <el-link type="primary" :underline="false" @click="handleDetails(row)">
              {{ row.name }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="标签" width="100">
          <template #default="{ row }">
            <el-tag v-for="tag in (row.tags || []).slice(0, 2)" :key="tag.key" size="small" class="tag-item">
              {{ tag.key }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="provider_type" label="平台" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="getPlatformType(row.provider_type)">
              {{ getPlatformLabel(row.provider_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="vpc_count" label="VPC数" width="100" />
        <el-table-column prop="account_name" label="云账户" width="150" />
        <el-table-column prop="owner_domain" label="所属域" width="150" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-dropdown trigger="click">
              <el-button size="small">
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleChangeDomain(row)">更改域</el-dropdown-item>
                  <el-dropdown-item divided @click="handleDelete(row)">删除</el-dropdown-item>
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
    </el-card>

    <!-- 创建全局VPC弹窗 -->
    <el-dialog
      title="创建全局VPC"
      v-model="createDialogVisible"
      width="500px"
    >
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="100px">
        <el-form-item label="云账户" prop="cloud_account_id">
          <el-select v-model="createForm.cloud_account_id" placeholder="选择云账户" style="width: 100%">
            <el-option v-for="acc in cloudAccounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入全局VPC名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="createForm.description" type="textarea" :rows="2" placeholder="请输入描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCreate" :loading="creating">创建</el-button>
      </template>
    </el-dialog>

    <!-- 详情弹窗 -->
    <el-dialog
      title="全局VPC详情"
      v-model="detailDialogVisible"
      width="700px"
    >
      <el-descriptions :column="2" border v-if="selectedGlobalVpc">
        <el-descriptions-item label="ID">{{ selectedGlobalVpc.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedGlobalVpc.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(selectedGlobalVpc.status)">{{ getStatusLabel(selectedGlobalVpc.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="平台">
          <el-tag size="small" :type="getPlatformType(selectedGlobalVpc.provider_type)">
            {{ getPlatformLabel(selectedGlobalVpc.provider_type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="VPC数">{{ selectedGlobalVpc.vpc_count }}</el-descriptions-item>
        <el-descriptions-item label="云账户">{{ selectedGlobalVpc.account_name }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ selectedGlobalVpc.owner_domain }}</el-descriptions-item>
        <el-descriptions-item label="描述">{{ selectedGlobalVpc.description || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Plus, ArrowDown, PriceTag } from '@element-plus/icons-vue'
import { getGlobalVPCs, createGlobalVPC, deleteGlobalVPC, batchDeleteGlobalVPCs, GlobalVPC } from '@/api/networkSync'
import { getCloudAccounts } from '@/api/cloud-account'
import type { CloudAccount } from '@/types'

// Data
const loading = ref(false)
const creating = ref(false)
const globalVpcs = ref<GlobalVPC[]>([])
const selectedGlobalVpcs = ref<GlobalVPC[]>([])
const cloudAccounts = ref<CloudAccount[]>([])

const createDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const selectedGlobalVpc = ref<GlobalVPC | null>(null)

const filters = reactive({
  name: '',
  status: '',
  platform: '',
  cloud_account_id: null as number | null
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const createForm = reactive({
  cloud_account_id: '',
  name: '',
  description: ''
})

const createRules = {
  cloud_account_id: [{ required: true, message: '请选择云账户', trigger: 'change' }],
  name: [{ required: true, message: '请输入全局VPC名称', trigger: 'blur' }]
}

// Methods
const getStatusType = (status: string) => {
  switch (status) {
    case 'available': return 'success'
    case 'creating': return 'warning'
    case 'error': return 'danger'
    default: return 'info'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'available': return '可用'
    case 'creating': return '创建中'
    case 'error': return '错误'
    default: return status
  }
}

const getPlatformLabel = (platform: string) => {
  const labels: Record<string, string> = {
    aliyun: '阿里云',
    alibaba: '阿里云',
    tencent: '腾讯云',
    Qcloud: '腾讯云',
    aws: 'AWS',
    azure: 'Azure'
  }
  return labels[platform] || platform || '未知'
}

const getPlatformType = (platform: string) => {
  const types: Record<string, 'primary' | 'warning' | 'success' | 'info'> = {
    aliyun: 'primary',
    alibaba: 'primary',
    tencent: 'warning',
    Qcloud: 'warning',
    aws: 'success',
    azure: 'info'
  }
  return types[platform] || 'info'
}

const handleSelectionChange = (selection: GlobalVPC[]) => {
  selectedGlobalVpcs.value = selection
}

const resetFilters = () => {
  filters.name = ''
  filters.status = ''
  filters.platform = ''
  filters.cloud_account_id = null
  pagination.page = 1
  fetchData()
}

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      ...filters
    }
    const res = await getGlobalVPCs(params)
    globalVpcs.value = res.items || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('Failed to fetch global VPCs:', error)
    ElMessage.error('获取全局VPC列表失败')
  } finally {
    loading.value = false
  }
}

const fetchCloudAccounts = async () => {
  try {
    const res = await getCloudAccounts({ page: 1, page_size: 100 })
    cloudAccounts.value = res.items || []
  } catch (error) {
    console.error('Failed to fetch cloud accounts:', error)
  }
}

const handleView = () => {
  ElMessage.info('视图切换功能开发中')
}

const handleTags = () => {
  ElMessage.info('标签管理功能开发中')
}

const handleCreate = () => {
  Object.assign(createForm, {
    cloud_account_id: '',
    name: '',
    description: ''
  })
  createDialogVisible.value = true
}

const confirmCreate = async () => {
  creating.value = true
  try {
    await createGlobalVPC({
      cloud_account_id: Number(createForm.cloud_account_id),
      name: createForm.name,
      description: createForm.description
    })
    ElMessage.success('全局VPC创建成功')
    createDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('创建失败')
  } finally {
    creating.value = false
  }
}

const handleBatchCommand = async (command: string) => {
  if (selectedGlobalVpcs.value.length === 0) return
  const actionNames = { tags: '设置标签', delete: '删除' }
  try {
    await ElMessageBox.confirm(
      `确定要批量${actionNames[command]}选中的 ${selectedGlobalVpcs.value.length} 个全局VPC吗？`,
      '批量操作确认',
      { type: 'warning' }
    )
    if (command === 'delete') {
      const ids = selectedGlobalVpcs.value.map(v => v.id)
      await batchDeleteGlobalVPCs(ids)
      ElMessage.success('批量删除成功')
    } else {
      ElMessage.info(`批量${actionNames[command]}功能开发中`)
    }
    selectedGlobalVpcs.value = []
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('批量操作失败')
  }
}

const handleEdit = (row: GlobalVPC) => {
  ElMessage.info('编辑功能开发中')
}

const handleChangeDomain = (row: GlobalVPC) => {
  ElMessage.info('更改域功能开发中')
}

const handleDetails = (row: GlobalVPC) => {
  selectedGlobalVpc.value = row
  detailDialogVisible.value = true
}

const handleDelete = async (row: GlobalVPC) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除全局VPC "${row.name}" 吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await deleteGlobalVPC(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    console.error('Failed to delete global VPC:', error)
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
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

// Lifecycle
onMounted(() => {
  fetchData()
  fetchCloudAccounts()
})
</script>

<style scoped>
.global-vpc-container {
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

.table-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
}

.tag-item {
  margin-right: 4px;
  margin-bottom: 2px;
}
</style>