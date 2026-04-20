<template>
  <el-dialog
    v-model="visible"
    :title="'云账户属性设置'"
    width="90%"
    top="5vh"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <!-- 顶部区域：云账号信息 + 快捷操作 -->
    <div class="drawer-header" v-if="account">
      <div class="account-icon">
        <el-avatar :size="48">
          <el-icon :size="32"><Cloudy /></el-icon>
        </el-avatar>
      </div>
      <div class="account-info">
        <h3>{{ account.name }}</h3>
        <div class="account-tags">
          <el-tag size="small" :type="getProviderType(account.provider_type)">
            {{ getProviderName(account.provider_type) }}
          </el-tag>
          <el-tag size="small" :type="getStatusType(account.status)">
            {{ getStatusText(account.status) }}
          </el-tag>
          <el-tag size="small" :type="account.enabled ? 'success' : 'info'">
            {{ account.enabled ? '启用' : '禁用' }}
          </el-tag>
        </div>
      </div>
      <div class="quick-actions">
        <el-button size="small" type="primary" @click="handleQuickSync">
          <el-icon><Refresh /></el-icon>
          同步
        </el-button>
        <el-button size="small" @click="handleQuickTest">
          <el-icon><Connection /></el-icon>
          连接测试
        </el-button>
        <el-button size="small" @click="handleQuickEdit">
          <el-icon><Edit /></el-icon>
          更新账号
        </el-button>
        <el-dropdown trigger="click" @command="handleQuickCommand">
          <el-button size="small">
            更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="statusSetting">状态设置</el-dropdown-item>
              <el-dropdown-item command="attributeSetting">属性设置</el-dropdown-item>
              <el-dropdown-item divided command="delete">
                <span style="color: var(--el-color-danger)">删除</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <el-tabs v-model="activeTab" type="border-card" @tab-click="handleTabClick" class="detail-tabs">
      <!-- 详情 Tab -->
      <el-tab-pane label="详情" name="detail">
        <DetailTab v-if="account" :account="account" :loading="detailLoading" @refresh="loadAccountDetail" />
      </el-tab-pane>

      <!-- 资源统计 Tab -->
      <el-tab-pane label="资源统计" name="resourceStats">
        <ResourceStatsTab v-if="account" :account-id="account.id" />
      </el-tab-pane>

      <!-- 订阅 Tab -->
      <el-tab-pane label="订阅" name="subscriptions">
        <SubscriptionTab v-if="account" :account-id="account.id" />
      </el-tab-pane>

      <!-- 云用户 Tab -->
      <el-tab-pane label="云用户" name="cloudUsers">
        <CloudUserTab v-if="account" :account-id="account.id" />
      </el-tab-pane>

      <!-- 云用户组 Tab -->
      <el-tab-pane label="云用户组" name="cloudUserGroups">
        <CloudUserGroupTab v-if="account" :account-id="account.id" />
      </el-tab-pane>

      <!-- 云上项目 Tab -->
      <el-tab-pane label="云上项目" name="cloudProjects">
        <CloudProjectTab v-if="account" :account-id="account.id" />
      </el-tab-pane>

      <!-- 定时任务 Tab -->
      <el-tab-pane label="定时任务" name="scheduledTasks">
        <ScheduledTaskTab v-if="account" :account-id="account.id" />
      </el-tab-pane>

      <!-- 操作日志 Tab -->
      <el-tab-pane label="操作日志" name="operationLogs">
        <OperationLogTab v-if="account" :account-id="account.id" />
      </el-tab-pane>
    </el-tabs>

    <template #footer>
      <el-button @click="handleClose">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Cloudy, Refresh, Connection, Edit, ArrowDown } from '@element-plus/icons-vue'
import { getCloudAccount, syncCloudAccount, testConnection, deleteCloudAccount } from '@/api/cloud-account'
import type { CloudAccount } from '@/types'
import DetailTab from './tabs/DetailTab.vue'
import ResourceStatsTab from './tabs/ResourceStatsTab.vue'
import SubscriptionTab from './tabs/SubscriptionTab.vue'
import CloudUserTab from './tabs/CloudUserTab.vue'
import CloudUserGroupTab from './tabs/CloudUserGroupTab.vue'
import CloudProjectTab from './tabs/CloudProjectTab.vue'
import ScheduledTaskTab from './tabs/ScheduledTaskTab.vue'
import OperationLogTab from './tabs/OperationLogTab.vue'

interface Props {
  modelValue: boolean
  accountId: number | null
  initialTab?: string
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

const activeTab = ref(props.initialTab || 'detail')
const account = ref<CloudAccount | null>(null)
const detailLoading = ref(false)

// 加载账户详情
watch(() => props.accountId, async (id) => {
  if (id && visible.value) {
    loadAccountDetail()
  }
}, { immediate: true })

async function loadAccountDetail() {
  if (!props.accountId) return
  detailLoading.value = true
  try {
    const res = await getCloudAccount(props.accountId)
    account.value = res
  } catch (error) {
    ElMessage.error('加载账户详情失败')
  } finally {
    detailLoading.value = false
  }
}

function handleTabClick(tab: any) {
  // Tab 切换时的逻辑处理
  activeTab.value = tab.paneName
}

// 平台类型映射
function getProviderType(type: string) {
  const map: Record<string, string> = {
    alibaba: 'warning',
    tencent: 'primary',
    aws: 'danger',
    azure: 'success'
  }
  return map[type] || 'info'
}

function getProviderName(type: string) {
  const map: Record<string, string> = {
    alibaba: '阿里云',
    tencent: '腾讯云',
    aws: 'AWS',
    azure: 'Azure'
  }
  return map[type] || type
}

function getStatusType(status: string) {
  const map: Record<string, string> = {
    active: 'success',
    inactive: 'info',
    error: 'danger'
  }
  return map[status] || 'info'
}

function getStatusText(status: string) {
  const map: Record<string, string> = {
    active: '已连接',
    inactive: '未连接',
    error: '连接错误'
  }
  return map[status] || status
}

// 快捷操作
async function handleQuickSync() {
  if (!account.value) return
  try {
    await syncCloudAccount(account.value.id, { mode: 'full', resource_types: ['all'] })
    ElMessage.success('同步已启动')
    loadAccountDetail()
  } catch (error) {
    ElMessage.error('同步启动失败')
  }
}

async function handleQuickTest() {
  if (!account.value) return
  try {
    const response = await testConnection(account.value.id)
    if (response.connected) {
      ElMessage.success('连接测试成功')
    } else {
      ElMessage.warning('连接测试失败: ' + (response.message || '无法连接'))
    }
  } catch (error) {
    ElMessage.error('连接测试失败')
  }
}

function handleQuickEdit() {
  // 打开编辑弹窗需要外部处理
  ElMessage.info('请在列表页使用"更新账号"功能')
}

async function handleQuickCommand(command: string) {
  if (command === 'delete') {
    if (!account.value) return
    try {
      await ElMessageBox.confirm(`确定要删除云账户 "${account.value.name}" 吗？`, '提示', { type: 'warning' })
      await deleteCloudAccount(account.value.id)
      ElMessage.success('删除成功')
      visible.value = false
      emit('refresh')
    } catch (e: any) {
      if (e !== 'cancel') {
        ElMessage.error('删除失败')
      }
    }
  }
}

function handleClose() {
  visible.value = false
  account.value = null
}
</script>

<style scoped>
.drawer-header {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--el-border-color);
  margin-bottom: 20px;
}

.account-icon {
  margin-right: 16px;
}

.account-info h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.account-tags {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.quick-actions {
  margin-left: auto;
  display: flex;
  gap: 8px;
}

.detail-tabs {
  border-radius: 4px;
}

.el-tabs--border-card {
  border-radius: 4px;
}

.el-tabs__content {
  padding: 20px;
  min-height: 400px;
}
</style>