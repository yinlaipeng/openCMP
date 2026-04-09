<template>
  <div class="project-context-notice">
    <el-alert
      :title="`当前项目: ${currentProject?.name || '加载中...'}`"
      type="info"
      show-icon
      :closable="false"
      class="project-alert"
    >
      <p>在此模式下，您只能查看和管理与此项目相关的安全告警。</p>
    </el-alert>
  </div>

  <el-card class="section-card">
    <template #header>
      <div class="card-header">
        <span>项目安全告警</span>
      </div>
    </template>

    <el-table :data="projectAlerts" v-loading="alertsLoading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="title" label="标题" width="200" />
      <el-table-column prop="severity" label="严重程度" width="120">
        <template #default="{ row }">
          <el-tag :type="getSeverityType(row.severity)">{{ row.severity }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="发生时间" width="180" />
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <el-button size="small" @click="viewAlertDetail(row)">查看详情</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { getProjectSecurityAlerts } from '@/api/iam'
import { getCurrentProjectId, getCurrentProjectName } from '@/utils/projectContext'

// 从项目上下文获取当前项目ID和名称
const selectedProjectId = computed(() => getCurrentProjectId())
const selectedProjectName = computed(() => getCurrentProjectName())

// 项目数据
const currentProject = ref<any>({
  id: selectedProjectId.value,
  name: selectedProjectName.value
})

const projectAlerts = ref<any[]>([])

// 加载状态
const alertsLoading = ref(false)

// 获取项目相关安全告警
const loadProjectAlerts = async () => {
  alertsLoading.value = true
  try {
    if (selectedProjectId.value) {
      // 从API获取项目相关的安全告警
      const response = await getProjectSecurityAlerts(selectedProjectId.value)
      projectAlerts.value = response.data.items || []
    }
  } catch (error) {
    console.error('Failed to load project alerts:', error)
    ElMessage.error('加载项目告警失败')
  } finally {
    alertsLoading.value = false
  }
}

// 根据严重程度返回标签类型
const getSeverityType = (severity: string) => {
  switch (severity.toLowerCase()) {
    case 'critical':
      return 'danger'
    case 'high':
      return 'warning'
    case 'medium':
      return 'info'
    case 'low':
      return 'success'
    default:
      return 'info'
  }
}

// 根据状态返回标签类型
const getStatusType = (status: string) => {
  return status === 'open' ? 'danger' : 'success'
}

// 查看告警详情
const viewAlertDetail = (row: any) => {
  ElMessage.info(`查看告警 ${row.id} 详情`)
}

onMounted(() => {
  loadProjectAlerts()
})
</script>

<style scoped>
.project-context-notice {
  margin-bottom: 20px;
}

.project-alert {
  border-radius: 4px;
}

.section-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

:deep(.el-card__header) {
  padding: 12px 20px;
  background-color: #fafafa;
  border-bottom: 1px solid #eee;
}
</style>