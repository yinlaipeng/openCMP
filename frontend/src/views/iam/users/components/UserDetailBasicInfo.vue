<template>
  <div class="basic-info-container" v-loading="loading">
    <div class="info-layout">
      <!-- 左列：Basic Info -->
      <div class="info-column left-column">
        <h4 class="section-title">Basic Info</h4>
        <el-descriptions :column="1" border class="info-descriptions">
          <el-descriptions-item label="ID">
            <span class="mono-text">{{ user?.id }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <span class="status-dot" :data-status="user?.enabled ? 'enabled' : 'disabled'">
              {{ user?.enabled ? '启用' : '禁用' }}
            </span>
          </el-descriptions-item>
          <el-descriptions-item label="名称">
            <span class="mono-text">{{ user?.name }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="所属域">
            <el-link type="primary" @click="handleDomainClick">{{ domain }}</el-link>
          </el-descriptions-item>
          <el-descriptions-item label="显示名">{{ user?.display_name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <span class="status-dot" :data-status="user?.enabled ? 'enabled' : 'disabled'">
              {{ user?.enabled ? '启用' : '禁用' }}
            </span>
          </el-descriptions-item>
          <el-descriptions-item label="控制台登录">
            <span class="status-dot" :data-status="user?.console_login ? 'enabled' : 'disabled'">
              {{ user?.console_login ? '允许' : '禁止' }}
            </span>
          </el-descriptions-item>
          <el-descriptions-item label="MFA">
            <span class="status-dot" :data-status="user?.mfa_enabled ? 'enabled' : 'disabled'">
              {{ user?.mfa_enabled ? '已开启' : '未开启' }}
            </span>
          </el-descriptions-item>
          <el-descriptions-item label="标签">
            <template v-if="user?.tags?.length">
              <el-tag v-for="tag in user.tags" :key="tag" size="small" class="tag-item">{{ tag }}</el-tag>
            </template>
            <template v-else>
              <span class="empty-text">-</span>
            </template>
          </el-descriptions-item>
          <el-descriptions-item label="组数">
            <span class="mono-text">{{ groupCount }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="项目数">
            <span class="mono-text">{{ projectCount }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="认证源">{{ authSource || '-' }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">
            <span class="mono-text">{{ formatDate(user?.created_at) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="更新时间">
            <span class="mono-text">{{ formatDate(user?.updated_at) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="备注">{{ user?.remark || '-' }}</el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- 右列：Other Info -->
      <div class="info-column right-column">
        <h4 class="section-title">Other Info</h4>
        <el-descriptions :column="1" border class="info-descriptions">
          <el-descriptions-item label="最后登录IP">
            <span class="mono-text">{{ user?.last_login_ip || '-' }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="最后访问方式">{{ user?.last_access_method || '-' }}</el-descriptions-item>
          <el-descriptions-item label="最后登录时间">
            <span class="mono-text">{{ formatDate(user?.last_login_time) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="密码过期时间">
            <span class="mono-text">{{ formatDate(user?.password_expire_time) }}</span>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import type { User } from '@/types/iam'

const props = defineProps<{
  user: User | null
  domain: string
  loading: boolean
}>()

const emit = defineEmits<{
  (e: 'domain-click', domainId: number): void
}>()

// Computed
const groupCount = computed(() => props.user?.group_count || 0)
const projectCount = computed(() => props.user?.project_count || 0)
const authSource = computed(() => props.user?.auth_source || 'Local')

// Methods
const formatDate = (date?: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const handleDomainClick = () => {
  if (props.user?.domain_id) {
    emit('domain-click', props.user.domain_id)
  }
}
</script>

<style scoped>
.basic-info-container {
  padding: var(--space-4, 16px);
}

.info-layout {
  display: flex;
  gap: var(--space-8, 32px);
}

.info-column {
  flex: 1;
  min-width: 300px;
}

.section-title {
  font-size: var(--font-size-sm, 14px);
  font-weight: var(--font-weight-semibold, 600);
  color: var(--color-muted, #64748B);
  margin-bottom: var(--space-3, 12px);
  padding-bottom: var(--space-2, 8px);
  border-bottom: 1px solid var(--color-border, #E2E8F0);
}

.info-descriptions {
  width: 100%;
}

.mono-text {
  font-family: var(--font-mono, 'Fira Code', monospace);
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

.empty-text {
  color: var(--color-muted, #64748B);
}

.tag-item {
  margin-right: var(--space-1, 4px);
}

@media (max-width: 768px) {
  .info-layout {
    flex-direction: column;
  }
}
</style>