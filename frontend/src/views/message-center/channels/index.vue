<template>
  <div class="channels-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">通知渠道设置</span>
          <!-- 工具栏按钮 - 参考 Cloudpods -->
          <div class="toolbar">
            <el-button @click="loadData" :icon="Refresh" circle title="刷新" />
            <el-button type="primary" @click="handleCreate">
              新 建
            </el-button>
            <el-button :disabled="selectedRows.length === 0" @click="handleBatchDelete">
              删 除
            </el-button>
            <el-button @click="showSettings" :icon="Setting" circle title="设置" />
          </div>
        </div>
      </template>

      <!-- 搜索栏 - 参考 Cloudpods 设计 -->
      <div class="search-bar">
        <div class="search-box">
          <el-select v-model="searchField" placeholder="搜索属性" style="width: 100px" class="search-field-select">
            <el-option label="名称" value="name" />
            <el-option label="创建时间" value="created_at" />
          </el-select>
          <el-input
            v-model="searchKeyword"
            placeholder="默认为名称搜索"
            style="width: 300px"
            clearable
            @keyup.enter="loadData"
          >
            <template #suffix>
              <el-icon class="search-icon" @click="loadData"><Search /></el-icon>
            </template>
          </el-input>
          <span class="search-hint">默认为名称搜索，自动匹配ID搜索项</span>
        </div>
        <div class="filter-box">
          <el-select v-model="filterForm.type" placeholder="类型" clearable style="width: 120px" @change="loadData">
            <el-option label="全部" value="" />
            <el-option label="邮件" value="email" />
            <el-option label="短信" value="sms" />
            <el-option label="钉钉" value="dingtalk" />
            <el-option label="飞书" value="feishu" />
            <el-option label="企业微信" value="workwx" />
            <el-option label="Lark" value="lark" />
          </el-select>
        </div>
      </div>

      <!-- 表格 - 参考 Cloudpods 表头设计 -->
      <el-table
        :data="channels"
        v-loading="loading"
        border
        stripe
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <!-- 选择列 -->
        <el-table-column type="selection" width="50" />
        <!-- 名称 -->
        <el-table-column label="名称" min-width="150" show-overflow-tooltip>
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">
              {{ row.name }}
            </el-button>
          </template>
        </el-table-column>
        <!-- 类型 -->
        <el-table-column label="类型" width="120">
          <template #default="{ row }">
            <el-tag>{{ typeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <!-- 所属范围 -->
        <el-table-column label="所属范围" width="120">
          <template #default="{ row }">
            <el-tag type="primary">系统</el-tag>
          </template>
        </el-table-column>
        <!-- 操作 -->
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="primary" link @click="handleTestConnection(row)">连接测试</el-button>
            <el-popconfirm title="确定删除该渠道吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button size="small" type="primary" link>删除</el-button>
              </template>
            </el-popconfirm>
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
        class="pagination"
        @size-change="loadData"
        @current-change="loadData"
      />
    </el-card>

    <!-- 创建/编辑对话框 - 参考 Cloudpods 设计 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="700px"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px" label-position="left">
        <!-- 基础信息表单 -->
        <el-form-item label="归属" prop="scope">
          <el-radio-group v-model="form.scope" :disabled="!!editingId">
            <el-radio-button label="system">系统</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input
            v-model="form.name"
            placeholder="2-128字符，包含字母、数字、连字符'-'，以字母开头且不能以'-'结尾"
          />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-radio-group v-model="form.type" @change="handleTypeChange">
            <el-radio-button label="email">邮件</el-radio-button>
            <el-radio-button label="sms">短信</el-radio-button>
            <el-radio-button label="dingtalk">钉钉</el-radio-button>
            <el-radio-button label="feishu">飞书</el-radio-button>
            <el-radio-button label="workwx">企业微信</el-radio-button>
          </el-radio-group>
        </el-form-item>

        <!-- 邮箱配置 - 参考 Cloudpods -->
        <template v-if="form.type === 'email'">
          <el-divider content-position="left">邮件服务器配置</el-divider>
          <el-form-item label="SMTP服务器" prop="config.smtp_host">
            <el-input v-model="form.config.smtp_host" placeholder="例如：smtp.gmail.com" />
          </el-form-item>
          <el-form-item label="SSL">
            <el-radio-group v-model="form.config.use_ssl">
              <el-radio-button :label="true">启用</el-radio-button>
              <el-radio-button :label="false">禁用</el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="端口" prop="config.smtp_port">
            <el-input v-model="form.config.smtp_port" placeholder="例如：465（SSL）或 25（无SSL）" />
          </el-form-item>
          <el-form-item label="用户名" prop="config.smtp_user">
            <el-input v-model="form.config.smtp_user" placeholder="邮箱地址" />
          </el-form-item>
          <el-form-item label="密码" prop="config.smtp_password">
            <el-input v-model="form.config.smtp_password" type="password" placeholder="邮箱密码或授权码" show-password />
          </el-form-item>
          <el-form-item label="发件人邮箱" prop="config.from_address">
            <el-input v-model="form.config.from_address" placeholder="一般与用户名邮箱地址保持一致" />
            <div class="form-help">发件人邮箱地址一般与用户名邮箱地址一致，如果邮箱有其他配置，请按实际情况填写</div>
          </el-form-item>
        </template>

        <!-- 短信配置 - 参考 Cloudpods -->
        <template v-if="form.type === 'sms'">
          <el-divider content-position="left">短信服务配置</el-divider>
          <el-form-item label="短信服务商">
            <el-radio-group v-model="form.config.provider">
              <el-radio-button label="aliyun">阿里云</el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="Access Key ID" prop="config.access_key_id">
            <el-input v-model="form.config.access_key_id" placeholder="例如：LTAI3ToIkpuRetxx" />
            <div class="form-help">获取位置：进入控制台 -> 点击用户头像 -> accesskeys</div>
          </el-form-item>
          <el-form-item label="Access Key Secret" prop="config.access_key_secret">
            <el-input v-model="form.config.access_key_secret" type="password" placeholder="例如：sdRnlxBrWZI64E5BK8R2dNyw3z4W" show-password />
            <div class="form-help">获取位置：进入控制台 -> 点击用户头像 -> accesskeys</div>
          </el-form-item>
          <el-form-item label="签名" prop="config.signature">
            <el-input v-model="form.config.signature" placeholder="一般为公司名称简称" />
            <div class="form-help">获取位置：进入控制台 -> 产品与服务 -> 短信服务 -> 签名管理</div>
          </el-form-item>

          <el-divider content-position="left">短信模板配置</el-divider>
          <el-form-item label="验证码模板" prop="config.verify_code_template">
            <el-input v-model="form.config.verify_code_template" placeholder="请输入模板CODE，例如：SMS_123456789" />
            <div class="form-help">模板类型：验证码，内容：您的验证码为${code}</div>
          </el-form-item>
          <el-form-item label="告警模板">
            <el-input v-model="form.config.alert_template" placeholder="请输入模板CODE，例如：SMS_123456789" />
            <div class="form-help">模板类型：短信通知，内容：${type}类型的${alert_name}触发告警，请及时登录平台查看</div>
          </el-form-item>
          <el-form-item label="异常登录模板">
            <el-input v-model="form.config.abnormal_login_template" placeholder="请输入模板CODE，例如：SMS_123456789" />
            <div class="form-help">模板类型：短信通知，内容：域${domain}的账号${user}因异常登录被锁定，请核实情况</div>
          </el-form-item>
        </template>

        <!-- 钉钉配置 - 参考 Cloudpods -->
        <template v-if="form.type === 'dingtalk'">
          <el-divider content-position="left">钉钉应用配置</el-divider>
          <el-form-item label="AgentId" prop="config.agent_id">
            <el-input v-model="form.config.agent_id" placeholder="例如：217947123" />
            <div class="form-help">
              获取位置：以管理员身份登录钉钉开放平台 -> 应用开发 -> 创建应用 -> 创建后应用详情中获取应用凭证
              <el-link type="primary" href="https://open.dingtalk.com/" target="_blank" class="help-link">
                <el-icon><Link /></el-icon>钉钉开放平台
              </el-link>
            </div>
          </el-form-item>
          <el-form-item label="AppKey" prop="config.app_key">
            <el-input v-model="form.config.app_key" placeholder="例如：dingo9s3gzs5123456" />
            <div class="form-help">获取位置：同 AgentId 获取位置</div>
          </el-form-item>
          <el-form-item label="AppSecret" prop="config.app_secret">
            <el-input v-model="form.config.app_secret" type="password" placeholder="例如：adjfkasjdfkjssadfzJPTZnyP6zxn23kiO..." show-password />
            <div class="form-help">获取位置：同 AgentId 获取位置</div>
          </el-form-item>
        </template>

        <!-- 飞书配置 - 参考 Cloudpods -->
        <template v-if="form.type === 'feishu'">
          <el-divider content-position="left">飞书应用配置</el-divider>
          <el-form-item label="AppID" prop="config.app_id">
            <el-input v-model="form.config.app_id" placeholder="例如：cli_9adbc25c4cb2020d" />
            <div class="form-help">
              获取位置：飞书开放平台 -> 开发者后台 -> 创建应用 -> 应用详情-凭证与基础信息中获取AppID和AppSecret
              <el-link type="primary" href="https://open.larksuite.com/" target="_blank" class="help-link">
                <el-icon><Link /></el-icon>飞书开放平台
              </el-link>
            </div>
          </el-form-item>
          <el-form-item label="AppSecret" prop="config.app_secret">
            <el-input v-model="form.config.app_secret" type="password" placeholder="例如：ccyaskdfjLKjkJN5jngseYBEnnBkbae" show-password />
            <div class="form-help">获取位置：同 AppID 获取位置</div>
          </el-form-item>
        </template>

        <!-- 企业微信配置 - 参考 Cloudpods -->
        <template v-if="form.type === 'workwx'">
          <el-divider content-position="left">企业微信应用配置</el-divider>
          <el-form-item label="CorpId" prop="config.corp_id">
            <el-input v-model="form.config.corp_id" placeholder="例如：ww2c41e47d2d3b13cb" />
            <div class="form-help">
              获取位置：以管理员身份登录企业微信管理后台 -> 我的企业 -> 获取企业ID（CorpId）
              <el-link type="primary" href="https://work.weixin.qq.com/" target="_blank" class="help-link">
                <el-icon><Link /></el-icon>企业微信管理后台
              </el-link>
            </div>
          </el-form-item>
          <el-form-item label="AgentId" prop="config.agent_id">
            <el-input v-model="form.config.agent_id" placeholder="例如：1000002" />
            <div class="form-help">获取位置：企业微信管理后台 -> 应用管理 -> 创建应用 -> 应用详情中获取应用AgentID</div>
          </el-form-item>
          <el-form-item label="Secret" prop="config.secret">
            <el-input v-model="form.config.secret" type="password" placeholder="例如：ZgyVyfr2Mvd0zzy6bE5prfKX25k4Wrgn4-1DSVDYXVo" show-password />
            <div class="form-help">获取位置：企业微信管理后台 -> 应用管理 -> 创建应用 -> 应用详情中获取应用凭证</div>
          </el-form-item>
        </template>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="handleTestConnectionInDialog" :loading="testingConnection">
            连接测试
          </el-button>
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Refresh, Setting, Search, Link } from '@element-plus/icons-vue'
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
const testingConnection = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('新建渠道')
const channels = ref<any[]>([])
const editingId = ref<number | null>(null)
const formRef = ref<FormInstance>()
const selectedRows = ref<any[]>([])

// 搜索相关
const searchField = ref('name')
const searchKeyword = ref('')
const filterForm = reactive({ type: '' })
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

// Initialize form config based on type
const initializeFormConfig = () => {
  return {
    // Email config
    smtp_host: '',
    smtp_port: '465',
    smtp_user: '',
    smtp_password: '',
    from_address: '',
    use_ssl: true,

    // SMS config
    provider: 'aliyun',
    access_key_id: '',
    access_key_secret: '',
    signature: '',
    verify_code_template: '',
    alert_template: '',
    abnormal_login_template: '',

    // DingTalk config
    agent_id: '',
    app_key: '',
    app_secret: '',

    // Feishu/Lark config
    app_id: '',
    app_secret: '',

    // WeCom config
    corp_id: '',
    secret: ''
  }
}

const form = reactive({
  name: '',
  type: '',
  scope: 'system',
  config: initializeFormConfig()
})

const rules = {
  name: [
    { required: true, message: '请输入名称', trigger: 'blur' },
    { min: 2, max: 128, message: '名称长度2-128字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$/, message: '以字母开头，只能包含字母、数字、连字符，不能以连字符结尾', trigger: 'blur' }
  ],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  scope: [{ required: true, message: '请选择归属', trigger: 'change' }],

  // Email config rules
  'config.smtp_host': [{ required: true, message: '请输入SMTP服务器', trigger: 'blur' }],
  'config.smtp_port': [{ required: true, message: '请输入端口', trigger: 'blur' }],
  'config.smtp_user': [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  'config.smtp_password': [{ required: true, message: '请输入密码', trigger: 'blur' }],
  'config.from_address': [{ required: true, message: '请输入发件人邮箱', trigger: 'blur' }],

  // SMS config rules
  'config.access_key_id': [{ required: true, message: '请输入Access Key ID', trigger: 'blur' }],
  'config.access_key_secret': [{ required: true, message: '请输入Access Key Secret', trigger: 'blur' }],
  'config.signature': [{ required: true, message: '请输入签名', trigger: 'blur' }],
  'config.verify_code_template': [{ required: true, message: '请输入验证码模板', trigger: 'blur' }],

  // DingTalk config rules
  'config.agent_id': [{ required: true, message: '请输入AgentId', trigger: 'blur' }],
  'config.app_key': [{ required: true, message: '请输入AppKey', trigger: 'blur' }],
  'config.app_secret': [{ required: true, message: '请输入AppSecret', trigger: 'blur' }],

  // Feishu config rules
  'config.app_id': [{ required: true, message: '请输入AppID', trigger: 'blur' }],
  'config.app_secret': [{ required: true, message: '请输入AppSecret', trigger: 'blur' }],

  // WeCom config rules
  'config.corp_id': [{ required: true, message: '请输入CorpId', trigger: 'blur' }],
  'config.agent_id': [{ required: true, message: '请输入AgentId', trigger: 'blur' }],
  'config.secret': [{ required: true, message: '请输入Secret', trigger: 'blur' }]
}

const typeLabel = (type: string) => {
  const map: Record<string, string> = {
    email: '邮件',
    sms: '短信',
    workwx: '企业微信',
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
      keyword: searchKeyword.value || undefined,
      search_field: searchField.value,
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

const handleSelectionChange = (rows: any[]) => {
  selectedRows.value = rows
}

const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请选择要删除的渠道')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRows.value.length} 个渠道吗？`,
      '删除确认',
      { type: 'warning' }
    )
    for (const row of selectedRows.value) {
      await deleteNotificationChannel(row.id)
    }
    ElMessage.success(`已删除 ${selectedRows.value.length} 个渠道`)
    loadData()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const showSettings = () => {
  ElMessage.info('设置功能开发中')
}

const handleCreate = () => {
  editingId.value = null
  dialogTitle.value = '新建渠道'
  form.name = ''
  form.type = ''
  form.scope = 'system'
  form.config = initializeFormConfig()
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  editingId.value = row.id
  dialogTitle.value = '编辑渠道'
  form.name = row.name
  form.type = row.type
  form.scope = row.scope || 'system'
  form.config = { ...initializeFormConfig(), ...row.config }
  dialogVisible.value = true
}

const handleTypeChange = () => {
  form.config = initializeFormConfig()
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    submitting.value = true

    const data = {
      name: form.name,
      type: form.type,
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
    if (e !== false) { // false means validation failed silently
      ElMessage.error(e.message || '操作失败')
    }
  } finally {
    submitting.value = false
  }
}

// 连接测试（弹窗内）
const handleTestConnectionInDialog = async () => {
  try {
    await formRef.value?.validate()
    testingConnection.value = true

    // 使用当前表单数据进行连接测试
    const data = {
      name: form.name || `temp-test-${Date.now()}`,
      type: form.type,
      config: form.config
    }

    await testNotificationChannel(0, data) // 传入数据而非ID
    ElMessage.success('连接测试成功')
  } catch (e: any) {
    if (e !== false) {
      ElMessage.error(e.message || '连接测试失败')
    }
  } finally {
    testingConnection.value = false
  }
}

// 连接测试（列表项）
const handleTestConnection = async (row: any) => {
  try {
    await testNotificationChannel(row.id)
    ElMessage.success('连接测试成功')
  } catch (e: any) {
    ElMessage.error(e.message || '连接测试失败')
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
.channels-page {
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

/* 搜索栏样式 */
.search-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 16px;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
}

.search-hint {
  color: #999;
  font-size: 12px;
}

.search-icon {
  cursor: pointer;
}

.filter-box {
  display: flex;
  gap: 8px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

/* 表单帮助文本样式 */
.form-help {
  margin-top: 4px;
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
}

.help-link {
  margin-left: 8px;
  font-size: 12px;
}

/* 弹窗底部按钮布局 */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

/* 分隔线样式 */
.el-divider {
  margin: 16px 0;
}

.el-divider__text {
  font-weight: 500;
  color: #303133;
}
</style>