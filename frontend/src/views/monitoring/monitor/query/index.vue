<template>
  <div class="monitoring-query-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">监控查询</span>
        </div>
      </template>

      <!-- Query Form -->
      <el-form :model="queryForm" inline class="query-form">
        <el-form-item label="资源类型">
          <el-select v-model="queryForm.resource_type" placeholder="请选择资源类型" style="width: 150px;">
            <el-option label="虚拟机" value="vm" />
            <el-option label="数据库" value="database" />
            <el-option label="负载均衡" value="loadbalancer" />
            <el-option label="存储" value="storage" />
          </el-select>
        </el-form-item>
        <el-form-item label="监控指标">
          <el-select v-model="queryForm.metric" placeholder="请选择监控指标" style="width: 200px;">
            <el-option label="CPU使用率" value="cpu_usage" />
            <el-option label="内存使用率" value="memory_usage" />
            <el-option label="磁盘使用率" value="disk_usage" />
            <el-option label="网络流量" value="network_traffic" />
            <el-option label="连接数" value="connections" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="queryForm.time_range"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            style="width: 350px;"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- Resource Selection -->
      <el-card shadow="hover" class="resource-card">
        <template #header>
          <span>选择资源</span>
        </template>
        <el-table :data="resources" @selection-change="handleSelectionChange" size="small">
          <el-table-column type="selection" width="55" />
          <el-table-column prop="name" label="资源名称" width="180" />
          <el-table-column prop="type" label="资源类型" width="100" />
          <el-table-column prop="ip" label="IP地址" width="140" />
          <el-table-column prop="platform" label="平台" width="100" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-card>

      <!-- Query Results -->
      <el-card shadow="hover" class="result-card">
        <template #header>
          <div class="card-header">
            <span>查询结果</span>
            <el-button type="primary" size="small" @click="handleExport">导出数据</el-button>
          </div>
        </template>
        <div class="result-chart">
          <div v-if="selectedResources.length === 0" class="empty-tip">
            <el-empty description="请选择资源后进行查询" />
          </div>
          <div v-else class="mock-chart-area">
            <div class="chart-title">{{ queryForm.metric }} 监控数据</div>
            <div class="chart-content">
              <div class="chart-line">
                <div class="line-label">最大值</div>
                <div class="line-value">85%</div>
              </div>
              <div class="chart-line">
                <div class="line-label">平均值</div>
                <div class="line-value">45%</div>
              </div>
              <div class="chart-line">
                <div class="line-label">最小值</div>
                <div class="line-value">12%</div>
              </div>
            </div>
          </div>
        </div>
      </el-card>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

interface Resource {
  id: string
  name: string
  type: string
  ip: string
  platform: string
  status: string
}

const queryForm = ref({
  resource_type: 'vm',
  metric: 'cpu_usage',
  time_range: []
})

const resources = ref<Resource[]>([])
const selectedResources = ref<Resource[]>([])

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'running':
    case 'active':
      return 'success'
    case 'stopped':
      return 'danger'
    default:
      return 'info'
  }
}

const loadResources = async () => {
  resources.value = [
    { id: 'vm-1', name: 'prod-web-01', type: '虚拟机', ip: '192.168.1.10', platform: '阿里云', status: 'Running' },
    { id: 'vm-2', name: 'prod-web-02', type: '虚拟机', ip: '192.168.1.11', platform: '阿里云', status: 'Running' },
    { id: 'vm-3', name: 'prod-api-01', type: '虚拟机', ip: '192.168.1.20', platform: '阿里云', status: 'Running' },
    { id: 'db-1', name: 'prod-mysql', type: '数据库', ip: '192.168.2.10', platform: '阿里云', status: 'Running' },
    { id: 'lb-1', name: 'prod-lb', type: '负载均衡', ip: '192.168.3.10', platform: '阿里云', status: 'Active' }
  ]
}

const handleSelectionChange = (selection: Resource[]) => {
  selectedResources.value = selection
}

const handleQuery = () => {
  if (selectedResources.value.length === 0) {
    ElMessage.warning('请选择至少一个资源')
    return
  }
  ElMessage.success('查询成功')
}

const handleReset = () => {
  queryForm.value = {
    resource_type: 'vm',
    metric: 'cpu_usage',
    time_range: []
  }
  selectedResources.value = []
}

const handleExport = () => {
  ElMessage.info('导出数据功能开发中')
}

onMounted(() => {
  loadResources()
})
</script>

<style scoped>
.monitoring-query-page {
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

.query-form {
  margin-bottom: 20px;
}

.resource-card {
  margin-bottom: 20px;
}

.result-card {
  margin-bottom: 20px;
}

.result-chart {
  height: 300px;
}

.empty-tip {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.mock-chart-area {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.chart-title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 20px;
  text-align: center;
}

.chart-content {
  display: flex;
  justify-content: space-around;
}

.chart-line {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.line-label {
  color: #909399;
  margin-bottom: 5px;
}

.line-value {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
}
</style>