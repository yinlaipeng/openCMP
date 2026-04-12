<template>
  <div class="buckets-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">存储桶列表</span>
        </div>
      </template>

      <el-table :data="buckets" v-loading="loading">
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

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
const detailDialogVisible = ref(false)
const selectedBucket = ref<Bucket | null>(null)

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
</style>