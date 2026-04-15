<template>
  <el-dialog
    v-model="visible"
    :title="'云账户属性设置 - ' + (account?.name || '')"
    width="90%"
    top="5vh"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-tabs v-model="activeTab" type="border-card" @tab-click="handleTabClick">
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
import { ElMessage } from 'element-plus'
import { getCloudAccount } from '@/api/cloud-account'
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

function handleClose() {
  visible.value = false
  account.value = null
}
</script>

<style scoped>
.el-tabs--border-card {
  border-radius: 4px;
}

.el-tabs__content {
  padding: 20px;
  min-height: 400px;
}
</style>