<template>
  <el-dialog
    v-model="visible"
    title="日志详情"
    width="700px"
    :before-close="handleClose"
  >
    <div class="log-detail-content" v-if="log">
      <!-- 详情两列展示 -->
      <div class="detail-layout">
        <div class="detail-column left-column">
          <el-descriptions :column="1" border size="small">
            <el-descriptions-item label="#ID">
              <span class="mono-text">{{ log.id }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="操作时间">
              <span class="mono-text">{{ formatDate(log.operation_time) }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="资源名称">{{ log.resource_name }}</el-descriptions-item>
            <el-descriptions-item label="资源类型">{{ log.resource_type }}</el-descriptions-item>
            <el-descriptions-item label="操作类型">{{ log.operation_type }}</el-descriptions-item>
            <el-descriptions-item label="服务类型">{{ log.service_type }}</el-descriptions-item>
            <el-descriptions-item label="风险级别">
              <el-tag :type="getRiskLevelType(log.risk_level)" size="small">{{ log.risk_level }}</el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <div class="detail-column right-column">
          <el-descriptions :column="1" border size="small">
            <el-descriptions-item label="事件类型">{{ log.event_type }}</el-descriptions-item>
            <el-descriptions-item label="结果">
              <span class="result-status" :class="log.result === '成功' ? 'success' : 'fail'">{{ log.result }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="发起人">{{ log.operator }}</el-descriptions-item>
            <el-descriptions-item label="所属项目">{{ log.project }}</el-descriptions-item>
            <el-descriptions-item label="请求ID">
              <span class="mono-text">{{ log.request_id || '-' }}</span>
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>

      <!-- JSON代码块查看器 -->
      <div class="json-section">
        <div class="json-header">
          <span class="json-title">请求详情 (JSON)</span>
          <el-button size="small" type="primary" link @click="handleCopy">
            <el-icon><CopyDocument /></el-icon>
            复制内容
          </el-button>
        </div>
        <pre class="json-code-block">{{ formatJson(log.details) }}</pre>
      </div>
    </div>

    <template #footer>
      <div class="footer-actions">
        <el-button @click="handlePrev" :disabled="!hasPrev">
          <el-icon><ArrowLeft /></el-icon>
          上一条
        </el-button>
        <el-button @click="handleNext" :disabled="!hasNext">
          下一条
          <el-icon><ArrowRight /></el-icon>
        </el-button>
        <el-button @click="handleClose">关闭</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { CopyDocument, ArrowLeft, ArrowRight } from '@element-plus/icons-vue'
import type { OperationLog } from '@/types/iam'

// Props
const props = defineProps<{
  modelValue: boolean
  log: OperationLog | null
  userId?: number
}>()

// Emits
const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'prev'): void
  (e: 'next'): void
}>()

// State
const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

// Navigation state (would be set by parent component)
const hasPrev = ref(false)
const hasNext = ref(false)

// Handlers
const handleClose = () => {
  visible.value = false
}

const handlePrev = () => {
  emit('prev')
}

const handleNext = () => {
  emit('next')
}

const handleCopy = async () => {
  if (!props.log) return
  try {
    await navigator.clipboard.writeText(formatJson(props.log.details))
    ElMessage.success('内容已复制')
  } catch (e) {
    ElMessage.error('复制失败')
  }
}

// Utils
const formatDate = (date?: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const formatJson = (details?: any) => {
  if (!details) return '{}'
  try {
    return JSON.stringify(details, null, 2)
  } catch (e) {
    return String(details)
  }
}

const getRiskLevelType = (level?: string) => {
  switch (level) {
    case '高': return 'danger'
    case '中': return 'warning'
    case '低': return 'success'
    default: return 'info'
  }
}
</script>

<style scoped>
.log-detail-content {
  padding: var(--space-2, 8px);
}

.detail-layout {
  display: flex;
  gap: var(--space-4, 16px);
  margin-bottom: var(--space-4, 16px);
}

.detail-column {
  flex: 1;
  min-width: 280px;
}

.mono-text {
  font-family: var(--font-mono, 'Fira Code', monospace);
}

.result-status.success {
  color: var(--color-success, #22C55E);
}

.result-status.fail {
  color: var(--color-danger, #EF4444);
}

.json-section {
  margin-top: var(--space-4, 16px);
}

.json-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-2, 8px);
}

.json-title {
  font-size: var(--font-size-sm, 14px);
  font-weight: var(--font-weight-semibold, 600);
  color: var(--color-muted, #64748B);
}

.json-code-block {
  background: var(--color-surface-dark, #0F172A);
  color: var(--color-foreground-dark, #F8FAFC);
  padding: var(--space-3, 12px);
  border-radius: var(--radius-md, 8px);
  font-family: var(--font-mono, 'Fira Code', monospace);
  font-size: var(--font-size-sm, 14px);
  overflow-x: auto;
  max-height: 300px;
}

.footer-actions {
  display: flex;
  gap: var(--space-2, 8px);
}

@media (max-width: 768px) {
  .detail-layout {
    flex-direction: column;
  }
}
</style>