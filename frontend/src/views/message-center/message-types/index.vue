<template>
  <div class="message-types-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">消息订阅设置</span>
          <!-- 工具栏按钮 - 参考 Cloudpods -->
          <div class="toolbar">
            <el-button @click="handleBatchEnable" :disabled="selectedRows.length === 0">
              启 用
            </el-button>
            <el-button @click="handleBatchDisable" :disabled="selectedRows.length === 0">
              禁 用
            </el-button>
          </div>
        </div>
      </template>

      <!-- 筛选区域 - 参考 Cloudpods 设计 -->
      <div class="filter-bar">
        <el-select v-model="filters.searchField" placeholder="搜索字段" style="width: 120px">
          <el-option label="名称" value="name" />
          <el-option label="描述" value="description" />
        </el-select>
        <el-input
          v-model="filters.keyword"
          placeholder="请输入关键词"
          style="width: 200px"
          clearable
          @keyup.enter="loadData"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-select v-model="filters.type" placeholder="类型" clearable style="width: 150px" @change="loadData">
          <el-option label="全部" value="" />
          <el-option label="安全信息" value="security" />
          <el-option label="资源" value="resource" />
          <el-option label="自动化流程" value="automated_process" />
        </el-select>
        <el-select v-model="filters.enabled" placeholder="状态" clearable style="width: 120px" @change="loadData">
          <el-option label="全部" value="" />
          <el-option label="启用" value="true" />
          <el-option label="禁用" value="false" />
        </el-select>
        <el-button type="primary" @click="loadData">
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
        :data="messageTypes"
        v-loading="loading"
        border
        stripe
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <!-- 选择列 -->
        <el-table-column type="selection" width="50" />
        <!-- 消息名称 -->
        <el-table-column label="消息名称" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">
              {{ row.display_name || row.name }}
            </el-button>
          </template>
        </el-table-column>
        <!-- 消息类型 -->
        <el-table-column label="消息类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" effect="plain">
              {{ getTypeName(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <!-- 启用状态 -->
        <el-table-column label="启用状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <!-- 操作 - 参考 Cloudpods：三个按钮始终显示，启用按钮在已启用状态时禁用 -->
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="handleReceiverManage(row)">
              接收管理
            </el-button>
            <el-button
              size="small"
              type="primary"
              link
              :disabled="row.enabled"
              @click="handleEnable(row)"
            >
              启用
            </el-button>
            <el-button
              size="small"
              type="primary"
              link
              :disabled="!row.enabled"
              @click="handleDisable(row)"
            >
              禁用
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
        @size-change="loadData"
        @current-change="loadData"
        class="pagination"
      />
    </el-card>

    <!-- 详情侧边页 - 参考 Cloudpods side-page 设计 -->
    <el-drawer
      v-model="detailVisible"
      :title="currentType?.display_name || currentType?.name"
      direction="rtl"
      size="50%"
      :before-close="handleDetailClose"
    >
      <template #header>
        <div class="drawer-header">
          <div class="drawer-header-info">
            <span class="drawer-header-label">消息订阅设置</span>
            <h4 class="drawer-header-title">{{ currentType?.display_name || currentType?.name }}</h4>
            <div class="drawer-header-actions">
              <el-button size="small" type="primary" link @click="handleReceiverManage(currentType!)">
                接收管理
              </el-button>
              <el-button size="small" type="primary" link :disabled="currentType?.enabled" @click="handleEnable(currentType!)">
                启用
              </el-button>
              <el-button size="small" type="primary" link :disabled="!currentType?.enabled" @click="handleDisable(currentType!)">
                禁用
              </el-button>
            </div>
          </div>
        </div>
      </template>

      <el-tabs v-model="detailTab" class="detail-tabs">
        <!-- 详情 Tab -->
        <el-tab-pane label="详情" name="detail">
          <div class="detail-content">
            <!-- 基本信息 - 左侧 -->
            <div class="detail-left">
              <div class="detail-section">
                <div class="detail-section-title">
                  <el-icon><Document /></el-icon>
                  <span>基本信息</span>
                </div>
                <div class="detail-items">
                  <div class="detail-item">
                    <span class="detail-item-title">ID</span>
                    <span class="detail-item-value">{{ currentType?.id }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-item-title">消息名称</span>
                    <span class="detail-item-value">{{ currentType?.display_name || currentType?.name }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-item-title">消息类型</span>
                    <span class="detail-item-value">
                      <el-tag :type="getTypeTag(currentType?.type)" effect="plain">
                        {{ getTypeName(currentType?.type) }}
                      </el-tag>
                    </span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-item-title">启用状态</span>
                    <span class="detail-item-value">
                      <el-tag :type="currentType?.enabled ? 'success' : 'info'">
                        {{ currentType?.enabled ? '启用' : '禁用' }}
                      </el-tag>
                    </span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-item-title">创建时间</span>
                    <span class="detail-item-value">{{ currentType?.created_at }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-item-title">更新时间</span>
                    <span class="detail-item-value">{{ currentType?.updated_at }}</span>
                  </div>
                </div>
              </div>
            </div>
            <!-- 资源列表 - 右侧 -->
            <div class="detail-right">
              <div class="detail-section">
                <div class="detail-section-title">
                  <el-icon><List /></el-icon>
                  <span>资源列表</span>
                </div>
                <div class="detail-items">
                  <div class="detail-item">
                    <span class="detail-item-title">资源类型</span>
                    <span class="detail-item-value">
                      <el-tag v-for="rt in parseResourceTypes(currentType?.resource_types || '')" :key="rt" style="margin: 2px">
                        {{ rt }}
                      </el-tag>
                      <span v-if="parseResourceTypes(currentType?.resource_types || '').length === 0">无</span>
                    </span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-item-title">描述</span>
                    <span class="detail-item-value">{{ currentType?.description || '无' }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-item-title">分组键</span>
                    <span class="detail-item-value">
                      <el-tag v-for="gk in parseResourceTypes(currentType?.group_keys || '')" :key="gk" type="warning" style="margin: 2px">
                        {{ gk }}
                      </el-tag>
                      <span v-if="parseResourceTypes(currentType?.group_keys || '').length === 0">无</span>
                    </span>
                  </div>
                  <div class="detail-item" v-if="currentType?.advance_days">
                    <span class="detail-item-title">提前天数</span>
                    <span class="detail-item-value">{{ parseAdvanceDays(currentType?.advance_days) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-tab-pane>

        <!-- 接收管理 Tab -->
        <el-tab-pane label="接收管理" name="receiver">
          <div class="receiver-tab-content">
            <el-button type="primary" @click="handleAddReceiver" style="margin-bottom: 16px">
              添加接收人
            </el-button>
            <el-table :data="receivers" v-loading="receiverLoading" border stripe>
              <el-table-column label="接收人" min-width="150">
                <template #default="{ row }">
                  {{ row.receiver_name }}
                </template>
              </el-table-column>
              <el-table-column label="类型" width="100">
                <template #default="{ row }">
                  <el-tag :type="getReceiverTypeTag(row.receiver_type)" effect="plain">
                    {{ getReceiverTypeName(row.receiver_type) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="通知渠道" min-width="150">
                <template #default="{ row }">
                  <el-tag v-if="row.inbox" size="small" style="margin: 2px">站内信</el-tag>
                  <el-tag v-if="row.email" size="small" type="success" style="margin: 2px">邮件</el-tag>
                  <el-tag v-if="row.wechat" size="small" type="warning" style="margin: 2px">企业微信</el-tag>
                  <el-tag v-if="row.dingtalk" size="small" type="danger" style="margin: 2px">钉钉</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="状态" width="80">
                <template #default="{ row }">
                  <el-tag :type="row.enabled ? 'success' : 'info'">
                    {{ row.enabled ? '启用' : '禁用' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="120">
                <template #default="{ row }">
                  <el-button size="small" type="primary" link @click="handleEditReceiver(row)">编辑</el-button>
                  <el-button size="small" type="danger" link @click="handleDeleteReceiver(row)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <!-- 操作日志 Tab -->
        <el-tab-pane label="操作日志" name="logs">
          <el-empty description="暂无操作日志" />
        </el-tab-pane>
      </el-tabs>

      <template #footer>
        <el-button @click="handleDetailClose">取 消</el-button>
        <el-button type="primary" @click="handleDetailClose">确 定</el-button>
      </template>
    </el-drawer>

    <!-- 添加/编辑接收人弹窗 -->
    <el-dialog v-model="receiverFormVisible" :title="receiverFormTitle" width="500px">
      <el-form ref="receiverFormRef" :model="receiverForm" :rules="receiverRules" label-width="100px">
        <el-form-item label="接收人类型" prop="receiver_type">
          <el-select v-model="receiverForm.receiver_type" placeholder="请选择类型" style="width: 100%">
            <el-option label="用户" value="user" />
            <el-option label="用户组" value="group" />
            <el-option label="角色" value="role" />
          </el-select>
        </el-form-item>
        <el-form-item label="接收人" prop="receiver_id">
          <el-select v-model="receiverForm.receiver_id" placeholder="请选择接收人" style="width: 100%" filterable>
            <el-option v-for="r in receiverOptions" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="通知渠道">
          <el-checkbox v-model="receiverForm.inbox">站内信</el-checkbox>
          <el-checkbox v-model="receiverForm.email">邮件</el-checkbox>
          <el-checkbox v-model="receiverForm.wechat">企业微信</el-checkbox>
          <el-checkbox v-model="receiverForm.dingtalk">钉钉</el-checkbox>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="receiverFormVisible = false">取消</el-button>
        <el-button type="primary" :loading="receiverSubmitting" @click="handleReceiverSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Search, RefreshRight, Document, List } from '@element-plus/icons-vue'
import request from '@/utils/request'

interface MessageType {
  id: number
  name: string
  display_name: string
  description: string
  type: string // security/resource/automated_process
  enabled: boolean
  title_cn: string
  title_en: string
  content_cn: string
  content_en: string
  resource_types: string
  group_keys: string
  advance_days: string
  is_system: boolean
  can_delete: boolean
  can_update: boolean
  created_at: string
  updated_at: string
}

interface Receiver {
  id: number
  topic_id: number
  receiver_type: string
  receiver_id: number
  receiver_name: string
  inbox: boolean
  email: boolean
  wechat: boolean
  dingtalk: boolean
  enabled: boolean
}

const loading = ref(false)
const detailVisible = ref(false)
const detailTab = ref('detail')
const currentType = ref<MessageType | null>(null)
const selectedRows = ref<MessageType[]>([])

const pagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

const filters = reactive({
  searchField: 'name',
  keyword: '',
  type: '',
  enabled: ''
})

const messageTypes = ref<MessageType[]>([])

// 接收管理相关
const receiverLoading = ref(false)
const receivers = ref<Receiver[]>([])
const receiverFormVisible = ref(false)
const receiverFormTitle = ref('添加接收人')
const receiverSubmitting = ref(false)
const receiverFormRef = ref<FormInstance>()
const editingReceiverId = ref<number | null>(null)
const receiverOptions = ref<any[]>([])

const receiverForm = reactive({
  receiver_type: 'user',
  receiver_id: null as number | null,
  inbox: true,
  email: false,
  wechat: false,
  dingtalk: false
})

const receiverRules = {
  receiver_type: [{ required: true, message: '请选择接收人类型', trigger: 'change' }],
  receiver_id: [{ required: true, message: '请选择接收人', trigger: 'change' }]
}

// 类型映射 - 参考 Cloudpods
const getTypeName = (type: string | undefined) => {
  if (!type) return ''
  const map: Record<string, string> = {
    security: '安全消息',
    resource: '资源消息',
    automated_process: '自动化流程'
  }
  return map[type] || type
}

const getTypeTag = (type: string | undefined) => {
  if (!type) return 'info'
  const map: Record<string, any> = {
    security: 'danger',
    resource: 'primary',
    automated_process: 'warning'
  }
  return map[type] || 'info'
}

// 接收人类型映射
const getReceiverTypeName = (type: string) => {
  const map: Record<string, string> = {
    user: '用户',
    group: '用户组',
    role: '角色'
  }
  return map[type] || type
}

const getReceiverTypeTag = (type: string) => {
  const map: Record<string, any> = {
    user: 'success',
    group: 'warning',
    role: 'primary'
  }
  return map[type] || 'info'
}

// 解析资源类型
const parseResourceTypes = (resourceTypes: string | undefined) => {
  if (!resourceTypes) return []
  try {
    return JSON.parse(resourceTypes)
  } catch {
    return []
  }
}

// 解析提前天数
const parseAdvanceDays = (advanceDays: string | undefined) => {
  if (!advanceDays) return ''
  try {
    const days = JSON.parse(advanceDays)
    return days.join('、') + ' 天'
  } catch {
    return advanceDays
  }
}

// 加载消息类型列表
const loadData = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit
    }

    // 搜索字段和关键词
    if (filters.keyword) {
      params.keyword = filters.keyword
      params.search_field = filters.searchField
    }

    if (filters.type) {
      params.type = filters.type
    }
    if (filters.enabled) {
      params.enabled = filters.enabled
    }

    const res = await request({ url: '/message-types', method: 'get', params })
    messageTypes.value = res.data || res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载消息类型失败')
  } finally {
    loading.value = false
  }
}

// 选择变化
const handleSelectionChange = (rows: MessageType[]) => {
  selectedRows.value = rows
}

// 批量启用 - 有确认弹窗
const handleBatchEnable = async () => {
  const nonSystemRows = selectedRows.value.filter(r => !r.is_system)
  if (nonSystemRows.length === 0) {
    ElMessage.warning('请选择非系统内置的消息类型')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要启用选中的 ${nonSystemRows.length} 条消息类型吗？`,
      '启用确认',
      { type: 'warning' }
    )
    for (const row of nonSystemRows) {
      await request({ url: `/message-types/${row.id}/enable`, method: 'post' })
    }
    ElMessage.success(`已启用 ${nonSystemRows.length} 条消息类型`)
    loadData()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '批量启用失败')
    }
  }
}

// 批量禁用 - 有确认弹窗
const handleBatchDisable = async () => {
  const nonSystemRows = selectedRows.value.filter(r => !r.is_system)
  if (nonSystemRows.length === 0) {
    ElMessage.warning('请选择非系统内置的消息类型')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要禁用选中的 ${nonSystemRows.length} 条消息类型吗？`,
      '禁用确认',
      { type: 'warning' }
    )
    for (const row of nonSystemRows) {
      await request({ url: `/message-types/${row.id}/disable`, method: 'post' })
    }
    ElMessage.success(`已禁用 ${nonSystemRows.length} 条消息类型`)
    loadData()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '批量禁用失败')
    }
  }
}

// 查看详情 - 参考 Cloudpods side-page
const handleView = (row: MessageType) => {
  currentType.value = row
  detailTab.value = 'detail'
  detailVisible.value = true
  loadReceivers(row.id)
}

// 关闭详情侧边页
const handleDetailClose = () => {
  detailVisible.value = false
  currentType.value = null
}

// 启用 - 有确认弹窗
const handleEnable = async (row: MessageType) => {
  try {
    await ElMessageBox.confirm(
      `确定要启用消息类型「${row.display_name || row.name}」吗？`,
      '启用确认',
      { type: 'info' }
    )
    await request({ url: `/message-types/${row.id}/enable`, method: 'post' })
    ElMessage.success('已启用')
    row.enabled = true
    if (currentType.value && currentType.value.id === row.id) {
      currentType.value.enabled = true
    }
    loadData()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '启用失败')
    }
  }
}

// 禁用 - 有确认弹窗
const handleDisable = async (row: MessageType) => {
  try {
    await ElMessageBox.confirm(
      `确定要禁用消息类型「${row.display_name || row.name}」吗？`,
      '禁用确认',
      { type: 'warning' }
    )
    await request({ url: `/message-types/${row.id}/disable`, method: 'post' })
    ElMessage.success('已禁用')
    row.enabled = false
    if (currentType.value && currentType.value.id === row.id) {
      currentType.value.enabled = false
    }
    loadData()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '禁用失败')
    }
  }
}

// 接收管理 - 在详情页切换到接收管理 tab
const handleReceiverManage = (row: MessageType) => {
  currentType.value = row
  detailTab.value = 'receiver'
  detailVisible.value = true
  loadReceivers(row.id)
  loadReceiverOptions()
}

const loadReceivers = async (topicId: number) => {
  receiverLoading.value = true
  try {
    const res = await request({ url: '/topic-receivers', method: 'get', params: { topic_id: topicId } })
    receivers.value = res.data || res.items || []
  } catch (e: any) {
    ElMessage.error(e.message || '加载接收人失败')
  } finally {
    receiverLoading.value = false
  }
}

const loadReceiverOptions = async () => {
  try {
    // 加载用户列表
    const usersRes = await request({ url: '/users', method: 'get', params: { limit: 100 } })
    const users = (usersRes.data || usersRes.items || []).map(u => ({ id: u.id, name: u.display_name || u.name, type: 'user' }))

    // 加载用户组列表
    const groupsRes = await request({ url: '/groups', method: 'get', params: { limit: 100 } })
    const groups = (groupsRes.data || groupsRes.items || []).map(g => ({ id: g.id, name: g.name, type: 'group' }))

    // 加载角色列表
    const rolesRes = await request({ url: '/roles', method: 'get', params: { limit: 100 } })
    const roles = (rolesRes.data || rolesRes.items || []).map(r => ({ id: r.id, name: r.display_name || r.name, type: 'role' }))

    receiverOptions.value = [...users, ...groups, ...roles]
  } catch {
    receiverOptions.value = []
  }
}

const handleAddReceiver = () => {
  editingReceiverId.value = null
  receiverFormTitle.value = '添加接收人'
  Object.assign(receiverForm, {
    receiver_type: 'user',
    receiver_id: null,
    inbox: true,
    email: false,
    wechat: false,
    dingtalk: false
  })
  receiverFormVisible.value = true
}

const handleEditReceiver = (row: Receiver) => {
  editingReceiverId.value = row.id
  receiverFormTitle.value = '编辑接收人'
  Object.assign(receiverForm, {
    receiver_type: row.receiver_type,
    receiver_id: row.receiver_id,
    inbox: row.inbox,
    email: row.email,
    wechat: row.wechat,
    dingtalk: row.dingtalk
  })
  receiverFormVisible.value = true
}

const handleReceiverSubmit = async () => {
  await receiverFormRef.value?.validate()
  receiverSubmitting.value = true
  try {
    const data = {
      topic_id: currentType.value?.id,
      ...receiverForm
    }

    if (editingReceiverId.value) {
      await request({ url: `/topic-receivers/${editingReceiverId.value}`, method: 'put', data })
      ElMessage.success('更新成功')
    } else {
      await request({ url: '/topic-receivers', method: 'post', data })
      ElMessage.success('添加成功')
    }

    receiverFormVisible.value = false
    loadReceivers(currentType.value?.id!)
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    receiverSubmitting.value = false
  }
}

const handleDeleteReceiver = async (row: Receiver) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除接收人「${row.receiver_name}」吗？`,
      '删除确认',
      { type: 'warning' }
    )
    await request({ url: `/topic-receivers/${row.id}`, method: 'delete' })
    ElMessage.success('删除成功')
    loadReceivers(currentType.value?.id!)
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const resetFilters = () => {
  filters.searchField = 'name'
  filters.keyword = ''
  filters.type = ''
  filters.enabled = ''
  pagination.page = 1
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.message-types-page {
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

/* 筛选区域 */
.filter-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

/* 详情侧边页样式 - 参考 Cloudpods side-page */
.drawer-header {
  display: flex;
  align-items: center;
}

.drawer-header-info {
  display: flex;
  flex-direction: column;
}

.drawer-header-label {
  color: #999;
  font-size: 14px;
}

.drawer-header-title {
  margin: 8px 0;
  font-size: 18px;
}

.drawer-header-actions {
  display: flex;
  gap: 8px;
}

.detail-tabs {
  height: 100%;
}

.detail-content {
  display: flex;
  gap: 20px;
  padding: 16px;
}

.detail-left {
  flex: 1;
  min-width: 300px;
}

.detail-right {
  flex: 1;
  min-width: 300px;
}

.detail-section {
  background: #f9f9f9;
  border-radius: 4px;
  padding: 16px;
}

.detail-section-title {
  display: flex;
  align-items: center;
  font-weight: bold;
  margin-bottom: 12px;
  color: #333;
}

.detail-section-title .el-icon {
  margin-right: 8px;
}

.detail-items {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.detail-item {
  display: flex;
  align-items: flex-start;
}

.detail-item-title {
  color: #666;
  min-width: 100px;
  flex-shrink: 0;
}

.detail-item-value {
  flex: 1;
  word-break: break-word;
}

.receiver-tab-content {
  padding: 16px;
}
</style>