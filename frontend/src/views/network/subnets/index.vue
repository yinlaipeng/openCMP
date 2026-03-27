<template>
  <div class="subnets-page">
    <el-card class="page-card">
      <template #header>
        <span class="title">子网管理</span>
      </template>
      <el-table :data="subnets" v-loading="loading">
        <el-table-column prop="id" label="ID" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="cidr" label="CIDR" />
        <el-table-column prop="zone_id" label="可用区" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getSubnets } from '@/api/network'

const subnets = ref([])
const loading = ref(false)

const loadSubnets = async () => {
  loading.value = true
  try {
    const res = await getSubnets()
    subnets.value = res.items || res
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadSubnets()
})
</script>

<style scoped>
.subnets-page {
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
