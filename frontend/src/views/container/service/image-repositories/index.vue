<template>
  <div class="image-repositories-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">镜像仓库列表</span>
        </div>
      </template>

      <el-table :data="repositories" v-loading="loading">
        <el-table-column label="名称" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="label" label="标签" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="spec" label="实例规格" width="120" />
        <el-table-column prop="domain" label="实例域名" width="200" />
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column prop="project" label="项目" width="120" />
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
    <el-dialog v-model="detailDialogVisible" title="镜像仓库详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedRepo?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedRepo?.name }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{ selectedRepo?.label }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedRepo?.status }}</el-descriptions-item>
        <el-descriptions-item label="实例规格">{{ selectedRepo?.spec }}</el-descriptions-item>
        <el-descriptions-item label="实例域名">{{ selectedRepo?.domain }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedRepo?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedRepo?.account }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedRepo?.project }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedRepo?.region }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedRepo?.created_at }}</el-descriptions-item>
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

interface ImageRepository {
  id: string
  name: string
  label: string
  status: string
  spec: string
  domain: string
  platform: string
  account: string
  project: string
  region: string
  created_at: string
}

const repositories = ref<ImageRepository[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedRepo = ref<ImageRepository | null>(null)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'running':
    case 'active':
      return 'success'
    case 'creating':
    case 'pending':
      return 'warning'
    case 'stopped':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const loadRepositories = async () => {
  loading.value = true
  try {
    // Mock data
    repositories.value = [
      {
        id: 'repo-1',
        name: 'prod-registry',
        label: 'production',
        status: 'Running',
        spec: '企业版',
        domain: 'registry.cn-hangzhou.aliyuncs.com',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project A',
        region: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'repo-2',
        name: 'dev-registry',
        label: 'development',
        status: 'Running',
        spec: '个人版',
        domain: 'registry.cn-shanghai.aliyuncs.com',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project B',
        region: 'cn-shanghai',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'repo-3',
        name: 'test-registry',
        label: 'test',
        status: 'Creating',
        spec: '企业版',
        domain: '-',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project A',
        region: 'cn-beijing',
        created_at: '2024-01-03 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    repositories.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: ImageRepository) => {
  selectedRepo.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row: ImageRepository) => {
  ElMessage.info(`编辑镜像仓库: ${row.name}`)
}

const handleDelete = async (row: ImageRepository) => {
  try {
    await ElMessageBox.confirm(`确认删除镜像仓库 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    repositories.value = repositories.value.filter(r => r.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadRepositories()
})
</script>

<style scoped>
.image-repositories-page {
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