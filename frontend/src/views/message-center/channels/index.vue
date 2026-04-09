<template>
  <div class="channels-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">通知渠道</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建渠道
          </el-button>
        </div>
      </template>

      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="类型">
          <el-select v-model="filterForm.type" placeholder="全部" clearable style="width: 120px">
            <el-option label="邮件" value="email" />
            <el-option label="短信" value="sms" />
            <el-option label="钉钉" value="dingtalk" />
            <el-option label="飞书" value="feishu" />
            <el-option label="企业微信" value="wechat" />
            <el-option label="Lark" value="lark" />
            <el-option label="Webhook" value="webhook" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadData">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="channels" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag>{{ typeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="scope" label="所属范围" width="120">
          <template #default="{ row }">
            <el-tag type="primary">系统</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="enabled" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="primary" @click="handleTest(row)">测试</el-button>
            <el-button size="small" :type="row.enabled ? 'warning' : 'success'" @click="handleToggle(row)">
              {{ row.enabled ? '禁用' : '启用' }}
            </el-button>
            <el-popconfirm title="确定删除该渠道吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        class="pagination"
        @size-change="loadData"
        @current-change="loadData"
      />
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="800px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item label="归属" prop="scope">
          <el-radio-group v-model="form.scope" :disabled="!!editingId">
            <el-radio label="system">系统</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入渠道名称" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择类型" style="width: 100%" @change="handleTypeChange">
            <el-option label="邮件" value="email" />
            <el-option label="短信" value="sms" />
            <el-option label="钉钉" value="dingtalk" />
            <el-option label="飞书" value="feishu" />
            <el-option label="企业微信" value="wechat" />
            <el-option label="Lark" value="lark" />
          </el-select>
        </el-form-item>

        <!-- 邮箱配置 -->
        <template v-if="form.type === 'email'">
          <el-form-item label="SMTP服务器" prop="config.smtp_host">
            <el-input v-model="form.config.smtp_host" placeholder="例如：smtp.gmail.com" />
          </el-form-item>
          <el-form-item label="SSL">
            <el-switch v-model="form.config.use_ssl" />
          </el-form-item>
          <el-form-item label="端口" prop="config.smtp_port">
            <el-input-number v-model="form.config.smtp_port" :min="1" :max="65535" style="width: 100%" />
          </el-form-item>
          <el-form-item label="用户名" prop="config.smtp_user">
            <el-input v-model="form.config.smtp_user" placeholder="邮箱地址" />
          </el-form-item>
          <el-form-item label="密码" prop="config.smtp_password">
            <el-input v-model="form.config.smtp_password" type="password" placeholder="邮箱密码或授权码" show-password />
          </el-form-item>
          <el-form-item label="发件人邮箱" prop="config.from_address">
            <el-input v-model="form.config.from_address" placeholder="一般与用户名邮箱地址保持一致" />
          </el-form-item>
        </template>

        <!-- 短信配置 -->
        <template v-if="form.type === 'sms'">
          <el-form-item label="短信服务商" prop="config.provider">
            <el-select v-model="form.config.provider" placeholder="请选择短信服务商" style="width: 100%">
              <el-option label="阿里云" value="aliyun" />
              <el-option label="华为云" value="huawei" />
            </el-select>
          </el-form-item>
          <el-form-item label="Access Key ID" prop="config.access_key_id">
            <el-input v-model="form.config.access_key_id" placeholder="请输入Access Key ID" />
          </el-form-item>
          <el-form-item label="Access Key Secret" prop="config.access_key_secret">
            <el-input v-model="form.config.access_key_secret" type="password" placeholder="请输入Access Key Secret" show-password />
          </el-form-item>
          <el-form-item label="签名" prop="config.signature">
            <el-input v-model="form.config.signature" placeholder="请输入短信签名" />
          </el-form-item>

          <!-- 短信模板部分 -->
          <el-collapse>
            <el-collapse-item title="短信模板配置">
              <el-tabs type="border-card">
                <el-tab-pane label="国内短信模板">
                  <el-form-item label="验证码模板ID">
                    <el-input v-model="form.config.domestic_templates.verify_code" placeholder="请输入验证码模板ID" />
                  </el-form-item>
                  <el-form-item label="告警模板ID">
                    <el-input v-model="form.config.domestic_templates.alert" placeholder="请输入告警模板ID" />
                  </el-form-item>
                  <el-form-item label="异常登录模板ID">
                    <el-input v-model="form.config.domestic_templates.abnormal_login" placeholder="请输入异常登录模板ID" />
                  </el-form-item>
                </el-tab-pane>
                <el-tab-pane label="国际/港澳台模板">
                  <el-form-item label="验证码模板ID">
                    <el-input v-model="form.config.intl_templates.verify_code" placeholder="请输入验证码模板ID" />
                  </el-form-item>
                  <el-form-item label="告警模板ID">
                    <el-input v-model="form.config.intl_templates.alert" placeholder="请输入告警模板ID" />
                  </el-form-item>
                  <el-form-item label="异常登录模板ID">
                    <el-input v-model="form.config.intl_templates.abnormal_login" placeholder="请输入异常登录模板ID" />
                  </el-form-item>
                </el-tab-pane>
              </el-tabs>
            </el-collapse-item>
          </el-collapse>
        </template>

        <!-- 钉钉/飞书/企业微信/Lark 配置 -->
        <template v-if="['dingtalk', 'feishu', 'wechat', 'lark'].includes(form.type)">
          <el-form-item label="Webhook地址" prop="config.webhook_url">
            <el-input v-model="form.config.webhook_url" type="textarea" :rows="4" placeholder="请输入Webhook地址" />
          </el-form-item>
          <el-form-item label="密钥（如有）" prop="config.secret">
            <el-input v-model="form.config.secret" type="password" placeholder="请输入密钥" show-password />
          </el-form-item>
        </template>

        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" />
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getNotificationChannels,
  createNotificationChannel,
  updateNotificationChannel,
  deleteNotificationChannel,
  enableNotificationChannel,
  disableNotificationChannel,
  testNotificationChannel
} from '@/api/message'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('新建渠道')
const channels = ref<any[]>([])
const editingId = ref<number | null>(null)
const formRef = ref<FormInstance>()

const filterForm = reactive({ type: '' })
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

// Initialize form with proper nested configuration structure
const initializeFormConfig = () => {
  return {
    smtp_host: '',
    smtp_port: 465,
    smtp_user: '',
    smtp_password: '',
    from_address: '',
    from_name: '',
    use_ssl: true,
    use_tls: true,
    provider: '',
    access_key_id: '',
    access_key_secret: '',
    signature: '',
    domestic_templates: {
      verify_code: '',
      alert: '',
      abnormal_login: ''
    },
    intl_templates: {
      verify_code: '',
      alert: '',
      abnormal_login: ''
    },
    webhook_url: '',
    secret: ''
  }
}

const form = reactive({
  name: '',
  type: '',
  scope: 'system',
  description: '',
  enabled: true,
  config: initializeFormConfig()
})

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  scope: [{ required: true, message: '请选择归属', trigger: 'change' }],
  // Specific validations for different config fields
  'config.smtp_host': [
    { required: true, message: '请输入SMTP服务器', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9.-]+$/, message: '请输入有效的服务器地址', trigger: 'blur' }
  ],
  'config.smtp_user': [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  'config.from_address': [
    { required: true, message: '请输入发件人邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  'config.smtp_port': [
    { required: true, message: '请输入端口', trigger: 'blur' },
    { type: 'number', message: '请输入有效端口号', trigger: 'blur' }
  ],
  'config.smtp_password': [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ],
  'config.provider': [
    { required: true, message: '请选择短信服务商', trigger: 'change' }
  ],
  'config.access_key_id': [
    { required: true, message: '请输入Access Key ID', trigger: 'blur' }
  ],
  'config.access_key_secret': [
    { required: true, message: '请输入Access Key Secret', trigger: 'blur' }
  ],
  'config.signature': [
    { required: true, message: '请输入签名', trigger: 'blur' }
  ],
  'config.webhook_url': [
    { required: true, message: '请输入Webhook地址', trigger: 'blur' },
    { type: 'url', message: '请输入有效的URL', trigger: 'blur' }
  ]
}

const typeLabel = (type: string) => {
  const map: Record<string, string> = {
    email: '邮件',
    sms: '短信',
    wechat: '企业微信',
    dingtalk: '钉钉',
    feishu: '飞书',
    lark: 'Lark',
    webhook: 'Webhook'
  }
  return map[type] || type
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getNotificationChannels({
      type: filterForm.type || undefined,
      page: pagination.page,
      page_size: pagination.pageSize
    })
    channels.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

const resetFilter = () => {
  filterForm.type = ''
  pagination.page = 1
  loadData()
}

const handleCreate = () => {
  editingId.value = null
  dialogTitle.value = '新建渠道'
  form.name = ''
  form.type = ''
  form.scope = 'system'
  form.description = ''
  form.enabled = true
  form.config = initializeFormConfig()
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  editingId.value = row.id
  dialogTitle.value = '编辑渠道'
  form.name = row.name
  form.type = row.type
  form.scope = row.scope || 'system'  // Default to system if not set
  form.description = row.description || ''
  form.enabled = row.enabled
  // Load config, merging with defaults to ensure all fields exist
  form.config = { ...initializeFormConfig(), ...row.config }
  dialogVisible.value = true
}

const handleTypeChange = () => {
  // Reset config when changing type to ensure clean defaults
  form.config = initializeFormConfig()
}

const handleSubmit = async () => {
  await formRef.value?.validate()
  submitting.value = true
  try {
    // Prepare data for submission
    const data = {
      name: form.name,
      type: form.type,
      description: form.description,
      enabled: form.enabled,
      config: form.config
    }

    if (editingId.value) {
      await updateNotificationChannel(editingId.value, data)
      ElMessage.success('更新成功')
    } else {
      await createNotificationChannel(data)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadData()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

const handleTest = async (row: any) => {
  try {
    await testNotificationChannel(row.id)
    ElMessage.success('测试消息已发送')
  } catch (e: any) {
    ElMessage.error(e.message || '测试失败')
  }
}

const handleToggle = async (row: any) => {
  try {
    if (row.enabled) {
      await disableNotificationChannel(row.id)
      ElMessage.success('已禁用')
    } else {
      await enableNotificationChannel(row.id)
      ElMessage.success('已启用')
    }
    loadData()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

const handleDelete = async (row: any) => {
  try {
    await deleteNotificationChannel(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (e: any) {
    ElMessage.error(e.message || '删除失败')
  }
}

onMounted(loadData)
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
.title { font-size: 16px; font-weight: bold; }
.filter-form { margin-bottom: 16px; }
.pagination { margin-top: 16px; justify-content: flex-end; }
</style>
