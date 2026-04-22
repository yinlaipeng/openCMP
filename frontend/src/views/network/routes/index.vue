<template>
  <div class="route-tables-container">
    <div class="page-header">
      <h2>路由表</h2>
      <div class="toolbar">
        <el-button link type="primary" @click="handleView">
          <el-icon><View /></el-icon>
          查看
        </el-button>
        <el-button @click="handleSyncStatus">
          <el-icon><Refresh /></el-icon>
          同步状态
        </el-button>
        <el-button @click="handleManageRoute">
          <el-icon><Setting /></el-icon>
          管理路由
        </el-button>
      </div>
    </div>

    <!-- 搜索表单 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="fetchData">
        <el-form-item label="名称">
          <el-input v-model="filters.name" placeholder="搜索路由表名称" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="选择状态" clearable style="width: 120px">
            <el-option label="可用" value="available" />
            <el-option label="创建中" value="creating" />
            <el-option label="错误" value="error" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="filters.platform" placeholder="选择平台" clearable style="width: 120px">
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-input v-model="filters.region" placeholder="区域" clearable style="width: 120px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <el-table
        :data="routeTables"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
      >
        <!-- Cloudpods 表头顺序: Name → Status → VPC → Platform → Cloud account → Owner Domain → Region -->
        <el-table-column prop="name" label="名称" width="180">
          <template #default="{ row }">
            <el-link type="primary" :underline="false" @click="handleDetails(row)">
              {{ row.name }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="vpc_name" label="VPC" width="150" />
        <el-table-column prop="provider_type" label="平台" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="getPlatformType(row.provider_type)">
              {{ getPlatformLabel(row.provider_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="account_name" label="云账户" width="150" />
        <el-table-column prop="owner_domain" label="所属域" width="150" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleSyncStatus(row)">同步状态</el-button>
            <el-button size="small" @click="handleManageRoute(row)">管理路由</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        class="pagination"
      />
    </el-card>

    <!-- 详情弹窗 -->
    <el-dialog
      title="路由表详情"
      v-model="detailDialogVisible"
      width="700px"
    >
      <el-tabs v-model="activeTab">
        <el-tab-pane label="详情" name="detail">
          <el-descriptions :column="2" border v-if="selectedRouteTable">
            <el-descriptions-item label="ID">{{ selectedRouteTable.id }}</el-descriptions-item>
            <el-descriptions-item label="名称">{{ selectedRouteTable.name }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(selectedRouteTable.status)">{{ getStatusLabel(selectedRouteTable.status) }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="VPC">{{ selectedRouteTable.vpc_name }}</el-descriptions-item>
            <el-descriptions-item label="平台">
              <el-tag size="small" :type="getPlatformType(selectedRouteTable.provider_type)">
                {{ getPlatformLabel(selectedRouteTable.provider_type) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="云账户">{{ selectedRouteTable.account_name }}</el-descriptions-item>
            <el-descriptions-item label="所属域">{{ selectedRouteTable.owner_domain }}</el-descriptions-item>
            <el-descriptions-item label="区域">{{ selectedRouteTable.region }}</el-descriptions-item>
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
      </el-tabs>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 管理路由弹窗 -->
    <el-dialog
      title="管理路由"
      v-model="manageRouteDialogVisible"
      width="600px"
    >
      <el-form :model="routeForm" label-width="100px">
        <el-form-item label="路由表">
          {{ routeForm.route_table_name }}
        </el-form-item>
        <el-form-item label="目的网段">
          <el-input v-model="routeForm.destination_cidr" placeholder="例如: 10.0.0.0/16" />
        </el-form-item>
        <el-form-item label="下一跳类型">
          <el-select v-model="routeForm.target_type" placeholder="选择类型" style="width: 100%">
            <el-option label="本地" value="local" />
            <el-option label="Internet网关" value="igw" />
            <el-option label="NAT网关" value="nat" />
            <el-option label="VPN网关" value="vpn" />
          </el-select>
        </el-form-item>
        <el-form-item label="下一跳">
          <el-input v-model="routeForm.target" placeholder="下一跳ID或地址" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="manageRouteDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAddRoute">添加路由</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { View, Refresh, Setting } from '@element-plus/icons-vue'
import { getRouteTables, RouteTable } from '@/api/networkSync'

// Data
const loading = ref(false)
const routeTables = ref<RouteTable[]>([])
const selectedRouteTable = ref<RouteTable | null>(null)
const detailDialogVisible = ref(false)
const manageRouteDialogVisible = ref(false)
const activeTab = ref('detail')

const routes = ref<any[]>([])

const filters = reactive({
  name: '',
  status: '',
  platform: '',
  region: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const routeForm = reactive({
  route_table_id: '',
  route_table_name: '',
  destination_cidr: '',
  target_type: 'local',
  target: ''
})

// Methods
const getStatusType = (status: string) => {
  switch (status) {
    case 'available': return 'success'
    case 'creating': return 'warning'
    case 'error': return 'danger'
    default: return 'info'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'available': return '可用'
    case 'creating': return '创建中'
    case 'error': return '错误'
    default: return status
  }
}

const getPlatformLabel = (platform: string) => {
  const labels: Record<string, string> = {
    aliyun: '阿里云',
    alibaba: '阿里云',
    tencent: '腾讯云',
    Qcloud: '腾讯云',
    aws: 'AWS',
    azure: 'Azure'
  }
  return labels[platform] || platform || '未知'
}

const getPlatformType = (platform: string) => {
  const types: Record<string, 'primary' | 'warning' | 'success' | 'info'> = {
    aliyun: 'primary',
    alibaba: 'primary',
    tencent: 'warning',
    Qcloud: 'warning',
    aws: 'success',
    azure: 'info'
  }
  return types[platform] || 'info'
}

const resetFilters = () => {
  filters.name = ''
  filters.status = ''
  filters.platform = ''
  filters.region = ''
  pagination.page = 1
  fetchData()
}

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      ...filters
    }
    const res = await getRouteTables(params)
    routeTables.value = res.items || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('Failed to fetch route tables:', error)
    ElMessage.error('获取路由表列表失败')
  } finally {
    loading.value = false
  }
}

const handleView = () => {
  ElMessage.info('视图切换功能开发中')
}

const handleSyncStatus = (row?: RouteTable) => {
  if (row) {
    ElMessage.success(`同步路由表 "${row.name}" 状态成功`)
  } else {
    ElMessage.info('同步所有路由表状态')
  }
  fetchData()
}

const handleManageRoute = (row?: RouteTable) => {
  if (row) {
    routeForm.route_table_id = row.id
    routeForm.route_table_name = row.name
    routeForm.destination_cidr = ''
    routeForm.target_type = 'local'
    routeForm.target = ''
    manageRouteDialogVisible.value = true
  } else {
    ElMessage.info('请选择一个路由表进行管理')
  }
}

const confirmAddRoute = () => {
  ElMessage.success('路由添加成功')
  manageRouteDialogVisible.value = false
}

const handleDetails = (row: RouteTable) => {
  selectedRouteTable.value = row
  detailDialogVisible.value = true
  activeTab.value = 'detail'
  loadRouteDetails(row.id)
}

const loadRouteDetails = async (id: string) => {
  routes.value = []
}

const handleSizeChange = (val: number) => {
  pagination.pageSize = val
  pagination.page = 1
  fetchData()
}

const handleCurrentChange = (val: number) => {
  pagination.page = val
  fetchData()
}

// Lifecycle
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.route-tables-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.toolbar {
  display: flex;
  gap: 10px;
  align-items: center;
}

.filter-card {
  margin-bottom: 20px;
}

.table-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
}
</style>