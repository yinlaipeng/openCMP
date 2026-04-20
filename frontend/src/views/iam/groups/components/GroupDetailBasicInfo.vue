<template>
  <div class="group-basic-info" v-loading="loading">
    <el-descriptions :column="2" border v-if="group">
      <el-descriptions-item label="ID">{{ group.id }}</el-descriptions-item>
      <el-descriptions-item label="用户组名">{{ group.name }}</el-descriptions-item>
      <el-descriptions-item label="备注" :span="2">{{ group.description || '-' }}</el-descriptions-item>
      <el-descriptions-item label="所属域">{{ domainName }}</el-descriptions-item>
      <el-descriptions-item label="创建时间">{{ formatDate(group.created_at) }}</el-descriptions-item>
      <el-descriptions-item label="更新时间" :span="2">{{ formatDate(group.updated_at) }}</el-descriptions-item>
    </el-descriptions>
    <el-empty v-else description="暂无数据" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Group, Domain } from '@/types/iam'

interface Props {
  group: Group | null
  loading: boolean
  domains?: Domain[]
}

const props = defineProps<Props>()

const domainName = computed(() => {
  if (!props.group?.domain_id) return '-'
  const domain = props.domains?.find(d => d.id === props.group!.domain_id)
  return domain?.name || `域#${props.group.domain_id}`
})

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}
</script>

<style scoped>
.group-basic-info {
  padding: 16px 0;
}
</style>