<template>
  <div class="eips-page">
    <el-card class="page-card">
      <template #header>
        <span class="title">弹性 IP</span>
      </template>
      <el-table :data="eips" v-loading="loading">
        <el-table-column prop="id" label="ID" />
        <el-table-column prop="address" label="IP 地址" />
        <el-table-column prop="bandwidth" label="带宽 (Mbps)" />
        <el-table-column prop="status" label="状态" />
        <el-table-column prop="region_id" label="区域" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getEIPs } from '@/api/network'

const eips = ref([])
const loading = ref(false)

const loadEIPs = async () => {
  loading.value = true
  try {
    const res = await getEIPs()
    eips.value = res.items || res
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadEIPs()
})
</script>

<style scoped>
.eips-page {
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
