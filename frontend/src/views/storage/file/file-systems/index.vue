<template>
  <div class="file-systems-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">文件系统列表</span>
          <div class="toolbar">
            <el-button @click="handleView">查看</el-button>
            <el-button type="primary" @click="handleCreate">创建</el-button>
            <el-dropdown :disabled="!hasSelection" @command="handleBatchCommand">
              <el-button>
                批量操作 <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="delete">批量删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button :disabled="!hasSelection" @click="handleTags">标签</el-button>
          </div>
        </div>
      </template>

      <!-- Search Filters -->
      <el-form :inline="true" class="search-form">
        <el-form-item label="名称">
          <el-input v-model="searchParams.name" placeholder="请输入名称" clearable @clear="loadFileSystems" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchParams.status" placeholder="请选择状态" clearable @clear="loadFileSystems">
            <el-option label="运行中" value="active" />
            <el-option label="创建中" value="creating" />
            <el-option label="已停止" value="stopped" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="searchParams.platform" placeholder="请选择平台" clearable @clear="loadFileSystems">
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-input v-model="searchParams.region" placeholder="请输入区域" clearable @clear="loadFileSystems" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadFileSystems">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="fileSystems" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column label="名称" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="tags" label="标签" width="150" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="fs_type" label="文件系统类型" width="120" />
        <el-table-column prop="storage_type" label="存储类型" width="120" />
        <el-table-column prop="protocol" label="协议类型" width="100" />
        <el-table-column prop="billing_method" label="计费方式" width="120" />
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column prop="domain" label="所属域" width="150" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="searchParams.page"
        v-model:page-size="searchParams.page_size"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        @size-change="loadFileSystems"
        @current-change="loadFileSystems"
      />
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="文件系统详情" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedFileSystem?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedFileSystem?.name }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{ selectedFileSystem?.tags }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedFileSystem?.status }}</el-descriptions-item>
        <el-descriptions-item label="文件系统类型">{{ selectedFileSystem?.fs_type }}</el-descriptions-item>
        <el-descriptions-item label="存储类型">{{ selectedFileSystem?.storage_type }}</el-descriptions-item>
        <el-descriptions-item label="协议类型">{{ selectedFileSystem?.protocol }}</el-descriptions-item>
        <el-descriptions-item label="计费方式">{{ selectedFileSystem?.billing_method }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedFileSystem?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedFileSystem?.account }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ selectedFileSystem?.domain }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedFileSystem?.region }}</el-descriptions-item>
        <el-descriptions-item label="容量">{{ selectedFileSystem?.capacity }} GB</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedFileSystem?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Create Modal -->
    <el-dialog v-model="createDialogVisible" title="创建文件系统" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="createForm.name" placeholder="请输入文件系统名称" />
        </el-form-item>
        <el-form-item label="存储类型">
          <el-select v-model="createForm.storage_type" placeholder="请选择存储类型">
            <el-option label="性能型" value="performance" />
            <el-option label="容量型" value="capacity" />
          </el-select>
        </el-form-item>
        <el-form-item label="协议类型">
          <el-select v-model="createForm.protocol" placeholder="请选择协议类型">
            <el-option label="NFS" value="nfs" />
            <el-option label="CIFS" value="cifs" />
          </el-select>
        </el-form-item>
        <el-form-item label="计费方式">
          <el-select v-model="createForm.billing_method" placeholder="请选择计费方式">
            <el-option label="按量付费" value="postpaid" />
            <el-option label="包年包月" value="prepaid" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCreate">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'

interface FileSystem {
  id: string
  name: string
  tags: string
  status: string
  fs_type: string
  storage_type: string
  protocol: string
  billing_method: string
  platform: string
  account: string
  domain: string
  region: string
  capacity?: number
  created_at: string
}

const fileSystems = ref<FileSystem[]>([])
const loading = ref(false)
const total = ref(0)
const selectedRows = ref<FileSystem[]>([])
const detailDialogVisible = ref(false)
const createDialogVisible = ref(false)
const selectedFileSystem = ref<FileSystem | null>(null)

const searchParams = ref({
  page: 1,
  page_size: 10,
  name: '',
  status: '',
  platform: '',
  region: ''
})

const createForm = ref({
  name: '',
  storage_type: '',
  protocol: '',
  billing_method: ''
})

const hasSelection = computed(() => selectedRows.value.length > 0)

const handleSelectionChange = (rows: FileSystem[]) => {
  selectedRows.value = rows
}

const resetSearch = () => {
  searchParams.value = {
    page: 1,
    page_size: 10,
    name: '',
    status: '',
    platform: '',
    region: ''
  }
  loadFileSystems()
}

const handleView = () => {
  if (selectedRows.value.length === 1) {
    viewDetail(selectedRows.value[0])
  } else {
    ElMessage.warning('请选择一条记录查看')
  }
}

const handleCreate = () => {
  createForm.value = { name: '', storage_type: '', protocol: '', billing_method: '' }
  createDialogVisible.value = true
}

const submitCreate = async () => {
  try {
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    loadFileSystems()
  } catch (e) {
    console.error(e)
    ElMessage.error('创建失败')
  }
}

const handleBatchCommand = async (command: string) => {
  if (!hasSelection.value) return
  if (command === 'delete') {
    try {
      await ElMessageBox.confirm(`确认删除选中的 ${selectedRows.value.length} 个文件系统？`, '警告', { type: 'warning' })
      fileSystems.value = fileSystems.value.filter(f => !selectedRows.value.some(r => r.id === f.id))
      ElMessage.success('删除成功')
      loadFileSystems()
    } catch (e) {
      console.error(e)
    }
  }
}

const handleTags = () => {
  if (!hasSelection.value) {
    ElMessage.warning('请选择需要设置标签的文件系统')
    return
  }
  ElMessage.info('标签功能开发中')
}

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
    case 'running':
      return 'success'
    case 'pending':
    case 'creating':
      return 'warning'
    case 'error':
    case 'stopped':
      return 'danger'
    default:
      return 'info'
  }
}

const loadFileSystems = async () => {
  loading.value = true
  try {
    // Mock data
    fileSystems.value = [
      {
        id: 'fs-1',
        name: 'nas-prod',
        tags: 'prod',
        status: 'Active',
        fs_type: 'NAS',
        storage_type: '性能型',
        protocol: 'NFS',
        billing_method: '包年包月',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        domain: 'Default Domain',
        region: 'cn-hangzhou',
        capacity: 1000,
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'fs-2',
        name: 'nas-dev',
        tags: 'dev',
        status: 'Active',
        fs_type: 'NAS',
        storage_type: '容量型',
        protocol: 'CIFS',
        billing_method: '按量付费',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        domain: 'Default Domain',
        region: 'cn-shanghai',
        capacity: 500,
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'fs-3',
        name: 'nas-test',
        tags: 'test',
        status: 'Creating',
        fs_type: 'NAS',
        storage_type: '性能型',
        protocol: 'NFS',
        billing_method: '按量付费',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        domain: 'Domain A',
        region: 'cn-beijing',
        capacity: 200,
        created_at: '2024-01-03 10:00:00'
      }
    ]
    total.value = 3
  } catch (e) {
    console.error(e)
    fileSystems.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: FileSystem) => {
  selectedFileSystem.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row: FileSystem) => {
  ElMessage.info(`编辑文件系统: ${row.name}`)
}

const handleDelete = async (row: FileSystem) => {
  try {
    await ElMessageBox.confirm(`确认删除文件系统 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    fileSystems.value = fileSystems.value.filter(f => f.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadFileSystems()
})
</script>

<style scoped>
.file-systems-page {
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

.toolbar {
  display: flex;
  gap: 8px;
}

.search-form {
  margin-bottom: 16px;
}
</style>