<template>
  <div class="ipv6-gateways-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">IPv6网关列表</span>
        </div>
      </template>

      <el-table :data="ipv6Gateways" v-loading="loading">
        <el-table-column label="名称" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="vpc" label="vpc" width="180" />
        <el-table-column prop="spec" label="规格" width="120" />
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="account" label="云账号" width="150" />
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
    <el-dialog v-model="detailDialogVisible" title="IPv6网关详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedIPv6Gateway?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedIPv6Gateway?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedIPv6Gateway?.status }}</el-descriptions-item>
        <el-descriptions-item label="VPC">{{ selectedIPv6Gateway?.vpc }}</el-descriptions-item>
        <el-descriptions-item label="规格">{{ selectedIPv6Gateway?.spec }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedIPv6Gateway?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedIPv6Gateway?.account }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedIPv6Gateway?.project }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedIPv6Gateway?.region }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedIPv6Gateway?.created_at }}</el-descriptions-item>
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

interface IPv6Gateway {
  id: string
  name: string
  status: string
  vpc: string
  spec: string
  platform: string
  account: string
  project: string
  region: string
  created_at: string
}

const ipv6Gateways = ref<IPv6Gateway[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedIPv6Gateway = ref<IPv6Gateway | null>(null)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'available':
    case 'active':
      return 'success'
    case 'pending':
    case 'creating':
      return 'warning'
    case 'error':
    case 'failed':
      return 'danger'
    default:
      return 'info'
  }
}

const loadIPv6Gateways = async () => {
  loading.value = true
  try {
    // Mock data
    ipv6Gateways.value = [
      {
        id: 'ipv6gw-1',
        name: 'IPv6网关 1',
        status: 'Available',
        vpc: 'vpc-1',
        spec: '中小企业版',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project A',
        region: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'ipv6gw-2',
        name: 'IPv6网关 2',
        status: 'Available',
        vpc: 'vpc-2',
        spec: '大型企业版',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project B',
        region: 'cn-shanghai',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'ipv6gw-3',
        name: 'IPv6网关 3',
        status: 'Creating',
        vpc: 'vpc-3',
        spec: '中小企业版',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project A',
        region: 'cn-beijing',
        created_at: '2024-01-03 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    ipv6Gateways.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: IPv6Gateway) => {
  selectedIPv6Gateway.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row: IPv6Gateway) => {
  ElMessage.info(`编辑IPv6网关: ${row.name}`)
}

const handleDelete = async (row: IPv6Gateway) => {
  try {
    await ElMessageBox.confirm(`确认删除IPv6网关 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    ipv6Gateways.value = ipv6Gateways.value.filter(g => g.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadIPv6Gateways()
})
</script>

<style scoped>
.ipv6-gateways-page {
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