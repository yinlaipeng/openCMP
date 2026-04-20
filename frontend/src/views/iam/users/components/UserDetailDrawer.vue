<template>
  <el-drawer
    v-model="visible"
    :title="user?.name || '用户详情'"
    direction="rtl"
    :size="drawerWidth"
    :before-close="handleClose"
    class="user-detail-drawer"
  >
    <!-- Header 区域 -->
    <template #header>
      <div class="drawer-header">
        <div class="header-left">
          <el-avatar :size="40" class="user-avatar">
            <el-icon><User /></el-icon>
          </el-avatar>
          <span class="user-name">{{ user?.display_name || user?.name }}</span>
        </div>
        <div class="header-right">
          <el-button size="small" @click="handleRefresh" :loading="loading">
            <el-icon><Refresh /></el-icon>
          </el-button>
          <el-button size="small" type="primary" @click="handleModifyAttributes">
            修改属性
          </el-button>
          <el-dropdown trigger="click" @command="handleCommand">
            <el-button size="small">
              更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="enable" v-if="!user?.enabled">启用</el-dropdown-item>
                <el-dropdown-item command="disable" v-if="user?.enabled">禁用</el-dropdown-item>
                <el-dropdown-item command="resetPassword">重置密码</el-dropdown-item>
                <el-dropdown-item command="resetMFA">重置MFA</el-dropdown-item>
                <el-dropdown-item command="delete" divided>
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
      <el-tab-pane label="详情" name="detail">
        <UserDetailBasicInfo :user="user" :domain="domainName" :loading="detailLoading" />
      </el-tab-pane>
      <el-tab-pane label="已加入项目" name="projects">
        <UserJoinedProjects :user-id="user?.id" @refresh="loadUserProjects" />
      </el-tab-pane>
      <el-tab-pane label="已加入组" name="groups">
        <UserJoinedGroups :user-id="user?.id" @refresh="loadUserGroups" />
      </el-tab-pane>
      <el-tab-pane label="操作日志" name="logs">
        <UserOperationLogs :user-id="user?.id" @view-log="handleViewLog" />
      </el-tab-pane>
    </el-tabs>
  </el-drawer>

  <!-- 子弹窗 -->
  <ModifyAttributesModal
    v-model="modifyModalVisible"
    :user="user"
    @success="handleModifySuccess"
  />
  <ResetPasswordModal
    v-model="resetPwdModalVisible"
    :user="user"
    @success="handleResetPwdSuccess"
  />
  <ResetMFAModal
    v-model="resetMFAModalVisible"
    :user="user"
    @success="handleResetMFASuccess"
  />
  <LogDetailModal
    v-model="logDetailModalVisible"
    :log="currentLog"
    :user-id="user?.id"
    @prev="handlePrevLog"
    @next="handleNextLog"
  />
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { User, Refresh, ArrowDown } from '@element-plus/icons-vue'
import type { User as UserType, Domain, OperationLog } from '@/types/iam'
import { getUser, getUserProjects, getUserGroups, enableUser, disableUser, deleteUser, getDomains } from '@/api/iam'

import UserDetailBasicInfo from './UserDetailBasicInfo.vue'
import UserJoinedProjects from './UserJoinedProjects.vue'
import UserJoinedGroups from './UserJoinedGroups.vue'
import UserOperationLogs from './UserOperationLogs.vue'
import ModifyAttributesModal from './ModifyAttributesModal.vue'
import ResetPasswordModal from './ResetPasswordModal.vue'
import ResetMFAModal from './ResetMFAModal.vue'
import LogDetailModal from './LogDetailModal.vue'

// Props
const props = defineProps<{
  modelValue: boolean
  userId?: number
}>()

// Emits
const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'refresh'): void
  (e: 'deleted'): void
}>()

// State
const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const activeTab = ref('detail')
const loading = ref(false)
const detailLoading = ref(false)
const user = ref<UserType | null>(null)
const domainName = ref('')
const allDomains = ref<Domain[]>([])

// Modal visibility
const modifyModalVisible = ref(false)
const resetPwdModalVisible = ref(false)
const resetMFAModalVisible = ref(false)
const logDetailModalVisible = ref(false)
const currentLog = ref<OperationLog | null>(null)

// Drawer width - responsive
const drawerWidth = computed(() => {
  if (window.innerWidth < 768) return '100%'
  if (window.innerWidth < 1024) return '70%'
  return '60%'
})

// Load user data when drawer opens
watch(() => props.modelValue, (val) => {
  if (val && props.userId) {
    loadUserData()
    loadDomains()
  }
})

const loadUserData = async () => {
  if (!props.userId) return
  loading.value = true
  detailLoading.value = true
  try {
    const res = await getUser(props.userId)
    user.value = res
    // Find domain name
    const domain = allDomains.value.find(d => d.id === user.value?.domain_id)
    domainName.value = domain?.name || `域#${user.value?.domain_id}`
  } catch (e: any) {
    ElMessage.error(e.message || '加载用户详情失败')
  } finally {
    loading.value = false
    detailLoading.value = false
  }
}

const loadDomains = async () => {
  try {
    const res = await getDomains({ limit: 100 })
    allDomains.value = res.items || []
  } catch (e: any) {
    console.error('加载域列表失败', e)
  }
}

const loadUserProjects = async () => {
  // Handled by UserJoinedProjects component
}

const loadUserGroups = async () => {
  // Handled by UserJoinedGroups component
}

const handleRefresh = () => {
  loadUserData()
}

const handleClose = () => {
  visible.value = false
}

const handleModifyAttributes = () => {
  modifyModalVisible.value = true
}

const handleCommand = async (command: string) => {
  switch (command) {
    case 'enable':
      await handleToggleEnable(true)
      break
    case 'disable':
      await handleToggleEnable(false)
      break
    case 'resetPassword':
      resetPwdModalVisible.value = true
      break
    case 'resetMFA':
      resetMFAModalVisible.value = true
      break
    case 'delete':
      await handleDelete()
      break
  }
}

const handleToggleEnable = async (enable: boolean) => {
  try {
    const action = enable ? '启用' : '禁用'
    await ElMessageBox.confirm(`确定要${action}该用户吗？`, '提示', { type: 'warning' })
    if (enable) {
      await enableUser(user.value!.id)
    } else {
      await disableUser(user.value!.id)
    }
    ElMessage.success(`${action}成功`)
    loadUserData()
    emit('refresh')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '操作失败')
    }
  }
}

const handleDelete = async () => {
  try {
    await ElMessageBox.confirm('确定要删除该用户吗？此操作不可恢复', '提示', { type: 'warning' })
    await deleteUser(user.value!.id)
    ElMessage.success('删除成功')
    visible.value = false
    emit('deleted')
    emit('refresh')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const handleViewLog = (log: OperationLog) => {
  currentLog.value = log
  logDetailModalVisible.value = true
}

const handlePrevLog = () => {
  // Navigation logic handled by LogDetailModal
}

const handleNextLog = () => {
  // Navigation logic handled by LogDetailModal
}

const handleModifySuccess = () => {
  loadUserData()
  emit('refresh')
}

const handleResetPwdSuccess = () => {
  ElMessage.success('密码重置成功')
}

const handleResetMFASuccess = () => {
  loadUserData()
}
</script>

<style scoped>
.user-detail-drawer {
  --el-drawer-padding-primary: 0;
}

.drawer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding-right: 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--space-3, 12px);
}

.user-avatar {
  background: var(--color-primary, #0F172A);
  color: white;
}

.user-name {
  font-size: var(--font-size-lg, 18px);
  font-weight: var(--font-weight-semibold, 600);
}

.header-right {
  display: flex;
  align-items: center;
  gap: var(--space-2, 8px);
}

.detail-tabs {
  padding: var(--space-4, 16px);
}

.detail-tabs :deep(.el-tabs__header) {
  margin-bottom: var(--space-4, 16px);
}
</style>