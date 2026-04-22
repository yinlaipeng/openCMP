<template>
  <div class="ssh-keys-container">
    <div class="page-header">
      <h2>密钥</h2>
      <div class="toolbar">
        <el-button link type="primary" @click="handleView">
          <el-icon><View /></el-icon>
          查看
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建
        </el-button>
        <el-dropdown trigger="click" @command="handleBatchCommand" :disabled="selectedKeys.length === 0">
          <el-button :disabled="selectedKeys.length === 0">
            批量操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="delete">批量删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- 搜索表单 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="fetchData">
        <el-form-item label="名称">
          <el-input v-model="filters.name" placeholder="搜索密钥名称" clearable style="width: 180px" />
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
        <el-form-item label="区域">
          <el-input v-model="filters.region" placeholder="区域" clearable style="width: 120px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <el-table
        :data="sshKeys"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="名称" width="200" fixed="left">
          <template #default="{ row }">
            <el-link type="primary" :underline="false" @click="handleDetails(row)">
              {{ row.name }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="public_scope" label="共享范围" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="row.public_scope === 'system' ? 'primary' : 'info'">
              {{ row.public_scope || '系统' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="public_key" label="公钥内容" width="300" show-overflow-tooltip />
        <el-table-column prop="fingerprint" label="指纹" width="200" show-overflow-tooltip />
        <el-table-column prop="keypair_type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ row.keypair_type || 'SSH' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="guest_cnt" label="关联主机数" width="120" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleDetails(row)">查看</el-button>
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

    <!-- 创建密钥弹窗 -->
    <el-dialog
      title="创建密钥"
      v-model="createDialogVisible"
      width="500px"
    >
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="100px">
        <el-form-item label="云账户" prop="cloud_account_id">
          <el-select v-model="createForm.cloud_account_id" placeholder="选择云账户" style="width: 100%" @change="handleAccountChange">
            <el-option v-for="acc in cloudAccounts" :key="acc.id" :label="acc.name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入密钥名称" />
        </el-form-item>
        <el-form-item label="公钥内容">
          <el-input v-model="createForm.public_key" type="textarea" :rows="4" placeholder="粘贴公钥内容或自动生成" />
        </el-form-item>
        <el-form-item label="自动生成">
          <el-switch v-model="createForm.auto_generate" />
          <span style="margin-left: 10px; color: #999;">开启后将自动生成密钥对</span>
        </el-form-item>
        <el-form-item label="项目">
          <el-select v-model="createForm.project_id" placeholder="选择项目" clearable style="width: 100%">
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
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
      v-model="detailDialogVisible"
      title="密钥详情"
      width="600px"
    >
      <el-descriptions :column="2" border v-if="selectedKey">
        <el-descriptions-item label="名称">{{ selectedKey.name }}</el-descriptions-item>
        <el-descriptions-item label="共享范围">{{ selectedKey.public_scope || '系统' }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ selectedKey.keypair_type || 'SSH' }}</el-descriptions-item>
        <el-descriptions-item label="关联主机数">{{ selectedKey.guest_cnt || 0 }}</el-descriptions-item>
        <el-descriptions-item label="平台">
          <el-tag size="small" :type="getPlatformType(selectedKey.provider_type)">
            {{ getPlatformLabel(selectedKey.provider_type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="云账户">{{ selectedKey.account_name }}</el-descriptions-item>
        <el-descriptions-item label="指纹">{{ selectedKey.fingerprint }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedKey.region }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedKey.project_name }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ selectedKey.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="公钥内容" :span="2">
          <div style="word-break: break-all; font-size: 12px;">
            {{ selectedKey.public_key }}
          </div>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Plus, ArrowDown } from '@element-plus/icons-vue'
import { getKeyPairs, createKeyPair, deleteKeyPair, batchDeleteKeyPairs, KeyPair } from '@/api/networkSync'
import { getCloudAccounts } from '@/api/cloud-account'
import type { CloudAccount } from '@/types'

// Data
const loading = ref(false)
const creating = ref(false)
const sshKeys = ref<KeyPair[]>([])
const selectedKeys = ref<KeyPair[]>([])
const cloudAccounts = ref<CloudAccount[]>([])
const projects = ref<any[]>([])

const createDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const selectedKey = ref<KeyPair | null>(null)

const filters = reactive({
  name: '',
  status: '',
  platform: '',
  cloud_account_id: null as number | null,
  region: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const createForm = reactive({
  cloud_account_id: '',
  name: '',
  public_key: '',
  auto_generate: true,
  project_id: ''
})

const createRules = {
  cloud_account_id: [{ required: true, message: '请选择云账户', trigger: 'change' }],
  name: [{ required: true, message: '请输入密钥名称', trigger: 'blur' }]
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

const formatDate = (dateStr: string | undefined) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}

const handleSelectionChange = (selection: KeyPair[]) => {
  selectedKeys.value = selection
}

const resetFilters = () => {
  filters.name = ''
  filters.status = ''
  filters.platform = ''
  filters.cloud_account_id = null
  filters.region = ''
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
    const res = await getKeyPairs(params)
    sshKeys.value = res.items || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('Failed to fetch SSH keys:', error)
    ElMessage.error('获取密钥列表失败')
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

const handleAccountChange = () => {
  // 根据云账户变化
}

const handleView = () => {
  ElMessage.info('视图切换功能开发中')
}

const handleCreate = () => {
  Object.assign(createForm, {
    cloud_account_id: '',
    name: '',
    public_key: '',
    auto_generate: true,
    project_id: ''
  })
  createDialogVisible.value = true
}

const confirmCreate = async () => {
  creating.value = true
  try {
    await createKeyPair({
      cloud_account_id: Number(createForm.cloud_account_id),
      name: createForm.name,
      public_key: createForm.public_key,
      auto_generate: createForm.auto_generate,
      project_id: createForm.project_id ? Number(createForm.project_id) : undefined
    })
    ElMessage.success('密钥创建成功')
    createDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('创建失败')
  } finally {
    creating.value = false
  }
}

const handleBatchCommand = async (command: string) => {
  if (selectedKeys.value.length === 0) return
  try {
    await ElMessageBox.confirm(
      `确定要批量删除选中的 ${selectedKeys.value.length} 个密钥吗？`,
      '批量操作确认',
      { type: 'warning' }
    )
    if (command === 'delete') {
      const ids = selectedKeys.value.map(k => k.id)
      await batchDeleteKeyPairs(ids)
      ElMessage.success('批量删除成功')
    }
    selectedKeys.value = []
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('批量操作失败')
  }
}

const handleEdit = (row: KeyPair) => {
  ElMessage.info('编辑功能开发中')
}

const handleDetails = (row: KeyPair) => {
  selectedKey.value = row
  detailDialogVisible.value = true
}

const handleDelete = async (row: KeyPair) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除SSH密钥 "${row.name}" 吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await deleteKeyPair(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    console.error('Failed to delete SSH key:', error)
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
.ssh-keys-container {
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
</style>