<template>
  <div class="lb-certificates-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">证书列表</span>
        </div>
      </template>

      <el-table :data="certificates" v-loading="loading">
        <el-table-column label="证书名称" width="180">
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
        <el-table-column prop="domain" label="证书域名" width="180" />
        <el-table-column prop="expire_time" label="过期时间" width="180" />
        <el-table-column prop="extension_domains" label="关联扩展域名" width="200" />
        <el-table-column prop="listeners_count" label="关联监听(数量)" width="150">
          <template #default="{ row }">
            {{ row.listeners_count }} 个
          </template>
        </el-table-column>
        <el-table-column prop="share_scope" label="共享范围" width="100" />
        <el-table-column prop="project" label="项目" width="120" />
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-dropdown>
              <el-button size="small">
                更多 <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="changeProject(row)">更改项目</el-dropdown-item>
                  <el-dropdown-item @click="setShare(row)">设置共享</el-dropdown-item>
                  <el-dropdown-item divided @click="handleDelete(row)">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="证书详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedCertificate?.id }}</el-descriptions-item>
        <el-descriptions-item label="证书名称">{{ selectedCertificate?.name }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{ selectedCertificate?.tags }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedCertificate?.status }}</el-descriptions-item>
        <el-descriptions-item label="证书域名">{{ selectedCertificate?.domain }}</el-descriptions-item>
        <el-descriptions-item label="过期时间">{{ selectedCertificate?.expire_time }}</el-descriptions-item>
        <el-descriptions-item label="关联扩展域名">{{ selectedCertificate?.extension_domains }}</el-descriptions-item>
        <el-descriptions-item label="关联监听">{{ selectedCertificate?.listeners_count }} 个</el-descriptions-item>
        <el-descriptions-item label="共享范围">{{ selectedCertificate?.share_scope }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedCertificate?.project }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedCertificate?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedCertificate?.account }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedCertificate?.region }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedCertificate?.created_at }}</el-descriptions-item>
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
import { ArrowDown } from '@element-plus/icons-vue'

interface Certificate {
  id: string
  name: string
  tags: string
  status: string
  domain: string
  expire_time: string
  extension_domains: string
  listeners_count: number
  share_scope: string
  project: string
  platform: string
  account: string
  region: string
  created_at: string
}

const certificates = ref<Certificate[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedCertificate = ref<Certificate | null>(null)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'normal':
    case 'active':
      return 'success'
    case 'expiring':
      return 'warning'
    case 'expired':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const loadCertificates = async () => {
  loading.value = true
  try {
    // Mock data
    certificates.value = [
      {
        id: 'cert-1',
        name: '证书 1',
        tags: 'prod',
        status: 'Normal',
        domain: 'example.com',
        expire_time: '2025-12-31',
        extension_domains: '*.example.com',
        listeners_count: 2,
        share_scope: '私有',
        project: 'Project A',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        region: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'cert-2',
        name: '证书 2',
        tags: 'api',
        status: 'Expiring',
        domain: 'api.example.com',
        expire_time: '2024-06-30',
        extension_domains: '-',
        listeners_count: 1,
        share_scope: '共享',
        project: 'Project B',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        region: 'cn-shanghai',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'cert-3',
        name: '证书 3',
        tags: 'test',
        status: 'Expired',
        domain: 'test.example.com',
        expire_time: '2024-01-01',
        extension_domains: '-',
        listeners_count: 0,
        share_scope: '私有',
        project: 'Project A',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        region: 'cn-beijing',
        created_at: '2023-01-01 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    certificates.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: Certificate) => {
  selectedCertificate.value = row
  detailDialogVisible.value = true
}

const changeProject = (row: Certificate) => {
  ElMessage.info(`更改证书项目: ${row.name}`)
}

const setShare = (row: Certificate) => {
  ElMessage.info(`设置证书共享: ${row.name}`)
}

const handleDelete = async (row: Certificate) => {
  try {
    await ElMessageBox.confirm(`确认删除证书 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    certificates.value = certificates.value.filter(c => c.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadCertificates()
})
</script>

<style scoped>
.lb-certificates-page {
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