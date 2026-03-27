<template>
  <div class="vms-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">虚拟机管理</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            创建虚拟机
          </el-button>
        </div>
      </template>
      
      <el-form :inline="true" :model="queryForm" class="query-form">
        <el-form-item label="云账户">
          <el-input v-model="queryForm.account_id" placeholder="请输入云账户 ID" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadVMs">查询</el-button>
        </el-form-item>
      </el-form>
      
      <el-table :data="vms" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="instance_type" label="实例类型" />
        <el-table-column prop="private_ip" label="内网 IP" />
        <el-table-column prop="public_ip" label="外网 IP" />
        <el-table-column prop="region_id" label="区域" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleAction(row, 'start')">启动</el-button>
            <el-button size="small" @click="handleAction(row, 'stop')">停止</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getVMs, vmAction, deleteVM } from '@/api/compute'
import type { VirtualMachine } from '@/types'

const vms = ref<VirtualMachine[]>([])
const loading = ref(false)
const queryForm = reactive({
  account_id: ''
})

const getStatusName = (status: string) => {
  const map: Record<string, string> = {
    Running: '运行中',
    Stopped: '已停止',
    Starting: '启动中',
    Stopping: '停止中',
    Pending: '创建中'
  }
  return map[status] || status
}

const getStatusType = (status: string) => {
  if (status === 'Running') return 'success'
  if (status === 'Stopped') return 'info'
  if (status === 'Pending' || status === 'Starting') return 'warning'
  return ''
}

const loadVMs = async () => {
  loading.value = true
  try {
    const res = await getVMs(queryForm as any)
    vms.value = res.items || res
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  ElMessage.info('创建虚拟机功能开发中')
}

const handleAction = async (row: VirtualMachine, action: string) => {
  try {
    await ElMessageBox.confirm(`确定要${action === 'start' ? '启动' : '停止'}该虚拟机吗？`, '提示', { type: 'warning' })
    await vmAction(row.id, 1, action as any)
    ElMessage.success('操作成功')
    loadVMs()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

const handleDelete = async (row: VirtualMachine) => {
  try {
    await ElMessageBox.confirm('确定要删除该虚拟机吗？', '提示', { type: 'warning' })
    await deleteVM(row.id, 1)
    ElMessage.success('删除成功')
    loadVMs()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

onMounted(() => {
  loadVMs()
})
</script>

<style scoped>
.vms-page {
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

.query-form {
  margin-bottom: 20px;
}
</style>
