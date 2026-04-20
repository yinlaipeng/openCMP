<template>
  <el-dialog v-model="visible" title="批量删除确认" width="600px" :before-close="handleClose">
    <el-alert type="warning" :closable="false" style="margin-bottom: 16px">
      <template #title>
        <span style="font-weight: 600">此操作将永久删除以下 {{ groups.length }} 个用户组，是否继续？</span>
      </template>
    </el-alert>

    <!-- 选中组列表 -->
    <el-table :data="groups" border stripe max-height="300">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" min-width="150" />
      <el-table-column prop="domain_name" label="所属域" width="120">
        <template #default="{ row }">
          {{ getDomainName(row.domain_id) }}
        </template>
      </el-table-column>
    </el-table>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="danger" @click="handleConfirm" :loading="submitting">确认删除</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Group, Domain } from '@/types/iam'

interface Props {
  modelValue: boolean
  groups: Group[]
  domains: Domain[]
  submitting: boolean
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'confirm'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const getDomainName = (domainId: number) => {
  const domain = props.domains.find(d => d.id === domainId)
  return domain?.name || `域#${domainId}`
}

const handleClose = () => {
  visible.value = false
}

const handleConfirm = () => {
  emit('confirm')
}
</script>

<style scoped>
</style>