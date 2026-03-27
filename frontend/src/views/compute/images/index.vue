<template>
  <div class="images-page">
    <el-card class="page-card">
      <template #header>
        <span class="title">镜像管理</span>
      </template>
      
      <el-table :data="images" v-loading="loading">
        <el-table-column prop="id" label="ID" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="os_name" label="操作系统" />
        <el-table-column prop="status" label="状态" />
        <el-table-column prop="size" label="大小">
          <template #default="{ row }">
            {{ (row.size / 1024 / 1024 / 1024).toFixed(2) }} GB
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getImages } from '@/api/compute'
import type { Image } from '@/types'

const images = ref<Image[]>([])
const loading = ref(false)

const loadImages = async () => {
  loading.value = true
  try {
    const res = await getImages()
    images.value = res.items || res
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadImages()
})
</script>

<style scoped>
.images-page {
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
