<template>
  <div class="global-vpc-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">全局VPC列表</span>
        </div>
      </template>

      <el-table :data="globalVpcs" v-loading="loading">
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
        <el-table-column prop="platform" label="平台" width="120" />
        <el-table-column prop="vpc_count" label="vpc数量" width="100" />
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column prop="domain" label="所属域" width="150" />
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="全局VPC详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedGlobalVpc?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedGlobalVpc?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedGlobalVpc?.status }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedGlobalVpc?.platform }}</el-descriptions-item>
        <el-descriptions-item label="VPC数量">{{ selectedGlobalVpc?.vpc_count }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ selectedGlobalVpc?.domain }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedGlobalVpc?.account }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedGlobalVpc?.created_at }}</el-descriptions-item>
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

interface GlobalVPC {
  id: string
  name: string
  status: string
  platform: string
  vpc_count: number
  account: string
  domain: string
  created_at: string
}

const globalVpcs = ref<GlobalVPC[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedGlobalVpc = ref<GlobalVPC | null>(null)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
      return 'success'
    case 'pending':
      return 'warning'
    case 'inactive':
    case 'failed':
      return 'danger'
    default:
      return 'info'
  }
}

const loadGlobalVpcs = async () => {
  loading.value = true
  try {
    // Mock data - in real implementation, call API
    globalVpcs.value = [
      {
        id: 'global-vpc-1',
        name: 'Global VPC 1',
        status: 'Active',
        platform: '阿里云',
        vpc_count: 5,
        account: 'Aliyun Account 1',
        domain: 'Default Domain',
        created_at: new Date().toISOString()
      },
      {
        id: 'global-vpc-2',
        name: 'Global VPC 2',
        status: 'Active',
        platform: 'AWS',
        vpc_count: 3,
        account: 'AWS Account 1',
        domain: 'Default Domain',
        created_at: new Date().toISOString()
      }
    ]
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: GlobalVPC) => {
  selectedGlobalVpc.value = row
  detailDialogVisible.value = true
}

const handleDelete = async (row: GlobalVPC) => {
  try {
    await ElMessageBox.confirm(`确认删除全局VPC ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    ElMessage.success('删除成功')
    loadGlobalVpcs()
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadGlobalVpcs()
})
</script>

<style scoped>
.global-vpc-page {
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