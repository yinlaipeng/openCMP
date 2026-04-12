<template>
  <div class="monitoring-dashboard-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">监控大盘</span>
        </div>
      </template>

      <!-- 统计卡片区域 -->
      <el-row :gutter="20" class="stat-cards">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #409eff;">
                <el-icon><Monitor /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">128</div>
                <div class="stat-label">监控资源数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #67c23a;">
                <el-icon><CircleCheck /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">120</div>
                <div class="stat-label">正常资源</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #e6a23c;">
                <el-icon><Warning /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">8</div>
                <div class="stat-label">告警资源</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #f56c6c;">
                <el-icon><CircleClose /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">0</div>
                <div class="stat-label">异常资源</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 图表区域 -->
      <el-row :gutter="20" class="chart-row">
        <el-col :span="12">
          <el-card shadow="hover">
            <template #header>
              <span>资源监控状态分布</span>
            </template>
            <div class="chart-placeholder">
              <el-progress type="dashboard" :percentage="93" :color="['#67c23a', '#e6a23c', '#f56c6c']">
                <template #default>
                  <span class="percentage-value">93%</span>
                  <span class="percentage-label">健康率</span>
                </template>
              </el-progress>
            </div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card shadow="hover">
            <template #header>
              <span>告警趋势</span>
            </template>
            <div class="chart-placeholder">
              <div class="mock-chart">
                <div class="chart-bar" style="height: 30%; background: #409eff;"></div>
                <div class="chart-bar" style="height: 45%; background: #409eff;"></div>
                <div class="chart-bar" style="height: 20%; background: #409eff;"></div>
                <div class="chart-bar" style="height: 60%; background: #e6a23c;"></div>
                <div class="chart-bar" style="height: 35%; background: #409eff;"></div>
                <div class="chart-bar" style="height: 25%; background: #409eff;"></div>
                <div class="chart-bar" style="height: 40%; background: #409eff;"></div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 告警列表 -->
      <el-card shadow="hover" class="alert-card">
        <template #header>
          <div class="card-header">
            <span>近期告警</span>
            <el-button type="primary" size="small" @click="goToAlerts">查看全部</el-button>
          </div>
        </template>
        <el-table :data="recentAlerts" size="small">
          <el-table-column prop="resource_name" label="资源名称" width="180" />
          <el-table-column prop="alert_type" label="告警类型" width="120" />
          <el-table-column prop="level" label="告警级别" width="100">
            <template #default="{ row }">
              <el-tag :type="getAlertLevelType(row.level)" size="small">{{ row.level }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="trigger_time" label="触发时间" width="160" />
          <el-table-column prop="platform" label="平台" width="100" />
        </el-table>
      </el-card>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Monitor, CircleCheck, Warning, CircleClose } from '@element-plus/icons-vue'

const router = useRouter()

const recentAlerts = ref([
  {
    resource_name: 'prod-web-01',
    alert_type: 'CPU使用率',
    level: '严重',
    status: '未处理',
    trigger_time: '2024-01-15 10:30:00',
    platform: '阿里云'
  },
  {
    resource_name: 'prod-db-01',
    alert_type: '内存使用率',
    level: '警告',
    status: '已处理',
    trigger_time: '2024-01-15 09:20:00',
    platform: '阿里云'
  },
  {
    resource_name: 'dev-api-02',
    alert_type: '磁盘使用率',
    level: '信息',
    status: '已屏蔽',
    trigger_time: '2024-01-15 08:15:00',
    platform: '阿里云'
  }
])

const getAlertLevelType = (level: string) => {
  switch (level) {
    case '严重':
      return 'danger'
    case '警告':
      return 'warning'
    case '信息':
      return 'info'
    default:
      return ''
  }
}

const getStatusType = (status: string) => {
  switch (status) {
    case '未处理':
      return 'danger'
    case '已处理':
      return 'success'
    case '已屏蔽':
      return 'info'
    default:
      return ''
  }
}

const goToAlerts = () => {
  router.push('/monitoring/alerts/resources')
}

onMounted(() => {
  // Load dashboard data
})
</script>

<style scoped>
.monitoring-dashboard-page {
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

.stat-cards {
  margin-bottom: 20px;
}

.stat-card {
  height: 100px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 28px;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.chart-row {
  margin-bottom: 20px;
}

.chart-placeholder {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.percentage-value {
  font-size: 24px;
  font-weight: bold;
}

.percentage-label {
  font-size: 14px;
  color: #909399;
}

.mock-chart {
  display: flex;
  align-items: flex-end;
  gap: 10px;
  height: 180px;
  padding: 10px;
}

.chart-bar {
  width: 40px;
  border-radius: 4px 4px 0 0;
}

.alert-card {
  margin-top: 20px;
}
</style>