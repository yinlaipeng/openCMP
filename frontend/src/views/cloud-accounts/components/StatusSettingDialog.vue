<template>
  <el-dialog
    v-model="visible"
    title="状态设置"
    width="500px"
    :close-on-click-modal="false"
  >
    <el-form :model="form" label-width="120px" v-loading="loading">
      <el-form-item label="云账号名称">
        <el-input v-model="form.name" readonly />
      </el-form-item>

      <el-form-item label="当前状态">
        <el-tag :type="getStatusType(form.currentStatus)">
          {{ getStatusText(form.currentStatus) }}
        </el-tag>
      </el-form-item>

      <el-form-item label="设置状态">
        <el-radio-group v-model="form.newStatus">
          <el-radio value="active">已连接</el-radio>
          <el-radio value="inactive">未连接</el-radio>
          <el-radio value="error">连接错误</el-radio>
        </el-radio-group>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button type="primary" @click="handleSave" :loading="saving">确认</el-button>
      <el-button @click="visible = false">取消</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { getCloudAccount, updateCloudAccount } from '@/api/cloud-account'

interface Props {
  modelValue: boolean
  accountId: number | null
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'saved'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const loading = ref(false)
const saving = ref(false)

const form = reactive({
  name: '',
  currentStatus: 'active',
  newStatus: 'active'
})

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    active: 'success',
    inactive: 'info',
    error: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: '已连接',
    inactive: '未连接',
    error: '连接错误'
  }
  return map[status] || status
}

const loadAccount = async () => {
  if (!props.accountId) return
  loading.value = true
  try {
    const account = await getCloudAccount(props.accountId)
    form.name = account.name
    form.currentStatus = account.status
    form.newStatus = account.status
  } catch (error) {
    ElMessage.error('加载云账号失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (!props.accountId) return
  saving.value = true
  try {
    await updateCloudAccount(props.accountId, { status: form.newStatus })
    ElMessage.success('状态设置成功')
    emit('saved')
    visible.value = false
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

watch(() => props.accountId, (id) => {
  if (id && visible.value) {
    loadAccount()
  }
}, { immediate: true })
</script>