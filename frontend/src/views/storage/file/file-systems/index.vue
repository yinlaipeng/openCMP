<template>
  <div class="file-systems-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">文件系统列表</span>
        </div>
      </template>

      <el-table :data="fileSystems" v-loading="loading">
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

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
const detailDialogVisible = ref(false)
const selectedFileSystem = ref<FileSystem | null>(null)

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
  } catch (e) {
    console.error(e)
    fileSystems.value = []
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
</style>