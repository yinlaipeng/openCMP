<template>
  <div class="vpc-interconnect-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">VPC互联列表</span>
        </div>
      </template>

      <el-table :data="interconnects" v-loading="loading">
        <el-table-column label="名称" width="200">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="vpc_count" label="vpc数量" width="100" />
        <el-table-column prop="domain" label="所属域" width="150" />
        <el-table-column prop="platform" label="平台" width="120" />
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="VPC互联详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedInterconnect?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedInterconnect?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedInterconnect?.status }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ selectedInterconnect?.type }}</el-descriptions-item>
        <el-descriptions-item label="带宽">{{ selectedInterconnect?.bandwidth }} Mbps</el-descriptions-item>
        <el-descriptions-item label="VPC数量">{{ selectedInterconnect?.vpc_count }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ selectedInterconnect?.domain }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedInterconnect?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedInterconnect?.account }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedInterconnect?.created_at }}</el-descriptions-item>
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

interface VPCInterconnectItem {
  id: string
  name: string
  status: string
  type: string
  bandwidth: number
  vpc_count: number
  domain: string
  platform: string
  account: string
  created_at: string
}

const interconnects = ref<VPCInterconnectItem[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedInterconnect = ref<VPCInterconnectItem | null>(null)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'available':
    case 'running':
    case 'active':
      return 'success'
    case 'pending':
      return 'warning'
    case 'failed':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const loadInterconnects = async () => {
  loading.value = true
  // Mock data - API not yet implemented
  interconnects.value = [
    {
      id: 'interconnect-1',
      name: 'VPC互联 1',
      status: 'Active',
      type: '专线',
      bandwidth: 1000,
      vpc_count: 3,
      domain: 'Default Domain',
      platform: '阿里云',
      account: 'Aliyun Account 1',
      created_at: new Date().toISOString()
    },
    {
      id: 'interconnect-2',
      name: 'VPC互联 2',
      status: 'Active',
      type: 'VPN',
      bandwidth: 500,
      vpc_count: 2,
      domain: 'Default Domain',
      platform: '阿里云',
      account: 'Aliyun Account 1',
      created_at: new Date().toISOString()
    }
  ]
  loading.value = false
}

const viewDetail = (row: VPCInterconnectItem) => {
  selectedInterconnect.value = row
  detailDialogVisible.value = true
}

const handleDelete = async (row: VPCInterconnectItem) => {
  try {
    await ElMessageBox.confirm(`确认删除 VPC互联 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    ElMessage.success('删除成功')
    loadInterconnects()
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadInterconnects()
})
</script>

<style scoped>
.vpc-interconnect-page {
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