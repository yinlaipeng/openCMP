<template>
  <div class="waf-policies-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">WAF策略列表</span>
        </div>
      </template>

      <el-table :data="wafPolicies" v-loading="loading">
        <el-table-column label="名称" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="tags" label="标签" width="150" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="120" />
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column prop="domain" label="所属域" width="150" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="WAF策略详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedWAFPolicy?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedWAFPolicy?.name }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{ selectedWAFPolicy?.tags }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedWAFPolicy?.status }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ selectedWAFPolicy?.type }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedWAFPolicy?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedWAFPolicy?.account }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ selectedWAFPolicy?.domain }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedWAFPolicy?.region }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedWAFPolicy?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

interface WAFPolicy {
  id: string
  name: string
  tags: string
  status: string
  type: string
  platform: string
  account: string
  domain: string
  region: string
  created_at: string
}

const wafPolicies = ref<WAFPolicy[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedWAFPolicy = ref<WAFPolicy | null>(null)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
    case 'enabled':
      return 'success'
    case 'pending':
    case 'configuring':
      return 'warning'
    case 'disabled':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const loadWAFPolicies = async () => {
  loading.value = true
  try {
    // Mock data
    wafPolicies.value = [
      {
        id: 'waf-1',
        name: 'WAF策略 1',
        tags: 'prod',
        status: 'Active',
        type: '自定义规则',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        domain: 'Default Domain',
        region: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'waf-2',
        name: 'WAF策略 2',
        tags: 'dev',
        status: 'Configuring',
        type: '托管规则',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        domain: 'Default Domain',
        region: 'cn-shanghai',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'waf-3',
        name: 'WAF策略 3',
        tags: 'test',
        status: 'Disabled',
        type: '自定义规则',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        domain: 'Domain A',
        region: 'cn-beijing',
        created_at: '2024-01-03 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    wafPolicies.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: WAFPolicy) => {
  selectedWAFPolicy.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row: WAFPolicy) => {
  ElMessage.info(`编辑WAF策略: ${row.name}`)
}

const handleDelete = async (row: WAFPolicy) => {
  try {
    await ElMessageBox.confirm(`确认删除WAF策略 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    wafPolicies.value = wafPolicies.value.filter(w => w.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadWAFPolicies()
})
</script>

<style scoped>
.waf-policies-page {
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