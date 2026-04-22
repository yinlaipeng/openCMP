<template>
  <div class="sync-policies-container">
    <!-- 页头 -->
    <div class="page-header">
      <h2>同步策略</h2>
      <div class="toolbar">
        <el-button size="default" @click="loadPolicies">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          新建策略
        </el-button>
        <el-dropdown trigger="click" @command="handleBatchCommand" :disabled="selectedIds.length === 0">
          <el-button>
            <el-icon><Operation /></el-icon>
            批量操作
            <el-badge v-if="selectedIds.length > 0" :value="selectedIds.length" type="primary" />
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="batchEnable" :icon="CircleCheck">批量启用</el-dropdown-item>
              <el-dropdown-item command="batchDisable" :icon="CircleClose">批量禁用</el-dropdown-item>
              <el-dropdown-item command="batchDelete" :icon="Delete" divided>
                <span style="color: var(--el-color-danger)">批量删除</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button @click="handleExport">
          <el-icon><Download /></el-icon>
          导出
        </el-button>
      </div>
    </div>

    <!-- 顶部分类Tab -->
    <el-tabs v-model="activeTab" class="status-tabs" @tab-change="handleTabChange">
      <el-tab-pane label="全部" name="all" />
      <el-tab-pane label="已启用" name="enabled" />
      <el-tab-pane label="已禁用" name="disabled" />
    </el-tabs>

    <!-- 筛选区 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="queryForm" @submit.prevent="loadPolicies">
        <el-form-item label="名称">
          <el-input v-model="queryForm.name" placeholder="支持策略名称、ID搜索" clearable style="width: 200px">
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="规则类型">
          <el-select v-model="queryForm.condition_type" placeholder="全部" clearable style="width: 140px">
            <el-option label="全部匹配标签" value="all_match" />
            <el-option label="至少一个标签" value="any_match" />
            <el-option label="Key匹配" value="key_match" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadPolicies">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-table
      ref="tableRef"
      :data="policies"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="50" />
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="策略名称" min-width="150">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleView(row)" class="name-link">
            {{ row.name }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column prop="rules" label="规则数量" width="100">
        <template #default="{ row }">
          <el-tag type="info" size="small">{{ row.rules?.length || 0 }} 条</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="scope" label="应用范围" width="120">
        <template #default="{ row }">
          {{ getScopeText(row.scope) }}
        </template>
      </el-table-column>
      <el-table-column prop="execution_count" label="执行次数" width="100">
        <template #default="{ row }">
          <span>{{ row.execution_count || 0 }} 次</span>
        </template>
      </el-table-column>
      <el-table-column prop="last_execution_time" label="最后执行" width="160">
        <template #default="{ row }">
          {{ row.last_execution_time ? formatDate(row.last_execution_time) : '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="enabled" label="状态" width="100">
        <template #default="{ row }">
          <el-switch
            v-model="row.enabled"
            size="small"
            @change="handleToggle(row)"
          />
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="160">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="140" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="handleView(row)">查看</el-button>
          <el-dropdown trigger="click" @command="(cmd: string) => handleDropdownCommand(cmd, row)" style="margin-left: 8px">
            <el-button size="small" link>
              更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="execute" :icon="PlayOne">执行策略</el-dropdown-item>
                <el-dropdown-item command="edit" :icon="EditPen">编辑</el-dropdown-item>
                <el-dropdown-item command="copy" :icon="CopyDocument">复制</el-dropdown-item>
                <el-dropdown-item divided command="toggle" :icon="row.enabled ? CircleClose : CircleCheck">
                  {{ row.enabled ? '禁用' : '启用' }}
                </el-dropdown-item>
                <el-dropdown-item command="delete" :icon="Delete">
                  <span style="color: var(--el-color-danger)">删除</span>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      class="pagination"
    />

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="showDialog"
      :title="isEdit ? '编辑策略' : '添加策略'"
      width="800px"
      destroy-on-close
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="策略名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入策略名称" />
        </el-form-item>

        <el-form-item label="备注">
          <el-input v-model="form.remarks" type="textarea" :rows="2" placeholder="请输入备注信息" />
        </el-form-item>

        <el-form-item label="所属域" prop="domain_id">
          <el-select v-model="form.domain_id" placeholder="请选择所属域" style="width: 100%">
            <el-option v-for="domain in domains" :key="domain.id" :label="domain.name" :value="domain.id" />
          </el-select>
        </el-form-item>

        <!-- 规则配置 -->
        <el-form-item label="同步规则">
          <div class="rules-container">
            <el-button type="primary" size="small" @click="addRule" style="margin-bottom: 10px">
              <el-icon><Plus /></el-icon>
              添加规则
            </el-button>

            <div v-for="(rule, index) in form.rules" :key="index" class="rule-item">
              <el-card shadow="hover" class="rule-card">
                <template #header>
                  <div class="rule-header">
                    <span>规则 {{ index + 1 }}</span>
                    <el-button type="danger" size="small" @click="removeRule(index)" :icon="Delete" circle />
                  </div>
                </template>

                <el-form-item label="条件类型">
                  <el-select v-model="rule.condition_type" placeholder="请选择条件类型" style="width: 100%">
                    <el-option label="全部匹配标签" value="all_match" />
                    <el-option label="至少一个标签" value="any_match" />
                    <el-option label="根据标签Key匹配" value="key_match" />
                  </el-select>
                </el-form-item>

                <el-form-item label="资源映射">
                  <el-select v-model="rule.resource_mapping" placeholder="请选择映射方式" style="width: 100%">
                    <el-option label="指定项目" value="specify_project" />
                    <el-option label="根据名称映射" value="specify_name" />
                  </el-select>
                </el-form-item>

                <el-form-item label="目标项目" v-if="rule.resource_mapping === 'specify_project'">
                  <el-select v-model="rule.target_project_id" placeholder="请选择目标项目" style="width: 100%">
                    <el-option v-for="project in projects" :key="project.id" :label="project.name" :value="project.id" />
                  </el-select>
                </el-form-item>

                <!-- 标签配置 -->
                <el-form-item label="匹配标签">
                  <div class="tags-container">
                    <el-button size="small" @click="addTag(index)" style="margin-bottom: 8px">
                      <el-icon><Plus /></el-icon>
                      添加标签
                    </el-button>

                    <div v-for="(tag, tagIndex) in rule.tags" :key="tagIndex" class="tag-item">
                      <el-input v-model="tag.tag_key" placeholder="标签Key" style="width: 150px; margin-right: 8px" />
                      <el-input v-model="tag.tag_value" placeholder="标签Value" style="width: 150px; margin-right: 8px" />
                      <el-button type="danger" size="small" @click="removeTag(index, tagIndex)" :icon="Delete" circle />
                    </div>
                  </div>
                </el-form-item>
              </el-card>
            </div>
          </div>
        </el-form-item>

        <el-form-item label="启用状态">
          <el-switch v-model="form.enabled" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 详情抽屉 -->
    <el-drawer
      v-model="showDetailDrawer"
      title="策略详情"
      size="60%"
      direction="rtl"
    >
      <!-- 顶部区域 -->
      <div class="drawer-header" v-if="currentPolicy">
        <div class="policy-icon">
          <el-avatar :size="48" :style="{ backgroundColor: '#409EFF' }">
            <el-icon :size="32"><Setting /></el-icon>
          </el-avatar>
        </div>
        <div class="policy-info">
          <h3>{{ currentPolicy.name }}</h3>
          <div class="policy-tags">
            <el-tag size="small">{{ getScopeText(currentPolicy.scope) }}</el-tag>
            <el-tag size="small" type="info">{{ currentPolicy.rules?.length || 0 }} 条规则</el-tag>
            <el-switch
              v-model="currentPolicy.enabled"
              size="small"
              @change="handleToggle(currentPolicy)"
            />
          </div>
        </div>
        <div class="quick-actions">
          <el-button size="small" type="primary" @click="handleExecute(currentPolicy)">
            <el-icon><VideoPlay /></el-icon>
            执行
          </el-button>
          <el-button size="small" @click="handleEdit(currentPolicy)">
            <el-icon><EditPen /></el-icon>
            编辑
          </el-button>
          <el-dropdown trigger="click" @command="(cmd: string) => handleDrawerCommand(cmd, currentPolicy)">
            <el-button size="small">
              更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="copy" :icon="CopyDocument">复制策略</el-dropdown-item>
                <el-dropdown-item command="delete" divided :icon="Delete">
                  <span style="color: var(--el-color-danger)">删除</span>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>

      <!-- 详情内容 -->
      <el-tabs v-model="detailTab" type="border-card" style="margin-top: 16px">
        <el-tab-pane label="规则概览" name="rulesOverview">
          <div class="rules-detail" v-if="currentPolicy?.rules?.length">
            <el-collapse>
              <el-collapse-item v-for="(rule, index) in currentPolicy.rules" :key="index" :title="`规则 ${index + 1}`">
                <el-descriptions :column="2" border size="small">
                  <el-descriptions-item label="条件类型">{{ getConditionText(rule.condition_type) }}</el-descriptions-item>
                  <el-descriptions-item label="资源映射">{{ getMappingText(rule.resource_mapping) }}</el-descriptions-item>
                  <el-descriptions-item label="目标项目" v-if="rule.target_project_id">{{ rule.target_project_name || rule.target_project_id }}</el-descriptions-item>
                  <el-descriptions-item label="匹配标签" :span="2">
                    <el-tag v-for="tag in rule.tags" :key="tag.tag_key" style="margin-right: 8px" size="small">
                      {{ tag.tag_key }}: {{ tag.tag_value }}
                    </el-tag>
                  </el-descriptions-item>
                </el-descriptions>
              </el-collapse-item>
            </el-collapse>
          </div>
          <el-empty v-else description="暂无规则" />
        </el-tab-pane>
        <el-tab-pane label="执行日志" name="executionLogs">
          <el-empty description="执行日志功能开发中" />
        </el-tab-pane>
        <el-tab-pane label="映射结果" name="mappingResults">
          <el-empty description="映射结果功能开发中" />
        </el-tab-pane>
      </el-tabs>

      <!-- 基本信息 -->
      <el-descriptions :column="2" border style="margin-top: 20px" title="基本信息">
        <el-descriptions-item label="策略ID">{{ currentPolicy?.id }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ currentPolicy?.domain_id }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ currentPolicy?.remarks || '无' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentPolicy?.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentPolicy?.updated_at) }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { Plus, Delete, Document, ArrowDown, EditPen, CircleCheck, CircleClose, Refresh, Operation, Download, Search, VideoPlay, CopyDocument, Setting } from '@element-plus/icons-vue'
import { getSyncPolicies, getSyncPolicy, createSyncPolicy, updateSyncPolicy, deleteSyncPolicy, updateSyncPolicyStatus } from '@/api/sync-policy'
import { getDomains } from '@/api/iam'
import { getProjects } from '@/api/project'
import type { SyncPolicy, CreateSyncPolicyRequest } from '@/types/sync-policy'

interface TagItem {
  tag_key: string
  tag_value: string
}

interface RuleForm {
  condition_type: 'all_match' | 'any_match' | 'key_match'
  resource_mapping: 'specify_project' | 'specify_name'
  target_project_id?: number
  target_project_name?: string
  tags: TagItem[]
}

const policies = ref<SyncPolicy[]>([])
const domains = ref<{ id: number; name: string }[]>([])
const projects = ref<{ id: number; name: string }[]>([])
const loading = ref(false)
const showDialog = ref(false)
const showDetailDrawer = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const tableRef = ref()
const currentPolicy = ref<SyncPolicy | null>(null)
const selectedIds = ref<number[]>([])

// 分页数据
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 顶部Tab
const activeTab = ref('all')

// 详情抽屉Tab
const detailTab = ref('rulesOverview')

// 查询表单
const queryForm = reactive({
  name: '',
  enabled: undefined as boolean | undefined,
  condition_type: '' as string
})

const form = reactive<CreateSyncPolicyRequest & { rules: RuleForm[] }>({
  name: '',
  remarks: '',
  scope: 'all',
  domain_id: 1,
  enabled: true,
  rules: []
})

const rules = {
  name: [{ required: true, message: '请输入策略名称', trigger: 'blur' }],
  domain_id: [{ required: true, message: '请选择所属域', trigger: 'change' }]
}

const loadPolicies = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pageSize.value,
      offset: (currentPage.value - 1) * pageSize.value
    }
    if (queryForm.name) params.name = queryForm.name
    if (queryForm.condition_type) params.condition_type = queryForm.condition_type

    // 根据顶部Tab过滤
    if (activeTab.value === 'enabled') params.enabled = true
    if (activeTab.value === 'disabled') params.enabled = false

    const res = await getSyncPolicies(params)
    policies.value = res.items || []
    total.value = res.total || 0
  } catch (e) {
    console.error(e)
    ElMessage.error('加载策略列表失败')
  } finally {
    loading.value = false
  }
}

const handleTabChange = (tab: string) => {
  currentPage.value = 1
  loadPolicies()
}

const handleSelectionChange = (selection: SyncPolicy[]) => {
  selectedIds.value = selection.map(p => p.id)
}

const handleBatchCommand = async (command: string) => {
  if (selectedIds.value.length === 0) {
    ElMessage.warning('请先选择策略')
    return
  }

  try {
    if (command === 'batchEnable') {
      await ElMessageBox.confirm(`确定要批量启用 ${selectedIds.value.length} 个策略吗？`, '提示', { type: 'info' })
      for (const id of selectedIds.value) {
        await updateSyncPolicyStatus(id, true)
      }
      ElMessage.success('批量启用成功')
    } else if (command === 'batchDisable') {
      await ElMessageBox.confirm(`确定要批量禁用 ${selectedIds.value.length} 个策略吗？`, '提示', { type: 'warning' })
      for (const id of selectedIds.value) {
        await updateSyncPolicyStatus(id, false)
      }
      ElMessage.success('批量禁用成功')
    } else if (command === 'batchDelete') {
      await ElMessageBox.confirm(`确定要批量删除 ${selectedIds.value.length} 个策略吗？此操作不可恢复。`, '提示', { type: 'warning' })
      for (const id of selectedIds.value) {
        await deleteSyncPolicy(id)
      }
      ElMessage.success('批量删除成功')
    }
    tableRef.value?.clearSelection()
    loadPolicies()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error('批量操作失败')
    }
  }
}

const handleExport = async () => {
  ElMessage.info('导出功能开发中')
}

const loadDomains = async () => {
  try {
    const res = await getDomains()
    domains.value = res.items || []
  } catch (e) {
    console.error(e)
  }
}

const loadProjects = async () => {
  try {
    const res = await getProjects()
    projects.value = res.items || []
  } catch (e) {
    console.error(e)
  }
}

const getScopeText = (scope?: string) => {
  const map: Record<string, string> = {
    'all': '全部云账号',
    'specified': '指定云账号',
    'resource_type': '指定资源类型'
  }
  return map[scope || ''] || scope || ''
}

const getConditionText = (type?: string) => {
  const map: Record<string, string> = {
    'all_match': '全部匹配标签',
    'any_match': '至少一个标签',
    'key_match': '根据标签Key匹配'
  }
  return map[type || ''] || type || ''
}

const getMappingText = (mapping?: string) => {
  const map: Record<string, string> = {
    'specify_project': '指定项目',
    'specify_name': '根据名称映射'
  }
  return map[mapping || ''] || mapping || ''
}

const formatDate = (date?: string) => {
  if (!date) return ''
  return new Date(date).toLocaleString('zh-CN')
}

const showCreateDialog = () => {
  isEdit.value = false
  resetForm()
  showDialog.value = true
}

const resetForm = () => {
  form.name = ''
  form.remarks = ''
  form.scope = 'all'
  form.domain_id = domains.value[0]?.id || 1
  form.enabled = true
  form.rules = []
}

const addRule = () => {
  form.rules.push({
    condition_type: 'all_match',
    resource_mapping: 'specify_project',
    tags: []
  })
}

const removeRule = (index: number) => {
  form.rules.splice(index, 1)
}

const addTag = (ruleIndex: number) => {
  form.rules[ruleIndex].tags.push({ tag_key: '', tag_value: '' })
}

const removeTag = (ruleIndex: number, tagIndex: number) => {
  form.rules[ruleIndex].tags.splice(tagIndex, 1)
}

const handleView = async (row: SyncPolicy) => {
  try {
    const detail = await getSyncPolicy(row.id)
    currentPolicy.value = detail
    showDetailDrawer.value = true
  } catch (e) {
    console.error(e)
    ElMessage.error('获取详情失败')
  }
}

const handleEdit = async (row: SyncPolicy) => {
  try {
    const detail = await getSyncPolicy(row.id)
    isEdit.value = true
    form.name = detail.name
    form.remarks = detail.remarks || ''
    form.scope = detail.scope
    form.domain_id = detail.domain_id
    form.enabled = detail.enabled
    form.rules = detail.rules?.map(r => ({
      condition_type: r.condition_type,
      resource_mapping: r.resource_mapping,
      target_project_id: r.target_project_id,
      target_project_name: r.target_project_name,
      tags: r.tags?.map(t => ({ tag_key: t.tag_key, tag_value: t.tag_value })) || []
    })) || []
    currentPolicy.value = detail
    showDialog.value = true
  } catch (e) {
    console.error(e)
    ElMessage.error('获取策略详情失败')
  }
}

const handleToggle = async (row: SyncPolicy) => {
  const newEnabled = row.enabled
  try {
    await updateSyncPolicyStatus(row.id, newEnabled)
    ElMessage.success(`${newEnabled ? '启用' : '禁用'}成功`)
  } catch (e) {
    console.error(e)
    row.enabled = !newEnabled // 恢复原状态
    ElMessage.error('更新状态失败')
  }
}

const handleExecute = async (row: SyncPolicy) => {
  ElMessage.info('执行策略功能开发中')
}

const handleCopy = async (row: SyncPolicy) => {
  try {
    const detail = await getSyncPolicy(row.id)
    isEdit.value = false
    form.name = detail.name + ' (副本)'
    form.remarks = detail.remarks || ''
    form.scope = detail.scope
    form.domain_id = detail.domain_id
    form.enabled = false // 复制的策略默认禁用
    form.rules = detail.rules?.map(r => ({
      condition_type: r.condition_type,
      resource_mapping: r.resource_mapping,
      target_project_id: r.target_project_id,
      target_project_name: r.target_project_name,
      tags: r.tags?.map(t => ({ tag_key: t.tag_key, tag_value: t.tag_value })) || []
    })) || []
    showDialog.value = true
    ElMessage.success('已复制策略，请修改后保存')
  } catch (e) {
    console.error(e)
    ElMessage.error('复制策略失败')
  }
}

const handleDelete = async (row: SyncPolicy) => {
  try {
    await ElMessageBox.confirm(`确定要删除策略 "${row.name}" 吗？此操作不可恢复。`, '提示', { type: 'warning' })
    await deleteSyncPolicy(row.id)
    ElMessage.success('删除成功')
    loadPolicies()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
      ElMessage.error('删除失败')
    }
  }
}

const handleDropdownCommand = (command: string, row: SyncPolicy) => {
  switch (command) {
    case 'execute':
      handleExecute(row)
      break
    case 'edit':
      handleEdit(row)
      break
    case 'copy':
      handleCopy(row)
      break
    case 'toggle':
      row.enabled = !row.enabled
      handleToggle(row)
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

const handleDrawerCommand = (command: string, row: SyncPolicy) => {
  if (command === 'copy') {
    handleCopy(row)
    showDetailDrawer.value = false
  } else if (command === 'delete') {
    handleDelete(row)
    showDetailDrawer.value = false
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      const submitData = {
        name: form.name,
        remarks: form.remarks,
        scope: form.scope,
        domain_id: form.domain_id,
        enabled: form.enabled,
        status: form.enabled ? 'active' : 'inactive',
        rules: form.rules.map(r => ({
          condition_type: r.condition_type,
          resource_mapping: r.resource_mapping,
          target_project_id: r.target_project_id,
          target_project_name: r.target_project_name,
          tags: r.tags
        }))
      }

      if (isEdit.value && currentPolicy.value) {
        await updateSyncPolicy(currentPolicy.value.id, submitData)
        ElMessage.success('更新成功')
      } else {
        await createSyncPolicy(submitData)
        ElMessage.success('创建成功')
      }
      showDialog.value = false
      loadPolicies()
    } catch (e) {
      console.error(e)
      ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
    } finally {
      submitting.value = false
    }
  })
}

const resetQuery = () => {
  queryForm.name = ''
  queryForm.condition_type = ''
  queryForm.enabled = undefined
  activeTab.value = 'all'
  currentPage.value = 1
  loadPolicies()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadPolicies()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadPolicies()
}

onMounted(() => {
  loadPolicies()
  loadDomains()
  loadProjects()
})
</script>

<style scoped>
.sync-policies-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.toolbar {
  display: flex;
  gap: 8px;
  align-items: center;
}

.status-tabs {
  margin-bottom: 16px;
}

.filter-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}

.name-link {
  font-weight: 500;
}

/* 抽屉顶部区域 */
.drawer-header {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--el-border-color);
}

.policy-icon {
  margin-right: 16px;
}

.policy-info h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.policy-tags {
  display: flex;
  gap: 8px;
  margin-top: 8px;
  align-items: center;
}

.quick-actions {
  margin-left: auto;
  display: flex;
  gap: 8px;
}

/* 规则编辑器样式 */
.rules-container {
  width: 100%;
}

.rule-item {
  margin-bottom: 16px;
}

.rule-card {
  margin-bottom: 0;
}

.rule-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tags-container {
  width: 100%;
}

.tag-item {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.rules-detail {
  margin-top: 16px;
}
</style>