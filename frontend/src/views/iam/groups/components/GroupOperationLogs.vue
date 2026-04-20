<template>
  <div class="group-operation-logs">
    <!-- 工具栏 -->
    <div class="tab-toolbar">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索日志..."
        clearable
        style="width: 200px"
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <!-- 日志表格 -->
    <el-table :data="logs" v-loading="loading" border stripe max-height="400">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="created_at" label="操作时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column prop="operation_type" label="操作类型" width="120" />
      <el-table-column prop="resource_type" label="资源类型" width="100" />
      <el-table-column prop="resource_name" label="资源名称" min-width="150" show-overflow-tooltip />
      <el-table-column prop="result" label="结果" width="80">
        <template #default="{ row }">
          <span class="status-dot" :data-status="row.result === 'success' ? 'enabled' : 'disabled'">
            {{ row.result === 'success' ? '成功' : '失败' }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="operator" label="操作人" width="120" />
    </el-table>

    <!-- 空状态 -->
    <el-empty v-if="!loading && logs.length === 0" description="暂无操作日志" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'

interface Props {
  groupId: number
  loading: boolean
  logs: any[]
}

interface Emits {
  (e: 'refresh'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const searchKeyword = ref('')

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const handleSearch = () => {
  emit('refresh')
}
</script>

<style scoped>
.group-operation-logs {
  padding: 16px 0;
}

.tab-toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 16px;
}

/* 状态圆点样式 */
.status-dot {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.status-dot::before {
  content: '';
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}

.status-dot[data-status="enabled"]::before {
  background: #22C55E;
}

.status-dot[data-status="disabled"]::before {
  background: #EF4444;
}
</style>