<template>
  <div class="domains-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">域</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建域
          </el-button>
        </div>
      </template>

      <!-- 筛选条件 -->
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="域名称">
          <el-input v-model="filterForm.name" placeholder="请输入域名称" clearable @keyup.enter="loadDomains" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterForm.enabled" placeholder="全部" clearable style="width: 100px">
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadDomains">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 域列表 -->
      <el-table :data="domains" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="域名称" min-width="180" show-overflow-tooltip />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="enabled" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="320" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">详情</el-button>
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button 
              size="small" 
              :type="row.enabled ? 'warning' : 'success'" 
              @click="handleToggleEnable(row)"
            >
              {{ row.enabled ? '禁用' : '启用' }}
            </el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
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
        @size-change="loadDomains"
        @current-change="loadDomains"
        class="pagination"
      />
    </el-card>

    <!-- 添加/编辑域对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑域' : '新建域'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="域名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入域名称" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入域描述" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 域详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="域详情" width="900px">
      <el-descriptions :column="2" border v-if="currentDomain">
        <el-descriptions-item label="ID">{{ currentDomain.id }}</el-descriptions-item>
        <el-descriptions-item label="域名称">{{ currentDomain.name }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentDomain.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentDomain.enabled ? 'success' : 'info'" size="small">
            {{ currentDomain.enabled ? '启用' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentDomain.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentDomain.updated_at) }}</el-descriptions-item>
      </el-descriptions>

      <el-tabs v-model="detailTab" style="margin-top: 20px">
        <el-tab-pane label="用户" name="users">
          <el-table :data="domainUsers" v-loading="usersLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="用户名" />
            <el-table-column prop="display_name" label="显示名" />
            <el-table-column prop="email" label="邮箱" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="用户组" name="groups">
          <el-table :data="domainGroups" v-loading="groupsLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="用户组名" />
            <el-table-column prop="description" label="描述" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="项目" name="projects">
          <el-table :data="domainProjects" v-loading="projectsLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="项目名" />
            <el-table-column prop="description" label="描述" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="角色" name="roles">
          <el-table :data="domainRoles" v-loading="rolesLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="角色名" />
            <el-table-column prop="display_name" label="显示名" />
            <el-table-column prop="type" label="类型" width="100">
              <template #default="{ row }">
                <el-tag size="small" :type="row.type === 'system' ? 'warning' : 'success'">
                  {{ row.type === 'system' ? '系统' : '自定义' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>

      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Domain } from '@/types/iam'
import {
  getDomains,
  createDomain,
  updateDomain,
  deleteDomain,
  enableDomain,
  disableDomain,
  getDomainUsers,
  getDomainGroups,
  getDomainProjects,
  getDomainRoles
} from '@/api/iam'

const domains = ref<Domain[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const detailTab = ref('users')
const isEdit = ref(false)
const submitting = ref(false)
const currentDomain = ref<Domain | null>(null)
const domainUsers = ref<Domain[]>([])
const domainGroups = ref<Domain[]>([])
const domainProjects = ref<Domain[]>([])
const domainRoles = ref<Domain[]>([])
const usersLoading = ref(false)
const groupsLoading = ref(false)
const projectsLoading = ref(false)
const rolesLoading = ref(false)
const formRef = ref<FormInstance>()

const filterForm = reactive({
  name: '',
  enabled: undefined as boolean | undefined
})

const pagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

const form = reactive({
  name: '',
  description: '',
  enabled: true
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入域名称', trigger: 'blur' }]
}

const loadDomains = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit
    }
    if (filterForm.name) params.keyword = filterForm.name
    if (filterForm.enabled !== undefined) params.enabled = filterForm.enabled

    const res = await getDomains(params)
    domains.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载域列表失败')
  } finally {
    loading.value = false
  }
}

const loadDomainUsers = async (domainId: number) => {
  usersLoading.value = true
  try {
    const res = await getDomainUsers(domainId)
    domainUsers.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    usersLoading.value = false
  }
}

const loadDomainGroups = async (domainId: number) => {
  groupsLoading.value = true
  try {
    const res = await getDomainGroups(domainId)
    domainGroups.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    groupsLoading.value = false
  }
}

const loadDomainProjects = async (domainId: number) => {
  projectsLoading.value = true
  try {
    const res = await getDomainProjects(domainId)
    domainProjects.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    projectsLoading.value = false
  }
}

const loadDomainRoles = async (domainId: number) => {
  rolesLoading.value = true
  try {
    const res = await getDomainRoles(domainId)
    domainRoles.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    rolesLoading.value = false
  }
}

const resetFilter = () => {
  filterForm.name = ''
  filterForm.enabled = undefined
  pagination.page = 1
  loadDomains()
}

const handleCreate = () => {
  isEdit.value = false
  form.name = ''
  form.description = ''
  form.enabled = true
  dialogVisible.value = true
}

const handleEdit = (row: Domain) => {
  isEdit.value = true
  currentDomain.value = row
  form.name = row.name
  form.description = row.description || ''
  form.enabled = row.enabled
  dialogVisible.value = true
}

const handleView = async (row: Domain) => {
  currentDomain.value = row
  detailTab.value = 'users'
  await loadDomainUsers(row.id)
  await loadDomainGroups(row.id)
  await loadDomainProjects(row.id)
  await loadDomainRoles(row.id)
  detailDialogVisible.value = true
}

const handleToggleEnable = async (row: Domain) => {
  try {
    const action = row.enabled ? '禁用' : '启用'
    await ElMessageBox.confirm(`确定要${action}该域吗？`, '提示', { type: 'warning' })

    if (row.enabled) {
      await disableDomain(row.id)
    } else {
      await enableDomain(row.id)
    }

    ElMessage.success(`${action}成功`)
    loadDomains()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || `${action}失败`)
    }
  }
}

const handleDelete = async (row: Domain) => {
  try {
    await ElMessageBox.confirm('确定要删除该域吗？', '提示', { type: 'warning' })
    await deleteDomain(row.id)
    ElMessage.success('删除成功')
    loadDomains()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (isEdit.value) {
        await updateDomain(currentDomain.value.id, form)
        ElMessage.success('更新成功')
      } else {
        await createDomain(form)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadDomains()
    } catch (e: any) {
      ElMessage.error(e.message || (isEdit.value ? '更新失败' : '创建失败'))
    } finally {
      submitting.value = false
    }
  })
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  loadDomains()
})
</script>

<style scoped>
.domains-page {
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

.filter-form {
  margin-bottom: 16px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
