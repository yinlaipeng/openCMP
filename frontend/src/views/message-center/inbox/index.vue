<template>
  <div class="inbox-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">站内信</span>
          <!-- 工具栏按钮 - 参考 Cloudpods -->
          <div class="toolbar">
            <el-button @click="handleViewMode">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button @click="handleRefresh">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button @click="handleMarkAllRead" :disabled="unreadCount === 0">
              <el-icon><Check /></el-icon>
              全部标记已读
            </el-button>
          </div>
        </div>
      </template>

      <!-- 统计卡片 - 参考 Cloudpods -->
      <el-row :gutter="16" class="stats-row">
        <el-col :span="3">
          <div class="stat-card total">
            <div class="stat-icon">
              <el-icon><Message /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ pagination.total }}</div>
              <div class="stat-label">总计</div>
            </div>
          </div>
        </el-col>
        <el-col :span="3">
          <div class="stat-card unread">
            <div class="stat-icon">
              <el-icon><BellFilled /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ unreadCount }}</div>
              <div class="stat-label">未读</div>
            </div>
          </div>
        </el-col>
        <el-col :span="3">
          <div class="stat-card read">
            <div class="stat-icon">
              <el-icon><CircleCheckFilled /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ readCount }}</div>
              <div class="stat-label">已读</div>
            </div>
          </div>
        </el-col>
        <el-col :span="3">
          <div class="stat-card important">
            <div class="stat-icon">
              <el-icon><WarningFilled /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ importantCount }}</div>
              <div class="stat-label">重要</div>
            </div>
          </div>
        </el-col>
      </el-row>

      <!-- 筛选区域 -->
      <div class="filter-bar">
        <el-select v-model="filters.status" placeholder="状态" clearable style="width: 120px" @change="loadMessages">
          <el-option label="全部" value="" />
          <el-option label="未读" value="unread" />
          <el-option label="已读" value="read" />
        </el-select>
        <el-select v-model="filters.priority" placeholder="级别" clearable style="width: 120px" @change="loadMessages">
          <el-option label="全部" value="" />
          <el-option label="重要" value="important" />
          <el-option label="普通" value="normal" />
          <el-option label="低" value="low" />
        </el-select>
        <el-select v-model="filters.topic_type" placeholder="类型" clearable style="width: 150px" @change="loadMessages">
          <el-option label="全部" value="" />
          <el-option label="资源" value="resource" />
          <el-option label="系统" value="system" />
          <el-option label="安全" value="security" />
        </el-select>
        <el-button type="primary" @click="loadMessages">
          <el-icon><Search /></el-icon>
          查询
        </el-button>
        <el-button @click="resetFilters">
          <el-icon><RefreshRight /></el-icon>
          重置
        </el-button>
      </div>

      <!-- 表格 - 中文表头，参考 Cloudpods -->
      <el-table
        :data="messages"
        v-loading="loading"
        border
        stripe
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <!-- 标题 -->
        <el-table-column label="标题" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">
              <span :class="{ 'unread-title': !row.read_at }">{{ row.title }}</span>
            </el-button>
          </template>
        </el-table-column>
        <!-- 严重级别 -->
        <el-table-column label="严重级别" width="100">
          <template #default="{ row }">
            <el-tag :type="getPriorityTag(row.priority)" effect="plain">
              {{ getPriorityName(row.priority) }}
            </el-tag>
          </template>
        </el-table-column>
        <!-- 接收人 -->
        <el-table-column label="接收人" width="150">
          <template #default="{ row }">
            <span>{{ getRecipients(row) }}</span>
          </template>
        </el-table-column>
        <!-- 类型 -->
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTopicTypeTag(row.topic_type)" size="small">
              {{ getTopicTypeName(row.topic_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <!-- 接收时间 -->
        <el-table-column label="接收时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.received_at || row.created_at) }}
          </template>
        </el-table-column>
        <!-- 内容 -->
        <el-table-column label="内容" min-width="250" show-overflow-tooltip>
          <template #default="{ row }">
            <span>{{ row.content }}</span>
          </template>
        </el-table-column>
        <!-- 操作 -->
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="handleView(row)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.limit"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadMessages"
        @current-change="loadMessages"
        class="pagination"
      />
    </el-card>

    <!-- 详情弹窗 - 参考 Cloudpods -->
    <el-dialog v-model="detailVisible" title="消息详情" width="700px">
      <div v-if="currentMessage" class="message-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="消息ID">{{ currentMessage.id }}</el-descriptions-item>
          <el-descriptions-item label="消息标题">{{ currentMessage.title }}</el-descriptions-item>
          <el-descriptions-item label="严重级别">
            <el-tag :type="getPriorityTag(currentMessage.priority)">
              {{ getPriorityName(currentMessage.priority) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="消息状态">
            <el-tag :type="currentMessage.status === 'ok' ? 'success' : 'warning'">
              {{ currentMessage.status === 'ok' ? '成功' : currentMessage.status }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="消息类型">
            <el-tag :type="getTopicTypeTag(currentMessage.topic_type)" size="small">
              {{ getTopicTypeName(currentMessage.topic_type) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="通知渠道">{{ currentMessage.contact_type || '站内信' }}</el-descriptions-item>
          <el-descriptions-item label="接收人" :span="2">{{ getRecipients(currentMessage) }}</el-descriptions-item>
          <el-descriptions-item label="接收时间">{{ formatTime(currentMessage.received_at || currentMessage.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="阅读状态">
            <el-tag :type="currentMessage.read_at ? 'success' : 'warning'">
              {{ currentMessage.read_at ? '已读' : '未读' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="消息内容" :span="2">
            <div style="white-space: pre-wrap;">{{ currentMessage.content }}</div>
          </el-descriptions-item>
          <!-- 接收详情 -->
          <el-descriptions-item label="接收详情" :span="2" v-if="currentMessage.receive_details">
            <el-table :data="currentMessage.receive_details" border size="small">
              <el-table-column prop="receiver_name" label="接收人" />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.status === 'sent_ok' ? 'success' : 'danger'" size="small">
                    {{ row.status === 'sent_ok' ? '发送成功' : row.status }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </el-descriptions-item>
        </el-descriptions>
      </div>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
        <el-button
          v-if="currentMessage && !currentMessage.read_at"
          type="success"
          @click="handleMarkRead(currentMessage)"
        >
          标记已读
        </el-button>
        <el-button
          v-if="currentMessage?.can_delete"
          type="danger"
          @click="handleDelete(currentMessage)"
        >
          删除
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Refresh, Check, Message, BellFilled, CircleCheckFilled, WarningFilled, Search, RefreshRight } from '@element-plus/icons-vue'
import { getMessages, getMessage, markRead, markAllRead, deleteMessage, getUnreadCount } from '@/api/message'

interface NotificationMessage {
  id: string | number
  title: string
  content: string
  contact_type?: string
  priority: string // important, normal, low
  status?: string // ok, pending, failed
  topic_type?: string // resource, system, security
  receive_details?: Array<{ receiver_name: string; status: string }>
  received_at?: string
  read_at?: string
  can_delete?: boolean
  can_update?: boolean
  created_at: string
}

const loading = ref(false)
const detailVisible = ref(false)
const currentMessage = ref<NotificationMessage | null>(null)
const selectedMessages = ref<NotificationMessage[]>([])
const unreadCount = ref(0)

const pagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

const filters = reactive({
  status: '',
  priority: '',
  topic_type: ''
})

const messages = ref<NotificationMessage[]>([])

// 计算属性
const readCount = computed(() => messages.value.filter(m => m.read_at).length)
const importantCount = computed(() => messages.value.filter(m => m.priority === 'important').length)

// 级别映射 - 参考 Cloudpods
const getPriorityName = (priority: string) => {
  const map: Record<string, string> = {
    important: '重要',
    normal: '普通',
    low: '低',
    critical: '致命'
  }
  return map[priority] || priority || '普通'
}

const getPriorityTag = (priority: string) => {
  const map: Record<string, any> = {
    important: 'warning',
    normal: 'info',
    low: 'info',
    critical: 'danger'
  }
  return map[priority] || 'info'
}

// 类型映射
const getTopicTypeName = (topicType: string) => {
  const map: Record<string, string> = {
    resource: '资源',
    system: '系统',
    security: '安全',
    task: '任务'
  }
  return map[topicType] || topicType || '系统'
}

const getTopicTypeTag = (topicType: string) => {
  const map: Record<string, any> = {
    resource: 'primary',
    system: 'info',
    security: 'danger',
    task: 'warning'
  }
  return map[topicType] || 'info'
}

// 获取接收人
const getRecipients = (msg: NotificationMessage) => {
  if (msg.receive_details && msg.receive_details.length > 0) {
    return msg.receive_details.map(r => r.receiver_name).join(', ')
  }
  return 'admin'
}

const formatTime = (time: string) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

const handleSelectionChange = (selection: NotificationMessage[]) => {
  selectedMessages.value = selection
}

const handleViewMode = () => {
  ElMessage.info('视图切换功能')
}

// 加载消息列表
const loadMessages = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit,
      details: 'true'
    }

    if (filters.status) {
      params.read = filters.status === 'read' ? 'true' : 'false'
    }
    if (filters.priority) {
      params.priority = filters.priority
    }
    if (filters.topic_type) {
      params.topic_type = filters.topic_type
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

// 加载未读数量
const loadUnreadCount = async () => {
  try {
    const res = await getUnreadCount()
    unreadCount.value = res.count || 0
  } catch (e: any) {
    console.error('加载未读数量失败', e)
  }
}

const handleView = async (row: NotificationMessage) => {
  try {
    // 获取详情
    const res = await getMessage(Number(row.id))
    currentMessage.value = res
    detailVisible.value = true

    // 自动标记已读
    if (!row.read_at) {
      await markRead(Number(row.id))
      row.read_at = new Date().toISOString()
      unreadCount.value--
    }
  } catch (e: any) {
    // 如果获取详情失败，直接显示当前数据
    currentMessage.value = row
    detailVisible.value = true
  }
}

const handleMarkRead = async (row: NotificationMessage) => {
  try {
    await markRead(Number(row.id))
    ElMessage.success('已标记为已读')
    row.read_at = new Date().toISOString()
    detailVisible.value = false
    unreadCount.value--
    loadMessages()
  } catch (e: any) {
    ElMessage.error(e.message || '标记失败')
  }
}

const handleMarkAllRead = async () => {
  try {
    await ElMessageBox.confirm('确定将所有消息标记为已读?', '确认', { type: 'warning' })
    await markAllRead()
    ElMessage.success('已全部标记为已读')
    unreadCount.value = 0
    loadMessages()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '操作失败')
    }
  }
}

const handleDelete = async (row: NotificationMessage) => {
  try {
    await ElMessageBox.confirm('确定删除此消息?', '确认', { type: 'warning' })
    await deleteMessage(Number(row.id))
    ElMessage.success('消息已删除')
    detailVisible.value = false
    loadMessages()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const handleRefresh = () => {
  loadMessages()
  loadUnreadCount()
  ElMessage.success('刷新成功')
}

const resetFilters = () => {
  filters.status = ''
  filters.priority = ''
  filters.topic_type = ''
  pagination.page = 1
  loadMessages()
}

onMounted(() => {
  loadUnreadCount()
  loadMessages()
})
</script>

<style scoped>
.inbox-page {
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

.toolbar {
  display: flex;
  gap: 8px;
}

/* 统计卡片样式 - 参考 Cloudpods */
.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-radius: 8px;
  background: #f5f7fa;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-card.total {
  background: linear-gradient(135deg, #f5f7fa 0%, #e5e7ea 100%);
  border-left: 4px solid #909399;
}

.stat-card.unread {
  background: linear-gradient(135deg, #fffaf0 0%, #ffe8d0 100%);
  border-left: 4px solid #e6a23c;
}

.stat-card.read {
  background: linear-gradient(135deg, #f0f9ff 0%, #d0e8ff 100%);
  border-left: 4px solid #409eff;
}

.stat-card.important {
  background: linear-gradient(135deg, #fff5f5 0%, #ffe0e0 100%);
  border-left: 4px solid #f56c6c;
}

.stat-icon {
  font-size: 28px;
  margin-right: 12px;
}

.stat-card.total .stat-icon {
  color: #909399;
}

.stat-card.unread .stat-icon {
  color: #e6a23c;
}

.stat-card.read .stat-icon {
  color: #409eff;
}

.stat-card.important .stat-icon {
  color: #f56c6c;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  line-height: 1.2;
}

.stat-label {
  font-size: 12px;
  color: #909399;
}

/* 筛选区域 */
.filter-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.unread-title {
  font-weight: bold;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.message-detail {
  padding: 10px 0;
}
</style>