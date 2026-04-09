<template>
  <div class="project-context-notice">
    <el-alert
      :title="`当前项目: ${currentProjectName || '加载中...'}`"
      type="info"
      show-icon
      :closable="false"
      class="project-alert"
    >
      <p>在此模式下，您只能查看和管理与此项目相关的资源和告警信息。</p>
    </el-alert>
  </div>

  <!-- 项目相关的安全告警 -->
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

  <!-- 项目相关的站内信 -->
  <el-card class="section-card">
    <template #header>
      <div class="card-header">
        <span>项目相关站内信</span>
      </div>
    </template>

    <el-table :data="projectMessages" v-loading="messagesLoading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="title" label="标题" width="200" />
      <el-table-column prop="type" label="类型" width="120">
        <template #default="{ row }">
          <el-tag>{{ row.type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getMessageStatusType(row.status)">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="时间" width="180" />
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <el-button size="small" @click="viewMessageDetail(row)">查看</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <!-- 项目相关的机器人管理 -->
  <el-card class="section-card">
    <template #header>
      <div class="card-header">
        <span>项目机器人管理</span>
      </div>
    </template>

    <el-table :data="projectRobots" v-loading="robotsLoading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" width="200" />
      <el-table-column prop="type" label="类型" width="150">
        <template #default="{ row }">
          <el-tag>{{ getRobotTypeLabel(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-switch
            v-model="row.enabled"
            @change="toggleRobotStatus(row)"
            :active-value="true"
            :inactive-value="false"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button size="small" @click="editRobot(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="deleteRobot(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getProject, getProjectSecurityAlerts, getProjectMessages, getProjectRobots, deleteRobot as deleteRobotApi, toggleRobotStatus as toggleRobotStatusApi } from '@/api/iam'
import { getCurrentProjectId, getCurrentProjectName } from '@/utils/projectContext'

// 从项目上下文获取项目ID和名称
const projectId = computed(() => getCurrentProjectId())
const currentProjectName = computed(() => getCurrentProjectName())

// 项目数据
const currentProject = ref<any>({ id: projectId.value, name: currentProjectName.value })
const projectAlerts = ref<any[]>([])
const projectMessages = ref<any[]>([])
const projectRobots = ref<any[]>([])

// 加载状态
const alertsLoading = ref(false)
const messagesLoading = ref(false)
const robotsLoading = ref(false)

// 获取项目详情和相关资源
const loadProjectData = async () => {
  try {
    // 获取项目详情
    if (projectId.value) {
      currentProject.value = await getProject(projectId.value)
    }

    // 加载项目相关资源
    await Promise.all([
      loadProjectAlerts(),
      loadProjectMessages(),
      loadProjectRobots()
    ])
  } catch (error) {
    console.error('Failed to load project data:', error)
    ElMessage.error('加载项目数据失败')
  }
}

// 加载项目安全告警
const loadProjectAlerts = async () => {
  alertsLoading.value = true
  try {
    if (projectId.value) {
      // 从API获取项目相关的安全告警
      const response = await getProjectSecurityAlerts(projectId.value)
      projectAlerts.value = response.data.items || []
    }
  } catch (error) {
    console.error('Failed to load project alerts:', error)
  } finally {
    alertsLoading.value = false
  }
}

// 加载项目站内信
const loadProjectMessages = async () => {
  messagesLoading.value = true
  try {
    if (projectId.value) {
      // 从API获取项目相关的消息
      const response = await getProjectMessages(projectId.value)
      projectMessages.value = response.data.items || []
    }
  } catch (error) {
    console.error('Failed to load project messages:', error)
  } finally {
    messagesLoading.value = false
  }
}

// 加载项目机器人
const loadProjectRobots = async () => {
  robotsLoading.value = true
  try {
    if (projectId.value) {
      // 从API获取项目相关的机器人
      const response = await getProjectRobots(projectId.value)
      projectRobots.value = response.data.items || []
    }
  } catch (error) {
    console.error('Failed to load project robots:', error)
  } finally {
    robotsLoading.value = false
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

// 根据消息状态返回标签类型
const getMessageStatusType = (status: string) => {
  return status === 'unread' ? 'warning' : 'info'
}

// 获取机器人类型标签
const getRobotTypeLabel = (type: string) => {
  const map: Record<string, string> = {
    'wechat_work': '企业微信',
    'dingtalk': '钉钉',
    'feishu': '飞书',
    'slack': 'Slack',
    'webhook': 'Webhook'
  }
  return map[type] || type
}

// 查看告警详情
const viewAlertDetail = (row: any) => {
  ElMessage.info(`查看告警 ${row.id} 详情`)
}

// 查看消息详情
const viewMessageDetail = (row: any) => {
  ElMessage.info(`查看消息 ${row.id} 详情`)
}

// 编辑机器人
const editRobot = (row: any) => {
  ElMessage.info(`编辑机器人 ${row.name}`)
}

// 删除机器人
const deleteRobot = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除机器人 "${row.name}" 吗？`, '提示', {
      type: 'warning'
    })
    // 调用API删除机器人
    await deleteRobotApi(row.id)
    ElMessage.success('机器人已删除')
    // 重新加载列表
    loadProjectRobots()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
      ElMessage.error('删除失败')
    }
  }
}

// 切换机器人状态
const toggleRobotStatus = async (row: any) => {
  try {
    // 调用API更新机器人状态
    await toggleRobotStatusApi(row.id, !row.enabled)
    const newStatus = !row.enabled
    ElMessage.success(`${row.name} 已${newStatus ? '启用' : '禁用'}`)
    // 更新本地数据
    row.enabled = newStatus
  } catch (error) {
    console.error('Failed to update robot status:', error)
    row.enabled = !row.enabled // 恢复原状态
    ElMessage.error('更新状态失败')
  }
}

onMounted(() => {
  loadProjectData()
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