<template>
  <el-dialog
    v-model="visible"
    title="更新云账号"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="editForm"
      :rules="editRules"
      label-width="130px"
      v-loading="loading"
    >
      <!-- 基本信息（只读） -->
      <div class="section-title">基本信息</div>
      <el-form-item label="名称">
        <el-input :value="account?.name" disabled />
      </el-form-item>
      <el-form-item label="平台">
        <el-tag :type="getProviderType(account?.provider_type)">
          {{ getProviderName(account?.provider_type) }}
        </el-tag>
      </el-form-item>
      <el-form-item label="状态">
        <el-tag :type="getStatusType(account?.status)">
          {{ getStatusText(account?.status) }}
        </el-tag>
      </el-form-item>

      <!-- 编辑信息 -->
      <div class="section-title">编辑信息</div>
      <el-form-item label="备注信息" prop="description">
        <el-input
          v-model="editForm.description"
          type="textarea"
          :rows="3"
          placeholder="请输入备注信息"
          maxlength="500"
          show-word-limit
        />
      </el-form-item>
      <el-form-item label="密钥ID" prop="accessKeyId">
        <el-input
          v-model="editForm.accessKeyId"
          placeholder="请输入 Access Key ID"
        />
        <template #extra>
          <span class="field-hint">Access Key ID</span>
        </template>
      </el-form-item>
      <el-form-item label="密钥密码" prop="accessKeySecret">
        <el-input
          v-model="editForm.accessKeySecret"
          type="password"
          placeholder="请输入 Access Key Secret"
          show-password
        />
        <template #extra>
          <span class="field-hint">Access Key Secret</span>
        </template>
      </el-form-item>

      <!-- 测试连接 -->
      <el-form-item label="连接验证">
        <div class="test-connection-wrapper">
          <el-button
            @click="handleTestConnection"
            :loading="testLoading"
            :disabled="!editForm.accessKeyId || !editForm.accessKeySecret"
          >
            {{ testLoading ? '测试中...' : '测试连接' }}
          </el-button>
          <span v-if="testResult" :class="testSuccess ? 'text-success' : 'text-danger'">
            <el-icon v-if="testSuccess"><CircleCheck /></el-icon>
            <el-icon v-else><CircleClose /></el-icon>
            {{ testResult }}
          </span>
        </div>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSave" :loading="saveLoading">
        保存更改
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { CircleCheck, CircleClose } from '@element-plus/icons-vue'
import {
  updateCloudAccount,
  testConnectionWithCredentials,
  getCloudAccount
} from '@/api/cloud-account'
import type { CloudAccount } from '@/types'

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

const formRef = ref()
const account = ref<CloudAccount | null>(null)
const loading = ref(false)
const testLoading = ref(false)
const saveLoading = ref(false)
const testResult = ref('')
const testSuccess = ref(false)

const editForm = ref({
  description: '',
  accessKeyId: '',
  accessKeySecret: ''
})

const editRules = {
  accessKeyId: [
    { required: true, message: '请输入密钥ID', trigger: 'blur' },
    { min: 10, message: '密钥ID长度至少10位', trigger: 'blur' }
  ],
  accessKeySecret: [
    { required: true, message: '请输入密钥密码', trigger: 'blur' },
    { min: 10, message: '密钥密码长度至少10位', trigger: 'blur' }
  ]
}

// 加载账户数据
watch(() => props.accountId, async (id) => {
  if (id && visible.value) {
    loading.value = true
    try {
      const res = await getCloudAccount(id)
      account.value = res
      editForm.value.description = res.description || ''
      // 从凭证中提取 access_key_id（如果存在）
      if (res.credentials) {
        try {
          const creds = JSON.parse(JSON.stringify(res.credentials))
          editForm.value.accessKeyId = creds.access_key_id || ''
        } catch {
          // 忽略解析错误
        }
      }
    } catch (error) {
      ElMessage.error('加载账户信息失败')
    } finally {
      loading.value = false
    }
  }
})

// 测试连接
async function handleTestConnection() {
  if (!props.accountId) return
  if (!editForm.value.accessKeyId || !editForm.value.accessKeySecret) {
    ElMessage.warning('请先输入密钥ID和密钥密码')
    return
  }

  testLoading.value = true
  testResult.value = ''
  try {
    const res = await testConnectionWithCredentials(props.accountId, {
      access_key_id: editForm.value.accessKeyId,
      access_key_secret: editForm.value.accessKeySecret
    })
    testSuccess.value = res.connected
    testResult.value = res.message
  } catch (error: any) {
    testSuccess.value = false
    testResult.value = error.message || '连接测试失败'
  } finally {
    testLoading.value = false
  }
}

// 保存更改
async function handleSave() {
  if (!props.accountId) return
  if (!formRef.value) return

  await formRef.value.validate()
  saveLoading.value = true
  try {
    // 构建更新数据
    const updateData = {
      description: editForm.value.description,
      credentials: {
        access_key_id: editForm.value.accessKeyId,
        access_key_secret: editForm.value.accessKeySecret,
        region_id: account.value?.credentials?.region_id || 'cn-hangzhou'
      }
    }
    await updateCloudAccount(props.accountId, updateData)
    ElMessage.success('更新成功')
    emit('saved')
    handleClose()
  } catch (error: any) {
    ElMessage.error(error.message || '更新失败')
  } finally {
    saveLoading.value = false
  }
}

// 关闭对话框
function handleClose() {
  visible.value = false
  editForm.value = {
    description: '',
    accessKeyId: '',
    accessKeySecret: ''
  }
  testResult.value = ''
  testSuccess.value = false
  account.value = null
}

// 辅助函数
function getProviderType(type: string | undefined): string {
  const types: Record<string, string> = {
    alibaba: '',
    tencent: 'warning',
    aws: 'success',
    azure: 'info'
  }
  return types[type || ''] || 'info'
}

function getProviderName(type: string | undefined): string {
  const names: Record<string, string> = {
    alibaba: '阿里云',
    tencent: '腾讯云',
    aws: 'AWS',
    azure: 'Azure'
  }
  return names[type || ''] || type || '未知'
}

function getStatusType(status: string | undefined): string {
  const types: Record<string, string> = {
    active: 'success',
    inactive: 'info',
    error: 'danger'
  }
  return types[status || ''] || 'info'
}

function getStatusText(status: string | undefined): string {
  const texts: Record<string, string> = {
    active: '已连接',
    inactive: '未连接',
    error: '错误'
  }
  return texts[status || ''] || status || '未知'
}
</script>

<style scoped>
.section-title {
  font-weight: 600;
  font-size: 14px;
  color: #303133;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e4e7ed;
}

.field-hint {
  font-size: 12px;
  color: #909399;
}

.test-connection-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
}

.text-success {
  color: #67c23a;
  display: flex;
  align-items: center;
  gap: 4px;
}

.text-danger {
  color: #f56c6c;
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>