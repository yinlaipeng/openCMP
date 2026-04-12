<template>
  <div class="blocked-resources-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">屏蔽资源</span>
        </div>
      </template>

      <el-table :data="blockedResources" v-loading="loading">
        <el-table-column prop="policy_name" label="策略名称" width="180" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_type" label="资源类型" width="100" />
        <el-table-column prop="resource_name" label="资源名称" width="180" />
        <el-table-column prop="block_time" label="屏蔽时间" width="160" />
        <el-table-column prop="block_reason" label="屏蔽原因" width="200" />
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="handleUnblock(row)">解除屏蔽</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

interface BlockedResource {
  id: string
  policy_name: string
  status: string
  resource_type: string
  resource_name: string
  block_time: string
  block_reason: string
}

const blockedResources = ref<BlockedResource[]>([])
const loading = ref(false)

const getStatusType = (status: string) => {
  switch (status) {
    case '屏蔽中':
      return 'warning'
    case '已解除':
      return 'success'
    default:
      return 'info'
  }
}

const loadBlockedResources = async () => {
  loading.value = true
  try {
    blockedResources.value = [
      { id: 'block-1', policy_name: '磁盘使用率告警', status: '屏蔽中', resource_type: '虚拟机', resource_name: 'dev-api-02', block_time: '2024-01-15 08:15:00', block_reason: '测试环境，临时屏蔽' },
      { id: 'block-2', policy_name: 'CPU使用率告警', status: '屏蔽中', resource_type: '虚拟机', resource_name: 'test-web-01', block_time: '2024-01-14 10:00:00', block_reason: '开发调试期间' },
      { id: 'block-3', policy_name: '网络流量告警', status: '屏蔽中', resource_type: '网络', resource_name: 'backup-link', block_time: '2024-01-13 15:30:00', block_reason: '数据迁移期间' },
      { id: 'block-4', policy_name: '连接数告警', status: '已解除', resource_type: '负载均衡', resource_name: 'dev-lb-01', block_time: '2024-01-12 09:00:00', block_reason: '压力测试期间' }
    ]
  } catch (e) {
    console.error(e)
    blockedResources.value = []
  } finally {
    loading.value = false
  }
}

const handleUnblock = async (row: BlockedResource) => {
  try {
    await ElMessageBox.confirm(`确认解除资源 ${row.resource_name} 的屏蔽？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })
    row.status = '已解除'
    ElMessage.success('解除屏蔽成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadBlockedResources()
})
</script>

<style scoped>
.blocked-resources-page {
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