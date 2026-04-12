<template>
  <div class="monitoring-dashboards-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">监控面板</span>
          <el-button type="primary" @click="handleCreateDashboard">新建面板</el-button>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="8" v-for="dashboard in dashboards" :key="dashboard.id">
          <el-card shadow="hover" class="dashboard-card" @click="viewDashboard(dashboard)">
            <template #header>
              <div class="dashboard-header">
                <span class="dashboard-name">{{ dashboard.name }}</span>
                <el-dropdown @command="handleDashboardCommand">
                  <el-icon><More /></el-icon>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item :command="`edit_${dashboard.id}`">编辑</el-dropdown-item>
                      <el-dropdown-item :command="`delete_${dashboard.id}`">删除</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </template>
            <div class="dashboard-preview">
              <div class="preview-item">
                <span class="preview-label">图表数量</span>
                <span class="preview-value">{{ dashboard.chart_count }}</span>
              </div>
              <div class="preview-item">
                <span class="preview-label">创建时间</span>
                <span class="preview-value">{{ dashboard.created_at }}</span>
              </div>
              <div class="preview-item">
                <span class="preview-label">最后更新</span>
                <span class="preview-value">{{ dashboard.updated_at }}</span>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>

    <!-- Create Dashboard Modal -->
    <el-dialog v-model="createDialogVisible" title="新建监控面板" width="500px">
      <el-form :model="newDashboard" label-width="100px">
        <el-form-item label="面板名称">
          <el-input v-model="newDashboard.name" placeholder="请输入面板名称" />
        </el-form-item>
        <el-form-item label="面板描述">
          <el-input v-model="newDashboard.description" type="textarea" placeholder="请输入面板描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCreateDashboard">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { More } from '@element-plus/icons-vue'

interface Dashboard {
  id: string
  name: string
  chart_count: number
  created_at: string
  updated_at: string
}

const dashboards = ref<Dashboard[]>([])
const createDialogVisible = ref(false)
const newDashboard = ref({
  name: '',
  description: ''
})

const loadDashboards = async () => {
  dashboards.value = [
    { id: 'dash-1', name: '主机监控大盘', chart_count: 12, created_at: '2024-01-01', updated_at: '2024-01-15' },
    { id: 'dash-2', name: '网络监控大盘', chart_count: 8, created_at: '2024-01-02', updated_at: '2024-01-14' },
    { id: 'dash-3', name: '数据库监控大盘', chart_count: 6, created_at: '2024-01-03', updated_at: '2024-01-13' },
    { id: 'dash-4', name: '应用性能监控', chart_count: 10, created_at: '2024-01-04', updated_at: '2024-01-12' },
    { id: 'dash-5', name: '资源使用趋势', chart_count: 4, created_at: '2024-01-05', updated_at: '2024-01-11' },
    { id: 'dash-6', name: '告警统计面板', chart_count: 5, created_at: '2024-01-06', updated_at: '2024-01-10' }
  ]
}

const viewDashboard = (dashboard: Dashboard) => {
  ElMessage.info(`查看面板: ${dashboard.name}`)
}

const handleCreateDashboard = () => {
  newDashboard.value = { name: '', description: '' }
  createDialogVisible.value = true
}

const confirmCreateDashboard = () => {
  if (!newDashboard.value.name) {
    ElMessage.warning('请输入面板名称')
    return
  }
  dashboards.value.push({
    id: `dash-${Date.now()}`,
    name: newDashboard.value.name,
    chart_count: 0,
    created_at: new Date().toISOString().split('T')[0],
    updated_at: new Date().toISOString().split('T')[0]
  })
  createDialogVisible.value = false
  ElMessage.success('面板创建成功')
}

const handleDashboardCommand = async (command: string) => {
  const [action, id] = command.split('_')
  const dashboard = dashboards.value.find(d => d.id === id)

  if (action === 'edit') {
    ElMessage.info(`编辑面板: ${dashboard?.name}`)
  } else if (action === 'delete') {
    try {
      await ElMessageBox.confirm(`确认删除面板 ${dashboard?.name}？`, '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      dashboards.value = dashboards.value.filter(d => d.id !== id)
      ElMessage.success('删除成功')
    } catch (e) {
      console.error(e)
    }
  }
}

onMounted(() => {
  loadDashboards()
})
</script>

<style scoped>
.monitoring-dashboards-page {
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

.dashboard-card {
  cursor: pointer;
  margin-bottom: 20px;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dashboard-name {
  font-weight: bold;
}

.dashboard-preview {
  padding: 10px 0;
}

.preview-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.preview-label {
  color: #909399;
}

.preview-value {
  color: #303133;
}
</style>