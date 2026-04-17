<template>
  <div class="sync-policies-container">
    <div class="page-header">
      <h2>同步策略</h2>
      <el-button type="primary" @click="showCreateDialog">
        <el-icon><Plus /></el-icon>
        添加策略
      </el-button>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="queryForm" @submit.prevent="loadPolicies">
        <el-form-item label="名称">
          <el-input v-model="queryForm.name" placeholder="策略名称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryForm.enabled" placeholder="全部" clearable style="width: 120px">
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadPolicies">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 空状态 -->
    <EmptyState
      v-if="!loading && policies.length === 0"
      title="暂无同步策略"
      description="当前没有任何同步策略，点击下方按钮创建策略"
      :icon="Document"
      createButtonText="添加策略"
      @create="showCreateDialog"
    />

    <!-- 数据表格 -->
    <el-table
      :data="policies"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
      v-if="policies.length > 0 || loading"
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" min-width="150" />
      <el-table-column prop="rules" label="规则数量" width="100">
        <template #default="{ row }">
          <el-tag type="info">{{ row.rules?.length || 0 }} 条</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="remarks" label="备注" min-width="200" show-overflow-tooltip />
      <el-table-column prop="enabled" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.enabled ? 'success' : 'info'">
            {{ row.enabled ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180">
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
                <el-dropdown-item command="edit" :icon="EditPen">编辑</el-dropdown-item>
                <el-dropdown-item command="toggle" :icon="row.enabled ? CircleClose : CircleCheck">
                  {{ row.enabled ? '禁用' : '启用' }}
                </el-dropdown-item>
                <el-dropdown-item command="delete" :icon="Delete" divided>删除</el-dropdown-item>
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

    <!-- 详情对话框 -->
    <el-dialog v-model="showDetailDialog" title="策略详情" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="策略名称">{{ currentPolicy?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentPolicy?.enabled ? 'success' : 'info'">
            {{ currentPolicy?.enabled ? '启用' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="应用范围">{{ getScopeText(currentPolicy?.scope) }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ currentPolicy?.domain_id }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ currentPolicy?.remarks || '无' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentPolicy?.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentPolicy?.updated_at) }}</el-descriptions-item>
      </el-descriptions>

      <div class="rules-detail" v-if="currentPolicy?.rules?.length">
        <h4 style="margin: 20px 0 10px">同步规则</h4>
        <el-collapse>
          <el-collapse-item v-for="(rule, index) in currentPolicy.rules" :key="index" :title="`规则 ${index + 1}`">
            <el-descriptions :column="2" border size="small">
              <el-descriptions-item label="条件类型">{{ getConditionText(rule.condition_type) }}</el-descriptions-item>
              <el-descriptions-item label="资源映射">{{ getMappingText(rule.resource_mapping) }}</el-descriptions-item>
              <el-descriptions-item label="目标项目" v-if="rule.target_project_id">{{ rule.target_project_name || rule.target_project_id }}</el-descriptions-item>
              <el-descriptions-item label="匹配标签" :span="2">
                <el-tag v-for="tag in rule.tags" :key="tag.tag_key" style="margin-right: 8px">
                  {{ tag.tag_key }}: {{ tag.tag_value }}
                </el-tag>
              </el-descriptions-item>
            </el-descriptions>
          </el-collapse-item>
        </el-collapse>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { Plus, Delete, Document, ArrowDown, EditPen, CircleCheck, CircleClose } from '@element-plus/icons-vue'
import { getSyncPolicies, getSyncPolicy, createSyncPolicy, updateSyncPolicy, deleteSyncPolicy, updateSyncPolicyStatus } from '@/api/sync-policy'
import { getDomains } from '@/api/iam'
import { getProjects } from '@/api/project'
import type { SyncPolicy, CreateSyncPolicyRequest } from '@/types/sync-policy'
import EmptyState from '@/components/common/EmptyState.vue'

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
const showDetailDialog = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const currentPolicy = ref<SyncPolicy | null>(null)

// 分页数据
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 查询表单
const queryForm = reactive({
  name: '',
  enabled: undefined as boolean | undefined
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
    if (queryForm.enabled !== undefined) params.enabled = queryForm.enabled

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
    showDetailDialog.value = true
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
      tags: r.tags?.map(t => ({ tag_key: t.key, tag_value: t.value })) || []
    })) || []
    currentPolicy.value = detail
    showDialog.value = true
  } catch (e) {
    console.error(e)
    ElMessage.error('获取策略详情失败')
  }
}

const handleToggle = async (row: SyncPolicy) => {
  const newEnabled = !row.enabled
  try {
    await updateSyncPolicyStatus(row.id, newEnabled)
    row.enabled = newEnabled
    ElMessage.success(`${newEnabled ? '启用' : '禁用'}成功`)
  } catch (e) {
    console.error(e)
    ElMessage.error('更新状态失败')
  }
}

const handleDelete = async (row: SyncPolicy) => {
  try {
    await ElMessageBox.confirm(`确定要删除策略 "${row.name}" 吗？`, '提示', { type: 'warning' })
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
    case 'edit':
      handleEdit(row)
      break
    case 'toggle':
      handleToggle(row)
      break
    case 'delete':
      handleDelete(row)
      break
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
  queryForm.enabled = undefined
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

.filter-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}

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