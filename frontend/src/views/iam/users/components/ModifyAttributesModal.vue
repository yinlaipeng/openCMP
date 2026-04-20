<template>
  <el-dialog
    v-model="visible"
    :title="`修改用户属性（${user?.name || ''}）`"
    width="500px"
    :before-close="handleClose"
  >
    <!-- 确认文案 -->
    <p class="confirm-text">
      您正在修改用户 <strong>{{ user?.name }}</strong> 的属性，请确认以下信息。
    </p>

    <!-- 迷你表格 -->
    <el-table :data="tableData" border size="small" class="mini-table">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="用户名" />
      <el-table-column prop="display_name" label="显示名" />
      <el-table-column prop="domain" label="所属域" />
    </el-table>

    <!-- 属性修改表单 -->
    <el-form :model="form" label-width="120px" class="modify-form">
      <el-form-item label="显示名">
        <el-input v-model="form.display_name" placeholder="请输入显示名" />
      </el-form-item>
      <el-form-item label="登录控制台">
        <el-switch v-model="form.console_login" :active-value="true" :inactive-value="false" />
      </el-form-item>
      <el-form-item label="启用MFA">
        <el-switch v-model="form.mfa_enabled" :active-value="true" :inactive-value="false" />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleConfirm" :loading="loading">
        确定
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { User } from '@/types/iam'
import { updateUser } from '@/api/iam'

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

const form = ref({
  display_name: '',
  console_login: true,
  mfa_enabled: false
})

const tableData = computed(() => {
  if (!props.user) return []
  return [{
    id: props.user.id,
    name: props.user.name,
    display_name: props.user.display_name || '-',
    domain: props.user.domain_name || `域#${props.user.domain_id}`
  }]
})

// Initialize form when user changes
watch(() => props.user, (user) => {
  if (user) {
    form.value = {
      display_name: user.display_name || '',
      console_login: user.console_login ?? true,
      mfa_enabled: user.mfa_enabled ?? false
    }
  }
}, { immediate: true })

// Handlers
const handleClose = () => {
  visible.value = false
}

const handleConfirm = async () => {
  if (!props.user) return
  loading.value = true
  try {
    await updateUser(props.user.id, {
      display_name: form.value.display_name,
      console_login: form.value.console_login,
      mfa_enabled: form.value.mfa_enabled
    })
    ElMessage.success('用户属性修改成功')
    visible.value = false
    emit('success')
  } catch (e: any) {
    ElMessage.error(e.message || '修改失败')
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

.modify-form {
  margin-top: var(--space-4, 16px);
}
</style>