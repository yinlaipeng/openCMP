<template>
  <div class="cdn-domains-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">CDN域名列表</span>
        </div>
      </template>

      <el-table :data="cdnDomains" v-loading="loading">
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
        <el-table-column prop="accelerate_region" label="加速区域" width="150" />
        <el-table-column prop="cname" label="CNAME" width="200" />
        <el-table-column prop="accelerate_type" label="加速类别" width="120" />
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column prop="project" label="项目" width="120" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="refreshPreheat(row)">刷新预热</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="CDN域名详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedCDNDomain?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedCDNDomain?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedCDNDomain?.status }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{ selectedCDNDomain?.tags }}</el-descriptions-item>
        <el-descriptions-item label="加速区域">{{ selectedCDNDomain?.accelerate_region }}</el-descriptions-item>
        <el-descriptions-item label="CNAME">{{ selectedCDNDomain?.cname }}</el-descriptions-item>
        <el-descriptions-item label="加速类别">{{ selectedCDNDomain?.accelerate_type }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedCDNDomain?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedCDNDomain?.account }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedCDNDomain?.project }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedCDNDomain?.created_at }}</el-descriptions-item>
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

interface CDNDomain {
  id: string
  name: string
  status: string
  tags: string
  accelerate_region: string
  cname: string
  accelerate_type: string
  platform: string
  account: string
  project: string
  created_at: string
}

const cdnDomains = ref<CDNDomain[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedCDNDomain = ref<CDNDomain | null>(null)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'online':
    case 'active':
      return 'success'
    case 'pending':
    case 'configuring':
      return 'warning'
    case 'offline':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const loadCDNDomains = async () => {
  loading.value = true
  try {
    // Mock data
    cdnDomains.value = [
      {
        id: 'cdn-1',
        name: 'cdn.example.com',
        status: 'Online',
        tags: 'prod',
        accelerate_region: '中国大陆',
        cname: 'cdn.example.com.w.cdngslb.com',
        accelerate_type: '全站加速',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project A',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'cdn-2',
        name: 'static.example.com',
        status: 'Configuring',
        tags: 'static',
        accelerate_region: '全球',
        cname: 'static.example.com.w.cdngslb.com',
        accelerate_type: '静态加速',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project B',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'cdn-3',
        name: 'video.example.com',
        status: 'Offline',
        tags: 'video',
        accelerate_region: '亚太地区',
        cname: 'video.example.com.w.cdngslb.com',
        accelerate_type: '视频点播',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project A',
        created_at: '2024-01-03 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    cdnDomains.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: CDNDomain) => {
  selectedCDNDomain.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row: CDNDomain) => {
  ElMessage.info(`编辑CDN域名: ${row.name}`)
}

const refreshPreheat = (row: CDNDomain) => {
  ElMessage.info(`刷新预热CDN域名: ${row.name}`)
}

const handleDelete = async (row: CDNDomain) => {
  try {
    await ElMessageBox.confirm(`确认删除CDN域名 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    cdnDomains.value = cdnDomains.value.filter(c => c.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadCDNDomains()
})
</script>

<style scoped>
.cdn-domains-page {
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