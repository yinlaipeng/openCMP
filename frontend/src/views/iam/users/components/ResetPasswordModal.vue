<template>
  <el-dialog
    v-model="visible"
    title="重置密码"
    width="500px"
    :before-close="handleClose"
  >
    <!-- 确认文案 -->
    <el-alert
      title="重置用户密码"
      description="重置后用户需要使用新密码重新登录"
      type="warning"
      show-icon
      :closable="false"
      style="margin-bottom: 20px"
    />

    <!-- 迷你表格 -->
    <el-table :data="tableData" border size="small" class="mini-table">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="用户名" />
      <el-table-column prop="display_name" label="显示名" />
    </el-table>

    <!-- 密码输入 -->
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px" class="password-form">
      <el-form-item label="新密码" prop="password">
        <el-input
          v-model="form.password"
          type="password"
          placeholder="请输入新密码（至少 8 位）"
          show-password
        >
          <template #suffix>
            <el-icon class="password-icon"><View /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleConfirm" :loading="loading">
        确定重置
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { View } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import type { User } from '@/types/iam'
import { resetUserPassword } from '@/api/iam'

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
const formRef = ref<FormInstance>()

const form = ref({
  password: ''
})

const rules: FormRules = {
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 8, message: '密码长度至少 8 位', trigger: 'blur' }
  ]
}

const tableData = computed(() => {
  if (!props.user) return []
  return [{
    id: props.user.id,
    name: props.user.name,
    display_name: props.user.display_name || '-'
  }]
})

// Reset form when dialog opens
watch(() => props.modelValue, (val) => {
  if (val) {
    form.value.password = ''
  }
})

// Handlers
const handleClose = () => {
  visible.value = false
}

const handleConfirm = async () => {
  if (!props.user) return

  await formRef.value?.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      await resetUserPassword(props.user!.id, form.value.password)
      ElMessage.success('密码重置成功')
      visible.value = false
      emit('success')
    } catch (e: any) {
      ElMessage.error(e.message || '密码重置失败')
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.mini-table {
  margin-bottom: var(--space-4, 16px);
}

.password-form {
  margin-top: var(--space-4, 16px);
}

.password-icon {
  cursor: pointer;
}
</style>