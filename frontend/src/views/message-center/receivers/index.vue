<template>
  <div class="receivers-container">
    <div class="page-header">
      <h2>接收人管理</h2>
      <div class="toolbar">
        <el-button @click="loadReceivers" :icon="Refresh" circle title="刷新" />
        <el-button type="primary" @click="handleCreate">
          新建
        </el-button>
        <el-button :disabled="selectedRows.length === 0" @click="handleBatchDelete">
          删除
        </el-button>
        <el-button @click="showSettings" :icon="Setting" circle title="设置" />
      </div>
    </div>
    <el-card class="filter-card">
      <el-form :inline="true" @submit.prevent="loadReceivers">
        <el-form-item>
          <el-select v-model="searchField" placeholder="搜索属性" style="width: 120px">
            <el-option label="用户名" value="name" />
            <el-option label="手机号" value="phone" />
            <el-option label="邮箱" value="email" />
            <el-option label="创建时间" value="created_at" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input
            v-model="searchKeyword"
            placeholder="默认为用户名搜索"
            style="width: 300px"
            clearable
            @keyup.enter="loadReceivers"
          >
            <template #suffix>
              <el-icon class="search-icon" @click="loadReceivers"><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
    </el-card>
    <el-card>
      <el-table
        :data="receivers"
        v-loading="loading"
        border
        stripe
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" />
        <el-table-column label="用户名" min-width="120" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.name }}
          </template>
        </el-table-column>
        <el-table-column label="启用状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ row.enabled ? '已启用' : '已禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="手机号" width="140">
          <template #default="{ row }">
            {{ row.phone || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="邮箱" width="200" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.email || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="通知渠道" width="150">
          <template #default="{ row }">
            <el-tag v-for="channel in row.notification_channels || []" :key="channel.id" size="small" style="margin-right: 4px;">
              {{ getChannelTypeName(channel.type) }}
            </el-tag>
            <span v-if="!(row.notification_channels && row.notification_channels.length)">-</span>
          </template>
        </el-table-column>
        <el-table-column label="所属域" width="120">
          <template #default="{ row }">
            {{ row.domain?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-dropdown trigger="click">
              <el-button size="small" type="primary" link>
                更多
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleEnable(row)" :disabled="row.enabled">
                    启用
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleDisable(row)" :disabled="!row.enabled">
                    禁用
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
        class="pagination"
        @size-change="loadReceivers"
        @current-change="loadReceivers"
      />
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px" label-position="left">
        <el-form-item label="域">
          <el-select v-model="form.domain_id" placeholder="请选择域" style="width: 300px" clearable>
            <el-option v-for="domain in domains" :key="domain.id" :label="domain.name" :value="domain.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="用户" prop="user_id">
          <el-select v-model="form.user_id" placeholder="请选择用户" style="width: 300px" @change="handleUserChange">
            <el-option v-for="user in users" :key="user.id" :label="user.name" :value="user.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <div class="mobile-input-group">
            <el-select v-model="form.country_code" style="width: 150px">
              <el-option label="中国大陆(+86)" value="+86" />
              <el-option label="香港(+852)" value="+852" />
              <el-option label="台湾(+886)" value="+886" />
              <el-option label="美国(+1)" value="+1" />
              <el-option label="日本(+81)" value="+81" />
            </el-select>
            <el-input v-model="form.phone_number" placeholder="手机号码" style="width: 150px" />
          </div>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="邮箱地址" style="width: 300px" />
        </el-form-item>
        <el-form-item>
          <template #label>
            <span>
              通知渠道
              <el-tooltip content="选择接收人可接收消息的渠道" placement="top">
                <el-icon style="margin-left: 4px; cursor: help;"><InfoFilled /></el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-checkbox-group v-model="selectedChannels">
            <el-checkbox label="webconsole" disabled>站内信（默认）</el-checkbox>
            <el-checkbox v-for="channel in channels" :key="channel.id" :label="channel.id">
              {{ channel.name }} ({{ getChannelTypeName(channel.type) }})
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
          <el-button @click="dialogVisible = false">取消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Refresh, Setting, Search, InfoFilled } from '@element-plus/icons-vue'
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
  user_id?: number
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
const dialogTitle = ref('')
const receivers = ref<Receiver[]>([])
const editingId = ref<number | null>(null)
const formRef = ref<FormInstance>()
const selectedRows = ref<Receiver[]>([])

const searchField = ref('name')
const searchKeyword = ref('')

const domains = ref<Domain[]>([])
const users = ref<User[]>([])
const channels = ref<NotificationChannel[]>([])
const selectedChannels = ref<string[]>(['webconsole'])

const form = reactive({
  name: '',
  email: '',
  phone: '',
  phone_number: '',
  country_code: '+86',
  domain_id: null as number | null,
  user_id: null as number | null
})

const rules = {
  user_id: [{ required: true, message: '请选择用户', trigger: 'change' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  phone_number: [
    { required: true, message: '请输入手机号', trigger: 'blur' }
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
      page_size: pagination.pageSize,
      keyword: searchKeyword.value || undefined,
      search_field: searchField.value
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
    const defaultDomain = domains.value.find(d => d.name === 'Default')
    if (defaultDomain && !form.domain_id) {
      form.domain_id = defaultDomain.id
    }
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

const handleSelectionChange = (rows: Receiver[]) => {
  selectedRows.value = rows
}

const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请选择要删除的接收人')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRows.value.length} 个接收人吗？`,
      '删除确认',
      { type: 'warning' }
    )
    for (const row of selectedRows.value) {
      await deleteReceiver(row.id)
    }
    ElMessage.success(`已删除 ${selectedRows.value.length} 个接收人`)
    loadReceivers()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const showSettings = () => {
  ElMessage.info('设置功能开发中')
}

const handleCreate = async () => {
  editingId.value = null
  dialogTitle.value = '新建接收人'
  resetForm()
  dialogVisible.value = true

  if (domains.value.length === 0) await loadDomains()
  if (users.value.length === 0) await loadUsers()
  if (channels.value.length === 0) await loadChannels()

  if (users.value.length > 0) {
    form.user_id = users.value[0].id
    handleUserChange(form.user_id)
  }
}

const handleEdit = async (row: Receiver) => {
  editingId.value = row.id
  dialogTitle.value = '编辑接收人'
  form.name = row.name
  form.email = row.email
  form.phone = row.phone
  form.domain_id = row.domain_id
  form.user_id = row.user_id || null

  if (row.phone) {
    const match = row.phone.match(/^(\+\d{1,3})\s*(\d+)$/)
    if (match) {
      form.country_code = match[1]
      form.phone_number = match[2]
    } else {
      form.phone_number = row.phone.replace(/\s/g, '')
    }
  }

  try {
    const res = await getReceiverChannels(row.id)
    const channelIds = (res.items || []).map((ch: any) => ch.id)
    selectedChannels.value = ['webconsole', ...channelIds.map(String)]
  } catch (e: any) {
    console.error('加载接收人通知渠道失败:', e)
    selectedChannels.value = ['webconsole']
  }

  dialogVisible.value = true
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
    form.phone_number = user.phone ? user.phone.replace(/\+\d{1,3}\s*/g, '') : ''
    form.name = user.name
  }
}

const resetForm = () => {
  form.name = ''
  form.email = ''
  form.phone = ''
  form.phone_number = ''
  form.country_code = '+86'
  form.domain_id = domains.value.find(d => d.name === 'Default')?.id || null
  form.user_id = null
  selectedChannels.value = ['webconsole']
}

const handleSubmit = async () => {
  await formRef.value?.validate()
  submitting.value = true
  try {
    const fullPhone = `${form.country_code} ${form.phone_number}`

    const data = {
      name: form.name,
      email: form.email,
      phone: fullPhone,
      domain_id: form.domain_id,
      user_id: form.user_id
    }

    if (editingId.value) {
      await updateReceiver(editingId.value, data)
      ElMessage.success('更新成功')

      const channelIds = selectedChannels.value.filter(c => c !== 'webconsole').map(Number)
      await setReceiverChannels(editingId.value, { channel_ids: channelIds })
    } else {
      const res = await createReceiver(data)
      const newId = res.id

      const channelIds = selectedChannels.value.filter(c => c !== 'webconsole').map(Number)
      await setReceiverChannels(newId, { channel_ids: channelIds })

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

const getChannelTypeName = (type: string) => {
  const typeMap: Record<string, string> = {
    email: '邮件',
    sms: '短信',
    webhook: 'Webhook',
    dingtalk: '钉钉',
    workwx: '企业微信',
    wechat: '企业微信',
    feishu: '飞书',
    lark: '飞书'
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
.receivers-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.page-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.toolbar {
  display: flex;
  gap: 8px;
}

.filter-card {
  margin-bottom: 16px;
}

.search-icon {
  cursor: pointer;
}

.pagination {
  margin-top: 16px;
  justify-content: flex-end;
}

.mobile-input-group {
  display: flex;
  gap: 8px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>