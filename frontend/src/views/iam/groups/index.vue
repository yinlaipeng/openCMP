<template>
  <div class="groups-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">用户组</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建用户组
          </el-button>
        </div>
      </template>

      <!-- 筛选条件 -->
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="用户组名">
          <el-input v-model="filterForm.name" placeholder="请输入用户组名" clearable @keyup.enter="loadGroups" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadGroups">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 用户组列表 -->
      <el-table :data="groups" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="用户组名" min-width="180" show-overflow-tooltip />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="domain_id" label="所属域" width="120">
          <template #default="{ row }">
            <span>{{ row.domain_id === 1 ? 'Default' : row.domain_id }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">详情</el-button>
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
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
        @size-change="loadGroups"
        @current-change="loadGroups"
        class="pagination"
      />
    </el-card>

    <!-- 添加/编辑用户组对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑用户组' : '新建用户组'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="用户组名" prop="name">
          <el-input v-model="form.name" placeholder="请输入用户组名" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入用户组描述" />
        </el-form-item>
        <el-form-item label="所属域" prop="domain_id" v-if="!isEdit">
          <el-select v-model="form.domain_id" placeholder="请选择所属域" style="width: 100%">
            <el-option label="Default" :value="1" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 用户组详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="用户组详情" width="800px">
      <el-descriptions :column="2" border v-if="currentGroup">
        <el-descriptions-item label="ID">{{ currentGroup.id }}</el-descriptions-item>
        <el-descriptions-item label="用户组名">{{ currentGroup.name }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentGroup.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="所属域">
          <span>{{ currentGroup.domain_id === 1 ? 'Default' : currentGroup.domain_id }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentGroup.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentGroup.updated_at) }}</el-descriptions-item>
      </el-descriptions>

      <el-tabs v-model="detailTab" style="margin-top: 20px">
        <el-tab-pane label="组成员" name="users">
          <div class="tab-toolbar">
            <span>用户组成员列表</span>
            <el-button size="small" type="primary" @click="handleAddUser">添加用户</el-button>
          </div>
          <el-table :data="groupUsers" v-loading="usersLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="用户名" />
            <el-table-column prop="display_name" label="显示名" />
            <el-table-column prop="email" label="邮箱" />
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button size="small" type="danger" @click="handleRemoveUser(row)">移除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>

      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 添加用户对话框 -->
    <el-dialog v-model="addUserDialogVisible" title="添加用户到用户组" width="600px">
      <el-form :model="addUserForm" label-width="100px">
        <el-form-item label="选择用户" required>
          <el-select
            v-model="addUserForm.user_id"
            placeholder="请选择用户"
            style="width: 100%"
            filterable
          >
            <el-option
              v-for="item in allUsers"
              :key="item.id"
              :label="`${item.name} (${item.display_name})`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addUserDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAddUserSubmit" :loading="addUserSubmitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Group, User } from '@/types/iam'
import {
  getGroups,
  createGroup,
  updateGroup,
  deleteGroup,
  getGroupUsers,
  addUserToGroup,
  removeUserFromGroup,
  getUsers
} from '@/api/iam'

const groups = ref<Group[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const addUserDialogVisible = ref(false)
const detailTab = ref('users')
const isEdit = ref(false)
const submitting = ref(false)
const addUserSubmitting = ref(false)
const currentGroupId = ref(0)
const currentGroup = ref<Group | null>(null)
const groupUsers = ref<User[]>([])
const allUsers = ref<User[]>([])
const usersLoading = ref(false)
const formRef = ref<FormInstance>()

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

const addUserForm = reactive({
  user_id: 0
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入用户组名', trigger: 'blur' }]
}

const loadGroups = async () => {
  loading.value = true
  try {
    const params: any = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit
    }
    if (filterForm.name) params.keyword = filterForm.name

    const res = await getGroups(params)
    groups.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载用户组列表失败')
  } finally {
    loading.value = false
  }
}

const loadGroupUsers = async (groupId: number) => {
  usersLoading.value = true
  try {
    const res = await getGroupUsers(groupId)
    groupUsers.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    usersLoading.value = false
  }
}

const loadAllUsers = async () => {
  try {
    const res = await getUsers({ limit: 100 })
    allUsers.value = res.items || []
  } catch (e: any) {
    console.error(e)
  }
}

const resetFilter = () => {
  filterForm.name = ''
  pagination.page = 1
  loadGroups()
}

const handleCreate = () => {
  isEdit.value = false
  form.name = ''
  form.description = ''
  form.domain_id = 1
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

const handleView = async (row: Group) => {
  currentGroup.value = row
  detailTab.value = 'users'
  await loadGroupUsers(row.id)
  detailDialogVisible.value = true
}

const handleAddUser = async () => {
  await loadAllUsers()
  addUserForm.user_id = 0
  addUserDialogVisible.value = true
}

const handleAddUserSubmit = async () => {
  if (!addUserForm.user_id) {
    ElMessage.error('请选择用户')
    return
  }

  addUserSubmitting.value = true
  try {
    await addUserToGroup(currentGroupId.value, addUserForm.user_id)
    ElMessage.success('添加用户成功')
    addUserDialogVisible.value = false
    await loadGroupUsers(currentGroupId.value)
  } catch (e: any) {
    ElMessage.error(e.message || '添加用户失败')
  } finally {
    addUserSubmitting.value = false
  }
}

const handleRemoveUser = async (row: User) => {
  try {
    await ElMessageBox.confirm('确定要移除该用户吗？', '提示', { type: 'warning' })
    await removeUserFromGroup(currentGroupId.value, row.id)
    ElMessage.success('移除用户成功')
    await loadGroupUsers(currentGroupId.value)
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '移除用户失败')
    }
  }
}

const handleDelete = async (row: Group) => {
  try {
    await ElMessageBox.confirm('确定要删除该用户组吗？', '提示', { type: 'warning' })
    await deleteGroup(row.id)
    ElMessage.success('删除成功')
    loadGroups()
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

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  loadGroups()
})
</script>

<style scoped>
.groups-page {
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

.tab-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
</style>
