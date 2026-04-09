<template>
  <div class="sync-policies-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">同步策略</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建同步策略
          </el-button>
        </div>
      </template>

      <el-table :data="policies" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称/备注" width="200">
          <template #default="{ row }">
            <div>
              <div>{{ row.name }}</div>
              <div v-if="row.remarks" class="remarks">{{ row.remarks }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
              {{ row.status === 'active' ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="启用状态" width="120">
          <template #default="{ row }">
            <el-switch v-model="row.enabled" @change="toggleEnabled(row)" />
          </template>
        </el-table-column>
        <el-table-column label="规则" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="manageRules(row)">管理规则</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="scope" label="应用范围" width="150" />
        <el-table-column prop="domain_id" label="所属域" width="100" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-dropdown>
              <el-button size="small">
                更多 <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="viewDetails(row)">详情</el-dropdown-item>
                  <el-dropdown-item @click="viewOperationLogs(row)">操作日志</el-dropdown-item>
                  <el-dropdown-item divided @click="toggleEnabled(row)">
                    {{ row.enabled ? '禁用' : '启用' }}
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleEdit(row)">编辑</el-dropdown-item>
                  <el-dropdown-item @click="handleDelete(row)" divided>删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.currentPage"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        class="pagination"
      />
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑同步策略' : '新建同步策略'"
      width="800px"
      :fullscreen="isMobile"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="域" prop="domain_id">
          <el-select v-model="form.domain_id" placeholder="请选择域" style="width: 100%">
            <el-option
              v-for="domain in domains"
              :key="domain.id"
              :label="domain.name"
              :value="domain.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入策略名称" />
        </el-form-item>

        <el-form-item label="备注" prop="remarks">
          <el-input v-model="form.remarks" type="textarea" placeholder="请输入备注" />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio label="active">启用</el-radio>
            <el-radio label="inactive">停用</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="启用状态" prop="enabled">
          <el-switch v-model="form.enabled" />
        </el-form-item>

        <el-form-item label="规则配置" prop="rules">
          <RuleConfiguration v-model="form.rules" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 规则管理对话框 -->
    <el-dialog
      v-model="rulesDialogVisible"
      title="管理规则"
      width="900px"
      :fullscreen="isMobile"
    >
      <el-table :data="currentRules" border stripe>
        <el-table-column prop="condition_type" label="规则" width="200">
          <template #default="{ row }">
            <el-tag>
              {{ getConditionTypeName(row.condition_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_mapping" label="资源映射" width="200">
          <template #default="{ row }">
            <el-tag type="info">
              {{ getResourceMappingName(row.resource_mapping) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="标签" min-width="200">
          <template #default="{ row }">
            <el-tag
              v-for="(tag, index) in row.tags || []"
              :key="index"
              size="small"
              style="margin-right: 5px; margin-bottom: 5px;"
            >
              {{ tag.key }}: {{ tag.value }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row, $index }">
            <el-button size="small" @click="editRule($index)">修改</el-button>
            <el-button size="small" type="danger" @click="deleteRule($index)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <template #footer>
        <el-button @click="rulesDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="saveRules">保存</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog
      v-model="detailsDialogVisible"
      title="详情"
      width="800px"
      :fullscreen="isMobile"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ currentPolicy.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ currentPolicy.name }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ currentPolicy.remarks }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentPolicy.status === 'active' ? 'success' : 'danger'">
            {{ currentPolicy.status === 'active' ? '启用' : '停用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="启用状态">
          <el-tag :type="currentPolicy.enabled ? 'success' : 'info'">
            {{ currentPolicy.enabled ? '启用' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="应用范围" :span="2">{{ currentPolicy.scope }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ currentPolicy.domain_id }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentPolicy.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentPolicy.updated_at) }}</el-descriptions-item>
      </el-descriptions>

      <div class="section">
        <h3>规则配置</h3>
        <div v-for="(rule, index) in currentPolicy.rules || []" :key="index" class="rule-detail">
          <h4>规则 {{ index + 1 }}</h4>
          <p><strong>条件类型:</strong> {{ getConditionTypeName(rule.condition_type) }}</p>
          <p><strong>资源映射:</strong> {{ getResourceMappingName(rule.resource_mapping) }}</p>
          <p v-if="rule.target_project_id"><strong>目标项目ID:</strong> {{ rule.target_project_id }}</p>
          <p v-if="rule.target_project_name"><strong>目标项目名称:</strong> {{ rule.target_project_name }}</p>
          <p><strong>标签:</strong></p>
          <ul>
            <li v-for="(tag, tagIndex) in rule.tags || []" :key="tagIndex">
              {{ tag.key }}: {{ tag.value }}
            </li>
          </ul>
        </div>
      </div>

      <template #footer>
        <el-button @click="detailsDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, ElDescriptions, ElDescriptionsItem } from 'element-plus'
import { Plus, ArrowDown } from '@element-plus/icons-vue'
import { getSyncPolicies, createSyncPolicy, updateSyncPolicy, deleteSyncPolicy, updateSyncPolicyStatus } from '@/api/sync-policy'
import { getDomains } from '@/api/iam'
import type { SyncPolicy, CreateSyncPolicyRequest, Rule } from '@/types/sync-policy'
import type { Domain } from '@/types'
import RuleConfiguration from '@/components/RuleConfiguration.vue'

const policies = ref<SyncPolicy[]>([])
const domains = ref<Domain[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const rulesDialogVisible = ref(false)
const detailsDialogVisible = ref(false)
const currentRules = ref<Rule[]>([])
const currentPolicy = ref<SyncPolicy>({} as SyncPolicy)

const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

const form = reactive<CreateSyncPolicyRequest>({
  name: '',
  remarks: '',
  status: 'active',
  enabled: true,
  rules: [],
  scope: '',
  domain_id: 1
})

const rules = {
  name: [{ required: true, message: '请输入策略名称', trigger: 'blur' }],
  domain_id: [{ required: true, message: '请选择所属域', trigger: 'change' }],
  rules: [{ required: true, message: '请至少配置一条规则', trigger: 'change' }]
}

// 检测移动端
const isMobile = computed(() => {
  return window.innerWidth < 768
})

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    // 并行加载数据
    await Promise.all([
      loadSyncPolicies(),
      loadDomains()
    ])
  } catch (e) {
    console.error(e)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

const loadSyncPolicies = async () => {
  const params = {
    page: pagination.currentPage,
    page_size: pagination.pageSize
  }
  const res = await getSyncPolicies(params)
  policies.value = res.items || []
  pagination.total = res.total || 0
}

const loadDomains = async () => {
  try {
    const res = await getDomains()
    domains.value = res.items || []
    // 如果domain_id没有选择且有可用的域，则默认选择第一个
    if (!form.domain_id && domains.value.length > 0) {
      form.domain_id = domains.value[0].id
    }
  } catch (e) {
    console.error('Failed to load domains:', e)
    ElMessage.error('加载域列表失败')
  }
}

const handleCreate = () => {
  isEdit.value = false
  Object.assign(form, {
    name: '',
    remarks: '',
    status: 'active',
    enabled: true,
    rules: [],
    scope: '',
    domain_id: domains.value.length > 0 ? domains.value[0].id : 1
  })
  dialogVisible.value = true
}

const handleEdit = (row: SyncPolicy) => {
  isEdit.value = true
  Object.assign(form, { ...row })
  dialogVisible.value = true
}

const handleDelete = async (row: SyncPolicy) => {
  try {
    await ElMessageBox.confirm(`确定要删除同步策略 "${row.name}" 吗？`, '提示', {
      type: 'warning'
    })
    await deleteSyncPolicy(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
      ElMessage.error('删除失败')
    }
  }
}

const toggleEnabled = async (row: SyncPolicy) => {
  try {
    await updateSyncPolicyStatus(row.id, !row.enabled)
    row.enabled = !row.enabled
    ElMessage.success(`${row.enabled ? '启用' : '禁用'}成功`)
    loadData()
  } catch (e) {
    console.error(e)
    ElMessage.error('更新状态失败')
    // 恢复开关状态
    row.enabled = !row.enabled
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      let res
      if (isEdit.value) {
        res = await updateSyncPolicy(form.id as number, form)
        ElMessage.success('更新成功')
      } else {
        res = await createSyncPolicy(form)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadData()
    } catch (e) {
      console.error(e)
      ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
    } finally {
      submitting.value = false
    }
  })
}

const manageRules = (row: SyncPolicy) => {
  currentPolicy.value = row
  currentRules.value = [...(row.rules || [])]
  rulesDialogVisible.value = true
}

const viewDetails = (row: SyncPolicy) => {
  currentPolicy.value = row
  detailsDialogVisible.value = true
}

const viewOperationLogs = (row: SyncPolicy) => {
  ElMessage.info('操作日志功能即将实现')
}

const getConditionTypeName = (type: string) => {
  const map: Record<string, string> = {
    'all_match': '全部匹配标签',
    'any_match': '至少一个标签',
    'key_match': '根据标签key匹配'
  }
  return map[type] || type
}

const getResourceMappingName = (mapping: string) => {
  const map: Record<string, string> = {
    'specify_project': '指定项目',
    'specify_name': '指定名称'
  }
  return map[mapping] || mapping
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

const editRule = (index: number) => {
  ElMessage.info('规则编辑功能请在主表单中操作')
}

const deleteRule = (index: number) => {
  currentRules.value.splice(index, 1)
  // 更新主策略的规则
  currentPolicy.value.rules = [...currentRules.value]
}

const saveRules = () => {
  rulesDialogVisible.value = false
  // 规则已在主表单中编辑，无需额外保存
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  loadData()
}

const handleCurrentChange = (page: number) => {
  pagination.currentPage = page
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.sync-policies-page {
  height: 100%;
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

.remarks {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
}

.section {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.rule-detail {
  margin: 10px 0;
  padding: 10px;
  background-color: #f5f5f5;
  border-radius: 4px;
}

.rule-detail h4 {
  margin-top: 0;
  color: #333;
}
</style>