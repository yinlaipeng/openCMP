<template>
  <div class="groups-container">
    <div class="page-header">
      <h2>用户组</h2>
      <div class="toolbar">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          新建用户组
        </el-button>
        <el-button :disabled="selectedGroups.length === 0" @click="handleBatchDelete">
          <el-icon><Delete /></el-icon>
          批量删除
        </el-button>
      </div>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" @submit.prevent="loadGroups">
        <el-form-item label="搜索">
          <el-select v-model="searchField" placeholder="选择字段" style="width: 100px">
            <el-option label="名称" value="name" />
            <el-option label="描述" value="description" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input v-model="searchKeyword" placeholder="输入关键词" clearable style="width: 180px" @keyup.enter="loadGroups">
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadGroups">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card>
      <el-table
        :data="groups"
        v-loading="loading"
        border
        stripe
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="名称" min-width="180" show-overflow-tooltip sortable>
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)" class="name-link">
              <span>{{ row.name }}</span>
            </el-button>
            <br>
            <small style="color: #999;">{{ row.description || '-' }}</small>
          </template>
        </el-table-column>
        <el-table-column label="所属域" width="120" sortable>
          <template #default="{ row }">
            <span>{{ getDomainName(row.domain_id) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="用户数" width="100" sortable>
          <template #default="{ row }">
            <span>{{ row.user_count || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column label="项目数" width="100" sortable>
          <template #default="{ row }">
            <span>{{ row.project_count || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="handleManageProjects(row)">管理项目</el-button>
            <el-button size="small" type="primary" link @click="handleManageUsers(row)">管理用户</el-button>
            <el-button size="small" type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.limit"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next"
        class="pagination"
      />
    </el-card>

    <!-- 添加/编辑用户组对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑用户组' : '新建用户组'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="用户组名" prop="name">
          <el-input v-model="form.name" placeholder="请输入用户组名" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="备注" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入用户组描述" />
        </el-form-item>
        <el-form-item label="域" prop="domain_id" v-if="!isEdit">
          <el-select v-model="form.domain_id" placeholder="请选择域" style="width: 100%">
            <el-option
              v-for="domain in allDomains"
              :key="domain.id"
              :label="domain.name"
              :value="domain.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 详情抽屉 -->
    <GroupDetailDrawer
      v-model="detailDrawerVisible"
      :group-id="currentGroupId"
      @refresh="loadGroups"
    />

    <!-- 批量删除确认弹窗 -->
    <BatchDeleteModal
      v-model="batchDeleteVisible"
      :groups="selectedGroups"
      :domains="allDomains"
      :submitting="batchDeleting"
      @confirm="handleBatchDeleteConfirm"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Refresh, Plus, Delete, Download, Setting, ArrowDown, Search } from '@element-plus/icons-vue'
import { Group, Domain } from '@/types/iam'
import {
  getGroups,
  createGroup,
  updateGroup,
  deleteGroup,
  getDomains
} from '@/api/iam'
import GroupDetailDrawer from './components/GroupDetailDrawer.vue'
import BatchDeleteModal from './components/BatchDeleteModal.vue'

const groups = ref<Group[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const detailDrawerVisible = ref(false)
const batchDeleteVisible = ref(false)
const batchDeleting = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const currentGroupId = ref(0)
const selectedGroups = ref<Group[]>([])
const allDomains = ref<Domain[]>([])
const formRef = ref<FormInstance>()

// 搜索相关
const searchField = ref('name')
const searchKeyword = ref('')

const currentFieldLabel = computed(() => {
  const labels: Record<string, string> = {
    name: '名称',
    description: '描述'
  }
  return labels[searchField.value] || '名称'
})

const searchPlaceholder = computed(() => {
  const placeholders: Record<string, string> = {
    name: '请输入名称关键词...',
    description: '请输入描述关键词...'
  }
  return placeholders[searchField.value] || '请输入名称关键词...'
})

const filterForm = reactive({
  name: ''
})

const pagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

const form = reactive({
  name: '',
  description: '',
  domain_id: 1
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入用户组名', trigger: 'blur' }]
}

const handleFieldChange = (field: string) => {
  searchField.value = field
  searchKeyword.value = ''
}

const loadGroups = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit,
      details: 'true'
    }

    // 根据搜索字段添加参数
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }

    const res = await getGroups(params)
    groups.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载用户组列表失败')
  } finally {
    loading.value = false
  }
}

const loadAllDomains = async () => {
  try {
    const res = await getDomains({ limit: 100 })
    allDomains.value = res.items || []
  } catch (e: any) {
    console.error(e)
  }
}

const handleRefresh = () => {
  loadGroups()
}

const resetFilter = () => {
  searchKeyword.value = ''
  searchField.value = 'name'
  pagination.page = 1
  loadGroups()
}

const handleSelectionChange = (selection: Group[]) => {
  selectedGroups.value = selection
}

const handleCreate = async () => {
  isEdit.value = false
  form.name = ''
  form.description = ''
  form.domain_id = 1
  await loadAllDomains()
  dialogVisible.value = true
}

const handleEdit = (row: Group) => {
  isEdit.value = true
  currentGroupId.value = row.id
  form.name = row.name
  form.description = row.description || ''
  form.domain_id = row.domain_id
  dialogVisible.value = true
}

const handleView = (row: Group) => {
  currentGroupId.value = row.id
  detailDrawerVisible.value = true
}

const handleManageProjects = (row: Group) => {
  currentGroupId.value = row.id
  detailDrawerVisible.value = true
}

const handleManageUsers = (row: Group) => {
  currentGroupId.value = row.id
  detailDrawerVisible.value = true
}

const handleDelete = async (row: Group) => {
  try {
    await ElMessageBox.confirm('此操作将永久删除该用户组，是否继续？', '删除确认', { type: 'warning' })
    await deleteGroup(row.id)
    ElMessage.success('删除成功')
    loadGroups()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const handleBatchDelete = () => {
  if (selectedGroups.value.length === 0) {
    ElMessage.warning('请先选择要删除的用户组')
    return
  }
  batchDeleteVisible.value = true
}

const handleBatchDeleteConfirm = async () => {
  batchDeleting.value = true
  try {
    // 逐个删除选中的用户组
    for (const group of selectedGroups.value) {
      await deleteGroup(group.id)
    }
    ElMessage.success(`成功删除 ${selectedGroups.value.length} 个用户组`)
    batchDeleteVisible.value = false
    selectedGroups.value = []
    loadGroups()
  } catch (e: any) {
    ElMessage.error(e.message || '批量删除失败')
  } finally {
    batchDeleting.value = false
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (isEdit.value) {
        await updateGroup(currentGroupId.value, form)
        ElMessage.success('更新成功')
      } else {
        await createGroup(form)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadGroups()
    } catch (e: any) {
      ElMessage.error(e.message || (isEdit.value ? '更新失败' : '创建失败'))
    } finally {
      submitting.value = false
    }
  })
}

const handleDownload = () => {
  ElMessage.info('导出功能待实现')
}

const handleSettings = () => {
  ElMessage.info('设置功能待实现')
}

const getDomainName = (domainId: number) => {
  const domain = allDomains.value.find(d => d.id === domainId);
  return domain ? domain.name : `域#${domainId}`;
}

onMounted(async () => {
  await loadAllDomains()
  loadGroups()
})
</script>

<style scoped>
.groups-container { padding: 20px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.page-header h2 { margin: 0; font-size: 18px; font-weight: 600; }
.filter-card { margin-bottom: 20px; }
.pagination { margin-top: 20px; text-align: right; }
.name-link { color: #409eff; cursor: pointer; }
</style>