<template>
  <div class="app-services-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">应用程序服务列表</span>
        </div>
      </template>

      <el-table :data="appServices" v-loading="loading">
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
        <el-table-column prop="tech_stack" label="技术栈" width="120" />
        <el-table-column prop="os" label="操作系统" width="120" />
        <el-table-column prop="endpoint" label="入站点地址" width="180" />
        <el-table-column prop="domain" label="域" width="150" />
        <el-table-column prop="plan" label="应用服务计划" width="150" />
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column prop="project" label="项目" width="120" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="应用程序服务详情" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedAppService?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedAppService?.name }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{ selectedAppService?.tags }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedAppService?.status }}</el-descriptions-item>
        <el-descriptions-item label="技术栈">{{ selectedAppService?.tech_stack }}</el-descriptions-item>
        <el-descriptions-item label="操作系统">{{ selectedAppService?.os }}</el-descriptions-item>
        <el-descriptions-item label="入站点地址">{{ selectedAppService?.endpoint }}</el-descriptions-item>
        <el-descriptions-item label="域">{{ selectedAppService?.domain }}</el-descriptions-item>
        <el-descriptions-item label="应用服务计划">{{ selectedAppService?.plan }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedAppService?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedAppService?.account }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedAppService?.region }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedAppService?.project }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedAppService?.created_at }}</el-descriptions-item>
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

interface AppService {
  id: string
  name: string
  tags: string
  status: string
  tech_stack: string
  os: string
  endpoint: string
  domain: string
  plan: string
  platform: string
  account: string
  region: string
  project: string
  created_at: string
}

const appServices = ref<AppService[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedAppService = ref<AppService | null>(null)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'running':
    case 'active':
      return 'success'
    case 'pending':
    case 'deploying':
      return 'warning'
    case 'stopped':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const loadAppServices = async () => {
  loading.value = true
  try {
    // Mock data
    appServices.value = [
      {
        id: 'app-1',
        name: '应用程序服务 1',
        tags: 'web',
        status: 'Running',
        tech_stack: 'Java',
        os: 'Linux',
        endpoint: 'https://app1.example.com',
        domain: 'Default Domain',
        plan: '标准版',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        region: 'cn-hangzhou',
        project: 'Project A',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'app-2',
        name: '应用程序服务 2',
        tags: 'api',
        status: 'Deploying',
        tech_stack: 'Python',
        os: 'Linux',
        endpoint: 'https://app2.example.com',
        domain: 'Default Domain',
        plan: '高级版',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        region: 'cn-shanghai',
        project: 'Project B',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'app-3',
        name: '应用程序服务 3',
        tags: 'backend',
        status: 'Stopped',
        tech_stack: 'Node.js',
        os: 'Windows',
        endpoint: 'https://app3.example.com',
        domain: 'Domain A',
        plan: '基础版',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        region: 'cn-beijing',
        project: 'Project A',
        created_at: '2024-01-03 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    appServices.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: AppService) => {
  selectedAppService.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row: AppService) => {
  ElMessage.info(`编辑应用程序服务: ${row.name}`)
}

const handleDelete = async (row: AppService) => {
  try {
    await ElMessageBox.confirm(`确认删除应用程序服务 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    appServices.value = appServices.value.filter(a => a.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadAppServices()
})
</script>

<style scoped>
.app-services-page {
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