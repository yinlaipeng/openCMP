<template>
  <div class="receivers-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">接收人管理</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建
          </el-button>
        </div>
      </template>

      <el-table :data="receivers" v-loading="loading" border stripe>
        <el-table-column prop="name" label="用户名" min-width="120" show-overflow-tooltip />
        <el-table-column prop="enabled" label="启用状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="phone" label="手机号码" width="140" />
        <el-table-column prop="email" label="邮箱" width="200" show-overflow-tooltip />
        <el-table-column label="通知渠道" width="200">
          <template #default="{ row }">
            <el-tag v-for="channel in row.notification_channels || []" :key="channel.id" size="small" style="margin-right: 4px;">
              {{ getChannelName(channel) }}
            </el-tag>
            <span v-if="!(row.notification_channels && row.notification_channels.length)">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="domain.name" label="所属域" width="120" />
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleViewDetails(row)">详情</el-button>
            <el-dropdown trigger="click">
              <el-button size="small">
                更多
                <el-icon><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleEnable(row)" :disabled="row.enabled">
                    启用
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleDisable(row)" :disabled="!row.enabled">
                    禁用
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleEdit(row)">
                    修改
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleDelete(row)" divided>
                    删除
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadReceivers"
        @current-change="loadReceivers"
        class="pagination"
      />
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="域" prop="domain_id">
          <el-select v-model="form.domain_id" placeholder="请选择域" style="width: 100%">
            <el-option v-for="domain in domains" :key="domain.id" :label="domain.name" :value="domain.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="用户" prop="user_id">
          <el-select v-model="form.user_id" placeholder="请选择用户" style="width: 100%" @change="handleUserChange">
            <el-option v-for="user in users" :key="user.id" :label="user.name" :value="user.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="手机号码" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号码" />
        </el-form-item>
        <el-form-item label="邮箱账号" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱账号" />
        </el-form-item>
        <el-form-item label="通知渠道">
          <el-checkbox-group v-model="selectedChannels">
            <el-checkbox v-for="channel in channels" :key="channel.id" :label="channel.id">
              {{ channel.name }} ({{ getChannelTypeName(channel.type) }})
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailsVisible" title="接收人详情" width="700px">
      <div v-if="currentReceiver" class="receiver-details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="用户名">{{ currentReceiver.name }}</el-descriptions-item>
          <el-descriptions-item label="启用状态">
            <el-tag :type="currentReceiver.enabled ? 'success' : 'info'">
              {{ currentReceiver.enabled ? '启用' : '禁用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="手机号码">{{ currentReceiver.phone || '-' }}</el-descriptions-item>
          <el-descriptions-item label="邮箱">{{ currentReceiver.email || '-' }}</el-descriptions-item>
          <el-descriptions-item label="所属域">{{ currentReceiver.domain?.name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDate(currentReceiver.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="通知渠道" :span="2">
            <el-tag v-for="channel in currentReceiver.notification_channels || []" :key="channel.id" size="small" style="margin-right: 4px;">
              {{ channel.name }} ({{ getChannelTypeName(channel.type) }})
            </el-tag>
            <span v-if="!(currentReceiver.notification_channels && currentReceiver.notification_channels.length)">-</span>
          </el-descriptions-item>
        </el-descriptions>
      </div>
      <template #footer>
        <el-button @click="detailsVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus, ArrowDown } from '@element-plus/icons-vue'
import {
  getReceivers,
  createReceiver,
  updateReceiver,
  deleteReceiver,
  enableReceiver,
  disableReceiver,
  getReceiverChannels,
  setReceiverChannels
} from '@/api/message'
import { getDomains } from '@/api/iam'
import { getUsers } from '@/api/iam'
import { getNotificationChannels } from '@/api/message'

interface Receiver {
  id: number
  name: string
  email: string
  phone: string
  domain_id: number
  enabled: boolean
  created_at: string
  domain?: {
    id: number
    name: string
  }
  notification_channels?: Array<{
    id: number
    name: string
    type: string
  }>
}

interface Domain {
  id: number
  name: string
}

interface User {
  id: number
  name: string
  email: string
  phone: string
}

interface NotificationChannel {
  id: number
  name: string
  type: string
}

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const detailsVisible = ref(false)
const dialogTitle = ref('')
const receivers = ref<Receiver[]>([])
const currentReceiver = ref<Receiver | null>(null)
const editingId = ref<number | null>(null)
const formRef = ref<FormInstance>()

const domains = ref<Domain[]>([])
const users = ref<User[]>([])
const channels = ref<NotificationChannel[]>([])
const selectedChannels = ref<number[]>([])

const form = reactive({
  name: '',
  email: '',
  phone: '',
  domain_id: 0,
  user_id: null as number | null,
  enabled: true
})

const rules = {
  domain_id: [{ required: true, message: '请选择域', trigger: 'change' }],
  email: [
    { required: true, message: '请输入邮箱账号', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号码', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ]
}

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const loadReceivers = async () => {
  loading.value = true
  try {
    const res = await getReceivers({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    receivers.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载接收人失败')
  } finally {
    loading.value = false
  }
}

const loadDomains = async () => {
  try {
    const res = await getDomains({ limit: 1000 })
    domains.value = res.items || []
  } catch (e: any) {
    console.error('加载域失败:', e)
  }
}

const loadUsers = async () => {
  try {
    const res = await getUsers({ limit: 1000 })
    users.value = res.items || []
  } catch (e: any) {
    console.error('加载用户失败:', e)
  }
}

const loadChannels = async () => {
  try {
    const res = await getNotificationChannels({ limit: 1000 })
    channels.value = res.items || []
  } catch (e: any) {
    console.error('加载通知渠道失败:', e)
  }
}

const handleCreate = () => {
  editingId.value = null
  dialogTitle.value = '新建接收人'
  resetForm()
  dialogVisible.value = true
}

const handleEdit = async (row: Receiver) => {
  editingId.value = row.id
  dialogTitle.value = '修改接收人'
  form.name = row.name
  form.email = row.email
  form.phone = row.phone
  form.domain_id = row.domain_id
  form.user_id = row.user_id || null
  form.enabled = row.enabled

  // Load associated channels
  try {
    const res = await getReceiverChannels(row.id)
    selectedChannels.value = (res.items || []).map((ch: any) => ch.id)
  } catch (e: any) {
    console.error('加载接收人通知渠道失败:', e)
    selectedChannels.value = []
  }

  dialogVisible.value = true
}

const handleViewDetails = async (row: Receiver) => {
  currentReceiver.value = row
  detailsVisible.value = true

  // Load full receiver details with channels if not already loaded
  if (!row.notification_channels || row.notification_channels.length === 0) {
    try {
      const res = await getReceiverChannels(row.id)
      if (currentReceiver.value) {
        currentReceiver.value.notification_channels = res.items || []
      }
    } catch (e: any) {
      console.error('加载接收人通知渠道失败:', e)
    }
  }
}

const handleEnable = async (row: Receiver) => {
  try {
    await enableReceiver(row.id)
    ElMessage.success('启用成功')
    loadReceivers()
  } catch (e: any) {
    ElMessage.error(e.message || '启用失败')
  }
}

const handleDisable = async (row: Receiver) => {
  try {
    await disableReceiver(row.id)
    ElMessage.success('禁用成功')
    loadReceivers()
  } catch (e: any) {
    ElMessage.error(e.message || '禁用失败')
  }
}

const handleDelete = async (row: Receiver) => {
  try {
    await ElMessageBox.confirm(`确定要删除接收人 "${row.name}" 吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await deleteReceiver(row.id)
    ElMessage.success('删除成功')
    loadReceivers()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const handleUserChange = (userId: number) => {
  const user = users.value.find(u => u.id === userId)
  if (user) {
    form.email = user.email || ''
    form.phone = user.phone || ''
    form.name = user.name
    // Find domain ID from user
    // Since we don't have a direct API to get domain by user, we'll need to get it from somewhere else
    // For now, we'll keep the selected domain as is
  }
}

const resetForm = () => {
  form.name = ''
  form.email = ''
  form.phone = ''
  form.domain_id = 0
  form.user_id = null
  form.enabled = true
  selectedChannels.value = []
}

const handleSubmit = async () => {
  await formRef.value?.validate()
  submitting.value = true
  try {
    const data = {
      name: form.name,
      email: form.email,
      phone: form.phone,
      domain_id: form.domain_id,
      user_id: form.user_id,
      enabled: form.enabled
    }

    if (editingId.value) {
      await updateReceiver(editingId.value, data)
      ElMessage.success('更新成功')

      // Update notification channels if they changed
      await setReceiverChannels(editingId.value, { channel_ids: selectedChannels.value })
    } else {
      const res = await createReceiver(data)
      const newId = res.id

      // Set notification channels for the new receiver
      await setReceiverChannels(newId, { channel_ids: selectedChannels.value })

      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    loadReceivers()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

const getChannelName = (channel: any) => {
  return channel.name || `渠道${channel.id}`
}

const getChannelTypeName = (type: string) => {
  const typeMap: Record<string, string> = {
    email: '邮件',
    sms: '短信',
    webhook: 'Webhook',
    dingtalk: '钉钉',
    wechat: '企业微信',
    feishu: '飞书',
    lark: 'Lark'
  }
  return typeMap[type] || type
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  try {
    return new Date(dateStr).toLocaleString('zh-CN')
  } catch {
    return dateStr
  }
}

onMounted(async () => {
  await Promise.all([
    loadReceivers(),
    loadDomains(),
    loadUsers(),
    loadChannels()
  ])
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.title {
  font-size: 16px;
  font-weight: bold;
}
.pagination {
  margin-top: 16px;
  justify-content: flex-end;
}
.receiver-details {
  max-height: 500px;
  overflow-y: auto;
}
</style>
