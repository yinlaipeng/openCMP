<template>
  <div class="vpc-peering-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">VPC对等连接列表</span>
        </div>
      </template>

      <el-table :data="peerings" v-loading="loading">
        <el-table-column label="名称" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="local_vpc_id" label="本端vpc" width="180" />
        <el-table-column prop="peer_vpc_id" label="对端vpc" width="180" />
        <el-table-column prop="peer_account" label="对端账号" width="180" />
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
    <el-dialog v-model="detailDialogVisible" title="VPC对等连接详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedPeering?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedPeering?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedPeering?.status }}</el-descriptions-item>
        <el-descriptions-item label="本端VPC">{{ selectedPeering?.local_vpc_id }}</el-descriptions-item>
        <el-descriptions-item label="对端VPC">{{ selectedPeering?.peer_vpc_id }}</el-descriptions-item>
        <el-descriptions-item label="对端账号">{{ selectedPeering?.peer_account }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ selectedPeering?.domain }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedPeering?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedPeering?.account }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedPeering?.created_at }}</el-descriptions-item>
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

interface VPCPeeringItem {
  id: string
  name: string
  status: string
  local_vpc_id: string
  peer_vpc_id: string
  peer_account: string
  domain: string
  platform: string
  account: string
  created_at: string
}

const peerings = ref<VPCPeeringItem[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedPeering = ref<VPCPeeringItem | null>(null)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'available':
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

const loadPeerings = async () => {
  loading.value = true
  // Mock data - API not yet implemented
  peerings.value = [
    {
      id: 'peering-1',
      name: 'VPC对等连接 1',
      status: 'Active',
      local_vpc_id: 'vpc-local-1',
      peer_vpc_id: 'vpc-peer-1',
      peer_account: 'account-2',
      domain: 'Default Domain',
      platform: '阿里云',
      account: 'Aliyun Account 1',
      created_at: new Date().toISOString()
    },
    {
      id: 'peering-2',
      name: 'VPC对等连接 2',
      status: 'Active',
      local_vpc_id: 'vpc-local-2',
      peer_vpc_id: 'vpc-peer-2',
      peer_account: 'account-3',
      domain: 'Default Domain',
      platform: '阿里云',
      account: 'Aliyun Account 1',
      created_at: new Date().toISOString()
    }
  ]
  loading.value = false
}

const viewDetail = (row: VPCPeeringItem) => {
  selectedPeering.value = row
  detailDialogVisible.value = true
}

const handleDelete = async (row: VPCPeeringItem) => {
  try {
    await ElMessageBox.confirm(`确认删除 VPC对等连接 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    ElMessage.success('删除成功')
    loadPeerings()
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadPeerings()
})
</script>

<style scoped>
.vpc-peering-page {
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