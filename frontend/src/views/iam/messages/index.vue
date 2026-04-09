<template>
  <div class="messages-page">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="menu-card">
          <el-menu :default-active="activeFilter" @select="handleFilter">
            <el-menu-item index="all">
              <el-icon><Message /></el-icon>
              全部消息
            </el-menu-item>
            <el-menu-item index="unread">
              <el-icon><Bell /></el-icon>
              未读消息
              <el-badge v-if="unreadCount > 0" :value="unreadCount" class="badge" />
            </el-menu-item>
            <el-menu-item index="read">
              <el-icon><Checked /></el-icon>
              已读消息
            </el-menu-item>
          </el-menu>
        </el-card>
      </el-col>
      
      <el-col :span="18">
        <el-card class="messages-card">
          <template #header>
            <div class="card-header">
              <span class="title">消息列表</span>
              <el-button @click="handleMarkAllRead">全部已读</el-button>
            </div>
          </template>
          
          <div v-loading="loading">
            <el-empty v-if="messages.length === 0" description="暂无消息" />

            <div v-else class="message-list">
              <div
                v-for="msg in messages"
                :key="msg.id"
                class="message-item"
                :class="{ unread: !msg.read_at }"
                @click="handleView(msg)"
              >
                <div class="message-icon">
                  <el-icon :size="24" :class="getLevelClass(msg.level)">
                    <component :is="getLevelIcon(msg.level)" />
                  </el-icon>
                </div>
                <div class="message-content">
                  <div class="message-title">{{ msg.title }}</div>
                  <div class="message-info">
                    <span class="message-type">{{ getTypeName(msg.type) }}</span>
                    <span class="message-time">{{ formatTime(msg.created_at) }}</span>
                  </div>
                </div>
                <el-tag v-if="!msg.read_at" size="small" type="danger">未读</el-tag>
                <el-button
                  size="small"
                  type="danger"
                  link
                  style="margin-left: 8px"
                  @click.stop="handleDelete(msg)"
                >
                  <el-icon><Delete /></el-icon>
                </el-button>
              </div>
            </div>

            <el-pagination
              v-if="pagination.total > pagination.pageSize"
              v-model:current-page="pagination.page"
              v-model:page-size="pagination.pageSize"
              :total="pagination.total"
              layout="prev, pager, next"
              class="msg-pagination"
              @current-change="loadMessages"
            />
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-dialog v-model="detailVisible" title="消息详情" width="500px">
      <div v-if="currentMessage" class="message-detail">
        <h3>{{ currentMessage.title }}</h3>
        <div class="meta">
          <el-tag size="small" style="margin-right: 8px">{{ getTypeName(currentMessage.type) }}</el-tag>
          <span>{{ formatTime(currentMessage.created_at) }}</span>
        </div>
        <div class="content">{{ currentMessage.content }}</div>
      </div>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
        <el-button type="danger" @click="() => { handleDelete(currentMessage); detailVisible = false }">删除</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Message, Bell, Checked, Warning, InfoFilled, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getMessages, getUnreadCount, markRead, markAllRead, deleteMessage } from '@/api/message'

const messages = ref<any[]>([])
const unreadCount = ref(0)
const detailVisible = ref(false)
const currentMessage = ref<any>(null)
const loading = ref(false)
const activeFilter = ref('all')
const pagination = ref({ page: 1, pageSize: 20, total: 0 })

// 从 localStorage 读取当前用户 ID
const currentUserId = computed(() => {
  try {
    const user = JSON.parse(localStorage.getItem('user') || '{}')
    return user.id || 0
  } catch { return 0 }
})

const getLevelClass = (level: string) => {
  const map: Record<string, string> = { info: 'text-info', warning: 'text-warning', error: 'text-danger' }
  return map[level] || ''
}

const getLevelIcon = (level: string) => {
  const map: Record<string, any> = { info: InfoFilled, warning: Warning, error: Warning }
  return map[level] || InfoFilled
}

const getTypeName = (type: string) => {
  const map: Record<string, string> = { system: '系统通知', security: '安全告警', resource: '资源通知', task: '任务通知' }
  return map[type] || type
}

const formatTime = (time: string) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

const loadMessages = async () => {
  if (!currentUserId.value) return
  loading.value = true
  try {
    const params: any = {
      user_id: currentUserId.value,
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    }
    if (activeFilter.value === 'unread') params.read = false
    if (activeFilter.value === 'read') params.read = true

    const res = await getMessages(params)
    messages.value = res.items || []
    pagination.value.total = res.total || 0

    const countRes = await getUnreadCount(currentUserId.value)
    unreadCount.value = countRes.count || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

const handleFilter = (filter: string) => {
  activeFilter.value = filter
  pagination.value.page = 1
  loadMessages()
}

const handleView = async (msg: any) => {
  currentMessage.value = msg
  detailVisible.value = true
  if (!msg.read_at) {
    try {
      await markRead(msg.id)
      msg.read_at = new Date().toISOString()
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    } catch { /* ignore */ }
  }
}

const handleMarkAllRead = async () => {
  if (!currentUserId.value) return
  try {
    await markAllRead(currentUserId.value)
    ElMessage.success('已全部标记为已读')
    loadMessages()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

const handleDelete = async (msg: any) => {
  try {
    await ElMessageBox.confirm('确定删除该消息吗？', '提示', { type: 'warning' })
    await deleteMessage(msg.id)
    ElMessage.success('删除成功')
    loadMessages()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error(e.message || '删除失败')
  }
}

onMounted(loadMessages)
</script>

<style scoped>
.messages-page {
  height: 100%;
}

.menu-card {
  margin-bottom: 20px;
}

.messages-card {
  min-height: 500px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 16px;
  font-weight: bold;
}

.badge {
  margin-left: auto;
}

.message-list {
  max-height: 500px;
  overflow-y: auto;
}

.message-item {
  display: flex;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  transition: background-color 0.3s;
}

.message-item:hover {
  background-color: #f5f7fa;
}

.message-item.unread {
  background-color: #ecf5ff;
}

.message-icon {
  margin-right: 15px;
}

.message-content {
  flex: 1;
}

.message-title {
  font-size: 14px;
  color: #333;
  margin-bottom: 5px;
}

.message-info {
  font-size: 12px;
  color: #999;
}

.message-type {
  margin-right: 10px;
}

.message-detail h3 {
  margin-bottom: 15px;
}

.message-detail .meta {
  margin-bottom: 20px;
  color: #999;
}

.message-detail .content {
  line-height: 1.8;
  color: #666;
}

.text-info {
  color: #409EFF;
}

.text-warning {
  color: #E6A23C;
}

.text-danger {
  color: #F56C6C;
}

.msg-pagination {
  margin-top: 16px;
  justify-content: center;
}
</style>
