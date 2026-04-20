<template>
  <div class="robots-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">机器人管理</span>
          <!-- 工具栏按钮 -->
          <div class="toolbar">
            <el-button @click="loadRobots" :icon="Refresh" circle title="刷新" />
            <el-button type="primary" @click="handleCreate">
              新建
            </el-button>
            <el-dropdown trigger="click" :disabled="selectedRows.length === 0">
              <el-button :disabled="selectedRows.length === 0">
                批量操作
                <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleBatchEnable">批量启用</el-dropdown-item>
                  <el-dropdown-item @click="handleBatchDisable">批量禁用</el-dropdown-item>
                  <el-dropdown-item @click="handleBatchDelete" divided>批量删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button @click="showSettings" :icon="Setting" circle title="设置" />
          </div>
        </div>
      </template>

      <!-- 搜索栏 -->
      <div class="search-bar">
        <div class="search-box">
          <el-select v-model="searchField" placeholder="搜索属性" style="width: 120px" class="search-field-select">
            <el-option label="名称" value="name" />
            <el-option label="状态" value="status" />
            <el-option label="启用状态" value="enabled" />
            <el-option label="类型" value="type" />
            <el-option label="所属项目" value="project" />
            <el-option label="创建时间" value="created_at" />
          </el-select>
          <el-input
            v-model="searchKeyword"
            placeholder="默认为名称搜索"
            style="width: 300px"
            clearable
            @keyup.enter="loadRobots"
          >
            <template #suffix>
              <el-icon class="search-icon" @click="loadRobots"><Search /></el-icon>
            </template>
          </el-input>
        </div>
      </div>

      <!-- 表格 -->
      <el-table
        :data="robots"
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
            {{ row.name }}
          </template>
        </el-table-column>
        <!-- 状态 -->
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'ready' ? 'success' : 'warning'">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <!-- 启用状态 -->
        <el-table-column label="启用状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ row.enabled ? '已启用' : '已禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <!-- 类型 -->
        <el-table-column label="类型" width="120">
          <template #default="{ row }">
            <el-tag>{{ getTypeName(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <!-- 所属项目 -->
        <el-table-column label="所属项目" width="120">
          <template #default="{ row }">
            {{ row.project?.name || '-' }}
          </template>
        </el-table-column>
        <!-- 共享范围 -->
        <el-table-column label="共享范围" width="100">
          <template #default="{ row }">
            {{ row.shared_scope || '私有' }}
          </template>
        </el-table-column>
        <!-- 创建时间 -->
        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <!-- 操作 -->
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-dropdown trigger="click">
              <el-button size="small" type="primary" link>
                更多
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleTest(row)">
                    测试
                  </el-dropdown-item>
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

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        class="pagination"
        @size-change="loadRobots"
        @current-change="loadRobots"
      />
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px" label-position="left">
        <!-- 项目 -->
        <el-form-item label="项目">
          <el-select v-model="form.project_id" placeholder="请选择项目" style="width: 300px" clearable>
            <el-option v-for="project in projects" :key="project.id" :label="project.name" :value="project.id" />
          </el-select>
        </el-form-item>
        <!-- 名称 -->
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="2-128字符，包含字母、数字、连字符，以字母开头，不能以连字符结尾" />
        </el-form-item>
        <!-- 类型 - Radio按钮组 -->
        <el-form-item label="类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio-button label="dingtalk">钉钉机器人</el-radio-button>
            <el-radio-button label="feishu">飞书机器人</el-radio-button>
            <el-radio-button label="workwx">企业微信机器人</el-radio-button>
            <el-radio-button label="webhook">Webhook</el-radio-button>
          </el-radio-group>
        </el-form-item>

        <!-- 钉钉/飞书/企业微信类型的Webhook地址 -->
        <el-form-item v-if="form.type !== 'webhook'" label="Webhook地址" prop="webhook_url">
          <el-input v-model="form.webhook_url" placeholder="请输入机器人Webhook地址" />
          <div class="form-extra">
            获取相关参数，请参考
            <a href="https://www.cloudpods.org/docs/guides/misc/notify/bot/" target="_blank" class="help-link">
              详细文档
              <el-icon><Link /></el-icon>
            </a>
          </div>
        </el-form-item>
        <!-- 钉钉/飞书/企业微信类型的密钥 -->
        <el-form-item v-if="form.type !== 'webhook'" label="密钥（如有）">
          <el-input v-model="form.secret" type="password" placeholder="请输入签名密钥" show-password />
        </el-form-item>

        <!-- Webhook类型的特殊字段 -->
        <el-form-item v-if="form.type === 'webhook'" label="URL" prop="webhook_url">
          <el-input v-model="form.webhook_url" placeholder="请输入Webhook URL地址" />
        </el-form-item>
        <el-form-item v-if="form.type === 'webhook'" label="请求头">
          <el-input v-model="form.header" placeholder="自定义请求头（JSON格式）" />
          <div class="form-extra">例如: {"Content-Type": "application/json", "Authorization": "Bearer xxx"}</div>
        </el-form-item>
        <el-form-item v-if="form.type === 'webhook'" label="请求体模板">
          <el-input v-model="form.body" type="textarea" :rows="3" placeholder="自定义请求体模板（JSON格式）" />
          <div class="form-extra" v-pre>可使用变量: {{message}}, {{title}}, {{timestamp}}</div>
        </el-form-item>
        <el-form-item v-if="form.type === 'webhook'" label="消息键">
          <el-input v-model="form.msg_key" placeholder="消息内容在请求体中的键名" />
          <div class="form-extra">例如: message、content、text</div>
        </el-form-item>
        <el-form-item v-if="form.type === 'webhook'" label="密钥">
          <el-input v-model="form.secret_key" type="password" placeholder="用于签名验证的密钥" show-password />
        </el-form-item>

        <!-- 描述 -->
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" placeholder="请输入描述" />
        </el-form-item>
        <!-- 启用 -->
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
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
import { Refresh, Setting, Search, ArrowDown, Link } from '@element-plus/icons-vue'
import {
  getRobots,
  createRobot,
  updateRobot,
  deleteRobot,
  enableRobot,
  disableRobot,
  testRobot
} from '@/api/message'
import { getProjects } from '@/api/iam'

interface Robot {
  id: number
  name: string
  type: string
  webhook_url: string
  secret: string
  header: string
  body: string
  msg_key: string
  secret_key: string
  description: string
  enabled: boolean
  status: string
  shared_scope: string
  created_at: string
  project_id: number
  project?: {
    id: number
    name: string
  }
}

interface Project {
  id: number
  name: string
}

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const robots = ref<Robot[]>([])
const editingId = ref<number | null>(null)
const formRef = ref<FormInstance>()
const selectedRows = ref<Robot[]>([])
const projects = ref<Project[]>([])

// 搜索相关
const searchField = ref('name')
const searchKeyword = ref('')

const form = reactive({
  name: '',
  type: 'dingtalk',
  webhook_url: '',
  secret: '',
  header: '',
  body: '',
  msg_key: '',
  secret_key: '',
  description: '',
  enabled: true,
  project_id: null as number | null
})

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  webhook_url: [
    { required: true, message: '请输入Webhook地址', trigger: 'blur' },
    { type: 'url', message: '请输入有效的URL', trigger: 'blur' }
  ]
}

// 类型映射常量
const TYPE_MAP: Record<string, string> = {
  dingtalk: '钉钉',
  feishu: '飞书',
  workwx: '企业微信',
  wechat: '企业微信',
  webhook: 'Webhook',
  lark: '飞书'
}

// 状态映射常量
const STATUS_MAP: Record<string, string> = {
  ready: '就绪',
  creating: '创建中',
  updating: '更新中',
  error: '错误'
}

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const loadRobots = async () => {
  loading.value = true
  try {
    const res = await getRobots({
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: searchKeyword.value || undefined,
      search_field: searchField.value
    })
    robots.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载机器人失败')
  } finally {
    loading.value = false
  }
}

const loadProjects = async () => {
  try {
    const res = await getProjects({ limit: 1000 })
    projects.value = res.items || []
    // 默认选中 system 项目
    const systemProject = projects.value.find(p => p.name === 'system')
    if (systemProject && !form.project_id) {
      form.project_id = systemProject.id
    }
  } catch (e: any) {
    console.error('加载项目失败:', e)
  }
}

// 选择变化处理
const handleSelectionChange = (rows: Robot[]) => {
  selectedRows.value = rows
}

// 批量启用
const handleBatchEnable = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请选择要启用的机器人')
    return
  }
  try {
    await ElMessageBox.confirm(
      `确定要启用选中的 ${selectedRows.value.length} 个机器人吗？`,
      '启用确认',
      { type: 'info' }
    )
    await Promise.all(selectedRows.value.map(row => enableRobot(row.id)))
    ElMessage.success(`已启用 ${selectedRows.value.length} 个机器人`)
    loadRobots()
  } catch (e: any) {
    if (e !== 'cancel' && e !== 'close') {
      ElMessage.error(e.message || '启用失败')
    }
  }
}

// 批量禁用
const handleBatchDisable = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请选择要禁用的机器人')
    return
  }
  try {
    await ElMessageBox.confirm(
      `确定要禁用选中的 ${selectedRows.value.length} 个机器人吗？`,
      '禁用确认',
      { type: 'warning' }
    )
    await Promise.all(selectedRows.value.map(row => disableRobot(row.id)))
    ElMessage.success(`已禁用 ${selectedRows.value.length} 个机器人`)
    loadRobots()
  } catch (e: any) {
    if (e !== 'cancel' && e !== 'close') {
      ElMessage.error(e.message || '禁用失败')
    }
  }
}

// 批量删除
const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请选择要删除的机器人')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRows.value.length} 个机器人吗？`,
      '删除确认',
      { type: 'warning' }
    )
    await Promise.all(selectedRows.value.map(row => deleteRobot(row.id)))
    ElMessage.success(`已删除 ${selectedRows.value.length} 个机器人`)
    loadRobots()
  } catch (e: any) {
    if (e !== 'cancel' && e !== 'close') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

// 设置弹窗
const showSettings = () => {
  ElMessage.info('设置功能开发中')
}

const handleCreate = () => {
  editingId.value = null
  dialogTitle.value = '新建机器人'
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: Robot) => {
  editingId.value = row.id
  dialogTitle.value = '编辑机器人'
  form.name = row.name
  form.type = row.type
  form.webhook_url = row.webhook_url
  form.secret = row.secret || ''
  form.header = row.header || ''
  form.body = row.body || ''
  form.msg_key = row.msg_key || ''
  form.secret_key = row.secret_key || ''
  form.description = row.description || ''
  form.enabled = row.enabled
  form.project_id = row.project_id
  dialogVisible.value = true
}

const handleTest = async (row: Robot) => {
  try {
    await testRobot(row.id)
    ElMessage.success('测试消息已发送')
  } catch (e: any) {
    ElMessage.error(e.message || '测试失败')
  }
}

const handleEnable = async (row: Robot) => {
  try {
    await enableRobot(row.id)
    ElMessage.success('启用成功')
    loadRobots()
  } catch (e: any) {
    ElMessage.error(e.message || '启用失败')
  }
}

const handleDisable = async (row: Robot) => {
  try {
    await disableRobot(row.id)
    ElMessage.success('禁用成功')
    loadRobots()
  } catch (e: any) {
    ElMessage.error(e.message || '禁用失败')
  }
}

const handleDelete = async (row: Robot) => {
  try {
    await ElMessageBox.confirm(`确定要删除机器人 "${row.name}" 吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await deleteRobot(row.id)
    ElMessage.success('删除成功')
    loadRobots()
  } catch (e: any) {
    if (e !== 'cancel' && e !== 'close') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const resetForm = () => {
  form.name = ''
  form.type = 'dingtalk'
  form.webhook_url = ''
  form.secret = ''
  form.header = ''
  form.body = ''
  form.msg_key = ''
  form.secret_key = ''
  form.description = ''
  form.enabled = true
  form.project_id = projects.value.find(p => p.name === 'system')?.id || null
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return // 验证失败，不继续执行
  }
  submitting.value = true
  try {
    const data = {
      name: form.name,
      type: form.type,
      webhook_url: form.webhook_url,
      secret: form.secret,
      header: form.header,
      body: form.body,
      msg_key: form.msg_key,
      secret_key: form.secret_key,
      description: form.description,
      enabled: form.enabled,
      project_id: form.project_id
    }

    if (editingId.value) {
      await updateRobot(editingId.value, data)
      ElMessage.success('更新成功')
    } else {
      await createRobot(data)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    loadRobots()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

const getTypeName = (type: string) => TYPE_MAP[type] || type

const getStatusName = (status: string) => STATUS_MAP[status] || status || '就绪'

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
    loadRobots(),
    loadProjects()
  ])
})
</script>

<style scoped>
.robots-page {
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

.search-icon {
  cursor: pointer;
}

/* 分页样式 */
.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

/* 弹窗底部按钮布局 */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

/* 表单附加信息 */
.form-extra {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.help-link {
  color: #409eff;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
</style>