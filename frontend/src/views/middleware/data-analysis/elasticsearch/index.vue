<template>
  <div class="elasticsearch-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">Elasticsearch实例列表</span>
        </div>
      </template>

      <el-table :data="esInstances" v-loading="loading">
        <el-table-column label="名称" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="label" label="标签" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="instance_type" label="类型" width="120" />
        <el-table-column prop="config" label="配置" width="150" />
        <el-table-column prop="version" label="版本" width="100" />
        <el-table-column prop="storage" label="存储" width="120" />
        <el-table-column prop="billing_method" label="计费方式" width="120" />
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column prop="project" label="项目" width="120" />
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
    <el-dialog v-model="detailDialogVisible" title="Elasticsearch实例详情" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedES?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedES?.name }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{ selectedES?.label }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedES?.status }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ selectedES?.instance_type }}</el-descriptions-item>
        <el-descriptions-item label="配置">{{ selectedES?.config }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ selectedES?.version }}</el-descriptions-item>
        <el-descriptions-item label="存储">{{ selectedES?.storage }}</el-descriptions-item>
        <el-descriptions-item label="计费方式">{{ selectedES?.billing_method }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedES?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedES?.account }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedES?.project }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedES?.region }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedES?.created_at }}</el-descriptions-item>
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

interface ESInstance {
  id: string
  name: string
  label: string
  status: string
  instance_type: string
  config: string
  version: string
  storage: string
  billing_method: string
  platform: string
  account: string
  project: string
  region: string
  created_at: string
}

const esInstances = ref<ESInstance[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedES = ref<ESInstance | null>(null)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'running':
    case 'active':
      return 'success'
    case 'creating':
    case 'pending':
      return 'warning'
    case 'stopped':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const loadESInstances = async () => {
  loading.value = true
  try {
    // Mock data
    esInstances.value = [
      {
        id: 'es-1',
        name: 'prod-search',
        label: 'production',
        status: 'Running',
        instance_type: '数据节点',
        config: '4核8GB',
        version: '7.10',
        storage: '500GB',
        billing_method: '包年包月',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project A',
        region: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'es-2',
        name: 'dev-search',
        label: 'development',
        status: 'Running',
        instance_type: '单节点',
        config: '2核4GB',
        version: '8.5',
        storage: '200GB',
        billing_method: '按量付费',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project B',
        region: 'cn-shanghai',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'es-3',
        name: 'test-search',
        label: 'test',
        status: 'Creating',
        instance_type: '集群版',
        config: '8核16GB',
        version: '6.8',
        storage: '1TB',
        billing_method: '按量付费',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project A',
        region: 'cn-beijing',
        created_at: '2024-01-03 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    esInstances.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: ESInstance) => {
  selectedES.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row: ESInstance) => {
  ElMessage.info(`编辑Elasticsearch实例: ${row.name}`)
}

const handleDelete = async (row: ESInstance) => {
  try {
    await ElMessageBox.confirm(`确认删除Elasticsearch实例 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    esInstances.value = esInstances.value.filter(e => e.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadESInstances()
})
</script>

<style scoped>
.elasticsearch-page {
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