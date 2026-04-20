<template>
  <div class="empty-state">
    <div class="empty-state-content">
      <el-icon :size="64" class="empty-icon">
        <component :is="icon" />
      </el-icon>
      <h3 class="empty-title">{{ title }}</h3>
      <p class="empty-description">{{ description }}</p>
      <el-button type="primary" @click="onCreate" v-if="showCreateButton">
        <el-icon><Plus /></el-icon>
        {{ createButtonText }}
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Component } from 'vue'
import { Plus, FolderOpened, Document, Timer } from '@element-plus/icons-vue'

interface Props {
  title?: string
  description?: string
  icon?: Component
  showCreateButton?: boolean
  createButtonText?: string
}

const props = withDefaults(defineProps<Props>(), {
  title: '暂无数据',
  description: '当前没有任何数据，点击下方按钮创建',
  icon: () => FolderOpened,
  showCreateButton: true,
  createButtonText: '添加'
})

const emit = defineEmits<{
  (e: 'create'): void
}>()

const onCreate = () => {
  emit('create')
}
</script>

<style scoped>
.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 60px 20px;
  min-height: 300px;
}

.empty-state-content {
  text-align: center;
  max-width: 400px;
}

.empty-icon {
  color: #c0c4cc;
  margin-bottom: 16px;
}

.empty-title {
  font-size: 16px;
  color: #606266;
  margin: 0 0 8px 0;
  font-weight: 500;
}

.empty-description {
  font-size: 14px;
  color: #909399;
  margin: 0 0 24px 0;
  line-height: 1.5;
}
</style>