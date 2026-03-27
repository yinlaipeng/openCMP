<template>
  <div class="security-groups-page">
    <el-card class="page-card">
      <template #header>
        <span class="title">安全组管理</span>
      </template>
      <el-table :data="securityGroups" v-loading="loading">
        <el-table-column prop="id" label="ID" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="vpc_id" label="VPC" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getSecurityGroups } from '@/api/network'

const securityGroups = ref([])
const loading = ref(false)

const loadSecurityGroups = async () => {
  loading.value = true
  try {
    const res = await getSecurityGroups()
    securityGroups.value = res.items || res
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadSecurityGroups()
})
</script>

<style scoped>
.security-groups-page {
  height: 100%;
}

.page-card {
  height: 100%;
}

.title {
  font-size: 18px;
  font-weight: bold;
}
</style>
