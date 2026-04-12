<template>
  <div class="table-storage-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">表格存储列表</span>
        </div>
      </template>

      <el-table :data="tableStorage" v-loading="loading">
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
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column prop="project" label="项目" width="120" />
        <el-table-column prop="region" label="区域" width="120" />
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="表格存储详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedTableStorage?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedTableStorage?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedTableStorage?.status }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedTableStorage?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedTableStorage?.account }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedTableStorage?.project }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedTableStorage?.region }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedTableStorage?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

interface TableStorage {
  id: string
  name: string
  status: string
  platform: string
  account: string
  project: string
  region: string
  created_at: string
}

const tableStorage = ref<TableStorage[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedTableStorage = ref<TableStorage | null>(null)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
    case 'normal':
      return 'success'
    case 'pending':
    case 'creating':
      return 'warning'
    case 'error':
    case 'disabled':
      return 'danger'
    default:
      return 'info'
  }
}

const loadTableStorage = async () => {
  loading.value = true
  try {
    // Mock data
    tableStorage.value = [
      {
        id: 'table-1',
        name: 'user-table',
        status: 'Active',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project A',
        region: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'table-2',
        name: 'order-table',
        status: 'Active',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project B',
        region: 'cn-shanghai',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'table-3',
        name: 'log-table',
        status: 'Disabled',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project A',
        region: 'cn-beijing',
        created_at: '2024-01-03 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    tableStorage.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: TableStorage) => {
  selectedTableStorage.value = row
  detailDialogVisible.value = true
}

onMounted(() => {
  loadTableStorage()
})
</script>

<style scoped>
.table-storage-page {
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