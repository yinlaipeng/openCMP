<template>
  <div class="regions-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">区域列表</span>
        </div>
      </template>

      <el-table :data="regions" v-loading="loading">
        <el-table-column label="名称" width="200">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="启用状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.status === 'enabled' ? 'success' : 'info'">
              {{ row.status === 'enabled' ? '已启用' : '未启用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="available_zone_count" label="可用区" width="100">
          <template #default="{ row }">
            <el-link @click="viewZones(row)">{{ row.available_zone_count }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="vpc_count" label="vpc" width="100" />
        <el-table-column prop="virtual_machine_count" label="虚拟机" width="100" />
        <el-table-column prop="platform" label="平台" width="120" />
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="区域详情" width="800px">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="详情" name="detail">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="区域ID">{{ selectedRegion?.id }}</el-descriptions-item>
            <el-descriptions-item label="区域名称">{{ selectedRegion?.name }}</el-descriptions-item>
            <el-descriptions-item label="启用状态">
              <el-tag :type="selectedRegion?.status === 'enabled' ? 'success' : 'info'">
                {{ selectedRegion?.status === 'enabled' ? '已启用' : '未启用' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="平台">{{ selectedRegion?.platform }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ selectedRegion?.created_at }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="可用区" name="zones">
          <el-table :data="zones" v-loading="zonesLoading">
            <el-table-column prop="id" label="可用区ID" width="180" />
            <el-table-column prop="name" label="可用区名称" width="200" />
            <el-table-column prop="status" label="状态" width="120" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="vpc" name="vpcs">
          <el-table :data="vpcs" v-loading="vpcsLoading">
            <el-table-column prop="id" label="VPC ID" width="180" />
            <el-table-column prop="name" label="VPC 名称" width="150" />
            <el-table-column prop="status" label="状态" width="100" />
            <el-table-column prop="cidr" label="CIDR" width="180" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="ip子网" name="subnets">
          <el-table :data="subnets" v-loading="subnetsLoading">
            <el-table-column prop="id" label="子网ID" width="200" />
            <el-table-column prop="name" label="子网名称" width="150" />
            <el-table-column prop="status" label="状态" width="100" />
            <el-table-column prop="cidr" label="CIDR" width="180" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="资源统计" name="stats">
          <el-row :gutter="20">
            <el-col :span="6">
              <el-statistic title="总虚拟机数" :value="selectedRegion?.virtual_machine_count || 0" />
            </el-col>
            <el-col :span="6">
              <el-statistic title="总VPC数" :value="selectedRegion?.vpc_count || 0" />
            </el-col>
            <el-col :span="6">
              <el-statistic title="总子网数" :value="subnets.length" />
            </el-col>
            <el-col :span="6">
              <el-statistic title="总可用区数" :value="selectedRegion?.available_zone_count || 0" />
            </el-col>
          </el-row>
        </el-tab-pane>
        <el-tab-pane label="操作日志" name="logs">
          <el-table :data="logs" v-loading="logsLoading">
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

const regions = ref<any[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedRegion = ref<any | null>(null)
const activeTab = ref('detail')

// Detail modal related data
const zones = ref<any[]>([])
const zonesLoading = ref(false)
const vpcs = ref<any[]>([])
const vpcsLoading = ref(false)
const subnets = ref<any[]>([])
const subnetsLoading = ref(false)
const logs = ref<any[]>([])
const logsLoading = ref(false)

const loadRegions = async () => {
  loading.value = true
  try {
    // Mock data - API requires valid cloud account credentials
    regions.value = [
      {
        id: 'cn-hangzhou',
        name: '华东1(杭州)',
        status: 'enabled',
        available_zone_count: 3,
        vpc_count: 5,
        virtual_machine_count: 15,
        platform: '阿里云',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'cn-shanghai',
        name: '华东2(上海)',
        status: 'enabled',
        available_zone_count: 2,
        vpc_count: 3,
        virtual_machine_count: 10,
        platform: '阿里云',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'cn-beijing',
        name: '华北2(北京)',
        status: 'enabled',
        available_zone_count: 3,
        vpc_count: 4,
        virtual_machine_count: 12,
        platform: '阿里云',
        created_at: '2024-01-01 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    regions.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = async (row: Region) => {
  selectedRegion.value = row
  detailDialogVisible.value = true
  activeTab.value = 'detail'
  await loadRegionDetails(row.id)
}

const viewZones = async (row: Region) => {
  selectedRegion.value = row
  detailDialogVisible.value = true
  activeTab.value = 'zones'
  await loadRegionDetails(row.id)
}

const loadRegionDetails = async (regionId: string) => {
  // Mock data - API requires valid cloud account credentials
  zonesLoading.value = true
  zones.value = [
    { id: 'zone-a', name: '可用区A', status: 'Available' },
    { id: 'zone-b', name: '可用区B', status: 'Available' },
    { id: 'zone-c', name: '可用区C', status: 'Available' }
  ]
  zonesLoading.value = false

  // Mock data for VPCs
  vpcsLoading.value = true
  vpcs.value = [
    { id: 'vpc-1', name: 'VPC 1', status: 'Available', cidr: '10.0.0.0/16' },
    { id: 'vpc-2', name: 'VPC 2', status: 'Available', cidr: '10.1.0.0/16' }
  ]
  vpcsLoading.value = false

  // Mock data for subnets
  subnetsLoading.value = true
  subnets.value = [
    { id: 'subnet-1', name: 'Subnet 1', status: 'Available', cidr: '10.0.1.0/24' },
    { id: 'subnet-2', name: 'Subnet 2', status: 'Available', cidr: '10.0.2.0/24' }
  ]
  subnetsLoading.value = false

  // Load logs (mock data)
  logsLoading.value = true
  logs.value = [
    { operation: '创建', operator: 'admin', result: '成功', timestamp: '2024-01-01 10:00:00' },
    { operation: '同步', operator: 'system', result: '成功', timestamp: '2024-01-02 12:00:00' }
  ]
  logsLoading.value = false
}

onMounted(() => {
  loadRegions()
})
</script>

<style scoped>
.regions-page {
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