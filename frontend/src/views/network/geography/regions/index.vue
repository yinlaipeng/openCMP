<template>
  <div class="regions-container">
    <div class="page-header">
      <h2>区域</h2>
      <div class="toolbar">
        <el-button link type="primary" @click="handleView">
          <el-icon><View /></el-icon>
          查看
        </el-button>
      </div>
    </div>

    <!-- 搜索表单 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="fetchData">
        <el-form-item label="名称">
          <el-input v-model="filters.name" placeholder="搜索区域名称" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item label="启用状态">
          <el-select v-model="filters.enabled" placeholder="选择状态" clearable style="width: 120px">
            <el-option label="已启用" value="true" />
            <el-option label="未启用" value="false" />
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
        <el-form-item>
          <el-button type="primary" @click="fetchData">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <el-table
        :data="regions"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
      >
        <!-- Cloudpods 表头顺序: Name → Enabled → Zones → VPC → Server → Platform -->
        <el-table-column prop="name" label="名称" width="200">
          <template #default="{ row }">
            <el-link type="primary" :underline="false" @click="handleDetails(row)">
              {{ row.name }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="启用状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ row.enabled ? '已启用' : '未启用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="zone_count" label="可用区" width="100">
          <template #default="{ row }">
            <el-link type="primary" @click="viewZones(row)">{{ row.zone_count || 0 }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="vpc_count" label="VPC" width="100" />
        <el-table-column prop="server_count" label="服务器" width="100" />
        <el-table-column prop="provider_type" label="平台" width="120">
          <template #default="{ row }">
            <el-tag size="small" :type="getPlatformType(row.provider_type)">
              {{ getPlatformLabel(row.provider_type) }}
            </el-tag>
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
      title="区域详情"
      v-model="detailDialogVisible"
      width="700px"
    >
      <el-tabs v-model="activeTab">
        <el-tab-pane label="详情" name="detail">
          <el-descriptions :column="2" border v-if="selectedRegion">
            <el-descriptions-item label="ID">{{ selectedRegion.id }}</el-descriptions-item>
            <el-descriptions-item label="名称">{{ selectedRegion.name }}</el-descriptions-item>
            <el-descriptions-item label="启用状态">
              <el-tag :type="selectedRegion.enabled ? 'success' : 'info'">
                {{ selectedRegion.enabled ? '已启用' : '未启用' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="平台">
              <el-tag size="small" :type="getPlatformType(selectedRegion.provider_type)">
                {{ getPlatformLabel(selectedRegion.provider_type) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="可用区数">{{ selectedRegion.zone_count }}</el-descriptions-item>
            <el-descriptions-item label="VPC数">{{ selectedRegion.vpc_count }}</el-descriptions-item>
            <el-descriptions-item label="服务器数">{{ selectedRegion.server_count }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="可用区" name="zones">
          <el-table :data="zones" v-loading="zonesLoading">
            <el-table-column prop="id" label="可用区ID" width="180" />
            <el-table-column prop="name" label="可用区名称" width="200" />
            <el-table-column prop="status" label="状态" width="120" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="VPC" name="vpcs">
          <el-table :data="vpcs" v-loading="vpcsLoading">
            <el-table-column prop="id" label="VPC ID" width="180" />
            <el-table-column prop="name" label="VPC名称" width="150" />
            <el-table-column prop="status" label="状态" width="100" />
            <el-table-column prop="cidr" label="CIDR" width="180" />
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
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { View } from '@element-plus/icons-vue'
import { getRegions, Region } from '@/api/networkSync'

// Data
const loading = ref(false)
const regions = ref<Region[]>([])
const selectedRegion = ref<Region | null>(null)
const detailDialogVisible = ref(false)
const activeTab = ref('detail')

const zones = ref<any[]>([])
const zonesLoading = ref(false)
const vpcs = ref<any[]>([])
const vpcsLoading = ref(false)

const filters = reactive({
  name: '',
  enabled: '',
  platform: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// Methods
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
  filters.enabled = ''
  filters.platform = ''
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
    const res = await getRegions(params)
    regions.value = res.items || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('Failed to fetch regions:', error)
    ElMessage.error('获取区域列表失败')
  } finally {
    loading.value = false
  }
}

const handleView = () => {
  ElMessage.info('视图切换功能开发中')
}

const handleDetails = (row: Region) => {
  selectedRegion.value = row
  detailDialogVisible.value = true
  activeTab.value = 'detail'
  loadRegionDetails(row.id)
}

const viewZones = (row: Region) => {
  selectedRegion.value = row
  detailDialogVisible.value = true
  activeTab.value = 'zones'
  loadRegionDetails(row.id)
}

const loadRegionDetails = async (regionId: string) => {
  zonesLoading.value = true
  zones.value = []
  zonesLoading.value = false

  vpcsLoading.value = true
  vpcs.value = []
  vpcsLoading.value = false
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
.regions-container {
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