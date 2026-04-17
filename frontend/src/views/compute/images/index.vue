<template>
  <div class="images-container">
    <div class="page-header">
      <h2>镜像管理</h2>
      <el-button type="primary" @click="handleSync">同步镜像</el-button>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="loadImages">
        <el-form-item label="云账号">
          <CloudAccountSelector
            v-model:value="filters.account_id"
            placeholder="选择云账号"
            @change="handleAccountChange"
          />
        </el-form-item>
        <el-form-item label="操作系统">
          <el-select v-model="filters.os_name" placeholder="选择操作系统" clearable>
            <el-option label="CentOS" value="CentOS" />
            <el-option label="Ubuntu" value="Ubuntu" />
            <el-option label="Windows" value="Windows" />
            <el-option label="Debian" value="Debian" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="filters.name" placeholder="镜像名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadImages">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table
      :data="images"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
    >
      <el-table-column prop="id" label="ID" min-width="150" />
      <el-table-column prop="name" label="名称" min-width="200" />
      <el-table-column label="平台/云账号" width="150">
        <template #default="{ row }">
          <div class="platform-cell">
            <el-tag size="small" :type="getPlatformType(row.platform)" effect="plain">
              {{ getPlatformLabel(row.platform) }}
            </el-tag>
            <span class="account-name">{{ row.account_name || '-' }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="os_name" label="操作系统" width="120" />
      <el-table-column prop="os_version" label="版本" width="100" />
      <el-table-column prop="architecture" label="架构" width="100" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="大小" width="100">
        <template #default="{ row }">
          {{ formatSize(row.size) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="viewDetail(row)">查看</el-button>
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

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="镜像详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedImage?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedImage?.name }}</el-descriptions-item>
        <el-descriptions-item label="操作系统">{{ selectedImage?.os_name }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ selectedImage?.os_version }}</el-descriptions-item>
        <el-descriptions-item label="架构">{{ selectedImage?.architecture }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedImage?.status }}</el-descriptions-item>
        <el-descriptions-item label="大小">{{ formatSize(selectedImage?.size) }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedImage?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import { getImages } from '@/api/compute'
import type { Image } from '@/types'

const images = ref<Image[]>([])
const loading = ref(false)
const detailVisible = ref(false)
const selectedImage = ref<Image | null>(null)

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 筛选条件
const filters = reactive({
  account_id: null as number | null,
  os_name: '',
  name: ''
})

// 平台类型映射
const platformLabels: Record<string, string> = {
  alibaba: '阿里云',
  tencent: '腾讯云',
  aws: 'AWS',
  azure: 'Azure'
}

const platformTypes: Record<string, 'primary' | 'warning' | 'success' | 'info'> = {
  alibaba: 'primary',
  tencent: 'warning',
  aws: 'success',
  azure: 'info'
}

const getPlatformLabel = (platform: string): string => {
  return platformLabels[platform] || platform || '未知'
}

const getPlatformType = (platform: string): 'primary' | 'warning' | 'success' | 'info' => {
  return platformTypes[platform] || 'info'
}

const getStatusType = (status: string) => {
  if (status === 'Available' || status === 'available') return 'success'
  if (status === 'Creating' || status === 'creating') return 'warning'
  if (status === 'Error' || status === 'error') return 'danger'
  return 'info'
}

const formatSize = (size?: number): string => {
  if (!size) return '-'
  return (size / 1024 / 1024 / 1024).toFixed(2) + ' GB'
}

const handleAccountChange = (accountId: number | null) => {
  filters.account_id = accountId
}

const loadImages = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      size: pagination.pageSize
    }
    if (filters.account_id) params.account_id = filters.account_id
    if (filters.os_name) params.os_name = filters.os_name
    if (filters.name) params.name = filters.name

    const res = await getImages(params)
    images.value = res.items || res
    pagination.total = res.total || images.value.length
  } catch (e) {
    console.error(e)
    ElMessage.error('加载镜像列表失败')
  } finally {
    loading.value = false
  }
}

const handleSync = () => {
  ElMessage.info('同步镜像功能开发中')
}

const viewDetail = (row: Image) => {
  selectedImage.value = row
  detailVisible.value = true
}

const resetFilters = () => {
  filters.account_id = null
  filters.os_name = ''
  filters.name = ''
  pagination.page = 1
  loadImages()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadImages()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadImages()
}

onMounted(() => {
  loadImages()
})
</script>

<style scoped>
.images-container {
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

.pagination {
  margin-top: 20px;
  text-align: right;
}

.platform-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.account-name {
  font-size: 12px;
  color: var(--text-color-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>