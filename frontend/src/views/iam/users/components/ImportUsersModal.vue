<template>
  <el-dialog
    v-model="visible"
    title="导入用户"
    width="600px"
    :before-close="handleClose"
  >
    <!-- 说明 -->
    <el-alert
      title="导入说明"
      description="支持 CSV 格式文件导入，每行一个用户。文件需包含：用户名、显示名、邮箱、密码等字段。"
      type="info"
      show-icon
      :closable="false"
      style="margin-bottom: 20px"
    />

    <!-- 文件上传 -->
    <el-upload
      ref="uploadRef"
      :auto-upload="false"
      :limit="1"
      accept=".csv"
      :on-change="handleFileChange"
      :on-exceed="handleExceed"
      drag
      class="upload-area"
    >
      <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
      <div class="el-upload__text">
        将 CSV 文件拖到此处，或 <em>点击上传</em>
      </div>
      <template #tip>
        <div class="el-upload__tip">
          只能上传 CSV 文件，且文件大小不超过 5MB
        </div>
      </template>
    </el-upload>

    <!-- 预览表格 -->
    <div v-if="previewData.length > 0" class="preview-section">
      <h4 class="preview-title">数据预览 (前 10 条)</h4>
      <el-table :data="previewData" border size="small" max-height="300">
        <el-table-column prop="name" label="用户名" width="120" />
        <el-table-column prop="display_name" label="显示名" width="120" />
        <el-table-column prop="email" label="邮箱" width="150" />
        <el-table-column prop="domain_id" label="所属域ID" width="100" />
        <el-table-column prop="password" label="密码" width="100">
          <template #default>
            <span class="password-mask">******</span>
          </template>
        </el-table-column>
      </el-table>
      <p class="preview-summary">
        共解析出 <strong>{{ totalCount }}</strong> 条用户数据
      </p>
    </div>

    <!-- 选项 -->
    <el-form :model="form" label-width="120px" class="options-form">
      <el-form-item label="所属域">
        <el-select v-model="form.domain_id" placeholder="请选择域">
          <el-option v-for="d in domains" :key="d.id" :label="d.name" :value="d.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="冲突处理">
        <el-radio-group v-model="form.conflict_mode">
          <el-radio value="skip">跳过已存在用户</el-radio>
          <el-radio value="update">更新已存在用户</el-radio>
        </el-radio-group>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button
        type="primary"
        @click="handleImport"
        :loading="loading"
        :disabled="previewData.length === 0"
      >
        开始导入
      </el-button>
    </template>
  </el-dialog>

  <!-- 导入结果弹窗 -->
  <el-dialog v-model="resultVisible" title="导入结果" width="500px">
    <el-result
      :icon="importResult.success ? 'success' : 'warning'"
      :title="importResult.title"
      :sub-title="importResult.message"
    >
      <template #extra>
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="成功导入">{{ importResult.successCount }}</el-descriptions-item>
          <el-descriptions-item label="跳过">{{ importResult.skipCount }}</el-descriptions-item>
          <el-descriptions-item label="失败">{{ importResult.failCount }}</el-descriptions-item>
        </el-descriptions>
      </template>
    </el-result>
    <template #footer>
      <el-button type="primary" @click="resultVisible = false">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import type { UploadInstance, UploadFile } from 'element-plus'
import type { Domain } from '@/types/iam'
import { getDomains, createUser } from '@/api/iam'

// Props
const props = defineProps<{
  modelValue: boolean
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
const uploadRef = ref<UploadInstance>()
const domains = ref<Domain[]>([])
const previewData = ref<any[]>([])
const totalCount = ref(0)
const resultVisible = ref(false)

const form = ref({
  domain_id: 1,
  conflict_mode: 'skip'
})

const importResult = ref({
  success: true,
  title: '',
  message: '',
  successCount: 0,
  skipCount: 0,
  failCount: 0
})

// Load domains when dialog opens
watch(() => props.modelValue, (val) => {
  if (val) {
    loadDomains()
    previewData.value = []
    totalCount.value = 0
  }
})

const loadDomains = async () => {
  try {
    const res = await getDomains({ limit: 100 })
    domains.value = res.items || []
    if (domains.value.length > 0) {
      form.value.domain_id = domains.value[0].id
    }
  } catch (e) {
    console.error('加载域失败', e)
  }
}

// Handlers
const handleClose = () => {
  visible.value = false
}

const handleFileChange = (file: UploadFile) => {
  if (!file.raw) return

  // Parse CSV file
  const reader = new FileReader()
  reader.onload = (e) => {
    const text = e.target?.result as string
    parseCSV(text)
  }
  reader.readAsText(file.raw)
}

const parseCSV = (text: string) => {
  const lines = text.split('\n').filter(line => line.trim())
  if (lines.length < 2) {
    ElMessage.error('CSV 文件格式不正确，需要包含表头和数据')
    return
  }

  // Parse header
  const headers = lines[0].split(',').map(h => h.trim().toLowerCase())

  // Parse data rows
  const data = []
  for (let i = 1; i < Math.min(lines.length, 11); i++) {
    const values = lines[i].split(',').map(v => v.trim())
    const row: any = {}
    headers.forEach((h, idx) => {
      row[h] = values[idx] || ''
    })
    data.push(row)
  }

  previewData.value = data
  totalCount.value = lines.length - 1
}

const handleExceed = () => {
  ElMessage.warning('只能上传一个文件')
}

const handleImport = async () => {
  if (previewData.value.length === 0) {
    ElMessage.error('请先上传 CSV 文件')
    return
  }

  loading.value = true
  const result = {
    successCount: 0,
    skipCount: 0,
    failCount: 0
  }

  try {
    // TODO: Replace with batch import API when backend implements
    // For now, import one by one
    for (const row of previewData.value) {
      try {
        await createUser({
          name: row.name || row.username,
          display_name: row.display_name || row.name,
          email: row.email,
          password: row.password,
          domain_id: form.value.domain_id,
          console_login: true,
          mfa_enabled: false
        })
        result.successCount++
      } catch (e: any) {
        if (e.message?.includes('已存在')) {
          result.skipCount++
        } else {
          result.failCount++
        }
      }
    }

    importResult.value = {
      success: result.failCount === 0,
      title: result.failCount === 0 ? '导入完成' : '导入完成（部分失败）',
      message: `共处理 ${totalCount.value} 条数据`,
      ...result
    }

    resultVisible.value = true
    emit('success')
  } catch (e: any) {
    ElMessage.error(e.message || '导入失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.upload-area {
  margin-bottom: var(--space-4, 16px);
}

.preview-section {
  margin-bottom: var(--space-4, 16px);
}

.preview-title {
  font-size: var(--font-size-sm, 14px);
  font-weight: var(--font-weight-semibold, 600);
  color: var(--color-muted, #64748B);
  margin-bottom: var(--space-2, 8px);
}

.preview-summary {
  margin-top: var(--space-2, 8px);
  font-size: var(--font-size-sm, 14px);
}

.password-mask {
  color: var(--color-muted, #64748B);
}

.options-form {
  margin-top: var(--space-4, 16px);
}
</style>