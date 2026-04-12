<template>
  <div class="zones-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">可用区列表</span>
        </div>
      </template>

      <el-table :data="zones" v-loading="loading">
        <el-table-column label="名称" width="200">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="l2_network_count" label="二层网络" width="120" />
        <el-table-column label="宿主机/可用宿主机" width="180">
          <template #default="{ row }">
            {{ row.host_count }}/{{ row.available_host_count }}
          </template>
        </el-table-column>
        <el-table-column label="物理机/可用物理机" width="180">
          <template #default="{ row }">
            {{ row.physical_host_count }}/{{ row.available_physical_host_count }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="可用区详情" width="800px">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="详情" name="detail">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="可用区ID">{{ selectedZone?.id }}</el-descriptions-item>
            <el-descriptions-item label="可用区名称">{{ selectedZone?.name }}</el-descriptions-item>
            <el-descriptions-item label="区域ID">{{ selectedZone?.region_id }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ selectedZone?.status }}</el-descriptions-item>
            <el-descriptions-item label="二层网络数">{{ selectedZone?.l2_network_count }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ selectedZone?.created_at }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="宿主机" name="hosts">
          <el-table :data="hosts" v-loading="hostsLoading">
            <el-table-column prop="id" label="宿主机ID" width="180" />
            <el-table-column prop="name" label="宿主机名称" width="150" />
            <el-table-column prop="status" label="状态" width="100" />
            <el-table-column prop="ip" label="IP地址" width="150" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="资源统计" name="stats">
          <el-row :gutter="20">
            <el-col :span="6">
              <el-statistic title="总宿主机数" :value="selectedZone?.host_count || 0" />
            </el-col>
            <el-col :span="6">
              <el-statistic title="可用宿主机" :value="selectedZone?.available_host_count || 0" />
            </el-col>
            <el-col :span="6">
              <el-statistic title="总物理机数" :value="selectedZone?.physical_host_count || 0" />
            </el-col>
            <el-col :span="6">
              <el-statistic title="可用物理机" :value="selectedZone?.available_physical_host_count || 0" />
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
import { ElMessage, ElMessageBox } from 'element-plus'

const zones = ref<any[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedZone = ref<any | null>(null)
const activeTab = ref('detail')

// Detail modal related data
const hosts = ref<any[]>([])
const hostsLoading = ref(false)
const logs = ref<any[]>([])
const logsLoading = ref(false)

const loadZones = async () => {
  loading.value = true
  try {
    // Mock data - API requires valid cloud account credentials
    zones.value = [
      {
        id: 'cn-hangzhou-a',
        name: '华东1可用区A',
        status: 'Available',
        region_id: 'cn-hangzhou',
        l2_network_count: 3,
        host_count: 15,
        available_host_count: 8,
        physical_host_count: 5,
        available_physical_host_count: 3,
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'cn-hangzhou-b',
        name: '华东1可用区B',
        status: 'Available',
        region_id: 'cn-hangzhou',
        l2_network_count: 2,
        host_count: 10,
        available_host_count: 5,
        physical_host_count: 3,
        available_physical_host_count: 2,
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'cn-shanghai-a',
        name: '华东2可用区A',
        status: 'Available',
        region_id: 'cn-shanghai',
        l2_network_count: 2,
        host_count: 12,
        available_host_count: 6,
        physical_host_count: 4,
        available_physical_host_count: 2,
        created_at: '2024-01-01 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    zones.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = async (row: Zone) => {
  selectedZone.value = row
  detailDialogVisible.value = true
  activeTab.value = 'detail'
  await loadZoneDetails(row.id)
}

const loadZoneDetails = async (zoneId: string) => {
  // Load hosts (mock data)
  hostsLoading.value = true
  hosts.value = [
    { id: 'host-1', name: 'Host 1', status: 'Running', ip: '192.168.1.1' },
    { id: 'host-2', name: 'Host 2', status: 'Running', ip: '192.168.1.2' }
  ]
  hostsLoading.value = false

  // Load logs (mock data)
  logsLoading.value = true
  logs.value = [
    { operation: '创建', operator: 'admin', result: '成功', timestamp: '2024-01-01 10:00:00' }
  ]
  logsLoading.value = false
}

const handleDelete = async (row: Zone) => {
  try {
    await ElMessageBox.confirm(`确认删除可用区 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    // In real implementation, call delete API
    ElMessage.success('删除成功')
    loadZones()
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadZones()
})
</script>

<style scoped>
.zones-page {
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