<template>
  <div class="project-context-notice">
    <el-alert
      :title="`当前项目: ${currentProject?.name || '加载中...'}`"
      type="info"
      show-icon
      :closable="false"
      class="project-alert"
    >
      <p>在此模式下，您只能查看和管理与此项目相关的站内信。</p>
    </el-alert>
  </div>

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
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { getProjectMessages } from '@/api/iam'
import { getCurrentProjectId, getCurrentProjectName } from '@/utils/projectContext'

// 从项目上下文获取当前项目ID和名称
const selectedProjectId = computed(() => getCurrentProjectId())
const selectedProjectName = computed(() => getCurrentProjectName())

// 项目数据
const currentProject = ref<any>({
  id: selectedProjectId.value,
  name: selectedProjectName.value
})

const projectMessages = ref<any[]>([])

// 加载状态
const messagesLoading = ref(false)

// 获取项目相关消息
const loadProjectMessages = async () => {
  messagesLoading.value = true
  try {
    if (selectedProjectId.value) {
      // 从API获取项目相关的消息
      const response = await getProjectMessages(selectedProjectId.value)
      projectMessages.value = response.items || []
    }
  } catch (error) {
    console.error('Failed to load project messages:', error)
    ElMessage.error('加载项目消息失败')
  } finally {
    messagesLoading.value = false
  }
}

// 根据消息状态返回标签类型
const getMessageStatusType = (status: string) => {
  return status === 'unread' ? 'warning' : 'info'
}

// 查看消息详情
const viewMessageDetail = (row: any) => {
  ElMessage.info(`查看消息 ${row.id} 详情`)
}

onMounted(() => {
  loadProjectMessages()
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