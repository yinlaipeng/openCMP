<template>
  <div class="route-tables-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">路由表列表</span>
        </div>
      </template>

      <el-table :data="routeTables" v-loading="loading">
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
        <el-table-column prop="vpc_id" label="vpc" width="180" />
        <el-table-column prop="platform" label="平台" width="120" />
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column prop="domain" label="所属域" width="150" />
        <el-table-column prop="region_id" label="区域" width="150" />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button size="small" type="success" @click="syncStatus(row)">同步状态</el-button>
            <el-button size="small" type="primary" @click="manageRoutes(row)">管理路由</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="路由表详情" width="800px">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="详情" name="detail">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">{{ selectedRouteTable?.id }}</el-descriptions-item>
            <el-descriptions-item label="名称">{{ selectedRouteTable?.name }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ selectedRouteTable?.status }}</el-descriptions-item>
            <el-descriptions-item label="VPC">{{ selectedRouteTable?.vpc_id }}</el-descriptions-item>
            <el-descriptions-item label="平台">{{ selectedRouteTable?.platform }}</el-descriptions-item>
            <el-descriptions-item label="云账号">{{ selectedRouteTable?.account }}</el-descriptions-item>
            <el-descriptions-item label="所属域">{{ selectedRouteTable?.domain }}</el-descriptions-item>
            <el-descriptions-item label="区域">{{ selectedRouteTable?.region_id }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ selectedRouteTable?.created_at }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="路由" name="routes">
          <el-table :data="routes">
            <el-table-column prop="destination_cidr" label="目的网段" width="200" />
            <el-table-column prop="target" label="下一跳" width="200" />
            <el-table-column prop="status" label="状态" width="120">
              <template #default="{ row }">
                <el-tag :type="row.status === 'Active' ? 'success' : 'info'">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="ip子网" name="subnets">
          <el-table :data="subnets">
            <el-table-column prop="id" label="子网ID" width="200" />
            <el-table-column prop="name" label="子网名称" width="150" />
            <el-table-column prop="status" label="状态" width="100" />
            <el-table-column prop="cidr" label="CIDR" width="180" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="操作日志" name="logs">
          <el-table :data="logs">
            <el-table-column prop="operation" label="操作" width="150" />
            <el-table-column prop="operator" label="操作员" width="150" />
            <el-table-column prop="result" label="结果" width="100" />
            <el-table-column prop="timestamp" label="时间" width="180" />
          </el-table>
        </el-tab-pane>
      </el-tabs>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

interface RouteTable {
  id: string
  name: string
  status: string
  vpc_id: string
  platform?: string
  account?: string
  domain?: string
  region_id?: string
  created_at?: string
}

const routeTables = ref<RouteTable[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedRouteTable = ref<RouteTable | null>(null)
const activeTab = ref('detail')

// Detail modal data
const routes = ref<any[]>([])
const subnets = ref<any[]>([])
const logs = ref<any[]>([])

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
    case 'available':
      return 'success'
    case 'pending':
      return 'warning'
    case 'failed':
      return 'danger'
    default:
      return 'info'
  }
}

const loadRouteTables = async () => {
  loading.value = true
  try {
    // Mock data - API requires valid cloud account credentials
    routeTables.value = [
      {
        id: 'rt-1',
        name: '路由表 1',
        status: 'Active',
        vpc_id: 'vpc-1',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        domain: 'Default Domain',
        region_id: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'rt-2',
        name: '路由表 2',
        status: 'Active',
        vpc_id: 'vpc-2',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        domain: 'Default Domain',
        region_id: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'rt-3',
        name: '路由表 3',
        status: 'Pending',
        vpc_id: 'vpc-3',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        domain: 'Default Domain',
        region_id: 'cn-shanghai',
        created_at: '2024-01-01 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    routeTables.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: RouteTable) => {
  selectedRouteTable.value = row
  detailDialogVisible.value = true
  activeTab.value = 'detail'
  loadRouteDetails(row.id)
}

const manageRoutes = (row: RouteTable) => {
  selectedRouteTable.value = row
  detailDialogVisible.value = true
  activeTab.value = 'routes'
  loadRouteDetails(row.id)
}

const loadRouteDetails = async (id: string) => {
  // Load routes (mock data)
  routes.value = [
    { destination_cidr: '10.0.0.0/16', target: 'local', status: 'Active' },
    { destination_cidr: '0.0.0.0/0', target: 'igw-12345', status: 'Active' }
  ]

  // Mock data for subnets (API requires valid cloud account)
  subnets.value = [
    { id: 'subnet-1', name: 'Subnet 1', status: 'Available', cidr: '10.0.1.0/24' },
    { id: 'subnet-2', name: 'Subnet 2', status: 'Available', cidr: '10.0.2.0/24' }
  ]

  // Load logs (mock)
  logs.value = [
    { operation: '创建', operator: 'admin', result: '成功', timestamp: '2024-01-01 10:00:00' }
  ]
}

const syncStatus = (row: RouteTable) => {
  ElMessage.info(`正在同步路由表 ${row.name} 的状态`)
}

onMounted(() => {
  loadRouteTables()
})
</script>

<style scoped>
.route-tables-page {
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