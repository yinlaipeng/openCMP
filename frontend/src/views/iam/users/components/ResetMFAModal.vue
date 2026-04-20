<template>
  <el-dialog
    v-model="visible"
    title="重置MFA"
    width="500px"
    :before-close="handleClose"
  >
    <!-- 黄色警告框 -->
    <el-alert
      title="警告"
      description="重置后，用户需在登录页面重新设置，之前信息将失效！"
      type="warning"
      show-icon
      :closable="false"
      style="margin-bottom: 20px"
    />

    <!-- 确认文案 -->
    <p class="confirm-text">
      您正在重置用户 <strong>{{ user?.name }}</strong> 的MFA设置，请确认操作。
    </p>

    <!-- 迷你表格 -->
    <el-table :data="tableData" border size="small" class="mini-table">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="用户名" />
      <el-table-column prop="display_name" label="显示名" />
      <el-table-column label="MFA状态">
        <template #default>
          <span class="status-dot" data-status="enabled">
            {{ user?.mfa_enabled ? '已开启' : '未开启' }}
          </span>
        </template>
      </el-table-column>
    </el-table>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="warning" @click="handleConfirm" :loading="loading">
        确定重置
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import type { User } from '@/types/iam'
// import { resetUserMFA } from '@/api/iam' // TODO: Add API when backend implements

// Props
const props = defineProps<{
  modelValue: boolean
  user: User | null
}>()

// Emits
const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

// State
const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const loading = ref(false)

const tableData = computed(() => {
  if (!props.user) return []
  return [{
    id: props.user.id,
    name: props.user.name,
    display_name: props.user.display_name || '-',
    mfa_enabled: props.user.mfa_enabled
  }]
})

// Handlers
const handleClose = () => {
  visible.value = false
}

const handleConfirm = async () => {
  if (!props.user) return
  loading.value = true
  try {
    // TODO: Call real API when backend implements
    // await resetUserMFA(props.user.id)

    // Mock success for now
    await new Promise(resolve => setTimeout(resolve, 500))
    ElMessage.success('MFA已重置，用户需在登录页面重新设置')
    visible.value = false
    emit('success')
  } catch (e: any) {
    ElMessage.error(e.message || '重置MFA失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.confirm-text {
  margin-bottom: var(--space-3, 12px);
  font-size: var(--font-size-base, 16px);
}

.confirm-text strong {
  color: var(--color-primary, #0F172A);
}

.mini-table {
  margin-bottom: var(--space-4, 16px);
}

.status-dot {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1, 4px);
}

.status-dot::before {
  content: '';
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-dot[data-status="enabled"]::before {
  background: var(--color-success, #22C55E);
}

.status-dot[data-status="disabled"]::before {
  background: var(--color-muted, #64748B);
}
</style>