<template>
  <div class="internal-messages-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">站内信</span>
        </div>
      </template>

      <!-- 消息列表 -->
      <el-table :data="messages" v-loading="loading" border stripe>
        <el-table-column prop="title" label="标题" min-width="150" show-overflow-tooltip />
        <el-table-column prop="level" label="级别" width="100">
          <template #default="{ row }">
            <el-tag :type="getLevelType(row.level)">
              {{ getLevelText(row.level) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="接收人" width="120">
          <template #default="{ row }">
            {{ getReceiverName(row.receiver_id) }}
          </template>
        </el-table-column>
        <el-table-column label="消息类型" width="120">
          <template #default="{ row }">
            {{ getMessageTypeName(row.type_id) }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="接收时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="内容" width="100">
          <template #default="{ row }">
            <el-button size="small" @click="showContent(row)">查看</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadMessages"
        @current-change="loadMessages"
        class="pagination"
      />
    </el-card>

    <!-- 消息内容弹窗 -->
    <el-dialog v-model="contentDialogVisible" title="消息内容" width="600px">
      <div class="message-content" v-if="currentMessage">
        <h4>{{ currentMessage.title }}</h4>
        <div class="message-meta">
          <el-tag :type="getLevelType(currentMessage.level)" style="margin-right: 10px;">
            {{ getLevelText(currentMessage.level) }}
          </el-tag>
          <span>接收人: {{ getReceiverName(currentMessage.receiver_id) }}</span>
          <span style="margin-left: 10px;">消息类型: {{ getMessageTypeName(currentMessage.type_id) }}</span>
          <span style="margin-left: 10px;">接收时间: {{ formatDate(currentMessage.created_at) }}</span>
        </div>
        <div class="message-body" style="margin-top: 20px; white-space: pre-wrap;">
          {{ currentMessage.content }}
        </div>
      </div>
      <template #footer>
        <el-button @click="contentDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getMessages, getMessage, markRead, getMessageTypes } from '@/api/message'
import { getUsers } from '@/api/iam' // Assuming we have an API to get users

interface Message {
  id: number
  title: string
  content: string
  level: string // info, warning, error
  sender_id: number
  receiver_id: number
  type_id: number
  read: boolean
  read_at?: string
  created_at: string
}

interface MessageType {
  id: number
  name: string
  display_name: string
}

interface User {
  id: number
  name: string
  display_name: string
}

const messages = ref<Message[]>([])
const messageTypes = ref<MessageType[]>([])
const users = ref<User[]>([])
const loading = ref(false)
const contentDialogVisible = ref(false)
const currentMessage = ref<Message | null>(null)

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 从 localStorage 获取当前用户 ID
const getCurrentUserId = () => {
  try {
    const user = JSON.parse(localStorage.getItem('user') || '{}')
    return user.id || 0
  } catch {
    return 0
  }
}

const loadMessages = async () => {
  loading.value = true
  try {
    const userId = getCurrentUserId()
    if (!userId) {
      ElMessage.error('无法获取当前用户信息')
      return
    }

    const params = {
      user_id: userId,
      page: pagination.page,
      page_size: pagination.pageSize
    }

    const res = await getMessages(params)
    messages.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载消息列表失败')
  } finally {
    loading.value = false
  }
}

const loadMessageTypes = async () => {
  try {
    const res = await getMessageTypes()
    messageTypes.value = res.items || [] // Changed from .data to .items
  } catch (e: any) {
    console.error('加载消息类型失败:', e)
  }
}

const loadUsers = async () => {
  try {
    const res = await getUsers({ limit: 1000 }) // 获取所有用户
    users.value = res.items || []
  } catch (e: any) {
    console.error('加载用户列表失败:', e)
  }
}

const getLevelType = (level: string) => {
  switch (level) {
    case 'info':
      return 'info'
    case 'warning':
      return 'warning'
    case 'error':
      return 'danger'
    case 'success':
      return 'success'
    default:
      return 'info'
  }
}

const getLevelText = (level: string) => {
  const levelMap: Record<string, string> = {
    info: '信息',
    warning: '警告',
    error: '错误',
    success: '成功'
  }
  return levelMap[level] || level
}

const getReceiverName = (receiverId: number) => {
  const user = users.value.find(u => u.id === receiverId)
  return user ? user.display_name || user.name : `用户${receiverId}`
}

const getMessageTypeName = (typeId: number) => {
  const msgType = messageTypes.value.find(t => t.id === typeId)
  return msgType ? msgType.display_name || msgType.name : `类型${typeId}`
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  try {
    return new Date(dateStr).toLocaleString('zh-CN')
  } catch {
    return dateStr
  }
}

const showContent = async (msg: Message) => {
  currentMessage.value = msg
  contentDialogVisible.value = true

  // 如果消息未读，自动标记为已读
  if (!msg.read_at) {
    try {
      await markRead(msg.id)
      msg.read_at = new Date().toISOString()
      // 更新本地消息列表中的已读状态
      const index = messages.value.findIndex(m => m.id === msg.id)
      if (index !== -1) {
        messages.value[index].read_at = msg.read_at
      }
    } catch (e) {
      console.error('标记已读失败:', e)
    }
  }
}

onMounted(async () => {
  await Promise.all([
    loadMessageTypes(),
    loadUsers(),
    loadMessages()
  ])
})
</script>

<style scoped>
.internal-messages-page {
  height: 100%;
  padding: 20px;
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

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.message-content h4 {
  margin-top: 0;
  margin-bottom: 15px;
}

.message-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 15px;
  font-size: 14px;
  color: #666;
}

.message-body {
  line-height: 1.6;
  color: #333;
}
</style>