<template>
  <div class="vpcs-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">VPC 管理</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            创建 VPC
          </el-button>
        </div>
      </template>
      
      <el-table :data="vpcs" v-loading="loading">
        <el-table-column prop="id" label="ID" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="cidr" label="CIDR" />
        <el-table-column prop="status" label="状态" />
        <el-table-column prop="region_id" label="区域" />
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <el-dialog v-model="dialogVisible" title="创建 VPC" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="CIDR">
          <el-input v-model="form.cidr" placeholder="10.0.0.0/16" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getVPCs, createVPC, deleteVPC } from '@/api/network'
import type { VPC } from '@/types'

const vpcs = ref<VPC[]>([])
const loading = ref(false)
const dialogVisible = ref(false)

const form = reactive({
  name: '',
  cidr: '',
  description: ''
})

const loadVPCs = async () => {
  loading.value = true
  try {
    const res = await getVPCs()
    vpcs.value = res.items || res
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  form.name = ''
  form.cidr = ''
  form.description = ''
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    await createVPC({ ...form, account_id: 1 })
    ElMessage.success('创建成功')
    dialogVisible.value = false
    loadVPCs()
  } catch (e) {
    console.error(e)
  }
}

const handleDelete = async (row: VPC) => {
  try {
    await deleteVPC(row.id, 1)
    ElMessage.success('删除成功')
    loadVPCs()
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadVPCs()
})
</script>

<style scoped>
.vpcs-page {
  height: 100%;
}

.page-card {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 18px;
  font-weight: bold;
}
</style>
