<template>
  <el-select
    :model-value="value"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    filterable
    clearable
    @update:model-value="handleUpdate"
    @change="handleChange"
  >
    <el-option
      v-for="account in accounts"
      :key="account.id"
      :label="account.name"
      :value="account.id"
      :disabled="account.health_status === 'unhealthy'"
    >
      <div class="account-option">
        <span class="account-name">{{ account.name }}</span>
        <div class="account-tags">
          <el-tag
            :type="getProviderTagType(account.provider_type)"
            size="small"
            class="provider-tag"
          >
            {{ getProviderLabel(account.provider_type) }}
          </el-tag>
          <el-tag
            :type="getHealthTagType(account.health_status)"
            size="small"
          >
            {{ getHealthLabel(account.health_status) }}
          </el-tag>
        </div>
      </div>
    </el-option>
  </el-select>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { getCloudAccounts } from '@/api/cloud-account'
import type { CloudAccount } from '@/types'

interface Props {
  value?: number | null
  disabled?: boolean
  placeholder?: string
}

interface Emits {
  (e: 'change', accountId: number | null, account: CloudAccount | null): void
  (e: 'update:value', accountId: number | null): void
}

const props = withDefaults(defineProps<Props>(), {
  value: null,
  disabled: false,
  placeholder: '请选择云账号'
})

const emit = defineEmits<Emits>()

const accounts = ref<CloudAccount[]>([])
const loading = ref(false)

// Provider type mapping
const providerLabels: Record<string, string> = {
  alibaba: '阿里云',
  tencent: '腾讯云',
  aws: 'AWS',
  azure: 'Azure'
}

const providerTagTypes: Record<string, 'primary' | 'warning' | 'success' | 'info'> = {
  alibaba: 'primary',
  tencent: 'warning',
  aws: 'success',
  azure: 'info'
}

// Get provider label
const getProviderLabel = (providerType: string): string => {
  return providerLabels[providerType] || providerType
}

// Get provider tag type
const getProviderTagType = (providerType: string): 'primary' | 'warning' | 'success' | 'info' => {
  return providerTagTypes[providerType] || 'info'
}

// Get health status tag type
const getHealthTagType = (healthStatus?: string): 'success' | 'danger' | 'info' => {
  if (healthStatus === 'healthy') return 'success'
  if (healthStatus === 'unhealthy') return 'danger'
  return 'info'
}

// Get health status label
const getHealthLabel = (healthStatus?: string): string => {
  if (healthStatus === 'healthy') return '正常'
  if (healthStatus === 'unhealthy') return '异常'
  return '未知'
}

// Load cloud accounts
const loadAccounts = async () => {
  loading.value = true
  try {
    const res = await getCloudAccounts()
    accounts.value = res.items || []
  } catch (error) {
    console.error('Failed to load cloud accounts:', error)
    ElMessage.error('加载云账号列表失败')
    accounts.value = []
  } finally {
    loading.value = false
  }
}

// Handle model value update
const handleUpdate = (val: number | null | undefined) => {
  emit('update:value', val ?? null)
}

// Handle selection change
const handleChange = (val: number | null | undefined) => {
  const accountId = val ?? null
  const account = accounts.value.find(a => a.id === accountId) || null
  emit('change', accountId, account)
}

// Load accounts on mount
onMounted(() => {
  loadAccounts()
})

// Expose reload method for parent components
defineExpose({
  reload: loadAccounts
})
</script>

<style scoped>
.account-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.account-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.account-tags {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
  margin-left: 12px;
}

.provider-tag {
  min-width: 50px;
  text-align: center;
}
</style>