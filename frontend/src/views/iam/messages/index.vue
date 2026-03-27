<template>
  <div class="messages-page">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="menu-card">
          <el-menu default-active="all">
            <el-menu-item index="all">
              <el-icon><Message /></el-icon>
              全部消息
            </el-menu-item>
            <el-menu-item index="unread">
              <el-icon><Bell /></el-icon>
              未读消息
              <el-badge :value="unreadCount" class="badge" />
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
          
          <el-empty v-if="messages.length === 0" description="暂无消息" />
          
          <div v-else class="message-list">
            <div
              v-for="msg in messages"
              :key="msg.id"
              class="message-item"
              :class="{ unread: !msg.read }"
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
              <el-tag v-if="!msg.read" size="small">未读</el-tag>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-dialog v-model="detailVisible" title="消息详情" width="500px">
      <div v-if="currentMessage" class="message-detail">
        <h3>{{ currentMessage.title }}</h3>
        <div class="meta">
          <el-tag size="small">{{ getTypeName(currentMessage.type) }}</el-tag>
          <span>{{ formatTime(currentMessage.created_at) }}</span>
        </div>
        <div class="content">{{ currentMessage.content }}</div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Message, Bell, Checked, Warning, InfoFilled, SuccessFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const messages = ref<any[]>([])
const unreadCount = ref(0)
const detailVisible = ref(false)
const currentMessage = ref<any>(null)

const unreadMessages = computed(() => messages.value.filter(m => !m.read))

const getLevelClass = (level: string) => {
  const map: Record<string, string> = {
    info: 'text-info',
    warning: 'text-warning',
    error: 'text-danger'
  }
  return map[level] || ''
}

const getLevelIcon = (level: string) => {
  const map: Record<string, any> = {
    info: InfoFilled,
    warning: Warning,
    error: Warning
  }
  return map[level] || InfoFilled
}

const getTypeName = (type: string) => {
  const map: Record<string, string> = {
    system: '系统通知',
    security: '安全告警',
    resource: '资源通知',
    task: '任务通知'
  }
  return map[type] || type
}

const formatTime = (time: string) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

const loadMessages = () => {
  // 模拟数据
  messages.value = [
    {
      id: 1,
      type: 'system',
      title: '系统升级通知',
      content: '系统将于今晚 22:00 进行例行维护，预计耗时 2 小时。',
      level: 'info',
      read: false,
      created_at: new Date().toISOString()
    },
    {
      id: 2,
      type: 'security',
      title: '登录异常提醒',
      content: '检测到您的账号在异地登录，请注意账号安全。',
      level: 'warning',
      read: false,
      created_at: new Date(Date.now() - 86400000).toISOString()
    },
    {
      id: 3,
      type: 'resource',
      title: '资源使用提醒',
      content: '您的云资源使用量已达到 80%，请注意监控。',
      level: 'info',
      read: true,
      created_at: new Date(Date.now() - 172800000).toISOString()
    }
  ]
  unreadCount.value = messages.value.filter(m => !m.read).length
}

const handleView = (msg: any) => {
  currentMessage.value = msg
  detailVisible.value = true
  if (!msg.read) {
    msg.read = true
    unreadCount.value--
  }
}

const handleMarkAllRead = () => {
  messages.value.forEach(m => m.read = true)
  unreadCount.value = 0
  ElMessage.success('已全部标记为已读')
}

onMounted(() => {
  loadMessages()
})
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
</style>
