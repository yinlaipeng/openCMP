<template>
  <div class="project-context-notice">
    <el-alert
      :title="`当前项目: ${currentProject?.name || '加载中...'}`"
      type="info"
      show-icon
      :closable="false"
      class="project-alert"
    >
      <p>在此模式下，您只能查看和管理与此项目相关的机器人。</p>
    </el-alert>
  </div>

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
import { getProjectRobots, createRobot, updateRobot, deleteRobot as deleteRobotApi, toggleRobotStatus as toggleRobotStatusApi } from '@/api/iam'
import { getCurrentProjectId, getCurrentProjectName } from '@/utils/projectContext'

// 从项目上下文获取当前项目ID和名称
const selectedProjectId = computed(() => getCurrentProjectId())
const selectedProjectName = computed(() => getCurrentProjectName())

// 项目数据
const currentProject = ref<any>({
  id: selectedProjectId.value,
  name: selectedProjectName.value
})

const projectRobots = ref<any[]>([])

// 加载状态
const robotsLoading = ref(false)

// 获取项目相关机器人
const loadProjectRobots = async () => {
  robotsLoading.value = true
  try {
    if (selectedProjectId.value) {
      // 从API获取项目相关的机器人
      const response = await getProjectRobots(selectedProjectId.value)
      projectRobots.value = response.data.items || []
    }
  } catch (error) {
    console.error('Failed to load project robots:', error)
    ElMessage.error('加载项目机器人失败')
  } finally {
    robotsLoading.value = false
  }
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
  loadProjectRobots()
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