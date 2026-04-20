<template>
  <el-drawer
    v-model="visible"
    :title="group?.name || '用户组详情'"
    direction="rtl"
    :size="drawerWidth"
    :before-close="handleClose"
  >
    <!-- Header 区域 -->
    <template #header>
      <div class="drawer-header">
        <div class="header-left">
          <el-icon class="group-icon" :size="24"><UserFilled /></el-icon>
          <span class="group-name">{{ group?.name }}</span>
        </div>
        <div class="header-right">
          <el-button @click="handleRefresh" :loading="refreshing">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
          <el-button @click="handleEdit">
            <el-icon><Edit /></el-icon>
            编辑
          </el-button>
          <el-dropdown trigger="click" @command="handleMoreCommand">
            <el-button>
              更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="delete">
                  <span style="color: #F56C6C">删除</span>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </template>

    <!-- Tabs 区域 -->
    <el-tabs v-model="activeTab" class="detail-tabs">
      <el-tab-pane label="详情" name="details">
        <GroupDetailBasicInfo :group="group" :loading="loading" />
      </el-tab-pane>
      <el-tab-pane label="已加入项目" name="projects">
        <GroupProjectsTab
          :group-id="groupId"
          :loading="projectsLoading"
          :projects="groupProjects"
          @refresh="loadGroupProjects"
        />
      </el-tab-pane>
      <el-tab-pane label="组内用户" name="users">
        <GroupUsersTab
          :group-id="groupId"
          :loading="usersLoading"
          :users="groupUsers"
          @refresh="loadGroupUsers"
        />
      </el-tab-pane>
      <el-tab-pane label="操作日志" name="logs">
        <GroupOperationLogs
          :group-id="groupId"
          :loading="logsLoading"
          :logs="groupLogs"
          @refresh="loadGroupLogs"
        />
      </el-tab-pane>
    </el-tabs>

    <!-- 编辑弹窗 -->
    <el-dialog v-model="editDialogVisible" title="编辑用户组" width="500px" append-to-body>
      <el-form :model="editForm" :rules="editRules" ref="editFormRef" label-width="100px">
        <el-form-item label="用户组名" prop="name">
          <el-input v-model="editForm.name" placeholder="请输入用户组名" disabled />
        </el-form-item>
        <el-form-item label="备注" prop="description">
          <el-input v-model="editForm.description" type="textarea" :rows="3" placeholder="请输入用户组描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleEditSubmit" :loading="editSubmitting">确定</el-button>
      </template>
    </el-dialog>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { UserFilled, Refresh, Edit, ArrowDown } from '@element-plus/icons-vue'
import { Group } from '@/types/iam'
import {
  getGroup,
  getGroupUsers,
  getGroupProjects,
  getResourceOperationLogs,
  updateGroup,
  deleteGroup
} from '@/api/iam'
import GroupDetailBasicInfo from './GroupDetailBasicInfo.vue'
import GroupProjectsTab from './GroupProjectsTab.vue'
import GroupUsersTab from './GroupUsersTab.vue'
import GroupOperationLogs from './GroupOperationLogs.vue'

interface Props {
  groupId: number
  modelValue: boolean
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'refresh'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const drawerWidth = computed(() => {
  if (window.innerWidth < 768) return '100%'
  if (window.innerWidth < 1024) return '70%'
  return '60%'
})

const activeTab = ref('details')
const group = ref<Group | null>(null)
const loading = ref(false)
const refreshing = ref(false)
const projectsLoading = ref(false)
const usersLoading = ref(false)
const logsLoading = ref(false)

const groupProjects = ref<any[]>([])
const groupUsers = ref<any[]>([])
const groupLogs = ref<any[]>([])

const editDialogVisible = ref(false)
const editSubmitting = ref(false)
const editFormRef = ref<FormInstance>()
const editForm = reactive({
  name: '',
  description: ''
})
const editRules: FormRules = {
  name: [{ required: true, message: '请输入用户组名', trigger: 'blur' }]
}

const loadGroup = async () => {
  if (!props.groupId) return
  loading.value = true
  try {
    const res = await getGroup(props.groupId)
    group.value = res
    editForm.name = res.name
    editForm.description = res.description || ''
  } catch (e: any) {
    ElMessage.error(e.message || '加载用户组详情失败')
  } finally {
    loading.value = false
  }
}

const loadGroupProjects = async () => {
  if (!props.groupId) return
  projectsLoading.value = true
  try {
    const res = await getGroupProjects(props.groupId)
    groupProjects.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    projectsLoading.value = false
  }
}

const loadGroupUsers = async () => {
  if (!props.groupId) return
  usersLoading.value = true
  try {
    const res = await getGroupUsers(props.groupId)
    groupUsers.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    usersLoading.value = false
  }
}

const loadGroupLogs = async () => {
  if (!props.groupId) return
  logsLoading.value = true
  try {
    const res = await getResourceOperationLogs('group', props.groupId)
    groupLogs.value = res.items || []
  } catch (e: any) {
    console.error(e)
  } finally {
    logsLoading.value = false
  }
}

const handleRefresh = async () => {
  refreshing.value = true
  try {
    await Promise.all([
      loadGroup(),
      loadGroupProjects(),
      loadGroupUsers(),
      loadGroupLogs()
    ])
    ElMessage.success('刷新成功')
  } finally {
    refreshing.value = false
  }
}

const handleEdit = () => {
  editForm.name = group.value?.name || ''
  editForm.description = group.value?.description || ''
  editDialogVisible.value = true
}

const handleEditSubmit = async () => {
  if (!editFormRef.value) return
  await editFormRef.value.validate(async (valid) => {
    if (!valid) return
    editSubmitting.value = true
    try {
      await updateGroup(props.groupId, editForm)
      ElMessage.success('更新成功')
      editDialogVisible.value = false
      await loadGroup()
      emit('refresh')
    } catch (e: any) {
      ElMessage.error(e.message || '更新失败')
    } finally {
      editSubmitting.value = false
    }
  })
}

const handleMoreCommand = async (command: string) => {
  if (command === 'delete') {
    try {
      await ElMessageBox.confirm('确定要删除该用户组吗？', '删除确认', { type: 'warning' })
      await deleteGroup(props.groupId)
      ElMessage.success('删除成功')
      visible.value = false
      emit('refresh')
    } catch (e: any) {
      if (e !== 'cancel') {
        ElMessage.error(e.message || '删除失败')
      }
    }
  }
}

const handleClose = () => {
  visible.value = false
}

watch(visible, (val) => {
  if (val && props.groupId) {
    loadGroup()
    loadGroupProjects()
    loadGroupUsers()
    loadGroupLogs()
  }
})
</script>

<style scoped>
.drawer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.group-icon {
  color: var(--el-color-primary);
}

.group-name {
  font-size: 18px;
  font-weight: 600;
}

.header-right {
  display: flex;
  gap: 8px;
}

.detail-tabs {
  padding: 0 16px;
}
</style>