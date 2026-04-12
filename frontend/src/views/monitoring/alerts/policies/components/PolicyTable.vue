<template>
  <el-table :data="policies" v-loading="loading">
    <el-table-column prop="name" label="名称" width="180" />
    <el-table-column prop="status" label="状态" width="100">
      <template #default="{ row }">
        <el-tag :type="getPolicyStatusType(row.status)">{{ row.status }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="enabled" label="启用状态" width="100">
      <template #default="{ row }">
        <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '已启用' : '已禁用' }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="resource_type" label="资源类型" width="100" />
    <el-table-column prop="detail" label="策略详情" width="150" />
    <el-table-column prop="level" label="告警级别" width="100">
      <template #default="{ row }">
        <el-tag :type="getAlertLevelType(row.level)" size="small">{{ row.level }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="owner" label="策略归属" width="120" />
    <el-table-column label="操作" width="150">
      <template #default="{ row }">
        <el-button size="small" @click="$emit('edit', row)">编辑</el-button>
        <el-button size="small" type="danger" @click="$emit('delete', row)">删除</el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue'

interface AlertPolicy {
  id: string
  name: string
  status: string
  enabled: boolean
  resource_type: string
  detail: string
  level: string
  owner: string
}

defineProps<{
  policies: AlertPolicy[]
  loading?: boolean
}>()

defineEmits<{
  (e: 'edit', row: AlertPolicy): void
  (e: 'delete', row: AlertPolicy): void
}>()

const getPolicyStatusType = (status: string) => {
  switch (status) {
    case '正常':
      return 'success'
    case '异常':
      return 'danger'
    default:
      return 'info'
  }
}

const getAlertLevelType = (level: string) => {
  switch (level) {
    case '严重':
      return 'danger'
    case '警告':
      return 'warning'
    case '信息':
      return 'info'
    default:
      return ''
  }
}
</script>