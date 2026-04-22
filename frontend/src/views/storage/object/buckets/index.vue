<template>
  <div class="buckets-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">存储桶列表</span>
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
          <el-input v-model="searchParams.name" placeholder="请输入名称" clearable @clear="loadBuckets" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchParams.status" placeholder="请选择状态" clearable @clear="loadBuckets">
            <el-option label="正常" value="active" />
            <el-option label="创建中" value="creating" />
            <el-option label="已删除" value="deleted" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="searchParams.platform" placeholder="请选择平台" clearable @clear="loadBuckets">
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-input v-model="searchParams.region" placeholder="请输入区域" clearable @clear="loadBuckets" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadBuckets">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="buckets" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column label="名称" width="200">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="tags" label="标签" width="150" />
        <el-table-column prop="permission" label="读写权限" width="120" />
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="project" label="项目" width="120" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
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
        @size-change="loadBuckets"
        @current-change="loadBuckets"
      />
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="存储桶详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedBucket?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedBucket?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedBucket?.status }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{ selectedBucket?.tags }}</el-descriptions-item>
        <el-descriptions-item label="读写权限">{{ selectedBucket?.permission }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedBucket?.platform }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedBucket?.project }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedBucket?.region }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedBucket?.created_at }}</el-descriptions-item>
        <el-descriptions-item label="存储用量">{{ selectedBucket?.storage_usage }} GB</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Create Modal -->
    <el-dialog v-model="createDialogVisible" title="创建存储桶" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="createForm.name" placeholder="请输入存储桶名称" />
        </el-form-item>
        <el-form-item label="读写权限">
          <el-select v-model="createForm.permission" placeholder="请选择读写权限">
            <el-option label="私有" value="private" />
            <el-option label="公共读" value="public-read" />
            <el-option label="公共读写" value="public-read-write" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-input v-model="createForm.region" placeholder="请输入区域" />
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

interface Bucket {
  id: string
  name: string
  status: string
  tags: string
  permission: string
  platform: string
  project: string
  region: string
  created_at: string
  storage_usage?: number
}

const buckets = ref<Bucket[]>([])
const loading = ref(false)
const total = ref(0)
const selectedRows = ref<Bucket[]>([])
const detailDialogVisible = ref(false)
const createDialogVisible = ref(false)
const selectedBucket = ref<Bucket | null>(null)

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
  permission: 'private',
  region: ''
})

const hasSelection = computed(() => selectedRows.value.length > 0)

const handleSelectionChange = (rows: Bucket[]) => {
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
  loadBuckets()
}

const handleView = () => {
  if (selectedRows.value.length === 1) {
    viewDetail(selectedRows.value[0])
  } else {
    ElMessage.warning('请选择一条记录查看')
  }
}

const handleCreate = () => {
  createForm.value = { name: '', permission: 'private', region: '' }
  createDialogVisible.value = true
}

const submitCreate = async () => {
  try {
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    loadBuckets()
  } catch (e) {
    console.error(e)
    ElMessage.error('创建失败')
  }
}

const handleBatchCommand = async (command: string) => {
  if (!hasSelection.value) return
  if (command === 'delete') {
    try {
      await ElMessageBox.confirm(`确认删除选中的 ${selectedRows.value.length} 个存储桶？`, '警告', { type: 'warning' })
      buckets.value = buckets.value.filter(b => !selectedRows.value.some(r => r.id === b.id))
      ElMessage.success('删除成功')
      loadBuckets()
    } catch (e) {
      console.error(e)
    }
  }
}

const handleTags = () => {
  if (!hasSelection.value) {
    ElMessage.warning('请选择需要设置标签的存储桶')
    return
  }
  ElMessage.info('标签功能开发中')
}

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
    case 'normal':
      return 'success'
    case 'pending':
    case 'creating':
      return 'warning'
    case 'error':
    case 'deleted':
      return 'danger'
    default:
      return 'info'
  }
}

const loadBuckets = async () => {
  loading.value = true
  try {
    // Mock data
    buckets.value = [
      {
        id: 'bucket-1',
        name: 'prod-bucket',
        status: 'Active',
        tags: 'prod',
        permission: '私有',
        platform: '阿里云',
        project: 'Project A',
        region: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00',
        storage_usage: 500
      },
      {
        id: 'bucket-2',
        name: 'dev-bucket',
        status: 'Active',
        tags: 'dev',
        permission: '公共读',
        platform: '阿里云',
        project: 'Project B',
        region: 'cn-shanghai',
        created_at: '2024-01-02 10:00:00',
        storage_usage: 100
      },
      {
        id: 'bucket-3',
        name: 'log-bucket',
        status: 'Creating',
        tags: 'log',
        permission: '私有',
        platform: '阿里云',
        project: 'Project A',
        region: 'cn-beijing',
        created_at: '2024-01-03 10:00:00',
        storage_usage: 0
      }
    ]
  } catch (e) {
    console.error(e)
    buckets.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: Bucket) => {
  selectedBucket.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row: Bucket) => {
  ElMessage.info(`编辑存储桶: ${row.name}`)
}

const handleDelete = async (row: Bucket) => {
  try {
    await ElMessageBox.confirm(`确认删除存储桶 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    buckets.value = buckets.value.filter(b => b.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadBuckets()
})
</script>

<style scoped>
.buckets-page {
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