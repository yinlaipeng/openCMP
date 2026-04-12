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
        <el-table-column label="名称" min-width="180" show-overflow-tooltip>
          <template #default="{ row }">
            <div>
              <el-button type="primary" link @click="handleView(row)" class="name-link">
                <span>{{ row.name }}</span>
              </el-button>
              <br>
              <small style="color: #999;">{{ row.description || '-' }}</small>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="所属域" width="120">
          <template #default="{ row }">
            <span>{{ getDomainName(row.domain_id) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="handleManageProjects(row)">管理项目</el-button>
            <el-button size="small" type="primary" link @click="handleManageUsers(row)">管理用户</el-button>
            <el-dropdown trigger="click" @command="(cmd) => handleCommand(cmd, row)">
              <el-button size="small" type="primary" link>
                更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="view">详情</el-dropdown-item>
                  <el-dropdown-item command="edit">编辑</el-dropdown-item>
                  <el-dropdown-item command="delete" divided>
                    <span style="color: #F56C6C">删除</span>
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

    <!-- 用户组详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="用户组详情" width="800px">
      <el-tabs v-model="detailTab" style="margin-top: 20px">
        <el-tab-pane label="详情" name="details">
          <el-descriptions :column="2" border v-if="currentGroup">
            <el-descriptions-item label="ID">{{ currentGroup.id }}</el-descriptions-item>
            <el-descriptions-item label="用户组名">{{ currentGroup.name }}</el-descriptions-item>
            <el-descriptions-item label="备注" :span="2">{{ currentGroup.description || '-' }}</el-descriptions-item>
            <el-descriptions-item label="所属域">
              <span>{{ getDomainName(currentGroup.domain_id) }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatDate(currentGroup.created_at) }}</el-descriptions-item>
            <el-descriptions-item label="更新时间">{{ formatDate(currentGroup.updated_at) }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="加入项目" name="projects">
          <div class="tab-toolbar">
            <span>用户组关联的项目</span>
            <el-button size="small" type="primary" @click="handleJoinProjectPopup(currentGroup!)">加入项目</el-button>
          </div>
          <el-table :data="groupProjects" v-loading="projectsLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="所属项目" min-width="150" />
            <el-table-column prop="role_name" label="角色" width="120">
              <template #default="{ row }">
                <span>{{ row.role_name || '未指定' }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="permissions" label="权限" width="120">
              <template #default="{ row }">
                <span>{{ row.permissions || '未指定' }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="domain_name" label="所属域" width="120">
              <template #default="{ row }">
                <span>{{ getDomainName(row.domain_id) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <el-button size="small" @click="handleEditProjectRole(row)">修改角色</el-button>
                <el-button size="small" type="danger" @click="handleRemoveProject(row)">移除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="组内用户" name="users">
          <div class="tab-toolbar">
            <span>用户组成员列表</span>
            <el-button size="small" type="primary" @click="handleAddGroupUser(currentGroup!)">添加用户</el-button>
          </div>
          <el-table :data="groupUsers" v-loading="usersLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="用户" min-width="120" />
            <el-table-column prop="enabled" label="启用状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.enabled ? 'success' : 'info'">
                  {{ row.enabled ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="domain_name" label="所属域" width="120">
              <template #default="{ row }">
                <span>{{ getDomainName(row.domain_id) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button size="small" type="danger" @click="handleRemoveUser(row)">移除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="操作日志" name="logs">
          <el-table :data="groupLogs" v-loading="logsLoading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="created_at" label="操作时间" width="180" />
            <el-table-column prop="resource_name" label="资源名称" min-width="150" />
            <el-table-column prop="resource_type" label="资源类型" width="120" />
            <el-table-column prop="operation_type" label="操作类型" width="120" />
            <el-table-column prop="service_type" label="服务类型" width="120" />
            <el-table-column prop="risk_level" label="风险级别" width="100" />
            <el-table-column prop="event_type" label="事件类型" width="120" />
            <el-table-column prop="result" label="结果" width="100">
              <template #default="{ row }">
                <el-tag :type="row.result === 'success' ? 'success' : 'danger'">
                  {{ row.result === 'success' ? '成功' : '失败' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="operator" label="发起人" width="120" />
            <el-table-column prop="project_name" label="所属项目" width="150" />
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
        <el-form-item label="选择用户" prop="user_id" :rules="[{ required: true, message: '请选择用户', trigger: 'change' }]">
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

    <!-- 添加项目对话框 -->
    <el-dialog v-model="addProjectDialogVisible" title="添加项目到用户组" width="600px">
      <el-form :model="addProjectForm" label-width="100px">
        <el-form-item label="选择项目" prop="project_id" :rules="[{ required: true, message: '请选择项目', trigger: 'change' }]">
          <el-select
            v-model="addProjectForm.project_id"
            placeholder="请选择项目"
            style="width: 100%"
            filterable
          >
            <el-option
              v-for="item in allProjects"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addProjectDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAddProjectSubmit" :loading="addProjectSubmitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 加入项目对话框 -->
    <el-dialog v-model="joinProjectDialogVisible" title="加入项目" width="600px">
      <el-form :model="joinProjectForm" label-width="100px">
        <el-form-item label="名称">
          <span>{{ joinProjectForm.name }}</span>
        </el-form-item>
        <el-form-item label="所属域">
          <span>{{ getDomainName(joinProjectForm.domain_id) }}</span>
        </el-form-item>
        <el-form-item label="项目" prop="project_id" :rules="[{ required: true, message: '请选择项目', trigger: 'change' }]">
          <el-select
            v-model="joinProjectForm.project_id"
            placeholder="请选择项目"
            style="width: 100%"
            filterable
          >
            <el-option
              v-for="item in allProjects"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="角色" prop="role_id" :rules="[{ required: true, message: '请选择角色', trigger: 'change' }]">
          <el-select
            v-model="joinProjectForm.role_id"
            placeholder="请选择角色"
            style="width: 100%"
            filterable
          >
            <el-option
              v-for="item in allRoles"
              :key="item.id"
              :label="item.display_name || item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="joinProjectDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleJoinProjectSubmit" :loading="joinProjectSubmitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 添加用户对话框 (from management view) -->
    <el-dialog v-model="addGroupUserDialogVisible" title="添加用户" width="600px">
      <el-form :model="addGroupUserForm" label-width="100px">
        <el-form-item label="名称">
          <span>{{ addGroupUserForm.name }}</span>
        </el-form-item>
        <el-form-item label="所属域">
          <span>{{ getDomainName(addGroupUserForm.domain_id) }}</span>
        </el-form-item>
        <el-form-item label="用户名" prop="user_id" :rules="[{ required: true, message: '请选择用户', trigger: 'change' }]">
          <el-select
            v-model="addGroupUserForm.user_id"
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
        <el-button @click="addGroupUserDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAddGroupUserSubmit" :loading="addGroupUserSubmitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 删除确认对话框 -->
    <el-dialog v-model="deleteConfirmDialogVisible" title="删除确认" width="600px">
      <el-form label-width="100px">
        <el-form-item label="名称">
          <span>{{ deleteConfirmInfo.name }}</span>
        </el-form-item>
        <el-form-item label="所属域">
          <span>{{ deleteConfirmInfo.domain }}</span>
        </el-form-item>
        <el-form-item label="创建时间">
          <span>{{ deleteConfirmInfo.created_at }}</span>
        </el-form-item>
      </el-form>
      <p style="color: #F56C6C;">此操作将永久删除该用户组，是否继续？</p>
      <template #footer>
        <el-button @click="deleteConfirmDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="handleDeleteConfirmed">确认删除</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Group, User, Project, Domain } from '@/types/iam'
import {
  getGroups,
  createGroup,
  updateGroup,
  deleteGroup,
  getGroupUsers,
  getGroupProjects,
  addUserToGroup,
  removeUserFromGroup,
  getUsers,
  getProjects,
  getDomains,
  addGroupToProject,
  removeGroupFromProject,
  getResourceOperationLogs
} from '@/api/iam'

const groups = ref<Group[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const addUserDialogVisible = ref(false)
const addProjectDialogVisible = ref(false)
const joinProjectDialogVisible = ref(false)  // New dialog for joining project
const addGroupUserDialogVisible = ref(false) // New dialog for adding user to group from management view
const deleteConfirmDialogVisible = ref(false) // New dialog for delete confirmation
const detailTab = ref('details') // Changed from 'users' to 'details' to match user's requirement
const isEdit = ref(false)
const submitting = ref(false)
const addUserSubmitting = ref(false)
const addProjectSubmitting = ref(false)
const joinProjectSubmitting = ref(false) // Submitting state for join project
const addGroupUserSubmitting = ref(false) // Submitting state for adding user to group
const currentGroupId = ref(0)
const currentGroup = ref<Group | null>(null)
const groupUsers = ref<User[]>([])
const groupProjects = ref<Project[]>([])
const groupLogs = ref<any[]>([]) // Operation logs
const allUsers = ref<User[]>([])
const allProjects = ref<Project[]>([])
const allDomains = ref<Domain[]>([])
const allRoles = ref<Role[]>([]) // Added for role selection
const usersLoading = ref(false)
const projectsLoading = ref(false)
const logsLoading = ref(false)
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

const addProjectForm = reactive({
  project_id: 0
})

const joinProjectForm = reactive({
  name: '',
  domain_id: 0,
  project_id: 0,
  role_id: 0
})

const addGroupUserForm = reactive({
  name: '',
  domain_id: 0,
  user_id: 0
})

const deleteConfirmInfo = reactive({
  name: '',
  domain: '',
  created_at: ''
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

const loadGroupProjects = async (groupId: number) => {
  projectsLoading.value = true
  try {
    const res = await getGroupProjects(groupId)
    groupProjects.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    projectsLoading.value = false
  }
}

const loadGroupLogs = async (groupId: number) => {
  logsLoading.value = true
  try {
    // Load operation logs for the group using the resource-specific API
    const res = await getResourceOperationLogs('group', groupId)
    // Format the logs to match the table structure
    groupLogs.value = (res.items || []).map((log: any) => ({
      id: log.id,
      created_at: log.created_at,
      resource_name: log.resource_name,
      resource_type: log.resource_type,
      operation_type: log.operation_type,
      service_type: log.service_type,
      risk_level: log.risk_level || '低',
      event_type: log.event_type || '-',
      result: log.result,
      operator: log.operator,
      project_name: log.project_name || '-'
    }))
  } catch (e: any) {
    console.error(e)
    // Use mock data in case of error
    groupLogs.value = []
  } finally {
    logsLoading.value = false
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

const loadAllProjects = async () => {
  try {
    const res = await getProjects({ limit: 100 })
    allProjects.value = res.items || []
  } catch (e: any) {
    console.error(e)
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

const loadAllRoles = async () => {
  try {
    const res = await getRoles({ limit: 100 })
    allRoles.value = res.items || []
  } catch (e: any) {
    console.error(e)
  }
}

// 当 domains 更新时，强制刷新表格视图
watch(allDomains, () => {
  // 简单地触发重新渲染
  if (groups.value.length > 0) {
    // 创建新数组引用以强制组件更新
    groups.value = [...groups.value]
  }
}, { deep: true })

const resetFilter = () => {
  filterForm.name = ''
  pagination.page = 1
  loadGroups()
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

const handleView = async (row: Group) => {
  currentGroup.value = row
  detailTab.value = 'details'
  await Promise.all([
    loadGroupUsers(row.id),
    loadGroupProjects(row.id),
    loadGroupLogs(row.id)
  ])
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

const handleAddProject = async () => {
  await loadAllProjects()
  addProjectForm.project_id = 0
  addProjectDialogVisible.value = true
}

const handleAddProjectSubmit = async () => {
  if (!addProjectForm.project_id) {
    ElMessage.error('请选择项目')
    return
  }

  addProjectSubmitting.value = true
  try {
    await addGroupToProject(currentGroupId.value, addProjectForm.project_id)
    ElMessage.success('添加项目成功')
    addProjectDialogVisible.value = false
    await loadGroupProjects(currentGroupId.value)
  } catch (e: any) {
    ElMessage.error(e.message || '添加项目失败')
  } finally {
    addProjectSubmitting.value = false
  }
}

const handleRemoveProject = async (row: Project) => {
  try {
    await ElMessageBox.confirm('确定要移除该项目吗？', '提示', { type: 'warning' })
    await removeGroupFromProject(currentGroupId.value, row.id)
    ElMessage.success('移除项目成功')
    await loadGroupProjects(currentGroupId.value)
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '移除项目失败')
    }
  }
}

const handleDelete = async (row: Group) => {
  handleDeleteWithDetails(row)
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

const handleCommand = (command: string, row: Group) => {
  switch (command) {
    case 'view': handleView(row); break
    case 'edit': handleEdit(row); break
    case 'delete': handleDelete(row); break
  }
}

const handleManageProjects = async (row: Group) => {
  currentGroup.value = row
  detailTab.value = 'projects'
  await Promise.all([
    loadGroupUsers(row.id),
    loadGroupProjects(row.id),
    loadGroupLogs(row.id)
  ])
  detailDialogVisible.value = true
}

const handleManageUsers = async (row: Group) => {
  currentGroup.value = row
  detailTab.value = 'users'
  await Promise.all([
    loadGroupUsers(row.id),
    loadGroupProjects(row.id),
    loadGroupLogs(row.id)
  ])
  detailDialogVisible.value = true
}

const handleViewLogs = async (row: Group) => {
  currentGroupId.value = row.id
  detailTab.value = 'logs'
  await loadGroupLogs(row.id)
  detailDialogVisible.value = true
}

// Handle Edit Project Role - Placeholder for future implementation
const handleEditProjectRole = async (row: Project) => {
  // This would open a dialog to edit the role assigned to the group for this project
  // For now, we'll just show a message indicating this feature is planned
  ElMessage.info('编辑项目角色功能待实现')
}

// Join Project popup handling (renaming to avoid conflict with handleAddProject)
const handleJoinProjectPopup = async (row: Group) => {
  currentGroup.value = row
  joinProjectForm.name = row.name
  joinProjectForm.domain_id = row.domain_id
  await Promise.all([loadAllProjects(), loadAllRoles()])
  joinProjectForm.project_id = 0
  joinProjectForm.role_id = 0
  joinProjectDialogVisible.value = true
}

const handleJoinProjectSubmit = async () => {
  if (!joinProjectForm.project_id) {
    ElMessage.error('请选择项目')
    return
  }

  joinProjectSubmitting.value = true
  try {
    // Add the group to the project
    await addGroupToProject(currentGroup.value!.id, joinProjectForm.project_id)
    ElMessage.success('加入项目成功')
    joinProjectDialogVisible.value = false
    await loadGroupProjects(currentGroup.value!.id)
  } catch (e: any) {
    ElMessage.error(e.message || '加入项目失败')
  } finally {
    joinProjectSubmitting.value = false
  }
}

// Add User to Group popup handling
const handleAddGroupUser = async (row: Group) => {
  currentGroup.value = row
  addGroupUserForm.name = row.name
  addGroupUserForm.domain_id = row.domain_id
  await loadAllUsers()
  addGroupUserForm.user_id = 0
  addGroupUserDialogVisible.value = true
}

const handleAddGroupUserSubmit = async () => {
  if (!addGroupUserForm.user_id) {
    ElMessage.error('请选择用户')
    return
  }

  addGroupUserSubmitting.value = true
  try {
    await addUserToGroup(currentGroup.value!.id, addGroupUserForm.user_id)
    ElMessage.success('添加用户成功')
    addGroupUserDialogVisible.value = false
    await loadGroupUsers(currentGroup.value!.id)
  } catch (e: any) {
    ElMessage.error(e.message || '添加用户失败')
  } finally {
    addGroupUserSubmitting.value = false
  }
}

// Delete Confirmation popup handling
const handleDeleteWithDetails = async (row: Group) => {
  deleteConfirmInfo.name = row.name
  deleteConfirmInfo.domain = getDomainName(row.domain_id)
  deleteConfirmInfo.created_at = formatDate(row.created_at)
  currentGroup.value = row
  deleteConfirmDialogVisible.value = true
}

const handleDeleteConfirmed = async () => {
  try {
    await deleteGroup(currentGroup.value!.id)
    ElMessage.success('删除成功')
    deleteConfirmDialogVisible.value = false
    loadGroups()
  } catch (e: any) {
    ElMessage.error(e.message || '删除失败')
  }
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const getDomainName = (domainId: number) => {
  const domain = allDomains.value.find(d => d.id === domainId);
  return domain ? domain.name : `域#${domainId}`;
}

onMounted(async () => {
  await loadAllDomains()  // 确保域名先加载完成
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
